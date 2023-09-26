# Hetzner DNS Based DynDns Service

A containerized microservice that utilizes [Hetzner DNS Zone](https://dns.hetzner.com/api-docs#tag/Zones) (manage your own at [Hetzner DNS](https://dns.hetzner.com/)) to act as a DynDNS service.
Provides a docker-container which spins up a single endpoint:
```
http://user:password@hetzner-dyndns:8053/?dnsRecordName={dynamicSubdomain}&ipv4={IPv4address}&ipv6=2{IPv6address}
```

## Getting Started

```
git clone git@github.com:simsibimsiwimsi/hetzner-dyndns.git

go run main.go
```
For the above command to yield a meaningful result, you need to create a _dyndns.yml_ config. Please refer to the __Deployment__ section below for a proper configuration.

### Prerequisites
* AMD64 / x86-64  system architecture (you probably have to build from source for ARM)
* _curl_ installed
* _jq_ installed

### Installing

```
curl -L -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" https://api.github.com/repos/simsibimsiwimsi/hetzner-dyndns/releases/latest | jq .assets[0].browser_download_url | xargs curl -L --output hetzner-dyndns.tar

docker load -i hetzner-dyndns.tar
```

## Deployment

For a docker compose deployment add a service like so:

_docker-compose.yml_
```
services:
  hetzner-dyndns:
    image: hetzner-dyndns:latest
    pull_policy: never
    container_name: hetzner-dyndns
    expose:
      - "8053"
    volumes:
      - type: bind
        source: /your/path/to/hetzner-dyndns.yml
        target: /var/dyndns/dyndns.yml
```

Create a yaml config in _/your/path/to/hetzner-dyndns.yml_
```
hetzner:
  dns: 
    zone-id: "REPLACE_WITH_YOUR_HETZNER_DNS_ZONE_ID"
    auth-api-token: "REPLACE_WITH_YOUR_HETZNER_DNS_AUTH_API_TOKEN"

users:
  REPLACE_WITH_YOUR_SUBDOMAIN_1:
    user: REPLACE_WITH_YOUR_USER_1
    password: REPLACE_WITH_BCRYPT_HASH_OF_USER_PASSWORD_1
  REPLACE_WITH_YOUR_SUBDOMAIN_2:
    user: REPLACE_WITH_YOUR_USER_2
    password: REPLACE_WITH_BCRYPT_HASH_OF_USER_PASSWORD_2
```

## Built With

* [Go](https://go.dev/) - Programming language allowing for resource-optimized services
* [Docker](https://docs.docker.com/language/golang/build-images/) - Simple container system

## Contributing

Have not thought about it, yet.

## Versioning

No specific approach to versioning defined. Latest should always be greatest. 
Please refer to https://github.com/simsibimsiwimsi/hetzner-dyndns/releases for a list of releases. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
