# Repository and Package Management Schemas

## Repo Table

| Column Number | 0       | 1      | 2               |
| ------------- | ------- | ------ | --------------- |
| Column Name   | id      | name   | instant_transit |
| Column Type   | INTEGER | STRING | BOOLEAN         |

## Link Table

| Column Number | 1        | 2           |
| ------------- | -------- | ----------- |
| Column Name   | repo\_id | archive\_id |
| Column Type   | INTEGER  | INTEGER     |

## Archive Table

| Column Number | 0       | 1       | 2      | 3       | 4    | 5       | 6           |
| ------------- | ------- | ------- | ------ | ------- | ---- | ------- | ----------- |
| Column Name   | id      | package | uri    | size    | hash | release | to\_release |
| Column Type   | INTEGER | STRING  | STRING | INTEGER | TEXT | INTEGER | INTEGER     |

### "release"

For a Package, the "release" column is the release number for this version of the package.
For a Delta, the "release" column is the release number for the version of the package it will install.

### "to\_release"

For a Package, the "to\_release" column is not used and should be set to NULL.
For a Delta, the "to\_release" is the release number for the version of the package will replace it.

### Exemplars

| id  | package | uri                                  | size    | hash | release | to\_release |
| --- | ------- | ------------------------------------ | ------- | ---- | ------- | ----------- |
| 1   | nano    | n/nano-4.5-116-1-x86\_64.eopkg       | 463356  | HASH | 116     | NULL        |
| 2   | nano    | n/nano-116-117-1-x86\_64.delta.eopkg | 463355  | HASH | 116     | 117         |
| 3   | nano    | n/nano-116-118-1-x86\_64.delta.eopkg | 463354  | HASH | 116     | 118         |
| 4   | nano    | n/nano-116-119-1-x86\_64.delta.eopkg | 463353  | HASH | 116     | 119         |

