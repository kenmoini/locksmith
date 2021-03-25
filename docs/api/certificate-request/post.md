# Create Certificate Request along Certificate Path

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/certificate-request`

**Method** : `POST`

**Content Type** : `JSON`

**Input Data Structure**

```
{
  "cn_path": string,
  "slug_path": string,
  "certificate_config": {
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
    "certificate_type": string, // optional
    "rsa_private_key": string, // optional
    "rsa_private_key_passphrase": string, // optional
    "expiration_date": []int, // [ years, months, days ]
    "san_data": { // optional
      "email_addresses": []string, // optional
      "uris": []string // optional
    }
  }
}
```

**Input Data examples**

```
{
  "cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority/Example Labs Signing Certificate Authority",
  "certificate_config" {
    "subject": {
      "common_name": "Example Labs OpenVPN Server",
      "organization": ["Example Labs"],
      "organizational_unit": ["Example Labs Cyber and Information Security"]
    },
    "expiration_date": [
      1, // Years
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
  --data '{"cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority/Example Labs Signing Certificate Authority", "certificate_config":{"subject": {"common_name": "Example Labs OpenVPN Server", "organization": ["Example Labs"], "organizational_unit": ["Example Labs Cyber and Information Security"]}, "expiration_date": [1,0,1], "san_data": {"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://$PKI_SERVER/locksmith/v1/certificate-request
```

## Success Responses

**Code** : `200 OK`

**Content example** : Response will reflect back the slugged ID of the certificate, the next certificate serial number, and the full representation of the generated CA Certificate.

```json
{
  "status": "certificate-request-created",
  "errors": [],
  "messages": [
    "Successfully created Certificate Request Example Labs OpenVPN Server in 'Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority/Example Labs Signing Certificate Authority'!"
  ],
  "csr_info": {
    "slug": "example-labs-openvpn-server",
    "csr_pem": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJRXVEQ0NBcUFDQVFBd2N6RVZNQk1HQTFVRUNoTU1SWGhoYlhCc1pTQk1ZV0p6TVRRd01nWURWUVFMRXl0RgplR0Z0Y0d4bElFeGhZbk1nUTNsaVpYSWdZVzVrSUVsdVptOXliV0YwYVc5dUlGTmxZM1Z5YVhSNU1TUXdJZ1lEClZRUURFeHRGZUdGdGNHeGxJRXhoWW5NZ1QzQmxibFpRVGlCVFpYSjJaWEl3Z2dJaU1BMEdDU3FHU0liM0RRRUIKQVFVQUE0SUNEd0F3Z2dJS0FvSUNBUURKZWQwcUlBNEc0c0pLMi9wQUc3KzlKOHA1WGdjTTdWZVpsK0Z2RVZ5TQpCNmJoaWZWR3lidGUvY3V0cWkwWmdvb3BhT2IvdGl0bmJZNnRtT1hoSDBXLy9yTXBjaXE3b29sNjUrbW1HZEIxCmxGZ0l6eStabjIvcDdZbUFWQms3ODZqN1Y3ZGkyUTdnQ3B3VGh5bEwzL2czeEhmRWRGWG9vTTJrMmIvZDVEazcKK0x2dWluamhzL3pPMDNhcWZzZGxwUFBlWWlydW1tTjlIOGlXRnZKelhYNEMxV0RoNWhCTTNUb1lIMWZWclAzRwpBYVNDRUEyL0pTNnR3Z1hJZEZTV3BrQ3dxaEZoQVBwYXp3bDhVdGJWUGNjay9WMlc4U05VN3VEa2RqL2pwWk16Ci9vWEtSMmY2R3RGOWpVZnN3M3VpSHoveVZHUzNnb0RIRHJidFR1bEhDYjF0ZG4waXhsaytpbW5OWVpIR0ZLOTUKK0U1RDljUmJ5ejdJOVhGU2M2b2gydjFZam5ZKzJNWUZzSDEvS08rUDMwRzdMTHMzN2k1YVlRazJGbkM4a1ZRQwplY3BQTyt3WGN0T0lsSkxWR20reWJFamQ4bTYwQVI4KzBZQmlZSFFLU0taSjRpK09WYzBtYWorNStWQTJnRGFvCjFKcnJkOWJmeU1mSEJBcDZwYnc3bE9UVkR4bEdmR28wanozTEJrZ1NQdWxDNkNlb2JKL2traHBqUmlMdHlvVFQKWHJaTGtTL2Y4bmx4N1I3cS9SakRqdURZOUdwUlhjaWkybEIzYWZBZWZtRWZFNkxveEVpaVpUUHlXNVpPT2c1Wgo3YkltaXh0QkZMTTdmM2JVRU5zUERZK2JhUXp0eThJQTlnc21CZU9vcXJTWGJuK2dmZDhDdnVHc1J5T3dZdURNCk9RSURBUUFCb0FBd0RRWUpLb1pJaHZjTkFRRU5CUUFEZ2dJQkFDVVR2OEtTQkJiY2VCN0VQdWd5WUhKeVdZb20KcmVweVB3bGhhdC9mT1Nza21pSzIxTlovRTJmOUVON3h5OGIwQXFibjJTRVYxLzhUeEhUVTFoQjNIM2dWeE9Kcwpxc1NsbjdPYmJMUVpKSzJmUDhTMDlHYjZlckVFT2cwR2NQNTkxWVE4VHRTUGFQeFdLcEprNlBmcnNYZzBSYVQ5CjA3MDV4Wk1KanV4VHN6eGVqeHdrbnhqNGZ2SkZkY3hQVDFuRURGUi9UNG9ZK21oTVQzVWxZOTF6c204eUQ0T20KR25hRkpyUzVlUWQzN21mcXVQR0htSW9QUDBOdS94Z3BmdkxVYUc2aGdwcWNONXpaT1VCUCtMQzdtRll0WmU0UQpTSlEyajcxb2h5c3lEWVhCdFFVKzBWK3g1WjFmcUVzQVI0Y0FUOGVPTlhYVnhpMjgxczEzampRRERlUm5hSFJICmFGVzNUdGxVN1Y1eS80OU45SldHZkthd3pSSm5KT3QweGFJWS9uL1BkbXRLS2wrM2p5bWQ2TmNwVG9ZR0RiNEkKclFIVk9SV2ttNU00Rk1CR2NhZkUxSXdYZWdLLzd5T1E3TWZDa05LYnlCaDhyMk9Vc1ozTzRYZGkvRnoyem4yVQpCK202ay9rYVNsL1JUa0ZNaWUwV1RnR3p6Ymx2VGhUODljVEc3Qy9qNTMzMUh6Y1R0SFRtdmMwemN0MDg0Y1BMCjdxMkx5ZFB2QlN5d2N3QWNDR2dhWmdOSTVsNDdYeHdRdFZQWkEyTU5QN2JCbmdFR3k2UE9vSDd1RGhTbDg5N1MKNFpsWDNCR1BiU1VMdWZVWm5IbGM0ajFERmMwaGJYOU1TVHI2MHN3Sm9menVMYnFDclNQeHdNdWxpSVVMVENHdgpkRTduZnd1YVFwWTM4bXJLCi0tLS0tRU5EIENFUlRJRklDQVRFIFJFUVVFU1QtLS0tLQo=",
    "certificate_request": {
      "Raw": "MIIEuDCCAqACAQAwczEVMBMGA1UEChMMRXhhbXBsZSBMYWJzMTQwMgYDVQQLEytFeGFtcGxlIExhYnMgQ3liZXIgYW5kIEluZm9ybWF0aW9uIFNlY3VyaXR5MSQwIgYDVQQDExtFeGFtcGxlIExhYnMgT3BlblZQTiBTZXJ2ZXIwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDJed0qIA4G4sJK2/pAG7+9J8p5XgcM7VeZl+FvEVyMB6bhifVGybte/cutqi0ZgoopaOb/titnbY6tmOXhH0W//rMpciq7ool65+mmGdB1lFgIzy+Zn2/p7YmAVBk786j7V7di2Q7gCpwThylL3/g3xHfEdFXooM2k2b/d5Dk7+Lvuinjhs/zO03aqfsdlpPPeYirummN9H8iWFvJzXX4C1WDh5hBM3ToYH1fVrP3GAaSCEA2/JS6twgXIdFSWpkCwqhFhAPpazwl8UtbVPcck/V2W8SNU7uDkdj/jpZMz/oXKR2f6GtF9jUfsw3uiHz/yVGS3goDHDrbtTulHCb1tdn0ixlk+imnNYZHGFK95+E5D9cRbyz7I9XFSc6oh2v1YjnY+2MYFsH1/KO+P30G7LLs37i5aYQk2FnC8kVQCecpPO+wXctOIlJLVGm+ybEjd8m60AR8+0YBiYHQKSKZJ4i+OVc0maj+5+VA2gDao1Jrrd9bfyMfHBAp6pbw7lOTVDxlGfGo0jz3LBkgSPulC6CeobJ/kkhpjRiLtyoTTXrZLkS/f8nlx7R7q/RjDjuDY9GpRXcii2lB3afAefmEfE6LoxEiiZTPyW5ZOOg5Z7bImixtBFLM7f3bUENsPDY+baQzty8IA9gsmBeOoqrSXbn+gfd8CvuGsRyOwYuDMOQIDAQABoAAwDQYJKoZIhvcNAQENBQADggIBACUTv8KSBBbceB7EPugyYHJyWYomrepyPwlhat/fOSskmiK21NZ/E2f9EN7xy8b0Aqbn2SEV1/8TxHTU1hB3H3gVxOJsqsSln7ObbLQZJK2fP8S09Gb6erEEOg0GcP591YQ8TtSPaPxWKpJk6PfrsXg0RaT90705xZMJjuxTszxejxwknxj4fvJFdcxPT1nEDFR/T4oY+mhMT3UlY91zsm8yD4OmGnaFJrS5eQd37mfquPGHmIoPP0Nu/xgpfvLUaG6hgpqcN5zZOUBP+LC7mFYtZe4QSJQ2j71ohysyDYXBtQU+0V+x5Z1fqEsAR4cAT8eONXXVxi281s13jjQDDeRnaHRHaFW3TtlU7V5y/49N9JWGfKawzRJnJOt0xaIY/n/PdmtKKl+3jymd6NcpToYGDb4IrQHVORWkm5M4FMBGcafE1IwXegK/7yOQ7MfCkNKbyBh8r2OUsZ3O4Xdi/Fz2zn2UB+m6k/kaSl/RTkFMie0WTgGzzblvThT89cTG7C/j5331HzcTtHTmvc0zct084cPL7q2LydPvBSywcwAcCGgaZgNI5l47XxwQtVPZA2MNP7bBngEGy6POoH7uDhSl897S4ZlX3BGPbSULufUZnHlc4j1DFc0hbX9MSTr60swJofzuLbqCrSPxwMuliIULTCGvdE7nfwuaQpY38mrK",
      "RawTBSCertificateRequest": "MIICoAIBADBzMRUwEwYDVQQKEwxFeGFtcGxlIExhYnMxNDAyBgNVBAsTK0V4YW1wbGUgTGFicyBDeWJlciBhbmQgSW5mb3JtYXRpb24gU2VjdXJpdHkxJDAiBgNVBAMTG0V4YW1wbGUgTGFicyBPcGVuVlBOIFNlcnZlcjCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMl53SogDgbiwkrb+kAbv70nynleBwztV5mX4W8RXIwHpuGJ9UbJu179y62qLRmCiilo5v+2K2dtjq2Y5eEfRb/+sylyKruiiXrn6aYZ0HWUWAjPL5mfb+ntiYBUGTvzqPtXt2LZDuAKnBOHKUvf+DfEd8R0VeigzaTZv93kOTv4u+6KeOGz/M7Tdqp+x2Wk895iKu6aY30fyJYW8nNdfgLVYOHmEEzdOhgfV9Ws/cYBpIIQDb8lLq3CBch0VJamQLCqEWEA+lrPCXxS1tU9xyT9XZbxI1Tu4OR2P+OlkzP+hcpHZ/oa0X2NR+zDe6IfP/JUZLeCgMcOtu1O6UcJvW12fSLGWT6Kac1hkcYUr3n4TkP1xFvLPsj1cVJzqiHa/ViOdj7YxgWwfX8o74/fQbssuzfuLlphCTYWcLyRVAJ5yk877Bdy04iUktUab7JsSN3ybrQBHz7RgGJgdApIpkniL45VzSZqP7n5UDaANqjUmut31t/Ix8cECnqlvDuU5NUPGUZ8ajSPPcsGSBI+6ULoJ6hsn+SSGmNGIu3KhNNetkuRL9/yeXHtHur9GMOO4Nj0alFdyKLaUHdp8B5+YR8ToujESKJlM/Jblk46DlntsiaLG0EUszt/dtQQ2w8Nj5tpDO3LwgD2CyYF46iqtJduf6B93wK+4axHI7Bi4Mw5AgMBAAGgAA==",
      "RawSubjectPublicKeyInfo": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyXndKiAOBuLCStv6QBu/vSfKeV4HDO1XmZfhbxFcjAem4Yn1Rsm7Xv3LraotGYKKKWjm/7YrZ22OrZjl4R9Fv/6zKXIqu6KJeufpphnQdZRYCM8vmZ9v6e2JgFQZO/Oo+1e3YtkO4AqcE4cpS9/4N8R3xHRV6KDNpNm/3eQ5O/i77op44bP8ztN2qn7HZaTz3mIq7ppjfR/Ilhbyc11+AtVg4eYQTN06GB9X1az9xgGkghANvyUurcIFyHRUlqZAsKoRYQD6Ws8JfFLW1T3HJP1dlvEjVO7g5HY/46WTM/6Fykdn+hrRfY1H7MN7oh8/8lRkt4KAxw627U7pRwm9bXZ9IsZZPoppzWGRxhSvefhOQ/XEW8s+yPVxUnOqIdr9WI52PtjGBbB9fyjvj99Buyy7N+4uWmEJNhZwvJFUAnnKTzvsF3LTiJSS1RpvsmxI3fJutAEfPtGAYmB0CkimSeIvjlXNJmo/uflQNoA2qNSa63fW38jHxwQKeqW8O5Tk1Q8ZRnxqNI89ywZIEj7pQugnqGyf5JIaY0Yi7cqE0162S5Ev3/J5ce0e6v0Yw47g2PRqUV3IotpQd2nwHn5hHxOi6MRIomUz8luWTjoOWe2yJosbQRSzO3921BDbDw2Pm2kM7cvCAPYLJgXjqKq0l25/oH3fAr7hrEcjsGLgzDkCAwEAAQ==",
      "RawSubject": "MHMxFTATBgNVBAoTDEV4YW1wbGUgTGFiczE0MDIGA1UECxMrRXhhbXBsZSBMYWJzIEN5YmVyIGFuZCBJbmZvcm1hdGlvbiBTZWN1cml0eTEkMCIGA1UEAxMbRXhhbXBsZSBMYWJzIE9wZW5WUE4gU2VydmVy",
      "Version": 0,
      "Signature": "JRO/wpIEFtx4HsQ+6DJgcnJZiiat6nI/CWFq3985KySaIrbU1n8TZ/0Q3vHLxvQCpufZIRXX/xPEdNTWEHcfeBXE4myqxKWfs5tstBkkrZ8/xLT0Zvp6sQQ6DQZw/n3VhDxO1I9o/FYqkmTo9+uxeDRFpP3TvTnFkwmO7FOzPF6PHCSfGPh+8kV1zE9PWcQMVH9Pihj6aExPdSVj3XOybzIPg6YadoUmtLl5B3fuZ+q48YeYig8/Q27/GCl+8tRobqGCmpw3nNk5QE/4sLuYVi1l7hBIlDaPvWiHKzINhcG1BT7RX7HlnV+oSwBHhwBPx441ddXGLbzWzXeONAMN5GdodEdoVbdO2VTtXnL/j030lYZ8prDNEmck63TFohj+f892a0oqX7ePKZ3o1ylOhgYNvgitAdU5FaSbkzgUwEZxp8TUjBd6Ar/vI5Dsx8KQ0pvIGHyvY5Sxnc7hd2L8XPbOfZQH6bqT+RpKX9FOQUyJ7RZOAbPNuW9OFPz1xMbsL+PnffUfNxO0dOa9zTNy3Tzhw8vurYvJ0+8FLLBzABwIaBpmA0jmXjtfHBC1U9kDYw0/tsGeAQbLo86gfu4OFKXz3tLhmVfcEY9tJQu59RmceVziPUMVzSFtf0xJOvrSzAmh/O4tuoKtI/HAy6WIhQtMIa90Tud/C5pCljfyaso=",
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
        "CommonName": "Example Labs OpenVPN Server",
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
            "Value": "Example Labs OpenVPN Server"
          }
        ],
        "ExtraNames": null
      },
      "Attributes": null,
      "Extensions": null,
      "ExtraExtensions": null,
      "DNSNames": null,
      "EmailAddresses": null,
      "IPAddresses": null,
      "URIs": null
    },
    "key_pair": {
      "public_key": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUNDZ0tDQWdFQXlYbmRLaUFPQnVMQ1N0djZRQnUvdlNmS2VWNEhETzFYbVpmaGJ4RmNqQWVtNFluMVJzbTcKWHYzTHJhb3RHWUtLS1dqbS83WXJaMjJPclpqbDRSOUZ2LzZ6S1hJcXU2S0pldWZwcGhuUWRaUllDTTh2bVo5dgo2ZTJKZ0ZRWk8vT28rMWUzWXRrTzRBcWNFNGNwUzkvNE44UjN4SFJWNktETnBObS8zZVE1Ty9pNzdvcDQ0YlA4Cnp0TjJxbjdIWmFUejNtSXE3cHBqZlIvSWxoYnljMTErQXRWZzRlWVFUTjA2R0I5WDFhejl4Z0drZ2hBTnZ5VXUKcmNJRnlIUlVscVpBc0tvUllRRDZXczhKZkZMVzFUM0hKUDFkbHZFalZPN2c1SFkvNDZXVE0vNkZ5a2RuK2hyUgpmWTFIN01ON29oOC84bFJrdDRLQXh3NjI3VTdwUndtOWJYWjlJc1paUG9wcHpXR1J4aFN2ZWZoT1EvWEVXOHMrCnlQVnhVbk9xSWRyOVdJNTJQdGpHQmJCOWZ5anZqOTlCdXl5N04rNHVXbUVKTmhad3ZKRlVBbm5LVHp2c0YzTFQKaUpTUzFScHZzbXhJM2ZKdXRBRWZQdEdBWW1CMENraW1TZUl2amxYTkptby91ZmxRTm9BMnFOU2E2M2ZXMzhqSAp4d1FLZXFXOE81VGsxUThaUm54cU5JODl5d1pJRWo3cFF1Z25xR3lmNUpJYVkwWWk3Y3FFMDE2MlM1RXYzL0o1CmNlMGU2djBZdzQ3ZzJQUnFVVjNJb3RwUWQybndIbjVoSHhPaTZNUklvbVV6OGx1V1Rqb09XZTJ5Sm9zYlFSU3oKTzM5MjFCRGJEdzJQbTJrTTdjdkNBUFlMSmdYanFLcTBsMjUvb0gzZkFyN2hyRWNqc0dMZ3pEa0NBd0VBQVE9PQotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg==",
      "private_key": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS1FJQkFBS0NBZ0VBeVhuZEtpQU9CdUxDU3R2NlFCdS92U2ZLZVY0SERPMVhtWmZoYnhGY2pBZW00WW4xClJzbTdYdjNMcmFvdEdZS0tLV2ptLzdZcloyMk9yWmpsNFI5RnYvNnpLWElxdTZLSmV1ZnBwaG5RZFpSWUNNOHYKbVo5djZlMkpnRlFaTy9PbysxZTNZdGtPNEFxY0U0Y3BTOS80TjhSM3hIUlY2S0ROcE5tLzNlUTVPL2k3N29wNAo0YlA4enROMnFuN0haYVR6M21JcTdwcGpmUi9JbGhieWMxMStBdFZnNGVZUVROMDZHQjlYMWF6OXhnR2tnaEFOCnZ5VXVyY0lGeUhSVWxxWkFzS29SWVFENldzOEpmRkxXMVQzSEpQMWRsdkVqVk83ZzVIWS80NldUTS82RnlrZG4KK2hyUmZZMUg3TU43b2g4LzhsUmt0NEtBeHc2MjdVN3BSd205YlhaOUlzWlpQb3BweldHUnhoU3ZlZmhPUS9YRQpXOHMreVBWeFVuT3FJZHI5V0k1MlB0akdCYkI5Znlqdmo5OUJ1eXk3Tis0dVdtRUpOaFp3dkpGVUFubktUenZzCkYzTFRpSlNTMVJwdnNteEkzZkp1dEFFZlB0R0FZbUIwQ2tpbVNlSXZqbFhOSm1vL3VmbFFOb0EycU5TYTYzZlcKMzhqSHh3UUtlcVc4TzVUazFROFpSbnhxTkk4OXl3WklFajdwUXVnbnFHeWY1SklhWTBZaTdjcUUwMTYyUzVFdgozL0o1Y2UwZTZ2MFl3NDdnMlBScVVWM0lvdHBRZDJud0huNWhIeE9pNk1SSW9tVXo4bHVXVGpvT1dlMnlKb3NiClFSU3pPMzkyMUJEYkR3MlBtMmtNN2N2Q0FQWUxKZ1hqcUtxMGwyNS9vSDNmQXI3aHJFY2pzR0xnekRrQ0F3RUEKQVFLQ0FnQXYvU0RhdWN2ZGhBRjNSekl5TnVuU3FqbWw4dW1IQUxsTzBraFY1aksrLzh1V0NRQXRIanZOQW5LVApLT2VaSGVpK3VFZmRQSXpXRTloYUxRTUVQaWlrOUl2RUlYZGdQZlMxRzZ3aGJpQ2pBUFIvRktwbjB2d2JJZ01RClYvZXl1ZlRUK1M2ckVyeGlUT1NrR2h1U0FRVGtjNTE3WTZKYXlJSnk4NUtwellSOGJtQ0ZEdUtBRUJqMVFwVVAKUXlkSFpLVFpvVlJNaE9XUmxoSjIvWHcrVWxTRFpFT3hTdFV6R2JhT3JGaUZnckRuaXRpZVNpaWNFTVV2aWZsSwpwN3JHTDA5VlJRemlxQkw1c2pxMkxCMFRxYVZYZ0NuY25BOG9XY1dqWlM4T2tBK2g0TXdKUkR6VjY3RmVVRnFECllJU0FOeVZLR1NQc09kOW1pbkN6MFlucTY3ejhldm0vdGlGTVk2bDZxcklWdDlETFBUa3RybUlqY3ZtY0RtQlgKM2hHaVN6S2pOZEtwOC92UEZLRlpYelkyVXY2UVQwalpZVURmTEhiOFlldGxtM0YwUkJsRTFITzQwdEdqVWlUcwo2aG42L3loTW5kZnpQR3BSS2pRSFdwN1hJRjNQTFlqaHpkYUU2L0lGRHhnOGRBZFUzSS9Fa3FvVDdrUWhqdGxiCndmTFh5NzlqdFNBaDZLYTRiNDNwZzVRU0tNeTZvcTZWMTRyNk5oQnZwb0dZUnltdjJ5dldEOEoxbU9NODg1RGUKL0JWa0xCcFBhYWswMnNNWGthYjluVXJrVHZpOHJ5YXdkRytka2FPZnRtQk1kMjZTaWFsby8vTXpFeldYSWJrMQphWXNxOFZnbjZmUlh3U2NHOXZMN2UzT3BHTnRrWW1Zc1pJUjlnaGhMeVBJU3oya05TUUtDQVFFQXpmYzVQWWxVCjd0WTNiQmYrM0VMRHJPczlTbS9FNXdaekhpb2srbjluMnNlQ0dscTVXbTFoNkpzM3BSNlRieGw2cm4xRnBSMzMKaEQzckFYNjVrbUQvQlFzOEZVVjRRc21rR0hhWHg0TVJyOTEvYllxNXZmc1J3UFNzaFhDRGZ0OU1zYkdsVVVnMAprNzBoZ1JGUlFZdVk2bjc3THhYNWd0ZTlkRjRUWWVVTDV1d1NVSzVZL2pTSkRuOW5Cb2ZZOWxCcmtNVitnSXNlCisvRmZyUzhYOVFCalg5SDRML1VNQnlrdjRaOTZuZkYrT2FKY3JjOVk2TXRpWnB3SHp2NHVBdG5vaW5NbzdCMFUKUDlqcW5kRlJ6RDg5S3lmd2xkZmJmeXAwTEVhNU0yb1VkU0hJTnpGbDhOMHdGZFRMdVRtUUwyV2thWXdjV2E3RgozS2o2ZTU3aEZybkJQd0tDQVFFQSttdHVpYy80T0hWbnBLVTdFZUI2RzIxWDBvZGRvMXV0SHQ5dHlVM2h0UldPCmF4YjNUQkdMdTBnbmp4c3ZwckF2TmRXaUZzYzhmM29aeFBOMmI3Sk9ES1lBTW94THREN3RUN3FETStvRXllRDEKaG8zVUUrSno1eXAvL1BWRFRNTmJFb01hZDZ4S3ByMFg2UVdhS1pUb2JQZ0JLemtxVy9wVzJvSEJERzZ3OXlrbQpIaWw1eTk1WEFjVnN6REpuTUVjdUFkcnNRYzQ1WW15LzZxc2JJSDllb1hPQ0QyTHp2dExlUVhTMFdYQ3NBUEVGCllNMVBMZW5sQnhCUndNbVIyN1pBemRBMkFLUy8wUXFWYjJEbXA0MUl4ZUxqeGVYRk1CdDhNVXpxeE1RL3NnOW8KVEthSUQyU1p6azMrZkpEVXZnTzRQdjlIT0t4eTlkOGUvemFOM3ZjY2h3S0NBUUVBb3NUTnA2UHdTdzlmblQwWgpYRmdtNjNDOGJ5Y0ZKTTRrQzZLaXRwUVpMdnljQk5mTnczak45MVV6RkhxbGFSZHBySnV6ZmxuQVVmSGMxc0dmClJkOEJxcXJHUU1rMTBSSXBiR3ZNWnc5ZDJ1M2cxbURiaVJmeFg5djh1emUvczNRazJBamI3UEJ6SEk1Sy9BVUQKZ2hrZ0w3RktNRnZkWTFtN2owc0paa1BzdEFHOE42YVJEZFBXdkc1U3JRYU9uNW5PYUFxcmZrcHpvZ2VPNVA1aQpvR3crSEd1REIrTlFMaGlPam0vS0p5ZkI2U28ycytVNURrQXM4NG83WVluZU1zS2kwMGRPLzhtN1J2blY4QUtMCjhpM0gyV01tN2tRNFlyYmFPR05yMlFYc1JPVDlwU0NVdjFVTnV6TUFETkZBOFRRU1NwYy9rR0JlWFpQczMrVWYKWFNaUFlRS0NBUUJPelNYQXFqZ1RGQ2JrTWJhUDNwS3VOTTlSQ1pYV3hRK0tTb2JTdFBaVXRJN1hkaWVsd1ZPMgpRSE5xWGdTMXNIVjZ2RnBBVHJ2ajVYbGNkN3lLVTVLcCtrYlBvVVJsV1BQMmhkdXBwM2VRUzRFWHNXUE9TaEZzCjZmdlNqeDk1ZFhRZ05DOU0vMk9TYXFpdWhEdkozL2p2Nlc3OHVnVnhZaXFZb2dJc3RseHJ1b0FyTjZRREdsbEkKem1aNUwxYzNZdjdBU0xMVjFsNUtjYXhHM1VjeEI4T3dqSmVkM2VhVDR1bGJzYXpiQkZDc0R5eEJGUHBZbVdTZgp2MnZxZmNPdlh2K3ZoRmlxQlMzelN5QlJKeTRPQmJDanpNMGVSanF4ODhRMkExMVJRK3hEVFFQbU12VlgyckZuCnNUVm0zM0NDeHNyZzBCWUthSUhZaXpqRzJOVDJGODdiQW9JQkFRQ3ZWcW5sWDA2MWJGZk00a2cwMHZpanIydEIKN1dkVjgzd2JseWR1eXpwZ2V5Q3FBR0pNWGJEcHNGa1BxMnZZQk9xUEs5YVIyQ2c2VmpuWjZQcGFQVEl5UTU5MwpDTGJUVTFzNXZVNXRWOW8zMi95UDg3WThhVGV4bzA2NElZVDhEMCtnUXQ5MXpFL25xWC9zMkRnOUZReUl0QlBRCnpZcE9mNkM2VGsxcVNhNTUzNDRZbE9jKy9TNHBON2dBcFppR0hBeWZvKy9ESE56L0FPNk9GVWQ4dWo0T1hTZlUKSUhsR2cxczZaUFBiRkFHUWtHOENUWVdpRXArUks3NGIvaEpOenBweDUxZlNoR2NTSFRrWEljNWRZZStaZllyOApoYzJqK2w0NkNEQUZMQzlUbitwS0dXMEVoQmhPZW9xbXZxMGk3Z3JsdHBmd0k2ZWZ1cituVWNEaHNkczkKLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K"
    }
  }
}
```