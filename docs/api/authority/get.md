# Read Certificate Authority along Certificate Path

Get the information of a Certificate Authority for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/authority`

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
curl --request GET -G --data-urlencode "cn_path=Example Labs Root Certificate Authority" "http://$PKI_SERVER/locksmith/v1/authority"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority" "http://$PKI_SERVER/locksmith/v1/authority"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority/Example Labs Intermediate Certificate Authority" "http://$PKI_SERVER/locksmith/v1/authority"
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Certificate Authority information for 'Example Labs Root Certificate Authority'"
  ],
  "slug": "example-labs-root-certificate-authority",
  "certificate_pem": "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",
  "certificate_information": {
    "Raw": "MIIHPDCCBSSgAwIBAgIBATANBgkqhkiG9w0BAQsFADB/MRUwEwYDVQQKEwxFeGFtcGxlIExhYnMxNDAyBgNVBAsTK0V4YW1wbGUgTGFicyBDeWJlciBhbmQgSW5mb3JtYXRpb24gU2VjdXJpdHkxMDAuBgNVBAMTJ0V4YW1wbGUgTGFicyBSb290IENlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTAzMTcwMDAwMDBaFw0zMTAzMTkwMDAwMDBaMH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtCASxdurIu/pv8QdMHdQ1IXFvD8NflP30k3Kfk8wuOWF9NGigXG5Qnt6Y7ZbfU527/WK3/9A5oS+sceoApkuDN7Sb1IiyoVoaOEIFzRzHXpVsjppPVPDKg4GoQnxv0YWjf5WujN56fSNPcck9F4tf6CiY7zRX6w0AFhbfMOIjkZQWA4N00IXzvtWX7BSIyNyKPGHskJQm2ehJCTVPSWxjdyk38Mq9f8FsDCXYoFt/ST+0hH8HeWjV6CEnUTq2pyLQP/4N+tt1iBlef7Porp/3/pQFcv2RAeFXas3bsJZiiulv1v28qfyD8XiBT4PvTey8A0IziaOfu5hpf18MjTwVa4z4sbpg8z1rK29XSOTeJbM2OPrPTyq6kMYWn0bfW+Rh0iDxucFjM9b7WtJe/2snETaFAePIM1UCelewmNuWJRVAQTE90h1FO4cmldV0rBKQLbmJ9XsTQ+d/zqKUPeLrJQQSvaR/pFiFxbrybA3Z3MX8EpX84M7MKwlfmG3cbnHUSEb3nKsl56XIB4rrnQs8mzjWJ/9+rKVIEbJnh6BQYydB4MmlPHQHFY42quCcQPbpiH+VAaqJciDzpZ/I/Xze3T2zaQXLuNTAdbmoB7iltTpTNWwsthUrbgflWYc7EqtDq0RnpwtZHTmLLHnFsTHYyflvNMlcKGO42LyY61JNjECAwEAAaOCAcEwggG9MA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRCpZJn/s2xSA0Y9fvGUKKwrGKyojAfBgNVHSMEGDAWgBRCpZJn/s2xSA0Y9fvGUKKwrGKyojBzBggrBgEFBQcBAQRnMGUwYwYIKwYBBQUHMAKGV2h0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jZXJ0cy9jYS5leGFtcGxlLmxhYnNfUm9vdF9DZXJ0aWZpY2F0aW9uX0F1dGhvcml0eS5jZXJ0LnBlbTBABgNVHREEOTA3gRdjZXJ0bWFzdGVyQGV4YW1wbGUubGFic4YcaHR0cHM6Ly9jYS5leGFtcGxlLmxhYnM6NDQzLzBhBgNVHR8EWjBYMFagVKBShlBodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMvY3JsL2NhLmV4YW1wbGUubGFic19Sb290X0NlcnRpZmljYXRpb25fQXV0aG9yaXR5LmNybDBABgNVHRIEOTA3gRdjZXJ0bWFzdGVyQGV4YW1wbGUubGFic4YcaHR0cHM6Ly9jYS5leGFtcGxlLmxhYnM6NDQzLzANBgkqhkiG9w0BAQsFAAOCAgEAR5HQgqIVkDXBRlwzvwsrsAc9b1AC3p+Of9U6JPxmfG7GzN7iY0j+yA5D++/h0Yjcvt51v2gdLGposEgFdIGFlvNKy6MQGahkHcdufwa9VWqn5BUBTEJfWyC9RQ0LRP9xM3xI95Ckg7DHqiHu+l+z2IBTVmYOzlV2z6wqu1P0hAyCVLBhylQbCiy+tIdVsXnW8He0BxfXzWypf7BNbd4Ybnls15eWxckK5w1ojWHVuY9YJokkIRuR3H9J7fiuFJK036V7+OkZsEUSwUslfI4ncB7vT13Gmy9B6soYJWBNM/3YZ5cRWZ9ZWYhwemGxlSYynsFa8i27hUAS1S8gZt8HjPxtaaKbimgIhfabkoRMQUge/vC3FAv29jOSIsLMXP2JvZyX6jwM7UmTtwI8ZwOPcAr2sHE92PFJmZ+0n0EGutfv1rK8N9Q+/rI/WOBGjkflOxhZRDVunXuseomxVveSYDVantw3k2jfEsl84taBPN87QGIFddns0yFdpcUvLH2NYh96gZu9xI4yqGyJlRVTLmkbfNWlRVou1XW1aAzRuf/xdEJzqXvKTPhjZFt/qNqvAM6b0nQmes7ViTfbqF9iZNtIFCzby1yDZBp2hpYSLv4LTnuqdDeiDVAmZ/EUwjdsK3GEpUt/7Ope5aza9O85taZObcVPBdHcsdwW08r7TB8=",
    "RawTBSCertificate": "MIIFJKADAgECAgEBMA0GCSqGSIb3DQEBCwUAMH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB4XDTIxMDMxNzAwMDAwMFoXDTMxMDMxOTAwMDAwMFowfzEVMBMGA1UEChMMRXhhbXBsZSBMYWJzMTQwMgYDVQQLEytFeGFtcGxlIExhYnMgQ3liZXIgYW5kIEluZm9ybWF0aW9uIFNlY3VyaXR5MTAwLgYDVQQDEydFeGFtcGxlIExhYnMgUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC0IBLF26si7+m/xB0wd1DUhcW8Pw1+U/fSTcp+TzC45YX00aKBcblCe3pjtlt9Tnbv9Yrf/0DmhL6xx6gCmS4M3tJvUiLKhWho4QgXNHMdelWyOmk9U8MqDgahCfG/RhaN/la6M3np9I09xyT0Xi1/oKJjvNFfrDQAWFt8w4iORlBYDg3TQhfO+1ZfsFIjI3Io8YeyQlCbZ6EkJNU9JbGN3KTfwyr1/wWwMJdigW39JP7SEfwd5aNXoISdROranItA//g3623WIGV5/s+iun/f+lAVy/ZEB4VdqzduwlmKK6W/W/byp/IPxeIFPg+9N7LwDQjOJo5+7mGl/XwyNPBVrjPixumDzPWsrb1dI5N4lszY4+s9PKrqQxhafRt9b5GHSIPG5wWMz1vta0l7/aycRNoUB48gzVQJ6V7CY25YlFUBBMT3SHUU7hyaV1XSsEpAtuYn1exND53/OopQ94uslBBK9pH+kWIXFuvJsDdncxfwSlfzgzswrCV+YbdxucdRIRvecqyXnpcgHiuudCzybONYn/36spUgRsmeHoFBjJ0HgyaU8dAcVjjaq4JxA9umIf5UBqolyIPOln8j9fN7dPbNpBcu41MB1uagHuKW1OlM1bCy2FStuB+VZhzsSq0OrRGenC1kdOYssecWxMdjJ+W80yVwoY7jYvJjrUk2MQIDAQABo4IBwTCCAb0wDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFEKlkmf+zbFIDRj1+8ZQorCsYrKiMB8GA1UdIwQYMBaAFEKlkmf+zbFIDRj1+8ZQorCsYrKiMHMGCCsGAQUFBwEBBGcwZTBjBggrBgEFBQcwAoZXaHR0cHM6Ly9jYS5leGFtcGxlLmxhYnM6NDQzL2NlcnRzL2NhLmV4YW1wbGUubGFic19Sb290X0NlcnRpZmljYXRpb25fQXV0aG9yaXR5LmNlcnQucGVtMEAGA1UdEQQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMvMGEGA1UdHwRaMFgwVqBUoFKGUGh0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jcmwvY2EuZXhhbXBsZS5sYWJzX1Jvb3RfQ2VydGlmaWNhdGlvbl9BdXRob3JpdHkuY3JsMEAGA1UdEgQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMv",
    "RawSubjectPublicKeyInfo": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtCASxdurIu/pv8QdMHdQ1IXFvD8NflP30k3Kfk8wuOWF9NGigXG5Qnt6Y7ZbfU527/WK3/9A5oS+sceoApkuDN7Sb1IiyoVoaOEIFzRzHXpVsjppPVPDKg4GoQnxv0YWjf5WujN56fSNPcck9F4tf6CiY7zRX6w0AFhbfMOIjkZQWA4N00IXzvtWX7BSIyNyKPGHskJQm2ehJCTVPSWxjdyk38Mq9f8FsDCXYoFt/ST+0hH8HeWjV6CEnUTq2pyLQP/4N+tt1iBlef7Porp/3/pQFcv2RAeFXas3bsJZiiulv1v28qfyD8XiBT4PvTey8A0IziaOfu5hpf18MjTwVa4z4sbpg8z1rK29XSOTeJbM2OPrPTyq6kMYWn0bfW+Rh0iDxucFjM9b7WtJe/2snETaFAePIM1UCelewmNuWJRVAQTE90h1FO4cmldV0rBKQLbmJ9XsTQ+d/zqKUPeLrJQQSvaR/pFiFxbrybA3Z3MX8EpX84M7MKwlfmG3cbnHUSEb3nKsl56XIB4rrnQs8mzjWJ/9+rKVIEbJnh6BQYydB4MmlPHQHFY42quCcQPbpiH+VAaqJciDzpZ/I/Xze3T2zaQXLuNTAdbmoB7iltTpTNWwsthUrbgflWYc7EqtDq0RnpwtZHTmLLHnFsTHYyflvNMlcKGO42LyY61JNjECAwEAAQ==",
    "RawSubject": "MH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
    "RawIssuer": "MH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
    "Signature": "R5HQgqIVkDXBRlwzvwsrsAc9b1AC3p+Of9U6JPxmfG7GzN7iY0j+yA5D++/h0Yjcvt51v2gdLGposEgFdIGFlvNKy6MQGahkHcdufwa9VWqn5BUBTEJfWyC9RQ0LRP9xM3xI95Ckg7DHqiHu+l+z2IBTVmYOzlV2z6wqu1P0hAyCVLBhylQbCiy+tIdVsXnW8He0BxfXzWypf7BNbd4Ybnls15eWxckK5w1ojWHVuY9YJokkIRuR3H9J7fiuFJK036V7+OkZsEUSwUslfI4ncB7vT13Gmy9B6soYJWBNM/3YZ5cRWZ9ZWYhwemGxlSYynsFa8i27hUAS1S8gZt8HjPxtaaKbimgIhfabkoRMQUge/vC3FAv29jOSIsLMXP2JvZyX6jwM7UmTtwI8ZwOPcAr2sHE92PFJmZ+0n0EGutfv1rK8N9Q+/rI/WOBGjkflOxhZRDVunXuseomxVveSYDVantw3k2jfEsl84taBPN87QGIFddns0yFdpcUvLH2NYh96gZu9xI4yqGyJlRVTLmkbfNWlRVou1XW1aAzRuf/xdEJzqXvKTPhjZFt/qNqvAM6b0nQmes7ViTfbqF9iZNtIFCzby1yDZBp2hpYSLv4LTnuqdDeiDVAmZ/EUwjdsK3GEpUt/7Ope5aza9O85taZObcVPBdHcsdwW08r7TB8=",
    "SignatureAlgorithm": 4,
    "PublicKeyAlgorithm": 1,
    "PublicKey": {
      "N": null,
      "E": 65537
    },
    "Version": 3,
    "SerialNumber": 1,
    "Issuer": {
      "Country": null,
      "Organization": [
        "Example Labs"
      ],
      "OrganizationalUnit": [
        "Example Labs Cyber and Information Security"
      ],
      "Locality": null,
      "Province": null,
      "StreetAddress": null,
      "PostalCode": null,
      "SerialNumber": "",
      "CommonName": "Example Labs Root Certificate Authority",
      "Names": [
        {
          "Type": [
            2,
            5,
            4,
            10
          ],
          "Value": "Example Labs"
        },
        {
          "Type": [
            2,
            5,
            4,
            11
          ],
          "Value": "Example Labs Cyber and Information Security"
        },
        {
          "Type": [
            2,
            5,
            4,
            3
          ],
          "Value": "Example Labs Root Certificate Authority"
        }
      ],
      "ExtraNames": null
    },
    "Subject": {
      "Country": null,
      "Organization": [
        "Example Labs"
      ],
      "OrganizationalUnit": [
        "Example Labs Cyber and Information Security"
      ],
      "Locality": null,
      "Province": null,
      "StreetAddress": null,
      "PostalCode": null,
      "SerialNumber": "",
      "CommonName": "Example Labs Root Certificate Authority",
      "Names": [
        {
          "Type": [
            2,
            5,
            4,
            10
          ],
          "Value": "Example Labs"
        },
        {
          "Type": [
            2,
            5,
            4,
            11
          ],
          "Value": "Example Labs Cyber and Information Security"
        },
        {
          "Type": [
            2,
            5,
            4,
            3
          ],
          "Value": "Example Labs Root Certificate Authority"
        }
      ],
      "ExtraNames": null
    },
    "NotBefore": "2021-03-17T00:00:00Z",
    "NotAfter": "2031-03-19T00:00:00Z",
    "KeyUsage": 96,
    "Extensions": [
      {
        "Id": [
          2,
          5,
          29,
          15
        ],
        "Critical": true,
        "Value": "AwIBBg=="
      },
      {
        "Id": [
          2,
          5,
          29,
          19
        ],
        "Critical": true,
        "Value": "MAMBAf8="
      },
      {
        "Id": [
          2,
          5,
          29,
          14
        ],
        "Critical": false,
        "Value": "BBRCpZJn/s2xSA0Y9fvGUKKwrGKyog=="
      },
      {
        "Id": [
          2,
          5,
          29,
          35
        ],
        "Critical": false,
        "Value": "MBaAFEKlkmf+zbFIDRj1+8ZQorCsYrKi"
      },
      {
        "Id": [
          1,
          3,
          6,
          1,
          5,
          5,
          7,
          1,
          1
        ],
        "Critical": false,
        "Value": "MGUwYwYIKwYBBQUHMAKGV2h0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jZXJ0cy9jYS5leGFtcGxlLmxhYnNfUm9vdF9DZXJ0aWZpY2F0aW9uX0F1dGhvcml0eS5jZXJ0LnBlbQ=="
      },
      {
        "Id": [
          2,
          5,
          29,
          17
        ],
        "Critical": false,
        "Value": "MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMv"
      },
      {
        "Id": [
          2,
          5,
          29,
          31
        ],
        "Critical": false,
        "Value": "MFgwVqBUoFKGUGh0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jcmwvY2EuZXhhbXBsZS5sYWJzX1Jvb3RfQ2VydGlmaWNhdGlvbl9BdXRob3JpdHkuY3Js"
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
    ],
    "ExtraExtensions": null,
    "UnhandledCriticalExtensions": null,
    "ExtKeyUsage": null,
    "UnknownExtKeyUsage": null,
    "BasicConstraintsValid": true,
    "IsCA": true,
    "MaxPathLen": -1,
    "MaxPathLenZero": false,
    "SubjectKeyId": "QqWSZ/7NsUgNGPX7xlCisKxisqI=",
    "AuthorityKeyId": "QqWSZ/7NsUgNGPX7xlCisKxisqI=",
    "OCSPServer": null,
    "IssuingCertificateURL": [
      "https://ca.example.labs:443/certs/ca.example.labs_Root_Certification_Authority.cert.pem"
    ],
    "DNSNames": null,
    "EmailAddresses": [
      "certmaster@example.labs"
    ],
    "IPAddresses": null,
    "URIs": [
      {
        "Scheme": "https",
        "Opaque": "",
        "User": null,
        "Host": "ca.example.labs:443",
        "Path": "/",
        "RawPath": "",
        "ForceQuery": false,
        "RawQuery": "",
        "Fragment": "",
        "RawFragment": ""
      }
    ],
    "PermittedDNSDomainsCritical": false,
    "PermittedDNSDomains": null,
    "ExcludedDNSDomains": null,
    "PermittedIPRanges": null,
    "ExcludedIPRanges": null,
    "PermittedEmailAddresses": null,
    "ExcludedEmailAddresses": null,
    "PermittedURIDomains": null,
    "ExcludedURIDomains": null,
    "CRLDistributionPoints": [
      "https://ca.example.labs:443/crl/ca.example.labs_Root_Certification_Authority.crl"
    ],
    "PolicyIdentifiers": null
  }
}
```