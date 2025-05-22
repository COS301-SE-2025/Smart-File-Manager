import pytest

from src.metadata_scraper import MetaDataScraper

# Check if any given file contains []
def test_standard_metadata():
    scraper = MetaDataScraper("testRootFolder/a.txt")
    print(scraper.get_standard_metadata())