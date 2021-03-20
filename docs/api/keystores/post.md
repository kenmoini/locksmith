# Create a new Key Store

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/keystores`

**Method** : `POST`

**Content Type** : `JSON`

**Input Data Structure**

```json
{
  "key_store_name": string
}
```

**Input Data examples**

```json
{
  "key_store_name": "Example Labs"
}
```

**Request Example**

A cURL request would look like this:

```
curl --header "Content-Type: application/json" --request POST \
  --data '{"key_store_name": "Example Labs"}' http://$PKI_SERVER/locksmith/v1/keystores
```

## Success Responses

**Code** : `200 OK`

**Content example** : Response will reflect back the slugged ID of the Key Store Name.

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Key Store successfully created!"
  ],
  "key_store_id": "example-labs"
}
```

## Error Response

**Condition** : If provided data is invalid, missing, or a system error occurs.

**Code** : `200 OK` - 200 OK is returned even on errors due to no specific HTTP error codes corellating to different processes taken place during the Key Store generation workflow.  Check the `status` field for specific status matches.

**Content example** :

```json
{
  "status": "key-store-name-missing",
  "errors": ["Key Store name parameter missing!  Pass with `key_store_name`"],
  "messages": []
}
```

## Return Statuses

Potential return statuses sent back via JSON are as follows:

- `success` - Successfully created Key Store
- `key-store-name-missing` - Missing Key Store Name parameter
- `key-store-creation-failed` - Errors are dependant on part of the workflow that failed, such as missing fields or system errors

