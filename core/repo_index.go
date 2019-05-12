//
// Copyright © 2017-2019 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package core

import (
	"encoding/xml"
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/db"
	"github.com/getsolus/libeopkg"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// initDistribution will look for the distribution.xml file which will define
// the all-important Obsoletes set
func (r *Repository) initDistribution() error {
	r.dist = nil

	dpath := filepath.Join(r.assetPath, "distribution.xml")
	if !PathExists(dpath) {
		log.Warnf("No distribution.xml defined for repo '%s'\n", r.ID)
		return nil
	}
	dist, err := libeopkg.NewDistribution(dpath)
	if err != nil {
		return err
	}
	r.dist = dist
	return nil
}

// emitDistribution is responsible for loading the distribution.xml file from
// the assets store and merging it into the final index
func (r *Repository) emitDistribution(encoder *xml.Encoder) error {
	elem := xml.StartElement{
		Name: xml.Name{
			Local: "Distribution",
		},
	}
	return encoder.EncodeElement(r.dist, elem)
}

// emitComponents is responsible for loading the components.xml file from
// the assets store and merging it into the final index
func (r *Repository) emitComponents(encoder *xml.Encoder) error {
	dpath := filepath.Join(r.assetPath, "components.xml")
	if !PathExists(dpath) {
		log.Warnf("No components.xml defined for repo '%s'\n", r.ID)
		return nil
	}
	comp, err := libeopkg.NewComponents(dpath)
	if err != nil {
		return err
	}

	elem := xml.StartElement{
		Name: xml.Name{
			Local: "Component",
		},
	}

	for i := range comp.Components {
		c := &comp.Components[i]
		if err := encoder.EncodeElement(c, elem); err != nil {
			return err
		}
	}
	// Now finalise the document
	return nil
}

// emitGroups is responsible for loading the groups.xml file from
// the assets store and merging it into the final index
func (r *Repository) emitGroups(encoder *xml.Encoder) error {
	dpath := filepath.Join(r.assetPath, "groups.xml")
	if !PathExists(dpath) {
		log.Warnf("No groups.xml defined for repo '%s'\n", r.ID)
		return nil
	}
	grp, err := libeopkg.NewGroups(dpath)
	if err != nil {
		return err
	}

	elem := xml.StartElement{
		Name: xml.Name{
			Local: "Group",
		},
	}

	for i := range grp.Groups {
		g := &grp.Groups[i]
		if err := encoder.EncodeElement(g, elem); err != nil {
			return err
		}
	}
	return nil
}

// pushDeltaPackages will insert all applicable (usable) delta packages from
// our repository into the emitted index
func (r *Repository) pushDeltaPackages(db db.Database, pool *Pool, entry *PoolEntry) error {
	// Get our local entry
	repoEntry, err := r.GetEntry(db, entry.Meta.Name)
	if err != nil {
		return err
	}

	var deltas []libeopkg.Delta

	// Find all delta IDs
	for _, id := range repoEntry.Deltas {
		poolEnt, err := pool.GetEntry(db, id)
		if err != nil {
			return err
		}
		if poolEnt.Delta == nil {
			return fmt.Errorf("invalid delta record, corruption: %s", id)
		}

		// Basically the ToRelease must be for this release
		if poolEnt.Delta.ToRelease != entry.Meta.GetRelease() {
			continue
		}

		// Insert delta and clone the pool entry meta for it
		deltas = append(deltas, libeopkg.Delta{
			ReleaseFrom: poolEnt.Delta.FromRelease,
			PackageURI:  poolEnt.Meta.PackageURI,
			PackageSize: poolEnt.Meta.PackageSize,
			PackageHash: poolEnt.Meta.PackageHash,
		})
	}

	// Deltas implicitly sorted by sort.Strings on their IDs
	if deltas != nil && len(deltas) > 0 {
		entry.Meta.DeltaPackages = &deltas
	}

	return nil
}

func (r *Repository) emitIndexPackage(db db.Database, pool *Pool, pkg string, encoder *xml.Encoder, entry *PoolEntry) error {
	// Wrap every output item as Package
	elem := xml.StartElement{
		Name: xml.Name{
			Local: "Package",
		},
	}

	// Retain compatibility with eopkg, auto-drop -dbginfo
	nom := entry.Meta.Name
	if strings.HasSuffix(nom, "-dbginfo") {
		nom = nom[0 : len(nom)-8]
	}

	// Check if its obsolete, if its automatically obsolete through our
	// dbginfo trick, warn in the console
	if r.dist != nil && r.dist.IsObsolete(nom) {
		if nom != entry.Name {
			log.Errorf("Abandoned obsolete package '%s' in repo '%s', please run 'trim obsolete'\n", pkg, r.ID)
		}
		return nil
	}

	// Warn that a package depends on an obsolete package so that it can be
	// purged from the repo (as it won't work!)
	if entry.Meta.RuntimeDependencies != nil && r.dist != nil {
		for _, p := range *entry.Meta.RuntimeDependencies {
			if r.dist.IsObsolete(p.Name) {
				log.Warnf("Encountered uninstallable package '%s' depending on obsolete package '%s' in repo '%s'. Please address.\n", entry.Name, p.Name, r.ID)
			}
		}
	}

	// Shove in the delta packages now
	if err := r.pushDeltaPackages(db, pool, entry); err != nil {
		return err
	}

	return encoder.EncodeElement(entry.Meta, elem)
}

// emitIndex does the heavy lifting of writing to the given file descriptor,
// i.e. serialising the DB repo out to the index file
func (r *Repository) emitIndex(db db.Database, pool *Pool, file *os.File) error {
	var pkgIds []string
	rootBucket := db.Bucket([]byte(DatabaseBucketRepo)).Bucket([]byte(r.ID)).Bucket([]byte(DatabaseBucketPackage))

	err := rootBucket.ForEach(func(k, v []byte) error {
		entry := RepoEntry{}
		if err := rootBucket.Decode(v, &entry); err != nil {
			return err
		}

		if r.dist != nil && r.dist.IsObsolete(entry.Name) {
			return nil
		}

		pkgIds = append(pkgIds, entry.Published)
		return nil
	})

	if err != nil {
		return err
	}

	// Ensure we'll emit in a sane order
	sort.Strings(pkgIds)

	encoder := xml.NewEncoder(file)
	encoder.Indent("    ", "    ")

	// Ensure we have the start element
	if err := encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PISI"}}); err != nil {
		return err
	}

	// Ensure distribution is at the head
	if err := r.emitDistribution(encoder); err != nil {
		return err
	}

	for _, pkg := range pkgIds {
		entry, err := pool.GetEntry(db, pkg)
		if err != nil {
			return err
		}
		if err = r.emitIndexPackage(db, pool, pkg, encoder, entry); err != nil {
			return err
		}
	}

	// Stick in the components
	if err := r.emitComponents(encoder); err != nil {
		return err
	}

	// Stick in the groups ..
	if err := r.emitGroups(encoder); err != nil {
		return err
	}

	// Now finalise the document
	if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PISI"}}); err != nil {
		return err
	}

	return encoder.Flush()
}

// Index will attempt to write the eopkg index out to disk
// This only requires a read-only database view
func (r *Repository) Index(db db.Database, pool *Pool) error {
	r.indexMut.Lock()
	defer r.indexMut.Unlock()
	var errAbort error

	indexPath := filepath.Join(r.path, "eopkg-index.xml.new")
	indexPathFinal := filepath.Join(r.path, "eopkg-index.xml")
	mapping := make(map[string]string)
	mapping[indexPath] = indexPathFinal

	defer func() {
		if errAbort != nil {
			for k := range mapping {
				log.Errorf("Removing potentially corrupt index file '%s' for repo '%s', reason: '%s'\n", k, r.ID, errAbort.Error())
				os.Remove(k)
			}
		}
	}()

	if err := r.initDistribution(); err != nil {
		return err
	}

	// Create index file
	f, err := os.Create(indexPath)
	if err != nil {
		errAbort = err
		return errAbort
	}

	// Write the index file
	errAbort = r.emitIndex(db, pool, f)
	f.Close()
	if errAbort != nil {
		return errAbort
	}

	// Sing the theme tune
	indexPathSha := filepath.Join(r.path, "eopkg-index.xml.sha1sum.new")
	indexPathShaFinal := filepath.Join(r.path, "eopkg-index.xml.sha1sum")
	mapping[indexPathSha] = indexPathShaFinal

	// Star in it
	if errAbort = WriteSha1sum(indexPath, indexPathSha); err != nil {
		return errAbort
	}

	// Write our XZ index out
	indexPathXz := filepath.Join(r.path, "eopkg-index.xml.new.xz")
	indexPathXzFinal := filepath.Join(r.path, "eopkg-index.xml.xz")
	mapping[indexPathXz] = indexPathXzFinal

	if errAbort = libeopkg.XzFile(indexPath, true); errAbort != nil {
		return errAbort
	}

	// Write sha1sum for our xz file
	indexPathXzSha := filepath.Join(r.path, "eopkg-index.xml.xz.sha1sum.new")
	indexPathXzShaFinal := filepath.Join(r.path, "eopkg-index.xml.xz.sha1sum")
	mapping[indexPathXzSha] = indexPathXzShaFinal

	// xz sha1
	if errAbort = WriteSha1sum(indexPathXz, indexPathXzSha); err != nil {
		return errAbort
	}

	for k, v := range mapping {
		if errAbort = os.Rename(k, v); errAbort != nil {
			return errAbort
		}
	}

	return nil
}
