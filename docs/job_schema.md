# Schema for Job Handling

## Job Types

### Check

#### Description:

    Compares the contents of a repo on Disk with the DB

#### Parameters:

- src

#### Results:

- Diff

#### Used By:

- Rescan Job

#### Followed By:

- N/A

---

### CherryPick

#### Description:

    Syncs a single package from one repo to another

#### Parameters:

- src
- dst
- pkg

#### Results:

- Diff

#### Followed By:

- Index (dst)

---

### Clone Repo

#### Description:
    Clone one repository into another

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Followed By:

- Index (dst)

---

### Compare

#### Description:

    Creates a diff of the contents of two repos

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Followed By:

- N/A

---

### Create

#### Description:
    Create a new repository by name

#### Parameters:

- dst

#### Results:

- N/A

#### Used By:

- Clone (dst)

#### Followed By:

- Index (dst)

---

### Delta

#### Description:

    Generates missing Delta Packages for an entire repo

#### Parameters:

- dst

#### Results:

- N/A

#### Followed By:

- Index (dst)

---

### Import

#### Description:

    Adds a new repo to the DB from an existing filepath

#### Parameters:

- dst

#### Results:

- Diff

#### Followed By:

- Index (dst)

---

### Index

#### Description:

    Update the index for a specific repo

#### Parameters:

- dst

#### Results:

- N/A

#### Used By:

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

#### Parameters:

- src

#### Results:

- N/A

#### Followed By:

- N/A

---

### Rescan

#### Description:

    Updates the DB with the contents of a repo on disk

#### Parameters:

- dst

#### Results:

- Diff

#### Used By:

- Import (dst)

#### Followed By:

- Index (dst)

---

### Sync

#### Description:

    Replicates the exact contents of one repo into another

#### Parameters:

- src
- dst

#### Results:

- Diff

#### Used By:

- Clone (src,dst)

#### Followed By:

- Index (dst)

---

### Transit Package

#### Description:

    Adds a new package to the Pool and all auto-transit repos

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

#### Parameters:

- dst

#### Results:

- Diff

#### Followed By:

- Index (dst)

---

### Trim Packages

#### Description:

    Remove old releases for packages in a repo

#### Parameters:

- dst
- max

#### Results:

- Diff

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

