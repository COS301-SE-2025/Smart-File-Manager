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
        self.minSize = 2 # hardcoded for now # even numbers are good


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
        root_dir = self.recDirCluster(full_vecs, files, 0, "Directory", builder)
        self.printDirectoryTree(root_dir)
        return root_dir
        
    def recDirCluster(self,full_vecs,files, depth, dir_prefix, builder):
        # Assign directory name
        dir_name = f"{dir_prefix}_{depth}"

        # Quit if not enough folders
        if len(full_vecs) < self.minSize or depth > 30:
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

        # idk lost in the sauce
        # if any(len(entries) <= self.minSize for entries in label_to_entries.values()):
        #     return builder.buildDirectory(dir_name, files, [])
 

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
                sub_dir = self.recDirCluster(sub_vecs,entries,depth+1,f"{dir_name}_{label}", builder)
                subdirs.append(sub_dir)
            # Not quite enough files to recluster so keep them together   
            else:
                retained_files.extend(entries)

        return builder.buildDirectory(dir_name, retained_files, subdirs)
   



    def printDirectoryTree(self, directory, indent=""):
        print(f"{indent}{directory.name}/")
        for file in directory.files:
            print(f"{indent}  - {file.name}")
        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")