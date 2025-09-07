# Bookkeeper

A tool that will do all the CBR/CBZ/PDF 

## Commands

### `bookkeeper scan <folder>`

Scan for comic books in the folder, recursively.
Each book file found will be reported as a one-line JSON entry.

> Partially implemented

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