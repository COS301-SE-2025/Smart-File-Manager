#RAKE seems like a decent option
#YAKE could be better for simple context based comparisons down the line

#pip install yake

#GeeksForGeeks example
# Create a KeywordExtractor instance

# # Text from which keywords will be extracted
# text = "YAKE (Yet Another Keyword Extractor) is a Python library for extracting keywords from text."

# # Extract keywords from the text
# keywords = kw_extractor.extract_keywords(text)

# # Print the extracted keywords and their scores
# for kw in keywords:
#     print("Keyword:", kw[0], "Score:", kw[1])
from yake import KeywordExtractor
from message_structure_pb2 import Directory, DirectoryRequest, File, MetadataEntry, Tag

# Receive data as D
class KeywordExtractor:
    def extract_kw(self, input):
        result = ""
        for file in input.files:    
            file_name = f"{file.original_path}"       
            self.open_file(file_name, 10)
        return result
    
    def get_kw(self, sentence):
        print(sentence)

    def split_by_delimiter(self, file_name, delimiter):
        buffer = ''
        with open(file_name, 'r', encoding='utf-8') as file:
            for line in file:
                buffer += line
                while '.' in buffer:
                    sentence, buffer = buffer.split('.', 1)
                    yield sentence.strip() + '.'
        if buffer.strip():
            yield buffer.strip()

    def open_file(self, file_name, maxSentences):
        counter = 0        
        for sentence in self.split_by_delimiter(file_name, '.'):
            if(counter > maxSentences):
                break
            self.get_kw(sentence)
            counter += 1
            

    

tag1 = Tag(name="ImFixed")
meta1 = MetadataEntry(key="author", value="johnny")

file1 = File(
    name="gopdoc.pdf",
    original_path="python/src/testFile.txt",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1]
)

dir1 = Directory(
    name="useless_files",
    path="/usr/trash",
    files=[file1],
    directories=[]
)
req = DirectoryRequest(root=dir1)        

if __name__ == "__main__":
    kw_extractor = KeywordExtractor()
    result = kw_extractor.extract_kw(req.root)
    print(result)

        



