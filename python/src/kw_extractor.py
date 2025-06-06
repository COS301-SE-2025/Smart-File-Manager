import time
import docx
from yake import KeywordExtractor
from pypdf import PdfReader
from message_structure_pb2 import Directory, DirectoryRequest, File, MetadataEntry, Tag

# Receive data as D
class KWExtractor:
    #Yake instance
    def __init__(self):
        self.yake_extractor = KeywordExtractor()

    #Main extractor function
    def extract_kw(self, input):
        result = []
        for file in input.files:    
            file_name = f"{file.original_path}"
            mime_type = next((entry.value for entry in file.metadata if entry.key == "mime_type"), None)
            result += self.open_file(file_name,mime_type, 1) #Seconds based time limit
        return self.list_to_map(result, 10)   
    
    #Make the result into a sorted map with max keywords
    def list_to_map(self, result, max_keywords):
        return_map = {}

        for file_name, keywords in result:
            #descending order
            sorted_keywords = sorted(keywords, key=lambda x: x[1], reverse=True)
            top_keywords = [kw for kw, _ in sorted_keywords[:max_keywords]]
            return_map[file_name] = top_keywords

        return return_map


    #open a file (check which type and send to be opened in the correct way)
    def open_file(self, file_name, file_type, max_duration_seconds=1):
        result = []
        if file_type == "application/pdf":
            keywords = self.pdf_extraction(file_name, '.', max_duration_seconds)
            result.append((file_name, keywords))
        elif file_type in ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
            keywords = self.docx_extraction(file_name, '.', max_duration_seconds)
            result.append((file_name, keywords))
        elif file_type == "text/plain":      
            keywords = self.def_extraction(file_name, '.', max_duration_seconds)    
            result.append((file_name, keywords))
        else:
            print("Unkown type, attempting extract...")
            try:
                keywords = self.def_extraction(file_name, '.', max_duration_seconds)
                result.append((file_name, keywords))
            except:
                print("Error occured trying to read unkown type")
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
        reader = PdfReader(file_name)
        for k in range(len(reader.pages)):
                page = reader.pages[k]
                for sentence in self.split_by_delimiter_pdf(page.extract_text(), delimiter):
                    elapsed = time.time() - start_time
                    if(elapsed > max_duration_seconds):
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
    
   

# if __name__ == "__main__":
#     kw_extractor = KWExtractor()
#     result = kw_extractor.extract_kw(req.root)


#     for filename, keywords in result.items():
#         print(f"\n== FILE: {filename} ==")
#         for kw in keywords:
#             print(kw)
    

        



