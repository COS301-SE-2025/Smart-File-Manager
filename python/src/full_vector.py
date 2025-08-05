import datetime
import numpy as np
from sklearn.preprocessing import LabelEncoder, MinMaxScaler, MultiLabelBinarizer
from typing import List, Dict, Optional


class FullVector:
    def __init__(self, transformer):
        self.scaler_size = MinMaxScaler()
        self.scaler_created = MinMaxScaler()
        self.model = transformer

    def create_full_vector(self, files: List[Dict]) -> None:
        features = ["size_bytes", "keywords", "created", "tags"]
        feature_data = self._gather_feature_data(files, features)

        # Normalize numerical features
        normalized_sizes = self._normalize_feature(feature_data["size_bytes"], self.scaler_size)
        normalized_created = self._normalize_feature(
            [self._to_unix_time(ts) for ts in feature_data["created"]],
            self.scaler_created
        )

        # One-hot encode tags
        tag_vectors = self._encode_multi_tags(feature_data["tags"])

        all_kw = []
        file_kw_lengths = []

        for kws in feature_data["keywords"]:
            file_kw_lengths.append(len(kws))
            all_kw.extend([kw for kw, _ in kws])

        # Batch encode kws
        if all_kw:
            all_kw_vectors = self.model.encode(all_kw, show_progress_bar=False, batch_size=64)
        else:
            all_kw_vectors = []

        # Reconstruct keyword vectors per file
        kw_vectors_per_file = []
        idx = 0
        for length, kws in zip(file_kw_lengths, feature_data["keywords"]):
            if length == 0:
                kw_vectors_per_file.append(None)
                continue

            vectors = []
            for i in range(length):
                kw_vec = all_kw_vectors[idx]
                _, score = kws[i]
                vectors.append(kw_vec * (1.0 / (1.0 + score)))
                idx += 1

            kw_vectors_per_file.append(np.mean(vectors, axis=0))

        # Construct final vectors
        embedding_dim = self.model.get_sentence_embedding_dimension()
        for i, file in enumerate(files):
            kw_vec = (
                kw_vectors_per_file[i].tolist()
                if kw_vectors_per_file[i] is not None
                else [0.0] * embedding_dim
            )

            weighted_tags = np.array(tag_vectors[i]) * 3
            full_vector = (
                kw_vec +
                [normalized_created[i]] +
                [normalized_sizes[i]] +
                weighted_tags.tolist()
            )

            file["full_vector"] = full_vector

    def _gather_feature_data(self, files: List[Dict], features: List[str]) -> Dict[str, List]:
        return {feat: [file[feat] for file in files] for feat in features}

    def _normalize_feature(self, values: List, scaler) -> List[float]:
        if not values:
            return []
        arr = np.array(values).reshape(-1, 1)
        return scaler.fit_transform(arr).flatten().tolist()

    def _encode_label(self, data: List[Optional[str]]) -> List[float]:
        clean_data = [(d if d else "__unknown__") for d in data]
        encoder = LabelEncoder()
        labels = encoder.fit_transform(clean_data)
        num_classes = len(set(labels))
        return [label / (num_classes - 1) if num_classes > 1 else 0.0 for label in labels]

    def _encode_multi_tags(self, tag_lists: List[List[str]]) -> List[List[float]]:
        mlb = MultiLabelBinarizer()
        return mlb.fit_transform(tag_lists).tolist()

    def _to_unix_time(self, iso_ts: str) -> float:
        return datetime.datetime.fromisoformat(iso_ts).timestamp()
