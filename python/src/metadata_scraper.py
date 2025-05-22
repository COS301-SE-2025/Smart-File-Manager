# Imports
import os
import mimetypes
import datetime
from pathlib import Path

import magic # Used for MIME-Type
from PIL import Image # Used for images
from PIL.ExifTags import TAGS
import mutagen # Used for audio
from pypdf import PdfReader # used for pdf (shocker)
import docx

# Used to extract metadata from various files
class MetaDataScraper:
    def __init__(self, filedir):
        self.filedir = filedir
        self.metadata = {}

    def set_file(self, new_filedir):
        self.filedir = new_filedir

    # Metadata that should be extracted from any file
    def get_standard_metadata(self):
        path = Path(self.filedir)
        stat = path.stat()

        self.metadata.update({
            "filename" : path.name,
            "absolute_path" : str(path.resolve()),
            "file_extension" : path.suffix.lower(),
            "size_bytes": stat.st_size,
            "created": datetime.datetime.fromtimestamp(stat.st_ctime).isoformat(),
            "modified": datetime.datetime.fromtimestamp(stat.st_mtime).isoformat(),
            "mime_type": magic.from_file(self.filedir, mime=True)
        })

    # Extracts image file metadata
    def get_image_metadata(filepath):
        try:
            with Image.open(filepath) as img:
                exif_data = img._getexif()
                if not exif_data:
                    return {}
                return {
                    TAGS.get(tag): value for tag, value in exif_data.items() if tag in TAGS
                }
        except Exception:
            return {}
        
    # Extract audio metadata
    def get_audio_metadata(filepath):
        try:
            audio = mutagen.File(filepath)
            if not audio:
                return {}
            return {k: str(v) for k, v in audio.items()}
        except Exception:
            return {}
    
    # Extract pdf metadata
    def get_pdf_metadata(filepath):
        try:
                pdf = PdfReader(filepath)
                return dict(pdf.metadata or {})
        except Exception:
            return {}
        
    # Extract .docx metadata
    def get_docx_metadata(filepath):
        try:
            doc = docx.Document(filepath)
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
        