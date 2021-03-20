# Create a new Root Certificate Authority

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/roots`

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
  "subject": {
    "common_name": "Example Labs Root Certificate Authority",
    "organization": ["Example Labs"],
    "organizational_unit": ["Example Labs Cyber and Information Security"]
  },
  "expiration_date": [
    10, // Years
    0, // Months
    1 // Days
  ],
  "san_data": {
    "email_addresses": ["certmaster@example.labs"],
    "uris": ["https://ca.example.labs:443/"]
  }
}
```

**Request Example**

A cURL request would look like this:

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"subject": {"common_name": "Example Labs Root Certificate Authority", "organization": ["Example Labs"], "organizational_unit": ["Example Labs Cyber and Information Security"]}, "expiration_date": [10,0,1], "san_data": {"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}' \
  http://$PKI_SERVER/locksmith/v1/roots
```

## Success Responses

**Code** : `200 OK`

**Content example** : Response will reflect back the slugged ID of the certificate, the next certificate serial number, and the full representation of the generated CA Certificate.

```json
{
  "status": "root-created",
  "errors": [],
  "messages": [
    "Root CA Example Labs Root Certificate Authority created!"
  ],
  "root": {
    "slug": "example-labs-root-certificate-authority",
    "next_serial": "2",
    "certificate": {
      "Raw": "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",
      "RawTBSCertificate": "MIIFJKADAgECAgEBMA0GCSqGSIb3DQEBCwUAMH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB4XDTIxMDMxNTAwMDAwMFoXDTMxMDMxNzAwMDAwMFowfzEVMBMGA1UEChMMRXhhbXBsZSBMYWJzMTQwMgYDVQQLEytFeGFtcGxlIExhYnMgQ3liZXIgYW5kIEluZm9ybWF0aW9uIFNlY3VyaXR5MTAwLgYDVQQDEydFeGFtcGxlIExhYnMgUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCxUP3aTFHCt03tJvG2Dlc//oG03lulwuzBDNxmoCGG8ber8n0Zh/ygqWOCmrOgZbVZp7yGgtjdvYFX7ZOcZgV8FjDRbWBpm1zpRMqIycS0eeJYab4qGvmC01ftsP6EPjPThz4l9fQGLTXbKbZGq3Enxrc4p5zIhZoSpTQZRw7BtUyHZ+jI4U17kNl79Iwy8Ccf7kz3REUPkI7QsJo04JEk7nefqUf0Vg3NA7Ydja1aPFgEvbal7agBtfpnnR8jGeR4NqY5v6nAOiyoVpMJdhW1wEeWz8oTMbVZuCY64zXD8M+sRAbf3pb6S8iuhdfSd+0RslnQ7O5Cj/QGFt6ye2aL6bPXJoVIRByG3GbuWnUzjLxCsXupVzjvl+qm476KxTbTu8YL/MHZ5HuaVetPyzgKQ5+k2EHxBVgFvFCTMLty49dq49Q7uSPaYAkSKAvADNGkM7E6NnJQZk6hwrT9gjaQ/FzCspMu/yiEnrRbdurrdSAWx6ABvjbnKrV0BXeJCFMLnDWho0qjVBi1fUijjVEHMA++KJ/XJC3VeQQ06lqT2iDKvfOwuG+RPbdE7kLMDZugo1WV7F6RcVGB7XBXoNWyAVKFZ2dIzbbd79uhvqh/fa13+LSyoQtWHt1h/R4y3r+72ks6AHRN2k+fX4bw0yscmgK9JPyIpwz3hCxqh2VA0QIDAQABo4IBwTCCAb0wDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFP567Eq+4BKPyV4E5iQp/iBFdAWwMB8GA1UdIwQYMBaAFP567Eq+4BKPyV4E5iQp/iBFdAWwMHMGCCsGAQUFBwEBBGcwZTBjBggrBgEFBQcwAoZXaHR0cHM6Ly9jYS5leGFtcGxlLmxhYnM6NDQzL2NlcnRzL2NhLmV4YW1wbGUubGFic19Sb290X0NlcnRpZmljYXRpb25fQXV0aG9yaXR5LmNlcnQucGVtMEAGA1UdEQQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMvMGEGA1UdHwRaMFgwVqBUoFKGUGh0dHBzOi8vY2EuZXhhbXBsZS5sYWJzOjQ0My9jcmwvY2EuZXhhbXBsZS5sYWJzX1Jvb3RfQ2VydGlmaWNhdGlvbl9BdXRob3JpdHkuY3JsMEAGA1UdEgQ5MDeBF2NlcnRtYXN0ZXJAZXhhbXBsZS5sYWJzhhxodHRwczovL2NhLmV4YW1wbGUubGFiczo0NDMv",
      "RawSubjectPublicKeyInfo": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsVD92kxRwrdN7Sbxtg5XP/6BtN5bpcLswQzcZqAhhvG3q/J9GYf8oKljgpqzoGW1Wae8hoLY3b2BV+2TnGYFfBYw0W1gaZtc6UTKiMnEtHniWGm+Khr5gtNX7bD+hD4z04c+JfX0Bi012ym2RqtxJ8a3OKecyIWaEqU0GUcOwbVMh2foyOFNe5DZe/SMMvAnH+5M90RFD5CO0LCaNOCRJO53n6lH9FYNzQO2HY2tWjxYBL22pe2oAbX6Z50fIxnkeDamOb+pwDosqFaTCXYVtcBHls/KEzG1WbgmOuM1w/DPrEQG396W+kvIroXX0nftEbJZ0OzuQo/0Bhbesntmi+mz1yaFSEQchtxm7lp1M4y8QrF7qVc475fqpuO+isU207vGC/zB2eR7mlXrT8s4CkOfpNhB8QVYBbxQkzC7cuPXauPUO7kj2mAJEigLwAzRpDOxOjZyUGZOocK0/YI2kPxcwrKTLv8ohJ60W3bq63UgFsegAb425yq1dAV3iQhTC5w1oaNKo1QYtX1Io41RBzAPviif1yQt1XkENOpak9ogyr3zsLhvkT23RO5CzA2boKNVlexekXFRge1wV6DVsgFShWdnSM223e/bob6of32td/i0sqELVh7dYf0eMt6/u9pLOgB0TdpPn1+G8NMrHJoCvST8iKcM94QsaodlQNECAwEAAQ==",
      "RawSubject": "MH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
      "RawIssuer": "MH8xFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEwMC4GA1UEAxMnRXhhbXBsZSBMYWJzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
      "Signature": "MTPJt5LBzIsmCuAVfO3flnCIZetdrYVpFcY0czXMn2mEf6NtdzcC9ISl9kJ+i78IL+RvVSnqYdQSVeNHVmWTMjIpIPfFTYj41DckHHZLwFVzOpvD5W07+vw3rMfUKMJ/MijmTwJQzm5TNjThfCvsxPkoWyqlB9R+u3uuwO9ks6+UlAsTjVuc9xgAfgHnot7JZGs8jniCzo1BVTwRO1072X9Om8QVehDdMFXGP7D0pZuOJXJa68WLP4yDJCi+jXfuqOK4BdvG3ZlPgrsdNlYishJhV4kOAJwgGiwW46jPIwDexXs7819nMcUcjzR0E4tvazFde++lySV/45QVB3BEisKHB1njKK3lSJKv5xlRDxL2aBMcd3DzqO7eRXDa2QT7jGRk6FavIOCALtbRWMj0LkKLMf71INFR0Ox+DhQQ96qnjK55cKph2WMIyiHL8s1QXqFbKMWWbdnPnUUj4obCstlJDXKjZcm4MlDpbm9jE4KKjIEJq9txLJdggrXwr/nFxatmg0lPINzmdbWSJm3IoqAzjZnnPtH0OmmCokM07tmTwKwZ/AcuMqTU4ZQAunDGbSI0dep7FzxgJy9O1q0+9Bwt1t09qIHPFd3S5x2M9Hw1J4f5DbxWELnoZGP5AuQfLTf6TiLOyaWh/0JuiPWjb00rWE7JzTJ/q+Gr/VGM27M=",
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
      "NotBefore": "2021-03-15T00:00:00Z",
      "NotAfter": "2031-03-17T00:00:00Z",
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
          "Value": "BBT+euxKvuASj8leBOYkKf4gRXQFsA=="
        },
        {
          "Id": [
            2,
            5,
            29,
            35
          ],
          "Critical": false,
          "Value": "MBaAFP567Eq+4BKPyV4E5iQp/iBFdAWw"
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
      "SubjectKeyId": "/nrsSr7gEo/JXgTmJCn+IEV0BbA=",
      "AuthorityKeyId": "/nrsSr7gEo/JXgTmJCn+IEV0BbA=",
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
  "status": "root-creation-error",
  "errors": ["cert-config-error"],
  "messages": ["Missing Expiration Date field"]
}
```

## Return Statuses

Potential return statuses sent back via JSON are as follows:

- `root-created` - Successful root CA generation
- `root-exists` - Root CA already exists with that CommonName slug
- `root-creation-error` - Errors are dependant on part of the workflow that failed, such as missing fields or system errors

