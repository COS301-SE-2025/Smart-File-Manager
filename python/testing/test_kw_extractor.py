import sys
from unittest.mock import patch, MagicMock
from types import SimpleNamespace
import pytest
import os
import tempfile

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))



TEST_DIR = os.path.dirname(__file__)
TEST_FILES_DIR = os.path.join(TEST_DIR, "test_files")

def get_test_file(name):
    return os.path.join(TEST_FILES_DIR, name)

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
    result = extractor.def_extraction("fake_file.txt", ".")
    expected_text = "Sentence1.Sentence2.Sentence3.Sentence4."

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
    result = extractor.docx_extraction("fake_file.docx", ".")

    expected_sentence = "Sentence1.Sentence2.Sentence3.Sentence4."

    extractor.get_kw.assert_called_once_with(expected_sentence)
    assert result == expected_sentence

#sentence extraction pdf
@patch("src.kw_extractor.fitz.open")
def test_pdf_extraction(mock_fitz_open):
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    # Mock document and pages
    mock_doc = MagicMock()
    mock_page1 = MagicMock()
    mock_page2 = MagicMock()

    mock_page1.get_text.return_value = "Page 1 text."
    mock_page2.get_text.return_value = "Page 2 text."

    mock_doc.__iter__.return_value = iter([mock_page1, mock_page2])
    mock_fitz_open.return_value = mock_doc

    # Mock the sentence splitter per page
    extractor.split_by_delimiter_pdf = MagicMock(side_effect=lambda text, delimiter: {
        "Page 1 text.": ["Sentence1", "Sentence2"],
        "Page 2 text.": ["Sentence3", "Sentence4"]
    }[text])

    # Mock the keyword extractor
    extractor.get_kw = MagicMock(side_effect=lambda s: s)

    result = extractor.pdf_extraction("fake_file.pdf", ".")

    expected_text = "Sentence1.Sentence2.Sentence3.Sentence4."

    extractor.get_kw.assert_called_once_with(expected_text)
    assert result == expected_text

# #open a file
def test_open_file():
    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()

    extractor.pdf_extraction = MagicMock(return_value=["pdf_kw1", "pdf_kw2"])
    extractor.docx_extraction = MagicMock(return_value=["docx_kw1", "docx_kw2"])
    extractor.def_extraction = MagicMock(return_value=["def_kw1", "def_kw2"])

    extractor.mime_handlers = {
    "application/pdf": extractor.pdf_extraction,
    "application/msword": extractor.docx_extraction,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": extractor.docx_extraction,
    "text/plain": extractor.def_extraction
    }

    # Test for pdf type
    result_pdf = extractor.open_file("file.pdf", "application/pdf")
    assert result_pdf == [("file.pdf", ["pdf_kw1", "pdf_kw2"])]

    # Test for docx types
    result_docx1 = extractor.open_file("file.doc", "application/msword")
    extractor.docx_extraction.assert_called_with("file.doc", '.',1)
    assert result_docx1 == [("file.doc", ["docx_kw1", "docx_kw2"])]

    result_docx2 = extractor.open_file("file.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
    extractor.docx_extraction.assert_called_with("file.docx", '.',1)
    assert result_docx2 == [("file.docx", ["docx_kw1", "docx_kw2"])]

    # Test for plain text
    result_text = extractor.open_file("file.txt", "text/plain")
    extractor.def_extraction.assert_called_with("file.txt", '.',1)
    assert result_text == [("file.txt", ["def_kw1", "def_kw2"])]

    # Test for unknown type calls def_extraction and handles exceptions
    extractor.def_extraction.reset_mock()
    extractor.def_extraction.side_effect = [["unk_kw1"]]
    result_unknown = extractor.open_file("file.unknown", "unknown/type")
    extractor.def_extraction.assert_called_with("file.unknown", '.',1)
    assert result_unknown == [("file.unknown", ["unk_kw1"])]

    # Test for unknown type that raises an exception in def_extraction
    extractor.def_extraction.side_effect = Exception("fail")
    result_unknown_fail = extractor.open_file("file.fail", "unknown/type",1)
    extractor.def_extraction.assert_called_with("file.fail", '.',1)
    assert result_unknown_fail == []

#kw_extract txt

def test_extract_kw():
    from src.kw_extractor import KWExtractor

    extractor = KWExtractor()

    # Create mock input object with required properties
    mock_input = MagicMock()
    mock_input.original_path = "mock/path/file1.txt"
    mock_input.metadata = [
        SimpleNamespace(key="mime_type", value="text/plain")
    ]

    # Mock the dependent methods
    extractor.open_file = MagicMock(return_value=[
        ("mock/path/file1.txt", [("kw1", 0.9), ("kw2", 0.8)])
    ])

    # Call the method under test
    result = extractor.extract_kw(mock_input)

    # Assertions
    extractor.open_file.assert_called_once_with("mock/path/file1.txt", "text/plain", 1)
    assert result == [("kw1", 0.9), ("kw2", 0.8)]


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
    mock_input.original_path = "mock/path/file.ext"
    mock_input.metadata = [
        SimpleNamespace(key="mime_type", value=mime_type)
    ]

    # Mock open_file and list_to_map
    mock_keywords = [("mock/path/file.ext", [("kw1", 0.9), ("kw2", 0.8)])]
    extractor.open_file = MagicMock(return_value=mock_keywords)

    result = extractor.extract_kw(mock_input)

    extractor.open_file.assert_called_once_with("mock/path/file.ext", mime_type,1)
    assert result == [("kw1", 0.9), ("kw2", 0.8)]



# < ------ INTEGRATION TESTING ------ >
#Tests the files as if they are passed through a single directory
def test_real_data_all_files():
    from src.message_structure_pb2 import Directory, File, Tag, MetadataEntry, DirectoryRequest
    from src.kw_extractor import KWExtractor
    tag1 = Tag(name="ImFixed")
    meta1 = MetadataEntry(key="author", value="johnny")
    meta4 = MetadataEntry(key="mime_type", value="text/plain")
    meta2 = MetadataEntry(key="mime_type", value="application/pdf")
    meta3 = MetadataEntry(key="mime_type", value="application/msword")

    file1_path = get_test_file("myPdf.pdf")
    file2_path = get_test_file("testFile.txt")
    file3_path = get_test_file("myWordDoc.docx")

    file1 = File(
        name="gopdoc.pdf",
        original_path=file1_path,
        new_path="/usr/trash/gopdoc.pdf",
        tags=[tag1],
        metadata=[meta1, meta2]
    )
    file2 = File(
        name="gopdoc2.pdf",
        original_path=file2_path,
        new_path="/usr/trash/gopdoc2.pdf",
        tags=[tag1],
        metadata=[meta1, meta4]
    )
    file3 = File(
        name="gopdoc3.pdf",
        original_path=file3_path,
        new_path="/usr/trash/gopdoc3.pdf",
        tags=[tag1],
        metadata=[meta1, meta3]
    )

    kw_extractor = KWExtractor()

    expected_pdf = {"project", "management", "proposal", "folders", "manager", "capstone", "southern", "cross"}
    expected_txt = {"assignment", "debugged", "midterms", "email", "alarm", "laptops", "evacuating", "java"}
    expected_docx = {"docker", "class", "diagram", "uml", "rest", "architecture", "deployment", "frontend"}

    def flatten_keywords(kws):
        flattened = set()
        for item in kws:
            if isinstance(item, tuple) and len(item) == 2:
                kw, _ = item
            else:
                kw = item
            flattened.add(kw.lower().strip())
        return flattened

    pdf_result = flatten_keywords(kw_extractor.extract_kw(file1))
    txt_result = flatten_keywords(kw_extractor.extract_kw(file2))
    docx_result = flatten_keywords(kw_extractor.extract_kw(file3))

    def normalize(kw):
        return kw.lower().replace("-", "").strip(".")

    def check_majority(expected_keywords, result_keywords, threshold=0):
        normalized_results = {normalize(kw) for kw in result_keywords}
        normalized_expected = {normalize(kw) for kw in expected_keywords}
        matches = normalized_expected & normalized_results
        ratio = len(matches) / len(normalized_expected)
        assert ratio >= threshold, f"Too few matches: {matches} (Ratio: {ratio:.2f})"

    check_majority(expected_txt, txt_result)
    check_majority(expected_docx, docx_result)
    check_majority(expected_pdf, pdf_result)


@pytest.mark.parametrize("filename, mime_type, expected_keywords", [
    ("myPdf.pdf", "application/pdf", {"project", "management", "proposal", "folders", "manager", "capstone", "southern", "cross"}),
    ("testFile.txt", "text/plain", {"assignment", "debugged", "midterms", "email", "alarm", "laptops", "evacuating", "java"}),
    ("myWordDoc.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", {"docker", "class", "diagram", "uml", "rest", "architecture", "deployment", "frontend"}),
])
def test_extract_kw_per_file_type(filename, mime_type, expected_keywords):
    from src.message_structure_pb2 import File, Tag, MetadataEntry, Directory, DirectoryRequest
    from src.kw_extractor import KWExtractor

    path = get_test_file(filename)

    tag = Tag(name="TestTag")
    meta = MetadataEntry(key="mime_type", value=mime_type)
    file = File(
        name=filename,
        original_path=path,
        new_path="/usr/fake/path/" + filename,
        tags=[tag],
        metadata=[meta]
    )


    kw_extractor = KWExtractor()
    result = kw_extractor.extract_kw(file)

    extracted_keywords = set(kw for kw, _ in result)

    def normalize(s): return s.lower().replace("-", "").strip(".")
    def check_majority(expected_keywords, result_keywords, threshold=0):
        normalized_results = {normalize(kw) for kw in result_keywords}
        normalized_expected = {normalize(kw) for kw in expected_keywords}
        matches = normalized_expected & normalized_results
        ratio = len(matches) / len(normalized_expected)
        assert ratio >= threshold, f"Too few matches: {matches} (Ratio: {ratio:.2f})"

    check_majority(expected_keywords, extracted_keywords)
