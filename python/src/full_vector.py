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

        # Catch bad data
        if len(self.tfidf_vec) != len(self.encoded_filetypes):
            return []
        if len(self.encoded_filetypes) != len(self.sizes):
            return []

        # Loop through each files
        for i in range(len(self.tfidf_vec)):
            self.full_vector.append((self.tfidf_vec[i], self.encoded_filetypes[i], self.sizes[i]))
        
        # for entry in self.full_vector:
        #     print("\n")
        #     print(entry)

        combined_points = []
        # combine into a single tuple
        for vec in self.full_vector:
            combined_points.append(list(vec[0]) + list(vec[1]) + [vec[2]])

        return combined_points

    #helper function
    def assignTF_IDF(self, result, vocabKW):

        kwclust = KWCluster()
        tfidf = []

        
        for filename, keywords in result.items():
            kwcluster = kwclust.createCluster(keywords,vocabKW)
            tfidf.append(kwcluster)

        return tfidf
    
    def assignFileTypeAndSize(self, root):
        types = []
        sizes = []
        for file in root.files:   
            file_type = next((entry.value for entry in file.metadata if entry.key == "file_extension"), None)
            # try catch incase conversion fails? -> add fixed data then?
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




