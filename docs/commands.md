# Commands

## Check

### Goal(s)

- Verified the contents of disk against the records in the DB.
- Report any and all discrepancies.

### Process

**Client**

1. Client requests a Check for a repo from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any inconsistencies that were detected.


**Daemon**

1. Daemon receives a request to Check a repo.
2. A new Job is created for the Check.
3. A worker is assigned to the Job, when available.
4. The worker requests a Check from the Repo Manager.
5. The Repo Manager requests the Repo Check itself.
6. The Repo compares each Archive record with the contents on disk, also keeping track of files not found in the DB.
7. If the Check fails with an error:
	1. The Error is returned through the Manager to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Check has completed, a Diff is returned.
9. The Manager reciprocates the Diff to the Worker.
10. The Worker encodes the Diff into the Results of the Job and retires it as Completed.


## Cherry-Pick

### Goals

1. Compare the available archives for a package in one repo with the same package in another.
2. Modify the second repo such that it contains an identical set of archives, as found in the first.

### Process

**Client**

1. Client requests a Cherry-Pick between two repos for a specific package from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any modifications that were needed.


**Daemon**

1. Daemon receives a request to Cherry Pick between repos for a package.
2. A new Job is created for the Cherry Pick.
3. A worker is assigned to the Job, when available.
4. The worker requests a Cherry Pick from the Repo Manager.
5. The Repo Manager Compares the first Repo with the second, calculating a Diff for only the requested Package.
6. The Repo Manager modifies the Archives in the second Repo to account for any differences.
7. If the Cherry Pick fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Cherry Pick has completed, a Diff is returned to the Worker.
9. The Worker encodes the Diff into the Results of the Job and retires it as Completed.

## Clone

### Goals

1. Create a Repository.
2. Sync the contents of an existing repository to the newly created repository.

### Process

**Client**

1. Client requests a Clone of an existing repo into a new repo from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any modifications that were needed.


**Daemon**

1. Daemon receives a request to Clone into a new repo.
2. A new Job is created for the Clone.
3. A worker is assigned to the Job, when available.
4. The worker requests a Clone from the Repo Manager.
5. The Repo Manager creates a new Repo.
6. The Repo Manager performs a Sync from the existing Repo to the new Repo, generating a Diff for all of new Archives.
7. If the Clone fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Clone has completed, a Diff is returned to the Worker.
9. The Worker encodes the Diff into the Results of the Job and retires it as Completed.

## Compare

### Goals

1. Determine the differences between two Repos.

### Process

**Client**

1. Client requests a Compare of two repos from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any inconsistencies that were found.


**Daemon**

1. Daemon receives a request to Compare into a new repo.
2. A new Job is created for the Compare.
3. A worker is assigned to the Job, when available.
4. The worker requests a Compare from the Repo Manager.
5. The Repo Manager Compares the first Repo with the second, calculating a Diff for all packages.
7. If the Compares fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Compare has completed, a Diff is returned to the Worker.
9. The Worker encodes the Diff into the Results of the Job and retires it as Completed.

## Create Repo

### Goals

1. Create a Repository.

### Process

**Client**

1. Client requests a Create Repo for a new repo from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.


**Daemon**

1. Daemon receives a request to Create Repo for the new repo.
2. A new Job is created for the Create Repo.
3. A worker is assigned to the Job, when available.
4. The worker requests a Create Repo from the Repo Manager.
5. The Repo Manager creates a new Repo.
6. If the Create Repo fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
7. When the Create Repo has completed, the Worker is notified of the success.
8. The Worker retires the Job as Completed.

## Daemon

### Goals

1. Start a new daemon process with the running client executable.

### Process


**Daemon**

1. Daemon Process begins its startup process.
2. Daemon Process sets up logging.
3. Daemon Process creates Base directory if missing.
4. Daemon Process creates Daemon.
5. Daemon Process tells the Daemon to Bind
	1. Daemon creates a lock file.
		1. If the lock file already exists, the Daemon terminates with an error.
	2. Daemon opens the JobStore.
		1. JobStore creates the Job DB if missing.
		2. JobStore adds any missing DB tables from Schema.
	3. Daemon creates the Repo Manager.
		1. Repo Manager creates the Repo DB if missing.
		2. Repo Manager adds any missing DB tables from Schema.
		3. Repo Manager creates a Worker Pool and starts execution.
		4. Repo Manager creates the "pool" Repo if missing.
	4. Daemon creates a Transit Listener.
		1. Transit Listener creates the transit path if missing.
		2. Transit Listener creates a FS watcher for the transit path.
		3. Transit Listener starts the FS watcher monitoring of the transit path.
	5. Daemon creates the API Listener.
		1. API Listener registers all of the API Endpoints with its Router.
	6. Daemon tells the API Listener to Bind.
		1. If running under systemd, the API Listener opens an existing Unix Socket for communication.
			1. If not, the API Listener creates a new Socket and sets the ownership and permissions.
6. Daemon Process tells the Daemon to begin serving requests.
	1. Daemon starts background listener for termination.
	2. If running under systemd, Daemon notifies systemd that it has started.
	3. Daemon tells the API Listener to Start serving requests.
		1. API Listener begins serving 


## Delta

### Goals

1. Generate any missing Delta Archives.
2. Remove any orphaned Delta Archives.

### Process

**Client**

1. Client requests a Delta for a repo from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any modifications that were made.


**Daemon**

1. Daemon receives a request to Delta a repo.
2. A new Job is created for the Delta.
3. A worker is assigned to the Job, when available.
4. The worker requests a Delta from the Repo Manager.
5. The Repo Manager creates missing Delta Archives for all packages in this repo.
6. The Repo Manager finds all orphaned Delta Archives and removes them from the DB, generates a new Index, the removes the archive files from disk.
7. If the Delta fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Delta has completed, a Diff is returned to the Worker.
9. The Worker encodes the Diff into the Results of the Job and retires it as Completed.

## Import

### Goals

1. Recreate a Repository.
2. Import the contents of disk to the newly created repository.

### Process

**Client**

1. Client requests an Import of an existing repo (disk) into a new repo (DB) from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.
5. If the Job was successful, a Diff will be printed for any modifications that were needed.


**Daemon**

1. Daemon receives a request to Import into a new repo.
2. A new Job is created for the Import.
3. A worker is assigned to the Job, when available.
4. The worker requests an Import from the Repo Manager.
5. The Repo Manager creates a new Repo.
6. The Repo Manager performs a Rescan of the repo on disk, generating a Diff for all of new Archives.
7. If the Import fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Import has completed, a Diff is returned to the Worker.
9. The Worker encodes the Diff into the Results of the Job and retires it as Completed.

## Index

### Goals

1. Update the Repo Index on disk and regenerate SHA sum files.

### Process

**Client**

1. Client requests an Import of an existing repo (disk) into a new repo (DB) from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.


**Daemon**

1. Daemon receives a request to Index an existing repo.
2. A new Job is created for the Index.
3. A worker is assigned to the Job, when available.
4. The worker requests an Index from the Repo Manager.
5. The Repo Manager requests an Index from the Repo.
6. The Repo performs an Index of the Repo, updating the SHA sums as well.
7. If the Index fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
8. When the Index has completed the Worker is notified.
9. The Worker retires the job as Completed.

## List Repo

### Goals

1. Provide a listing of all of the ferryd managed repos.
2. For each Repo provide:
	1. A count of the available Package Archives
	2. A count of the available Delta Archives
	3. An estimate of the total size of the Archives on disk.
	4. The total used space of the filesystem where the Repo is stored.
	5. The total free space of the filesystem where the Repo is stored.

### Process

**Client**

1. Client requests a List Repos of all the repos.
2. Client receives a response with the summary or an error.
3. If a summary was returned, it is printed for the User.

**Daemon**

1. Daemon receives a request for a Repo listing.
2. Daemon requests a listing from the Repo Manager.
3. The Repo Manager requests a summary from each Repo.
6. Each Repo summarizes its contents and checks the available space of its filesystem.
7. If the List Repos fails with an error:
	1. The Error is returned to the Repo Manager.
	2. The API encodes the error in the Response.
8. The API encodes the Summaries into the Response.

## Remove Repo

### Goals

1. Remove a Repository from the ferryd DB.

### Process

**Client**

1. Client requests a Remove Repo for an existing repo from the Daemon.
2. Client receives a response with the Job ID or an error.
3. If there is no errors, the Client periodically requests the Job by ID from the Daemon.
4. When the Job is completed or has failed, the Job is summarized for the User.


**Daemon**

1. Daemon receives a request to Remove Repo for an existing repo.
2. A new Job is created for the Remove Repo.
3. A worker is assigned to the Job, when available.
4. The worker requests a Remove Repo from the Repo Manager.
5. The Repo Manager removes the Repo and all of its packages from the DB.
6. If the Remove Repo fails with an error:
	1. The Error is returned to the Worker.
	2. The Worker encodes the error into the Message of the Job and retires it as Failed.
7. When the Create Repo has completed, the Worker is notified of the success.
8. The Worker retires the Job as Completed.


## Rescan

## Reset Completed

## Reset Failed

## Reset Queued

## Status

## Sync

## Trim Obsoletes

## Trim Packages

## Version