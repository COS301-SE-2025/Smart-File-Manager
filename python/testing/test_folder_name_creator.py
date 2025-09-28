import pytest
from src.create_folder_name import FolderNameCreator  

@pytest.fixture(scope="module")
def dummy_model():
    # Fake model returns index-based unit vectors (useful for cosine similarity)
    class DummyModel:
        def encode(self, texts):
            return [[i + 1 for _ in range(5)] for i in range(len(texts))]

    return DummyModel()

# <------ Unit TESTING ----->

@pytest.fixture
def creator(dummy_model):
    return FolderNameCreator(model=dummy_model, case_convention="CAMEL")


def test_lemmatize_deduplicates(creator):
    result = creator.lemmatize(["running.cases", "run-case", "run"])
    assert result == [['run', 'case'], ['run', 'case'], ['run']]
    assert len(result) == 3


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
    assert creator.generateFolderName([]) == "misc"

def test_format_case(dummy_model):
    words = ["test", "case"]

    assert FolderNameCreator(dummy_model, "CAMEL").format_case(words) == "testCase"
    assert FolderNameCreator(dummy_model, "SNAKE").format_case(words) == "test_case"
    assert FolderNameCreator(dummy_model, "PASCAL").format_case(words) == "TestCase"
    assert FolderNameCreator(dummy_model, "KEBAB").format_case(words) == "test-case"

