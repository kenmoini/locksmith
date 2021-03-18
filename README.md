# Locksmith - PKI over an API

[![Container Repository on Quay](https://quay.io/repository/kenmoini/locksmith/status "Container Repository on Quay")](https://quay.io/repository/kenmoini/locksmith) [![release](https://github.com/kenmoini/locksmith/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/kenmoini/locksmith/actions/workflows/release.yml) [![Tests](https://github.com/kenmoini/locksmith/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kenmoini/locksmith/actions/workflows/test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/kenmoini/locksmith.svg)](https://pkg.go.dev/github.com/kenmoini/locksmith) [![Go Report Card](https://goreportcard.com/badge/github.com/kenmoini/locksmith)](https://goreportcard.com/report/github.com/golang-standards/project-layout) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://github.com/kenmoini/locksmith/tree/main/LICENSE)

Locksmith is a simple Golang application, which when supplied a `config.yml` file will start a RESTful API via an HTTP server that will allow the management of Public Key Infrastructure.

- [API Documentation](https://github.com/kenmoini/locksmith/tree/main/docs/api)
- [Contribution Guide](https://github.com/kenmoini/locksmith/tree/main/docs/CONTRIBUTING.md)
- [Basic Usage](#how-to-use-locksmith)
- [Deployment Options](#deployment-options)
- [FAQs](#faqs)
- [Testing](#testing)

## How to Use Locksmith

```bash
$ ./locksmith [-config file]
```

### 1. Generate the Locksmith `config.yml` file

A sample `config.yml` can be found in this repository at [configs/config.yml.example](https://github.com/kenmoini/locksmith/tree/main/configs/config.yml.example)

### 2. Run Locksmith

Running Locksmith will do the following:

1. Create a PKI Base Directory
2. Start an HTTP Server
3. Respond to requests, generating & serving authorities/certificates/keys/requests/revocations.

### 3. Make Requests to the API

The API is served at the HTTP endpoint base path as defined in the configuration YAML.

You can find the API documentation in the [docs/apis/](https://github.com/kenmoini/locksmith/tree/main/docs/api) folder.

---

## Deployment Options

You can run Locksmith on almost any system due to it being a simple Golang binary.  There are also resources to build a container easily, or you could alternatively pull it from Quay.

### Deployment - From Releases *[First release coming soon...]*

You can pull the pre-built application from the [GitHub Releases page](https://github.com/kenmoini/locksmith/releases) of this repository.

### Deployment - As a SystemD Service

Since Locksmith is a binary that takes minimal execution time configuration it's perfect for use with SystemD as a persistent service that can start at system boot.

You can find an example SystemD service file located at [init/locksmith.service](https://github.com/kenmoini/locksmith/tree/main/init/locksmith.service).  You can install it on a system with the following:

```bash
# Get the service file
sudo wget -O /etc/systemd/system/locksmith.service https://raw.githubusercontent.com/kenmoini/locksmith/main/init/locksmith.service
# Reload SystemD
sudo systemctl daemon-reload
# Enable/start the service - configuration required at /etc/locksmith/config.yml and Locksmith binary in $PATH
sudo systemctl enable locksmith
sudo systemctl start locksmith
```

### Deployment - As a Container

Locksmith comes with a `Containerfile` that can be built with Docker or Podman with the following command:

```bash
# Build the container
podman build -f Containerfile -t locksmith .
# Create the config
mkdir container-config
cp configs/config.yml.example container-config/config.yml
# Run the container
podman run -p 8080:8080 -v container-config/:/etc/locksmith locksmith
```

If you prefer to just use a pre-built container you can pull it from Quay via the following:

```bash
# Optional, pre-pull the image
podman pull quay.io/kenmoini/locksmith
# Create the config
mkdir container-config
cp configs/config.yml.example container-config/config.yml
# Run the container
podman run -p 8080:8080 -v container-config/:/etc/locksmith quay.io/kenmoini/locksmith
```

### Deployment - A Container...as a SystemD Service

In case you'd like to run the Container-as-a-Service for the features containerization provides, you can utilize the resources named `init/caas-*`

The Podman container runtime is used by default - you can change it if you'd like in the [init/caas-locksmith-vars.sh](https://github.com/kenmoini/locksmith/tree/main/init/caas-locksmith-vars.sh) file.

The CaaS launcher also provides the ability for assigning resource limits, static networking, and other functions easily - modify for your needs prior to deploying.

```bash
# Get the Container-as-a-Service service files
sudo wget -O /etc/systemd/system/locksmith.service https://raw.githubusercontent.com/kenmoini/locksmith/main/init/caas-locksmith.service
sudo wget -O /etc/locksmith/caas-start.sh https://raw.githubusercontent.com/kenmoini/locksmith/main/init/caas-locksmith-start.sh
sudo wget -O /etc/locksmith/caas-stop.sh https://raw.githubusercontent.com/kenmoini/locksmith/main/init/caas-locksmith-stop.sh
sudo wget -O /etc/locksmith/caas-vars.sh https://raw.githubusercontent.com/kenmoini/locksmith/main/init/caas-locksmith-vars.sh

# Reload SystemD
sudo systemctl daemon-reload
# Enable/start the service - configuration required at /etc/locksmith/config.yml
sudo systemctl enable locksmith
sudo systemctl start locksmith
```

### Deployment - Building From Source

Since this is just a Golang application, as long as you have Golang v1.15+ then the following commands will do the job:

```bash
# Create the config
cp configs/config.yml.example config.yml
# Build the application (Golang 1.15+)
make build
# Run the application
./locksmith
```

### Deployment - To Kubernetes

Of course since this is easily containerized and a tasty Golang binary, it's very easy to deploy on Kubernetes and can scale with little effort.

In the [deploy/kubernetes](https://github.com/kenmoini/locksmith/tree/main/deploy/kubernetes) directory you can find a set of manifests to deploy onto a standard Kubernetes + Nginx Ingress + CertManager cluster.

```bash
# Apply the manifests
kubectl apply -f deploy/kubernetes/
```

---

## FAQs

- **Does this include any sort of authentication, rate limiting, etc?**

  No, that's the job of an API Gateway - this is more of a microservice so manage and secure accordingly.

- **Has this been architected for multi-tenancy?**

  Multiple root certificates and trusted signers?  Yeah, sure.

  Multiple customers/entities/non-trusted orgs? That's a horrible idea, so: no.  This is a small binary service that is deployed first-class via containers, authenticated at an API Gateway, easily scaled out in a Kubernetes cluster.  So your multi-tenancy would be better set at the PaaS layer with different namespaces/PVs/SAs/etc.

---

## Testing

For the purposes of checking the generation of PKI via Locksmith/Golang against a standard OpenSSL generated PKI there are a set of resources to generate and compare the chains.

### 1. Generate OpenSSL PKI Chain

This can easily be done by running the following command:

```bash
./scripts/generate_test_pki.openssl.sh
```

With the default settings it will create a PKI chain with a Root CA, Intermediate CA, and Server Certificate with CRL in the `.test_pki_root` directory.

The OpenSSL configuration files used to generate this PKI can be found in the `/openssl_extras/` directory.

### 2. Launch Locksmith & Generate PKI Chain

There is also a quick and easy way to generate a comparable chain via Locksmith by running the following:

```bash
./scripts/generate_test_pki.locksmith.sh
```

***NOTE***: This requires Locksmith to be available in the local directory - you can build it from source by running `make build`

Running that script will start Locksmith with the `configs/config.yml.example` configuration, listening on port 8080.  It will then run the required cURL requests locally to generate the PKI Chain that is available in the `./.generated` directory.

### 3. Compare PKI Chains

Another script can make your life easier when comparing PKI Chains to ensure the Subject, Issuer, Capabilities, and so on are aligned closely.

```bash
./scripts/generate_test_pki.compare.sh
```

This script will compare the two different PKI chains that were generated in the previous two steps.

### Bonus: Bundled Scripts!

You can run all three testing scripts with the following command:

```bash
./scripts/generate_test_pki.bundle.sh
```
