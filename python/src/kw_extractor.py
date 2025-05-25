#RAKE seems like a decent option
#YAKE could be better for simple context based comparisons down the line

#pip install yake
#pip install pypdf

#GeeksForGeeks example
# Create a KeywordExtractor instance

# # Text from which keywords will be extracted
# text = "YAKE (Yet Another Keyword Extractor) is a Python library for extracting keywords from text."

# # Extract keywords from the text
# keywords = kw_extractor.extract_keywords(text)

# # Print the extracted keywords and their scores
# for kw in keywords:
#     print("Keyword:", kw[0], "Score:", kw[1])
import re
from yake import KeywordExtractor
from pypdf import PdfReader
from message_structure_pb2 import Directory, DirectoryRequest, File, MetadataEntry, Tag

# Receive data as D
class KWExtractor:
    #Yake instance
    yake_extractor = KeywordExtractor()

    #Main extractor function
    def extract_kw(self, input):
        result = []
        for file in input.files:    
            file_name = f"{file.original_path}"
            mime_type = next((entry.value for entry in file.metadata if entry.key == "mime_type"), None)
            result += self.open_file(file_name, 1, mime_type) #Only process x sentences per file #May be larger but just for now so that not too many lines are process
        return result
    
    #open a file (check which type and send to be opened in the correct way)
    def open_file(self, file_name, max_sentences, file_type):
        result = []
        if file_type == "application/pdf":
            print("Pdf extraction")
            result += self.pdf_extraction(file_name, '.', max_sentences)
        elif file_type in ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
            print("word extraction")
            result += ["word"]
        elif file_type == "text/plain":      
            print("plain text extraction")
            result += self.def_extraction(file_name, '.', max_sentences)    
        else:
            print("Unkown type, attempting extract...")
            #result += [f"Unsupported File Type: {file_type}"]     
            try:
                result += self.def_extraction(file_name, '.', max_sentences)
            except:
                print("Error occured trying to read unkown type")
        return result
    
    #open a file with no mime_type (txt) or "text/plain"
    def def_extraction(self, file_name, delimiter, max_sentences):
        counter = 0
        result = []
        for sentence in self.split_by_delimiter_def(file_name, '.'):
            if(counter > max_sentences):
                break
            counter += 1
            result += self.get_kw(sentence)
        return result

    #Open a PDF file
    def pdf_extraction(self, file_name, delimiter, max_sentences):
        counter = 0
        result = []
        reader = PdfReader(file_name)
        for k in range(len(reader.pages)):
                page = reader.pages[k]
                for sentence in self.split_by_delimiter_pdf(page.extract_text(), '.'):
                    if(counter > max_sentences):
                        return result
                    result += self.get_kw(sentence)
                    counter += 1
        return result

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

    #extract keywords from the given sentence using yake
    def get_kw(self, sentence):
        keyword = self.yake_extractor.extract_keywords(sentence)        
        return keyword

        

    


            

    

tag1 = Tag(name="ImFixed")
meta1 = MetadataEntry(key="author", value="johnny")
meta4 = MetadataEntry(key="mime_type", value="text/plain")
meta2 = MetadataEntry(key="mime_type", value="application/pdf")
meta3 = MetadataEntry(key="mime_type", value="application/msword")

file1 = File(
    name="gopdoc.pdf",
    original_path="python/testing/test_files/myPdf.pdf",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta2]
)
file2 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/testFile.txt",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta4]
)
file3 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/myWordDoc.docx",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta3]
)

dir1 = Directory(
    name="useless_files",
    path="/usr/trash",
    files=[file1, file2, file3],
    directories=[]
)
req = DirectoryRequest(root=dir1)        

if __name__ == "__main__":
    kw_extractor = KWExtractor()
    result = kw_extractor.extract_kw(req.root)
    #for kw in result:
     #  print("Keyword:", kw[0], "Score:", kw[1])

        



