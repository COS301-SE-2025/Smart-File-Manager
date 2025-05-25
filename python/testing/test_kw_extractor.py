from unittest.mock import patch, MagicMock

# < ------ UNIT TESTING ------>
#kw extraction from sentence
@patch("src.kw_extractor.kw.extract")
def test_kw_extract(mock_kw_extract):
    mock_sentence = MagicMock()
    mock_sentence.__enter.__return_value = mock_sentence
    mock_sentence.__getexif.return_value = {
        ""
    }
    mock_kw_extract.return_value = mock_sentence

    from src.kw_extractor import KWExtractor
    extractor = KWExtractor()
    keywords = extractor.extract_kw("Here is a very short. Broken up. Sentence")

    assert keywords == []

# kw extraction plain text

#kw extraction docx

#kw extraction pdf

#kw extraction unkown
