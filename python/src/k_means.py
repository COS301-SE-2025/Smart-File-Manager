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


    def cluster(self,files):
        self.kmeans.fit(files)
        return self.kmeans.labels_

    def predict(self, points):
        predictions = self.kmeans.predict(points)
        centers = self.kmeans.cluster_centers_
        centers_rounded = np.round(centers, 4) # rounded to get mostly matching        
        return predictions, centers_rounded
    
    def dirCluster(self,full_vecs,files):
        # labels = self.recDirCluster(full_vecs,files,self.n_clusters)   
        # #print(labels)    
        # self.printDirectoryTree(labels)    
        # return labels
        builder = DirectoryCreator("Root",files)
        root_dir = self.recDirCluster(full_vecs, files, self.n_clusters, "Directory", builder)
        self.printDirectoryTree(root_dir)
        return root_dir
        
    def recDirCluster(self,full_vecs,files, depth, dir_prefix, builder):
        dir_name = f"{dir_prefix}_{depth}"

          
        if len(full_vecs) <= depth:
            return builder.buildDirectory(dir_name, files, []) 

        self.kmeans = KMeans(
            n_clusters=depth,
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
 

    def labelToFileName(self, labels, files):
        label_to_filenames = defaultdict(list)
        for index, file in enumerate(files):
            label = labels[index]
            filename = file["filename"]
            fv = file["full_vector"]
            label_to_filenames[label].append((filename,fv))
        
     
        return label_to_filenames        


    def printTree(self,tree, prefix=""):
        if "files" in tree:
            for file_entry in tree["files"]:
                if isinstance(file_entry, dict):
                    print(f"{prefix}  - {file_entry['filename']}")
                else:
                    print(f"{prefix}  - INVALID FILE ENTRY: {file_entry}")
        for child in tree.get("children", []):
            print(f"{prefix}{child['label']}:")
            self.printTree(child, prefix + "  ")

    def printDirectoryTree(self, directory, indent=""):
        print(f"{indent}{directory.name}/")
        for file in directory.files:
            print(f"{indent}  - {file.new_path}")
        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")