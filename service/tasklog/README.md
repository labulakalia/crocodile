# Tasklog Service

This is the Tasklog service

Generated with

```
micro new crocodile/service/tasklog --namespace=crocodile --alias=tasklog --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: crocodile.srv.tasklog
- Type: srv
- Alias: tasklog

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
./tasklog-srv
```

Build a docker image
```
make docker
```