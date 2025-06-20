import datetime
import numpy as np
import pandas as pd
from sklearn.preprocessing import LabelEncoder, MinMaxScaler
from typing import List, Dict, Optional, Tuple
from sentence_transformers import SentenceTransformer
from sklearn.preprocessing import MultiLabelBinarizer

class FullVector:
    def __init__(self):
        self.model = SentenceTransformer('all-MiniLM-L6-v2')

    def create_full_vector(self, files: List[Dict]) -> None:

        # List of features we want to extract 
        features = ["size_bytes", "keywords", "created", "tags"]

        # Get a dict of feature => list of all those features across every file (i.e. collect features)
        feature_data = self._gather_feature_data(files, features)

        # Normalize numerical data
        normalized_sizes = self._normalize_feature(feature_data["size_bytes"]) #
        normalized_created = self._normalize_feature([
            self._to_unix_time(ts) for ts in feature_data["created"]
        ])

        # Normalize categorical data
        tag_vectors = self._encode_multi_tags(feature_data["tags"])
        print("\n\n\n")
        print(tag_vectors)


        # Build final vector per file
        for idx, file in enumerate(files):

            # Sentence Transformer for keywords
            kw_text = " ".join([kw for kw, _ in file["keywords"]])
            embedding = self.model.encode(kw_text).tolist()

            weighted_tags = [x * 3 for x in tag_vectors[idx]]

            full_vector = (
                embedding +
                [normalized_created[idx]] +
                [normalized_sizes[idx]] +
                weighted_tags
            )
            file["full_vector"] = full_vector

    # Normalization methods and preprocessing
    def _gather_feature_data(self, files: List[Dict], features: List[str]) -> Dict[str, List]:
        result = {key: [] for key in features}
        for file in files:
            for feat in features:
                result[feat].append(file[feat])
        return result

    # Normalize numerical data via MinMaxScaling
    def _normalize_feature(self, values: List) -> List[float]:
        if not values:
            return []
        arr = np.array(values).reshape(-1, 1)
        return MinMaxScaler().fit_transform(arr).flatten().tolist()

    # Normalize categorical data via LabelEncoding
    def _encode_label(self, data: List[Optional[str]]) -> List[float]:
        
        # Empty tags
        clean_data = [(d if d else "__unknown__") for d in data]
    
        encoder = LabelEncoder()
        labels = encoder.fit_transform(clean_data)
        num_classes = len(set(labels))
    
        # Normalize in range [0, 1]
        return [label / (num_classes - 1) if num_classes > 1 else 0.0 for label in labels]


    def _encode_multi_tags(self, tag_lists: List[List[str]]) -> List[List[float]]:
        mlb = MultiLabelBinarizer()
        one_hot = mlb.fit_transform(tag_lists)
        return one_hot.tolist()




    def _to_unix_time(self, iso_ts: str) -> float:
        return datetime.datetime.fromisoformat(iso_ts).timestamp()

