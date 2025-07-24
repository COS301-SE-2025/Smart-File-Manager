import os
import re
from pathlib import Path
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np
import nltk
from nltk.stem import WordNetLemmatizer

from collections import defaultdict

nltk.download('wordnet', quiet=True)

class FolderNameCreator:
    def __init__(self, model):
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

    # Remove all types of extensions - .png, .tar.gz, etc.
    def remove_all_extensions(self,filename):
        while True:
            filename, ext = os.path.splitext(filename)
            if not ext:
                break
        return filename
    
    def generateFolderName(self, files) -> str:
            # No files no name
            if not files:
                return "Untitled"           

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
              #  print(parent_name_scores)
                # Assign keywords scores
                for kw,score in file["keywords"]:
                    if kw.lower() not in self.keyword_scores:
                        self.keyword_scores[kw.lower()] = 0
                    self.keyword_scores[kw.lower()] += (score * self.weights["keywords"])

           # print(parent_name_scores)
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

    def assignParentScores(self, absolute_path, parent_name_scores):
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



    def combine_lists(self, keywords, filenames, parent_names):
        scores = defaultdict(float)

        for kw,score in keywords:
            scores[kw] += score

        for fn,score in filenames:
            scores[fn] += score 

        for pn,score in parent_names:
            scores[pn] += score

        return sorted(scores.items(), key=lambda x: -x[1])

    def getRepresentativeKeywords(self, scores):
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

    
    def lemmatize(self, folder_keyword):
        seen = set()
        normalized_keywords = []
        for kw in folder_keyword:
            #words = kw.lower().replace(".", "_").split()
            words = re.split(r'[\s\._\-]+', kw.lower())
            for word in words:
                lemma = self.lemmatizer.lemmatize(word)
                if lemma not in seen:
                    seen.add(lemma)
                    normalized_keywords.append(lemma)
        return normalized_keywords

    def lemmatize_with_scores(self, folder_keywords_with_scores):
        seen = {}
        for kw, score in folder_keywords_with_scores:
            #words = kw.lower().replace(".", "").split()
            words = re.split(r'[\s\._\-]+', kw.lower())
            for word in words:
                lemma = self.lemmatizer.lemmatize(word)
                if lemma not in seen or score > seen[lemma]:
                    seen[lemma] = score
        return sorted(seen.items(), key=lambda x: -x[1])

"""
    def generateWithKeywords(self, keyword_scores):
        # Top weighted keywrods
        sorted_keywords = sorted(keyword_scores.items(), key=lambda x: x[1], reverse=True)
        top_keywords = [kw for kw, _ in sorted_keywords[:self.max_keywords]]

        # Encode and find a centroid
        embeddings = self.model.encode(top_keywords)
        centroid = np.mean(embeddings, axis=0,keepdims=True)

        # Best keyword
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx =sims.argsort()[-self.foldername_length:][::-1] 

        folder_keyword = [top_keywords[i] for i in best_idx]

        return folder_keyword
    
    def generateWithoutKeywords(self, folder_names):

        # Encode file names with sentence transformer
        embeddings = self.model.encode(folder_names)
        centroid = np.mean(embeddings, axis=0, keepdims=True)

        # Find most representative name (closest to centroid)
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx = sims.argsort()[-self.foldername_length:][::-1]

        # folder_keyword = folder_names[best_idx]
        folder_keyword = [os.path.splitext(folder_names[i])[0] for i in best_idx]

        return folder_keyword
"""