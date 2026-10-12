#!/bin/bash

set -euo pipefail

# STEP 1

echo "STEP 1: Set Deployment Values"


export AWS_REGION=us-west-1
export AWS_PROFILE=cloud-platform-kit-login
export IMAGE_TAG=$(git rev-parse --short HEAD)
export ECR_REGISTRY=${API_GATEWAY_REPOSITORY_URL%%/*}

# STEP 2
echo "STEP 2: Retrieve ECR Repo URL's"
cd terraform/main

export API_GATEWAY_REPOSITORY_URL=$(terraform output -raw ecr_api_gateway_url)
export VULNERABILITY_SCANNER_REPOSITORY_URL=$(terraform output -raw ecr_vulnerability_scanner_url)
export ECR_REGISTRY=${API_GATEWAY_REPOSITORY_URL%%/*}
cd ../..
# STEP 3 

echo "STEP 3: Authenticate Podman to ECR Registry"

aws ecr get-login-password \
  --region "$AWS_REGION" \
  --profile "$AWS_PROFILE" | \
podman login \
  --username AWS \
  --password-stdin \
  "$ECR_REGISTRY"

# STEP 4 
echo "STEP 4: Build and Push API Gateway"

podman build --platform linux/amd64 \
  -t api-gateway:"$IMAGE_TAG" \
  ./services/api-gateway

podman tag api-gateway:"$IMAGE_TAG" \
  "$API_GATEWAY_REPOSITORY_URL:$IMAGE_TAG"

podman push "$API_GATEWAY_REPOSITORY_URL:$IMAGE_TAG"

# STEP 5 

echo "STEP 5: Build and Push Vulnerability Scanner"

podman build --platform linux/amd64 \
  -t vulnerability-scanner:"$IMAGE_TAG" \
  ./services/vulnerability-scanner

podman tag vulnerability-scanner:"$IMAGE_TAG" \
  "$VULNERABILITY_SCANNER_REPOSITORY_URL:$IMAGE_TAG"

podman push "$VULNERABILITY_SCANNER_REPOSITORY_URL:$IMAGE_TAG"

# STEP 6

echo "STEP 6: Verify Both Images"

aws ecr describe-images \
  --region "$AWS_REGION" \
  --profile "$AWS_PROFILE" \
  --repository-name api-gateway \
  --image-ids imageTag="$IMAGE_TAG"

aws ecr describe-images \
  --region "$AWS_REGION" \
  --profile "$AWS_PROFILE" \
  --repository-name vulnerability-scanner \
  --image-ids imageTag="$IMAGE_TAG"

echo "STEP 7: Record these Image URI's"
echo "$API_GATEWAY_REPOSITORY_URL:$IMAGE_TAG"
echo "$VULNERABILITY_SCANNER_REPOSITORY_URL:$IMAGE_TAG"
