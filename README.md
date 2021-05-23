# CRDT
Go implementation of some [CRDTs](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type)

## WIP
This is still a work in progress in the early stage. A lot can (and probably) will change. Currently, the implemented CRDTs are:
- GCounter: grow only counter
- GSet: grow only set
- Counter: a regular counter (increase and decrease)

The next in my pipeline are:
- Set: a regular set
- Register: a register (key-value pair)
- Dictionary: map of key-value pair

A How-to will follow as soon as I am satisfied with the initial implementation.