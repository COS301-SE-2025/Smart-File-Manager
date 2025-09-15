import os
import re
import unicodedata
from pathlib import Path
from typing import Dict
from collections import Counter, defaultdict

import numpy as np
import nltk
from nltk.corpus import wordnet
from nltk import pos_tag
from nltk.stem import WordNetLemmatizer
from sklearn.metrics.pairwise import cosine_similarity


class FolderNameCreator:
    def __init__(self, model, case_convention: str):
        self.model = model
        self.max_keywords = 200
        self.foldername_length = 2
        self.lemmatizer = WordNetLemmatizer()

        # Weighting
        self.weights = {
            "keywords": 0.9,
            "filename": 0.002,
            "tags": 1.5,
            "original_parent": 0.01,
            "created": 0.5,
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
        text = re.sub(r"([a-z])([A-Z])", r"\1 \2", text)  # camelCase → camel Case
        text = re.sub(r"[^a-zA-Z0-9]+", " ", text)        # keep alphanum only
        return [w for w in text.split() if w]

    # ---------- CORE ----------
    def generateFolderName(self, files):
        tags, keywords = [], []

        for file in files:
            tags.extend(file.get("tags", []))
            keywords.extend([kw for kw, _ in file.get("keywords", [])])

        candidates = []
        if tags:
            most_common_tag, _ = Counter(map(self._normalize_text, tags)).most_common(1)[0]
            candidates = [most_common_tag]
        elif keywords:
            most_common_kw, _ = Counter(map(self._normalize_text, keywords)).most_common(1)[0]
            candidates = [most_common_kw]

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
