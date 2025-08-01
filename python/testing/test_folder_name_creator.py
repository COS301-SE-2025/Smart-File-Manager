import pytest
from src.create_folder_name import FolderNameCreator  
from unittest.mock import MagicMock

@pytest.fixture(scope="module")
def dummy_model():
    # Fake model returns index-based unit vectors (useful for cosine similarity)
    class DummyModel:
        def encode(self, texts):
            return [[i + 1 for _ in range(5)] for i in range(len(texts))]

    return DummyModel()

@pytest.fixture
def creator(dummy_model):
    return FolderNameCreator(model=dummy_model)

def test_remove_all_extensions(creator):
    assert creator.remove_all_extensions("file.tar.gz") == "file"
    assert creator.remove_all_extensions("doc.pdf") == "doc"
    assert creator.remove_all_extensions("image.backup.png") == "image"

def test_assign_parent_scores(creator):
    scores = {}
    creator.assignParentScores("/home/user/project/docs/file.txt", scores)
    assert "docs" in scores
    assert scores["docs"] > 0
    assert scores["home"] < scores["docs"]  # deeper folders weighted more

def test_combine_lists_merges_and_sorts(creator):
    k = [("alpha", 0.3), ("beta", 0.1)]
    f = [("beta", 0.2)]
    p = [("alpha", 0.2), ("gamma", 0.4)]
    combined = creator.combine_lists(k, f, p)
    assert combined[0][0] == "alpha"  # alpha has 0.3+0.2 = 0.5
    assert combined[1][0] == "gamma"
    assert combined[2][0] == "beta"

def test_lemmatize_deduplicates(creator):
    result = creator.lemmatize(["running.cases", "run-case", "run"])
    assert set(result) =={"run", "case"}    
    assert len(result) == 2

def test_lemmatize_with_scores_prefers_high_score(creator):
    result = creator.lemmatize_with_scores([("running", 1.0), ("run", 2.0)])
    assert result[0][0] == "run"
    assert result[0][1] == 2.0

def test_generate_folder_name_basic(creator):
    files = [
        {"filename": "MeetingNotes.txt", "absolute_path": "/docs/work", "keywords": [("meeting", 0.9)]},
        {"filename": "Summary.docx", "absolute_path": "/docs/work", "keywords": [("summary", 0.7)]}
    ]
    name = creator.generateFolderName(files)
    assert isinstance(name, str)
    assert name != "Untitled"
    assert "_" not in name or len(name.split("_")) <= creator.foldername_length

def test_generate_folder_name_empty(creator):
    assert creator.generateFolderName([]) == "Untitled"
