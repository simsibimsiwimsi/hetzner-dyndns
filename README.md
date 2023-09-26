# Hetzner DNS Based DynDns Service

A containerized microservice that utilizes [Hetzner DNS Zone](https://dns.hetzner.com/api-docs#tag/Zones) (manage your own at [Hetzner DNS](https://dns.hetzner.com/)) to act as a DynDNS service.
Provides a docker-container which spins up a single endpoint:
```
http://user:password@hetzner-dyndns:8053/?dnsRecordName={dynamicSubdomain}&ipv4={IPv4address}&ipv6=2{IPv6address}
```

## Getting Started

TODO

### Prerequisites

TODO

### Installing

TODO

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Go](https://go.dev/) - Programming language allowing for resource-optimized services
* [Docker](https://docs.docker.com/language/golang/build-images/) - Simple container system

## Contributing

TODO

## Versioning

TODO 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

