//
// Copyright © 2017 Solus Project
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

package libdb

import (
	"sync"
)

// DbForeachFunc is used in the root (untyped buckets)
type DbForeachFunc func(key, val []byte) error

// A Closable is a handle or database that can be closed
type Closable interface {
	// Close the database
	Close()
}

// Database is the opaque interface to the underlying database implementation
type Database interface {

	// Close the database (might no-op)
	Close()

	// Return a subset of the database for usage
	Bucket(id []byte) Database

	// Put an object into storage (unique key)
	PutObject(id []byte, o interface{}) error

	// Get an object from storage
	GetObject(id []byte, o interface{}) error

	// Attempt to decode the input into the given output pointer
	Decode(input []byte, o interface{}) error

	// For every key value pair, run the given function
	ForEach(f DbForeachFunc) error
}

// Private helper to add sync locks to the interfaces
type closable struct {
	closed bool
	mut    *sync.Mutex
}

func (c *closable) initClosable() {
	c.closed = false
	c.mut = &sync.Mutex{}
}

func (c *closable) close() bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.closed {
		return false
	}
	c.closed = true
	return true
}

// Open will return an opaque representation of the underlying database
// implementation suitable for usage within ferryd
func Open(path string) (Database, error) {
	// For now we're just using leveldb
	return newLevelDBHandle(path)
}