class KWCluster:

    # only add exact matches for now
    def createCluster(self, keywords, vocab):
        kw_score_map = {kw: score for kw, score in keywords}
        cluster_vector = []
        for kw in vocab:
            cluster_vector.append(kw_score_map.get(kw, 0))  # score if exists else 0

        return cluster_vector