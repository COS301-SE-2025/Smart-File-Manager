from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
from kw_extractor import KWExtractor
from full_vector import FullVector
from vocabulary import Vocabulary
import os

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves):
        self.slaves = ThreadPoolExecutor(maxSlaves)
        self.scraper = MetaDataScraper()
        self.kw_extractor = KWExtractor()
        self.vocab = Vocabulary()
        self.full_vec = FullVector()


    # Takes gRPC request's root and sends it to be processed by a slave
    def submitTask(self, request : DirectoryRequest):
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Modifies directory request by adding metadata
        self.scrapeMetadata(request.root)
        response_directory = request.root
        kw_response = self.kw_extractor.extract_kw(request.root)
        vocabKW = self.vocab.createVocab(kw_response)
        cluster_vec = self.full_vec.assignTF_IDF(kw_response,vocabKW)
        for vec in cluster_vec:
            print("\n")
            print(vec)
            print("\n")

        response = DirectoryResponse(root=response_directory)
        return response
    
        '''
        What is this method supposed to really do?
        * Extract metadata
        * Perform clustering
        * Generate new tree strucutre as gRPC Directory

        Preferably each of the above should be handled by its own class for seperation of concerns
        '''

    # Traverses Directory recursively and extracts metadata for each file
    def scrapeMetadata(self, currentDirectory : Directory) -> None:

        # Extract metadata
        for curFile in currentDirectory.files:
            # Ensure file path is valid
            try:
                self.scraper.set_file(os.path.abspath(curFile.original_path))
            except ValueError:
                # Invalid path => add error tag to metadata entry
                meta_error = MetadataEntry(key="Error", value="File does not exist - could not extract metadata")
                curFile.metadata.append(meta_error)
                continue
            else:
                # Valid path => scrape
                self.scraper.get_metadata()
                extracted_metadata = self.scraper.metadata
                for k,v in extracted_metadata.items():

                    meta_entry = MetadataEntry(key=str(k), value = str(v))
                    curFile.metadata.append(meta_entry)

        # Recurisve call
        if len(currentDirectory.directories) != 0:
            for curDir in currentDirectory.directories:
                self.scrapeMetadata(curDir)

########################
tag1 = Tag(name="ImFixed")
meta1 = MetadataEntry(key="author", value="johnny")
meta4 = MetadataEntry(key="mime_type", value="text/plain")
meta2 = MetadataEntry(key="mime_type", value="application/pdf")
meta3 = MetadataEntry(key="mime_type", value="application/msword")

file1 = File(
    name="gopdoc.pdf",
    original_path="python/testing/test_files/myPdf.pdf",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta2]
)
file2 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/testFile.txt",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta4]
)
file3 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/myWordDoc.docx",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta3]
)

dir1 = Directory(
    name="useless_files",
    path="/usr/trash",
    files=[file1, file2,file3],
    directories=[]
)
req = DirectoryRequest(root=dir1) 

if __name__ == "__main__":
    master = Master(1)
    master.process(req)