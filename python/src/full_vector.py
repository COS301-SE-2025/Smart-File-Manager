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
    def createFullVector(self, files, vocabKW):
        sizes = [file["size_bytes"] for file in files]
        filetypes = sorted(set(file["filetype"] for file in files))
        scaler = MinMaxScaler()   
        norm_sizes = scaler.fit_transform(np.array(sizes).reshape(-1,1))
        filetype_to_onehot = {ft:self.oneHotEncoding(ft,filetypes) for ft in filetypes}
        for idx, file in enumerate(files):
            self.full_vector = []
            tfidf_vec = self.assignTF_IDF(file["keywords"],vocabKW)            
            encoded_filetype = filetype_to_onehot[file["filetype"]]     
            normalized_size = norm_sizes[idx].tolist()
            if isinstance(tfidf_vec,np.ndarray):
                tfidf_vec = tfidf_vec.flatten().tolist()

            file["full_vector"] = tfidf_vec + encoded_filetype+ normalized_size



    #helper function
    def assignTF_IDF(self, result, vocabKW):

        kwclust = KWCluster()       

        return kwclust.createCluster(result,vocabKW)    

    

    def oneHotEncoding(self, filetype, filetypes):
        df = pd.DataFrame(filetypes, columns=["filetype"])
        encoder = OneHotEncoder(sparse_output=False)
        one_hot_encoded = encoder.fit_transform(df[["filetype"]])

        # Find the index of the requested filetype
        index = df[df["filetype"] == filetype].index[0]
        return one_hot_encoded[index].tolist()




