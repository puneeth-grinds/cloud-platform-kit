# Terraform State Backend

Terraform state is stored in a private S3 bucket managed by the configuration in `terraform/state-backend`.

## Backend Configuration

- Region: `us-west-1`
- Bucket name: `cloud-platform-kit-tf-state-<account-id>`
- State key: `state-backend/terraform.tfstate`
- Encryption: AES256
- Versioning: enabled
- Public access: blocked
- State locking: native S3 lock file with `use_lockfile = true`
- DynamoDB locking: not used

The bucket has `prevent_destroy = true` and must not be included in the normal apply-test-destroy workflow.

## Authentication

Follow [AWS Authentication](aws-authentication.md) and use the `cloud-platform-kit-terraform` profile for Terraform commands.

## Normal Workflow

```bash
aws login --profile cloud-platform-kit-login
cd terraform/state-backend
AWS_PROFILE=cloud-platform-kit-terraform terraform init
AWS_PROFILE=cloud-platform-kit-terraform terraform plan
```

The initial local state has already been migrated to S3. Run state migration again only if the backend location changes.

## Verification

Use `terraform output` to confirm Terraform can read the remote state. `terraform state pull` can display the complete state, but its output must be treated as sensitive.

Accessing the state object through its ordinary HTTPS URL should return `AccessDenied` because the bucket is private.

## Safety

- Never make the state bucket or state object public.
- Never commit state files, saved plans, or `.terraform/`.
- Never delete the state object or bucket manually.
- Keep `.terraform.lock.hcl` committed.
