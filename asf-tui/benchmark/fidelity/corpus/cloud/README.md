# Cloud Architecture Benchmark

## Architecture Description

A cloud-native SaaS application on AWS.

## Facts

### Security Controls
- MFA is enabled for all IAM users
- Encryption is enabled (KMS)
- AWS GuardDuty is enabled
- CloudTrail is enabled
- VPC is configured with private subnets
- Security groups are restricted
- IAM policies use least privilege
- AWS WAF is enabled
- S3 buckets are encrypted and versioned
- RDS has automated backups
- CloudWatch monitoring is enabled
- AWS Config is enabled

### Components
- API Gateway (AWS)
- Lambda Functions
- RDS Database
- S3 Storage
- CloudFront CDN
- EC2 Instances
- EKS Cluster
- DynamoDB

### Relationships
- API Gateway -> Lambda (invokes)
- Lambda -> RDS (queries)
- Lambda -> S3 (reads/writes)
- Lambda -> DynamoDB (queries)
- CloudFront -> API Gateway (routes)
- EKS -> RDS (queries)
- EC2 -> S3 (reads)

## Expected Assumptions

1. KMS keys are rotated regularly
2. IAM policies are reviewed periodically
3. Security groups are audited
4. CloudTrail logs are monitored
5. GuardDuty findings are remediated
6. S3 access logs are analyzed
7. RDS backup restore is tested
8. Lambda function permissions are minimal
9. EKS pod security policies are enforced
10. DynamoDB backup is configured
11. CloudFront cache invalidation is secure
12. EC2 instances are patched regularly
13. Identity federation is configured
14. Resource tags are used for cost/security tracking
15. Cross-account access is controlled

## Expected Contradictions

None

## Metrics

- Expected fidelity score: 90%+
- Expected assumption count: 15+
- Expected contradiction count: 0
