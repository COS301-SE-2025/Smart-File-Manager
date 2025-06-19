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
        return self.kmeans.labels_

    def predict(self, points):
        predictions = self.kmeans.predict(points)
        centers = self.kmeans.cluster_centers_
        centers_rounded = np.round(centers, 4) # rounded to get mostly matching        
        return predictions, centers_rounded

 


