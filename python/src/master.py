from typing import Callable
from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
from kw_extractor import KWExtractor
from full_vector import FullVector
import os
from k_means import KMeansCluster 
from collections import defaultdict

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves):
        self.slaves = ThreadPoolExecutor(maxSlaves)
        self.scraper = MetaDataScraper()
        self.kw_extractor = KWExtractor()
        self.full_vec = FullVector()  

    # Takes gRPC request's root and sends it to be processed by a slave
    def submitTask(self, request : DirectoryRequest):
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    # Handles request appropriately 
    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Map request type to method and call
        requestHandler = {
            "CLUSTERING" : self.handleClusteringRequest,
            "METADATA" : self.handleMetadatRequest
        }

        handler = requestHandler.get(request.requestType.upper())

        if not handler == None:
            return handler(request)
        else:
            reponse =  DirectoryResponse()
            reponse.response_code = 400
            reponse.response_msg = "Unknown Request type: Must be in  [CLUSTERING, METADATA]"
            return reponse

    
    def handleClusteringRequest(self, request : DirectoryRequest) -> DirectoryResponse:
        
        # List of map where map contains keywords, tags and metadata
        files = []

        # Modifies directory request by adding metadata and creates map of files with metadata and keywords
        if self.extract_metadata(request.root, files, metadata_fn=self.scraper.get_standard_metadata, build_file_entry=True):

            # Modifies file list to add additional entry to each map i.e. full vector which contains all encoded data required for clustering
            self.full_vec.create_full_vector(files)

            # Append all full vectors
            full_vecs = []
            for file in files:
                full_vecs.append(file["full_vector"])

            # Recursively cluster and return a directory
            kmeans = KMeansCluster(int(len(full_vecs)*(1/6)))
            response_directory = kmeans.dirCluster(full_vecs,files)
            kmeans.printDirectoryTree(response_directory) 
            response = DirectoryResponse(root=response_directory, response_code=200, response_msg="Files successfully clustered")
            return response
        else:
            response = DirectoryResponse(response_code=400, response_msg="No file could be opened")
            return response

    def handleMetadatRequest(self, request : DirectoryRequest) -> DirectoryResponse:
        
        if self.extract_metadata(request.root, files=[], metadata_fn=self.scraper.get_standard_metadata, build_file_entry=False):
            response = DirectoryResponse(root=request.root, response_code=200, response_msg="Successfully extracted at least some metadata")
            return response
        else:
            response = DirectoryResponse(root=request.root, response_code=400, response_msg="No file could be opened")


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
                files.append(file_entry)

        for curDir in currentDirectory.directories:
            if self.extract_metadata(curDir, files, metadata_fn, build_file_entry):
                success = True # Ensure success did not change during rec call

        return success



