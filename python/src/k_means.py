# the actual clustering that will make a directory to send back to go
import os

from sklearn.cluster import KMeans
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np

from directory_builder import DirectoryCreator
from collections import defaultdict

class KMeansCluster:
    def __init__(self, numClusters, depth, model):
        self.kmeans = KMeans(
            n_clusters=numClusters,
            n_init=50, 
            init="k-means++",
            max_iter=500,
            tol=1e-5           
            )
        self.n_clusters = numClusters
        self.minSize = 2 # hardcoded for now # even numbers are good
        self.maxDepth = depth
        self.maxKeywords = 10 # for folder name creation
        # temporary
        self.model = model


    def cluster(self,files):
        self.kmeans.fit(files)
        return self.kmeans.labels_

    def predict(self, points):
        predictions = self.kmeans.predict(points)
        centers = self.kmeans.cluster_centers_
        centers_rounded = np.round(centers, 4) # rounded to get mostly matching        
        return predictions, centers_rounded
    
    def dirCluster(self,full_vecs,files):
        builder = DirectoryCreator("Root",files) # instead of root it should be the parent folder
        root_dir = self.recDirCluster(full_vecs, files, 0, "Directory", builder)
        #print(root_dir)
        #self.printMetaData(root_dir)
        return root_dir
        
    def generateFolderName(self, files, max_keywords = 10):
        all_keywords = []
        keyword_scores = {}

        for file in files:
            for kw,score in file["keywords"]:
                if kw not in keyword_scores:
                    keyword_scores[kw] = 0
                keyword_scores[kw] += 1.0 / (1.0 + score)

        # Top weighted keywrods
        sorted_keywords = sorted(keyword_scores.items(), key=lambda x: x[1], reverse=True)
        top_keywords = [kw for kw, _ in sorted_keywords[:max_keywords]]

        if not top_keywords: # probably an image. Can use a suffix if there are multiple of these folders. Or we can use their names in a sentence transformer
            return "Misc"
        
        # Encode and find a centroid
        embeddings = self.model.encode(top_keywords)
        centroid = np.mean(embeddings, axis=0,keepdims=True)

        # Best keyword
        sims = cosine_similarity(centroid, embeddings)[0]
        best_idx = int(np.argmax(sims))
        folder_keyword = top_keywords[best_idx]
        folder_keyword.replace(".","")

        return folder_keyword.replace(" ", "_")


    def recDirCluster(self,full_vecs,files, depth, dir_prefix, builder):
        # Assign directory name
        folder_name = self.generateFolderName(files)        
        #dir_name = f"{dir_prefix}/{folder_name}"
        if folder_name.lower() not in [part.lower() for part in dir_prefix.split(os.sep)]:
            dir_name = os.path.join(dir_prefix, folder_name)
        else:
            dir_name = dir_prefix
        

            

        # Quit if not enough folders
        if len(full_vecs) < self.minSize or depth > self.maxDepth: # depth can be changed on init of kmeans
            return builder.buildDirectory(dir_name, files, []) 
        
        # if there are too many clusters reduce it 
        if len(full_vecs) < self.n_clusters:
            self.n_clusters = len(full_vecs)

        # Min clusters (x leaves, x dirs per level)
        if self.n_clusters <= self.minSize:
            self.n_clusters = self.minSize
        else:
            self.n_clusters -= 1

        # create the clustering
        self.kmeans = KMeans(
            n_clusters=self.n_clusters,
            random_state=42,
            n_init="auto",
        )

        # cluster and get labels
        labels = self.cluster(full_vecs)

        # label -> files
        label_to_entries = {}
        for i, label in enumerate(labels):
            label_to_entries.setdefault(label, []).append(files[i])
 

        subdirs = []
        retained_files = []


        for label, entries in label_to_entries.items():
            # if a label has one entry then clustering is pretty good
            # To avoid having files in leaves we return all of the files used in this clustering
            # -> can defnitely backfire but lets hope the clustering is goated
            if len(entries) <= 1:
                # Flavour 1 (millions of dirs)
                # sub_vecs = [e["full_vector"] for e in entries]
                # sub_dir = self.recDirCluster(sub_vecs,entries,depth,f"{dir_name}_{label}", builder)
                # subdirs.append(sub_dir)
                # Flavour 2 (More files per dir)
                return builder.buildDirectory(dir_name, files, [])                
            # Good number of entries, atleast minsize so recursively check (if exactly minsize it will make a dir of these two folders)
            elif len(entries) >= self.minSize:                  
                sub_vecs = [e["full_vector"] for e in entries]
                sub_dir = self.recDirCluster(sub_vecs,entries,depth+1,f"{dir_name}", builder)
                subdirs.append(sub_dir)
            # Not quite enough files to recluster so keep them together   
            else:
                retained_files.extend(entries)

        return builder.buildDirectory(dir_name, retained_files, subdirs)
   



    def printDirectoryTree(self, directory, indent=""):
       # print(f"{indent}{directory.name}/")
        for file in directory.files:
           # print(f"{file.name} ")
          #  print(f"{indent} - {file.original_path} ")
            print(f'"{file.new_path}",')
        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")

