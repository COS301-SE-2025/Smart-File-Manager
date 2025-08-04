import os
import re
from pathlib import Path
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np
import nltk

from nltk.corpus import wordnet
from nltk import pos_tag
from nltk.stem import WordNetLemmatizer

#adding these imports could be slow
from sentence_transformers import SentenceTransformer
from typing import List, Dict

from collections import defaultdict


class FolderNameCreator:
    def __init__(self, model : SentenceTransformer):
        self.model = model
        self.max_keywords = 5000 # for folder name creation
        self.lemmatizer = WordNetLemmatizer()
        # Weighting of different vars
        self.weights = {
            "keywords":0.7, # If there are keywords they should really be different to not be together
            "filename":0.002, # Should only make a small difference compared to keywrods (when they are used)
            "tags":1.5, # Even though they should already be weighted significantly
            "original_parent":0.01,
            # Metadata which can be considered
            "created":0.5
        }
        self.foldername_length = 2
        self.filename_scores = {}
        self.parent_name_scores = {}
        self.keyword_scores = {}

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

    # Remove all types of extensions - .png, .tar.gz, etc.
    def remove_all_extensions(self,filename : str) -> str:
        while True:
            filename, ext = os.path.splitext(filename)
            if not ext:
                break
        return filename
    
    def generateFolderName(self, files : List[Dict]) -> str:
            # No files no name
            if not files:
                return "Untitled"           

            # Reset the vals
            self.filename_scores = {}
            self.parent_name_scores = {}
            self.keyword_scores = {}
            # Assign scores with weightings
            for file in files:
                # file name
                fn = self.remove_all_extensions(file["filename"]).lower()
                if not (fn.startswith("~") or fn.endswith(".tmp")):
                    if fn not in self.filename_scores:
                        self.filename_scores[fn] = 0
                    self.filename_scores[fn] += (self.weights["filename"])


                # parent name assigned
                self.assignParentScores(file["absolute_path"], self.parent_name_scores)

                # Assign keywords scores
                for kw,score in file["keywords"]:
                    if kw.lower() not in self.keyword_scores:
                        self.keyword_scores[kw.lower()] = 0
                    self.keyword_scores[kw.lower()] += (score * self.weights["keywords"])

            # Extend by adding metadata as another arg
            combined = self.combine_lists(
                self.getRepresentativeKeywords(self.keyword_scores),
                self.getRepresentativeKeywords(self.filename_scores),
                self.getRepresentativeKeywords(self.parent_name_scores)
            )
            lemmatized = self.lemmatize_with_scores(combined)
            folder_name = "_".join([word for word, _ in lemmatized[:self.foldername_length]])
            folder_name = re.sub(r'[^\w\d_]+', '', folder_name)  # remove any accidental punct
            folder_name = folder_name.strip("_")

            return folder_name

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

        for fn,score in filenames:
            scores[fn] += score 

        for pn,score in parent_names:
            scores[pn] += score

        return sorted(scores.items(), key=lambda x: -x[1])

    def getRepresentativeKeywords(self, scores : Dict[str, float]):
        if not scores:
            return []

        # Top weighted keywrods
        sorted_keywords = sorted(scores.items(), key=lambda x: x[1], reverse=True)
        top_keywords = [kw for kw, _ in sorted_keywords[:self.max_keywords]]

        # Encode file names with sentence transformer
        embeddings = self.model.encode(top_keywords)
        centroid = np.mean(embeddings, axis=0, keepdims=True)

        # Find most representative name (closest to centroid)
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx = sims.argsort()[-self.foldername_length:][::-1]

        # folder_keyword = folder_names[best_idx]
        folder_keyword = [top_keywords[i] for i in best_idx]

        return [(word, scores[word]) for word in folder_keyword]

    
    # Helper function to map NLTK POS to WordNet POS
    def get_wordnet_pos(self, treebank_tag):
        if treebank_tag.startswith('J'):
            return wordnet.ADJ
        elif treebank_tag.startswith('V'):
            return wordnet.VERB
        elif treebank_tag.startswith('N'):
            return wordnet.NOUN
        elif treebank_tag.startswith('R'):
            return wordnet.ADV
        else:
            return wordnet.NOUN  # fallback to noun

    def lemmatize(self, folder_keyword):
        seen = set()
        normalized_keywords = []
        for kw in folder_keyword:
            words = re.split(r'[\s\._\-]+', kw.lower())
            # Tag parts of speech
            tagged = pos_tag(words)
            for word, tag in tagged:
                wn_tag = self.get_wordnet_pos(tag)
                lemma = self.lemmatizer.lemmatize(word, pos=wn_tag)
                if lemma not in seen:
                    seen.add(lemma)
                    normalized_keywords.append(lemma)
        return normalized_keywords

    def lemmatize_with_scores(self, folder_keywords_with_scores):
        seen = {}
        for kw, score in folder_keywords_with_scores:
            words = re.split(r'[\s\._\-]+', kw.lower())
            tagged = pos_tag(words)
            for word, tag in tagged:
                wn_tag = self.get_wordnet_pos(tag)
                lemma = self.lemmatizer.lemmatize(word, pos=wn_tag)
                if lemma not in seen or score > seen[lemma]:
                    seen[lemma] = score
        return sorted(seen.items(), key=lambda x: -x[1])
