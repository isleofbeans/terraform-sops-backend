## How dose the terraform-sops-backend work

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![explanation](../assets/breadcrum-explanation.drawio.svg)](./index.md)

terraform-sops-backend works as a intermediate terraform HTTP backend using a regular terraform HTTP backend but encrypting the terraform state before it is forwarded to the actual backend.  
It supports all default methods defined by the [terraform remote backend type HTTP](https://developer.hashicorp.com/terraform/language/backend/http) which are:

* GET to fetch the state
* POST to update the state
* LOCK to acquire a state lock
* UNLOCK to release a state lock

## Fetch the state

* A incoming GET request is forwarded to the configured backend
* On Status 200 the response body is tried to be decrypted using the SOPS meta information contained in the body.
* The updated response body is responded to the calling client

## Update the state

* A incoming POST request body is encrypted using the configured SOPS key(s)
* The incoming POST request is forwarded to the configured backend with the updated body.
* The backend response is responded to the calling client

## Acquire a state lock

* The incoming LOCK request is forwarded to the configured backend using the configured lock method (default: LOCK)
* The backend response is responded to the calling client

## Release a state lock

* The incoming UNLOCK request is forwarded to the configured backend using the configured unlock method (default: UNLOCK)
* The backend response is responded to the calling client
