# Read Certificate Authority's Certificate Revocation List along Certificate Path

Get the Certificate Revocation List of a Certificate Authority for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/revocations`

**Method** : `GET`

**Data required** : Certificate Authority Path as a Slash-Delimited String

## Input Parameters

When operating against the PKI Chain there is a Certificate Authority Path that is needed - this can be the CommonName chain or the slugged version of the CommonName Chain, eg with the following 3 chained CAs:

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
curl --request GET -G --data-urlencode "cn_path=Example Labs Root Certificate Authority" "http://$PKI_SERVER/locksmith/v1/revocations"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority" "http://$PKI_SERVER/locksmith/v1/revocations"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority/Example Labs Intermediate Certificate Authority" "http://$PKI_SERVER/locksmith/v1/revocations"
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Certificate Revocation List for 'Example Labs Root Certificate Authority'"
  ],
  "slug": "Example Labs Root Certificate Authority",
  "crl_pem": "MIIDPDCCASQCAQEwDQYJKoZIhvcNAQELBQAwfzEVMBMGA1UEChMMRXhhbXBsZSBMYWJzMTQwMgYDVQQLEytFeGFtcGxlIExhYnMgQ3liZXIgYW5kIEluZm9ybWF0aW9uIFNlY3VyaXR5MTAwLgYDVQQDEydFeGFtcGxlIExhYnMgUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkXDTIxMDMyMzA1NDkyNloXDTIyMDMyMzA1NDkyNlqgcTBvMB8GA1UdIwQYMBaAFM2LYr20tBEnlVDSDcLTg7vBaWkbMAoGA1UdFAQDAgEAMEAGA1UdEgQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMvMA0GCSqGSIb3DQEBCwUAA4ICAQBUSiEueNchv0DIO9HyQpJ4ygXBLBvr/1EDZ1RQQJogutXeFWGgV9/i7MSVgA75z/x5T+BaZiuxN1fiF5u687EZTkrMeZbm2NVzPC/2RQw8nPa5SN2ViKtm1J/mOcP2n69qNCU8agMBiNIfQk4j3MmMNR6xKbZuyXK1JuMV7KytvSR+OAVrC0QutaZk1A7UhNJtIFKvHxMz/kB3Iq4DVnC1nxFv5gKshMTXlTNJqxdfdSUaKEN9kW4gmVzNk6Lp4VfA7pAfi29eeEMsgf+AYFv1r+h5evQUCIIgDD/n0cae6qq7uGbxWJ8VJ0E0EVo3IE42tfCfwJrXmleTeUymZe8WqnVTusjvJfmblnpehVxt/gMvR26nozT6bbsoExu11R5lRYkyKOtITU0bGuZGaKnS+cxhWv6Uc6JVJtmZD+K6YEI0dlwWyMoA2MhqIm/BZRxqODo8JOiB6Iik2U1h0WCrWbHGj8fk2NKh09ry3Xb1miYAvXusewi7IbmQkNlPY6+zpnXvROSxbNQVkbUbA6w1f1G2jal5OEWoqVRpMIMoPjS/D2Ae8HmqTEeGU4yiYm76z/6zdQOnfMHKuKu5GQLXYdG60peDzUmBoGE29ab/HNQHJcriQsDfubdW5LYdfPjfOJHksPN46UC1hK0HNnXBLP0h1RTMEQ4vwatUnZmiuQ==",
  "crl_list": {
    "TBSCertList": {
      "Raw": "MIIBJAIBATANBgkqhkiG9w0BAQsFADB/MRUwEwYDVQQKEwxFeGFtcGxlIExhYnMxNDAyBgNVBAsTK0V4YW1wbGUgTGFicyBDeWJlciBhbmQgSW5mb3JtYXRpb24gU2VjdXJpdHkxMDAuBgNVBAMTJ0V4YW1wbGUgTGFicyBSb290IENlcnRpZmljYXRlIEF1dGhvcml0eRcNMjEwMzIzMDU0OTI2WhcNMjIwMzIzMDU0OTI2WqBxMG8wHwYDVR0jBBgwFoAUzYtivbS0ESeVUNINwtODu8FpaRswCgYDVR0UBAMCAQAwQAYDVR0SBDkwN4EXY2VydG1hc3RlckBleGFtcGxlLmxhYnOGHGh0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My8=",
      "Version": 1,
      "Signature": {
        "Algorithm": [
          1,
          2,
          840,
          113549,
          1,
          1,
          11
        ],
        "Parameters": {
          "Class": 0,
          "Tag": 5,
          "IsCompound": false,
          "Bytes": "",
          "FullBytes": "BQA="
        }
      },
      "Issuer": [
        [
          {
            "Type": [
              2,
              5,
              4,
              10
            ],
            "Value": "Example Labs"
          }
        ],
        [
          {
            "Type": [
              2,
              5,
              4,
              11
            ],
            "Value": "Example Labs Cyber and Information Security"
          }
        ],
        [
          {
            "Type": [
              2,
              5,
              4,
              3
            ],
            "Value": "Example Labs Root Certificate Authority"
          }
        ]
      ],
      "ThisUpdate": "2021-03-23T05:49:26Z",
      "NextUpdate": "2022-03-23T05:49:26Z",
      "RevokedCertificates": null,
      "Extensions": [
        {
          "Id": [
            2,
            5,
            29,
            35
          ],
          "Critical": false,
          "Value": "MBaAFM2LYr20tBEnlVDSDcLTg7vBaWkb"
        },
        {
          "Id": [
            2,
            5,
            29,
            20
          ],
          "Critical": false,
          "Value": "AgEA"
        },
        {
          "Id": [
            2,
            5,
            29,
            18
          ],
          "Critical": false,
          "Value": "MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMv"
        }
      ]
    },
    "SignatureAlgorithm": {
      "Algorithm": [
        1,
        2,
        840,
        113549,
        1,
        1,
        11
      ],
      "Parameters": {
        "Class": 0,
        "Tag": 5,
        "IsCompound": false,
        "Bytes": "",
        "FullBytes": "BQA="
      }
    },
    "SignatureValue": {
      "Bytes": "VEohLnjXIb9AyDvR8kKSeMoFwSwb6/9RA2dUUECaILrV3hVhoFff4uzElYAO+c/8eU/gWmYrsTdX4hebuvOxGU5KzHmW5tjVczwv9kUMPJz2uUjdlYirZtSf5jnD9p+vajQlPGoDAYjSH0JOI9zJjDUesSm2bslytSbjFeysrb0kfjgFawtELrWmZNQO1ITSbSBSrx8TM/5AdyKuA1ZwtZ8Rb+YCrITE15UzSasXX3UlGihDfZFuIJlczZOi6eFXwO6QH4tvXnhDLIH/gGBb9a/oeXr0FAiCIAw/59HGnuqqu7hm8VifFSdBNBFaNyBONrXwn8Ca15pXk3lMpmXvFqp1U7rI7yX5m5Z6XoVcbf4DL0dup6M0+m27KBMbtdUeZUWJMijrSE1NGxrmRmip0vnMYVr+lHOiVSbZmQ/iumBCNHZcFsjKANjIaiJvwWUcajg6PCTogeiIpNlNYdFgq1mxxo/H5NjSodPa8t129ZomAL17rHsIuyG5kJDZT2Ovs6Z170TksWzUFZG1GwOsNX9Rto2peThFqKlUaTCDKD40vw9gHvB5qkxHhlOMomJu+s/+s3UDp3zByriruRkC12HRutKXg81JgaBhNvWm/xzUByXK4kLA37m3VuS2HXz43ziR5LDzeOlAtYStBzZ1wSz9IdUUzBEOL8GrVJ2Zork=",
      "BitLength": 4096
    }
  }
}
```