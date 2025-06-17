# the actual clustering that will make a directory to send back to go

from sklearn.cluster import KMeans
import numpy as np
import json # for reading a json file. Will not be used


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


if __name__ == "__main__":
    #Temporary read data from file
    X = []
    with open("python/testing/exampleDataKMeans.txt", "r") as f:
        for line in f:
            entry = json.loads(line.strip())   
            X.append(entry["full_vector"])     

    X = np.array(X)  # Convert list of lists to numpy array
         
    k_means = KMeansCluster(3)
    labels = k_means.cluster(X)
    print("Labels: ", labels)
        # Here, vectors have length = 3 keyword + 3 filetype + 1 size = 7 features
    new_points = [
        [0.5, 0.5, 0.5, 1, 0, 0, 0.1],
        [0.1, 0.2, 0.3, 0, 1, 0, 0.4],
        [0.9, 0.9, 0.9, 0, 0, 1, 0.7]
    ]

    predictions,centers = k_means.predict(new_points)
    print("Pred: ", predictions)
    print("Centers: ", centers)    


