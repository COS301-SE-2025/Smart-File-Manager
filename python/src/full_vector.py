from create_kw_cluster import KWCluster

class FullVector:
    def __init__(self):
        self.full_vector = []

    def createFullVector(self):
        pass
        # assign_IDF         
        # one hot encoding 
        # 
        # for all files:
        #   self.full.vec.append(clustervec[index], filetype[index], size[index])
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
        clusterVec = []
        # for kw in vocabKW:
        #     print(kw)
        for filename, keywords in result.items():
            kwcluster = kwclust.createCluster(keywords,vocabKW)
            clusterVec.append(kwcluster)
            #print(clusterVec)

        return clusterVec


