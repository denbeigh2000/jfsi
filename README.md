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
 - Tools for adding/removing nodes to an existing cluster + providing rebalancing tools
 - Configuration nodes
 - Configuration file support for applications
 - Some way to account for node failures in replication (either hinted handoff, some simple metadata storage in coordinator for later fulfillment)
