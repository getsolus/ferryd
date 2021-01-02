# Pool

## What is the "pool"?

The pool is a pseudo-repository which contains every archive present in any of the managed ferryd repos.

## When is the "pool" used?

When a new manifest is transited, all of the archives are first added to the pool. From there, all of the
archives are either hardlinked (same FS) or copied (different FS) to any Repo marked for instant transit.

## Why does the "pool" exist?

The pool provides a singular location for creating any hardlinks. This avoids the need to always find out
if two repositories share the same FS, while also saving significant disk space when they do. If all
repos are on different FS, then the Pool at least acts as a local backup of all of the archives, useful
for an emergency recreation of the repositories.

## What maintenance needs to be done to the "pool"?

Any archives no longer used by any of the managed repos may be safely removed from the Pool to save disk
space and keep the Repo DB small.