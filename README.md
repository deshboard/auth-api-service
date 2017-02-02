# Auth API Service

[![Build Status](https://img.shields.io/travis/deshboard/auth-api-service.svg?style=flat-square)](https://travis-ci.org/deshboard/auth-api-service)


## Prerequisites

- up to date [Docker](https://www.docker.com/) (1.13.0 at the moment)
- up to date [Docker Compose](https://docs.docker.com/compose/) (1.10.0 at the moment)
- [Glide](http://glide.sh/)
- make


## Setup

Before starting any Docker containers you might want to setup your own port forwarding rules based on `docker-compose.override.yml.example`.

``` bash
$ cp docker-compose.override.yml.example docker-compose.override.yml
```

Run the following commands to prepare and start the Docker environment:

``` bash
$ make setup
$ make start
```


## Development

To install Go dependencies locally run the following commands:

``` bash
$ make install
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
