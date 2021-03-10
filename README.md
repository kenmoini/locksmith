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

## Testing

For the purposes of checking the generation of PKI via Locksmith/Golang against a standard OpenSSL generated PKI there are a set of resources to generate and compare the chains.

### 1. Generate OpenSSL PKI Chain

This can easily be done by running the following command:

```bash
./generate_test_pki.openssl.sh
```

With the default settings it will create a PKI chain with a Root CA, Intermediate CA, and Server Certificate with CRL in the `.test_pki_root` directory.

The OpenSSL configuration files used to generate this PKI can be found in the `openssl_extras/` directory.

### 2. Launch Locksmith & Generate PKI Chain

There is also a quick and easy way to generate a comparable chain via Locksmith by running the following:

```bash
./generate_test_pki.locksmith.sh
```

***NOTE***: This requires Locksmith to be available in the local directory - you can build it from source by running `go build`

Running that script will start Locksmith with the `config.yml.example` configuration, listening on port 8080.  It will then run the required cURL requests locally to generate the PKI Chain that is available in the `./.generated` directory.

### 3. Compare PKI Chains

Another script can make your life easier when comparing PKI Chains to ensure the Subject, Issuer, Capabilities, and so on are aligned closely.

```bash
./generate_test_pki.compare.sh
```

This script will compare the two different PKI chains that were generated in the previous two steps.

### Bonus: Bundled Scripts!

You can run all three testing scripts with the following command:

```bash
./generate_test_pki.bundle.sh
```