from message_structure_pb2 import Directory, File, Tag, MetadataEntry
import os
from collections import defaultdict

TEST_DIR = os.path.dirname(__file__)

class DirectoryCreator:
    # construct the map of map of list of file
    def __init__(self, directoryName, fileMap):
        self.directory_name_idx = 0 
        self.directory_name = "Directory"
        #self.FILE_DIR = os.path.join(TEST_DIR, directoryName)
        self.FILE_DIR = directoryName
        self.file_map = defaultdict()
        for file in fileMap:
            self.file_map[file["filename"]] = file


    def get_path(self, name):
        return os.path.join(self.FILE_DIR, name)


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
        # new_file_path = self.get_path(f"{dirName}/{filename}")
        new_file_path = os.path.join(dirName, filename)


        my_tags = self.createTags(file_info)

        return File(
            name=filename,
            original_path=original_file_path,
            new_path=new_file_path,
            tags=my_tags,  
            metadata=self.createMetaData(file_info)
        )


    def createTags(self, file_info):
        if "tags" in file_info: 
            tags = []
            for tag in file_info["tags"]:
               # print(tag)
               tags.append(Tag(name=tag)) 
            return tags
        else:
            return []


    def createMetaData(self, file_info):
        return [
            MetadataEntry(key=key, value=str(val))
            for key, val in file_info.items()
            if key not in ("keywords", "full_vector")
        ]
    
    