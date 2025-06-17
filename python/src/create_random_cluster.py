import numpy as np
import random

def generate_random_point():
    # First 60 random floats as np.float64
    floats_part = [np.float64(random.uniform(0, 1)) for _ in range(60)]
    
    
    # One-hot encoding for filetype (3 categories)
    filetype_index = random.randint(0, 2)
    filetype_part = [0.0, 0.0, 0.0]
    filetype_part[filetype_index] = 1.0
    
    # Random file size (integer) between 0 and 100,000
    size = random.randint(0, 100000)
    
    # Combine all parts
    full_point = floats_part + filetype_part + [size]
    return full_point

# Generate example data points
random_data = [generate_random_point() for _ in range(5)]

for i, point in enumerate(random_data, 1):
    print(point, ",")
