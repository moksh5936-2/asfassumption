#!/usr/bin/env python3
"""
Generate assumption_gold_standard.csv for Phase 6 Human Validation.
Maps ASF assumption patterns to 20 reference architectures with 2-4 assumptions per pattern.
"""

import csv
import os

PATTERN_IDS = {
    "Authentication (MFA)": "AUTH-MFA",
    "Authentication (SSO)": "AUTH-SSO",
    "Availability & Resilience": "AVAIL-RESIL",
    "Backup & Recovery": "BACKUP-RECOV",
    "Change Management": "CHG-MGMT",
    "Cloud Security (IAM)": "CLOUD-IAM",
    "Container Security": "CONT-SEC",
    "Data Flow & Classification": "DATAFLOW",
    "Encryption at Rest": "ENCR-REST",
    "Encryption in Transit (TLS)": "ENCR-TLS",
    "Endpoint Security": "ENDPOINT",
    "Human Factors & Process": "HUMAN-FACT",
    "Identity Lifecycle": "ID-LIFECYCLE",
    "Incident Response": "INC-RESP",
    "Least Privilege": "LEAST-PRIV",
    "Monitoring & Alerting": "MONITOR",
    "Network Segmentation": "NET-SEG",
    "Physical Security": "PHYS-SEC",
    "Supply Chain Security": "SUPPLY-CHAIN",
    "Third-party Dependency": "THIRD-PARTY",
}

ARCHITECTURES = {
    1: {
        "name": "VPN \u2192 Payroll DB",
        "topology": "User Laptop --VPN-- VPN Gateway --TLS-- Internal Web App --SQL-- Payroll Database (RDS)",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Endpoint Security",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Third-party Dependency"
        ],
        "components": {
            "user": "User Laptop",
            "gateway": "VPN Gateway",
            "app": "Internal Web App",
            "database": "Payroll Database (RDS)",
            "identity": "AD Directory",
            "network": "VPN Tunnel",
            "auth": "AD Authentication",
        }
    },
    2: {
        "name": "Web App \u2192 LB \u2192 App Server \u2192 RDS",
        "topology": "Browser --HTTPS-- ALB --HTTPS-- EC2 App Server (x3) --SQL-- RDS Primary + Replica --Logs-- CloudWatch",
        "patterns": [
            "Authentication (SSO)", "Availability & Resilience", "Backup & Recovery",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Browser",
            "gateway": "ALB",
            "app": "EC2 App Server",
            "database": "RDS Primary + Replica",
            "monitoring": "CloudWatch",
            "network": "VPC Private Subnet",
            "auth": "SSO Provider",
        }
    },
    3: {
        "name": "Mobile \u2192 API Gateway \u2192 Lambda \u2192 DynamoDB",
        "topology": "Mobile App --HTTPS-- API Gateway --Event-- Lambda Function (xN) --SDK-- DynamoDB Table, Cognito User Pool",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Incident Response", "Least Privilege", "Monitoring & Alerting",
            "Network Segmentation", "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Mobile App",
            "gateway": "API Gateway",
            "app": "Lambda Function",
            "database": "DynamoDB Table",
            "identity": "Cognito User Pool",
            "network": "API Gateway VPC Endpoint",
            "auth": "Cognito Auth",
        }
    },
    4: {
        "name": "SSO \u2192 IdP \u2192 SAML Federation",
        "topology": "User Browser --SAML-- Okta IdP --SAML Assertion-- Service Provider Apps (x5), AD Directory",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Cloud Security (IAM)", "Data Flow & Classification", "Encryption at Rest",
            "Encryption in Transit (TLS)", "Human Factors & Process", "Identity Lifecycle",
            "Incident Response", "Least Privilege", "Monitoring & Alerting",
            "Network Segmentation", "Third-party Dependency"
        ],
        "components": {
            "user": "User Browser",
            "gateway": "Okta IdP",
            "app": "Service Provider Apps",
            "database": "AD Directory",
            "identity": "Okta IdP",
            "network": "SAML Federation",
        }
    },
    5: {
        "name": "Microservices \u2192 Mesh \u2192 K8s \u2192 Istio",
        "topology": "Ingress Gateway --mTLS-- Service A --mTLS-- Service B --mTLS-- Service C, Istio Pilot, Citadel CA, K8s API, StatefulSet DB",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Container Security", "Data Flow & Classification", "Encryption at Rest",
            "Encryption in Transit (TLS)", "Human Factors & Process", "Identity Lifecycle",
            "Incident Response", "Least Privilege", "Monitoring & Alerting",
            "Network Segmentation", "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "External Client",
            "gateway": "Ingress Gateway",
            "app": "Service A/B/C",
            "database": "StatefulSet DB",
            "identity": "Citadel CA",
            "network": "Istio Service Mesh",
            "control_plane": "K8s API Server",
            "storage": "Persistent Volume",
        }
    },
    6: {
        "name": "E-commerce \u2192 Payment \u2192 PCI",
        "topology": "Browser --HTTPS-- Web App --API-- Payment Processor (Stripe), Token Vault, Order DB",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Container Security", "Data Flow & Classification", "Encryption at Rest",
            "Encryption in Transit (TLS)", "Human Factors & Process", "Identity Lifecycle",
            "Incident Response", "Least Privilege", "Monitoring & Alerting",
            "Network Segmentation", "Physical Security", "Supply Chain Security",
            "Third-party Dependency"
        ],
        "components": {
            "user": "Browser",
            "gateway": "Web App",
            "app": "Web App",
            "database": "Order DB",
            "identity": "Customer Auth",
            "third_party": "Payment Processor (Stripe)",
            "vault": "Token Vault",
            "network": "PCI Enclave",
        }
    },
    7: {
        "name": "Multi-region \u2192 Active/Passive \u2192 DR",
        "topology": "Route53 --Failover-- Region A (Active) App+DB --Cross-Region-- Region B (Passive) App+DB Replica, S3 Primary/Replica",
        "patterns": [
            "Authentication (SSO)", "Availability & Resilience", "Backup & Recovery",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Third-party Dependency"
        ],
        "components": {
            "app": "Application (Region A/B)",
            "database": "DB Primary/Replica",
            "storage": "S3 Primary/Replica",
            "dns": "Route53",
            "network": "Cross-Region Link",
            "monitoring": "Route53 Health Checks",
        }
    },
    8: {
        "name": "CI/CD \u2192 Artifact \u2192 Deploy",
        "topology": "Developer --Push-- GitHub --CI-- GitHub Actions --Push-- Artifact Registry (ECR) --Sync-- ArgoCD --Deploy-- K8s Cluster",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Change Management", "Cloud Security (IAM)", "Container Security",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Developer",
            "gateway": "GitHub",
            "app": "K8s Cluster",
            "ci": "GitHub Actions",
            "registry": "Artifact Registry (ECR)",
            "cd": "ArgoCD",
            "network": "Deployment Pipeline",
        }
    },
    9: {
        "name": "Vendor SaaS \u2192 API \u2192 Internal",
        "topology": "Vendor SaaS (Salesforce) --OAuth-- API Gateway --OAuth Token-- Internal CRM Sync --SQL-- Customer DB",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "vendor": "Vendor SaaS (Salesforce)",
            "gateway": "API Gateway",
            "app": "Internal CRM Sync",
            "database": "Customer DB",
            "auth": "OAuth Token Service",
            "network": "Vendor API Integration",
        }
    },
    10: {
        "name": "Kafka \u2192 S3 \u2192 Redshift",
        "topology": "Producer Apps --Produce-- Kafka Cluster --Consume-- Spark Job --Write-- S3 Data Lake --Query-- Redshift",
        "patterns": [
            "Authentication (SSO)", "Availability & Resilience", "Backup & Recovery",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Human Factors & Process",
            "Identity Lifecycle", "Incident Response", "Least Privilege",
            "Monitoring & Alerting", "Network Segmentation", "Supply Chain Security",
            "Third-party Dependency"
        ],
        "components": {
            "app": "Producer Apps",
            "gateway": "Kafka Cluster",
            "processor": "Spark Job",
            "database": "Redshift",
            "storage": "S3 Data Lake",
            "network": "Kafka Data Pipeline",
        }
    },
    11: {
        "name": "Healthcare \u2192 PHI \u2192 HIPAA",
        "topology": "Patient Portal --Auth0-- App Server --SQL-- PHI Database, Audit Logs --SIEM--",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Physical Security", "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Patient Portal",
            "app": "App Server",
            "database": "PHI Database",
            "identity": "Auth0",
            "monitoring": "SIEM",
            "logs": "Audit Logs",
            "network": "HIPAA Enclave",
        }
    },
    12: {
        "name": "Fintech \u2192 Ledger \u2192 SOX",
        "topology": "User --Trading App-- Ledger Service --SQL-- Accounting DB, Market Data API, Audit Trail Service --Write-- Immutable Audit DB",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Trader",
            "app": "Trading App",
            "gateway": "Ledger Service",
            "database": "Accounting DB",
            "audit": "Audit Trail Service",
            "audit_db": "Immutable Audit DB",
            "third_party": "Market Data API",
            "network": "SOX Enclave",
        }
    },
    13: {
        "name": "Partner B2B \u2192 Federation \u2192 API",
        "topology": "Partner A IdP --SAML-- API Gateway --OAuth-- Partner A Resources, Partner B IdP --SAML-- API Gateway --OAuth-- Partner B Resources",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Human Factors & Process",
            "Identity Lifecycle", "Incident Response", "Least Privilege",
            "Monitoring & Alerting", "Network Segmentation", "Supply Chain Security",
            "Third-party Dependency"
        ],
        "components": {
            "partner_a": "Partner A IdP",
            "partner_b": "Partner B IdP",
            "gateway": "API Gateway",
            "resource_a": "Partner A Resources",
            "resource_b": "Partner B Resources",
            "auth": "SAML Federation",
            "network": "B2B API Exchange",
        }
    },
    14: {
        "name": "Hybrid Cloud \u2192 VPN \u2192 Direct Connect",
        "topology": "On-Prem DC --Direct Connect-- AWS VPC, Corporate FW, Transit Gateway --Prod/Staging/Dev VPCs--",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Human Factors & Process",
            "Identity Lifecycle", "Incident Response", "Least Privilege",
            "Monitoring & Alerting", "Network Segmentation", "Third-party Dependency"
        ],
        "components": {
            "on_prem": "On-Prem DC",
            "gateway": "Corporate FW",
            "router": "Transit Gateway",
            "vpc_prod": "Prod VPC",
            "vpc_staging": "Staging VPC",
            "vpc_dev": "Dev VPC",
            "direct_connect": "Direct Connect",
            "vpn": "VPN Backup",
            "network": "Hybrid Cloud Network",
        }
    },
    15: {
        "name": "IoT \u2192 MQTT \u2192 Gateway \u2192 Cloud",
        "topology": "IoT Device --MQTT/TLS-- IoT Gateway --Stream-- Kinesis --Processing Lambda-- TimeSeries DB, Device Registry / CA",
        "patterns": [
            "Authentication (MFA)", "Availability & Resilience", "Backup & Recovery",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Human Factors & Process",
            "Identity Lifecycle", "Incident Response", "Least Privilege",
            "Monitoring & Alerting", "Network Segmentation", "Supply Chain Security",
            "Third-party Dependency"
        ],
        "components": {
            "device": "IoT Device",
            "gateway": "IoT Gateway",
            "stream": "Kinesis",
            "app": "Processing Lambda",
            "database": "TimeSeries DB",
            "identity": "Device Registry / CA",
            "network": "MQTT/TLS Tunnel",
        }
    },
    16: {
        "name": "ML \u2192 Training \u2192 Serving",
        "topology": "Feature Store --Training Job (SageMaker)-- Model Registry --Serving Endpoint-- Production App, Data Lake (S3)",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "feature_store": "Feature Store",
            "training": "Training Job (SageMaker)",
            "model_registry": "Model Registry",
            "endpoint": "Serving Endpoint",
            "app": "Production App",
            "data_lake": "Data Lake (S3)",
            "network": "ML Pipeline",
        }
    },
    17: {
        "name": "Multi-tenant \u2192 Tenant Isolation",
        "topology": "Tenant A/B --API Gateway (Shared)-- App Service (Shared)-- Tenant A/B DB (Isolated), JWT Claims",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Container Security", "Data Flow & Classification", "Encryption at Rest",
            "Encryption in Transit (TLS)", "Human Factors & Process", "Identity Lifecycle",
            "Incident Response", "Least Privilege", "Monitoring & Alerting",
            "Network Segmentation", "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "tenant_a": "Tenant A",
            "tenant_b": "Tenant B",
            "gateway": "API Gateway (Shared)",
            "app": "App Service (Shared)",
            "db_a": "Tenant A DB",
            "db_b": "Tenant B DB",
            "auth": "JWT Claims",
            "network": "Tenant Isolation Boundary",
        }
    },
    18: {
        "name": "CDN \u2192 WAF \u2192 Origin \u2192 DB",
        "topology": "Global Users --CloudFront CDN-- WAF --ALB-- EC2 Origin --RDS--, Lambda@Edge (Auth)",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Global Users",
            "cdn": "CloudFront CDN",
            "waf": "WAF",
            "gateway": "ALB",
            "app": "EC2 Origin",
            "database": "RDS",
            "edge": "Lambda@Edge",
            "network": "CDN Edge Network",
        }
    },
    19: {
        "name": "Secrets \u2192 Vault \u2192 App \u2192 Rotation",
        "topology": "App Pod --mTLS-- Vault Agent Sidecar --API-- Vault Server --KMS (Unseal)--, Database (stores secrets)",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Change Management", "Cloud Security (IAM)", "Data Flow & Classification",
            "Encryption at Rest", "Encryption in Transit (TLS)", "Human Factors & Process",
            "Identity Lifecycle", "Incident Response", "Least Privilege",
            "Monitoring & Alerting", "Network Segmentation", "Supply Chain Security",
            "Third-party Dependency"
        ],
        "components": {
            "app": "App Pod",
            "gateway": "Vault Agent Sidecar",
            "vault": "Vault Server",
            "kms": "KMS (Unseal)",
            "database": "Secrets Storage Database",
            "identity": "Vault Identity",
            "network": "Vault mTLS Channel",
        }
    },
    20: {
        "name": "ERP \u2192 SOX \u2192 Financial Reporting",
        "topology": "Finance Team --ERP Web App-- ERP Backend --SQL-- Financial DB, Approval Workflow, Audit Logs, Reporting Engine --Auditor Access--",
        "patterns": [
            "Authentication (MFA)", "Authentication (SSO)", "Availability & Resilience",
            "Backup & Recovery", "Change Management", "Cloud Security (IAM)",
            "Data Flow & Classification", "Encryption at Rest", "Encryption in Transit (TLS)",
            "Human Factors & Process", "Identity Lifecycle", "Incident Response",
            "Least Privilege", "Monitoring & Alerting", "Network Segmentation",
            "Physical Security", "Supply Chain Security", "Third-party Dependency"
        ],
        "components": {
            "user": "Finance Team",
            "app": "ERP Web App",
            "backend": "ERP Backend",
            "database": "Financial DB",
            "workflow": "Approval Workflow",
            "logs": "Audit Logs",
            "reporting": "Reporting Engine",
            "auditor": "Auditor Access",
            "network": "SOX Enclave",
        }
    },
}

PATTERN_TEMPLATES = {
    "Authentication (MFA)": [
        {
            "component_key": "gateway",
            "depends_on": "MFA Provider",
            "text": "MUST be available for {component} to enforce step-up authentication.",
            "category": "Dependency",
            "risk": "Authentication bypass during MFA provider outage",
            "verify": "Monitor MFA provider uptime SLAs and test failover",
        },
        {
            "component_key": "gateway",
            "depends_on": "MFA Enrollment",
            "text": "MUST be enforced for all users accessing {component} without exception.",
            "category": "Explicit",
            "risk": "Non-enforcement creates MFA-free access paths",
            "verify": "Audit MFA enforcement policy across all access points monthly",
        },
        {
            "component_key": "gateway",
            "depends_on": "MFA Token Distribution",
            "text": "trusts that MFA tokens are provisioned before users need access to {component}.",
            "category": "Trust",
            "risk": "Users cannot authenticate on day one without provisioned token",
            "verify": "Automate MFA token provisioning in onboarding workflow",
        },
        {
            "component_key": "user",
            "depends_on": "Phishing Resistance",
            "text": "assumes that MFA factor is not interceptable via phishing for {component}.",
            "category": "Implicit",
            "risk": "Phishing sites can proxy MFA session in real time",
            "verify": "Deploy phishing-resistant MFA (WebAuthn/FIDO2) and test",
        },
    ],
    "Authentication (SSO)": [
        {
            "component_key": "identity",
            "depends_on": "SSO Provider",
            "text": "MUST be enforced for all application access through {component} without exception.",
            "category": "Explicit",
            "risk": "Non-enforcement creates credential-free access paths",
            "verify": "Audit all applications for SSO integration coverage",
        },
        {
            "component_key": "gateway",
            "depends_on": "SSO Provider Availability",
            "text": "MUST be available to authenticate users at {component}.",
            "category": "Dependency",
            "risk": "Complete access denial during SSO provider outage",
            "verify": "Test SSO failover mechanism and monitor provider uptime",
        },
        {
            "component_key": "identity",
            "depends_on": "Session Management",
            "text": "assumes that SSO sessions are properly scoped and timed at {component}.",
            "category": "Implicit",
            "risk": "Long-lived sessions weaken access control after credential change",
            "verify": "Review SSO session timeout configuration and test revocation",
        },
        {
            "component_key": "gateway",
            "depends_on": "Federation Trust",
            "text": "trusts that the SSO identity provider assertions are cryptographically verified at {component}.",
            "category": "Trust",
            "risk": "Forged SAML assertions grant unauthorized access",
            "verify": "Validate signature verification on all SAML assertions",
        },
    ],
    "Availability & Resilience": [
        {
            "component_key": "app",
            "depends_on": "Auto-scaling",
            "text": "MUST be available for {component} to handle peak load without degradation.",
            "category": "Dependency",
            "risk": "Service degradation or outage during traffic spikes",
            "verify": "Load test to validate auto-scaling thresholds and capacity",
        },
        {
            "component_key": "database",
            "depends_on": "Database Replication",
            "text": "architecturally requires replica failover for {component} to maintain availability.",
            "category": "Architectural",
            "risk": "Primary failure causes extended write outage",
            "verify": "Perform failover drill and measure cutover time",
        },
        {
            "component_key": "gateway",
            "depends_on": "Redundant Infrastructure",
            "text": "assumes that {component} is deployed across multiple availability zones.",
            "category": "Implicit",
            "risk": "Single-AZ failure takes down the entire service",
            "verify": "Audit infrastructure deployment for multi-AZ distribution",
        },
        {
            "component_key": "network",
            "depends_on": "Network Path Redundancy",
            "text": "MUST be tested for failover scenarios involving {component}.",
            "category": "Derived",
            "risk": "Undiscovered network path failures cause silent downtime",
            "verify": "Execute chaos engineering experiments on network paths",
        },
    ],
    "Backup & Recovery": [
        {
            "component_key": "database",
            "depends_on": "Backup Schedule",
            "text": "MUST be tested via restore drill for {component} at least quarterly.",
            "category": "Derived",
            "risk": "Untested backups are equivalent to no backups",
            "verify": "Perform quarterly restore test and measure RTO/RPO",
        },
        {
            "component_key": "database",
            "depends_on": "Backup Storage",
            "text": "MUST be available for {component} in a separate geographic region.",
            "category": "Dependency",
            "risk": "Regional disaster destroys both primary and backup data",
            "verify": "Validate cross-region backup replication status",
        },
        {
            "component_key": "storage",
            "depends_on": "Data Integrity",
            "text": "assumes that {component} backup data is not silently corrupted.",
            "category": "Implicit",
            "risk": "Silent corruption restores unusable data after disaster",
            "verify": "Implement checksum verification on all backup files",
        },
    ],
    "Change Management": [
        {
            "component_key": "app",
            "depends_on": "Change Approval",
            "text": "MUST be documented and approved before deployment to {component}.",
            "category": "Operational",
            "risk": "Untracked changes introduce misconfiguration vulnerabilities",
            "verify": "Audit change tickets against actual deployment history",
        },
        {
            "component_key": "gateway",
            "depends_on": "Change Testing",
            "text": "MUST be tested in a staging environment before promoting to {component}.",
            "category": "Derived",
            "risk": "Untested changes cause production outages",
            "verify": "Verify staging deployment gate in CI/CD pipeline",
        },
        {
            "component_key": "network",
            "depends_on": "Network Change Review",
            "text": "assumes that changes to {component} do not bypass security controls.",
            "category": "Implicit",
            "risk": "Security group or firewall changes create unauthorized access paths",
            "verify": "Automated drift detection on network security configuration",
        },
    ],
    "Cloud Security (IAM)": [
        {
            "component_key": "app",
            "depends_on": "IAM Role Configuration",
            "text": "MUST be enforced with least-privilege IAM roles for {component}.",
            "category": "Explicit",
            "risk": "Over-privileged roles allow lateral movement",
            "verify": "Review IAM policy boundaries and use Access Analyzer",
        },
        {
            "component_key": "identity",
            "depends_on": "IAM Policy Review",
            "text": "assumes that IAM policies attached to {component} do not grant unintended access.",
            "category": "Implicit",
            "risk": "IAM policy complexity hides privilege escalation paths",
            "verify": "Run IAM policy simulation and least-privilege analysis",
        },
        {
            "component_key": "gateway",
            "depends_on": "IAM Federation",
            "text": "trusts that IAM roles assumed by {component} are scoped to trusted entities.",
            "category": "Trust",
            "risk": "Cross-account role trust allows external compromise",
            "verify": "Audit IAM trust policy documents for external principals",
        },
        {
            "component_key": "app",
            "depends_on": "Service Control Policies",
            "text": "MUST be available to prevent actions on {component} outside of approved regions.",
            "category": "Dependency",
            "risk": "Resource creation in unapproved regions bypasses compliance",
            "verify": "Validate SCP boundary enforcement across all accounts",
        },
    ],
    "Container Security": [
        {
            "component_key": "app",
            "depends_on": "Container Image Scanning",
            "text": "assumes that container images deployed to {component} have no known vulnerabilities.",
            "category": "Implicit",
            "risk": "Known CVEs in base images are deployed to production",
            "verify": "Integrate image scanning into CI/CD pipeline with fail-on-critical",
        },
        {
            "component_key": "app",
            "depends_on": "Pod Security Context",
            "text": "MUST be enforced to prevent privileged containers in {component}.",
            "category": "Explicit",
            "risk": "Privileged containers enable host-level compromise",
            "verify": "Audit pod security admission controller configuration",
        },
        {
            "component_key": "app",
            "depends_on": "Runtime Security",
            "text": "MUST be validated for runtime behavior anomalies in {component}.",
            "category": "Derived",
            "risk": "Compromised container runs undetected cryptominers or data exfiltration",
            "verify": "Deploy runtime security agent and monitor for drift",
        },
    ],
    "Data Flow & Classification": [
        {
            "component_key": "database",
            "depends_on": "Data Classification",
            "text": "assumes that data flowing to {component} is properly classified by sensitivity.",
            "category": "Implicit",
            "risk": "Unclassified sensitive data lacks required protection controls",
            "verify": "Audit data classification tags and data loss prevention rules",
        },
        {
            "component_key": "network",
            "depends_on": "Data Flow Mapping",
            "text": "MUST be documented for all data paths through {component}.",
            "category": "Operational",
            "risk": "Undocumented data flows create blind spots for security monitoring",
            "verify": "Maintain and review data flow diagrams quarterly",
        },
        {
            "component_key": "app",
            "depends_on": "Data Leak Prevention",
            "text": "MUST be enforced to prevent unauthorized data exfiltration from {component}.",
            "category": "Explicit",
            "risk": "Sensitive data exfiltration via authorized channels goes undetected",
            "verify": "Deploy DLP controls and test with realistic data patterns",
        },
        {
            "component_key": "gateway",
            "depends_on": "Protocol Validation",
            "text": "assumes that {component} validates all data schemas before processing.",
            "category": "Implicit",
            "risk": "Malformed data payloads trigger processing pipeline failures",
            "verify": "Implement schema validation gateway and monitor rejections",
        },
    ],
    "Encryption at Rest": [
        {
            "component_key": "database",
            "depends_on": "Encryption Key",
            "text": "MUST be enforced for all data stored in {component} using AES-256.",
            "category": "Explicit",
            "risk": "Unencrypted data at rest exposed via physical storage access",
            "verify": "Audit encryption status of all storage volumes and snapshots",
        },
        {
            "component_key": "database",
            "depends_on": "Key Management",
            "text": "assumes that encryption keys for {component} are rotated on schedule.",
            "category": "Implicit",
            "risk": "Stale encryption keys increase blast radius of key compromise",
            "verify": "Automate key rotation and audit last rotation date",
        },
        {
            "component_key": "storage",
            "depends_on": "Key Access Control",
            "text": "MUST be available via KMS with strict access controls for {component}.",
            "category": "Dependency",
            "risk": "KMS unavailability prevents decryption of critical data",
            "verify": "Test KMS failover and validate key access audit logs",
        },
    ],
    "Encryption in Transit (TLS)": [
        {
            "component_key": "gateway",
            "depends_on": "TLS Certificate",
            "text": "MUST be enforced with TLS 1.2+ for all connections to {component}.",
            "category": "Explicit",
            "risk": "Weak or absent TLS allows man-in-the-middle attacks",
            "verify": "Scan endpoints with TLS scanner and verify minimum version",
        },
        {
            "component_key": "gateway",
            "depends_on": "Certificate Validity",
            "text": "assumes that TLS certificates for {component} are renewed before expiry.",
            "category": "Implicit",
            "risk": "Expired certificates cause service interruption",
            "verify": "Automate certificate renewal with 30-day expiry monitoring",
        },
        {
            "component_key": "network",
            "depends_on": "Certificate Pinning",
            "text": "trusts that the CA chain for {component} has not been compromised.",
            "category": "Trust",
            "risk": "Compromised CA issues fraudulent certificates for traffic interception",
            "verify": "Implement certificate pinning for critical internal services",
        },
        {
            "component_key": "app",
            "depends_on": "mTLS Configuration",
            "text": "MUST be tested to verify mutual TLS is correctly configured for {component}.",
            "category": "Derived",
            "risk": "One-way TLS only authenticates server, not client",
            "verify": "Verify mTLS handshake with invalid client cert is rejected",
        },
    ],
    "Endpoint Security": [
        {
            "component_key": "user",
            "depends_on": "Endpoint Protection",
            "text": "assumes that {component} has up-to-date endpoint detection and response installed.",
            "category": "Implicit",
            "risk": "Compromised endpoint pivots to internal application",
            "verify": "Audit EDR agent installation and last check-in across all endpoints",
        },
        {
            "component_key": "user",
            "depends_on": "Patch Management",
            "text": "MUST be enforced for OS and browser patches on {component}.",
            "category": "Explicit",
            "risk": "Unpatched endpoint vulnerabilities enable remote compromise",
            "verify": "Automate patch deployment and audit patch compliance",
        },
    ],
    "Human Factors & Process": [
        {
            "component_key": "user",
            "depends_on": "Security Training",
            "text": "assumes that personnel accessing {component} have completed security awareness training.",
            "category": "Implicit",
            "risk": "Untrained users fall victim to phishing or social engineering",
            "verify": "Mandate annual security training with phishing simulation",
        },
        {
            "component_key": "app",
            "depends_on": "Procedure Documentation",
            "text": "MUST be documented for incident response procedures involving {component}.",
            "category": "Operational",
            "risk": "Undocumented procedures delay incident response",
            "verify": "Review and tabletop-test incident response runbooks quarterly",
        },
        {
            "component_key": "user",
            "depends_on": "User Behavior Monitoring",
            "text": "assumes that users of {component} do not share credentials or bypass security controls.",
            "category": "Implicit",
            "risk": "Credential sharing undermines non-repudiation and access control",
            "verify": "Monitor for anomalous access patterns and shared account usage",
        },
    ],
    "Identity Lifecycle": [
        {
            "component_key": "identity",
            "depends_on": "User Provisioning",
            "text": "MUST be enforced to provision users in {component} within 24 hours of hire.",
            "category": "Explicit",
            "risk": "Delayed provisioning blocks productivity or creates unauthorized shadow access",
            "verify": "Audit time-to-provision metrics and automate identity lifecycle",
        },
        {
            "component_key": "identity",
            "depends_on": "User Deprovisioning",
            "text": "MUST be enforced to deactivate access to {component} within 1 hour of termination.",
            "category": "Explicit",
            "risk": "Former employees retain access to sensitive systems",
            "verify": "Audit deprovisioning SLA compliance and test offboarding flow",
        },
        {
            "component_key": "gateway",
            "depends_on": "Access Recertification",
            "text": "MUST be documented for quarterly access review of {component}.",
            "category": "Operational",
            "risk": "Stale access accumulates and expands insider threat surface",
            "verify": "Verify recertification completion and orphaned account remediation",
        },
        {
            "component_key": "identity",
            "depends_on": "Identity Source of Truth",
            "text": "trusts that the HR system feeds accurate identity data to {component}.",
            "category": "Trust",
            "risk": "Inaccurate HR data causes incorrect access decisions",
            "verify": "Reconcile identity records between HR and IAM systems monthly",
        },
    ],
    "Incident Response": [
        {
            "component_key": "monitoring",
            "depends_on": "Security Monitoring",
            "text": "MUST be tested to detect compromise of {component} within SLA.",
            "category": "Derived",
            "risk": "Undetected breaches extend dwell time beyond acceptable limits",
            "verify": "Conduct purple team exercises and measure detection time",
        },
        {
            "component_key": "logs",
            "depends_on": "Log Retention",
            "text": "assumes that logs from {component} are retained for forensic investigation.",
            "category": "Implicit",
            "risk": "Deleted logs prevent post-incident root cause analysis",
            "verify": "Verify log retention policy meets compliance and investigation requirements",
        },
        {
            "component_key": "app",
            "depends_on": "Containment Plan",
            "text": "MUST be documented with containment steps for compromised {component}.",
            "category": "Operational",
            "risk": "Slow containment allows lateral movement to critical systems",
            "verify": "Tabletop-test incident containment procedures for scenario",
        },
        {
            "component_key": "audit_db",
            "depends_on": "Audit Trail Integrity",
            "text": "assumes that {component} logs cannot be tampered with by attackers.",
            "category": "Implicit",
            "risk": "Log tampering hides attacker activity from investigators",
            "verify": "Implement write-once-read-many (WORM) storage for audit logs",
        },
    ],
    "Least Privilege": [
        {
            "component_key": "app",
            "depends_on": "Permission Boundaries",
            "text": "MUST be enforced for all service accounts accessing {component}.",
            "category": "Explicit",
            "risk": "Over-privileged service accounts enable data exfiltration",
            "verify": "Audit service account permissions and reduce to minimum required",
        },
        {
            "component_key": "database",
            "depends_on": "Database Access Control",
            "text": "assumes that database accounts for {component} are scoped to specific tables.",
            "category": "Implicit",
            "risk": "Broad database grants expose all data via any compromised application",
            "verify": "Review database role grants and revoke unnecessary permissions",
        },
        {
            "component_key": "user",
            "depends_on": "Just-in-Time Access",
            "text": "MUST be validated via privilege escalation review for {component}.",
            "category": "Derived",
            "risk": "Standing high privileges increase insider threat impact",
            "verify": "Implement just-in-time access and audit elevation requests",
        },
        {
            "component_key": "gateway",
            "depends_on": "API Key Scope",
            "text": "assumes that API keys for {component} are scoped to the minimum operations required.",
            "category": "Implicit",
            "risk": "Over-scoped API keys allow unintended resource access",
            "verify": "Audit API key permissions and rotate over-scoped keys",
        },
    ],
    "Monitoring & Alerting": [
        {
            "component_key": "monitoring",
            "depends_on": "Alert Configuration",
            "text": "assumes that security alerts from {component} are configured to detect attack patterns.",
            "category": "Implicit",
            "risk": "Misconfigured alerting misses attacker activity",
            "verify": "Validate alert rules against MITRE ATT&CK coverage map",
        },
        {
            "component_key": "monitoring",
            "depends_on": "Alert Response",
            "text": "MUST be documented with escalation paths for alerts from {component}.",
            "category": "Operational",
            "risk": "Alerts with no responder are ignored",
            "verify": "Verify on-call rotation coverage and alert acknowledgement SLA",
        },
        {
            "component_key": "database",
            "depends_on": "Anomaly Detection",
            "text": "MUST be tested to validate baseline deviations for {component}.",
            "category": "Derived",
            "risk": "Gradual data exfiltration below alert threshold goes undetected",
            "verify": "Tune anomaly detection baselines and test with data exfiltration scenarios",
        },
        {
            "component_key": "logs",
            "depends_on": "Log Ingestion",
            "text": "MUST be available to capture all events from {component}.",
            "category": "Dependency",
            "risk": "Log ingestion bottlenecks drop security-relevant events",
            "verify": "Stress test log pipeline and measure event loss rate",
        },
    ],
    "Network Segmentation": [
        {
            "component_key": "network",
            "depends_on": "Firewall Rules",
            "text": "MUST be enforced to isolate {component} from untrusted networks.",
            "category": "Explicit",
            "risk": "Broad network access allows lateral movement from compromise",
            "verify": "Audit security group and firewall rule effectiveness",
        },
        {
            "component_key": "network",
            "depends_on": "Micro-segmentation",
            "text": "assumes that {component} is isolated at the workload level, not just perimeter.",
            "category": "Implicit",
            "risk": "Flat network allows any-to-any communication once inside",
            "verify": "Review network policy configuration and test east-west traffic",
        },
        {
            "component_key": "gateway",
            "depends_on": "Traffic Filtering",
            "text": "MUST be documented showing all allowed traffic flows through {component}.",
            "category": "Operational",
            "risk": "Unreviewed network rules accumulate and expand attack surface",
            "verify": "Maintain and annually review network segmentation matrix",
        },
        {
            "component_key": "network",
            "depends_on": "Network Isolation Testing",
            "text": "MUST be tested for tenant isolation at {component}.",
            "category": "Derived",
            "risk": "Tenant isolation failure leaks cross-tenant data",
            "verify": "Penetration test network isolation between tenants",
        },
    ],
    "Physical Security": [
        {
            "component_key": "database",
            "depends_on": "Data Center Access",
            "text": "assumes that physical access to {component} is restricted to authorized personnel.",
            "category": "Implicit",
            "risk": "Unauthorized physical access allows direct data theft",
            "verify": "Audit data center access logs and badge access records",
        },
        {
            "component_key": "user",
            "depends_on": "Workstation Security",
            "text": "MUST be enforced to prevent unauthorized physical access to {component}.",
            "category": "Explicit",
            "risk": "Unlocked workstations allow unauthorized data access",
            "verify": "Enforce screen lock policy and audit compliance",
        },
        {
            "component_key": "network",
            "depends_on": "Physical Port Security",
            "text": "assumes that physical network ports for {component} are not publicly accessible.",
            "category": "Implicit",
            "risk": "Public network ports allow direct network access bypassing security controls",
            "verify": "Survey physical network port accessibility and secure exposed ports",
        },
    ],
    "Supply Chain Security": [
        {
            "component_key": "gateway",
            "depends_on": "Dependency Scanning",
            "text": "assumes that software dependencies used by {component} have no known vulnerabilities.",
            "category": "Implicit",
            "risk": "Compromised upstream library introduces vulnerability into build",
            "verify": "Automate dependency scanning with Software Bill of Materials",
        },
        {
            "component_key": "app",
            "depends_on": "Image Signing",
            "text": "MUST be enforced to verify artifact signatures before deploying to {component}.",
            "category": "Explicit",
            "risk": "Untrusted images with backdoors deployed to production",
            "verify": "Verify image signature verification in deployment pipeline",
        },
        {
            "component_key": "ci",
            "depends_on": "CI/CD Integrity",
            "text": "MUST be documented with traceable provenance for all artifacts deployed to {component}.",
            "category": "Operational",
            "risk": "Untraceable artifacts cannot be verified or audited",
            "verify": "Implement attestation and provenance tracking in build pipeline",
        },
        {
            "component_key": "network",
            "depends_on": "Vendor Risk Assessment",
            "text": "trusts that third-party components used by {component} are not malicious.",
            "category": "Trust",
            "risk": "Supply chain attack via trusted third-party component",
            "verify": "Conduct vendor security assessment and monitor for breaches",
        },
    ],
    "Third-party Dependency": [
        {
            "component_key": "app",
            "depends_on": "Third-party SLA",
            "text": "MUST be available for {component} to function as specified in the SLA.",
            "category": "Dependency",
            "risk": "Third-party outage causes cascading application failure",
            "verify": "Monitor third-party uptime and test degraded-mode operation",
        },
        {
            "component_key": "gateway",
            "depends_on": "Third-party Security Posture",
            "text": "trusts that the third-party provider for {component} maintains security certifications.",
            "category": "Trust",
            "risk": "Third-party security incident exposes shared data",
            "verify": "Review SOC 2 Type II reports and penetration test results",
        },
        {
            "component_key": "network",
            "depends_on": "Data Residency",
            "text": "assumes that third-party {component} processes data in approved geographic regions.",
            "category": "Implicit",
            "risk": "Data processed in unapproved regions violates compliance requirements",
            "verify": "Audit third-party data center locations and contractual commitments",
        },
        {
            "component_key": "database",
            "depends_on": "Third-party Access",
            "text": "MUST be documented and have data processing agreements for {component}.",
            "category": "Operational",
            "risk": "Undocumented third-party data access breaches compliance",
            "verify": "Maintain and annually review third-party data processing register",
        },
    ],
}


def get_component(arch, key, fallback_pattern):
    """Get component name from architecture, with fallback."""
    if key in arch["components"]:
        return arch["components"][key]
    components = list(arch["components"].values())
    if components:
        return components[0]
    return fallback_pattern


def generate_rows():
    rows = []
    for arch_id in sorted(ARCHITECTURES.keys()):
        arch = ARCHITECTURES[arch_id]
        for pattern_name in arch["patterns"]:
            pattern_id = PATTERN_IDS[pattern_name]
            templates = PATTERN_TEMPLATES[pattern_name]
            for t in templates:
                comp = get_component(arch, t["component_key"], pattern_name)
                text = t["text"].format(component=comp)
                rows.append({
                    "Architecture_ID": arch_id,
                    "Architecture_Name": arch["name"],
                    "Topology": arch["topology"],
                    "Pattern_ID": pattern_id,
                    "Pattern_Name": pattern_name,
                    "Component": comp,
                    "Depends_On": t["depends_on"],
                    "Assumption_Text": text,
                    "Ontology_Category": t["category"],
                    "Risk": t["risk"],
                    "Verification_Method": t["verify"],
                    "Human_Match_ID": "",
                    "Human_Match_Flag": "",
                    "Reviewer_Notes": "",
                })
    return rows


def main():
    output_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)),
        "assumption_gold_standard.csv"
    )

    fieldnames = [
        "Architecture_ID", "Architecture_Name", "Topology",
        "Pattern_ID", "Pattern_Name", "Component", "Depends_On",
        "Assumption_Text", "Ontology_Category", "Risk",
        "Verification_Method", "Human_Match_ID", "Human_Match_Flag",
        "Reviewer_Notes"
    ]

    rows = generate_rows()

    with open(output_path, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)

    file_size = os.path.getsize(output_path)
    print(f"Total rows: {len(rows)}")
    print(f"File size: {file_size:,} bytes ({file_size/1024:.1f} KB)")

    # Architecture coverage summary
    from collections import Counter
    arch_counts = Counter(r["Architecture_ID"] for r in rows)
    print("\nArchitecture coverage:")
    for aid in sorted(arch_counts):
        print(f"  Architecture {aid} ({ARCHITECTURES[aid]['name']}): {arch_counts[aid]} rows")


if __name__ == "__main__":
    main()
