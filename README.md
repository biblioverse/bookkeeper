# Bookkeeper

A tool that will do all the CBR/CBZ/PDF 

## Commands

### `bookkeeper scan <folder>`

Scan for comic books in the folder, recursively.
Each book file found will be reported as a one-line JSON entry.

> Partially implemented

```bash
❯ ./bookkeeper scan fixtures
Scanning: .../fixtures
{"path":"Full of Fun/Full_Of_Fun_001__c2c___1957___ABPC_.cbr","status":"success","size":15666637,"hash":"","book":{"title":"Full_Of_Fun_001__c2c___1957___ABPC_","pages":36}}
{"path":"Full of Fun/Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz","status":"success","size":44292901,"hash":"","book":{"title":"Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc","pages":37}}
{"path":"testfile.pdf","status":"success","size":6012,"hash":"","book":{"title":"Title of the Book","pages":1,"authors":["The Author"],"keywords":["book","fantasy"]}}
Scanned all files
```

### `bookeeper extract <book> <extractTo>`

Extract all pages to a folder, will also create a `pages.json` to list all pages

```bash
 ./bookkeeper extract fixtures/Full\ of\ Fun/Full_Of_Fun_001__c2c___1957___ABPC_.cbr out
Extraction complete. 36 files extracted to .../out
```

```json
[
  "Full Of Fun 001 (c2c) (1957) (ABPC)/FullF101.jpg",
  "Full Of Fun 001 (c2c) (1957) (ABPC)/FullF102.jpg",
  "Full Of Fun 001 (c2c) (1957) (ABPC)/FullF103.jpg",
  // ...
]
```