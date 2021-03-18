# Key Stores

This API endpoint will return the Key Store IDs managed by Locksmith.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/keystores`

**Method** : `GET`

**Data required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
# List Key Store IDs
curl http://pkiserver/locksmith/v1/keystores
```

And the data returned would be the minified version of the following JSON:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Listings of Key Stores"
  ],
  "key_stores": [
    "default",
    "example-labs"
  ]
}
```