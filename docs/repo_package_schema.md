# Repository and Package Management Schemas

## Repo Table

| Column Number | 0       | 1      | 2               |
| ------------- | ------- | ------ | --------------- |
| Column Name   | id      | name   | instant_transit |
| Column Type   | INTEGER | STRING | BOOLEAN         |

## Package Table

| Column Number | 0       | 1           | 2      | 3       | 4    | 5       | 6    |
| ------------- | ------- | ----------- | ------ | ------- | ---- | ------- | ---- |
| Column Name   | id      | name        | uri    | size    | hash | release | meta |
| Column Type   | INTEGER | STRING      | STRING | INTEGER | TEXT | INTEGER | BLOB |

## Delta Table

| Column Number | 0       | 1             | 2      | 3       | 4    | 5         | 6         |
| ------------- | ------- | ------------- | ------ | ------- | ---- | --------- | --------- |
| Column Name   | id      | package\_name | uri    | size    | hash | from\_rel | to\_rel   |
| Column Type   | INTEGER | STRING        | STRING | INTEGER | TEXT | INTEGER   | INTEGER   |

## Repo Packages Table

| Column Number | 1       | 2           |
| ------------- | ------- | ----------- |
| Column Name   | repo_id | package\_id |
| Column Type   | INTEGER | INTEGER     |


## Repo Deltas Table

| Column Number | 1       | 2           |
| ------------- | ------- | ----------- |
| Column Name   | repo_id | delta\_id   |
| Column Type   | INTEGER | INTEGER     |

