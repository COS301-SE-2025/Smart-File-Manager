import hdbscan

class HDBSCANCluster:
    def __init__(self, min_cluster_size=3, metric="euclidean"):
        self.model = hdbscan.HDBSCAN(min_cluster_size=min_cluster_size, metric=metric)

    def cluster(self, vectors):
        return self.model.fit_predict(vectors)
