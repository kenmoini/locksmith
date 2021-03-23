# Read Certificate Request along Certificate Path

Get the information of a Certificate Request for a given Certificate Path.

The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/certificate-request`

**Method** : `GET`

**Data required** : Certificate Authority Path as a Slash-Delimited String and Certificate Request ID.

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

### Certificate Request ID

The Certificate Request ID is the CommonName or the slugged CommonName of the Certificate Request - pass the `certificate_request_id` parameter.

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
curl --request GET -G --data-urlencode "cn_path=Example Labs Root Certificate Authority" --data-urlencode "certificate_request_id=Example Labs Intermediate Certificate Authority" "http://$PKI_SERVER/locksmith/v1/certificate-request"
curl --request GET -G --data-urlencode "slug_path=example-labs-root-certificate-authority/Example Labs Intermediate Certificate Authority" --data-urlencode "certificate_request_id=OpenVPN Server" "http://$PKI_SERVER/locksmith/v1/certificate-request"
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Certificate Request information for 'Example Labs Root Certificate Authority'"
  ],
  "csr_pem": "MIIE8jCCAtoCAQAwgYcxFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTE4MDYGA1UEAxMvRXhhbXBsZSBMYWJzIEludGVybWVkaWF0ZSBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDsi/quMy3PrX7w+Rh9Tt+rk7sycF44vzITGqT+A33rPtLMtOi1JOrH+hs58VJWCH+Ga5CnFWeC3joqhMC7QrtnknN3vvbxfYHN2KctJSycAzpdvO8GfRicUgKf3cqboJgJtAxGTc/9WhNlxpZ2q5QqeylLT/SKC2h/XyDFMBhYBQZrACK4XYwAPLLBNjyxD0esa0kwc+fdj6RUva8cOTEgLMe0Y4mnY+Dnd0lp7DD/yed+faRh+WsvDsj5tD1boxii7bjS1L1IYwdD+E7rywztbyki2jnC6tHUp4sTvrQV9KTptbi3+eVs5dETQRPw6xGskN3Wj+x05UkwzhmYLqMLdnVMv1cBqU4nc5IOnDCFBOM/6lxVyQaWfxqEPcm2f4DEgqKw4NkhF+8w1slFq1Z1Bb7mB0ibm4qYp3ejdDjV+KtOctDlmj/T78I2NK3iNoCW5AXJBnZEXaMQcCsUxMZqrp8CdwcCeBr+7EhwPwt3cDYUN2TsEPM5AJ7358UJH/4Na/mNmv39T9SwpqFUNlSyjBGxG/qh2AZNUrYTQE+p0P9qHajLotSEelCHEi+LvEKXPxSuzizxDZM8oPJSNFa5sINxBvO7b0XeWO+R0jKU799rkxJVtWmClYjlAVyO+Df2GWTK4/KP3Gplp1RwvrlIhNWz6Ryh+zFqQSgm08zVGwIDAQABoCUwIwYJKoZIhvcNAQkOMRYwFDASBgNVHRMBAf8ECDAGAQH/AgEAMA0GCSqGSIb3DQEBDQUAA4ICAQCTVw9F1gTTxhBH+IL2UYI0acb2YKOTSo9/CyJ1zgjc0a7Yzh3E6yfUj42V11PBBbkRgWdgJ3rffrkXhYNjUdgU4L631RPq8nhyBz7N0T66cJWQipnKEcsXGJsaGB4oPte8uL9yKwHfj+6RiDFrpGK0CQkV0L76tI0dMEGJBbi16QCgp4WnQUDp/RMrNVbr8C2k7HWDEVn2dYJz7ZQnnPDpmeIWXJn4U8vMyzLqhtJQ1+DLN220Kxg9ZlsX7iw5/s6mNC3d6roQEjPdDCES/S5pkQdEKxPDoAYTW2jfBwfFrGlXbxVkpR9/hqlAVIEfRAu4Bff2q+2D2DtDu2TyG/9rC0gVJ63yxq3m+GVX4AzabuJaxxKOWrfsuPwxoTFfxcrZm54q322xhl/Stl/aPOgJFrNPljvRJlhD8+lAe7R0NlDqb2OEyfg1DOUveFAI5DwTDFl6yIa2pufm/N8shsBfYhz9Effxn97FiwnZ6b4RZ3+hMIqTHHhsCZE7hE0Ke+p59cYMjG4weJn6/SZVKq2RcBKz3YCp+K0C7BWWlHpJpjjhYYfAWKuKJny8MQlMxRtNXBxELV7A4sZNBBDuucXIZAOin7O6IhzOu/ihwhvh90wJsVEPYYqfj/q2/+EAjkdp9NiChSQ8YNn23/2bvw1l6isqc6qi5ucXUPYxf5Hyxw==",
  "certificate_request": {
    "Raw": "MIIE8jCCAtoCAQAwgYcxFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTE4MDYGA1UEAxMvRXhhbXBsZSBMYWJzIEludGVybWVkaWF0ZSBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDsi/quMy3PrX7w+Rh9Tt+rk7sycF44vzITGqT+A33rPtLMtOi1JOrH+hs58VJWCH+Ga5CnFWeC3joqhMC7QrtnknN3vvbxfYHN2KctJSycAzpdvO8GfRicUgKf3cqboJgJtAxGTc/9WhNlxpZ2q5QqeylLT/SKC2h/XyDFMBhYBQZrACK4XYwAPLLBNjyxD0esa0kwc+fdj6RUva8cOTEgLMe0Y4mnY+Dnd0lp7DD/yed+faRh+WsvDsj5tD1boxii7bjS1L1IYwdD+E7rywztbyki2jnC6tHUp4sTvrQV9KTptbi3+eVs5dETQRPw6xGskN3Wj+x05UkwzhmYLqMLdnVMv1cBqU4nc5IOnDCFBOM/6lxVyQaWfxqEPcm2f4DEgqKw4NkhF+8w1slFq1Z1Bb7mB0ibm4qYp3ejdDjV+KtOctDlmj/T78I2NK3iNoCW5AXJBnZEXaMQcCsUxMZqrp8CdwcCeBr+7EhwPwt3cDYUN2TsEPM5AJ7358UJH/4Na/mNmv39T9SwpqFUNlSyjBGxG/qh2AZNUrYTQE+p0P9qHajLotSEelCHEi+LvEKXPxSuzizxDZM8oPJSNFa5sINxBvO7b0XeWO+R0jKU799rkxJVtWmClYjlAVyO+Df2GWTK4/KP3Gplp1RwvrlIhNWz6Ryh+zFqQSgm08zVGwIDAQABoCUwIwYJKoZIhvcNAQkOMRYwFDASBgNVHRMBAf8ECDAGAQH/AgEAMA0GCSqGSIb3DQEBDQUAA4ICAQCTVw9F1gTTxhBH+IL2UYI0acb2YKOTSo9/CyJ1zgjc0a7Yzh3E6yfUj42V11PBBbkRgWdgJ3rffrkXhYNjUdgU4L631RPq8nhyBz7N0T66cJWQipnKEcsXGJsaGB4oPte8uL9yKwHfj+6RiDFrpGK0CQkV0L76tI0dMEGJBbi16QCgp4WnQUDp/RMrNVbr8C2k7HWDEVn2dYJz7ZQnnPDpmeIWXJn4U8vMyzLqhtJQ1+DLN220Kxg9ZlsX7iw5/s6mNC3d6roQEjPdDCES/S5pkQdEKxPDoAYTW2jfBwfFrGlXbxVkpR9/hqlAVIEfRAu4Bff2q+2D2DtDu2TyG/9rC0gVJ63yxq3m+GVX4AzabuJaxxKOWrfsuPwxoTFfxcrZm54q322xhl/Stl/aPOgJFrNPljvRJlhD8+lAe7R0NlDqb2OEyfg1DOUveFAI5DwTDFl6yIa2pufm/N8shsBfYhz9Effxn97FiwnZ6b4RZ3+hMIqTHHhsCZE7hE0Ke+p59cYMjG4weJn6/SZVKq2RcBKz3YCp+K0C7BWWlHpJpjjhYYfAWKuKJny8MQlMxRtNXBxELV7A4sZNBBDuucXIZAOin7O6IhzOu/ihwhvh90wJsVEPYYqfj/q2/+EAjkdp9NiChSQ8YNn23/2bvw1l6isqc6qi5ucXUPYxf5Hyxw==",
    "RawTBSCertificateRequest": "MIIC2gIBADCBhzEVMBMGA1UEChMMRXhhbXBsZSBMYWJzMTQwMgYDVQQLEytFeGFtcGxlIExhYnMgQ3liZXIgYW5kIEluZm9ybWF0aW9uIFNlY3VyaXR5MTgwNgYDVQQDEy9FeGFtcGxlIExhYnMgSW50ZXJtZWRpYXRlIENlcnRpZmljYXRlIEF1dGhvcml0eTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAOyL+q4zLc+tfvD5GH1O36uTuzJwXji/MhMapP4Dfes+0sy06LUk6sf6GznxUlYIf4ZrkKcVZ4LeOiqEwLtCu2eSc3e+9vF9gc3Ypy0lLJwDOl287wZ9GJxSAp/dypugmAm0DEZNz/1aE2XGlnarlCp7KUtP9IoLaH9fIMUwGFgFBmsAIrhdjAA8ssE2PLEPR6xrSTBz592PpFS9rxw5MSAsx7Rjiadj4Od3SWnsMP/J5359pGH5ay8OyPm0PVujGKLtuNLUvUhjB0P4TuvLDO1vKSLaOcLq0dSnixO+tBX0pOm1uLf55Wzl0RNBE/DrEayQ3daP7HTlSTDOGZguowt2dUy/VwGpTidzkg6cMIUE4z/qXFXJBpZ/GoQ9ybZ/gMSCorDg2SEX7zDWyUWrVnUFvuYHSJubipind6N0ONX4q05y0OWaP9PvwjY0reI2gJbkBckGdkRdoxBwKxTExmqunwJ3BwJ4Gv7sSHA/C3dwNhQ3ZOwQ8zkAnvfnxQkf/g1r+Y2a/f1P1LCmoVQ2VLKMEbEb+qHYBk1SthNAT6nQ/2odqMui1IR6UIcSL4u8Qpc/FK7OLPENkzyg8lI0Vrmwg3EG87tvRd5Y75HSMpTv32uTElW1aYKViOUBXI74N/YZZMrj8o/camWnVHC+uUiE1bPpHKH7MWpBKCbTzNUbAgMBAAGgJTAjBgkqhkiG9w0BCQ4xFjAUMBIGA1UdEwEB/wQIMAYBAf8CAQA=",
    "RawSubjectPublicKeyInfo": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA7Iv6rjMtz61+8PkYfU7fq5O7MnBeOL8yExqk/gN96z7SzLTotSTqx/obOfFSVgh/hmuQpxVngt46KoTAu0K7Z5Jzd7728X2BzdinLSUsnAM6XbzvBn0YnFICn93Km6CYCbQMRk3P/VoTZcaWdquUKnspS0/0igtof18gxTAYWAUGawAiuF2MADyywTY8sQ9HrGtJMHPn3Y+kVL2vHDkxICzHtGOJp2Pg53dJaeww/8nnfn2kYflrLw7I+bQ9W6MYou240tS9SGMHQ/hO68sM7W8pIto5wurR1KeLE760FfSk6bW4t/nlbOXRE0ET8OsRrJDd1o/sdOVJMM4ZmC6jC3Z1TL9XAalOJ3OSDpwwhQTjP+pcVckGln8ahD3Jtn+AxIKisODZIRfvMNbJRatWdQW+5gdIm5uKmKd3o3Q41firTnLQ5Zo/0+/CNjSt4jaAluQFyQZ2RF2jEHArFMTGaq6fAncHAnga/uxIcD8Ld3A2FDdk7BDzOQCe9+fFCR/+DWv5jZr9/U/UsKahVDZUsowRsRv6odgGTVK2E0BPqdD/ah2oy6LUhHpQhxIvi7xClz8Urs4s8Q2TPKDyUjRWubCDcQbzu29F3ljvkdIylO/fa5MSVbVpgpWI5QFcjvg39hlkyuPyj9xqZadUcL65SITVs+kcofsxakEoJtPM1RsCAwEAAQ==",
    "RawSubject": "MIGHMRUwEwYDVQQKEwxFeGFtcGxlIExhYnMxNDAyBgNVBAsTK0V4YW1wbGUgTGFicyBDeWJlciBhbmQgSW5mb3JtYXRpb24gU2VjdXJpdHkxODA2BgNVBAMTL0V4YW1wbGUgTGFicyBJbnRlcm1lZGlhdGUgQ2VydGlmaWNhdGUgQXV0aG9yaXR5",
    "Version": 0,
    "Signature": "k1cPRdYE08YQR/iC9lGCNGnG9mCjk0qPfwsidc4I3NGu2M4dxOsn1I+NlddTwQW5EYFnYCd63365F4WDY1HYFOC+t9UT6vJ4cgc+zdE+unCVkIqZyhHLFxibGhgeKD7XvLi/cisB34/ukYgxa6RitAkJFdC++rSNHTBBiQW4tekAoKeFp0FA6f0TKzVW6/AtpOx1gxFZ9nWCc+2UJ5zw6ZniFlyZ+FPLzMsy6obSUNfgyzdttCsYPWZbF+4sOf7OpjQt3eq6EBIz3QwhEv0uaZEHRCsTw6AGE1to3wcHxaxpV28VZKUff4apQFSBH0QLuAX39qvtg9g7Q7tk8hv/awtIFSet8sat5vhlV+AM2m7iWscSjlq37Lj8MaExX8XK2ZueKt9tsYZf0rZf2jzoCRazT5Y70SZYQ/PpQHu0dDZQ6m9jhMn4NQzlL3hQCOQ8EwxZesiGtqbn5vzfLIbAX2Ic/RH38Z/exYsJ2em+EWd/oTCKkxx4bAmRO4RNCnvqefXGDIxuMHiZ+v0mVSqtkXASs92AqfitAuwVlpR6SaY44WGHwFiriiZ8vDEJTMUbTVwcRC1ewOLGTQQQ7rnFyGQDop+zuiIczrv4ocIb4fdMCbFRD2GKn4/6tv/hAI5HafTYgoUkPGDZ9t/9m78NZeorKnOqoubnF1D2MX+R8sc=",
    "SignatureAlgorithm": 6,
    "PublicKeyAlgorithm": 1,
    "PublicKey": {
      "N": null,
      "E": 65537
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
    "Attributes": [
      {
        "Type": [
          1,
          2,
          840,
          113549,
          1,
          9,
          14
        ],
        "Value": [
          [
            {
              "Type": [
                2,
                5,
                29,
                19
              ],
              "Value": null
            }
          ]
        ]
      }
    ],
    "Extensions": [
      {
        "Id": [
          2,
          5,
          29,
          19
        ],
        "Critical": true,
        "Value": "MAYBAf8CAQA="
      }
    ],
    "ExtraExtensions": null,
    "DNSNames": null,
    "EmailAddresses": null,
    "IPAddresses": null,
    "URIs": null
  }
}
```