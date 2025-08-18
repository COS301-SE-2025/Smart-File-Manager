from message_structure_pb2 import Directory, File, Tag, MetadataEntry
import os
from collections import defaultdict

TEST_DIR = os.path.dirname(__file__)

class DirectoryCreator:
    # construct the map of map of list of file
    def __init__(self, directoryName, fileMap):
        self.directory_name_idx = 0 
        self.directory_name = directoryName
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

    def merge(self, dir1, dir2):
        if dir1.name != dir2.name:
            raise ValueError("Cannot merge directories with different names")

        merged_dir = Directory()
        merged_dir.name = dir1.name
        merged_dir.path = self.get_path(dir1.name)

        merged_dir.files.extend(list(dir1.files))
        merged_dir.files.extend(list(dir2.files))

        merged_dir.directories.extend(list(dir1.directories))
        merged_dir.directories.extend(list(dir2.directories))

        return merged_dir



    def createFile(self, filename, dirName):
        file_info = self.file_map[filename]
        original_file_path = file_info["absolute_path"]
        new_file_path = self.get_path(f"{dirName}/{filename}")
        is_locked = file_info.get("is_locked", False)

        my_tags = self.createTags(file_info)
        return File(
            name=filename,
            original_path=original_file_path,
            new_path=new_file_path,
            tags=my_tags,  
            metadata=self.createMetaData(file_info),
            is_locked = is_locked

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
    
    
