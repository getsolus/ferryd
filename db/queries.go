package db

const repoPackages = `
WITH ids AS (
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name = :name
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id
`

const sharedPackages = `
WITH ids AS(
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name=:name1 OR repos.name=:name2
    GROUP BY package_id
    HAVING count(*) > 1
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id;
`

const uniquePackages = `
WITH ids AS(
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name=:name1 OR repos.name=:name2
    GROUP BY package_id
    HAVING count(*) = 1 AND repos.name=:name1
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id
`
