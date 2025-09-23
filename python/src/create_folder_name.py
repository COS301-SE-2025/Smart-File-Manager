import os
from pathlib import Path
import re
from typing import Dict
import unicodedata
from collections import Counter, defaultdict

import nltk
from nltk.corpus import wordnet
from nltk import pos_tag
from nltk.stem import WordNetLemmatizer


class FolderNameCreator:
    def __init__(self, model, case_convention: str):
        self.model = model
        self.max_keywords = 200
        self.foldername_length = 2
        self.lemmatizer = WordNetLemmatizer()

        # Weighting
        self.weights = {
            "keywords":0.9, # If there are keywords they should really be different to not be together
            "filename":1.0, # Should only make a small difference compared to keywrods (when they are used)
            "tags":1.5, # Even though they should already be weighted significantly
            "original_parent":0.01,
            # Metadata which can be considered
            "created":0.5
        }

        supported_cases = ["CAMEL", "SNAKE", "PASCAL", "KEBAB"]
        cleaned_case_convention = case_convention.strip().upper()
        self.case_convention = (
            cleaned_case_convention if cleaned_case_convention in supported_cases else "CAMEL"
        )

        self._ensure_nltk_data()

    def _ensure_nltk_data(self):
        try:
            nltk.data.find("tokenizers/punkt")
            nltk.data.find("taggers/averaged_perceptron_tagger")
            nltk.data.find("corpora/wordnet")
            nltk.data.find("corpora/omw-1.4")
        except LookupError:
            nltk.download("punkt", quiet=True)
            nltk.download("averaged_perceptron_tagger_eng", quiet=True)
            nltk.download("wordnet", quiet=True)
            nltk.download("omw-1.4", quiet=True)

    # ---------- CLEANING ----------
    def _normalize_text(self, text: str) -> str:
        """Normalize unicode, strip accents, and lowercase"""
        text = unicodedata.normalize("NFKD", text)
        text = text.encode("ascii", "ignore").decode("utf-8", "ignore")
        return text.lower().strip()

    def _split_words(self, text: str):
        """Split on camelCase, PascalCase, snake_case, kebab-case, dots, numbers"""
        text = re.sub(r"([a-z])([A-Z])", r"\1 \2", text)  # camelCase -> camel Case
        text = re.sub(r"[^a-zA-Z0-9]+", " ", text)        # Discard non alpha numeric 
        return [w for w in text.split() if w]

    # ---------- CORE ----------
    def generateFolderName(self, files):

        # Collect weighted candidates
        scores = defaultdict(float)

        for file in files:
            # Tags
            for tag in file.get("tags", []):
                if isinstance(tag, str):
                    scores[self._clean_name(tag)] += self.weights["tags"]
                elif hasattr(tag, "name"):  # gRPC Tag object
                    scores[self._clean_name(tag.name)] += self.weights["tags"]

            # Keywrods
            for kw, _ in file.get("keywords", []):
                if isinstance(kw, str):
                    scores[self._clean_name(kw)] += self.weights["keywords"]

            # Filename
            if "filename" in file:
                base_name = self.remove_all_extensions(file["filename"])
                clean_fn = self._clean_name(base_name)
                scores[clean_fn] += self.weights["filename"]

        # Nothing found
        if not scores:
            return "Group"

        # Pick top X by score
        sorted_candidates = sorted(scores.items(), key=lambda x: -x[1])
        top_candidates = [name for name, _ in sorted_candidates[:self.foldername_length]]

        # Encode and pick representative (semantic centroid)
        embeddings = self.model.encode(top_candidates)
        centroid = np.mean(embeddings, axis=0, keepdims=True)
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx = sims.argsort()[-self.foldername_length:][::-1]

        folder_keywords = [top_candidates[i] for i in best_idx]

        # Combine into folder name
        return "_".join(folder_keywords)

    def _clean_name(self, name: str) -> str:
        # Remove non-alphanumeric characters and collapse spaces
        name = re.sub(r'[^a-zA-Z0-9\s]', '', name)
        name = re.sub(r'\s+', '_', name).strip('_')
        name = name.lower()

        # Truncate to avoid absurdly long names
        return name[:30] if name else "Group"


    def assignParentScores(self, absolute_path : str, parent_name_scores : Dict[str,float]) -> None:
        path = Path(absolute_path)
        parents = list(path.parents)
        depth = 1

        for parent in parents:
            name = parent.name.lower()
            if name in self.keyword_scores or name in self.filename_scores:
                continue
            if name == "":
                continue
            if name not in parent_name_scores:
                parent_name_scores[name] = 0
            parent_name_scores[name] += self.weights["original_parent"] ** depth
            depth += 1


    def combine_lists(self, keywords : Dict[str, float], filenames : Dict[str,float], parent_names:Dict[str,float]) -> Dict[str,float]:
        scores = defaultdict(float)

        for kw,score in keywords:
            scores[kw] += score

        if not candidates:
            candidates = ["misc"]

        words = self.lemmatize(candidates)
        return self.format_case(words)

    def lemmatize(self, raw_keywords):
        seen = set()
        results = []
        for kw in raw_keywords:
            tokens = self._split_words(kw)
            tagged = pos_tag(tokens)
            lemmas = []
            for word, tag in tagged:
                wn_tag = self.get_wordnet_pos(tag)
                lemma = self.lemmatizer.lemmatize(word, pos=wn_tag)
                lemmas.append(lemma)
            if lemmas:
                formatted = self.format_case(lemmas)
                if formatted not in seen:
                    seen.add(formatted)
                    results.append(formatted)
        return results or ["misc"]

    def get_wordnet_pos(self, treebank_tag):
        if treebank_tag.startswith("J"):
            return wordnet.ADJ
        elif treebank_tag.startswith("V"):
            return wordnet.VERB
        elif treebank_tag.startswith("N"):
            return wordnet.NOUN
        elif treebank_tag.startswith("R"):
            return wordnet.ADV
        return wordnet.NOUN

    def format_case(self, words):
        if not words:
            return "misc"

        if self.case_convention == "CAMEL":
            return words[0].lower() + "".join(w.capitalize() for w in words[1:])
        if self.case_convention == "SNAKE":
            return "_".join(w.lower() for w in words)
        if self.case_convention == "PASCAL":
            return "".join(w.capitalize() for w in words)
        if self.case_convention == "KEBAB":
            return "-".join(w.lower() for w in words)
