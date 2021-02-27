# Locksmith - PKI over an API

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