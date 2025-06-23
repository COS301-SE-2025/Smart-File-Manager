import unittest
from unittest.mock import patch, MagicMock
import datetime
import numpy as np
from sklearn.preprocessing import LabelEncoder, MinMaxScaler
from typing import List, Dict, Optional
from sentence_transformers import SentenceTransformer
from sklearn.preprocessing import MultiLabelBinarizer
import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))
from src.full_vector import FullVector


# < ------ UNIT TESTING ------>
class TestFullVector(unittest.TestCase):

    @patch('src.full_vector.SentenceTransformer')
    def test_create_full_vector(self, mock_model_class):
        mock_model = MagicMock()
        mock_model.encode.return_value = np.array([1.0, 2.0, 3.0])
        mock_model.get_sentence_embedding_dimension.return_value = 3
        mock_model_class.return_value = mock_model

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

        full_vector_instance = FullVector()
        full_vector_instance.create_full_vector(files)

        for file in files:
            self.assertIn("full_vector", file)
            self.assertIsInstance(file["full_vector"], list)
            # kw_vector (3) + created (1) + size (1) + tag vector (3)
            self.assertEqual(len(file["full_vector"]), 8)


    # Test feature gathering
    def test_gather_feature_data(self):
        files = [
            {"size_bytes": 1000, "keywords": [("key1", 0.5)], "created": "2023-01-01T00:00:00", "tags": ["tag1", "tag2"]},
            {"size_bytes": 2000, "keywords": [], "created": "2023-02-01T00:00:00", "tags": ["tag3"]}
        ]
        features = ["size_bytes", "keywords", "created", "tags"]

        full_vector_instance = FullVector()
        result = full_vector_instance._gather_feature_data(files, features)

        # Assertions
        self.assertEqual(len(result["size_bytes"]), 2)
        self.assertEqual(len(result["keywords"]), 2)
        self.assertEqual(len(result["created"]), 2)
        self.assertEqual(len(result["tags"]), 2)

    def test_normalize_feature(self):
        values = [1000, 2000]
        scaler_mock = MagicMock(spec=MinMaxScaler)
        scaler_mock.fit_transform.return_value = np.array([0.5, 1.0])

        full_vector_instance = FullVector()
        result = full_vector_instance._normalize_feature(values, scaler_mock)

        # Assertions
        self.assertEqual(len(result), 2)
        self.assertAlmostEqual(result[0], 0.5)
        self.assertAlmostEqual(result[1], 1.0)

    def test_encode_multi_tags(self):
        tag_lists = [["tag1", "tag2"], ["tag3"]]

        full_vector_instance = FullVector()
        result = full_vector_instance._encode_multi_tags(tag_lists)

        # Assertions
        self.assertEqual(len(result), 2)
        self.assertEqual(len(result[0]), 3)  
        self.assertEqual(len(result[1]), 3)

    def test_to_unix_time(self):
        iso_ts = "2025-06-22T00:00:00"

        full_vector_instance = FullVector()
        result = full_vector_instance._to_unix_time(iso_ts)

        # Assertions
        self.assertAlmostEqual(result, datetime.datetime(2025, 6, 22, 0, 0).timestamp())