import sys
import unittest
from unittest.mock import patch
import pytest
import os
import tempfile

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))



TEST_DIR = os.path.dirname(__file__)
TEST_FILES_DIR = os.path.join(TEST_DIR, "test_files_3")

def get_test_file(name):
    return os.path.join(TEST_FILES_DIR, name)

# < ------ UNIT TESTING ------>
# get path
@patch("src.directory_builder.TEST_DIR", new=os.path.dirname(__file__))
def test_get_path():
    from src.directory_builder import DirectoryCreator
    creator = DirectoryCreator("mock_dir", [])
    path = creator.get_path("my_file.txt")
    expected = os.path.join("mock_dir", "my_file.txt")
    assert path == expected


# createMetaData filters keywords and full_vector
def test_create_metadata_filters_correctly():
    file_info = {
        "filename": "file.txt",
        "absolute_path": "/abs/path/file.txt",
        "filetype": "txt",
        "size_kb": 100,
        "keywords": ["a", "b"],
        "full_vector": [[1, 2, 3]],
        "created": "today"
    }
    from src.directory_builder import DirectoryCreator
    creator = DirectoryCreator("mock_dir", [file_info])
    metadata = creator.createMetaData(file_info)

    from src.message_structure_pb2 import MetadataEntry
    assert all(isinstance(entry, MetadataEntry) for entry in metadata)
    keys = [entry.key for entry in metadata]
    assert "keywords" not in keys
    assert "full_vector" not in keys
    assert "filetype" in keys
    assert "size_kb" in keys
    assert "created" in keys

# createFile returns a File proto with expected values
def test_create_file():
    file_info = {
        "filename": "test.txt",
        "absolute_path": "/abs/path/test.txt",
        "size_kb": 123,
        "filetype": "txt",
        "created": "yesterday"
    }
    from src.directory_builder import DirectoryCreator
    creator = DirectoryCreator("mock_dir", [file_info])
    file_proto = creator.createFile("test.txt", "dir_1")


    from src.message_structure_pb2 import File
    assert isinstance(file_proto, File)
    assert file_proto.name == "test.txt"
    assert file_proto.original_path == "/abs/path/test.txt"
    assert file_proto.new_path.endswith("mock_dir/dir_1/test.txt")
    assert len(file_proto.metadata) > 0
    assert file_proto.tags == []

# buildDirectory creates directory with given name, files, children
def test_build_directory_structure():
    files = [{
        "filename": "doc1.pdf",
        "absolute_path": "/original/path/doc1.pdf",
        "size_kb": 88,
        "filetype": "pdf",
        "created": "2025-01-01"
    }]
    from src.directory_builder import DirectoryCreator
    creator = DirectoryCreator("mock_dir", files)
    child_dir = creator.buildDirectory("child_dir", [], [])
    parent_dir = creator.buildDirectory("parent_dir", files, [child_dir])

    from src.message_structure_pb2 import Directory
    assert isinstance(parent_dir, Directory)
    assert parent_dir.name == "parent_dir"
    assert len(parent_dir.files) == 1
    assert len(parent_dir.directories) == 1
    assert parent_dir.directories[0].name == "child_dir"

from src.directory_builder import DirectoryCreator
from src.message_structure_pb2 import Directory

MOCK_FILES = [
    {
        "filename": "file1.txt",
        "absolute_path": "/original/path/file1.txt",
        "size_kb": 100,
        "filetype": "txt",
        "created": "2024-01-01"
    },
    {
        "filename": "file2.pdf",
        "absolute_path": "/original/path/file2.pdf",
        "size_kb": 250,
        "filetype": "pdf",
        "created": "2024-01-02"
    }
]

# <------ INTEGRATION TESTING ------>
@pytest.fixture
def temp_base_dir():
    with tempfile.TemporaryDirectory() as tempdir:
        yield tempdir  

# to visually see it does actually work
def printDirectoryTree(directory, indent=""):
    print(f"{indent}{directory.name}/")
    for file in directory.files:
        print(f"{indent}  - {file.name}")
    for subdir in directory.directories:
        printDirectoryTree(subdir, indent + "  ")

def test_build_single_directory():
    creator = DirectoryCreator("testdir", MOCK_FILES)

    proto = creator.buildDirectory("testdir", MOCK_FILES, [])

    assert isinstance(proto, Directory)
    assert proto.name == "testdir"
    assert len(proto.files) == 2
    assert proto.directories == []

    # Validate file proto fields
    f1 = proto.files[0]
    assert f1.name == "file1.txt"
    assert f1.original_path == "/original/path/file1.txt"
    assert f1.new_path.endswith("testdir/testdir/file1.txt")  # dirName is used in new_path

    # Metadata should not include keywords/full_vector
    meta_keys = {entry.key for entry in f1.metadata}
    assert "size_kb" in meta_keys
    assert "filetype" in meta_keys
    assert "keywords" not in meta_keys
    assert "full_vector" not in meta_keys

def test_build_nested_directory_structure():
    creator = DirectoryCreator("parentdir", MOCK_FILES)

    # Build nested children
    child = creator.buildDirectory("childdir", [MOCK_FILES[1]], [])
    parent = creator.buildDirectory("parentdir", [MOCK_FILES[0]], [child])


    assert parent.name == "parentdir"
    assert len(parent.files) == 1
    assert len(parent.directories) == 1

    child_proto = parent.directories[0]
    assert child_proto.name == "childdir"
    assert child_proto.files[0].name == "file2.pdf"