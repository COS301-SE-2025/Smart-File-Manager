# Smart File Manager – Flutter API Client Documentation
Version 1.0.0


> **Base URI**: `http://localhost:51000/`

---

## loadTreeData

Fetches the file tree structure for a Smart Manager.

**Parameters:**

* `name`: The name of the Smart Manager.

**Endpoint:**

```
GET /loadTreeData?name={name}
```

**Returns:**

* A nested JSON structure representing the full directory tree and tags. Example Below of response:

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
      "metadata" : {
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
      "metadata" : {
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
          "metadata" : {
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
              "metadata" : {
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
              "metadata" : {
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

**Throws:**

* Exception if the request fails.

---
## sortTree

Fetches the sorted version of file tree structure for a Smart Manager.

**Parameters:**

* `name`: The name of the Smart Manager.

**Endpoint:**

```
GET /sortTree?name={name}
```

**Returns:**

* A nested JSON structure representing the full directory tree and tags. Example Below of response:

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
      "metadata" : {
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
      "metadata" : {
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
          "metadata" : {
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
              "metadata" : {
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
              "metadata" : {
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

**Throws:**

* Exception if the request fails.

---

## addSmartManager

Creates and registers a new Smart Manager by mounting a specified folder path.

**Parameters:**

* `name`: Unique identifier for the Smart Manager.
* `path`: Absolute path to the folder to mount.

**Endpoint:**

```
POST /addDirectory?name={name}&path={path}
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if creation fails.

---

## deleteSmartManager

Deletes a Smart Manager and any associated data.

**Parameters:**

* `name`: Name of the Smart Manager to delete.

**Endpoint:**

```
POST /deleteDirectory?name={name}
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if deletion fails.

---

## addTagToFile

Adds a tag to a specific file under a Smart Manager.

**Parameters:**

* `path`: Path to the file.
* `tag`: Tag to assign.

**Endpoint:**

```
POST /addTag?path={path}&tag={tag}
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if tagging fails.

---

## deleteTagFromFile

Removes a tag from a specific file.

**Parameters:**

* `path`: Path to the file.
* `tag`: Tag to remove.

**Endpoint:**

```
POST /removeTag?path={path}&tag={tag}
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if removal fails.

## Locking

Locks a file or folder. When locking a folder all sub-folders, -files are also locked.

**Endpoint**

```
/lock?path=../../testRootFolder
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if removal fails.
  
## Unlocking

Unlocks a file or folder. When unlocking a folder all sub-folders, -files are also unlocked

**Endpoint**
```
/unlock?path=../../testRootFolder
```

**Returns:**

* `true` on success.

**Throws:**

* Exception if removal fails.
---

## Notes

* All `POST` endpoints use query parameters; no request body is sent.
* Make sure to URI-encode query parameter values when needed.
* The backend must be running at the defined base URI for requests to succeed.

Let me know if you'd like a PDF or markdown version of this.



# startUp

** No parameters **

** Endpoint: **
get /startUp

** response: **

{
  "responseMessage": "Request successful!, composites: 1",
  "managerNames": [
    "first2"
  ]
}
