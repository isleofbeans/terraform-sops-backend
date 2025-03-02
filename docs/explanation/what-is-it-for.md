# What is the terraform-sops-backend useful for

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)[![explanation](../assets/breadcrum-explanation.drawio.svg)](./index.md)

Whenever a platform engineer is using a terraform HTTP backend provided by a public provider as <https://gitlab.com> [the dilemma occur](https://github.com/isleofbeans/terraform-sops-backend) that sensitive data are potentially visible to that public provider.  
terraform-sops-backend is addressing this by using [SOPS](https://github.com/getsops/sops) to encrypt the fields in the terraform state with [AGE](https://github.com/FiloSottile/age) and /or [Vault](https://developer.hashicorp.com/vault)

## Unencrypted terraform state

```json
{
  "version": 4,
  "terraform_version": "1.10.4",
  "serial": 34,
  "lineage": "45cb6745-6910-e348-8a9b-499a6d516669",
  "outputs": {
    "password": {
      "value": "zu::i=He0-qOvRsFG:G1F",
      "type": "string",
      "sensitive": true
    }
  },
  "resources": [
    {
      "mode": "managed",
      "type": "random_password",
      "name": "this",
      "provider": "provider[\"registry.terraform.io/hashicorp/random\"]",
      "instances": [
        {
          "schema_version": 3,
          "attributes": {
            "bcrypt_hash": "$2a$10$uEpkHlu3ZMekF1H5LDbYS.mSbbokonskxDav1qWCaXc2OnjazdZ6y",
            "id": "none",
            "keepers": null,
            "length": 21,
            "lower": true,
            "min_lower": 0,
            "min_numeric": 0,
            "min_special": 0,
            "min_upper": 0,
            "number": true,
            "numeric": true,
            "override_special": null,
            "result": "zu::i=He0-qOvRsFG:G1F",
            "special": true,
            "upper": true
          },
          "sensitive_attributes": [
            [
              {
                "type": "get_attr",
                "value": "bcrypt_hash"
              }
            ],
            [
              {
                "type": "get_attr",
                "value": "result"
              }
            ]
          ]
        }
      ]
    }
  ],
  "check_results": null
}
```

## Encrypted terraform state

```json
{
  "version": 4,
  "terraform_version": "1.10.4",
  "serial": 34,
  "lineage": "45cb6745-6910-e348-8a9b-499a6d516669",
  "outputs": {
    "password": {
      "value": "ENC[AES256_GCM,data:rGbUi2Soeigi/wlProelX8jaFC/k,iv:G0GhwbVHI/NgPs16CySzo692SHByDQfAb1PIV+vscb8=,tag:jo60mbQf7WQFU+vhznLjpA==,type:str]",
      "type": "ENC[AES256_GCM,data:/wVnH0OM,iv:7DKsWAlnvT4UUmPvknFGOAF8MlgUGZO8RlXW511xjZo=,tag:+UVdIcDn3hbkVIBffjSMhg==,type:str]",
      "sensitive": "ENC[AES256_GCM,data:xa1NiQ==,iv:QCGz1RUaYHmd/Hbfq0uvC+glNw6p7R1XH0I8iPEjGQ8=,tag:y3NTQt4NV/l3hyMtuQHTCA==,type:bool]"
    }
  },
  "resources": [
    {
      "mode": "ENC[AES256_GCM,data:E2X4K1nOsA==,iv:hOP9U866V9bZSeuVhOFbJZOJB3nQNmJQXiT8asIueTs=,tag:v3ahz0U5JmOXcswVS80rhg==,type:str]",
      "type": "ENC[AES256_GCM,data:5TaorTJkO85ZfHEX6ck1,iv:o5lPnyrVGdllvjTHBPkwyU0K/efvwADXsr9cG8tSz7I=,tag:5OUh3c6Mz+mAbV55VD142w==,type:str]",
      "name": "ENC[AES256_GCM,data:XF+0YA==,iv:YErjDWY4qY9bBecJJ8kSOW8O3/2KQQrWzpC1EL6rjWE=,tag:2d9ODjGt+Sz577khJbCYiw==,type:str]",
      "provider": "ENC[AES256_GCM,data:2Rf30hzwxQOfFOcmtf8V+HPfwmeUoZVbfG/VwTwvdqv2sVIepMiWBuHJ+vcOXkwsuGM=,iv:9sqtUBCVeaE0xE4tLxkJFtpGHxcAoHPnoaj+hbgUrKI=,tag:xrGaO8VkXIGf/1I525g4dw==,type:str]",
      "instances": [
        {
          "schema_version": "ENC[AES256_GCM,data:yg==,iv:Nah5xV2X3cZ8m0FxvyJ9w87sn47eSGFqMPYkNCnEK9U=,tag:eiF7bqngivpSiEpqfAO7ow==,type:float]",
          "attributes": {
            "bcrypt_hash": "ENC[AES256_GCM,data:AxCtCMIi6IE/5HBxXTypETs08POg8ohzTHRVWg0pHnswOaTK9CaqP9DDBAcKjSdl7Ly3W/wH5Lyr93u5,iv:ejO7dgM0NObK80szMGvduDkX6UH88GVF74lxUUEp9po=,tag:To72g59K1m0MtRvXM4RWMw==,type:str]",
            "id": "ENC[AES256_GCM,data:EP31xg==,iv:xjZKHDa2ZHVIpA+CB+7mowgB4YiAh+uY0EZVJdK+XlU=,tag:94Zwk6UCXhepi1aQsrMPVw==,type:str]",
            "keepers": null,
            "length": "ENC[AES256_GCM,data:fis=,iv:+obs6aSnzLrsbqrDSUYSk1qEK0zPv1gKkJ8xBVy0rEw=,tag:zY0KTlf+2aLk9jaoB12OuQ==,type:float]",
            "lower": "ENC[AES256_GCM,data:nwsqtw==,iv:6PNlDqySiy/1+viK/iKA9IN4DL3aTxvFE3NK1rtD98c=,tag:/JNt4vGEh0KWCfP/0IKZQQ==,type:bool]",
            "min_lower": "ENC[AES256_GCM,data:SA==,iv:JBMqVbPP+MxQMphV7NWbPr3CUQOldY6445eNahcBzDE=,tag:MqIGvfLbLnP83uH8EXXesg==,type:float]",
            "min_numeric": "ENC[AES256_GCM,data:AQ==,iv:09oxWeUgNnKblhdRE0q5AedqxN7goUqVwPM6Z+wT0fk=,tag:t1PG/6F9P2uRo2qujyfHpA==,type:float]",
            "min_special": "ENC[AES256_GCM,data:Sg==,iv:bmu95PxOfrgs63G3aE8LEJ1A0xUQkhzjF+dv37Y6Ji8=,tag:XjFkQb/VNChCZLLWFMuj+g==,type:float]",
            "min_upper": "ENC[AES256_GCM,data:NA==,iv:wol7jC4AjjuxGVeg1+IYbtvbSOYDwjZeDlXr2kMvMJ4=,tag:dz1oJ9iobme9IdhZlh0OSQ==,type:float]",
            "number": "ENC[AES256_GCM,data:s6JZTg==,iv:BP8c3HZ7SEkWwLzi3LEnnEDOKlaCX4U8J6Sthj0G9DA=,tag:Y/mxGFvn9MkCPU8Hw1MJcg==,type:bool]",
            "numeric": "ENC[AES256_GCM,data:PK6ihQ==,iv:Fr6Pxur1enKxd1A7eDAnsGh/ryLjkY4qT3iNV4JNx2k=,tag:5C2Wzaa46Kf3NM2QsUOUTQ==,type:bool]",
            "override_special": null,
            "result": "ENC[AES256_GCM,data:GE1n9GAbNGyxsmZ3Kp6HAdElJ295,iv:STsFtgvBWPUketfLVZsAI1VEzjJDtNt+QCic+yelOso=,tag:a6jvko3C11k2vqsFAMUV8g==,type:str]",
            "special": "ENC[AES256_GCM,data:T44sPw==,iv:uKQTv0p1tEVtvs1CUBGhXCH1b1m74TTDunpZWUnKsdU=,tag:HKIrzr+yVLKQExN3Jxgj/A==,type:bool]",
            "upper": "ENC[AES256_GCM,data:/HQapw==,iv:XHUngHPaYcud7HIXXuov7apZsbKpRBb12hlKechs6ik=,tag:oyUhMTpn4O9H2bhRZrrdJg==,type:bool]"
          },
          "sensitive_attributes": [
            [
              {
                "type": "ENC[AES256_GCM,data:x2ZR5yb8K5I=,iv:tA68qfLDHoxE2y3DK3nLYmHCJ823QCp2/DkBkVOIjeg=,tag:pYjI0YFWmYpzR9LmF8LbHA==,type:str]",
                "value": "ENC[AES256_GCM,data:PF84BM2cedOqcsU=,iv:HgCbUJW4rbp1I2dMR5qJb/wpkdjNNSgyZqs+suhjzQ8=,tag:3zfRGc2Ut1FszqoJUInq+w==,type:str]"
              }
            ],
            [
              {
                "type": "ENC[AES256_GCM,data:H68x/YTmd6w=,iv:Xbv1IY+mVroSxnCDgoIWdXfb2uHbYejtFHEb3bLEhFs=,tag:YXsKtoY5AF35HhajV/0Mcg==,type:str]",
                "value": "ENC[AES256_GCM,data:9ZRkFlD9,iv:IOSdn/tiAl0jMGQ13c9e2565LZDiDRGrQyaebC0mB84=,tag:Ptt630qU5umg12dK3dXLfw==,type:str]"
              }
            ]
          ]
        }
      ]
    }
  ],
  "check_results": null,
  "sops": {
    "kms": null,
    "gcp_kms": null,
    "azure_kv": null,
    "hc_vault": [
      {
        "vault_address": "http://127.0.0.1:8200",
        "engine_path": "sops",
        "key_name": "terraform",
        "created_at": "2025-02-09T12:00:11Z",
        "enc": "vault:v1:wV6Y3OxeS1OyhljIpPI41sGSgsKicfRDyjCQjxxYhqcji1wyP+NgfmF7YqCMgeAMbl6DIpeDk/eYcydd"
      }
    ],
    "age": [
      {
        "recipient": "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
        "enc": "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSA2L3U5c0U4UzdHcmExVXhn\nQ2I5VGdpQVBObkVRQ3ZESkJ5eEJPSFJ2T3pRCktSUk50bkF5d29ZMEFXaGRrS3dW\nWXM1MkJSS2Z3V09LQ2Iyb1M4VkdiVUUKLS0tIE1ST2EwelBZc01LVDR6MzkvOVgr\nTTVhTTMrN0hrV25DTERJMjJjdldpcWsK/OSP4VPRLRGzeKNC5R1Q4fl5sc3Wwxa1\nNw8opp4vnfXtfgQRwl8enPI8ZsP/qfKkj1BWTWhCuha6dNlT+DYLCg==\n-----END AGE ENCRYPTED FILE-----\n"
      }
    ],
    "lastmodified": "2025-02-09T12:00:11Z",
    "mac": "ENC[AES256_GCM,data:jg6uYvGsHAxcYaFHy+ck7Aq0kEtcjnOsQf4pvs2yJjt3Dd4ZdO2DdnAruftMFInJ8/rurwXloJHdP9BIbEuZc4d4EVIu7xIiLJeYYucpQ6Pb+WOchW//i0Kfbr+ffhZDvz1ZePWym+vBLHluUG87GltuuwCJrmQ7mKYaehrNUSk=,iv:FZdf08Y6Tu6JI0b233X4BNIhkeVrh4JC9v+4qLEb1mc=,tag:UxKv89s/tvLWhMgqSmq1lg==,type:str]",
    "pgp": null,
    "unencrypted_regex": "^(version|terraform_version|serial|lineage)$",
    "version": "3.9.3"
  }
}
```
