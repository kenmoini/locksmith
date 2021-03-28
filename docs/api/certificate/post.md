# Create Certificate along Certificate Path

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/certificate`

**Method** : `POST`

**Content Type** : `JSON`

**Input Data Structure**

```
{
  "cn_path": string,
  // or
  "slug_path": string,

  "csr_input": {
    "from_pem": string,
    // or
    "from_ca_path": {
      "target": string,
      "cn_path": string,
    }
  },
  "signing_key_passphrase": string // optional
}
```

**Input Data examples**

```
{
  "cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority/Example Labs Signing Certificate Authority",
  "csr_input" {
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
  http://$PKI_SERVER/locksmith/v1/certificate
```

## Success Responses

**Code** : `200 OK`

**Content example** : Response will reflect back the slugged ID of the certificate, Certificate PEM encoded in Base64, and the full representation of the generated Certificate.

```json

```