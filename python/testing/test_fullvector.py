from unittest.mock import MagicMock
from src.full_vector import FullVector
import numpy as np
import datetime
from sklearn.preprocessing import MinMaxScaler
import unittest

class TestFullVector(unittest.TestCase):

    def setUp(self):
        self.mock_transformer = MagicMock()
        self.mock_transformer.get_sentence_embedding_dimension.return_value = 3
        self.mock_transformer.encode.side_effect = lambda *args, **kwargs: [
            np.array([1.0, 2.0, 3.0]) for _ in args[0]
        ]

    def test_create_full_vector(self):
        files = [
            {
                "size_bytes": 1000,
                "keywords": [("keyword1", 0.5), ("keyword2", 1.0)],
                "created": "2023-01-01T00:00:00",
                "tags": ["tag1", "tag2"]
            },
            {
                "size_bytes": 2000,
                "keywords": [],
                "created": "2023-02-01T00:00:00",
                "tags": ["tag3"]
            }
        ]

        fv = FullVector(self.mock_transformer)
        fv.create_full_vector(files)

        for file in files:
            self.assertIn("full_vector", file)
            self.assertIsInstance(file["full_vector"], list)
            self.assertEqual(len(file["full_vector"]), 8)  # 3 (kw) + 1 (created) + 1 (size) + 3 (tags)

    def test_gather_feature_data(self):
        files = [
            {"size_bytes": 1000, "keywords": [("key1", 0.5)], "created": "2023-01-01T00:00:00", "tags": ["tag1", "tag2"]},
            {"size_bytes": 2000, "keywords": [], "created": "2023-02-01T00:00:00", "tags": ["tag3"]}
        ]
        features = ["size_bytes", "keywords", "created", "tags"]

        fv = FullVector(self.mock_transformer)
        result = fv._gather_feature_data(files, features)

        self.assertEqual(len(result["size_bytes"]), 2)
        self.assertEqual(len(result["keywords"]), 2)
        self.assertEqual(len(result["created"]), 2)
        self.assertEqual(len(result["tags"]), 2)

    def test_normalize_feature(self):
        values = [1000, 2000]
        scaler_mock = MagicMock(spec=MinMaxScaler)
        scaler_mock.fit_transform.return_value = np.array([0.5, 1.0])

        fv = FullVector(self.mock_transformer)
        result = fv._normalize_feature(values, scaler_mock)

        self.assertEqual(len(result), 2)
        self.assertAlmostEqual(result[0], 0.5)
        self.assertAlmostEqual(result[1], 1.0)

    def test_encode_multi_tags(self):
        tag_lists = [["tag1", "tag2"], ["tag3"]]

        fv = FullVector(self.mock_transformer)
        result = fv._encode_multi_tags(tag_lists)

        self.assertEqual(len(result), 2)
        self.assertEqual(len(result[0]), 3)
        self.assertEqual(len(result[1]), 3)

    def test_to_unix_time(self):
        iso_ts = "2025-06-22T00:00:00"

        fv = FullVector(self.mock_transformer)
        result = fv._to_unix_time(iso_ts)

        self.assertAlmostEqual(result, datetime.datetime(2025, 6, 22, 0, 0).timestamp())
