# Explanations

[![readme](../assets/breadcrum-readme.drawio.svg)](../../README.md)

## What is the terraform-sops-backend useful for

When ether a platform engineer is using a terraform HTTP backend provided by a public provider as <https://gitlab.com> he has the dilemma that sensitive data are potentially visible to that public provider.  
[Continue...](./what-is-it-for.md)

## How dose the terraform-sops-backend work

terraform-sops-backend works as a intermediate terraform HTTP backend using a regular terraform HTTP backend but encrypting the terraform state before it is forwarded the the actual backend.  
[Continue...](./how-dose-it-work.md)

## How dose it look from the deployment perspective

The deployment is separated into two environments.  
One is controlled by the user of terraform-sops-backend and the other is controlled by the provider of the used terraform HTTP backend.  
[Continue...](./deployment.md)
