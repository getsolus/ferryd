# CLI

- [ ] check
- [ ] cherry-pick
- [ ] clone
- [ ] compare
- [x] create-repo
- [x] daemon
- [ ] delta
- [x] help
- [ ] import
- [ ] index
- [x] list-repo
- [x] remove-repo
- [ ] rescan
- [x] reset-completed
- [x] reset-failed
- [x] reset-queued
- [x] status
- [ ] sync
- [ ] trim-obsoletes
- [ ] trim-packages
- [x] version

# Repos --- In-Progress

## Check

- [ ] Check Repo

## Import

- [x] Create Repo
- [ ] Rescan Repo
- [ ] Index Repo

## Rescan

- [ ] Rescan Repo
  - [ ] Check Repo
  - [ ] Import Missing Packages

## Compare

- [ ] Full Repo Diff

## Cherry-Pick

- [ ] Single Package Sync
  - [x] Single Package Diff
  - [x] Remove a specific package from the DB
  - [ ] Remove a specific package from disk
  - [ ] Link a package between repos

## Sync

- [ ] Full Repo Sync
  - [ ] Full Repo Diff
  - [ ] Single Package Sync

## Clone

- [x] Create Repo
- [ ] Full Repo Sync

## Delta

- [ ] Full Repo Delta
  - [ ] Single Package Delta

## Transit Package

- [ ] Add Package
  - [x] Adding Package to the Database
  - [ ] Adding Package to disk
- [ ] Single Package Delta
  - [x] Adding deltas to the DB
  - [ ] Adding Deltas to disk

## Trim Obsoletes

- [ ] Remove Package
    - [ ] Remove Release in a Repo

## Trim Packages

- [ ] Remove Release in a Repo
- [ ] Remove Release on Disk
  - [ ] Remove Package on Disk
  - [ ] Remove Repo on Disk

# Repos --- Done

## List
- [x] Summarize Repos
  - [x] Summarize Single Repo

## Create

- [x] Create Repo
  - [x] Adding a repo to DB
  - [x] Adding a repo to disk
  - [x] Create missing repos directory
  - [x] Use Create to add "pool" when missing
  - [x] Testing

## Remove

- [x] Remove Repo
  - [x] Remove a repo from the DB
  - [x] Remove all package for a repo from the DB
  - [x] Testing

