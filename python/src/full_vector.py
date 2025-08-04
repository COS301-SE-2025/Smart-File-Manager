import datetime
import numpy as np
from sklearn.preprocessing import LabelEncoder, MinMaxScaler
from typing import List, Dict, Optional, Tuple
from sentence_transformers import SentenceTransformer
from sklearn.preprocessing import MultiLabelBinarizer

class FullVector:
    def __init__(self, transformer):
        # self.model = SentenceTransformer('all-MiniLM-L6-v2', local_files_only=True)
        self.scaler_size = MinMaxScaler()
        self.scaler_created = MinMaxScaler()
        self.model = transformer

    def create_full_vector(self, files: List[Dict]) -> None:

        # List of features we want to extract 
        features = ["size_bytes", "keywords", "created", "tags"]

        # Get a dict of feature => list of all those features across every file (i.e. collect features)
        feature_data = self._gather_feature_data(files, features)

        # Normalize numerical data
        normalized_sizes = self._normalize_feature(feature_data["size_bytes"], self.scaler_size) 
        normalized_created = self._normalize_feature([
            self._to_unix_time(ts) for ts in feature_data["created"]
        ], self.scaler_created)

        # Normalize categorical data    
        tag_vectors = self._encode_multi_tags(feature_data["tags"])


        # Build final vector per file
        for idx, file in enumerate(files):            

            # Take keywords associated with mean score
            kw_embeddings = [self.model.encode(kw) * (1.0 / (1.0 + score)) for kw, score in file["keywords"]]
            if kw_embeddings:
                kw_vector = np.mean(kw_embeddings, axis=0).tolist()
            else:
                kw_vector = [0] * self.model.get_sentence_embedding_dimension()

            # Tags must carry significant weight
            weighted_tags = [x * 3 for x in tag_vectors[idx]]

            full_vector = (
                kw_vector +
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

    # Normalize numerical data via passed scaler (note must use same scaler for same values)
    def _normalize_feature(self, values: List, scaler) -> List[float]:
        if not values:
            return []
        arr = np.array(values).reshape(-1, 1)
        return scaler.fit_transform(arr).flatten().tolist()

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
    
