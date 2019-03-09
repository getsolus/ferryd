# Schema for Job Handling

## Job Types


### Bulk Add

#### Description:
    Adds one or more packages to a repository

#### Type:

- Serial

#### Parameters:

- repoID
- Package(s)

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Clone Repo

#### Description:
    Clone one repository into another

#### Type:

- Serial

#### Parameters:

- repoID
- newClone
- mode

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Copy Source

#### Description:
    Copy a particular source from one repo to another

#### Type:

- Serial

#### Parameters:

- repoID
- target
- source
- release

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Create Repo

#### Description:
    Create a new repository by name

#### Type:

- Serial

#### Parameters:

- id

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Delete Repo

#### Description:
    Delete an existing repository by name

#### Type:

- Serial

#### Parameters:

- id

#### Created By:
- API Call

#### Followed By:
- None

---

### Delta

#### Description:
    Create delta packages for a single package in a repo

#### Type:

- Parallel

#### Parameters:

- repoID
- packageID

#### Created By:
- Delta Repo Job

#### Followed By:
- None

---

### Delta Index

#### Description:
    Create delta packages for a single package in a repo and re-index after

#### Type:

- Parallel

#### Parameters:

- repoID
- packageID

#### Created By:
- Transit Process
- Pull Repo Job

#### Followed By:
- Direct Index

---

### Delta Repo

#### Description:
    Create delta packages for all packages in a repo

#### Type:

- Serial

#### Parameters:

- repoID

#### Created By:
- API Call

#### Followed By:
- None

---

### Index Repo

#### Description:
    Update the index for a specific repo

#### Type:

- Serial

#### Parameters:

- repoID

#### Created By:
- API Call

#### Followed By:
- None

---

### Pull Repo

#### Description:
    Pull one repository into another

#### Type:

- Serial

#### Parameters:

- sourceID
- targetID

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Remove Source

#### Description:
    Remove a specific package or release and its sub-packages/deltas

#### Type:

- Serial

#### Parameters:

- repoID
- source
- release

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Transit

#### Description:
    Add a new package or package release to a repo

#### Type:

- Serial

#### Parameters:

- path

#### Created By:
- FS Notify

#### Followed By:
- Direct Index

---

### Trim Obsoletes

#### Description:
    Remove obsoleted packages from the repo

#### Type:

- Serial

#### Parameters:

- repoID

#### Created By:
- API Call

#### Followed By:
- Direct Index

---

### Trim Packages

#### Description:
    Remove old releases for packages in a repo

#### Type:

- Serial

#### Parameters:

- repoID
- maxKeep

#### Created By:
- API Call

#### Followed By:
- Direct Index


## SQLite Schema

| Column Number | 0       | 1        | 2         | 3         | 4       | 5       | 6         | 7       |
| ------------- | ------- | -------- | --------- | --------- | ------- | ------- | --------- | ------- |
| Column Name   | id      | job_type | src\_repo | dst\_repo | sources | release | max\_keep | mode    |
| Column Type   | INTEGER | INTEGER  | STRING    | STRING    | TEXT    | INTEGER | INTEGER   | INTEGER |

| Column Number | 8        | 9        | 10       | 11        | 12      |
| ------------- | -------- | -------- | -------- | --------- | ------- |
| Column Name   | created  | started  | finished | status    | message |
| Column Type   | DATETIME | DATETIME | DATETIME | INTEGER   | TEXT    |

