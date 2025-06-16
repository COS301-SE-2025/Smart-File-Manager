# the actual clustering that will make a directory to send back to go

from sklearn.cluster import KMeans
import numpy as np


class KMeansCluster:
    def __init__(self, numClusters):
        self.kmeans = KMeans(
            n_clusters=numClusters,
            random_state=0,
            n_init="auto"
            )

    def cluster(self,files):
        self.kmeans.fit(files)
        self.kmeans.labels_
        np.array([1, 1, 1, 0, 0, 0], dtype=np.int32)

    def predict(self):
        self.kmeans.predict([[0,0],[12,3]])
        np.array([1,0], dtype=np.int32)
        self.kmeans.cluster_centers_
        np.array([[10.,2.],[1.,2.]])


if __name__ == "__main__":
    X = np.array([[1, 2], [1, 4], [1, 0],
                [10, 2], [10, 4], [10, 0]])
    k_means = KMeansCluster(2)
    k_means.cluster(X)
    result = k_means.predict()
    print(result)
