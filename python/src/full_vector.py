import datetime
import numpy as np
import pandas as pd
from sklearn.preprocessing import LabelEncoder, MinMaxScaler
from typing import List, Dict, Tuple
from sentence_transformers import SentenceTransformer

class FullVector:
    def __init__(self):
        self.model = SentenceTransformer('all-MiniLM-L6-v2')

    def create_full_vector(self, files: List[Dict]) -> None:
        features = ["size_bytes", "mime_type", "keywords", "created"]

        # Get a dict of feature => list of all those features across every file (i.e. collect features)
        feature_data = self._gather_feature_data(files, features)

        # Normalize numerical data
        normalized_sizes = self._normalize_feature(feature_data["size_bytes"])
        normalized_created = self._normalize_feature([
            self._to_unix_time(ts) for ts in feature_data["created"]
        ])

        # Normalize categorical data
        filetype_encoded_map = self._encode_label(feature_data["mime_type"])

        # Build final vector per file
        for idx, file in enumerate(files):
            # Sentence Transformer
            # Use real keywords for filetypes where they could be extracted
            kw_text = " ".join([kw for kw, _ in file["keywords"]])
            embedding = self.model.encode(kw_text).tolist()

            # Scale certain features
            scaled_type = [val * 0.3 for val in filetype_encoded_map[file["mime_type"]]]

            full_vector = (
                embedding +
                [normalized_created[idx]] +
                [normalized_sizes[idx]] 
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
    def _normalize_feature(self, values: List[float]) -> List[float]:
        if not values:
            return []
        arr = np.array(values).reshape(-1, 1)
        return MinMaxScaler().fit_transform(arr).flatten().tolist()

    # Normalize categorical data via LabelEncoding
    def _encode_label(self, data: List[str]) -> Dict[str, List[float]]:
        encoder = LabelEncoder()
        labels = encoder.fit_transform(data)
        label_map = {
            ext: [label / len(set(labels))]  # normalize the label to [0, 1]
            for ext, label in zip(data, labels)
        }
        return label_map

    def _to_unix_time(self, iso_ts: str) -> float:
        return datetime.datetime.fromisoformat(iso_ts).timestamp()

