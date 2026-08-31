# AWS Authentication

This project uses browser-based AWS CLI authentication with temporary credentials. Long-Llived AWS access keys are not created or stored

## Authentication Model
Two profiles are made use:

| Profile | Purpose |
|---|---|
| `cloud-platform-kit-login` | Holds the browser-authenticated AWS login session |
| `cloud-platform-kit-personal` | Supplies temporary credentials to Terraform |

## Prerequisites

- AWS CLI version 2.32.0 or newer
- IAM user named `cloud-platform-kit`
- Console access enabled for the IAM user
- MFA enabled
- AWS-managed `SignInLocalDevelopmentAccess` policy attached

