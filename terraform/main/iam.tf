data "aws_iam_policy_document" "ecs_tasks_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "github-actions" {
  statement {
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]
    #  RULE 1: The token MUST be explicitly built for AWS.
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    #  RULE 2: The token MUST come from this exact repository AND this exact branch.
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:puneeth-grinds/cloud-platform-kit:ref:refs/heads/main"]
    }
  }
}

resource "aws_iam_openid_connect_provider" "github" {
  url            = "https://token.actions.githubusercontent.com"
  client_id_list = ["sts.amazonaws.com"]

}

resource "aws_iam_role" "ecs_task_execution" {
  name               = "ecs-task-execution"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_role.json
}

resource "aws_iam_role" "github-actions-deploy" {
  name               = "github-actions-deploy"
  assume_role_policy = data.aws_iam_policy_document.github-actions.json
}

resource "aws_iam_policy" "github-actions-deploy" {
  name        = "github-actions-ecs-ecr-least-privilege"
  description = "Least-privilege policy for GitHub Actions to authenticate, push to ECR, deploy to ECS, and pass specific roles."

  policy = jsonencode({
    version = "2012-10-17"
    statement = [
      {
        sid      = "ECRGetAuthorizationToken"
        Effect   = "Allow"
        Action   = ["ecr:GetAuthorizationToken"]
        Resource = "*"
      },
      {
        Sid    = "ECRPushActions"
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:CompleteLayerUpload",
          "ecr:InitiateLayerUpload",
          "ecr:PutImage",
          "ecr:UploadLayerPart"
        ]
        Resource = [
          aws_ecr_repository.ecr_api_gateway.arn,
          aws_ecr_repository.ecr_vulnerability_scanner.arn,
        ]

      },
      {
        Sid    = "ECSDeploymentActions"
        Effect = "Allow"
        Action = [
          "ecs:RegisterTaskDefinition",
          "ecs:DescribeTaskDefinition",
          "ecs:DescribeServices",
          "ecs:UpdateService"
        ]
        Resource = "*"
      },
      {
        Sid    = "IAMPassRoleScoped"
        Effect = "Allow"
        Action = ["iam:PassRole"]
        Resource = [
          aws_iam_role.ecs_task_execution.arn,
          aws_iam_role.api_gateway_task.arn,
          aws_iam_role.vulnerability_scanner_task.arn,
        ]
        Condition = {
          StringEquals = {
            "iam:PassedToService" = "ecs-tasks.amazonaws.com"
          }
        }
      }
    ]

  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "github-actions-deploy" {
  role = aws_iam_role.github-actions-deploy.name
  policy_arn = "aws_iam_policy.github-actions-deploy.arn"
}

data "aws_iam_policy_document" "ecs_task_execution_ssm" {
  statement {
    effect    = "Allow"
    actions   = ["ssm:GetParameters"]
    resources = [aws_ssm_parameter.rds_db_connection_string.arn, aws_ssm_parameter.api_gateway_api_key_ssm.arn]
  }
}

resource "aws_iam_role_policy" "ecs_task_execution_ssm" {
  name   = "read-service-parameters"
  role   = aws_iam_role.ecs_task_execution.id
  policy = data.aws_iam_policy_document.ecs_task_execution_ssm.json
}

resource "aws_iam_role" "api_gateway_task" {
  name               = "api-gateway-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_role.json
}

resource "aws_iam_role" "vulnerability_scanner_task" {
  name               = "vulnerability-scanner-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_role.json
}

data "aws_iam_policy_document" "vulnerability_scanner_s3" {
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
    ]
    resources = ["${aws_s3_bucket.s3_storage.arn}/*"]
  }
}

resource "aws_iam_role_policy" "vulnerability_scanner_s3" {
  name   = "read-write-application-bucket-objects"
  role   = aws_iam_role.vulnerability_scanner_task.id
  policy = data.aws_iam_policy_document.vulnerability_scanner_s3.json
}
