# use numpy to normalize most data types, will probably be in kb

import numpy as np

def normalize(vec):
    vec = np.array(vec, dtype=float)
    min_val = vec.min()
    max_val = vec.max()
    if min_val == max_val:
        return np.zeros_like(vec)
    return (vec - min_val) / (max_val - min_val)

# if name == main
# arr ={}
# normalize(arr)
