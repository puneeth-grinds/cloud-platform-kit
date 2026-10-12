#!/bin/bash

set -euo pipefail

echo "STEP 1: Retrieve ECR Repo URL's"
cd terraform/main

export API_GATEWAY_REPOSITORY_URL=$(terraform output -raw ecr_api_gateway_url)
export VULNERABILITY_SCANNER_REPOSITORY_URL=$(terraform output -raw ecr_vulnerability_scanner_url)

echo "STEP 2: Set Deployment Values"

cd ../..

export AWS_REGION=us-west-1
export AWS_PROFILE=cloud-platform-kit-login
export IMAGE_TAG=$(git rev-parse --short HEAD)
export ECR_REGISTRY=${API_GATEWAY_REPOSITORY_URL%%/*}