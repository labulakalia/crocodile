# Actuator Service

This is the Actuator service

Generated with

```
micro new crocodile/web/actuator --namespace=crocodile --alias=actuator --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: crocodile.web.actuator
- Type: web
- Alias: actuator

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
./actuator-web
```

Build a docker image
```
make docker
```