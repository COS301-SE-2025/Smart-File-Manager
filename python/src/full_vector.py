from create_kw_cluster import KWCluster

class FullVector:
    def __init__(self):
        self.full_vector = []
        self.tfidf_vector = []

    def createFullVector(self, keywords, vocab, filetypes, sizes):
        filenames = list(keywords.keys())
        self.tfidf_vec = self.assignTF_IDF(keywords,vocab)
        self.filetypes = self.assignFileType(filetypes)
        # assign_IDF         
        # one hot encoding 
        # 
        # for all files:
        #   self.full.vec.append(clustervec[index], filetype[index], size[index])
        for vec in self.full_vector:
            print(vec)
            print("\n")
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
    
    def assignFileType(self, filetypes):
        full_vec = []
        for i in range(len(self.tfidf_vector)):
            full_vec = self.tfidf_vector[i] + filetypes[i]
        return full_vec
        
# scrape metadata -> modify root -> return none
# extract kw -> modify none -> return Map<filename, keywords(kw, tf idf)>
# create vocab -> modify none -> return list(vocabulary)
# full_vector: 
# assign tf idf -> modify none -> return list(tf idf)
# assign file type -> modify none -> return filetypes (redundant but this is my bandaid)
# assign size -> modify none -> return sizes (redundant aswell)
# modify full vec to combine [tfidf, type, size]



