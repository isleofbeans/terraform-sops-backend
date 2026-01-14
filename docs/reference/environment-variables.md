# Configuration using environment variables

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![reference](../assets/breadcrum-reference.drawio.svg)](./index.md)

|                                    |                                         |                                                                |             |
| ---------------------------------- |-----------------------------------------|----------------------------------------------------------------| ----------- |
| TRANSFORM_AGE_PRIVATE_KEY          | optional                                | private AGE key to decrypt terraform state                     |             |
| TRANSFORM_AGE_PUBLIC_KEY           | required                                | public AGE key to encrypt terraform state                      |             |
| BACKEND_LOCK_METHOD                | optional                                | lock method to use with the backend terraform state server     | "LOCK"      |
| BACKEND_UNLOCK_METHOD              | optional                                | unlock method to use with the backend terraform state server   | "UNLOCK"    |
| BACKEND_URL                        | required                                | base url to connect with the backend terraform state server    |             |
| BACKEND_MTLS_CERT                  | optional                                | cert data for mTLS authentication                              |             |
| BACKEND_MTLS_CERT_FILE             | optional                                | certificate file for mTLS authentication                       |             |
| BACKEND_MTLS_KEY                   | optional                                | key data for mTLS authentication                               |             |
| BACKEND_MTLS_KEY_FILE              | optional                                | key file for mTLS authentication                               |             |
| LOG_JSON                           | optional                                | if logging has to use json format                              |             |
| LOG_LEVEL                          | optional                                | active log level one of [TRACE, DEBUG, INFO, WARN, ERROR, OFF] | "INFO"      |
| SERVER_PORT                        | optional                                | port the service is listening to                               | "8080"      |
| TRANSFORM_VAULT_ADDRESS            | optional                                | vault address to de- and encrypt terraform state               |             |
| TRANSFORM_VAULT_APP_ROLE_ID        | optional / required if vault addr != "" | AppRole ID to authenticate with vault                          |             |
| TRANSFORM_VAULT_APP_ROLE_SECRET_ID | optional / required if vault addr != "" | AppRole secret ID to authenticate with vault                   |             |
| TRANSFORM_VAULT_TRANSIT_MOUNT      | optional                                | mount point of the transit engine to use                       | "sops"      |
| TRANSFORM_VAULT_TRANSIT_NAME       | optional                                | name of the transit engine secret to use                       | "terraform" |
