# UUID Service

Simple service for creating & getting UUIDs, provided a key.

## Running locally

Get a local instance up and running:

```bash
vagrant up
```

To create a new `UUID`, run:

```bash
curl -XPOST -d '{"key":"test"}' http://localhost:3000/uuid
```

To get an existing `UUID`, run:

```bash
curl http://localhost:3000/uuid/test
```
