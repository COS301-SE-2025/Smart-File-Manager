# IF THIS IS HERE IT ISNT ADDED TO README
# pip install nltk

import os
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np
import nltk
from nltk.stem import WordNetLemmatizer
nltk.download('wordnet', quiet=True)

class FolderNameCreator:
    def __init__(self, model):
        self.model = model
        self.max_keywords = 5000 # for folder name creation
        self.lemmatizer = WordNetLemmatizer()


    def generateFolderName(self, files) -> str:
            if not files:
                return "Untitled"
            
            keyword_scores = {}
 #           print("-----------------------------")
            for file in files:
#                print("Considering file: ", file["filename"])
                for kw,score in file["keywords"]:
                    if kw not in keyword_scores:
                        keyword_scores[kw] = 0
                    keyword_scores[kw] += 1.0 / (1.0 + score)


            if not keyword_scores: # probably an image. Can use a suffix if there are multiple of these folders. Or we can use their names in a sentence transformer
                folder_names = [os.path.basename(file["filename"]) for file in files]

                folder_keyword = self.generateWithoutKeywords(folder_names)

            else:
                folder_keyword = self.generateWithKeywords(keyword_scores)

            # Normalize: lowercase + lemmatize + deduplicate
            normalized_keywords = self.lemmatize(folder_keyword)
            # Join top 3 after normalization
            folder_name = "_".join(normalized_keywords[:3])
           # folder_name = "_".join(kw.replace(" ", "_").replace(".","") for kw in folder_keyword)
            
#            print("Folder name chosen: ", folder_name)
#            print("-----------------------------")
            return folder_name


    def generateWithKeywords(self, keyword_scores):
        # Top weighted keywrods
        sorted_keywords = sorted(keyword_scores.items(), key=lambda x: x[1], reverse=True)
        top_keywords = [kw for kw, _ in sorted_keywords[:self.max_keywords]]

        # Encode and find a centroid
        embeddings = self.model.encode(top_keywords)
        centroid = np.mean(embeddings, axis=0,keepdims=True)

        # Best keyword
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx =sims.argsort()[-3:][::-1] 

        folder_keyword = [top_keywords[i] for i in best_idx]

        return folder_keyword
    
    def generateWithoutKeywords(self, folder_names):

        # Encode file names with sentence transformer
        embeddings = self.model.encode(folder_names)
        centroid = np.mean(embeddings, axis=0, keepdims=True)

        # Find most representative name (closest to centroid)
        sims = cosine_similarity(centroid, embeddings).flatten()
        best_idx = sims.argsort()[-2:][::-1]

        # folder_keyword = folder_names[best_idx]
        folder_keyword = [os.path.splitext(folder_names[i])[0] for i in best_idx]

        return folder_keyword
    
    def lemmatize(self, folder_keyword):
        seen = set()
        normalized_keywords = []
        for kw in folder_keyword:
            words = kw.lower().replace(".", "").split()
            for word in words:
                lemma = self.lemmatizer.lemmatize(word)
                if lemma not in seen:
                    seen.add(lemma)
                    normalized_keywords.append(lemma)
        return normalized_keywords

