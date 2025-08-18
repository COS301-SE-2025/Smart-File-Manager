from message_structure_pb2 import Directory, File, Tag, MetadataEntry
import os
from collections import defaultdict

TEST_DIR = os.path.dirname(__file__)

class DirectoryCreator:
    """
    Root = self.FILE_DIR (e.g. 'test_files_3').

    Conventions:
    - buildDirectory(rel_path, ...) takes a RELATIVE path under the root.
      * rel_path == ""   -> the root directory node
      * rel_path == "a"  -> child "a"
      * rel_path == "a/b"-> grandchild "a/b"
    - For compatibility with print, the root directory's
      Directory.path is rendered as '<root>/<root>'.
    - File.new_path is always '<root>/<full relative path>/<filename>'.
      For the root level files: '<root>/<root>/<filename>'.
    """

    def __init__(self, directoryName, fileMap):
        self.directory_name_idx = 0
        self.directory_name = directoryName   # display name of the root
        self.FILE_DIR = directoryName         # root folder
        self.file_map = defaultdict()
        for file in fileMap:
            self.file_map[file["filename"]] = file

    def get_path(self, rel_name: str) -> str:

        return os.path.join(self.FILE_DIR, rel_name) if rel_name else self.FILE_DIR

    def buildDirectory(self, rel_path: str, files, children):
        """
        rel_path: '' for root, otherwise 'a', 'a/b', ...
        """
        # Display name for the printing
        display_name = self.directory_name if rel_path == "" else os.path.basename(rel_path)

        # Path formatting:
        #  - Root dir path shows as '<root>/<root>' (matches your sample)
        #  - Others: '<root>/<rel_path>'
        dir_path = self.get_path(self.directory_name if rel_path == "" else rel_path)

        file_objs = [self.createFile(f["filename"], rel_path) for f in files]

        return Directory(
            name=display_name,
            path=dir_path,
            files=file_objs,
            directories=children
        )

    def merge(self, dir1, dir2):
        if dir1.name != dir2.name:
            raise ValueError("Cannot merge directories with different names")

        merged_dir = Directory()
        merged_dir.name = dir1.name
        merged_dir.path = self.get_path(self.directory_name)  # keep root '<root>/<root>' style

        merged_dir.files.extend(list(dir1.files))
        merged_dir.files.extend(list(dir2.files))

        merged_dir.directories.extend(list(dir1.directories))
        merged_dir.directories.extend(list(dir2.directories))

        return merged_dir

    def createFile(self, filename: str, rel_path: str):
        """
        rel_path is the relative directory path for this file.
        """
        file_info = self.file_map[filename]
        original_file_path = file_info["absolute_path"]
        is_locked = file_info.get("is_locked", False)

        # For root, new files appear under '<root>/<root>/<filename>'
        # For others: '<root>/<rel_path>/<filename>'
        path_base = self.directory_name if rel_path == "" else rel_path
        new_file_path = self.get_path(os.path.join(path_base, filename))

        return File(
            name=filename,
            original_path=original_file_path,
            new_path=new_file_path,
            tags=self.createTags(file_info),
            metadata=self.createMetaData(file_info),
            is_locked=is_locked
        )

    def createTags(self, file_info):
        if "tags" in file_info:
            return [Tag(name=tag) for tag in file_info["tags"]]
        return []

    def createMetaData(self, file_info):
        return [
            MetadataEntry(key=key, value=str(val))
            for key, val in file_info.items()
            if key not in ("keywords", "full_vector")
        ]