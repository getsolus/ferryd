# API

## Generic Response

A `GenericResponse` is returned any time there is an error during an API request. It consists of an array of strings, one per error.

``` JSON
{
	"errors"    : [strings]
}
```

## /api/v1/status

Reports the status of Ferryd and the `JobStore`

### GET

On success, this endpoint returns a `StatusResponse` which contains the time the daemon started, the version number of `ferryd`, and then lists of all of the most recent jobs. `Current` will contain up to 10 jobs, with the currently running jobs listed first, and queued jobs after. `Failed` will contain up to 10 of the most recently failed jobs. `Completed` will contain up to 10 of the most recently finished jobs.

``` JSON
{
	"time_started" : "2020-12-30T15:51:00Z",
	"version"      : "1.0.0",
	"current"      : [Jobs],
	"failed"       : [Jobs],
	"completed"    : [Jobs]
}
```

## /api/v1/daemon?action=:action

### PATCH

**NOTE:** This endpoint is not yet implemented.

The `daemon` endpoint is intended to be used for things like restarting, reloading, or stopping the daemon. Because the current usage of `ferryd` is under a `systemd` environment, this is not needed. However, the intent of the `action` query parameter is to instruct the API as to which action shoudl be taken.


## /api/v1/repos

### GET

On success, GET will return a list of all of the available repos:

```JSON
[
	{
		"name"     : "pool",
		"packages" : 2398409,
		"deltas"   : 240823049,
		"size"     : 20394029340294
	},
	{
		"name"     : "unstable",
		"packages" : 2398409,
		"deltas"   : 240823049,
		"size"     : 20394029340294
	}
]
```

## /api/v1/repo/:left?instant=:instant&import=:import&clone=:clone

### POST

#### Create

Create a new repo named ":left", and set it to instant transit if ":instant" is `true`. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

#### Clone (clone=:clone)

Create a new repo named ":left" by copying the contents of the repo named ":clone". This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

#### Import (import=:import)

Create a new repo named ":left" by rereading the contents of former repo ":left" from disk. The query argument ":import" must be set for all imports. It will be set for instant transit if ":instant" is `true`.  This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

### PATCH

#### Check (action="check")

Compares the contents of disk with the contents of the database for the repo name ":left" and generates a `repo.Diff` of any inconsistencies. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

#### Delta (action="delta")

Generates any missing delta packages and cleans up old deltas for the repo named ":left" and generates a `repo.Diff` of any inconsistencies. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

#### Index (action="index")

Regenerates the repo Index and the corresponding SHA hashsum files. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

#### Rescan (action="rescan")

Compares the contents of disk with the contents of the database for the repo named ":left" and generates a `repo.Diff` of any inconsistencies. If inconsistencies are found, they are repaired. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

#### Trim Obsoletes (action="trim-obsoletes")

Remove all package archives (deltas included) from the repo named ":left", as indicated in its `distribution.xml` in the Assets directory and generates a `repo.Diff` of any removals. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

#### Trim Packages (action="trim-packages"&max=:max)

Remove all old package archives (deltas included) from the repo named ":left", up to and excluding the ":max" number of relases specified and generates a `repo.Diff` of any removals. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

### DELETE

Deletes all references to the repo named ":left" in the package DB and cleans up any entries in the Pool that are no longer needed. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

## /api/v1/repos/:left/cherrypick/:right?package=":package"

### PATCH

Copies all of the package archives (including deltas) from the repo named ":left" to the repo named ":right" for the package named ":package" and generates a `repo.Diff` of the additions. This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

## /api/v1/repos/:left/compare/:right

### GET

Generates a `repo.Diff` from all of the inconsistencies between the repo named ":left" and the repo named ":right". This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.

## /api/v1/repos/:left/sync/:right

### PATCH

Generates a `repo.Diff` from all of the inconsistencies between the repo named ":left" and the repo named ":right" and then correct all of those inconsistencies in ":right" such that ":right" is then identica to ":left". This creates a job in the `JobStore` which can later be accessed from the JobID specified in the body of the response, e.g.:

```
12345
```

The completed Job will contain the JSON encoded `repo.Diff` in its "results" field.


## /api/v1/jobs?status=":status"

### DELETE

#### Completed Jobs (status="completed")

Removes all jobs from the `JobStore` which have already been completed.

#### Failed Jobs (status="failed")

Removes all jobs from the `JobStore` which previously failed to complete.

#### Queued Jobs (status="queued")

Removes all jobs from the `JobStore` which have not yet been started.


## /api/v1/jobs/:id

### GET

Retrieves a Job matching the integer `:id` as returned from a previous API request. The response body will contain a JSON encoded `jobs.Job` if a job exists for that ID:

```JSON
{
	"id"       : 12345,
	"type"     : 1,
	"src"      : "unstable",
	"dst"      : "shannon",
	"pkg"      : "nano",
	"max"      : 3,
	"created"  : "2020-12-31T11:05:00Z",
	"started"  : "2020-12-31T11:05:00Z",
	"finished" : "2020-12-31T11:05:00Z",
	"status"   : 1,
	"message"  : "Something went terribly wrong",
	"Results"  : "<base64>"
}
```

