import os
import numpy as np
from collections import defaultdict
from sklearn.cluster import KMeans

from directory_builder import DirectoryCreator
from create_folder_name import FolderNameCreator


class KMeansCluster:
    def __init__(self, num_clusters, max_depth, model, parent_folder):
        self.base_clusters = num_clusters
        self.max_depth = max_depth
        self.min_size = 2  # Minimum cluster size
        self.parent_folder = parent_folder

        self.kmeans = None
        self.folder_namer = FolderNameCreator(model)

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

    def predict(self, points):
        preds = self.kmeans.predict(points)
        centers = np.round(self.kmeans.cluster_centers_, 4)
        return preds, centers

        if depth > 0:
            # Assign directory name
            folder_name = self.folder_namer.generateFolderName(files)        
            dir_name = folder_name

    def cluster(self, vectors):
        return self.kmeans.fit_predict(vectors)

    def dirCluster(self, full_vecs, files):
        builder = DirectoryCreator(self.parent_folder, files)
        return self._recursive_clustering(full_vecs, files, depth=0, dir_prefix=self.parent_folder, builder=builder)


    def _recursive_clustering(self, vectors, files, depth, dir_prefix, builder):
        # Base condition: shallow depth or too few vectors
        if len(vectors) < self.min_size or depth > self.max_depth:
            return builder.buildDirectory(dir_prefix, files, [])

        # Assign new directory name at deeper levels
        if depth > 0:
            folder_name = self.folder_namer.generateFolderName(files)
            if folder_name in dir_prefix:
                return builder.buildDirectory(dir_prefix, files, [])
            dir_name = folder_name
        else:
            dir_name = dir_prefix

        # Determine the number of clusters to use
        n_clusters = min(len(vectors), self.base_clusters)
        if n_clusters <= self.min_size:
            return builder.buildDirectory(dir_name, files, [])

        labels = self.fit_kmeans(vectors, n_clusters)

        # Organize files by cluster label
        label_to_entries = defaultdict(list)
        for i, label in enumerate(labels):
            label_to_entries[label].append(files[i])

        subdirs = []
        retained_files = []

        for entries in label_to_entries.values():
            if len(entries) < self.min_size:
                # Keep small clusters in the current directory
                retained_files.extend(entries)
                continue

            sub_vecs = [entry["full_vector"] for entry in entries]
            sub_dir = self._recursive_clustering(sub_vecs, entries, depth + 1, dir_name, builder)
            subdirs.append(sub_dir)

        return builder.buildDirectory(dir_name, retained_files, subdirs)


    def printDirectoryTree(self, directory, indent=""):

        for file in directory.files:
            print(f'"{file.new_path}",')

        for subdir in directory.directories:
            self.printDirectoryTree(subdir, indent + "  ")

