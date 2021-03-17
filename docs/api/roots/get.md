# List Root Certificate Authorities

Get the slugged list of Root Certificate Authorities.  The slug is a DNS/file-safe filter on the CA CommonName and used to query and use the Root CA in other workflows.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/roots`

**Method** : `GET`

**Data required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
curl http://pkiserver/locksmith/v1/roots
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [],
  "roots": [
    "example-labs-root-certificate-authority",
    "example-labs-iit-root-ca",
    "example-labs-skunkworks-certificate-authority"
  ]
}
```

## Notes

* This function is basically just a directory listing.