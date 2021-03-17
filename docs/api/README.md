# Locksmith APIDocs

Where full URLs are provided in responses they will be rendered as if the service is running on 'http://pkiserver/'.

The base of the API URL path is variably defined via the Configuration YAML - in these documents it's represented as `/locksmith`.

All endpoints are open - authentication is handled by an external API Gateway and is outside of the scope of Locksmith.

## Root Certificate Authorities

* [List Root Certificate Authorities](roots/get.md) : `GET /locksmith/roots`
* [Create New Root CA](roots/post.md) : `POST /locksmith/roots`

## Intermediate Certificate Authorities

The PKI chain managed by Locksmith is theoretically unlimited in how many Intermediate CAs you could chain along - the first limitation would likely be your file system.

When operating against Intermediate CAs there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

- Example Labs Root CA
  - Example Labs Intermediate CA
    - Example Labs Server Signing CA

The CommonName chain would be represented as: `Example Labs Root CA/Example Labs Intermediate CA/Example Labs Server Signing CA`
The slugged CommonName chain (what is stored in the filesystem) would be: `example-labs-root-ca/example-labs-intermediate-ca/example-labs-server-signing-ca`

* [List Intermediate Certificate Authorities](intermediates/get.md) : `GET /locksmith/intermediates`
* [Create New Intermediate Certificate Authority](intermediates/post.md) : `POST /locksmith/intermediates`

---

## Key Pairs

Key Pairs provide key pair management outside of the scope of x509 PKI - this is useful when you want key pairs for CSRs and Clients.

* [List Key Pairs](keys/get.md) : `GET /locksmith/keys`
* [Create New Key Pairs](keys/post.md) : `POST /locksmith/keys`