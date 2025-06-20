# the actual clustering that will make a directory to send back to go

from sklearn.cluster import KMeans
import numpy as np

from directory_builder import DirectoryCreator
from collections import defaultdict

class KMeansCluster:
    def __init__(self, numClusters):
        self.kmeans = KMeans(
            n_clusters=numClusters,
            random_state=42,
            n_init="auto",            
            )
        self.n_clusters = numClusters
        self.minSize = 2 # hardcoded for now


    def cluster(self,files):
        self.kmeans.fit(files)
        return self.kmeans.labels_

    def predict(self, points):
        predictions = self.kmeans.predict(points)
        centers = self.kmeans.cluster_centers_
        centers_rounded = np.round(centers, 4) # rounded to get mostly matching        
        return predictions, centers_rounded
    
    def dirCluster(self,full_vecs,files):
        builder = DirectoryCreator("Root",files)
        root_dir = self.recDirCluster(full_vecs, files, self.n_clusters, "Directory", builder)
        self.printDirectoryTree(root_dir)
        return root_dir
        
    def recDirCluster(self,full_vecs,files, depth, dir_prefix, builder):
        # Assign directory name
        dir_name = f"{dir_prefix}_{depth}"

        # Quit if not enough folders
        if len(full_vecs) <= self.minSize:
            return builder.buildDirectory(dir_name, files, []) 
        
        if self.n_clusters <= self.minSize:
            self.n_clusters = self.minSize
        else:
            self.n_clusters -= 1

        self.kmeans = KMeans(
            n_clusters=self.n_clusters,
            random_state=42,
            n_init="auto",
        )

        
        labels = self.cluster(full_vecs)

        
        label_to_entries = {}
        for i, label in enumerate(labels):
            label_to_entries.setdefault(label, []).append(files[i])

        subdirs = []
        retained_files = []
   

        for label, entries in label_to_entries.items():
            if len(entries) > 2:                  
                sub_vecs = [e["full_vector"] for e in entries]
                sub_dir = self.recDirCluster(sub_vecs,entries,depth+1,f"{dir_name}_{label}", builder)
                subdirs.append(sub_dir)
                
            else:
                retained_files.extend(entries)

        return builder.buildDirectory(dir_name, retained_files, subdirs)
   



    def printDirectoryTree(self, directory, indent=""):
        print(f"{indent}{directory.name}/")
        for file in directory.files:
            print(f"{indent}  - {file.name}")
        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")