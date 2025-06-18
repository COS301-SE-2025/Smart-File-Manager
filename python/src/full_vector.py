from create_kw_cluster import KWCluster
from sklearn.preprocessing import OneHotEncoder
import pandas as pd
from sklearn.preprocessing import MinMaxScaler
import numpy as np

class FullVector:
    def __init__(self):
        self.full_vector = []
        self.seen_filetypes = []

#File map [
#     {
#         "filename" : "name",
#         "keywords" : [],
#         "size_kb" : int,
#         "filetype" : "type",
#         # ADD
#         "full_vector": [["keywords"], ["filetype"], "size"]
#     },
#     ...
# ]
    def createFullVector(self, files, vocabKW, filetypes, sizes):
        for file in files:
            self.full_vector = []
            tfidf_vec = self.assignTF_IDF(file["keywords"],vocabKW)            
            encoded_filetype = self.oneHotEncoding(file["filetype"], filetypes)        
            scaler = MinMaxScaler() #############################      
            scaler.fit_transform(np.array(sizes))
            file["full_vector"] = tfidf_vec



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
    

    def oneHotEncoding(self, filetype, filetypes):
        df = pd.DataFrame(filetypes, columns=["filetype"])
        encoder = OneHotEncoder(sparse_output=False)
        one_hot_encoded = encoder.fit_transform(df[["filetype"]])
        one_hot_df = pd.DataFrame(one_hot_encoded, columns=encoder.get_feature_names_out(["filetype"]))
        # print(f"encoded data: \n{one_hot_df}")
        return one_hot_df.values.tolist()




