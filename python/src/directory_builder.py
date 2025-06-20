from message_structure_pb2 import DirectoryResponse, Directory
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
import os
#temp
from collections import defaultdict

TEST_DIR = os.path.dirname(__file__)

class DirectoryCreator:
    def __init__(self, directoryName, fileMap):
        self.directory_name_idx = 0 
        self.directory_name = "Directory"
        self.FILE_DIR = os.path.join(TEST_DIR, directoryName)
        self.file_map = defaultdict()
        for file in fileMap:
            self.file_map[file["filename"]] = file


    def get_path(self, name):
        return os.path.join(self.FILE_DIR, name)

    # Recursive function
    # Create directory for each pass
    # Create directory within passes
    # Create files at end of map
    # label map is a map of maps
    # map
    # [
    # clusterpassone [0,1,2]
    # -> 0: file1, file2, ...
    # -> 1: ...
    # -> 2: ...
    # clusterpasstwo [0,1,2,3]
    # -> 0: file1, file2, ...
    # -> 1: ...
    # -> 2: ...
    # -> 3: ...
    # ...
    #]


    def buildDirectory(self, name, files, children):
        dir_path = self.get_path(name)
        file_objs = [self.createFile(f["filename"], name) for f in files]
        return Directory(
            name=name,
            path=dir_path,
            files=file_objs,
            directories=children
        )

    def createFile(self, filename, dirName):
        file_info = self.file_map[filename]
        original_file_path = file_info["absolute_path"]
        new_file_path = self.get_path(f"{dirName}/{filename}")

        return File(
            name=filename,
            original_path=original_file_path,
            new_path=new_file_path,
            tags=self.createTags(),
            metadata=self.createMetaData(file_info)
        )

    def createTags(self):
        return []

    def createMetaData(self, file_info):
        return [
            MetadataEntry(key=key, value=str(val))
            for key, val in file_info.items()
            if key not in ("keywords", "full_vector")
        ]