# the actual clustering that will make a directory to send back to go
import os

from sklearn.cluster import KMeans
import numpy as np

from directory_builder import DirectoryCreator
from collections import defaultdict

from create_folder_name import FolderNameCreator

class KMeansCluster:
    def __init__(self, numClusters, depth, model, parent_folder):
        self.kmeans = KMeans(
            n_clusters=numClusters,
            n_init=50, 
            init="k-means++",
            max_iter=500,
            tol=1e-5           
            )
        self.n_clusters = numClusters
        self.min_size = 2 # hardcoded for now # even numbers are good
        self.max_depth = depth
        self.folder_namer = FolderNameCreator(model)
        self.parent_folder = parent_folder


    def cluster(self,files):
        self.kmeans.fit(files)
        return self.kmeans.labels_

    def predict(self, points):
        predictions = self.kmeans.predict(points)
        centers = self.kmeans.cluster_centers_
        centers_rounded = np.round(centers, 4) # rounded to get mostly matching        
        return predictions, centers_rounded
    
    def dirCluster(self,full_vecs,files):
        builder = DirectoryCreator(self.parent_folder,files) # instead of root it should be the parent folder
        root_dir = self.recDirCluster(full_vecs, files, 0, self.folder_namer.generateFolderName(files), builder)
        #print(root_dir)
        #self.printMetaData(root_dir)
        return root_dir
        
    


    def recDirCluster(self,full_vecs,files, depth, dir_prefix, builder):

        if depth > 0:
            # Assign directory name
            folder_name = self.folder_namer.generateFolderName(files)        
            dir_name = os.path.join(dir_prefix, folder_name)

            if folder_name in dir_prefix:
                return builder.buildDirectory(dir_prefix,files,[])
        else:
            dir_name = dir_prefix
                    

        # Quit if not enough folders
        if len(full_vecs) < self.min_size or depth > self.max_depth: # depth can be changed on init of kmeans
            return builder.buildDirectory(dir_name, files, []) 
        
        # if there are too many clusters reduce it 
        if len(full_vecs) < self.n_clusters:
            self.n_clusters = len(full_vecs)

        # Min clusters (x leaves, x dirs per level)
        if self.n_clusters <= self.min_size:
            self.n_clusters = self.min_size
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
            elif len(entries) >= self.min_size:                  
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

