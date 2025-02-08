# Configuration file

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![reference](../assets/breadcrum-reference.drawio.svg)](./index.md)

```yaml
---
server:
  port: "8080"            # (optional) port the service is listening to
backend:
  url: ""                 # (required) base url to connect with the backend terraform state server
  lock_method: "LOCK"     # (optional) lock method to use with the backend terraform state server
  unlock_method: "UNLOCK" # (optional) unlock method to use with the backend terraform state server
transform:
  age:
    public_key: ""        # (required) public AGE key to encrypt terraform state
    private_key: ""       # (optional) private AGE key to decrypt terraform state
  vault:
    address: ""           # (optional) vault address to de- and encrypt terraform state
    app_role:
      id: ""              # (optional) (required if --vault-addr != "") AppRole ID to authenticate with vault
      secret_id: ""       # (optional) (required if --vault-addr != "") AppRole secret ID to authenticate with vault
    transit:
      mount: "sops"       # (optional) mount point of the transit engine to use
      name: "terraform"   # (optional) name of the transit engine secret to use
log:
  json: false             # (optional) if logging has to use json format
  level: "INFO"           # (optional) active log level one of [TRACE, DEBUG, INFO, WARN, ERROR, OFF]
```
