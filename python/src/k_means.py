# the actual clustering that will make a directory to send back to go
from collections import Counter
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
        
        unlocked_dirs = self._recursive_clustering(full_vecs, files, 0, "", builder)

        locked_dirs = self.buildLockedDirs(self.locked_files, builder)


        root_dir = builder.merge(unlocked_dirs, locked_dirs)
        return root_dir 
        

    def _recursive_clustering(self, full_vecs, files, depth, rel_path, builder):
        """
        rel_path: relative path from the root ("" at root).
        """
        # Base case: Directory too small or directory too deep
        if len(full_vecs) < self.min_size or depth > self.max_depth:
            return builder.buildDirectory(rel_path, files, [])

        # Find K (number of clusters) using elbow method
        bias_factor = (1 / (depth + 1)) * 0.02
        k = self.get_num_clusters(
            full_vecs,
            k_min=self.min_size,
            bias_factor=bias_factor,
            cluster_fraction=5
        )
        if k <= 1:
            return builder.buildDirectory(rel_path, files, [])

        # Cluster
        labels = self.fit_kmeans(full_vecs, k)

        # Group files by label (carry vectors along)
        groups = {}
        for i, label in enumerate(labels):
            groups.setdefault(label, []).append((files[i], full_vecs[i]))

        # Separate valid clusters vs small ones (retained in parent)
        valid_groups = []
        retained_files = []
        for entries in groups.values():
            if len(entries) < self.min_size:
                retained_files.extend([f for f, _ in entries])
            else:
                valid_groups.append(entries)

        # Stop splitting for clusters that only have 1 cluster (so we don't get group/group/etc...)
        if len(valid_groups) <= 1:
            return builder.buildDirectory(rel_path, files, [])

        # Build children
        parent_name = os.path.basename(rel_path) if rel_path else self.parent_folder
        sibling_counts = {}
        subdirs = []

        for entries in valid_groups:
            child_files = [f for f, _ in entries]
            child_vecs = [v for _, v in entries]

            # Try to generateActualName
            proposed = (self.folder_namer.generateFolderName(child_files) or "").strip()

            # Avoid duplicates
            child_name = self._normalize_child_name(proposed, parent_name, child_files)

            # If name would duplicate parent (or ends up empty) -> skip splitting this cluster and keep files at the current level.
            if child_name is None:
                retained_files.extend(child_files)
                continue

            # Keep track of same names using map so it becomes group_1, group_2 etc...
            key = child_name.lower()
            sibling_counts[key] = sibling_counts.get(key, 0) + 1
            if sibling_counts[key] > 1:
                child_name = f"{child_name}_{sibling_counts[key]}"

            child_rel_path = os.path.join(rel_path, child_name) if rel_path else child_name
            sub_dir = self._recursive_clustering(child_vecs, child_files, depth + 1, child_rel_path, builder)
            subdirs.append(sub_dir)

        # If subdirs are useless then stop creating deeper firs
        if not subdirs:
            return builder.buildDirectory(rel_path, files, [])

        return builder.buildDirectory(rel_path, retained_files, subdirs)


    def _normalize_child_name(self, candidate, parent_name, entry_files):
        """
        Returns a safe child folder name or None if we should NOT split to avoid 'X/X' or generic noise.
        - If candidate is empty/generic ('group', 'cluster'), derive from file extensions.
        - If equal to parent (case-insensitive), return None to avoid duplicate nesting.
        """
        name = (candidate or "").strip()

        # Names like group or cluster should rather use filetypes
        if not name or name.lower() in {"group", "cluster"}:
            name = self._name_from_extensions(entry_files)

        # Don't split if empty
        if not name:
            return None

        # Don't split if nested group names
        if parent_name and name.lower() == parent_name.lower():
            return None

        return name


    def _name_from_extensions(self, entry_files):

        exts = []
        for f in entry_files:
            _, ext = os.path.splitext(f["filename"])
            if ext:
                exts.append(ext.lstrip(".").lower())

        if not exts:
            return None

        top_ext, _ = Counter(exts).most_common(1)[0]

        # Map file types
        buckets = {
            "png": "images", "jpg": "images", "jpeg": "images", "gif": "images", "webp": "images",
            "pdf": "pdfs",
            "doc": "docs", "docx": "docs",
            "txt": "notes", "md": "notes", "rtf": "notes",
            "csv": "data", "xlsx": "data", "xls": "data",
            "ppt": "slides", "pptx": "slides"
        }
        return buckets.get(top_ext, top_ext)


    def sigmoid(self, x):
        return 1 / (1 + np.exp(-x))


    def printDirectoryTree(self, directory, indent=""):

        def human_bytes(n):
            try:
                n = int(n)
            except Exception:
                return None
            units = ["B", "KiB", "MiB", "GiB", "TiB", "PiB"]
            i = 0
            f = float(n)
            while f >= 1024 and i < len(units) - 1:
                f /= 1024.0
                i += 1
            if i == 0:
                return f"{int(f)} {units[i]}"
            return f"{f:.2f} {units[i]}"

        def dir_label(d):
            name = getattr(d, "name", None) or "<unnamed>"
            path = getattr(d, "path", "") or ""
            is_locked = bool(
                getattr(d, "is_locked", False) or getattr(d, "locked", False)
            )
            subdirs = getattr(d, "directories", None) or []
            files = getattr(d, "files", None) or []

            label = f"[DIR] {name}"
            if path and path != name:
                label += f" - {path}"
            label += f" [dirs={len(subdirs)} files={len(files)}]"
            if is_locked:
                label += " [locked]"
            return label

        def file_label(f):
            # Name
            name = getattr(f, "name", None)
            if name is None and isinstance(f, dict):
                name = f.get("name") or f.get("filename")
            if not name or str(name).strip() == "":
                name = "<unnamed>"

            # Path-ish
            path = getattr(f, "new_path", None)
            if path is None and isinstance(f, dict):
                path = f.get("new_path") or f.get("path") or f.get("original_path")

            # Size (optional)
            size = getattr(f, "size", None)
            if size is None and isinstance(f, dict):
                size = f.get("size") or f.get("size_bytes") or f.get("bytes")

            # Locked (optional)
            locked = getattr(f, "is_locked", None)
            if locked is None and isinstance(f, dict):
                locked = f.get("is_locked") or f.get("locked")

            label = f"[FILE] {name}"
            if size not in (None, 0, "0"):
                hb = human_bytes(size)
                if hb:
                    label += f" ({hb})"
            if path and path != name:
                label += f" - {path}"
            if locked:
                label += " [locked]"
            return label

        def sort_key_dir(sd):
            n = getattr(sd, "name", "") or getattr(sd, "path", "") or ""
            return str(n).lower()

        def sort_key_file(f):
            n = getattr(f, "name", None)
            if n is None and isinstance(f, dict):
                n = f.get("name") or f.get("filename")
            if not n:
                n = getattr(f, "new_path", None)
                if n is None and isinstance(f, dict):
                    n = (
                        f.get("new_path")
                        or f.get("path")
                        or f.get("original_path")
                        or ""
                    )
            return str(n).lower()

        def walk(d, prefix):
            if d is None:
                print(f"{prefix}(nil directory)")
                return

            subdirs = list(getattr(d, "directories", None) or [])
            files = list(getattr(d, "files", None) or [])

            subdirs.sort(key=sort_key_dir)
            files.sort(key=sort_key_file)

            entries = [("dir", sd) for sd in subdirs] + [
                ("file", f) for f in files
            ]

            if not entries:
                print(f"{prefix}└── (empty)")
                return

            for i, (kind, item) in enumerate(entries):
                is_last = i == len(entries) - 1
                branch = "└── " if is_last else "├── "
                next_prefix = prefix + ("    " if is_last else "│   ")

                if kind == "dir":
                    print(f"{prefix}{branch}{dir_label(item)}")
                    walk(item, next_prefix)
                else:
                    print(f"{prefix}{branch}{file_label(item)}")

        base = indent or ""
        if directory is None:
            print("(nil directory)")
            return

        # Print root and recurse
        print(f"{base}{dir_label(directory)}")
        walk(directory, base)

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
        # Build a nested tree keyed by relative path parts
        dir_tree = {}
    
        for f in files:
            full_path = os.path.normpath(f["original_path"])
            parts = full_path.split(os.sep)
            try:
                parent_index = parts.index(self.parent_folder)
                # Relative parts: directories under root (exclude filename at the end)
                relative_parts = parts[parent_index + 1 : -1]
            except (ValueError, IndexError):
                relative_parts = ["Unknown"]
    
            node = dir_tree
            for part in relative_parts:
                node = node.setdefault(part, {})
            node.setdefault("_files", []).append(f)
    
        def build_dirs(node, rel_prefix=""):
            children = []
            for name, sub in node.items():
                if name == "_files":
                    continue
                rel_path = os.path.join(rel_prefix, name) if rel_prefix else name
                subdirs = build_dirs(sub, rel_path)
                files_here = sub.get("_files", [])
                children.append(builder.buildDirectory(rel_path, files_here, subdirs))
            return children

        root_files = dir_tree.get("_files", [])
        return builder.buildDirectory("", root_files, build_dirs(dir_tree, ""))




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

