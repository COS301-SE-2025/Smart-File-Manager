import re
import nltk
nltk.download('punkt')       # For word tokenization
nltk.download('wordnet')     # For lemmatization
nltk.download('stopwords')   #
from nltk.stem import PorterStemmer
from difflib import get_close_matches
import numpy as np
from nltk.stem import WordNetLemmatizer
from nltk.corpus import stopwords

class KWCluster:
    def __init__(self):
        self.ps = PorterStemmer()
        self.lemmatizer = WordNetLemmatizer()
        self.stops = set(stopwords.words("english"))
        pass


    # only add exact matches for now
    def createCluster(self, keywords, vocab):
        kw_score_map = {self.normalize(kw): score for kw, score in keywords}
        cluster_vector = []
        for kw in vocab:
            match = get_close_matches(self.normalize(kw), kw_score_map.keys(), n=2, cutoff=0.8)
            if match:
                cluster_vector.append(kw_score_map[match[0]])  # score if exists else 0
            else:
                cluster_vector.append(0)


        return cluster_vector
    

    # def normalize(self, text):
    #     text = text.lower()
    #     text = re.sub(r'\W+', ' ', text).strip()
    #     return self.ps.stem(text)
    



    def normalize(self,kw: str) -> str:
        kw = kw.lower()
        kw = re.sub(r'[^\w\s]', '', kw)        # remove punctuation
        kw = re.sub(r'\s+', ' ', kw).strip()   # normalize whitespace
        words = kw.split()
        words = [self.lemmatizer.lemmatize(w) for w in words if w not in self.stops]
        return ' '.join(words)
    