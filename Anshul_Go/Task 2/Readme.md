# Assignment: Build an Application Protocol over TCP/UDP

## Objective

Build a simple application-layer protocol on top of raw TCP and UDP sockets.
---

## Requirements

### 1. Design a Custom Protocol

Create a request-response protocol that supports multiple operations (similar in spirit to HTTP).

At minimum, support:

* `GET` – retrieve data/resource from the server
* `POST` – send data/resource to the server

You may define additional methods if needed.

Example request formats:

```text
GET /file.txt
```

```text
POST /file.txt

<file contents>
```

Example responses:

```text
STATUS 200

<requested data>
```

```text
STATUS 404

Resource not found
```

The exact protocol format is up to you. However, it should be documented and consistently followed by both client and server.

---

### 2. File Transfer Support

Use the protocol to transfer files between devices.

Example operations:

* Upload a file to the server
* Download a file from the server

The protocol should clearly specify:

* Request format
* Response format
* How file data is transmitted

---

### 3. TCP Implementation

Implement the protocol over TCP sockets.

The client should be able to communicate with the server using the protocol you designed and perform all supported operations.

---