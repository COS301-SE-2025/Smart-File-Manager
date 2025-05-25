from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
import os

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves):
        self.slaves = ThreadPoolExecutor(maxSlaves)
        self.scraper = MetaDataScraper()

    # Takes gRPC request's root and sends it to be processed by a slave
    def submitTask(self, request : DirectoryRequest):
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Modifies directory request by adding metadata
        self.scrapeMetadata(request.root)
        response_directory = request.root
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
            self.scraper.set_file(os.path.abspath(curFile.original_path))
            self.scraper.get_metadata()
            extracted_metadata = self.scraper.metadata
            for k,v in extracted_metadata:
                meta_entry = MetadataEntry(key=k, value=v)
                curFile.metadata = meta_entry
            curFile.metadata = extracted_metadata

        # Recurisve call
        if len(currentDirectory.directories) != 0:
            for curDir in currentDirectory.directories:
                self.scrapeMetadata(curDir)