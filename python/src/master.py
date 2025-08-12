from typing import Callable, Tuple
from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry, Keyword
from kw_extractor import KWExtractor
from full_vector import FullVector
import os
from k_means import KMeansCluster 

import time

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves, transformer):
        self.slaves = ThreadPoolExecutor(maxSlaves)
        self.scraper = MetaDataScraper()
        self.kw_extractor = KWExtractor()
        self.full_vec = FullVector(transformer)  

    # Takes gRPC request's root and sends it to be processed by a slave
    def submit_task(self, request : DirectoryRequest):
        self.start_time = time.time()
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    # Handles request appropriately 
    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Map request type to method and call
        requestHandler = {
            "CLUSTERING" : self.handle_clustering_request,
            "METADATA" : self.handle_metadata_request,
            "KEYWORDS" : self.handle_keyword_request
        }

        self.request_submitted = time.time()
        print("Submitting request: " + str(self.request_submitted - self.start_time))
        handler = requestHandler.get(request.requestType.upper())

        if not handler == None:
            return handler(request)
        else:
            reponse =  DirectoryResponse()
            reponse.response_code = 400
            reponse.response_msg = "Unknown Request type: Must be in  [CLUSTERING, METADATA, KEYWORDS]"
            return reponse

    
    def handle_clustering_request(self, request : DirectoryRequest) -> DirectoryResponse:
        
        # List of map where map contains keywords, tags and metadata
        files = []

        # Modifies directory request by adding metadata and creates map of files with metadata and keywords
        if self.extract_metadata(request.root, files, metadata_fn=self.scraper.get_standard_metadata, build_file_entry=True):

            self.extraction = time.time()
            print("Metadata and keywords extracted: " +str(self.extraction - self.start_time))
            # Modifies file list to add additional entry to each map i.e. full vector which contains all encoded data required for clustering
            self.full_vec.create_full_vector(files)
            self.full_vector_time = time.time()
            print("Full vector created: " + str(self.full_vector_time - self.start_time))

            # Append all full vectors
            full_vecs = []
            for file in files:
                full_vecs.append(file["full_vector"])
            self.full_vector_time_2 = time.time()
            print("Full vectors appended: " + str(self.full_vector_time_2 - self.start_time)) 

            # Recursively cluster and return a directory
            kmeans = KMeansCluster(5, 10, self.full_vec.model, request.root.name)
            response_directory = kmeans.dirCluster(full_vecs,files)
            self.clustering_time = time.time()
            print("Clustering complete: " + str(self.clustering_time - self.start_time))
            kmeans.printDirectoryTree(response_directory) 
            response = DirectoryResponse(root=response_directory, response_code=200, response_msg="Files successfully clustered")
            # print(response)
            self.response_time = time.time()
            print("Sending response: " + str(self.response_time - self.start_time))
            return response
        else:
            response = DirectoryResponse(response_code=400, response_msg="No file could be opened")
            return response

    def handle_metadata_request(self, request : DirectoryRequest) -> DirectoryResponse:
        
        if self.extract_metadata(request.root, files=[], metadata_fn=self.scraper.get_standard_metadata, build_file_entry=False):
            response = DirectoryResponse(root=request.root, response_code=200, response_msg="Successfully extracted at least some metadata")
            return response
        else:
            response = DirectoryResponse(root=request.root, response_code=400, response_msg="No file could be opened")

    def handle_keyword_request(self, request : DirectoryRequest) -> DirectoryResponse:

        self.kw_extractor.set_n(1)        

        if self.__keyword_extractor__(request.root):
            response = DirectoryResponse(root=request.root, response_code=200, response_msg="Successfully extracted keywords for at least 1 file")
        else:
            response = DirectoryResponse(root=request.root, response_code=400, response_msg="Could not extract any keywords")

        self.kw_extractor.set_n(3)
        return response 

    def __keyword_extractor__(self, currentDirectory : Directory) -> bool:

        success = False

        for curFile in currentDirectory.files:
            # Invariant: if file could not be opened extract_kw returns empty list
            keywords = self.kw_extractor.extract_kw(curFile)
            for word in keywords:
                success = True
                curFile.keywords.append(Keyword(keyword=word[0].lower(), score=word[1]))

        for curDir in currentDirectory.directories:
            if self.__keyword_extractor__(curDir):
                success = True

        return success

    def extract_metadata(
        self,
        currentDirectory: Directory,
        files: list,
        metadata_fn: Callable[[], dict],
        build_file_entry: bool = False
    ) -> bool:
        
        success = False  # Track if information for at least one file could be extracted

        for curFile in currentDirectory.files:
            try:
                self.scraper.set_file(os.path.abspath(curFile.original_path))
            except ValueError:
                meta_error = MetadataEntry(key="Error", value="File does not exist - could not extract metadata")
                curFile.metadata.append(meta_error)
                continue

            metadata_fn()
            metadata = self.scraper.metadata
            if not metadata:
                meta_error = MetadataEntry(key="Error", value="Metadata function returned None or empty")
                curFile.metadata.append(meta_error)
                continue

            # Update success
            success = True

            for k, v in metadata.items():
                curFile.metadata.append(MetadataEntry(key=str(k), value=str(v)))

            if build_file_entry:
                file_entry = dict(metadata)
                file_entry["keywords"] = self.kw_extractor.extract_kw(curFile)
                file_entry["tags"] = [tag.name.strip().lower() for tag in curFile.tags if tag.name]
                file_entry["is_locked"] = curFile.is_locked
                file_entry["original_path"] = curFile.original_path
                files.append(file_entry)

        for curDir in currentDirectory.directories:
            if self.extract_metadata(curDir, files, metadata_fn, build_file_entry):
                success = True # Ensure success did not change during rec call

        return success



