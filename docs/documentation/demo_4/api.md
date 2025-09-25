# Smart File Manager – Flutter API Client Documentation
**Version:** 3.0.0.0  
**Prepared By:** Spark Industries  
**Prepared For:** Southern Cross Solutions & Personal Development team use  

## Content
* [Introduction](#introduction)
* [A note on security](#a-note-on-security)
* [General Notes](#general-notes)
* [loadTreeData](#loadtreedata)
* [sortTree](#sorttree)
* [moveDiretory](#movedirectory)
* [addSmartManager](#addsmartmanager)
* [deleteManager](#deletemanager)
* [addTagToFile](#addtagtofile)
* [deleteTagFromFile](#deletetagfromfile)
* [deleteFile](#deletefile)
* [deleteFolder](#deletefolder)
* [bulkDeleteFiles](#bulkdeletefiles)
* [bulkDeleteFolders](#bulkdeletefolders)
* [bulkAddTag & bulkRemoveTag](#bulk-add-and-bulk-remove-tags)
* [findDuplicates](#findduplicates)
* [returnType](#returntype)
* [returnStats](#returnstats)
* [Locking](#locking)
* [Unlocking](#unlocking)
* [StartUp](#startup)
* [Search](#search)
* [keywordSearch](#keywordsearch)
* [isKeywordSearchReady](#iskeywordsearchready)


## Introduction
Our team makes uses of standard Rest API endpoints to connect our filesystem server (go backend) to our frontend. The reasons we did not use gRPC for these endpoints (like our connection from the filesystem server to the clustering server) is as follows:

* The size of requests are much smaller than the ones sent to the filesystem server.
* It is easier to use traditional JSON based payloads for frontend development.

> **Base URI**: `http://localhost:51000/`

## A note on security
Our application is deployed as a standalone desktop application which does not rely on (or use) internet connection. That being said we require all requests to include an apiSecret header. This value is set during the installation process.

## General Notes 

* Make sure to URI-encode query parameter values when needed.
* The backend must be running at the defined base URI for requests to succeed.
* apiSecret header must be included for each request

---

## loadTreeData

Fetches the file tree structure for a Smart Manager.

### Usage
```GET /loadTreeData?name={name}```

### Parameters

* name: The name of the Smart Manager.

### Returns

A nested JSON structure representing the full directory tree and tags.
Example response:

```json
{
  "name": "root",
  "isFolder": true,
  "children": [
    {
      "name": "file1.txt",
      "path": "c:://",
      "isFolder": false,
      "tags": ["work", "important"],
      "metadata": {
        "size": "12KB",
        "dateCreated": "2023-11-08T14:20:00Z",
        "mimeType": "Tiaan/Bosman",
        "lastModified": "2024-02-20T10:10:00Z"
      },
    },
    {
      "name": "file2.docx",
      "path": "c:://",
      "isFolder": false,
      "tags": ["document"],
      "metadata": {
        "size": "12KB",
        "dateCreated": "2023-11-08T14:20:00Z",
        "mimeType": "Tiaan/Bosman",
        "lastModified": "2024-02-20T10:10:00Z"
      },
    },
    {
      "name": "Documents",
      "isFolder": true,
      "children": [
        {
          "name": "resume.pdf",
          "path": "c:://",
          "isFolder": false,
          "tags": ["personal", "career"],
          "metadata": {
            "size": "12KB",
            "dateCreated": "2023-11-08T14:20:00Z",
            "mimeType": "Tiaan/Bosman",
            "lastModified": "2024-02-20T10:10:00Z"
          },
        },
        {
          "name": "Reports",
          "isFolder": true,
          "children": [
            {
              "name": "annual_2023.pdf",
              "path": "c:://",
              "isFolder": false,
              "tags": ["report", "finance"],
              "metadata": {
                "size": "12KB",
                "dateCreated": "2023-11-08T14:20:00Z",
                "mimeType": "Tiaan/Bosman",
                "lastModified": "2024-02-20T10:10:00Z"
              },
            },
            {
              "name": "q1_summary.docx",
              "path": "c:://",
              "isFolder": false,
              "tags": ["summary", "q1"],
              "metadata": {
                "size": "12KB",
                "dateCreated": "2023-11-08T14:20:00Z",
                "mimeType": "Tiaan/Bosman",
                "lastModified": "2024-02-20T10:10:00Z"
              },
            }
          ]
        }
      ]
    }
  ]
}
```

### Throws

Exception if the request fails.

## sortTree

Fetches the sorted version of the file tree structure for a Smart Manager.

### Usage
```GET /sortTree?name={name}```

### Parameters

* name: The name of the Smart Manager.

### Returns

A nested JSON structure representing the full directory tree and tags.
Example response:

```json
{
  "name": "root",
  "isFolder": true,
  "children": [
    {
      "name": "file1.txt",
      "path": "c:://",
      "isFolder": false,
      "tags": ["work", "important"],
      "metadata": {
        "size": "12KB",
        "dateCreated": "2023-11-08T14:20:00Z",
        "mimeType": "Tiaan/Bosman",
        "lastModified": "2024-02-20T10:10:00Z"
      },
    },
    {
      "name": "file2.docx",
      "path": "c:://",
      "isFolder": false,
      "tags": ["document"],
      "metadata": {
        "size": "12KB",
        "dateCreated": "2023-11-08T14:20:00Z",
        "mimeType": "Tiaan/Bosman",
        "lastModified": "2024-02-20T10:10:00Z"
      },
    },
    {
      "name": "Documents",
      "isFolder": true,
      "children": [
        {
          "name": "resume.pdf",
          "path": "c:://",
          "isFolder": false,
          "tags": ["personal", "career"],
          "metadata": {
            "size": "12KB",
            "dateCreated": "2023-11-08T14:20:00Z",
            "mimeType": "Tiaan/Bosman",
            "lastModified": "2024-02-20T10:10:00Z"
          },
        },
        {
          "name": "Reports",
          "isFolder": true,
          "children": [
            {
              "name": "annual_2023.pdf",
              "path": "c:://",
              "isFolder": false,
              "tags": ["report", "finance"],
              "metadata": {
                "size": "12KB",
                "dateCreated": "2023-11-08T14:20:00Z",
                "mimeType": "Tiaan/Bosman",
                "lastModified": "2024-02-20T10:10:00Z"
              },
            },
            {
              "name": "q1_summary.docx",
              "path": "c:://",
              "isFolder": false,
              "tags": ["summary", "q1"],
              "metadata": {
                "size": "12KB",
                "dateCreated": "2023-11-08T14:20:00Z",
                "mimeType": "Tiaan/Bosman",
                "lastModified": "2024-02-20T10:10:00Z"
              },
            }
          ]
        }
      ]
    }
  ]
}
```

### Throws

Exception if the request fails.

## moveDirectory

After sorting a Smart Manager, the user can choose to actually move the files that were sorted into a new sorted directory.

### Usage
````GET /moveDirectory?name={name}```

### Parameters

* name: Smart Manager name

### Returns

true on success.

false on failure.

### Throws

Exception if invalid.

## addSmartManager

Creates and registers a new Smart Manager by mounting a specified folder path.

### Usage
```GET /addDirectory?name={name}&path={path}```

### Parameters

* name: Unique identifier for the Smart Manager.

* path: Absolute path to the folder to mount.

### Returns

true on success.

### Throws

Exception if creation fails.

## deleteManager

Deletes a Smart Manager and any associated data.

### Usage
```GET /deleteManager?name={name}```

### Parameters

* name: Name of the Smart Manager to delete.

### Returns

true on success.

Manager not found on failure.

### Throws

Exception if deletion fails.

## deleteFile

Deletes a file contained in specified manager.

### Usage
```GET /deleteFile?name={name}&path={path}```

### Parameters

* name: Name of the Smart Manager to delete from.
* path : Path to the file to delete

### Returns

A JSON structure consisting of all files

```json
{
  "name": "TagManager",
  "isFolder": true,
  "rootPath": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers",
  "children": [
    {
      "name": "Chapter-9-World-War-I-1914-1918.pdf",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
      "isFolder": false,
      "tags": [
        "history",
        "world-war-1",
        "pdf",
        "academic",
        "chapter-9"
      ],
      "metadata": {
        "size": "10297185",
        "dateCreated": "2025-06-27 07:57",
        "mimeType": ".pdf",
        "lastModified": "2025-06-27 07:57"
      },
      "locked": false
    }
  ]
}
```

'false' if not found.

### Throws

Exception if deletion fails.

## deleteFolder

Deletes a folder contained in specified manager.

### Usage
```GET /deleteFile?name={name}&path={path}```

### Parameters

* name: Name of the Smart Manager to delete from.
* path : Path to the folder to delete

### Returns

A JSON structure consisting of all files

```json
{
  "name": "TagManager",
  "isFolder": true,
  "rootPath": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers",
  "children": [
    {
      "name": "Chapter-9-World-War-I-1914-1918.pdf",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
      "isFolder": false,
      "tags": [
        "history",
        "world-war-1",
        "pdf",
        "academic",
        "chapter-9"
      ],
      "metadata": {
        "size": "10297185",
        "dateCreated": "2025-06-27 07:57",
        "mimeType": ".pdf",
        "lastModified": "2025-06-27 07:57"
      },
      "locked": false
    }
  ]
}
```

'false' if not found.

### Throws

Exception if deletion fails.

## bulkDeleteFiles

Deletes a list of files contained in specified manager.

### Usage
```GET /bulkDeleteFiles?name={name}, json body```

### Parameters

* name: Name of the Smart Manager to delete from.
* JSON body containing list of files to delete

```json
[

{
"file_path": "/home/user/documents/report.pdf",
},
{
"file_path": "/home/user/photos/vacation.jpg",
},
{
"file_path": "/home/user/music/song.mp3",
}
]
```


### Returns

A JSON structure consisting of all files

```json
{
  "name": "TagManager",
  "isFolder": true,
  "rootPath": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers",
  "children": [
    {
      "name": "Chapter-9-World-War-I-1914-1918.pdf",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
      "isFolder": false,
      "tags": [
        "history",
        "world-war-1",
        "pdf",
        "academic",
        "chapter-9"
      ],
      "metadata": {
        "size": "10297185",
        "dateCreated": "2025-06-27 07:57",
        "mimeType": ".pdf",
        "lastModified": "2025-06-27 07:57"
      },
      "locked": false
    }
  ]
}
```

'false' if not found.

### Throws

Exception if deletion fails.

## bulkDeleteFolders

Deletes a list of folders contained in specified manager.

### Usage
```GET /bulkDeleteFolders?name={name}, json body```

### Parameters

* name: Name of the Smart Manager to delete from.
* JSON body containing list of folders to delete

```json
[

{
"file_path": "/home/user/documents",
},
{
"file_path": "/home/user/photos",
},
{
"file_path": "/home/user/music",
}
]
```


### Returns

A JSON structure consisting of all files

```json
{
  "name": "TagManager",
  "isFolder": true,
  "rootPath": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers",
  "children": [
    {
      "name": "Chapter-9-World-War-I-1914-1918.pdf",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
      "isFolder": false,
      "tags": [
        "history",
        "world-war-1",
        "pdf",
        "academic",
        "chapter-9"
      ],
      "metadata": {
        "size": "10297185",
        "dateCreated": "2025-06-27 07:57",
        "mimeType": ".pdf",
        "lastModified": "2025-06-27 07:57"
      },
      "locked": false
    }
  ]
}
```

'false' if not found.

### Throws

Exception if deletion fails.

## addTagToFile

Adds a tag to a specific file under a Smart Manager.

### Usage
```GET /addTag?path={path}&tag={tag}```

### Parameters

path: Path to the file.

tag: Tag to assign.

### Returns

true on success.

### Throws

Exception if tagging fails.

## deleteTagFromFile

Removes a tag from a specific file.

### Usage
```GET /removeTag?path={path}&tag={tag}```

### Parameters

* path: Path to the file.

* tag: Tag to remove.

### Returns

true on success.

### Throws

Exception if removal fails.

## Locking

Locks a file or folder. When locking a folder all sub-folders and files are also locked.

### Usage
```/lock?name={name}&path=../../testRootFolder```

### Returns

true on success.

### Throws

Exception if removal fails.

## Unlocking

Unlocks a file or folder. When unlocking a folder all sub-folders and files are also unlocked.

### Usage
```/unlock?name={name}&path=../../testRootFolder```

### Returns

true on success.

### Throws

Exception if removal fails.

## startUp

Loads smart managers into memory and returns the respective names

### Usage
GET /startUp

### Parameters

None

### Response
```json
{
  "responseMessage": "Request successful!, composites: 1",
  "managerNames": [
    "first2"
  ]
}
```

## findDuplicate

Identifies and returns a list of files that are considered duplicates based on file size, hash, and content matching.

### Usage
```GET /findDuplicates?name={name}```

### Parameters

* name: The name of the Smart Manager to search for duplicates in.

### Returns

A JSON array of duplicate file pairs. Each object includes:

* name: The common file name.

* original: The full path to the original file.

* duplicate: The full path to the detected duplicate.

Example Response
```json
[
  {
    "name": "The Man of Steel.docx",
    "original": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers2/The Man of Steel (another copy).docx",
    "duplicate": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers2/The Man of Steel.docx"
  },
  {
    "name": "pyramid-technology.pdf",
    "original": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers2/pyramid-technology (copy).pdf",
    "duplicate": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers2/pyramid-technology.pdf"
  }
]
```

### Throws

Exception if the request fails or if the Smart Manager name is invalid.

## returnType

Identifies and returns a list of files of the same file type.

### Usage
```GET /returnType?name={name}&type{filetype}&umbrella{bool}```

### Parameters

* name: The name of the Smart Manager to retrieve the file types from.
* type: File types to retrieve, can be specific(pdf) or umbrella(Documents).
* umbrella: Boolean that determines if umbrella types should be returned.

### Returns

A JSON array of files that fit the parameters. Each object includes:

* file_name: The common file name.

* file_path: The full path to the file.


Example Response
```json
[
  {
    "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/egypts-pyramids.pdf",
    "file_name": "egypts-pyramids.pdf"
  },
  {
    "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
    "file_name": "Chapter-9-World-War-I-1914-1918.pdf"
  },
  {
    "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Origins-of-World-War-I-Essay-ojznse.pdf",
    "file_name": "Origins-of-World-War-I-Essay-ojznse.pdf"
  },
  {
    "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/World_War_1_The_Great_War_and_its_Impact_OA_edition.pdf",
    "file_name": "World_War_1_The_Great_War_and_its_Impact_OA_edition.pdf"
  },
  {
    "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/pyramid-technology.pdf",
    "file_name": "pyramid-technology.pdf"
  }
]
```

### Throws

Exception if the request fails or if the Smart Manager name is invalid.

## returnStats

Returns a set of statistics from all managers for the dashboard page

### Usage
```GET /returnStats```

### Parameters
* N/A

### Returns

A JSON array of objects containing stats per manager. Each object includes:

* manager_name: Name of the manager
* size: Size of the manager(kB)
* folders: Amount of folders in manager.
* files: Amount of files in manager.
* recent: List of 5 most recently accessed files(lim 5).
* largest: List of 5 largest files(lim 5).
* oldest: List of 5 leasr recently accessed files(lim 5).


Example Response
```json
[
  {
    "manager_name": "manager3",
    "size": 16289179,
    "folders": 0,
    "files": 45,
    "recent": [
      {
        "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/test_files_3/Main.form",
        "file_name": "Main.form"
      }...
    ],
    "largest": [
      {
        "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/test_files_3/~WRL1847.tmp",
        "file_name": "~WRL1847.tmp"
      }...
    ],
    "oldest": [
      {
        "file_path": "/home/henco/Documents/University-Files/Third-year/COS-301/test_files_3/Apr18 meeting.txt",
        "file_name": "Apr18 meeting.txt"
      }...
    ],
    "umbrella_counts": [
      25,
      15,
      0,
      0,
      0,
      1,
      0,
      4
    ]
  }
]
```

### Throws

N/A returns empty structure if no managers found

## Bulk add and Bulk Remove tags

This functionality allows adding and removing tags in bulk.

### Usage
```
POST /bulkAddTag?name{manager name}, json body
POST /bulkRemoveTag?name{manager name}, json body
```

### Parameters

* name: Name of the Smart Manager.

* JSON body of all files that require tags.

Example JSON body:
```json
[
  {
    "file_path": "/home/user/documents/report.pdf",
    "tags": ["work", "important", "pdf"]
  },
  {
    "file_path": "/home/user/photos/vacation.jpg",
    "tags": ["holiday", "family", "2025"]
  },
  {
    "file_path": "/home/user/music/song.mp3",
    "tags": ["music", "mp3", "favorites"]
  }
]
```

### Returns

A JSON structure consisting of all files with the included tags.

```json
{
  "name": "TagManager",
  "isFolder": true,
  "rootPath": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers",
  "children": [
    {
      "name": "Chapter-9-World-War-I-1914-1918.pdf",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/Papers/Chapter-9-World-War-I-1914-1918.pdf",
      "isFolder": false,
      "tags": [
        "history",
        "world-war-1",
        "pdf",
        "academic",
        "chapter-9"
      ],
      "metadata": {
        "size": "10297185",
        "dateCreated": "2025-06-27 07:57",
        "mimeType": ".pdf",
        "lastModified": "2025-06-27 07:57"
      },
      "locked": false
    }
  ]
}
```

### Throws

Exception "Manager not found" if Smart Manager could not be found.

## Search

Finds files of a given name inside a specific Smart Manager.

### Usage
```GET /search?compositeName={name}&searchText={fileToSearch}```

### Parameters

* compositeName: Unique identifier for the Smart Manager.

* searchText: File name to search for.

### Returns

A JSON array of all the files that match the specified search criteria. Each file includes:

* name: string

* path: string

* isFolder: boolean

* tags: string[] (All user-defined tags for the file)

* metadata: string[] (All metadata extracted for the file)

```json
{
  "name": "Josh",
  "isFolder": true,
  "children": [
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_2/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    },
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_1/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    }
  ]
}
```

### Throws

Not applicable (returns empty structure if nothing found).


## Search

Finds files of a given name inside a specific Smart Manager.

### Usage
```GET /search?compositeName={name}&searchText={fileToSearch}```

### Parameters

* compositeName: Unique identifier for the Smart Manager.

* searchText: File name to search for.

### Returns

A JSON array of all the files that match the specified search criteria. Each file includes:

* name: string

* path: string

* isFolder: boolean

* tags: string[] (All user-defined tags for the file)

* metadata: string[] (All metadata extracted for the file)

```json
{
  "name": "Josh",
  "isFolder": true,
  "children": [
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_2/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    },
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_1/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    }
  ]
}
```

### Throws

Not applicable (returns empty structure if nothing found).

## keywordSearch

Extracts keywords from files found in the smartfile manager, it then stores the keywords and tags and locks in a json file. on startup it loads the composite in, reads the stored json and reads over the keywords.

### Usage
```GET /keywordSearch?compositeName={name}&searchText={fileToSearch}```

### Parameters
* compositeName: Unique identifier for the Smart Manager.
* searchText: File name to search for.

### Return
A JSON array of all the files that match the specified search criteria. Each file includes:
* name: string
* path: string
* isFolder: boolean
* tags: string[] (All user-defined tags for the file)
* metadata: string[] (All metadata extracted for the file)

```json
{
  "name": "Josh",
  "isFolder": true,
  "children": [
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_2/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    },
    {
      "name": "tiaan.jpeg",
      "path": "/home/henco/Documents/University-Files/Third-year/COS-301/GIT/Smart-File-Manager/docs/documentation/demo_1/assets/readmeAssets/tiaan.jpeg",
      "isFolder": false,
      "metadata": {
        "size": "202624",
        "dateCreated": "2025-07-08 13:19",
        "mimeType": "",
        "lastModified": "2025-07-08 13:19"
      },
      "locked": false
    }
  ]
}
```

### Throws
Not applicable (returns empty structure if nothing found).


## isKeywordSearchReady
Extracts keywords from files found in the smartfile manager, it then stores the keywords and tags and locks in a json file. on startup it loads the composite in, reads the stored json and reads over the keywords.

### Usage
```GET /isKeywordSearchReady?compositeName={name}```

### Parameters
* compositeName: Unique identifier for the Smart Manager.

### Returns
Boolean value based on if the composite has retrieved keywords

### Throws
Bad request if manager not found.
