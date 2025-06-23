# Imports
import os
import datetime
from pathlib import Path
import magic # Used for MIME-Type
from PIL import Image # Used for images
from PIL.ExifTags import TAGS, GPSTAGS
import mutagen # Used for audio and video
from pypdf import PdfReader # used for pdf (shocker)
import docx
from pymediainfo import MediaInfo

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
            "mime_type": magic.from_file(self.__filedir, mime=True),
            "accessed": datetime.datetime.fromtimestamp(stat.st_atime).isoformat(),
            "owner_uid": stat.st_uid,
            "owner_gid": stat.st_gid,
            "mode": oct(stat.st_mode), # permissions idk why but apparently usually read as octal
            "inode": stat.st_ino
        })

    # Extracts image file metadata
    def get_image_metadata(self) -> dict:
        return_data = {}
        try:
            with Image.open(self.__filedir) as img:
                # exif
                exif_data = img._getexif()
                if exif_data is not None:
                    for tag, value in exif_data.items():
                        decoded = TAGS.get(tag, tag)
                        return_data[decoded] = value

                # Geolocation
                if "GPSInfo" in return_data:
                    gps_info = return_data["GPSInfo"]
                    gps_decoded = {}
                    for t in gps_info:
                        sub_decoded = GPSTAGS.get(t, t)
                        gps_decoded[sub_decoded] = gps_info[t]
                    return_data["GPSInfo"] = gps_decoded


                # dimensions
                return_data.update({
                    "image_width": img.width,
                    "image_height": img.height,
                    "image_mode": img.mode,
                    "image_format": img.format,
                })

                return return_data

        except Exception:
            return {}
        
    # Extract audio metadata
    def get_audio_video_metadata(self) -> dict:
        try:
            audio = mutagen.File(self.__filedir)
            if not audio:
                return {}
            return {k: str(v) for k, v in audio.items()}
        except Exception:
            return {}
    
    # Extract pdf metadata
    def get_pdf_metadata(self) -> dict:
        try:    
                returndata = {}
                pdf = PdfReader(self.__filedir)
                returndata = dict(pdf.metadata) or {}
                returndata["num_pages"] = len(pdf.pages)
                return returndata
        except Exception:
            return {}
        
    # Extract .docx metadata
    def get_docx_metadata(self) -> dict:
        try:
            doc = docx.Document(self.__filedir)
            core_props = doc.core_properties
            return {
                "author": core_props.author,
                "title": core_props.title,
                "subject": core_props.subject,
                "created": core_props.created.isoformat() if core_props.created else None,
                "modified": core_props.modified.isoformat() if core_props.modified else None,
                "category": core_props.category,
                "comments": core_props.comments,
                "keywords": core_props.keywords,
                "language": core_props.language,
                "last_modified_by": core_props.last_modified_by,
                "revision": core_props.revision,
                "version": core_props.version,

            }
        except Exception:
            return {}
        
    # Extract video metadata
    def get_video_metadata(self):
        try:
            media_info = MediaInfo.parse(self.__filedir)
            video_data = {}
            for track in media_info.tracks:
                if track.track_type == 'Video':
                    video_data.update({
                        "video_format": track.format,
                        "duration_ms": track.duration,
                        "frame_rate": track.frame_rate,
                        "width": track.width,
                        "height": track.height,
                        "bit_rate": track.bit_rate,
                        "codec": track.codec_id,
                    })
            return video_data
        except Exception:
            return {}


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
        elif "audio" in mime_type:
            self.metadata.update(self.get_audio_video_metadata())
        elif "video" in mime_type:
            self.metadata.update(self.get_video_metadata())
        elif mime_type == "application/pdf":
            self.metadata.update(self.get_pdf_metadata())
        elif mime_type in ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
            self.metadata.update(self.get_docx_metadata())