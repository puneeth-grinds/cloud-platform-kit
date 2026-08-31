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

Check the aws CLI version:
```
aws --version
```
> The `SignInLocalDevelopmentAccess` policy only permits browser-based authentication. It does not grant permission to create or manage AWS resources.

## Sign in
Start the browser authentication flow:
```bash
aws login --profile cloud-platform-kit-login
```
When prompted:
1. Choose `us-west-1` as the region.
2. Sign in to the correct personal AWS account.
3. Use the `cloud-platform-kit` IAM username and its console password.
4. Complete MFA verification.
5. Confirm that the selected session belongs to the expected account and IAM user.

## Configure Terraform access

This only needs to be configured once:

```bash
aws configure set credential_process \
  'aws configure export-credentials --profile cloud-platform-kit-login --format process' \
  --profile cloud-platform-kit-personal
```

```bash
aws configure set region us-west-1 \
  --profile cloud-platform-kit-personal
```

## Verify the identity

```bash
aws sts get-caller-identity \
  --profile cloud-platform-kit-personal
```

Confirm that the expected personal AWS account and the `cloud-platform-kit` IAM user are shown.

## Run Terraform

```bash
AWS_PROFILE=cloud-platform-kit-personal terraform plan
```

Using `AWS_PROFILE` for each command prevents this project from affecting the AWS profiles used by other repositories.

## Sign out

```bash
aws logout --profile cloud-platform-kit-login
```

Temporary credentials are cached locally by the AWS CLI and expire automatically. No credentials are stored in this repository.