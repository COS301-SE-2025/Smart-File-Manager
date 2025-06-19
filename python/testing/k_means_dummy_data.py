import numpy as np
import random
import json

# Set a seed for reproducibility
np.random.seed(42)

# File type one-hot encoding 
filetypes = {
    ".txt": [1, 0, 0],
    ".pdf": [0, 1, 0],
    ".docx": [0, 0, 1]
}

# Number of samples to generate
num_files = 10 

files = []
for _ in range(num_files):
    #individual components
    keyword_vectors = [random.uniform(0,1) for _ in range(3)]  # 3 keyword features per file
    rand_filetype = random.choice(list(filetypes.keys()))
    filetype_vectors = filetypes[rand_filetype]
    size_vectors = [random.uniform(0,1)]  # single normalized size feature

    # Combine all into full_vectors
    full_vector = keyword_vectors + filetype_vectors + size_vectors

    # (Optional) Convert to dicts
    files.append({"full_vector":full_vector})


with open("exampleDataKMeans.txt", "w") as f:
    for entry in files:
        f.write(json.dumps(entry) + "\n")
