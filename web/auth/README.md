# Auth Service

This is the Auth service

Generated with

```
micro new crocodile/web/auth --namespace=crocodile --alias=auth --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: crocodile.web.auth
- Type: web
- Alias: auth

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./auth-web
```

Build a docker image
```
make docker
```