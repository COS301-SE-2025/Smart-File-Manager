import os
import datetime
from unittest.mock import patch, MagicMock

import pytest

from src.metadata_scraper import MetaDataScraper

TEST_DIR = os.path.dirname(__file__)
TEST_FILES_DIR = os.path.join(TEST_DIR, "test_files")

def get_test_file(name):
    return os.path.join(TEST_FILES_DIR, name)

# < ------ UNIT TESTING ------>
# Image files
@patch("src.metadata_scraper.Image.open")
def test_image_metadata(mock_image_open):
    mock_image = MagicMock()
    mock_image.__enter__.return_value = mock_image
    mock_image._getexif.return_value = {
        306: "2020:01:01 00:00:00",  # DateTime
        270: "Test Image",  # ImageDescription
    }
    mock_image_open.return_value = mock_image

    from src.metadata_scraper import MetaDataScraper
    metadata = MetaDataScraper.get_image_metadata("fake.jpg")

    assert metadata["DateTime"] == "2020:01:01 00:00:00"
    assert metadata["ImageDescription"] == "Test Image"

# Audio Files
@patch("src.metadata_scraper.mutagen.File")
def test_audio_metadata(mock_mutagen_file):
    mock_audio = {"TPE1": ["Artist Name"], "TIT2": ["Track Title"]}
    mock_mutagen_file.return_value = mock_audio

    from src.metadata_scraper import MetaDataScraper
    metadata = MetaDataScraper.get_audio_metadata("myAudio.mp3")

    assert metadata["TPE1"] == "['Artist Name']"
    assert metadata["TIT2"] == "['Track Title']"

# Pdf files
@patch("src.metadata_scraper.PdfReader")
def test_pdf_metadata(mock_pdf_reader):
    mock_reader = MagicMock()
    mock_reader.metadata = {
        "/Author": "John Doe",
        "/Title": "Sample PDF"
    }
    mock_pdf_reader.return_value = mock_reader

    from src.metadata_scraper import MetaDataScraper
    metadata = MetaDataScraper.get_pdf_metadata("myPdf.pdf")

    assert metadata["/Author"] == "John Doe"
    assert metadata["/Title"] == "Sample PDF"

# Word documents
@patch("src.metadata_scraper.docx.Document")
def test_docx_metadata(mock_docx_doc):
    mock_doc = MagicMock()
    mock_doc.core_properties.author = "Jane Doe"
    mock_doc.core_properties.title = "Test Doc"
    mock_doc.core_properties.subject = "Testing"
    mock_doc.core_properties.created = datetime.datetime(2021, 1, 1)
    mock_doc.core_properties.modified = datetime.datetime(2021, 2, 2)

    mock_docx_doc.return_value = mock_doc

    from src.metadata_scraper import MetaDataScraper
    metadata = MetaDataScraper.get_docx_metadata("myWordDoc.docx")

    assert metadata["author"] == "Jane Doe"
    assert metadata["title"] == "Test Doc"
    assert metadata["created"] == "2021-01-01T00:00:00"

# < ------ INTEGRATION TESTING ------ >
# Real pdf
# def test_real_pdf():
#     test_file = get_test_file("myPdf.pdf")

#     scraper = MetaDataScraper(test_file)

#     metadata = scraper.metadata
#     print(metadata)
    
#     # assert metadata["filename"] == "myPdf.pdf"
#     # assert metadata["file_extension"] == ".pdf"
#     # assert metadata["mime_type"] == "application/pdf"
#     # assert "created" in metadata and isinstance(metadata["created"], str)
#     # assert "modified" in metadata and isinstance(metadata["modified"], str)
#     # assert metadata["size_bytes"] > 0

# # Real img
# def test_real_img():
#     test_file = get_test_file("myImg.jpg")

#     scraper = MetaDataScraper(test_file)
#     print(scraper.metadata)
