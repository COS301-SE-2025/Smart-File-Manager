from create_kw_cluster import KWCluster
from sklearn.preprocessing import OneHotEncoder
import pandas as pd

class FullVector:
    def __init__(self):
        self.full_vector = []
        self.tfidf_vector = []
        self.filetypes = []
        self.sizes = []
        self.encoded_filetypes = []

    def createFullVector(self, keywords, vocab, root):
        self.tfidf_vec = self.assignTF_IDF(keywords,vocab)
        self.filetypes, self.sizes = self.assignFileTypeAndSize(root)
        self.encoded_filetypes = self.oneHotEncoding(self.filetypes)

        # print(len(self.tfidf_vec))
        # print(len(self.sizes))
        # print(len(self.encoded_filetypes))
        # Catch bad data
        if len(self.tfidf_vec) != len(self.encoded_filetypes):
            return []
        if len(self.encoded_filetypes) != len(self.sizes):
            return []

        # Loop through each files
        for i in range(len(self.tfidf_vec)):
            self.full_vector.append((self.tfidf_vec[i], self.encoded_filetypes[i], self.sizes[i]))
        
        for entry in self.full_vector:
            print("\n")
            print(entry)

        # print(self.encoded_filetypes)
        # assign_IDF         
        # one hot encoding 
        # 
        # for all files:
        #   self.full.vec.append(clustervec[index], filetype[index], size[index])
        # for vec in self.tfidf_vec:
        #     print("\n")
        #     print(vec)
        # for vec in self.filetypes:
        #     print("\n")
        #     print(vec)
        # for vec in self.sizes:
        #     print("\n")
        #     print(vec)
    #helper function
    def assignTF_IDF(self, result, vocabKW):
        # process request root
        # Return vector of TF-IDF
        # Compare direct words, similiar words (cat - cats), contained words (AI - AI-driven)
        # Assign TF-IDF score
        # returns:
        # [
        # [0,0,0,0.51,0,0.1],
        # [...]
        # ]
        kwclust = KWCluster()
        tfidf = []
        # for kw in vocabKW:
        #     print(kw)
        
        for filename, keywords in result.items():
            kwcluster = kwclust.createCluster(keywords,vocabKW)
            tfidf.append(kwcluster)

        return tfidf
    
    def assignFileTypeAndSize(self, root):
        types = []
        sizes = []
        for file in root.files:   
            file_type = next((entry.value for entry in file.metadata if entry.key == "file_extension"), None)
            file_size = next((int(entry.value) for entry in file.metadata if entry.key == "size_bytes"), None)
            types.append(file_type)
            sizes.append(file_size)
        return types, sizes
    
    def oneHotEncoding(self, filetypes):
        df = pd.DataFrame(filetypes, columns=["filetype"])

        encoder = OneHotEncoder(sparse_output=False)

        one_hot_encoded = encoder.fit_transform(df[["filetype"]])

        one_hot_df = pd.DataFrame(one_hot_encoded, columns=encoder.get_feature_names_out(["filetype"]))

        # print(f"encoded data: \n{one_hot_df}")

        return one_hot_df.values.tolist()
# scrape metadata -> modify root -> return none
# extract kw -> modify none -> return Map<filename, keywords(kw, tf idf)>
# create vocab -> modify none -> return list(vocabulary)
# full_vector: 
# assign tf idf -> modify none -> return list(tf idf)
# assign file type -> modify none -> return filetypes (redundant but this is my bandaid)
# assign size -> modify none -> return sizes (redundant aswell)
# modify full vec to combine [tfidf, type, size]



