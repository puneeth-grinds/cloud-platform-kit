# AWS Authentication

This project uses browser-based AWS login with temporary credentials. No long-lived AWS access keys are created or stored.

## Authentication Flow

```text
IAM user
  -> browser login
  -> temporary credentials
  -> Terraform role
  -> AWS resources
```

Three local AWS profiles are used:

| Profile | Purpose |
|---|---|
| `cloud-platform-kit-login` | Holds the browser-authenticated login session |
| `cloud-platform-kit-personal` | Makes the temporary login credentials available to other tools |
| `cloud-platform-kit-terraform` | Assumes the Terraform execution role |

The IAM user is named `cloud-platform-kit`. It uses console access, MFA, and the AWS-managed `SignInLocalDevelopmentAccess` policy.

The IAM user can assume the `cloud-platform-kit-terraform` role. Terraform receives its AWS permissions from this role instead of using the IAM user's permissions directly.

## Sign In

```bash
aws login --profile cloud-platform-kit-login
```

Sign in as the `cloud-platform-kit` IAM user and complete MFA verification.

## One-Time Local Configuration

Make the temporary login credentials available to Terraform:

```bash
aws configure set credential_process \
  'aws configure export-credentials --profile cloud-platform-kit-login --format process' \
  --profile cloud-platform-kit-personal
```

```bash
aws configure set region us-west-1 \
  --profile cloud-platform-kit-personal
```

Configure the Terraform role profile:

```bash
aws configure set role_arn \
  arn:aws:iam::512297269123:role/cloud-platform-kit-terraform \
  --profile cloud-platform-kit-terraform
```

```bash
aws configure set source_profile \
  cloud-platform-kit-personal \
  --profile cloud-platform-kit-terraform
```

```bash
aws configure set region us-west-1 \
  --profile cloud-platform-kit-terraform
```

These commands configure local profiles only. They do not store credentials in the repository.

## Verify the Terraform Identity

```bash
aws sts get-caller-identity \
  --profile cloud-platform-kit-terraform
```

Confirm that the ARN contains:

```text
assumed-role/cloud-platform-kit-terraform
```

## Run Terraform

```bash
AWS_PROFILE=cloud-platform-kit-terraform terraform plan
```

Selecting the profile for each command prevents this project from changing the AWS identity used by other repositories.

## Sign Out

```bash
aws logout --profile cloud-platform-kit-login
```

The AWS CLI caches temporary credentials locally and they expire automatically. No credentials are stored in this repository.
