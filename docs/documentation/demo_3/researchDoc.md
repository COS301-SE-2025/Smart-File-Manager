# Capstone SFM Research Document

# Document Overview

This document aims to serve as a research document that discusses details regarding the implementation of the SFM project. It mainly focuses on the approach that will be taken to implement the “auto place” functionality used to determine where in a folder structure a file should be placed based on metadata.

# File Organization

## How are file usually organized and its characteristics

The standard way of organizing files is based on the context in which the file is used. For example for University, files may be organized in subdirectories by Study Year ⇒ Semester ⇒ Module ⇒ Assignment Type ⇒ etc… 

A structured approach as outlined allows for files to be found in a directory not by memorizing the location but by answering simple questions about the use of said file, serving as directions to its location.

Importantly it should then be noted that the files are not organized by some single metric but rather based on the “semantic value” that the file has. This forms the centre of our definition of “smart” in the context of a file manager. A smart file manager should be able to autonomously detect the correct context of a given file while adhering to best organizational practices.

We define a series of such organizational practices as follows:

1. The relation between entries (files or subdirectories) in a directory should be visible at first glance
2. Directories should impose a maximum limit of the amount of entries (should be user defined)
3. Given a directory tree, files should preferably be located at a similar “level” i.e. a directory should not contain only subdirectories with a single file, as this would seem out of place.

# Collecting Information On Files

In order to make informed decisions on where a file must be located we first need adequate data to base such a decision on. 

**Metadata**

Modern files contain metadata which contains useful information to characterize a file. Metadata may be extracted from a file using Python modules such as Pillow or Piexif. The type of metadata that may be extracted from a file varies depending on the file type (metadata in itself) but for a .pdf includes:

- checksum
- file_name
- file_size
- file_type
- mime_type
- pdf_version
- linearized
- page_count
- page_mode
- title
- creator
- create_date
- etc…

**Keywords**

Text based files such as .pdf, .txt, .rtf, .html etc… may further be scanned for common keywords. A high frequency of certain keywords may then also be used as keywords common among certain files may imply a logical connection between them. Note that care must be taken to sanitize keywords to remove redundant information (commonly used words such as indefinites etc.).

**Tags**

Users may also wish to apply their own tags to a file to help the system better label the file. 

# Extracting Information from File Data

Having looked at a variety of sources that may be used to gain information on the file it follows that some process must then be followed to organize the files based on this data. This involves parsing the data to determine a similarity between different files. Files may then be grouped via this similarity into a directory tree structure.

These types of classification based problems lend themselves to AI based solutions. We specifically consider clustering approaches for the problem.

## K-means clustering

K means is an unsupervised machine learning model that can be used to group unlabelled data (files) into different clusters (directories). This is the same approach used by online retailers to provide functionality such as “users also bought”. The primary challenge with this approach involves converting our data (vectors of metadata and keywords) into a form for which a “distance” metric may be assigned.

### Encoding Metadata

Note: To achieve all functionality in this explanation we make use of scikit-learn open source library for machine learning in python

Suppose we are given 3 files for which we have extracted the following metadata 

```json
files = [
    {
        "filename": "invoice_q1.pdf",
        "keywords": ["invoice", "payment", "q1"],
        "size_kb": 120,
        "filetype": "pdf"
    },
    {
        "filename": "meeting_notes.docx",
        "keywords": ["meeting", "project", "minutes"],
        "size_kb": 85,
        "filetype": "docx"
    },
    {
        "filename": "holiday_photo.jpg",
        "keywords": ["holiday", "beach", "family"],
        "size_kb": 2048,
        "filetype": "jpg"
    }
]
```

**Vectorizing keywords** 

We first vectorize the keywords into a numerical format TF-IDF (term frequency-inverse document frequency) - a numerical statistic which can be used to evaluate the importance of word across a collection of documents. 

To do this we combine the keywords across all files into a single array and use the *TfidfVectorizer.* Keywords may then be represented as floating points values for example:

```json
files = [
    {
        "filename": "invoice_q1.pdf",
        "keywords_vectorized": [0.32, 0.58, 0.392],

    },
    {
        "filename": "meeting_notes.docx",
        "keywords_vectorized": [0.42, 0.2421, 0.99],

    },
    {
        "filename": "holiday_photo.jpg",
        "keywords_vectorized": [0.23, 0.69, 0.92],
    }
]
```

   Note: Removed other keys for brevity

Some pseudocode

```python
from sklearn.feature_extraction.text import TfidfVectorizer

# Join keywords into a single string for each file
docs = [" ".join(file["keywords"]) for file in files]

vectorizer = TfidfVectorizer()
keyword_vectors = vectorizer.fit_transform(docs).toarray()
```

**Vectorizing categorical data**

Categorical data such as filetypes may be encoded using “One-hot” encoding. An encoding scheme to represent categorical data as binary vectors. This takes each category and converts it into a column of a vector with a 1 indicating the category is present and 0 indicating an absence.

Suppose we support only the file types “.pdf”, “.jpg”, “.docx” then file type could be encoded as

```json
files = [
		{
			"filename" : "invoice_q1.pdf",
			"filetype_encoded" : [1, 0, 0] // 1 in index 0 ==> .pdf
		}
		{
			"filename" : "meeting_notes.docx",
			"filetype_encoded" : [0, 0, 1] // 1 in index 2 ==> .docx
		}
]
```

Note: Removed other keys for brevity

Some pseudocode

```python
from sklearn.preprocessing import OneHotEncoder
import numpy as np

filetypes = [[file["filetype"]] for file in files]
encoder = OneHotEncoder(sparse=False)
filetype_vectors = encoder.fit_transform(filetypes)
```

**Normalization of other data**

Other numerical data such as the file size should also be normalized to fit in the range [0..1]. 

Some pseudocode

```python
from sklearn.preprocessing import MinMaxScaler

sizes = np.array([[file["size_kb"]] for file in files])
scaler = MinMaxScaler()
size_vectors = scaler.fit_transform(sizes)

```

**Bringing it all together**

Finally we concatenate all vectors to form a n-dimensional real number based representation of all our data on the files. This data may then be used to perform K-means clustering on

Some pseudocode

```python
full_vectors = np.hstack([keyword_vectors, filetype_vectors, size_vectors])
```

All data should then be represented as follows

```json
files = [
    {
			"full_vector" : [0.56, 0.87, 0.34   1,0,0,   0.02]
    },
    {
			"full_vector" : [0.26, 0.17, 0.54   0,0,1,   0.24]

    },
    {
			"full_vector" : [0.23, 0.45, 0.45   0,1,0,   0.22]
    }
]
```

In reality these vectors will be considerably longer

### Running K-means

The sci-kit learn K-means algorithm can then be applied to the data

[https://scikit-learn.org/stable/modules/generated/sklearn.cluster.KMeans.html](https://scikit-learn.org/stable/modules/generated/sklearn.cluster.KMeans.html)

This function also includes parameters for number of clusters to generate. How to choose initial centroid clusters etc…

The K-means algorithm may also be used recursively for generating a directory tree up to a certain depth. Consider the following high level overview of such an implementation

```python
# I wrote this while very tired so if it makes no sense sorry...
def findClusters(files, nrClusters):
	
	if nrSubTrees > 0:
		Run K-means with nrClusters
		for each cluster:
			Create directory with appropriate name
			findClusters(cluster_files, nrClusters)
	else:
		place files in current cluster

findClusters(original_files, 3) # Creates folder-strucutre with 3 sub-dirs for each dir
	
```

## Alternative and additional ways of extracting information

### HDBSCAN

Another built in method in scikit-learn, this clusters data using hierarchical density based clustering.

Subject to the same encoding scheme as describe above

## Cosine Similarity

If keywords lists may be converted into term-frequency vectors cosine similarity 

$Consine(A,B) = (A.B) / ||A||||B||$ 

which may capture overlap and similarity between two files’ keywords.

### Semantic Embedding (BERT)

Idk what this is exactly yet but seems cool