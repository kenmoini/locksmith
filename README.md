# Locksmith - PKI over an API

[![Tests](https://github.com/kenmoini/locksmith/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kenmoini/locksmith/actions/workflows/test.yml) [![release](https://github.com/kenmoini/locksmith/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/kenmoini/locksmith/actions/workflows/release.yml)

Locksmith is a simple Golang application, which when supplied a `config.yml` file will start a RESTful API via an HTTP server that will allow the management of Public Key Infrastructure.

## How to Use Locksmith

```bash
$ ./locksmith [-config file]
```

### 1. Generate the Locksmith `config.yml` file

A sample `config.yml` can be found in this repository at [config.yml.example](config.yml.example)

### 2. Run Locksmith

Running Locksmith will do the following:

1. Create a PKI Root Directory
2. Start an HTTP Server
3. Respond to requests, serve certificates and requests

### 3. Make Requests to the API

The API is served at the HTTP endpoint base path as defined in the configuration YAML.

*OpenAPI Spec v3 coming soon...*

## Deployment - As a Container

Locksmith comes with a `Containerfile` that can be built with Docker or Podman with the following command:

```bash
podman build -f Containerfile -t locksmith .
podman run -p 8080:8080 -v config/:/etc/locksmith locksmith
```

## Deployment - Building From Source

Since this is just a Golang application, as long as you have Golang v1.15+ then the following commands will do the job:

```bash
cp config.yml.example config.yml
go build
./locksmith
```

## FAQs

- **Does this include any sort of authentication, rate limiting, etc?**

  No, that's the job of an API Gateway - this is more of a microservice so manage and secure accordingly.