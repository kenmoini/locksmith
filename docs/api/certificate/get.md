# Read Certificate along Certificate Path [WIP]

Get the information of a Certificate for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/certificate`

**Method** : `GET`

**Data required** : Certificate Authority Path as a Slash-Delimited String and Certificate ID.

## Input Parameters

### Certificate Authority Path

When operating against PKI Chain there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

- Example Labs Root CA
  - Example Labs Intermediate CA
    - Example Labs Signing CA

The CommonName chain would be represented as: `Example Labs Root CA/Example Labs Intermediate CA/Example Labs Signing CA`
The slugged CommonName chain (what is stored in the filesystem) would be: `example-labs-root-ca/example-labs-intermediate-ca/example-labs-signing-ca`

To use a CommonName chain, pass the `parent_cn_path` parameter.
To use a slugged CommonName chain, pass the `parent_slug_path` parameter.

### Certificate ID

The Certificate ID is the CommonName or the slugged CommonName of the Certificate - pass the `certificate_id` parameter.

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" --data-urlencode "certificate_id=Example Labs Intermediate Certificate Authority" "http://$PKI_SERVER/locksmith/v1/certificate"
curl --request GET -G --data-urlencode "parent_slug_path=example-labs-root-certificate-authority/Example Labs Intermediate Certificate Authority" --data-urlencode "certificate_id=OpenVPN Server" "http://$PKI_SERVER/locksmith/v1/certificate"
```

And the data returned would be the minified version of the following JSON:

```json

```