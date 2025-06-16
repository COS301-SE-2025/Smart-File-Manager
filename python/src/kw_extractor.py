import time
import docx
import math
from yake import KeywordExtractor
from pypdf import PdfReader
from message_structure_pb2 import Directory, DirectoryRequest, File, MetadataEntry, Tag

# TEMP
from vocabulary import Vocabulary
from create_kw_cluster import KWCluster
##

# Receive data as D
class KWExtractor:
    #Yake instance
    def __init__(self):
        self.yake_extractor = KeywordExtractor(lan="en", n=1)

    #Main extractor function
    def extract_kw(self, input):
        result = []
        for file in input.files:    
            file_name = f"{file.original_path}"
            mime_type = next((entry.value for entry in file.metadata if entry.key == "mime_type"), None)
            result += self.open_file(file_name,mime_type, 1) #Seconds based time limit
        return self.list_to_map(result, 100)   
    
    #Make the result into a sorted map with max keywords
    def list_to_map(self, result, max_keywords):
        return_map = {}

        for file_name, keywords in result:
            #descending order
            sorted_keywords = sorted(keywords, key=lambda x: x[1], reverse=True)
            top_keywords = sorted_keywords[:max_keywords]
            # Extract scores and normalize them (L2 norm)
            scores = [score for _, score in top_keywords]
            norm = math.sqrt(sum(s**2 for s in scores)) or 1  # avoid div by zero

            normalized_keywords = [(kw, score / norm) for kw, score in top_keywords]

            return_map[file_name] = normalized_keywords

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

    vocab = Vocabulary()
    vocabKW = vocab.createVocab(result)
    print(vocabKW)

    kwclust = KWCluster()
    clusterVec = []
    # for kw in vocabKW:
    #     print(kw)
    for filename, keywords in result.items():
        print(f"\n== FILE: {filename} ==")
        # for kw, score in keywords:
        #     print("Keyword: ", kw, "\tScore: ", score)
        clusterVec.append(kwclust.createCluster(keywords,vocabKW))
        print(clusterVec)
    for vec in clusterVec:
        print("\n")
        print(vec)
        print("\n")
    

    

        



