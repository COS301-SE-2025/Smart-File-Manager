from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
from kw_extractor import KWExtractor
from full_vector import FullVector
import os
from k_means import KMeansCluster 
from hdb_cluster import HDBSCANCluster
#temp
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

    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Contains list of objects which stores metadata, keywords etc...
        files = []

        # Modifies directory request by adding metadata and creates map of files with metadata and keywords
        self.getFileInfo(request.root, files)
        response_directory = request.root
        #print(files)

        self.full_vec.create_full_vector(files)
        full_vecs = []
        for file in files:
            print(file["full_vector"])
            print(file["filename"])
            full_vecs.append(file["full_vector"])

        
        kmeans = KMeansCluster(6)
        labels = kmeans.cluster(full_vecs)
        print(labels)
        # hdbscan_cluster = HDBSCANCluster(min_cluster_size=3)
        # labels = hdbscan_cluster.cluster(full_vecs)

        
        label_to_filenames = defaultdict(list)

        for index, file in enumerate(files):
            label = labels[index]
            filename = file["filename"]
            label_to_filenames[label].append(filename)

        # Optional: print the grouped result
        for label in sorted(label_to_filenames):
            print(f"Label {label}:")
            for filename in label_to_filenames[label]:
                print(f"  - {filename}")


        # Metadata Request ==> Return Directory with metadata attached
        response = DirectoryResponse(root=response_directory)
        return response
    

    # Traverses Directory recursively and extracts metadata and keywords for each file
    def getFileInfo(self, currentDirectory: Directory, files: list) -> None:
        for curFile in currentDirectory.files:
            try:
                self.scraper.set_file(os.path.abspath(curFile.original_path))
            except ValueError:
                meta_error = MetadataEntry(key="Error", value="File does not exist - could not extract metadata")
                curFile.metadata.append(meta_error)
                continue

            # Extract metadata
            self.scraper.get_standard_metadata()
            extracted_metadata = self.scraper.metadata

            # Add to curFile.metadata
            for k, v in extracted_metadata.items():
                curFile.metadata.append(MetadataEntry(key=str(k), value=str(v)))

            # Build output entry
            file_entry = {}
            file_entry.update(extracted_metadata)

            # Extract keywords
            keywords = self.kw_extractor.extract_kw(curFile)
            file_entry["keywords"] = keywords

            files.append(file_entry)

        for curDir in currentDirectory.directories:
            self.getFileInfo(curDir, files)


import numpy as np
if __name__ == "__main__":
#    master = Master(1)
    kmeans = KMeansCluster(3)
    X = [
        [np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), np.float64(0.1761190379942198), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.12721378280709897],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, np.float64(0.14985554132463938), np.float64(0.14985554132463938), np.float64(0.1453011207686066), np.float64(0.13775170326151098), np.float64(0.13588457081573124), np.float64(0.132439301673504), np.float64(0.12174767379894379), np.float64(0.12051367106078721), np.float64(0.11645295367216063), np.float64(0.09314785652067586), 0.0, 0.0, 1.0, 0.0, 0.0, 0.9924333268750425],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.7575546001939815],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.02837284236359479],
        [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0]
    ]
    newX = []
    for f in X:
        newX.append(f[0:20])
    labels = kmeans.cluster(newX)
    print(labels)

# docx, jpg, pdf, webm, webp, jpg
# todo, myimg, mypdf, myvid, thumb, holiday

#todo(docx),thumb(webp) -> closest in size
#myimg, holiday -> same types
#mypdf, myvid -> closest in size?

