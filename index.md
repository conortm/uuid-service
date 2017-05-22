# UUID Service [![Build Status](https://img.shields.io/travis/conortm/uuid-service.svg)](https://travis-ci.org/conortm/uuid-service) [![Coverage Status](https://img.shields.io/coveralls/conortm/uuid-service.svg)](https://coveralls.io/r/conortm/uuid-service?branch=master)

A simple service for `PUT`ting & `GET`ting UUIDs, provided a key.

This is a [Go](http://golang.org/) web service, backed by [MongoDB](https://www.mongodb.org/),
running on [Docker](https://www.docker.com/), with [Vagrant](https://www.vagrantup.com/)
for local development.

## How to use this image

Start a [uuid-service](https://hub.docker.com/r/conortm/uuid-service/) instance
and link it to a [MongoDB](https://hub.docker.com/_/mongo/) instance:

```bash
docker run --name some-uuid-service --link some-mongo:mongo -d conortm/uuid-service
```

## How to use this service

Once you have the service running at `http://my-uuid-service.com`, do the following:

To create a new `UUID`, run:

```bash
curl -v -X PUT http://my-uuid-service.com/uuid/my-key
```

**Note:** A response status of `201 Created` indicates a successfully created UUID.
A status of `200 OK` indicates that the UUID already exists for the provided key.

To get an existing `UUID`, run:

```bash
curl -v http://my-uuid-service.com/uuid/my-key
```

## Vagrant

Get a local instance up and running with:

```bash
vagrant up
```

**Note:** Substitute `my-uuid-service.com` above with `localhost:3000`.

## License

[MIT](https://github.com/conortm/uuid-service/blob/master/LICENSE) © [Conor McNamara](http://conortm.io)
