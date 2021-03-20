# Key Pairs

This API endpoint will return either the Key Pair IDs in the store, the Public Key of a specified Key Pair ID, or if also supplied the Passphrase will return the full Key Pair.

They Key Pair ID is what is represented on the file system and is created from the `key_pair_id` parameter - in all cases the passed string is slugged.

**API Version** : Version 1 (v1)

**URL** : `/locksmith/v1/keys`

**Method** : `GET`

**Data required** : None

**Optional Data** : Key Store ID, Key Pair ID, and the Passphrase

## Input Parameters

To get a list of the Key Pair IDs in a Key Store only the `key_store_id` parameter is needed.  There is a `default` Key Store created at initialization that will be used if the Key Store ID is not specified.

To get the Public Key of a specific Key Pair ID pass the `key_store_id` and `key_pair_id` parameters.

To obtain the full Key Pair, both Private and Public pass the `key_store_id`, `key_pair_id`, and `passphrase`  parameters.

## Success Response

**Code** : `200 OK`

**Content examples**

A cURL request would look like this:

```
# List Key Pair IDs in the 'default' Key Store 
curl http://$PKI_SERVER/locksmith/v1/keys

# Get Public Key for OpenVPN Server (openvpn-server) in the 'networking' Key Store ID
curl --request GET -G --data-urlencode "key_store_id=networking" --data-urlencode "key_pair_id=OpenVPN Server" http://$PKI_SERVER/locksmith/v1/keys

# Get full key pair for OpenVPN Server (openvpn-server)
curl --request GET --data-urlencode "key_pair_id=OpenVPN Server" --data-urlencode "passphrase=s3cr3t" http://$PKI_SERVER/locksmith/v1/keys
```

And the data returned would be the minified version of the following JSON, respective of the 3 different example cURLs:

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Listing of Key Pair IDs in Key Store 'default'"
  ],
  "key_pairs": [
    "mykeypair",
    "openvpn-server",
    "vdi-terminal",
    "personal-site"
  ]
}
```

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Public Key for Key Pair ID 'OpenVPN Server' (openvpn-server) in Key Store 'networking'"
  ],
  "key_pair": {
    "public_key": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUNDZ0tDQWdFQXhYMDVpSjNPSjZJaWlvMUluMUVyVUx0d1dDNG04cmVBK2w0K1QxRTRNT3RBL1FuZTU1amoKL3ZjYlQ3b0s0ZGZpajRvZXowYlBPd2xVRy9UNlhMSDgxbUN4bFEwNGt3eXEwdWFhaFIySFAwdFhRUkEvQ096OApaY1BmT3c5N0dvcnNhWWtVMFIxOW5IT1M3RVE4aFlVbTFIbG43QTl1dktTVUIrdzErR25QSGZJeTFLK0NMYVR5CnNSMnFGZjhnZzBxQng1MC9HeitoL3N1cG1kN2s3NEF3M1Q1QVZwTGsxS1dURFRsSWg2eTF3RnNRaCtGMEdWYnEKczgycUpGREJmRmpJRXBIV3hXNFJieWpCRk04WUJoNTdydHVOakZMVEtRVXBkR29tRE8yZEVTOWk3eHRGL3AzNgoyMDVRUEVJYWI5aG41ZnRDcjI5ZWlSSjRKNy9YOGNRRG5hRlpTT2Y0VVpIcGhranZ5RmN2a3JpZXpCZmRNYnFhClFqL0loZW9kU201SG5PUXpqbzcydndrK21MRGpMNHZkVyt4WHBHbEpiZEJWT1JYdEZkUDhiM2JQVlFyQlVFNlcKME9MaEZGdEZGRDRKU2JSVXdYNFhta2xkYU55VjFkOXNHQkQvT2hFWnRjeHlOMzEwTWM4MXR2c1FGdWtJQ09VdQpQd0JGR00rWlRyOU5QYkNJeFZlTXV5YzMvTzQ5ZGRRdVlGbnRNWmZnNUR5YWNQaWpaVzZBZW0vRmJUL0pxaHFDCmswWkJRSlhzR0dVTHpkaENwVUdrbloxMm9JdHphemtyemxTeWNQRUovMmJxS2N6c05VRU8ycWlrcGU3UHBrbnkKazF1SUpRQTd0UXEyRGdQV2plRjNpOVhRZjJNdzNsdjJUQW9kM2ttT0ZuTnFqQnVaNmxrVmk0VUNBd0VBQVE9PQotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg=="
  }
}
```

```json
{
  "status": "success",
  "errors": [],
  "messages": [
    "Loaded Key Pair ID 'OpenVPN Server' (openvpn-server) from Key Store 'default'!"
  ],
  "key_pair": {
    "public_key": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUNDZ0tDQWdFQTFld0VneGk0dVZ1WXpOSFJmaXp4MzBPNDF0aTVGaWpSbHp2SUF3MUw3RjAvU2o0azNiWU8KbDNld2JGNVBqaVdWc0tjZ00wNmdMWUlJZlFtN1k4QUxUcS83TzlDVG5FU0RrTWVNb1RNV2wwMzVid0F3cmhvUwpNdnpRcmZhb1U3MHIrbjhpUndoaGtNcHdMOXNZR3M1MHp2SndRWEVnV3hENHlBeHZMTGU1NGpIZUdXeXVvRFcwCkpTc3RlZThzdkd4R2gydWt4cmtJWjY4RUx3UEJldGNkenVuczVmemd0eE92K2wydDhMZ2FEY2RoVXgreTJQM2IKZ1R2MFBFZW92cTB2bFdDTU5FQ2doQUdMclFVN1RQY0h0d3d5OWZHK09VOGd4WjFxb3NQSlJVQW85YnYxamRnbApCSjZsRDhtdmU2T3pnTllpT0dlYkxIY2pYNXV4QXV2YUZzV1NpRkxKbmlQWHI2YWVRRElmemZYNFVaZnhscFk0CldqaGZTcjhMcDUxWStpOVR6S04zVjM5VUtuRDVvaE5kUHFoTkNJLzV1UzkrVVc0elJwMkFuU2w5MFBhSnZ1Q0wKSlJIYlBKZEw0SmxxZVFLY21xOCt5VU1aOEo0RTlzNzhKWU9qbS9QMGFiOFFKZEFZV2wyMm5nU0ZvblREZmRUVQpEMGpZUmNWYmNZOVBBQzN4MTh1R1pvWGxzZHdIb2lxRzBQQlZETmFxc09KeEU1VytFTFRYVTZSN0E4dmp5Vk80CmcraDZPU1FNNjNjRWZhakNmNzl2bnhIK29NRDBaRHdYKzVkYW81RmVDQi9YRmN0YW1URzhETDlxdGpJWjc3ci8KWXoxOUxIYkdRWUlEdnorM0ZoQzlNVWlSK2hPWk5RN1ZlbFdOV3Rzc0lDczBvbEVmdVFFaWpXc0NBd0VBQVE9PQotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg==",
    "private_key": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS2dJQkFBS0NBZ0VBMWV3RWd4aTR1VnVZek5IUmZpengzME80MXRpNUZpalJsenZJQXcxTDdGMC9TajRrCjNiWU9sM2V3YkY1UGppV1ZzS2NnTTA2Z0xZSUlmUW03WThBTFRxLzdPOUNUbkVTRGtNZU1vVE1XbDAzNWJ3QXcKcmhvU012elFyZmFvVTcwcituOGlSd2hoa01wd0w5c1lHczUwenZKd1FYRWdXeEQ0eUF4dkxMZTU0akhlR1d5dQpvRFcwSlNzdGVlOHN2R3hHaDJ1a3hya0laNjhFTHdQQmV0Y2R6dW5zNWZ6Z3R4T3YrbDJ0OExnYURjZGhVeCt5CjJQM2JnVHYwUEVlb3ZxMHZsV0NNTkVDZ2hBR0xyUVU3VFBjSHR3d3k5ZkcrT1U4Z3haMXFvc1BKUlVBbzlidjEKamRnbEJKNmxEOG12ZTZPemdOWWlPR2ViTEhjalg1dXhBdXZhRnNXU2lGTEpuaVBYcjZhZVFESWZ6Zlg0VVpmeApscFk0V2poZlNyOExwNTFZK2k5VHpLTjNWMzlVS25ENW9oTmRQcWhOQ0kvNXVTOStVVzR6UnAyQW5TbDkwUGFKCnZ1Q0xKUkhiUEpkTDRKbHFlUUtjbXE4K3lVTVo4SjRFOXM3OEpZT2ptL1AwYWI4UUpkQVlXbDIybmdTRm9uVEQKZmRUVUQwallSY1ZiY1k5UEFDM3gxOHVHWm9YbHNkd0hvaXFHMFBCVkROYXFzT0p4RTVXK0VMVFhVNlI3QTh2agp5Vk80ZytoNk9TUU02M2NFZmFqQ2Y3OXZueEgrb01EMFpEd1grNWRhbzVGZUNCL1hGY3RhbVRHOERMOXF0aklaCjc3ci9ZejE5TEhiR1FZSUR2eiszRmhDOU1VaVIraE9aTlE3VmVsV05XdHNzSUNzMG9sRWZ1UUVpaldzQ0F3RUEKQVFLQ0FnRUF1M3hWSUFpa3JWK0g5Y3JXam4wWnB4R1ZpRWI5UUZ5YUJLL1NSa3A3QmpkYlp0ZzhPMHg2VVdvRwo1NU5vcWk1cW1SNkFiRGMyejJ1dHdOaXNzV240L3dmaGFyVU5DZUpLWkxOZm4xQkZObXFTZUNSMGhjSTN2UlF5CldLVmJOYmtRT0VVQVo3MEN0WUdXL1hwS0VBUnQvNG9mdEZ0UGZrRExxWmlzUDBidTFUM2JaL0VHdzBjT0VaMWMKQ0FnRTcwYitNV1c0VHFxUW9UNlVyaGZlbEtqQWFUNC83L25IZ096eVNMMmQydUdmaEFBQVhuZmpxYVlqb1lwSQpMaVNuMGlXN0ZISS9ydlFOT21TWVpCVzN1V2F3RGsvVXdoVVRJT3Ntejh0OVVCWEQ2cnVtcW1nSGVEKzlnZVVBCnlXdDhMeUowQ2pDaVlVajhpT1lKekp1SmNnUC9hR2Ira2U0V3c0UVVrTmdJUzE4VHFWcGNXeU5vQUpHTHI4Wm4KM3pwOWhLdG9ER1JrZkgrbjRoRVRqbEM0eDZlaGtsOGJHZmVjTnhzOG9CN0ZJdzAxM3k0VzYwbHpvM1ZkT3lOQQpNOEtOKzFGLzh6VjdJdlJWcGlnZjNSUGFOWGJ5K09mZDFmMFJVNlA5NTY5ZzdwSXpnQ2NTODRFcUhVT0VPTzR5Cm1hUEtyMVJQODY2aU5MQnZFdHlnMVY3bWJoNlZkVThkeTdhQlhMRnNKYWdkUmRXdzFPL3YzT3h0WnVuUUQ1eUcKYUZFWFhqb1pncGJaNmIvNEFJVFMxZE9mckdSM3Vxa3pCT1NyWlRNQXJVVDhvR3ZEN1l4UWdDVHpEbG5NdGx2Nwp5MWl2VmovUVVxMlhKVWFITFR4S09YME11cGFTSFZOWUdvSXdxR2E0SUUveFNGd2RUUkVDZ2dFQkFQZGtrMGlzCkdsWlBoWmNGWm1FNWd1SkJYdXNLVDhTVEJxV2t6N3RDTmUzTHRFRzJRUUNPZUYvcjhKcFVTRHp1VHJ0cFpENlcKS0RuZ3l0MjVFSEEybDV1T0t3QURxTkluQ3Z4MUlENGI1TU9Fc1ZBYm1va0RGdzlTbUhicUlRZFZ0bWFueHhRVgpPeUd1MTZjRXdZU0lQVnBRVytnbVVWcjhKZHdEdUliTm40T2hXaGIzSVh3cmZFYWMvNy9DN3AwdnNDZlZ5RU1FCllZb2ovRWo1MS9vQ2o3akk3emdrN3pWMnhZMVNCT2w3Z0pScHBwMVc5eXgrOTdTaExnVVJFZzVvWkFndnlsN24KUVhBVHA1YWNuRGxNRFIyd3NIc3hTbTk5d01wczc1TlRQOUlYc3J4bUZuUElNOHR1RGRXdlJld3VyaFpsRVpQTgpwYnRIeHdKTmw5ekhPczhDZ2dFQkFOMWRWS3lvZ3hGeWJZNEV5S0g0emo3dVNiMDEvS3kwaVdVdC80Y3hjbTR1CnZkUHdkWW04Skg3cHVNTlo2UkMzWFBVM20wandEQXU3YUhrNVk5NUtIalhyVStyUXpJVS9KcnI0UGQ3VG9IT3kKSEpIUzVDVGdKY1hUTFQ0Q0Nodkx4Qzh1eG9UZEEwbDc0MS9abnZNS1RSMHZNWlZIcjl4R0RQc2t5OEdXZGtmbgplNHB6aHF2Q0g3Mi9aeWQwempOSTFKNm90SWFjMW9XSkpnRzVnZHNMRHJTVHVYbWszS0tyV1kzdVRkbUlObUpqCk8wNmxtdkk0OCtWeXkxR2VNN3lWeGEwS0cyY0podmZERXAzWTI0TytTTFVSMlhvWjd5QWVEbGRISnB4QUdvMVEKeGtEbU50MEFkR0J1VFVoT2Q2V1lNRFhUdGhUTnl3eUY5NWR4MUtMZGVxVUNnZ0VCQU53VTd5blJZVVN6VGNiQwpHUWdaSDZTa1B3cWRpOFQyZncxUkJ6UXhmTVJsV2FDendEUDhpbjNhNlpxQnJCbjRicll1MWUwUHJBMkJPemZ3CmNQMUNzN0RBMHVRYVhVOUhTSEM2eWNvM0NsWWRiNWd5VmxIWkcrU1h0K2JoOWl3T1Jrd2dxZXZsejByeHZndHAKSWJjRGRJRXB4L2xJVFV2QjBQUmZvd0xaWGpTOWorV3FTSEdzUmN5VDBya0hjenNHdDVGWWorVit1ajhvTUVIRApjaTJKcGMzZmcyRFJDclRuU211a01aWjhOakRScEZXSWppOVpiSWVXYzlneURYd2Z4ZzI2WmkvelRyV2o1bzBJCkdicW5PMnZVU2N2dVY2ZkRtWVQ3VUU0aDJ6N05za1lFRTZsQXkwTUlUdXB3R0tZNkNNa0hkSkdtZXUrV3RTWUoKWFRZZFR1VUNnZ0VCQU02OFdWWWUwcm02bW1KbVNWSXI0Y0tZSExuZTc2b0R1Y1dLM2ZoT3o4WGpWVm5ZV28xVgo3dWV3TStRTjFrTE1YTDZQUGpFeUxxM09TdFhjS1U3eS9aL3h0Wis1ZlNoOFFCbWh1WGFmUWx5SzNXKzYrMk1OCkMzbmpyWDhadklNVkhKWE1JNDcyTWhtdzREc21MUEppam41UkV3ZU51Y29JaWhzSzFGaHB3dkdJV0xLSERpRGUKM1hJQ2piNGxzbVhuQU50a1I2VG9XTmpCcTRNMDB2ZlZMZGlybGk3ckx6dWt0N0I3L0t6S0w3QlhhSTRjejhhawpOZlAxNzdpNy9TbUUzdWFxWjhrazlxM3h1ek03MGxjSm9UR3FCK2Vtek5LNy96eTNzSEdBMU10aHdxWGQyeU12ClI0Qy93dUZpbHc5S1FNd2tld3FXMzZsRWZHVXQ1QjV1cGhrQ2dnRUFITzNRSm1BQXk1Z2hhRUttbE9kR21FdXIKaFV5WDN4ZTA0M1NvbWlDUGRvRGpuUE1ETEhlYjdGc0dQckFYdytMRGpSMXl2M2xTTUJubTRFSWZrVU5VandTeApsZFlNQzQxRmNHRjJrNVlEVWFnNGNaMnEyeUhrdCs0WFpya1Qwc1RKL083SXZmSm80Vjc4QWtwVXBoSUtKSnZFCitSS3EydW9nRUo1SkNWMW4xT1liUHlvTjhkaUVmSlZzMmk3NHMwQmFnY21zMlRRZ1MyWitoSTVaNE5TZ3BPQmMKZXBvR1JaOXdtandPc0E4VWg2QXhJWFkzM3kvWE1ob2tySU5CWFNSWmMwcEl2TFVMVjZ0ejNvbCtEUXY2akIxcwpqVEhvQ3FLV2pVem1EVkxZRnNVYk5ZM21UTWhsem1MZzk1Umkyd3FQWG00OWFNTHVlVUFBRWJaUFFWc1B4UT09Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="
  }
}
```

## Error Response

**Condition** : If provided data is invalid, missing, or a system error occurs.

**Code** : `200 OK` - 200 OK is returned even on errors due to no specific HTTP error codes corellating to different processes taken place during the Key Store generation workflow.  Check the `status` field for specific status matches.

**Content example** :

```json
{
  "status": "key-pair-id-missing",
  "errors": ["Key Pair ID parameter missing!  Pass with `key_pair_id`"],
  "messages": []
}
```

## Return Statuses

Potential return statuses sent back via JSON are as follows:

- `suceess` - Successfully retrieved Key Pair
- `invalid-key-store` - Provided Key Store is invalid
- `private-key-decryption-error` - Issue decrypting Private Key
- `no-private-key` - No Private Key stored for Key Pair
- `invalid-key-pair-id` - Invalid Key Pair ID in specified Key Store
- `empty-key-store` - Key store is empty