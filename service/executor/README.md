# Executor Service

This is the Executor service

Generated with

```
micro new github.com/labulaka521/crocodile/service/executor --namespace=crocodile --alias=executor --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: crocodile.srv.executor
- Type: srv
- Alias: executor

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
./executor-srv
```

Build a docker image
```
make docker
```