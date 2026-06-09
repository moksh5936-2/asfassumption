"""100 security policies organized by domain with structured metadata."""
from __future__ import annotations
from benchmark.data import Policy


def get_all_policies() -> list[Policy]:
    policies: list[Policy] = []
    pid = 0

    def add(text: str, domain: str, *tags: str):
        nonlocal pid
        policies.append(Policy(id=f"policy_{pid:03d}", text=text, domain=domain, tags=list(tags)))
        pid += 1

    # ── ACCESS CONTROL (20) ──────────────────────────────────────
    add("Only Finance employees may access the payroll system.", "access", "payroll", "finance", "restriction")
    add("SSH access to production servers is restricted to the SRE team.", "access", "ssh", "production", "sre")
    add("Database credentials are scoped to least privilege per application.", "access", "database", "credentials", "least-privilege")
    add("Production environment access requires an approved change ticket.", "access", "production", "change-management")
    add("API keys are scoped to the minimum permissions required.", "access", "api-keys", "least-privilege")
    add("Root access to all systems is logged and reviewed weekly.", "access", "root", "audit", "logging")
    add("Service accounts are restricted to their designated function only.", "access", "service-accounts", "restriction")
    add("VPN access is granted only to active employees.", "access", "vpn", "employee", "authentication")
    add("Customer data access requires a documented business justification.", "access", "customer-data", "justification")
    add("The admin console is accessible only to the security team.", "access", "admin-console", "security-team")
    add("Code repository access is granted based on team membership.", "access", "code-repo", "team-membership")
    add("Cloud console access requires SSO authentication.", "access", "cloud-console", "sso")
    add("Database read replicas are restricted to reporting tools only.", "access", "database", "read-replica", "reporting")
    add("Backup storage access is limited to the backup service account.", "access", "backup", "storage", "service-account")
    add("Monitoring dashboards are view-only for engineering staff.", "access", "monitoring", "read-only", "engineering")
    add("Secret store access is restricted to authorized applications.", "access", "secrets", "vault", "authorization")
    add("Partner portal access is limited to authorized vendors.", "access", "partner", "vendor", "portal")
    add("API gateway access is controlled by API keys per service.", "access", "api-gateway", "api-keys")
    add("Kubernetes cluster admin access is restricted to platform team.", "access", "kubernetes", "cluster-admin", "platform")
    add("File shares are mounted with read-only access for most users.", "access", "file-share", "read-only", "permissions")

    # ── IDENTITY & AUTHENTICATION (15) ──────────────────────────
    add("MFA is required for all administrative access to production.", "identity", "mfa", "administration", "production")
    add("Passwords must be at least 12 characters with complexity requirements.", "identity", "password", "complexity")
    add("All internal applications integrate with SSO.", "identity", "sso", "internal-apps")
    add("Data center physical access requires biometric authentication.", "identity", "biometric", "data-center", "physical")
    add("Service-to-service communication uses certificate-based authentication.", "identity", "certificate", "service-auth")
    add("Session timeout is enforced after 15 minutes of inactivity.", "identity", "session", "timeout", "inactivity")
    add("Account lockout occurs after 5 consecutive failed login attempts.", "identity", "lockout", "failed-login")
    add("Password rotation is required every 90 days.", "identity", "password-rotation", "compliance")
    add("New employee accounts are provisioned within 24 hours of hire.", "identity", "provisioning", "employee")
    add("Orphaned accounts are disabled within 30 days of termination.", "identity", "orphaned-accounts", "offboarding")
    add("Privileged access requires just-in-time approval.", "identity", "jit", "privileged-access")
    add("API authentication uses OAuth2 tokens with refresh rotation.", "identity", "oauth2", "api-auth", "tokens")
    add("Smart card authentication is required for classified systems.", "identity", "smart-card", "classified")
    add("Identity federation is configured between trusted partner organizations.", "identity", "federation", "partners")
    add("Break-glass accounts are monitored and reviewed monthly.", "identity", "break-glass", "emergency", "monitoring")

    # ── NETWORK (15) ────────────────────────────────────────────
    add("The production network is segmented from non-production environments.", "network", "segmentation", "production")
    add("A DMZ isolates public-facing web services from internal networks.", "network", "dmz", "web", "isolation")
    add("Database servers are not directly accessible from the internet.", "network", "database", "internet", "isolation")
    add("A WAF is deployed in front of all web applications.", "network", "waf", "web", "protection")
    add("Internal traffic between services is encrypted.", "network", "encryption", "internal", "traffic")
    add("VPN is required for all remote network access.", "network", "vpn", "remote-access")
    add("Network egress traffic is filtered by application.", "network", "egress", "filtering", "application")
    add("DNS traffic is routed through internal resolvers.", "network", "dns", "internal-resolution")
    add("Load balancers terminate TLS for all incoming connections.", "network", "tls", "load-balancer")
    add("Network ACLs restrict traffic between subnets.", "network", "acl", "subnet", "restriction")
    add("Wireless networks are isolated from the corporate wired network.", "network", "wireless", "isolation")
    add("Bastion hosts are used for all SSH access to internal systems.", "network", "bastion", "ssh", "jump-host")
    add("VPC peering is limited to production and monitoring accounts.", "network", "vpc", "peering", "production")
    add("DDoS protection is enabled on all public-facing endpoints.", "network", "ddos", "protection", "public")
    add("Network flow logs are enabled for all critical network segments.", "network", "flow-logs", "monitoring")

    # ── CONFIGURATION (15) ──────────────────────────────────────
    add("All data is encrypted at rest using AES-256.", "configuration", "encryption", "aes-256", "at-rest")
    add("Backups are encrypted before transfer to storage.", "configuration", "backup", "encryption", "transfer")
    add("Backup retention policy: 30 days of daily, 12 months of monthly backups.", "configuration", "backup", "retention")
    add("Backup restore testing is performed quarterly.", "configuration", "backup", "restore-testing")
    add("Logging is enabled for all production services.", "configuration", "logging", "production")
    add("Audit logs are retained for a minimum of 1 year.", "configuration", "audit", "retention", "logs")
    add("System configurations are baselined and monitored for drift.", "configuration", "baseline", "drift-monitoring")
    add("Security patches are applied within 30 days of release.", "configuration", "patching", "vulnerability")
    add("Anti-malware is installed on all endpoints.", "configuration", "antimalware", "endpoints")
    add("Full disk encryption is enabled on all company laptops.", "configuration", "disk-encryption", "laptops")
    add("Container images are scanned for vulnerabilities before deployment.", "configuration", "container", "vulnerability-scan")
    add("Infrastructure as code is used for all production deployments.", "configuration", "iac", "terraform", "deployment")
    add("Configuration drift alerts are enabled for all critical systems.", "configuration", "drift-alerts", "critical")
    add("TLS 1.2 is the minimum enforced version on all services.", "configuration", "tls", "encryption", "standard")
    add("Secrets are rotated every 90 days.", "configuration", "secrets", "rotation", "credentials")

    # ── PROCESS (15) ────────────────────────────────────────────
    add("All production changes require approval from the change board.", "process", "change-management", "approval")
    add("Incident response procedures are tested quarterly.", "process", "incident-response", "testing")
    add("Disaster recovery plan is tested annually.", "process", "disaster-recovery", "testing")
    add("Code review is required before merging to the main branch.", "process", "code-review", "merge")
    add("Penetration testing is performed annually by an external firm.", "process", "pentest", "external")
    add("Vulnerability scanning is run weekly on all production systems.", "process", "vulnerability-scan", "weekly")
    add("Third-party security assessment is completed before integration.", "process", "third-party", "assessment")
    add("Data retention policies are enforced quarterly.", "process", "data-retention", "enforcement")
    add("User access reviews are conducted quarterly.", "process", "access-review", "quarterly")
    add("Vendor risk assessment is completed before engagement.", "process", "vendor-risk", "assessment")
    add("Security awareness training is completed by all employees annually.", "process", "security-training", "annual")
    add("An insider threat program is in place.", "process", "insider-threat", "program")
    add("Forensics readiness plan is maintained and updated.", "process", "forensics", "readiness")
    add("Data classification labels are enforced on all documents.", "process", "data-classification", "labels")
    add("Business continuity plan is tested every 6 months.", "process", "business-continuity", "testing")

    # ── GOVERNANCE (10) ─────────────────────────────────────────
    add("SOC 2 Type II audit is completed annually.", "governance", "soc2", "audit", "compliance")
    add("GDPR consent management is enforced for all EU customer data.", "governance", "gdpr", "consent", "privacy")
    add("PCI DSS compliance is maintained for all payment data.", "governance", "pci-dss", "compliance", "payments")
    add("HIPAA controls are in place for all protected health information.", "governance", "hipaa", "health", "compliance")
    add("ISO 27001 certification is maintained.", "governance", "iso27001", "certification")
    add("SOX controls are enforced on all financial reporting systems.", "governance", "sox", "financial", "controls")
    add("Data processing agreements are in place with all vendors.", "governance", "dpa", "vendors", "data-processing")
    add("Board-level security risk reporting is delivered quarterly.", "governance", "board-reporting", "risk")
    add("The security policy is reviewed and updated annually.", "governance", "policy-review", "annual")
    add("Acceptable use policy is signed by all employees.", "governance", "aup", "acceptable-use")

    # ── DOCUMENTATION (5) ───────────────────────────────────────
    add("Runbooks are maintained for all critical production systems.", "documentation", "runbooks", "operations")
    add("Architecture diagrams are updated quarterly.", "documentation", "architecture", "diagrams")
    add("Network topology is documented and reviewed annually.", "documentation", "network-topology", "documentation")
    add("Data flow diagrams are maintained for all data processing.", "documentation", "data-flow", "privacy")
    add("Incident response playbooks are kept current.", "documentation", "playbooks", "incident-response")

    # ── DEPENDENCY (5) ──────────────────────────────────────────
    add("All third-party vendors have security SLAs in their contracts.", "dependency", "vendors", "sla", "third-party")
    add("Software supply chain is verified through signed commits.", "dependency", "supply-chain", "signed-commits")
    add("Open source dependencies are scanned for vulnerabilities weekly.", "dependency", "oss", "vulnerability-scan")
    add("Service-level dependencies are documented for all critical applications.", "dependency", "dependencies", "documentation")
    add("Critical system dependencies are mapped and reviewed quarterly.", "dependency", "dependency-map", "review")

    assert len(policies) == 100, f"Expected 100 policies, got {len(policies)}"
    return policies
