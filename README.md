# jfsi

> **j**ust **f**ucking **s**tore **i**t

This software is still in a development stage.

jfsi is a scalable, distributed blob storage engine with a RESTful API.

### What is jfsi?
jfsi is a very simple key-value blob store which supports CRUD
operations. The blobs are stored on a configurable set of storage
nodes which are running the storage http server.

jfsi is not designed to be a user-facing application, rather sit beneath
your own application and serve as a storage engine. I built it because
I wanted something that could store large blobs like S3, but provide a
simple key-value interface - something that would allow me to
"just fucking store it".

It does have:
 - RESTful HTTP API
 - Replication
 - Sharding
 - Chunked storage
 - Dynamically-scalable frontend

It does not:
 - Store metadata
 - Provide any kind of authentication
 - Provide any hierarchical directory structure

### Structure

Response load is distributed between a number of application nodes - these
are horizontally scalable and can have more added/removed at any time.

Storage load is distributed between a number of storage nodes - these
need to be rebalanced when changing the number of nodes in a pool.

(Soon) A number of metadata nodes store information about the blobs, namely
the mapping of blob uuid to chunk uuid. The metadata nodes are the source of
truth for whether a node exists, and the chunk uuid mapping will
deterministically map the nodes in the cluster that the chunk/s can be found
on. This functionality is currently implemented with Redis as the backend.

(Soon) Configuration is spread by controller nodes that serve JSON over HTTP,
which can be marshaled into an `application.StorageConfig` type.

### Usage:

#### API

| **Method**    | **Endpoint**      | **Description**   |
|---------------|-------------------|-------------------|
| POST          | /                 | Upload a blob     |
| PUT           | /&lt;blobID&gt;   | Update a blob     |
| GET           | /&lt;blobID&gt;   | Download a blob   |
| DELETE        | /&lt;blobID&gt;   | Delete a blob     |

Stuff that is still 
- Tools for adding/removing nodes to an existing cluster + providing rebalancing tools
- Configuration nodes
- Configuration file support for applications
- Some way to account for node failures in replication (either hinted handoff, some simple metadata storage in coordinator for later fulfillment)


#### Binaries
Storage node:
```
storage-http -port 8000
```

Application node:
```
application-http -port 8080`
```

### TODO
 - Caching layers around store/application interfaces (wrap-around any implementation)
 - Metadata store for storing chunk info (shard/replicate using same manner)
 - Periodic health-check polling/mark node unhealthy in http clients
 - Periodic check (from FE nodes? Coordinator nodes?) that the storage nodes are in
   sync with metadata nodes, and properly replicated
 - Tools for adding/removing nodes to an existing cluster + providing rebalancing tools
 - Configuration nodes
 - Configuration file support for applications
 - Some way to account for node failures in replication (either hinted handoff, some simple metadata storage in coordinator for later fulfillment)
