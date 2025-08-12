import time
from typing import List, Tuple
import docx
from yake import KeywordExtractor
import fitz  # PyMuPDF
# from pypdf import PdfReader
from message_structure_pb2 import File

# Keyword extractor class
# Given a file as input extracts the top 10 keywords along with their value from file
# Supports extraction for: text/plain, pdf, docx 

class KWExtractor:
    #Yake instance
    def __init__(self):
        self.yake_extractor = KeywordExtractor(lan="en", n=3)
        self.mime_handlers = {
        "application/pdf": self.pdf_extraction,
        "application/msword": self.docx_extraction,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": self.docx_extraction,
        "text/plain": self.def_extraction
        }
        self.confidence_threshold = 0.085
        self.supported_mime_types = set(self.mime_handlers.keys())


    def set_n(self, new_val : int):
        if new_val > 0:
            self.yake_extractor.n = new_val

    #Main extractor function
    def extract_kw(self, input: File) -> List[Tuple[str, float]]:
        file_name = input.original_path
        mime_type = next((entry.value for entry in input.metadata if entry.key == "mime_type"), None)

        if mime_type is None:
            if file_name.endswith(".pdf"):
                mime_type = "application/pdf"
            elif file_name.endswith(".docx") or file_name.endswith(".doc"):
                mime_type = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
            elif file_name.endswith(".txt"):
                mime_type = "text/plain"
            else:
                return []  

        if mime_type not in self.supported_mime_types:
            return []  

        result = self.open_file(file_name, mime_type, 1)
        if not result:
            return []

        _, keywords = result[0]
        filtered = [(kw, score) for kw, score in keywords if score < self.confidence_threshold]
        if len(filtered) < 3:
            filtered = keywords[:3]
        return filtered

    

    #Open a file (check which type and send to be opened in the correct way)
    def open_file(self, file_name, file_type, max_duration_seconds=1):
        result = []

        handler = self.mime_handlers.get(file_type, self.def_extraction)
        
        try:
            keywords = handler(file_name, '.', max_duration_seconds)
            result.append((file_name, keywords))
        except Exception as e:
            # print(f"Could not extract keywords from {file_name} | error {e}")
            pass
        
        return result

    
    #open a file with no mime_type (txt) or "text/plain"
    def def_extraction(self, file_name, delimiter, max_duration_seconds=1):
        start_time = time.time()
        final_sentence = ""
        for sentence in self.split_by_delimiter_def(file_name, delimiter):
            elapsed = time.time() - start_time
            if(elapsed > max_duration_seconds):
                return self.get_kw(final_sentence)
            final_sentence += sentence + delimiter
        return self.get_kw(final_sentence)
        

    #Open a PDF file

    def pdf_extraction(self, file_name, delimiter, max_duration_seconds=1):
        start_time = time.time()
        final_sentence = ""
        doc = fitz.open(file_name)
        for page in doc:
            text = page.get_text()
            for sentence in self.split_by_delimiter_pdf(text, delimiter):
                if time.time() - start_time > max_duration_seconds:
                    return self.get_kw(final_sentence)
                final_sentence += sentence + delimiter
        return self.get_kw(final_sentence)


        
    
    #open docx file
    def docx_extraction(self, file_name, delimiter, max_duration_seconds=1):
        start_time = time.time()
        final_sentence = ""
        doc = docx.Document(file_name)     
        for paragraph in doc.paragraphs:
            for sentence in self.split_by_delimiter_docx(paragraph.text, delimiter):
                elapsed = time.time() - start_time
                if(elapsed > max_duration_seconds):
                    return self.get_kw(final_sentence)
                final_sentence += sentence + delimiter
        return self.get_kw(final_sentence)
        
                

    #default delimiter
    def split_by_delimiter_def(self, file_name, delimiter):
        buffer = ''
        with open(file_name, 'r', encoding='utf-8') as file:
            for line in file:
                buffer += line
                while delimiter in buffer:
                    sentence, buffer = buffer.split(delimiter, 1)
                    yield sentence.strip() + delimiter
        if buffer.strip():
            yield buffer.strip()

    #pdf delimiter
    def split_by_delimiter_pdf(self, page_text, delimiter):
        buffer = page_text
        while delimiter in buffer:
            sentence, buffer = buffer.split(delimiter, 1)
            yield sentence.strip() + delimiter
        if buffer.strip():
            yield buffer.strip()

    #docx delimiter
    def split_by_delimiter_docx(self, page_text, delimiter):
        buffer = page_text
        while delimiter in buffer:
            sentence, buffer = buffer.split(delimiter, 1)
            yield sentence.strip() + delimiter
        if buffer.strip():
            yield buffer.strip()

    #extract keywords from the given sentence using yake
    def get_kw(self, sentence):
        keyword = self.yake_extractor.extract_keywords(sentence)        
        return keyword
    