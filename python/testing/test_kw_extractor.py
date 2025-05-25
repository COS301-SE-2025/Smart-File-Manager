import unittest
from unittest.mock import patch, MagicMock
import pytest

# < ------ UNIT TESTING ------>
#kw extraction from sentence
@patch("src.kw_extractor.KeywordExtractor")
def test_kw_extract(mock_yake_extract):
    mock_keywords = MagicMock()
    mock_keywords.extract_keywords.return_value = [("kw1", 0.1), ("kw2", 0.2)]
    mock_yake_extract.return_value = mock_keywords
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()
    keywords = extractor.get_kw("dummy sentence")

    assert keywords == [("kw1", 0.1), ("kw2", 0.2)]

#kw split by delimiter default
import os
import tempfile

def test_split_by_delimiter_def():
    test_text = "Hello world. This is a test file.\nIt has multiple sentences. End."
    with tempfile.NamedTemporaryFile('w+', delete=False) as tmpfile:
        tmpfile.write(test_text)
        tmpfile_path = tmpfile.name
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()  

    sentences = list(extractor.split_by_delimiter_def(tmpfile_path, '.'))

    os.unlink(tmpfile_path)
    expected = [
        "Hello world.",
        "This is a test file.",
        "It has multiple sentences.",
        "End."
    ]

    assert sentences == expected

#kw split by delimiter docx
def test_split_by_delimiter_docx():
    page_text = "Hello world. This is a test file.\nIt has multiple sentences. End."

    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()  

    sentences = list(extractor.split_by_delimiter_docx(page_text, '.'))

    expected = [
        "Hello world.",
        "This is a test file.",
        "It has multiple sentences.",
        "End."
    ]

    assert sentences == expected

#kw split by delimiter pdf
def test_split_by_delimiter_pdf():
    page_text = "Hello world. This is a test file.\nIt has multiple sentences. End."

    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()  

    sentences = list(extractor.split_by_delimiter_pdf(page_text, '.'))

    expected = [
        "Hello world.",
        "This is a test file.",
        "It has multiple sentences.",
        "End."
    ]

    assert sentences == expected

#sentence extraction plain text
def test_def_extraction():

    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    extractor.split_by_delimiter_def = MagicMock(return_value=iter([
        "Sentence1",
        "Sentence2",
        "Sentence3",
        "Sentence4"
    ]))

    extractor.get_kw = MagicMock(side_effect=lambda s: s)
    result = extractor.def_extraction("fake_file.txt", ".", max_sentences=3)
    expected_text = "Sentence1.Sentence2.Sentence3."

    extractor.get_kw.assert_called_once_with(expected_text)

    assert result == expected_text

#sentence extraction extraction docx
@patch("src.kw_extractor.docx")  
def test_docx_extraction(mock_docx):
    from src.kw_extractor import KWExtractor

    extractor = KWExtractor()

    mock_doc = MagicMock()
    mock_doc.paragraphs = [MagicMock(text="Paragraph one."), MagicMock(text="Paragraph two.")]

    mock_docx.Document.return_value = mock_doc

    def fake_split(text, delimiter):
        if text == "Paragraph one.":
            yield "Sentence1"
            yield "Sentence2"
        elif text == "Paragraph two.":
            yield "Sentence3"
            yield "Sentence4"

    extractor.split_by_delimiter_docx = fake_split

    extractor.get_kw = MagicMock(side_effect=lambda s: s)

    # Call method with max_sentences = 3
    result = extractor.docx_extraction("fake_file.docx", ".", 3)

    expected_sentence = "Sentence1.Sentence2.Sentence3."

    extractor.get_kw.assert_called_once_with(expected_sentence)
    assert result == expected_sentence

#sentence extraction pdf
@patch("src.kw_extractor.PdfReader")
def test_pdf_extraction(mock_pdf_reader):
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    mock_reader_instance = MagicMock()
    mock_pdf_reader.return_value = mock_reader_instance

    mock_page1 = MagicMock()
    mock_page2 = MagicMock()
    mock_page1.extract_text.return_value = "Page 1 text."
    mock_page2.extract_text.return_value = "Page 2 text."

    mock_reader_instance.pages = [mock_page1, mock_page2]

    def fake_split(text, delimiter):
        if text == "Page 1 text.":
            yield "Sentence1"
            yield "Sentence2"
        elif text == "Page 2 text.":
            yield "Sentence3"
            yield "Sentence4"

    extractor.split_by_delimiter_pdf = fake_split

    extractor.get_kw = MagicMock(side_effect=lambda s: s)

    result = extractor.pdf_extraction("fake_file.pdf", ".", max_sentences=3)

    expected_text = "Sentence1.Sentence2.Sentence3."

    extractor.get_kw.assert_called_once_with(expected_text)

    assert result == expected_text

#open a file
def test_open_file():
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    extractor.pdf_extraction = MagicMock(return_value=["pdf_kw1", "pdf_kw2"])
    extractor.docx_extraction = MagicMock(return_value=["docx_kw1", "docx_kw2"])
    extractor.def_extraction = MagicMock(return_value=["def_kw1", "def_kw2"])

    # Test for pdf type
    result_pdf = extractor.open_file("file.pdf", 5, "application/pdf")
    extractor.pdf_extraction.assert_called_once_with("file.pdf", '.', 5)
    assert result_pdf == [("file.pdf", ["pdf_kw1", "pdf_kw2"])]

    # Test for docx types
    result_docx1 = extractor.open_file("file.doc", 3, "application/msword")
    extractor.docx_extraction.assert_called_with("file.doc", '.', 3)
    assert result_docx1 == [("file.doc", ["docx_kw1", "docx_kw2"])]

    result_docx2 = extractor.open_file("file.docx", 2, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
    extractor.docx_extraction.assert_called_with("file.docx", '.', 2)
    assert result_docx2 == [("file.docx", ["docx_kw1", "docx_kw2"])]

    # Test for plain text
    result_text = extractor.open_file("file.txt", 4, "text/plain")
    extractor.def_extraction.assert_called_with("file.txt", '.', 4)
    assert result_text == [("file.txt", ["def_kw1", "def_kw2"])]

    # Test for unknown type calls def_extraction and handles exceptions
    extractor.def_extraction.reset_mock()
    extractor.def_extraction.side_effect = [["unk_kw1"]]
    result_unknown = extractor.open_file("file.unknown", 1, "unknown/type")
    extractor.def_extraction.assert_called_with("file.unknown", '.', 1)
    assert result_unknown == [("file.unknown", ["unk_kw1"])]

    # Test for unknown type that raises an exception in def_extraction
    extractor.def_extraction.side_effect = Exception("fail")
    result_unknown_fail = extractor.open_file("file.fail", 1, "unknown/type")
    extractor.def_extraction.assert_called_with("file.fail", '.', 1)
    assert result_unknown_fail == []

#kw_extract txt
def test_extract_kw():
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    mock_input = MagicMock()
    mock_file = MagicMock()
    mock_file.original_path = "mock/path/file1.txt"
    mock_file.metadata = [
        MagicMock(key="mime_type", value="text/plain")
    ]
    mock_input.files = [mock_file]

    extractor.open_file = MagicMock(return_value=[("mock/path/file1.txt", ["kw1", "kw2"])])
    extractor.list_to_map = MagicMock(return_value={"mock/path/file1.txt": ["kw1", "kw2"]})

    result = extractor.extract_kw(mock_input)

    # Assertions
    extractor.open_file.assert_called_once_with("mock/path/file1.txt", 10, "text/plain")
    extractor.list_to_map.assert_called_once_with([("mock/path/file1.txt", ["kw1", "kw2"])], 10)
    assert result == {"mock/path/file1.txt": ["kw1", "kw2"]}

#all file types
@pytest.mark.parametrize("mime_type", [
    "application/pdf",
    "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain",
    "unknown/type"
])
def test_extract_kw_all_types(mime_type):
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    mock_input = MagicMock()
    mock_file = MagicMock()
    mock_file.original_path = "mock/path/file.ext"
    mock_file.metadata = [MagicMock(key="mime_type", value=mime_type)]
    mock_input.files = [mock_file]

    # Mock open_file and list_to_map
    mock_keywords = [("mock/path/file.ext", ["kw1", "kw2"])]
    extractor.open_file = MagicMock(return_value=mock_keywords)
    extractor.list_to_map = MagicMock(return_value={"mock/path/file.ext": ["kw1", "kw2"]})

    result = extractor.extract_kw(mock_input)

    extractor.open_file.assert_called_once_with("mock/path/file.ext", 10, mime_type)
    extractor.list_to_map.assert_called_once_with(mock_keywords, 10)
    assert result == {"mock/path/file.ext": ["kw1", "kw2"]}



# < ------ INTEGRATION TESTING ------ >
