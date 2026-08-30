
# cloud-platform-kit
 
A production-pattern cloud platform engineering project built to demonstrate
hands-on proficiency with the full modern platform engineering stack — infrastructure
as code, containerised microservices, CI/CD automation, distributed tracing, and
security hardening.
 
---
 
## What this is
 
Two Go microservices running on AWS ECS Fargate, fully provisioned with Terraform,
deployed via GitHub Actions, and observed through OpenTelemetry and Grafana Cloud.
The platform does real useful work: it accepts HTTP requests, routes them through an
API gateway, and runs Trivy vulnerability scans against Docker images, persisting
results to Postgres.
 
This is not a tutorial follow-along. Every architectural decision — single-AZ to
control cost, ECS-native deploys over ArgoCD, Grafana Cloud over self-hosted Grafana
— is a deliberate trade-off documented here.
 
---
 
## Architecture
 
```
Internet / user
      │
      ▼
┌─────────────────────────────────────────────────────┐
│  AWS VPC                          Terraform-managed  │
│                                                      │
│  ┌─────────────────────────────────────────────┐    │
│  │  Public subnet                              │    │
│  │  ┌───────────────────────────────────────┐  │    │
│  │  │  ALB  (Application Load Balancer)     │  │    │
│  │  └───────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────┘    │
│                       │                             │
│  ┌─────────────────────────────────────────────┐    │
│  │  Private subnet  ·  ECS Fargate             │    │
│  │  ┌──────────────┐    ┌──────────────────┐   │    │
│  │  │ api-gateway  │───▶│vulnerability-    │   │    │
│  │  │              │    │scanner           │   │    │
│  │  └──────┬───────┘    └────────┬─────────┘   │    │
│  │         │                    │              │    │
│  │  ┌──────▼───────┐    ┌───────▼──────────┐   │    │
│  │  │ RDS Postgres │    │ S3               │   │    │
│  │  └──────────────┘    └──────────────────┘   │    │
│  └─────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────┘
      │                                    │
      │  GitHub Actions                    │ OTel traces
      │  build → push → deploy             ▼
      │  ┌─────────────────┐    ┌──────────────────────┐
      └─▶│ ECR             │    │ Grafana Cloud         │
         │ Container images│    │ Dashboards · Alerts   │
         └─────────────────┘    └──────────────────────┘
```
 
---
 
## Services
 
### `api-gateway`
 
The front door to the platform. Receives all inbound HTTP traffic from the ALB,
validates the caller via API key auth, applies rate limiting, and forwards scan
requests to `vulnerability-scanner`. All distributed traces originate here.
 
**Key responsibilities:**
- API key authentication middleware
- Per-client rate limiting
- Request routing and proxy to downstream services
- OpenTelemetry trace initiation and request logging
### `vulnerability-scanner`
 
The core business logic. Receives a scan request containing a Docker image name,
shells out to Trivy, parses the resulting CVE findings, writes a summary record to
RDS Postgres, uploads the full raw report to S3, and returns a structured JSON
response.
 
**Key responsibilities:**
- Trivy scan execution and output parsing
- CVE result persistence to RDS Postgres via `pgx`
- Raw report storage in S3
- Scan history retrieval API
---
 
## Tech stack
 
| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Infrastructure as Code | Terraform 1.7+ |
| Container runtime | Docker |
| Compute | AWS ECS Fargate |
| Container registry | AWS ECR |
| Network | AWS VPC, ALB |
| Database | AWS RDS Postgres (single-AZ) |
| Object storage | AWS S3 |
| Secrets | AWS SSM Parameter Store |
| CI/CD | GitHub Actions |
| Observability | OpenTelemetry SDK + Grafana Cloud |
| Security scanning | Trivy (in CI and as core feature) |
 
---
 
## Repository structure
 
```
cloud-platform-kit/
├── infra/
│   ├── bootstrap/          # S3 state bucket + DynamoDB lock table (applied once)
│   └── main/               # All project infrastructure (VPC, ECS, RDS, ALB, etc.)
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── networking.tf
│       ├── compute.tf
│       ├── data.tf
│       └── security.tf
├── services/
│   ├── api-gateway/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── proxy/
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── vulnerability-scanner/
│       ├── cmd/
│       │   └── main.go
│       ├── internal/
│       │   ├── handler/
│       │   ├── scanner/
│       │   ├── store/
│       │   └── report/
│       ├── Dockerfile
│       └── go.mod
├── .github/
│   └── workflows/
│       ├── api-gateway.yml
│       └── vulnerability-scanner.yml
└── README.md
```
 
---
 
## Prerequisites
 
- AWS account with programmatic access configured (`aws configure`)
- Go 1.22+
- Terraform 1.7+
- Docker Desktop
- Trivy
---
 
## Getting started
 
### 1. Bootstrap remote state (one-time)
 
```bash
cd infra/bootstrap
terraform init
terraform apply
```
 
Note the output values — you will need the S3 bucket name and DynamoDB table name
for all subsequent Terraform commands.
 
### 2. Provision infrastructure
 
```bash
cd infra/main
terraform init \
  -backend-config="bucket=<state-bucket-name>" \
  -backend-config="key=cloud-platform-kit/terraform.tfstate" \
  -backend-config="region=eu-west-1" \
  -backend-config="dynamodb_table=cloud-platform-kit-tf-locks"
 
terraform plan
terraform apply
```
 
### 3. Build and push images manually (first time)
 
```bash
# Authenticate Docker to ECR
aws ecr get-login-password --region eu-west-1 | \
  docker login --username AWS --password-stdin \
  <account-id>.dkr.ecr.eu-west-1.amazonaws.com
 
# api-gateway
docker build -t cloud-platform-kit/api-gateway ./services/api-gateway
docker tag cloud-platform-kit/api-gateway:latest \
  <account-id>.dkr.ecr.eu-west-1.amazonaws.com/api-gateway:latest
docker push <account-id>.dkr.ecr.eu-west-1.amazonaws.com/api-gateway:latest
 
# vulnerability-scanner
docker build -t cloud-platform-kit/vulnerability-scanner ./services/vulnerability-scanner
docker tag cloud-platform-kit/vulnerability-scanner:latest \
  <account-id>.dkr.ecr.eu-west-1.amazonaws.com/vulnerability-scanner:latest
docker push <account-id>.dkr.ecr.eu-west-1.amazonaws.com/vulnerability-scanner:latest
```
 
### 4. Verify the stack
 
```bash
# Get ALB DNS name
terraform -chdir=infra/main output alb_dns_name
 
# Hit the health endpoint
curl https://<alb-dns-name>/health
```
 
### 5. Tear down when done
 
```bash
cd infra/main
terraform destroy
```
 
> The `apply → test → destroy` workflow keeps AWS costs near zero between
> development sessions. Tear down after every session; provision fresh next time.
 
---
 
## CI/CD
 
On every push to `main`, GitHub Actions runs:
 
1. **Test** — `go test ./...` for the changed service
2. **Build** — multi-stage Docker build producing a minimal production image
3. **Push** — image tagged with the commit SHA pushed to ECR
4. **Deploy** — new ECS task definition registered with the new image digest,
   ECS service updated with a rolling deploy (zero downtime)
GitHub authenticates to AWS using OIDC — no long-lived access keys stored in
GitHub secrets.
 
---
 
## Observability
 
Both services are instrumented with the OpenTelemetry Go SDK. Traces, metrics,
and logs are exported via OTLP to Grafana Cloud (free tier).
 
**Dashboards:**
- Request latency (p50, p95, p99) per service
- Error rate per service and per route
- Scan throughput (scans/minute, CVEs found/scan)
**Alerts:**
- Scanner error rate > 5% for 5 minutes
- API gateway p99 latency > 2s for 3 minutes
---
 
## Design decisions
 
**Single-AZ** — RDS and ECS run in one availability zone to eliminate NAT Gateway
and data transfer costs during development. A production deployment would span
two AZs minimum.
 
**ECS-native CI/CD over ArgoCD** — ArgoCD requires Kubernetes. Rather than running
k3s on a free-tier EC2 to get GitOps, GitHub Actions directly updates ECS task
definitions on deploy. This is a legitimate pattern and simpler to reason about for
a two-service platform.
 
**Grafana Cloud over self-hosted** — self-hosting Grafana + Prometheus on Fargate
costs money and adds operational burden. Grafana Cloud's free tier covers everything
needed here with zero infrastructure overhead.
 
**SSM Parameter Store over Secrets Manager** — SSM free tier covers the handful of
secrets this platform needs. Secrets Manager costs $0.40/secret/month — unnecessary
at this scale.
 
---
 
## Labels for issue tracking
 
```
pillar:terraform
pillar:docker-ecs
pillar:go-services
pillar:cicd
pillar:observability
pillar:security
pillar:docs
 
stage:0-bootstrap
stage:1-network
stage:2-security-groups
stage:3-ecr
stage:4-data-layer
stage:5-compute
stage:6-go-services
stage:7-dockerfiles
stage:8-ecs-wiring
stage:9-smoke-test
stage:10-cicd
stage:11-otel
stage:12-security-hardening
stage:13-docs
 
type:infra
type:code
type:config
type:test
type:docs
```
 
---
 
## Status
 
🚧 Under active development. Building stage by stage