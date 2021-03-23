# List Intermediate Certificate Authorities along Certificate Path

Get the slugged list of Intermediate Certificate Authorities for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/intermediates`

**Method** : `GET`

**Data required** : Certificate Authority Path as a Slash-Delimited String

## Input Parameters

When operating against Intermediate CAs there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

- Example Labs Root CA
  - Example Labs Intermediate CA
    - Example Labs Signing CA

The CommonName chain would be represented as: `Example Labs Root CA/Example Labs Intermediate CA/Example Labs Signing CA`
The slugged CommonName chain (what is stored in the filesystem) would be: `example-labs-root-ca/example-labs-intermediate-ca/example-labs-signing-ca`

To use a CommonName chain, pass the `cn_path` parameter.
To use a slugged CommonName chain, pass the `slug_path` parameter.

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
curl --request GET -G --data-urlencode "cn_path=Example Labs Root Certificate Authority" "http://$PKI_SERVER/locksmith/v1/intermediates"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority" "http://$PKI_SERVER/locksmith/v1/intermediates"
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Listing of Intermediate Certificate Authorities under Example Labs Root Certificate Authority"
  ],
  "intermediate_certificate_authorities": [
    "example-labs-intermediate-certificate-authority",
    "example-labs-apac-certificate-authority",
    "example-labs-emea-certificate-authority"
  ]
}
```

## Notes

* This function is basically just a directory listing, but with a little splitting/joining of strings to map a CA Path to the file system.