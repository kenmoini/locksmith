# Locksmith APIDocs

Where full URLs are provided in responses they will be rendered as if the service is running on `http://$PKI_SERVER/`.

This way you can `export PKI_SERVER="localhost:8080"` or wherever it is your PKI server is and run the cURL examples in your terminal directly from the API Documentation.

The base of the API URL path is variably defined via the Configuration YAML - in these documents it's represented as `/locksmith`.

All endpoints are open - authentication is handled by an external API Gateway and is outside of the scope of Locksmith.

Make note of pluralism - for example, you query a list of Certificates but will request information about a single Certificate.  The API endpoints reflect this pluralism.

- [Root Certificate Authorities](#root-certificate-authorities)
- [Intermediate Certificate Authorities](#intermediate-certificate-authorities)
- [Authority](#authority)
- [Certificate Requests](#certificate-requests)
- [Certificate Request](#certificate-request)
- [Certificates](#certificates)
- [Certificate](#certificate)
- Renewals
- [Certificate Revocations](#certificate-revocations)
- [Key Pairs](#key-pairs)
- [Key Stores](#key-stores)

## Root Certificate Authorities

* [List Root Certificate Authorities](roots/get.md) : `GET /locksmith/roots`
* [Create New Root CA](root/post.md) : `POST /locksmith/root`

## Intermediate Certificate Authorities

The PKI chain managed by Locksmith is theoretically unlimited in how many Intermediate CAs you could chain along - the first limitation would likely be your file system.

When operating against Intermediate CAs there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

- Example Labs Root CA
  - Example Labs Intermediate CA
    - Example Labs Signing CA

The CommonName chain would be represented as: `Example Labs Root CA/Example Labs Intermediate CA/Example Labs Signing CA`

The Slugged CommonName chain (what is stored in the filesystem) would be: `example-labs-root-ca/example-labs-intermediate-ca/example-labs-signing-ca`

You can address the CA Path with either the CommonName Chain or Slugged CommonName Chain - you could even mix and match since input is slugged anyway.

* [List Intermediate Certificate Authorities](intermediates/get.md) : `GET /locksmith/intermediates`
* [Create New Intermediate Certificate Authority](intermediate/post.md) : `POST /locksmith/intermediate`

## Authority

Authority is the structure used to read any Certificate Authority, Root or Intermediate, up and down a CA Path.

* [Read Certificate Authority](authority/get.md) : `GET /locksmith/authority`

## Certificate Requests

* [List Certificate Requests](certificate-requests/get.md) : `GET /locksmith/certificate-requests`

## Certificate Request

* [Read Certificate Request](certificate-request/get.md) : `GET /locksmith/certificate-request`
* [Create New Certificate Request](certificate-request/post.md) : `GET /locksmith/certificate-request`

## Certificate Revocations

Certificate Revocations provides reading of a Certificate Authority's Certificate Revocation List

* [Read Certificate Authority CRL](revocations/get.md) : `GET /locksmith/revocations`

---

## Key Pairs

Key Pairs provide key pair management outside of the scope of x509 PKI - this is useful when you want key pairs for CSRs, Servers, and Clients.

* [List Key Pairs](keys/get.md) : `GET /locksmith/keys`

* [Create New Key Pairs](key/post.md) : `POST /locksmith/key`
* [Retrieve Key Pair](key/get.md) : `GET /locksmith/key`

## Key Stores

Key Stores organize groups of Key Pairs.

* [List Key Stores](keystores/get.md) : `GET /locksmith/keystores`

* Retrieve Key Store Information : `GET /locksmith/keystore`
* [Create New Key Store](keystore/post.md) : `POST /locksmith/keystore`