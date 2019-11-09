# Schema for Job Handling

## Job Types

### Check

#### Description:

    Compares the contents of a repo on Disk with the DB

#### Type:

- Parallel

#### Parameters:

- src

#### Results:

- Diff

#### Created By:

- API Call
- Rescan Job

#### Followed By:

- N/A

---

### CherryPick

#### Description:

    Syncs a single package from one repo to another

#### Type:

- Serial

#### Parameters:

- src
- dst
- pkg

#### Results:

- Diff

#### Created By:

- API Call

#### Followed By:

- Index (dst)

---

### Clone Repo

#### Description:
    Clone one repository into another

#### Type:

- Serial

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Created By:

- API Call

#### Followed By:

- Index (dst)

---

### Compare

#### Description:

    Creates a diff of the contents of two repos

#### Type:

- Parallel

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Created By:

- API Call

#### Followed By:

- N/A

---

### Create

#### Description:
    Create a new repository by name

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- N/A

#### Created By:

- API Call
- Clone (dst)

#### Followed By:

- Index (dst)

---

### Delta

#### Description:

    Generates missing Delta Packages for an entire repo

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- N/A

#### Created By:

- API Call

#### Followed By:

- DeltaPackage (dst, all pkgs)
- Index (dst)

---

### DeltaPackage

#### Description:

    Create delta packages for a single package in a repo

#### Type:

- Parallel

#### Parameters:

- dst
- pkg

#### Results:

- N/A

#### Created By:

- Delta

#### Followed By:

- Index (dst)

---

### Import

#### Description:

    Adds a new repo to the DB from an existing filepath

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- Diff

#### Created By:

- API Call

#### Followed By:

- Index (dst)

---

### Index

#### Description:

    Update the index for a specific repo

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- N/A

#### Created By:

- API Call
- Clone (dst)
- Create (dst)
- Delta (dst)
- Import (dst)
- Rescan (dst)
- Sync (dst)
- Trim Obsoletes (dst)
- Trim Packages (dst)
- Transit Package (dst)

#### Followed By:

- None

---

### Remove

#### Description:

    Remove removes a repo from the DB but not its contents on disk

#### Type:

- Serial

#### Parameters:

- src

#### Results:

- N/A

#### Created By:

- API Call

#### Followed By:

- N/A

---

### Rescan

#### Description:

    Updates the DB with the contents of a repo on disk

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- Diff

#### Created By:

- API Call
- Import (dst)

#### Followed By:

- Index (dst)

---

### Sync

#### Description:

    Replicates the exact contents of one repo into another

#### Type:

- Serial

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Created By:

- API Call
- Clone (src,dst)

#### Followed By:

- Index (dst)

---

### Transit Package

#### Description:

    Adds a new package to the Pool and all auto-transit repos

#### Type:

- Serial

#### Parameters:

- pkg

#### Results:

- N/A

#### Created By:

- Transit Listener

#### Followed By:

- Index (dst)

---

### Trim Obsoletes

#### Description:

    Remove obsoleted packages from the repo

#### Type:

- Serial

#### Parameters:

- dst

#### Results:

- Diff

#### Created By:
- API Call

#### Followed By:

- Index (dst)

---

### Trim Packages

#### Description:

    Remove old releases for packages in a repo

#### Type:

- Serial

#### Parameters:

- dst
- max

#### Results:

- Diff

#### Created By:

- API Call

#### Followed By:
- Index (dst)


## SQLite Schema

| Column Number | 0       | 1       | 2      | 3      | 4      | 5       |
| ------------- | ------- | ------- | ------ | ------ | ------ | ------- |
| Column Name   | id      | type    | src    | dst    | pkg    | max     |
| Column Type   | INTEGER | INTEGER | STRING | STRING | STRING | INTEGER |

| Column Number | 6        | 7        | 8        | 9         | 10      | 11      |
| ------------- | -------- | -------- | -------- | --------- | ------- | ------- |
| Column Name   | created  | started  | finished | status    | message | results |
| Column Type   | DATETIME | DATETIME | DATETIME | INTEGER   | TEXT    | BLOB    |

