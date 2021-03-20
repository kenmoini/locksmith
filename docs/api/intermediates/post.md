# Create a new Intermediate Certificate Authority

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/intermediates`

**Method** : `POST`

**Content Type** : `JSON`

**Input Data Structure**

```json
{
  "subject": {
    "common_name": string,
    "organization": []string,
    "organizational_unit": []string, // optional
    "country": []string, // optional
    "province": []string, // optional
    "locality": []string, // optional
    "street_address": []string, // optional
    "postal_code": []string, // optional
  },
  "rsa_private_key_passphrase": string, // optional
  "expiration_date": []int, // [ years, months, days ]
  "san_data": {
    "email_addresses": []string,
    "uris": []string
  }
}
```

**Input Data examples**

```json
{
  "parent_cn_path": "Example Labs Root Certificate Authority",
  "certificate_config" {
    "subject": {
      "common_name": "Example Labs Intermediate Certificate Authority",
      "organization": ["Example Labs"],
      "organizational_unit": ["Example Labs Cyber and Information Security"]
    },
    "expiration_date": [
      3, // Years
      0, // Months
      1 // Days
    ],
    "san_data": {
      "email_addresses": ["certmaster@example.labs"],
      "uris": ["https://ca.example.labs:443/"]
    }
  }
}
```

**Request Example**

A cURL request would look like this:

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"parent_cn_path": "Example Labs Root Certificate Authority", "certificate_config":{"subject": {"common_name": "Example Labs Intermediate Certificate Authority", "organization": ["Example Labs"], "organizational_unit": ["Example Labs Cyber and Information Security"]}, "expiration_date": [3,0,1], "san_data": {"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://$PKI_SERVER/locksmith/v1/roots
```

## Success Responses

**Code** : `200 OK`

**Content example** : Response will reflect back the slugged ID of the certificate, the next certificate serial number, and the full representation of the generated CA Certificate.

```json
{
  "status": "intermed-ca-created",
  "errors": [],
  "messages": [
    "Successfully created Intermediate CA Example Labs Intermediate Certificate Authority!"
  ],
  "root": {
    "slug": "example-labs-intermediate-certificate-authority",
    "next_serial": "02",
    "certificate": {
      "Raw": "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",
      "RawTBSCertificate": "MIIFLaADAgECAgECMA0GCSqGSIb3DQEBCwUAMH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB4XDTIxMDMxNTAwMDAwMFoXDTI0MDMxNzAwMDAwMFowgYcxFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTE4MDYGA1UEAxMvRXhhbXBsZSBMYWJzIEludGVybWVkaWF0ZSBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDFpe0D5KMem9Eww6ceJXxoPy3Km1HiSkJdgM4neahKJgG4LzkmKZ7R0nbnbl4IbzI/vTFTCref9XhnUBBoD9TrnBj5lnBgYIMwQ96F9852n/RldTSfIx/UOXNI+criefnzEqfLgUc0KN060R/8EeTR1HkOT1ZQPPcngbF7xUHoX3PrzkrRtPAkgIOmni5AM0YhTFuuFQrVGzU4sZr/z6M1gkR6EIyAe+4JyRxAmsYeQZemi0i58nsC2wm3s9g5+KUDz07yu6afUBrjiSfJ2TzAjZYA6m5q4sS8IaiUExuzo9+Tplw+Ymj1/xxvafV2y6fmifC7EKoZGpYRtUN08Zj+fmdZIL43RqcIcKIsQGPoKAxOWd+joEafRBQm1insc+w4hZwp/JnB54jETW1CO0+oXTXBjNlT81XPeNCEboI6jmkGMg2Gm2vkGDY878zk1qti18RO2dM5CvcSXNURl6OPauM4xaoDYoWYuAyNZKCafQpjSnt3GOAGVKSm4V2i/LmDeSTHuUwmapssZRdP82YlaSO+eIOoIVXUoZYTVhaW7Csn3umf8KFb8zsrAh8cAokBArpfn0i6ji9o9Gdg6rYafIQQVUT6fAIQAYC7ihAjQFF14/DYUv0Bnsfa0GWoBgs5pNKCVQEoLFSeIRGcFeHFeudvQWaaqzBQJKfQEf1a0QIDAQABo4IBwTCCAb0wDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHsVgP8sd5IU8bKmNlYgL5xA3dC1MB8GA1UdIwQYMBaAFGcskynsixDzMFBaxDFi490RO64tMHMGCCsGAQUFBwEBBGcwZTBjBggrBgEFBQcwAoZXaHR0cHM6Ly9jYS5leGFtcGxlLmxhYnM6NDQzL2NlcnRzL2NhLmV4YW1wbGUubGFic19Sb290X0NlcnRpZmljYXRpb25fQXV0aG9yaXR5LmNlcnQucGVtMEAGA1UdEQQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMvMGEGA1UdHwRaMFgwVqBUoFKGUGh0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jcmwvY2EuZXhhbXBsZS5sYWJzX1Jvb3RfQ2VydGlmaWNhdGlvbl9BdXRob3JpdHkuY3JsMEAGA1UdEgQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMv",
      "RawSubjectPublicKeyInfo": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxaXtA+SjHpvRMMOnHiV8aD8typtR4kpCXYDOJ3moSiYBuC85Jime0dJ2525eCG8yP70xUwq3n/V4Z1AQaA/U65wY+ZZwYGCDMEPehffOdp/0ZXU0nyMf1DlzSPnK4nn58xKny4FHNCjdOtEf/BHk0dR5Dk9WUDz3J4Gxe8VB6F9z685K0bTwJICDpp4uQDNGIUxbrhUK1Rs1OLGa/8+jNYJEehCMgHvuCckcQJrGHkGXpotIufJ7AtsJt7PYOfilA89O8rumn1Aa44knydk8wI2WAOpuauLEvCGolBMbs6Pfk6ZcPmJo9f8cb2n1dsun5onwuxCqGRqWEbVDdPGY/n5nWSC+N0anCHCiLEBj6CgMTlnfo6BGn0QUJtYp7HPsOIWcKfyZweeIxE1tQjtPqF01wYzZU/NVz3jQhG6COo5pBjINhptr5Bg2PO/M5NarYtfETtnTOQr3ElzVEZejj2rjOMWqA2KFmLgMjWSgmn0KY0p7dxjgBlSkpuFdovy5g3kkx7lMJmqbLGUXT/NmJWkjvniDqCFV1KGWE1YWluwrJ97pn/ChW/M7KwIfHAKJAQK6X59Iuo4vaPRnYOq2GnyEEFVE+nwCEAGAu4oQI0BRdePw2FL9AZ7H2tBlqAYLOaTSglUBKCxUniERnBXhxXrnb0FmmqswUCSn0BH9WtECAwEAAQ==",
      "RawSubject": "MIGHMRUwEwYDVQQKEwxFeGFtcGxlIExhYnMxNDAyBgNVBAsTK0V4YW1wbGUgTGFicyBDeWJlciBhbmQgSW5mb3JtYXRpb24gU2VjdXJpdHkxODA2BgNVBAMTL0V4YW1wbGUgTGFicyBJbnRlcm1lZGlhdGUgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
      "RawIssuer": "MH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
      "Signature": "IovTkwuqF+EzrsTZevy0Rl/usd4W026jharpgRapeDYdbsK52bYmBp/JdXjRx2w5xxqLNcOBz6pUeQhVTJjmE0PcyiZ3QlQYyC9C7gGe43+4OqQahKHp1QeWPsdFSadZZyY0CvRPHn7RZoMvegIhIpfmHJLhOFH1T7jn+sVMMpzW3l1ts1VS4zZaaNtM1sMxFvcveLuyhCabrfqe/l7wpNdfZCxAupCfqqyQhc8flqjAyEyDNWKzwrw7sI4FRKML6IsiS+KsQ/bvZcAXl+NuJlBMdKwoJRF+4VwaDPDBhQP6phqTFlfWsnTAJc8uiF+0lNUwwN60KJKqTmSVBI+N/XMFsprC3ZrXziHOtNjH8Aa4JdpeG29pJ/UmNsmI3VMzlfSZHNDwT2eu9R4zYkPwwWNLaI/KUuJ9LGqn61UlT7L/M3gcGeGFz+uXKB+/Jgkh5WsaQMx8zDfXyuktutcUV5s8MGElzVNEgoWI8qbuw19izgqc9pY6IplT5zGWkArrqYhI9MaaRI/BacLyYJ72bSJfqg3v1PWx/i/2UKyanL386a2F4y75XmzguU97CjOP/2pKsDFAAduR16RgM+Ib9erhV+aRfHP7U+NKIt3N4M0+9yNJyRpL64x42CsAwGzKD+16Qz9sZokz5sw23WgHLYEuTcYBGbmBHTzJQ3VBukg=",
      "SignatureAlgorithm": 4,
      "PublicKeyAlgorithm": 1,
      "PublicKey": {
        "N": null,
        "E": 65537
      },
      "Version": 3,
      "SerialNumber": 2,
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
        "CommonName": "Example Labs Intermediate Certificate Authority",
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
            "Value": "Example Labs Intermediate Certificate Authority"
          }
        ],
        "ExtraNames": null
      },
      "NotBefore": "2021-03-15T00:00:00Z",
      "NotAfter": "2024-03-17T00:00:00Z",
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
          "Value": "BBR7FYD/LHeSFPGypjZWIC+cQN3QtQ=="
        },
        {
          "Id": [
            2,
            5,
            29,
            35
          ],
          "Critical": false,
          "Value": "MBaAFGcskynsixDzMFBaxDFi490RO64t"
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
      "SubjectKeyId": "exWA/yx3khTxsqY2ViAvnEDd0LU=",
      "AuthorityKeyId": "ZyyTKeyLEPMwUFrEMWLj3RE7ri0=",
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
}
```

## Error Response

**Condition** : If provided data is invalid, missing, or a system error occurs.

**Code** : `200 OK` - 200 OK is returned even on errors due to no specific HTTP error codes corellating to different processes taken place during the CA generation workflow.  Check the `status` field for specific status matches.

**Content example** :

```json
{
  "status": "intermed-ca-creation-error",
  "errors": ["cert-config-error"],
  "messages": ["Missing Expiration Date field"]
}
```

## Return Statuses

Potential return statuses sent back via JSON are as follows:

- `intermed-ca-created` - Successful Intermediate CA generation
- `intermed-ca-exists` - Intermediate CA already exists with that CommonName slug at the specified Certificate Authority Chain Path
- `intermed-ca-creation-error` - Errors are dependant on part of the workflow that failed, such as missing fields or system errors
- `invalid-parent-path` - Invalid parent path, no chain exists or the chain is invalid

