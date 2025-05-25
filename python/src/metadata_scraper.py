# Imports
import os
import mimetypes
import datetime
from pathlib import Path

import magic # Used for MIME-Type
from PIL import Image # Used for images
from PIL.ExifTags import TAGS
import mutagen # Used for audio and video
from pypdf import PdfReader # used for pdf (shocker)
import docx

# Used to extract metadata from various files
class MetaDataScraper:
    def __init__(self):
        self.metadata = {}
        self.__filedir = ""

    def set_file(self, new_filedir) -> None:
        if os.path.exists(new_filedir):
            self.__filedir = new_filedir
            self.metadata.clear()
        else:
            raise ValueError("New path does not exist on this filesystem")

    # Metadata that should be extracted from any file
    def get_standard_metadata(self) -> None:
        path = Path(self.__filedir)
        stat = path.stat()

        self.metadata.update({
            "filename" : path.name,
            "absolute_path" : str(path.resolve()),
            "file_extension" : path.suffix.lower(),
            "size_bytes": stat.st_size,
            "created": datetime.datetime.fromtimestamp(stat.st_ctime).isoformat(),
            "modified": datetime.datetime.fromtimestamp(stat.st_mtime).isoformat(),
            "mime_type": magic.from_file(self.__filedir, mime=True)
        })

    # Extracts image file metadata
    def get_image_metadata(self) -> None:
        try:
            with Image.open(self.__filedir) as img:
                exif_data = img._getexif()
                if not exif_data:
                    return {}
                return {
                    TAGS.get(tag): value for tag, value in exif_data.items() if tag in TAGS
                }
        except Exception:
            return {}
        
    # Extract audio metadata
    def get_audio_video_metadata(self) -> None:
        try:
            audio = mutagen.File(self.__filedir)
            if not audio:
                return {}
            return {k: str(v) for k, v in audio.items()}
        except Exception:
            return {}
    
    # Extract pdf metadata
    def get_pdf_metadata(self) -> None:
        try:
                pdf = PdfReader(self.__filedir)
                return dict(pdf.metadata or {})
        except Exception:
            return {}
        
    # Extract .docx metadata
    def get_docx_metadata(self) -> None:
        try:
            doc = docx.Document(self.__filedir)
            core_props = doc.core_properties
            return {
                "author": core_props.author,
                "title": core_props.title,
                "subject": core_props.subject,
                "created": core_props.created.isoformat() if core_props.created else None,
                "modified": core_props.modified.isoformat() if core_props.modified else None,
            }
        except Exception:
            return {}
        
    # Extract video metadata

    # Extract all metadata possible
    def get_metadata(self) -> None:

        # Returns empty metadata if directory not set
        if self.__filedir == "":
            return

        self.get_standard_metadata()

        # Check file MIME type to decide what other metadata to scan
        mime_type = str(self.metadata["mime_type"])
        if "image" in mime_type:
            self.metadata.update(self.get_image_metadata())
        elif "audio" in mime_type or "video" in mime_type:
            self.metadata.update(self.get_audio_video_metadata())
        elif mime_type == "application/pdf":
            self.metadata.update(self.get_pdf_metadata())
        elif mime_type in ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
            self.metadata.update(self.get_docx_metadata())