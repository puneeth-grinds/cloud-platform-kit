#!/bin/bash

set -euo pipefail

echo "step1: Retrieve ECR Repo URL's"
cd terraform/main

export API_GATEWAY_REPOSITORY_URL=$(terraform output -raw ecr_api_gateway_url)
export VULNERABILITY_SCANNER_REPOSITORY_URL=$(terraform output -raw ecr_vulnerability_scanner_url)