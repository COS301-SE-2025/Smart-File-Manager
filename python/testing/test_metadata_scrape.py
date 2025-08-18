import os



TEST_DIR = os.path.dirname(__file__)
TEST_FILES_DIR = os.path.join(TEST_DIR, "test_files")

def get_test_file(name):
    return os.path.join(TEST_FILES_DIR, name)

# < ------ UNIT TESTING ------>
# Image file
"""
@patch("src.metadata_scraper.os.path.exists", return_value=True)
@patch("src.metadata_scraper.Image.open")
def test_image_metadata(mock_image_open, mock_path_exist):
    mock_image = MagicMock()
    mock_image.__enter__.return_value = mock_image
    mock_image._getexif.return_value = {
        306: "2020:01:01 00:00:00",  # DateTime
        270: "Test Image",  # ImageDescription
    }
    mock_image_open.return_value = mock_image

    from src.metadata_scraper import MetaDataScraper
    scraper = MetaDataScraper()
    scraper.set_file("fake.jpg")
    metadata = scraper.get_image_metadata()

    assert metadata["DateTime"] == "2020:01:01 00:00:00"
    assert metadata["ImageDescription"] == "Test Image"

# Audio and video
@patch("src.metadata_scraper.os.path.exists", return_value=True)
@patch("src.metadata_scraper.mutagen.File")
def test_audio_video_metadata(mock_mutagen_file, mock_path_exist):
    from src.metadata_scraper import MetaDataScraper

    mock_audio = {"TPE1": ["Artist Name"], "TIT2": ["Track Title"]}
    mock_video = {"©nam": ["Test Video"], "©ART": ["Video Creator"]}

    mock_mutagen_file.side_effect = [mock_audio, mock_video] # audio then video for mock

    # Test audio 
    audio_scraper = MetaDataScraper()
    audio_scraper.set_file("myAudio.mp3")
    audio_metadata = audio_scraper.get_audio_video_metadata()
    assert audio_metadata["TPE1"] == "['Artist Name']"
    assert audio_metadata["TIT2"] == "['Track Title']"

    # Test video 
    video_scraper = MetaDataScraper()
    video_scraper.set_file("myVideo.mp4")
    video_metadata = video_scraper.get_audio_video_metadata()
    assert video_metadata["©nam"] == "['Test Video']"
    assert video_metadata["©ART"] == "['Video Creator']"


# Pdf files
@patch("src.metadata_scraper.os.path.exists", return_value=True)
@patch("src.metadata_scraper.PdfReader")
def test_pdf_metadata(mock_pdf_reader, mock_path_exist):
    mock_reader = MagicMock()
    mock_reader.metadata = {
        "/Author": "John Doe",
        "/Title": "Sample PDF"
    }
    mock_pdf_reader.return_value = mock_reader

    from src.metadata_scraper import MetaDataScraper
    scraper = MetaDataScraper()
    scraper.set_file("myPdf.pdf")
    metadata = scraper.get_pdf_metadata()

    assert metadata["/Author"] == "John Doe"
    assert metadata["/Title"] == "Sample PDF"

# Word documents
@patch("src.metadata_scraper.os.path.exists", return_value=True)
@patch("src.metadata_scraper.docx.Document")
def test_docx_metadata(mock_docx_doc, mock_path_exist):
    mock_doc = MagicMock()
    mock_doc.core_properties.author = "Jane Doe"
    mock_doc.core_properties.title = "Test Doc"
    mock_doc.core_properties.subject = "Testing"
    mock_doc.core_properties.created = datetime.datetime(2021, 1, 1)
    mock_doc.core_properties.modified = datetime.datetime(2021, 2, 2)

    mock_docx_doc.return_value = mock_doc

    from src.metadata_scraper import MetaDataScraper
    scraper = MetaDataScraper()
    scraper.set_file("myWordDoc.docx")
    metadata = scraper.get_docx_metadata()

    assert metadata["author"] == "Jane Doe"
    assert metadata["title"] == "Test Doc"
    assert metadata["created"] == "2021-01-01T00:00:00"

# Test non existing file
@patch("src.metadata_scraper.os.path.exists", return_value=False)
@patch("src.metadata_scraper.docx.Document")
def test_nonexisting_file(mock_docx_doc, mock_path_exist):
    mock_doc = MagicMock()
    mock_doc.core_properties.modified = datetime.datetime(2021, 2, 2)
    mock_docx_doc.return_value = mock_doc

    from src.metadata_scraper import MetaDataScraper
    scraper = MetaDataScraper()

    with pytest.raises(ValueError, match="New path does not exist on this filesystem"):
        scraper.set_file("myWordDoc.docx")


# < ------ INTEGRATION TESTING ------ >
# Real pdf
def test_real_pdf():
    test_file = get_test_file("myPdf.pdf")

    scraper = MetaDataScraper()
    scraper.set_file(test_file)
    scraper.get_metadata()
    metadata = scraper.metadata
    
    assert metadata["filename"] == "myPdf.pdf"
    assert metadata["file_extension"] == ".pdf"
    assert metadata["mime_type"] == "application/pdf"
    assert "created" in metadata and isinstance(metadata["created"], str)
    assert "modified" in metadata and isinstance(metadata["modified"], str)
    assert metadata["size_bytes"] > 0
    assert len(metadata) >= 7

# Real img
def test_real_img():
    test_file = get_test_file("myImg.jpg")

    scraper = MetaDataScraper()
    scraper.set_file(test_file)
    scraper.get_metadata()

    metadata = scraper.metadata

    assert metadata["filename"] == "myImg.jpg"
    assert metadata["file_extension"] == ".jpg"
    assert metadata["mime_type"] == "image/jpeg"
    assert "created" in metadata and isinstance(metadata["created"], str)
    assert "modified" in metadata and isinstance(metadata["modified"], str)
    assert metadata["size_bytes"] > 0
    assert len(metadata) >= 7

# Real word doc
def test_real_word():
    test_file = get_test_file("myWordDoc.docx")

    scraper = MetaDataScraper()
    scraper.set_file(test_file)
    scraper.get_metadata()

    metadata = scraper.metadata

    assert metadata["filename"] == "myWordDoc.docx"
    assert metadata["file_extension"] == ".docx"
    assert metadata["mime_type"] == "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    assert "created" in metadata and isinstance(metadata["created"], str)
    assert "modified" in metadata and isinstance(metadata["modified"], str)
    assert metadata["size_bytes"] > 0
    #    assert metadata["author"] == "Philipp du Plessis"
    #    assert metadata["title"] == ""
    #    assert metadata["subject"] == ""
    assert len(metadata) >= 7

# Real audio file
def test_real_audio():
    test_file = get_test_file("myAudio.m4a")

    scraper = MetaDataScraper()
    scraper.set_file(test_file)
    scraper.get_metadata()

    metadata = scraper.metadata
    
    assert metadata["filename"] == "myAudio.m4a"
    assert metadata["file_extension"] == ".m4a"
    assert "audio" in metadata["mime_type"]  
    assert "created" in metadata and isinstance(metadata["created"], str)
    assert "modified" in metadata and isinstance(metadata["modified"], str)
    assert metadata["size_bytes"] > 0
    assert len(metadata) >= 7

# Real video file
def test_real_video():
    test_file = get_test_file("myVideo.webm")

    scraper = MetaDataScraper()
    scraper.set_file(test_file)
    scraper.get_metadata()

    metadata = scraper.metadata
     
    assert metadata["filename"] == "myVideo.webm"
    assert metadata["file_extension"] == ".webm"
    assert metadata["mime_type"] == "video/webm"
    assert "created" in metadata and isinstance(metadata["created"], str)
    assert "modified" in metadata and isinstance(metadata["modified"], str)
    assert metadata["size_bytes"] > 0
"""