
Repository
  - Packages
    - Releases
      - Deltas

ReleaseSet
  - Release:    Package Release
  - []Releases: Deltas for that Release

ReleaseSet (nano, release 116)
  - Release:   Package nano 116
  - []Release: Deltas
    - Release: 116 from 115
    - Release: 116 from 114
    - Release: 116 from 113

PackageSet
  - string:       Name
  - Repo:         Parent
  - []ReleaseSets Sets

Repo A: PackageSet "nano" <-> Repo B: PackageSet "nano"

Repo A: ReleaseSet (nano, release 116) <-> Repo B: ReleaseSet (nano, release 116)