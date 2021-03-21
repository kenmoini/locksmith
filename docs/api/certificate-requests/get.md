# List Certificate Requests along Certificate Path

Get the slugged list of Certificate Requests for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/certificate-requests`

**Method** : `GET`

**Data required** : Certificate Authority Path as a Slash-Delimited String

## Input Parameters

When operating against Certificate Requests there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

- Example Labs Root CA
  - Example Labs Intermediate CA
    - Example Labs Signing CA

The full CommonName chain would be a string represented as: `Example Labs Root CA/Example Labs Intermediate CA/Example Labs Signing CA`
The full slugged CommonName chain (what is stored in the filesystem) would be this string: `example-labs-root-ca/example-labs-intermediate-ca/example-labs-signing-ca`

To use a CommonName chain, pass the `parent_cn_path` string parameter.
To use a slugged CommonName chain, pass the `parent_slug_path` string parameter.

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" "http://$PKI_SERVER/locksmith/v1/certificate-requests"
curl --request GET -G --data-urlencode "parent_slug_path=example-labs-root-certificate-authority" "http://$PKI_SERVER/locksmith/v1/certificate-requests"
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority/example-labs-intermediate-certificate-authority" "http://$PKI_SERVER/locksmith/v1/certificate-requests"
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Listing of Certificate Requests under Example Labs Root Certificate Authority"
  ],
  "certificate_requests": [
    "example-labs-intermediate-certificate-authority",
    "example-labs-apac-certificate-authority",
    "example-labs-emea-certificate-authority"
  ]
}
```

## Notes

* This function is basically just a directory listing, but with a little splitting/joining of strings to map a CA Path to the file system.