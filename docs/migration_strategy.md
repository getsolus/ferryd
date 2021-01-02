# Migrating to Ferryd 1.0.0

1. rsync "unstable" to the new package server
    - dodges any "goofiness" with the existing "shannon"
2. import "unstable" to populate the pool and DB
    - avoids needed to do some sort of remote import
3. delta "unstable" to remove orphaned deltas and generate missing deltas
    - reduce the number of objects for the clone, while also freeing up space from unused deltas
    - end up with a "pristine" copy of "unstable"
4. clone new "shannon" from "unstable"
    - fully dodges a inconsistency from the very start
5. sync "unstable" to "shannon" as a sanity check
