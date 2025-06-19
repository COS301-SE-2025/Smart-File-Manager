import numpy as np
import random

def generate_random_point():
    # First 5000 random floats as np.float64
    floats_part = [np.float64(random.uniform(0, 1)) for _ in range(500)]
    
    
    # One-hot encoding for filetype 
    filetype_index = random.randint(0, 9)
    filetype_part = [0.0,0.0,0.0, 0.0,0.0,0.0, 0.0,0.0,0.0, 0.0]
    filetype_part[filetype_index] = 1.0
    
    # Random file size (integer) between 0 and 100,000
    size = random.randint(0, 100000)
    
    # Combine all parts
    full_point = floats_part + filetype_part + [size]
    return full_point

# Generate example data points
random_data = [generate_random_point() for _ in range(1000)]


with open("output.txt", "w") as f:
    for i, point in enumerate(random_data, 1):
        f.write(f"{point},\n")
