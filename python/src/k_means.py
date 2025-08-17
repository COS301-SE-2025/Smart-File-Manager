# the actual clustering that will make a directory to send back to go
import os

from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
import numpy as np
from math import ceil

from directory_builder import DirectoryCreator

from create_folder_name import FolderNameCreator

# Set random seed (I'm really not sure if this helps)
np.random.seed(42)

class KMeansCluster:
    def __init__(self, numClusters, max_depth, model, parent_folder):
        self.base_clusters = numClusters
        self.min_size = 2 # hardcoded for now # even numbers are good
        self.max_depth = max_depth
        self.folder_namer = FolderNameCreator(model)
        self.parent_folder = parent_folder

        self.locked_files = []

        self.kmeans = None


    def fit_kmeans(self, vectors, num_clusters):
        # Only reinitialize the model when necessary
        if self.kmeans is None or self.kmeans.n_clusters != num_clusters:
            self.kmeans = KMeans(
                n_clusters=num_clusters,
                init="k-means++",
                max_iter=500,
                tol=1e-5,
                random_state=42
            )
        self.kmeans.fit(vectors)
        return self.kmeans.labels_


    def cluster(self, vectors):
        return self.kmeans.fit_predict(vectors)

    def predict(self, points):
        preds = self.kmeans.predict(points)
        centers = np.round(self.kmeans.cluster_centers_, 4)
        return preds, centers
    
    def dirCluster(self,full_vecs,files):
        builder = DirectoryCreator(self.parent_folder,files) # instead of root it should be the parent folder
        self.remove_locked_files(files,full_vecs)
        print("Locked files ", len(self.locked_files))
        
        unlocked_dirs = self._recursive_clustering(full_vecs, files, 0, self.parent_folder, builder)

        locked_dirs = self.buildLockedDirs(self.locked_files, builder)


        root_dir = builder.merge(unlocked_dirs, locked_dirs)
        #print(root_dir)
        #self.printMetaData(root_dir)
        return root_dir 
        
    


    def _recursive_clustering(self,full_vecs,files, depth, dir_prefix, builder):

        # Quit if not enough folders  # Base condition: shallow depth or too few vectors
        if len(full_vecs) < self.min_size or depth > self.max_depth: # depth can be changed on init of kmeans
            return builder.buildDirectory(dir_prefix, files, []) 

        if depth > 0:
            # Assign directory name
            folder_name = self.folder_namer.generateFolderName(files)        
            dir_name = folder_name 

            if os.path.basename(dir_prefix) == folder_name:
                return builder.buildDirectory(dir_prefix,files,[])
        else:
            dir_name = dir_prefix
                    

        
        bias_factor = (1 / (depth + 1)) * 0.02
#        get optimal amount of clusters
        k = self.get_num_clusters(
                full_vecs, 
                k_min = self.min_size, 
                bias_factor = bias_factor,
                cluster_fraction = 5
                )
       # print("Best k values found: ", k)
        if k <= 1:
            return builder.buildDirectory(dir_name, files, [])

        
        # cluster and get labels
        labels = self.fit_kmeans(full_vecs, k)

        # label -> files
        label_to_entries = {}
        for i, label in enumerate(labels):
            label_to_entries.setdefault(label, []).append(files[i])
 

        subdirs = []
        retained_files = []


        for label, entries in label_to_entries.items():
            if len(entries) < self.min_size:
                retained_files.extend(entries)
                continue

            sub_vecs = [entry["full_vector"] for entry in entries]
            sub_dir = self._recursive_clustering(sub_vecs, entries, depth + 1, dir_name, builder)
            subdirs.append(sub_dir)

        return builder.buildDirectory(dir_name, retained_files, subdirs)
   

    def sigmoid(self, x):
        return 1 / (1 + np.exp(-x))


    def printDirectoryTree(self, directory, indent=""):

        for file in directory.files:
           # print(f"{file.name} ")
          #  print(f"{indent} - {file.original_path} ")
            print(f'"{file.new_path}",')

        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")

    def remove_locked_files(self, files, full_vecs):        
        # Iterate in reversed since we are modifying indices
        for i in reversed(range(len(files))):
            file = files[i]
            if file.get("is_locked", False):
                #print("Found locked file ", file["filename"])
                self.locked_files.append(file)
                del files[i]
                del full_vecs[i]

    def buildLockedDirs(self, files, builder):
        def add_to_tree(tree, parts, file):
            current = tree
            current = current.setdefault(os.path.join(*parts), {})
            current.setdefault("_files", []).append(file)

        def build_dirs_from_tree(tree):
            dirs = []
            for name, subtree in tree.items():
                if name == "_files":
                    continue
                subdirs = build_dirs_from_tree(subtree) 
                files = subtree.get("_files", [])
                dirs.append(builder.buildDirectory(name, files, subdirs))
            return dirs

        dir_tree = {}
        for f in files:
            full_path = os.path.normpath(f["original_path"])
            path_parts = full_path.split(os.sep)
            try:
                parent_index = path_parts.index(self.parent_folder)
                relative_parts = path_parts[parent_index + 1 : -1]
                if not relative_parts:
                    relative_parts = ["(root)"]
            except (ValueError, IndexError):
                relative_parts = ["Unkown"]

            add_to_tree(dir_tree, relative_parts,f)

        return builder.buildDirectory(self.parent_folder, [], build_dirs_from_tree(dir_tree))

    def get_num_clusters(self, X, k_min = 2, k_max = None, random_state = 42, bias_factor = 0.01, bad_threshold=0.2, cluster_fraction = 4):
        """
        Determine amount of clusters using silhouette_score and elbow method
        X - features
        k_min - min clusters set on init
        k_max - max clusters set on init
        random_state - random num for reproducibility

        returns:
            dict with best_k, silhouette_score, inertias
        """
        num_samples = len(X)
        #print("Num samples: ", num_samples)
        if num_samples < 2:
            return 1

        if k_max is None:
            k_max = min(num_samples, max(k_min, ceil(num_samples // cluster_fraction)))

        silhouette_scores = []
        valid_k = []
        k_values = range(k_min, min(k_max, num_samples)+1)
        #print("Range of val: ", k_values)

        for k in k_values:
            if k<= 1 or k>= num_samples:
            #    print("Continue")
                continue
            try:
                kmeans = KMeans(n_clusters=k, random_state=42, n_init="auto", init="k-means++", tol=1e-5)
                labels = kmeans.fit_predict(X)

                score = silhouette_score(X, labels)
                score += k * bias_factor

                silhouette_scores.append(score)
                valid_k.append(k)

            except Exception:
                pass

        if not silhouette_scores:
           # print("Sil scores empty")
            return 1


        best_index = int(np.argmax(silhouette_scores))
        best_k = valid_k[best_index]

        # if the best score is very bad dont split
        best_score = silhouette_scores[best_index] - best_k * bias_factor
        if best_score < bad_threshold:
            return 1

        return best_k

