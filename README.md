# UUID Service [![Build Status](https://img.shields.io/travis/conortm/uuid-service.svg)](https://travis-ci.org/conortm/uuid-service) [![Coverage Status](https://img.shields.io/coveralls/conortm/uuid-service.svg)](https://coveralls.io/r/conortm/uuid-service?branch=master)

A simple service for creating & getting UUIDs, provided a key.

## Running locally

Get a local instance up and running with [Vagrant](https://www.vagrantup.com/):

```bash
vagrant up
```

To create a new `UUID`, run:

```bash
curl -v -X PUT http://localhost:3000/uuid/my-key
```

**Note:** A status of `201 Created` indicates a successfully created UUID. A status of `200 OK`
indicates that the UUID already exists for the provided key.

To get an existing `UUID`, run:

```bash
curl -v http://localhost:3000/uuid/my-key
```
