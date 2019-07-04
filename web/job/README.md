# Log Service

This is the Log service

Generated with

```
micro new github.com/labulaka521/crocodile/web/job --namespace=crocodile --alias=log --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: crocodile.web.log
- Type: web
- Alias: log

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
./log-web
```

Build a docker image
```
make docker
```