#RAKE seems like a decent option
#YAKE could be better for simple context based comparisons down the line

#pip install yake
#pip install pypdf
#pip install python-docx

import docx
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
            result += self.open_file(file_name, 10, mime_type) #Only process x sentences per file #May be larger but just for now so that not too many lines are process
        return result   
    

    #open a file (check which type and send to be opened in the correct way)
    def open_file(self, file_name, max_sentences, file_type):
        result = []
        if file_type == "application/pdf":
            print("Pdf extraction")
            keywords = self.pdf_extraction(file_name, '.', max_sentences)
            result.append((file_name, keywords))
        elif file_type in ["application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
            print("word extraction")
            keywords = self.docx_extraction(file_name, '.', max_sentences)
            result.append((file_name, keywords))
        elif file_type == "text/plain":      
            print("plain text extraction")
            keywords = self.def_extraction(file_name, '.', max_sentences)    
            result.append((file_name, keywords))
        else:
            print("Unkown type, attempting extract...")
            try:
                keywords = self.def_extraction(file_name, '.', max_sentences)
                result.append((file_name, keywords))
            except:
                print("Error occured trying to read unkown type")
        return result
    
    #open a file with no mime_type (txt) or "text/plain"
    def def_extraction(self, file_name, delimiter, max_sentences):
        counter = 0
        result = []
        final_sentence = ""
        for sentence in self.split_by_delimiter_def(file_name, delimiter):
            if(counter > max_sentences):
                break
            counter += 1
            final_sentence += sentence + delimiter
        result = self.get_kw(final_sentence)
        return result

    #Open a PDF file
    def pdf_extraction(self, file_name, delimiter, max_sentences):
        counter = 0
        result = []
        final_sentence = ""
        reader = PdfReader(file_name)
        for k in range(len(reader.pages)):
                page = reader.pages[k]
                for sentence in self.split_by_delimiter_pdf(page.extract_text(), delimiter):
                    if(counter > max_sentences):
                        break
                    counter += 1
                    final_sentence += sentence + delimiter
        result = self.get_kw(final_sentence)
        return result
    #open docx file
    def docx_extraction(self, file_name, delimiter, max_sentences):
        counter = 0
        result = []
        final_sentence = ""
        doc = docx.Document(file_name)     
        for paragraph in doc.paragraphs:
            for sentence in self.split_by_delimiter_docx(paragraph.text, delimiter):
                if(counter > max_sentences):
                    break
                counter += 1
                final_sentence += sentence + delimiter
        result = self.get_kw(final_sentence)
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
    files=[file1, file2,file3],
    directories=[]
)
req = DirectoryRequest(root=dir1)        

if __name__ == "__main__":
    kw_extractor = KWExtractor()
    result = kw_extractor.extract_kw(req.root)
for file_name, keywords in result:
    print(f"\n== FILE: {file_name} ==")
    for kw, score in keywords:
        print("Keyword:", kw, "Score:", score)

        



