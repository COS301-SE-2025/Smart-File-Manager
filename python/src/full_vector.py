from create_kw_cluster import KWCluster
from sklearn.preprocessing import OneHotEncoder
import pandas as pd
from sklearn.preprocessing import MinMaxScaler
import numpy as np
from vocabulary import Vocabulary
from message_structure_pb2 import File

class FullVector:
    def __init__(self):
        self.vocab = Vocabulary()
        
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
    def createFullVector(self, files : File) -> None:
        filetypes, sizes, vocabKW = self.assignSizeFileTypeKeywords(files)
 
        scaler = MinMaxScaler()   
        norm_sizes = scaler.fit_transform(np.array(sizes).reshape(-1,1)).tolist()
        filetype_to_onehot = {ft:self.oneHotEncoding(ft,filetypes) for ft in filetypes}
        for idx, file in enumerate(files):
            tfidf_vec = self.assignTF_IDF(file["keywords"],vocabKW)    
            encoded_filetype = filetype_to_onehot[file["file_extension"]]
            normalized_size = norm_sizes[idx]
            if isinstance(tfidf_vec,np.ndarray):
                tfidf_vec = tfidf_vec.flatten().tolist()

            file["full_vector"] = [tfidf_vec,encoded_filetype,normalized_size]            



    #helper function
    def assignTF_IDF(self, result : list[tuple], vocabKW : list[str]) -> list[float]:

        kwclust = KWCluster()       

        return kwclust.createCluster(result,vocabKW)    

    

    def oneHotEncoding(self, filetype, filetypes):
        df = pd.DataFrame(filetypes, columns=["file_extension"])
        encoder = OneHotEncoder(sparse_output=False)
        one_hot_encoded = encoder.fit_transform(df[["file_extension"]])
        # Get the index of the filetype in the original list
        try:
            index = df[df["file_extension"] == filetype].index[0]
        except IndexError:
            raise ValueError(f"Filetype '{filetype}' not found in the list.")
        
        return one_hot_encoded[index].tolist()

    def assignSizeFileTypeKeywords(self, files):
        sizes = []
        filetypes = set()
        keywords = []

        for file in files:
            sizes.append(file["size_bytes"])
            filetypes.add(file["file_extension"])
            keywords.append(file["keywords"])

        vocabKW = self.vocab.createVocab(keywords)
        
        return sorted(filetypes), sizes, vocabKW


if __name__ == "__main__":
    files = [{'filename': 'todo.docx', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/todo.docx', 'file_extension': '.docx', 'size_bytes': 16581, 'created': '2025-06-19T08:39:48.704454', 'modified': '2025-05-26T11:13:50.652646', 'mime_type': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'accessed': '2025-06-19T08:39:48.704454', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 25332747904318331, 'keywords': [('to-do', 0.1761190379942198), ('list', 0.1761190379942198), ('groceries.Call', 0.1761190379942198), ('Dad.Have', 0.1761190379942198), ('meeting', 0.1761190379942198), ('boss.Finish', 0.1761190379942198), ('assignments.Blah', 0.1761190379942198), ('blah', 0.1761190379942198), ('blah.Things', 0.1761190379942198), ('forgot', 0.1761190379942198)]}, {'filename': 'myImg.jpg', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/UniFiles/myImg.jpg', 'file_extension': '.jpg', 'size_bytes': 235226, 'created': '2025-06-19T08:39:48.636931', 'modified': '2025-05-26T11:13:50.428722', 'mime_type': 'image/jpeg', 'accessed': '2025-06-19T08:39:48.636931', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 17169973579585759, 'keywords': []}, {'filename': 'myPdf.pdf', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/UniFiles/myPdf.pdf', 'file_extension': '.pdf', 'size_bytes': 1722297, 'created': '2025-06-19T08:39:48.652774', 'modified': '2025-05-26T11:13:50.630410', 'mime_type': 'application/pdf', 'accessed': '2025-06-19T08:39:48.652774', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 7036874418003109, 'keywords': [('Advanced', 0.14985554132463938), ('Real-time', 0.14985554132463938), ('Dropbox', 0.1453011207686066), ('search', 0.13775170326151098), ('customization', 0.13588457081573124), ('system', 0.132439301673504), ('Manager', 0.12174767379894379), ('AI-powered', 0.12051367106078721), ('user', 0.11645295367216063), ('folders', 0.09314785652067586)]}, {'filename': 'myVideo.webm', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/PersonalFiles/myVideo.webm', 'file_extension': '.webm', 'size_bytes': 1318606, 'created': '2025-06-19T08:39:48.581326', 'modified': '2025-05-26T11:13:50.320314', 'mime_type': 'video/webm', 'accessed': '2025-06-19T08:39:48.581326', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 4503599627604377, 'keywords': []}, {'filename': 'thumbbig-708440.webp', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/PersonalFiles/thumbbig-708440.webp', 'file_extension': '.webp', 'size_bytes': 65346, 'created': '2025-06-19T08:39:48.621100', 'modified': '2025-05-26T11:13:50.365673', 'mime_type': 'image/webp', 'accessed': '2025-06-19T08:39:48.621100', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 10133099161817879, 'keywords': []}, {'filename': 'holiday.JPG', 'absolute_path': '/mnt/c/Users/Phili/OneDrive/PhilippDup/University/Year_3/COS301/CapstoneSFM/Code/Smart-File-Manager/python/testing/test_files_2/UserFiles/PersonalFiles/holiday.JPG', 'file_extension': '.jpg', 'size_bytes': 1735302, 'created': '2025-06-19T08:39:48.527258', 'modified': '2025-05-26T11:13:50.148073', 'mime_type': 'image/jpeg', 'accessed': '2025-06-19T08:39:48.527258', 'owner_uid': 1000, 'owner_gid': 1000, 'mode': '0o100777', 'inode': 3940649674183008, 'keywords': []}]
    fv = FullVector()
    fv.createFullVector(files)
    for f in files:
        print(f, "\n")
