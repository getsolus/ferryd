# Repository and Package Management Schemas

## Repo Table

| Column Number | 0       | 1      | 2               |
| ------------- | ------- | ------ | --------------- |
| Column Name   | id      | name   | instant_transit |
| Column Type   | INTEGER | STRING | BOOLEAN         |

## Release Table

| Column Number | 0       | 1       | 2      | 3       | 4    | 5       | 6             | 7    |
| ------------- | ------- | ------- | ------ | ------- | ---- | ------- | ------------- | ---- |
| Column Name   | id      | package | uri    | size    | hash | release | from\_release | meta |
| Column Type   | INTEGER | STRING  | STRING | INTEGER | TEXT | INTEGER | INTEGER       | BLOB |

### "release"

For a Package, the "release" column is the release number for this version of the package.
For a Delta, the "release" column is the release number for the version of the package it will install.

### "from\_release"

For a Package, the "from\_release" column is not used and should be set to NULL.
For a Delta, the "from\_release" is the release number for the version of the package it will replace.

### Meta

For a Package, the "meta" column includes a copy of the metadata needed to generate the index.
For a Delta, the "meta" column is not used.

### Exemplars

| id  | package | uri                                  | size    | hash | release | from\_release | meta |
| --- | ------- | ------------------------------------ | ------- | ---- | ------- | ------------- | ---- |
| 1   | nano    | n/nano-4.5-116-1-x86\_64.eopkg       | 463356  | HASH | 116     | NULL          | BLOB |
| 2   | nano    | n/nano-116-115-1-x86\_64.delta.eopkg | 463355  | HASH | 116     | 115           | NULL |
| 3   | nano    | n/nano-116-114-1-x86\_64.delta.eopkg | 463354  | HASH | 116     | 114           | NULL |
| 4   | nano    | n/nano-116-113-1-x86\_64.delta.eopkg | 463353  | HASH | 116     | 113           | NULL |

## Packages Table

| Column Number | 1        | 2           |
| ------------- | -------- | ----------- |
| Column Name   | repo\_id | release\_id |
| Column Type   | INTEGER  | INTEGER     |
