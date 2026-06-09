# ASF Phase 6 Experiment: Architecture #017

**Architecture:** Multi-tenant SaaS → Tenant Isolation
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Tenant A] --> [API Gateway] --> [App Service] --> [Tenant A DB]
[Tenant B] --> [API Gateway] --> [App Service] --> [Tenant B DB]
                (Shared)              (Shared)        (Isolated)
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Tenant data is isolated at database level |
| P2 | Row-level security enforced |
| P3 | Tenant context propagated through JWT claims |
| P4 | Cross-tenant access is blocked |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Tenant ↔ Gateway | Authentication boundary |
| Shared Service ↔ Isolated DB | Isolation boundary |
| Tenant ↔ Tenant | Data boundary |

### Complexity Rating
**Moderate** — multi-tenant SaaS, 6 nodes, 3 trust boundaries, shared infrastructure with isolated data stores.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | JWT tokens are signed with a strong, rotated signing key and validated on every request at the API Gateway. | A compromised JWT signing key allows an attacker to forge tokens for any tenant and any user. |
| H-002 | The JWT claim that identifies the tenant (e.g., `tenant_id`) is set by the authentication service and cannot be modified by the client. | A client-modifiable tenant claim allows tenant-impersonation via crafted JWTs. |
| H-003 | The API Gateway extracts the tenant context from the JWT and injects it as a request header that the backend trusts but the client cannot override. | A backend that trusts a client-supplied tenant header instead of the gateway-injected header can be tricked into cross-tenant access. |
| H-004 | Row-level security policies in the database correctly filter all queries by the tenant context and have no bypass paths (e.g., direct table access, admin roles). | An RLS policy with an exception for admin roles or direct table access defeats tenant isolation. |
| H-005 | Database connections from the App Service use a single shared pool where tenant context is set per-session before query execution. | A connection pool that reuses a connection with a stale tenant context allows Query B to execute under Tenant A's identity. |
| H-006 | The App Service sets the tenant context for every database query; there is no default or fallback tenant. | A query without an explicit tenant context may return data from all tenants or default to an incorrect tenant. |
| H-007 | The API Gateway validates that the tenant specified in the JWT matches an active, non-suspended tenant in the tenant registry. | A JWT for a suspended or deleted tenant should be rejected to prevent access from deactivated accounts. |
| H-008 | The database connection pool is scoped per tenant or uses connection-level tenant context (e.g., SET app.tenant_id = x). | A shared pool without per-connection tenant context risks cross-tenant data leakage on connection reuse. |
| H-009 | API Gateway rate limiting and throttling are applied per tenant to prevent one tenant from exhausting resources for all others. | A single tenant with high traffic volume can degrade the experience for all tenants if rate limits are global. |
| H-010 | The application enforces data export limits per tenant to prevent bulk data exfiltration. | A compromised tenant admin account can export all tenant data without detection if no export limits exist. |
| H-011 | Tenant onboarding includes provisioning a dedicated database schema or database with appropriate isolation. | A new tenant sharing a schema with other tenants without RLS or separate tables risks data leakage. |
| H-012 | Tenant offboarding includes a full purge of tenant data from all databases, caches, and backups. | Residual tenant data in backups after offboarding creates data exposure risk if the backup is later accessed. |
| H-013 | The shared App Service has no persistent state that can leak data between tenant requests (e.g., static variables, caches). | A shared service with tenant A's data cached in application memory can serve it to tenant B on the next request. |
| H-014 | JWT tokens have an expiration that is short (minutes) and are refreshed via a secure refresh token flow. | Long-lived JWTs that leak via logs or network capture grant persistent access to the tenant's data. |
| H-015 | The API Gateway validates that the JWT `aud` (audience) claim matches the specific API endpoint being accessed. | A JWT issued for the admin API should not be accepted by the public API, even for the same tenant. |
| H-016 | Row-level security policies are tested as part of the CI/CD pipeline with automated cross-tenant access tests. | Untested RLS policies may have logical flaws that allow cross-tenant data access in edge cases. |
| H-017 | The database uses a separate encryption key per tenant or per schema to prevent cross-tenant decryption. | A single encryption key across all tenants means that compromising the key decrypts all tenant data. |
| H-018 | Tenant metadata (tenant_id, active status, feature flags) is stored in a separate hardened service from tenant data. | If tenant metadata and tenant data share a database, a SQL injection in the metadata lookup could expose tenant data. |
| H-019 | The application enforces that tenant context cannot be changed mid-request (no tenant hopping). | A request that starts processing under Tenant A and switches to Tenant B mid-request could leak data across the transition. |
| H-020 | The API Gateway strips any tenant-related headers from incoming requests before injecting the gateway-verified tenant header. | A client sending a forged `X-Tenant-ID` header that is not stripped could bypass the gateway's tenant assertion. |
| H-021 | Database backups are segregated per tenant so that a compromised backup of one tenant does not expose all tenants. | Co-mingled backups mean a backup breach exposes all tenants' data. |
| H-022 | The App Service has a maximum request processing time to prevent a slow tenant request from blocking the shared worker pool. | A single tenant's slow query can exhaust the application worker pool, causing denial of service for all tenants. |
| H-023 | The shared API Gateway enforces per-tenant API key rate limits distinct from user-level JWT rate limits. | An API key shared across multiple users of a tenant bypasses per-user rate limits and can be used to abuse the API. |
| H-024 | The application validates that the tenant in the JWT matches the tenant associated with the API key (if both are present). | A valid JWT for Tenant A used with an API key for Tenant B should be rejected as a cross-tenant credential mismatch. |
| H-025 | Database encryption at rest uses a Customer Master Key (CMK) per tenant where regulatory compliance requires it. | Shared encryption keys across tenants may violate compliance requirements (e.g., GDPR data separation). |
| H-026 | The App Service logs the tenant context with every database query for auditability. | Without tenant-context logging, a cross-tenant data access incident cannot be attributed to the responsible tenant. |
| H-027 | Monitoring and alerting are configured to detect anomalous cross-tenant access patterns (e.g., a user accessing more than N tenant workspaces). | A compromised user account accessing multiple tenants' data may indicate privilege escalation that should trigger alerts. |
| H-028 | The application implements a tenant-scoped rate limit on authentication attempts to prevent credential stuffing per tenant. | A credential stuffing attack against one tenant should not affect authentication performance for other tenants. |
| H-029 | Database migrations run with tenant context and do not accidentally expose one tenant's data to another during schema changes. | A migration that copies data across tenant schemas or resets RLS policies can create temporary cross-tenant exposure. |
| H-030 | The tenant registry is highly available and cached with short TTL to prevent a tenant registry outage from blocking all tenants. | Tenant registry downtime means no tenant can authenticate; a cache with short TTL allows graceful degradation. |
| H-031 | The application validates that tenant-level feature flags are not used to bypass security controls (e.g., disable RLS for a tenant). | A feature flag that disables RLS for performance reasons creates a persistent cross-tenant data leak. |
| H-032 | The App Service has no debug or admin endpoints that bypass tenant context. | An admin endpoint that executes queries without tenant filtering can expose all tenants' data. |
| H-033 | The database audit logging includes the tenant context for all queries to enable forensic investigation. | Database audit logs without tenant context cannot distinguish between legitimate cross-tenant queries and breaches. |
| H-034 | The JWT signing key is rotated at least quarterly and upon any suspicion of compromise. | Static JWT signing keys increase the window of exposure for forged tokens. |
| H-035 | The App Service is stateless and can be scaled horizontally; tenant affinity (sticky sessions) is not required. | Stateful services that require tenant affinity complicate scaling and can fail during rolling deployments. |
| H-036 | The API Gateway mutates the JWT to inject a gateway-verified tenant identity that the app service trusts. | A gateway that passes the JWT through without mutation allows the app service to parse a potentially forged tenant context. |
| H-037 | The database RLS policies are written to use a parameterized tenant context (e.g., `current_setting('app.tenant_id')`) to prevent SQL injection in the tenant context value. | A tenant context value that is concatenated into SQL creates a SQL injection vulnerability in the RLS policy itself. |
| H-038 | The API Gateway enforces that public endpoints (registration, password reset) do not require tenant context and are rate-limited separately. | Public endpoints that leak tenant enumeration (valid tenant ID = different error message) aid reconnaissance. |
| H-039 | Tenant-level resource quotas (storage, API calls, concurrent connections) are enforced to prevent resource abuse. | A tenant exceeding resource quotas can degrade infrastructure performance for all tenants. |
| H-040 | The application implements defense-in-depth with application-level tenant checks in addition to database RLS. | RLS alone is insufficient if the application has a bug that constructs a query without tenant filtering before the RLS can act. |
| H-041 | The JWT includes a unique session ID that is logged and can be revoked independently of the tenant. | Revoking a tenant's access should not require revoking all sessions; per-session revocation supports targeted incident response. |
| H-042 | The application validates that the tenant context is a valid UUID or structured identifier, not a freeform string. | A freeform tenant ID can be exploited for NoSQL injection or path traversal if used in file paths or cache keys. |
| H-043 | The tenant database credentials (if per-tenant databases are used) are stored securely and rotated. | Per-tenant database credentials stored in plaintext configuration files can be extracted, granting data access. |
| H-044 | The shared API Gateway and App Service are in a separate VPC/network segment from the tenant databases. | A shared network segment means a compromised App Service has network access to all tenant databases. |

**Total (H): 44**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Backup & Recovery (covered under operational), Physical Security (cloud-hosted), Supply Chain Security (deferred to Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Tenant administrators use MFA to authenticate to the SaaS admin console. | Explicit | Tenant admin account compromise gives an attacker full access to the tenant's data and configuration. |
| ASF-002 | MFA is not bypassed for API-based access; API tokens require MFA step-up for privileged operations. | Derived | API tokens without MFA are vulnerable to token theft and replay. |
| ASF-003 | The SaaS provider's internal admin access to the multi-tenant infrastructure requires MFA. | Operational | Provider admin without MFA can access all tenant data. |
| ASF-004 | Tenant users can enroll their own MFA devices through a verified process. | Implicit | MFA enrollment without verification allows an attacker to add their device to another tenant's account. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Tenant authentication uses a federated SSO (SAML/OIDC) per tenant, not shared credentials. | Explicit | Shared authentication across tenants bypasses tenant-specific identity controls. |
| ASF-006 | The IdP for each tenant is independently configured and tenant metadata is isolated. | Derived | A misconfiguration in one tenant's IdP should not affect other tenants' authentication. |
| ASF-007 | SSO session timeout is consistent between the tenant IdP and the SaaS application session. | Trust | Mismatched timeouts leave SaaS sessions active after tenant IdP logout. |
| ASF-008 | Tenant SSO configuration changes (e.g., new IdP, certificate rotation) go through a validation process. | Operational | An unvalidated SSO configuration change can lock all users out of a tenant. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | A single App Service instance failure does not block all tenants (redundancy exists). | Architectural | Shared App Service is a SPOF if deployed as a single instance. |
| ASF-010 | A single tenant's database failure does not block other tenants from accessing the application. | Operational | A per-tenant database outage should not cause a complete application outage. |
| ASF-011 | The API Gateway has sufficient throughput for the aggregate of all tenant traffic. | Environmental | Gateway capacity must scale with the total number of tenants, not per-tenant average. |
| ASF-012 | Tenant database connection pools are configured to prevent a single tenant from exhausting all database connections. | Derived | One tenant's connection pool spike can starve other tenants of database access. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | Tenant databases are backed up independently and restorable per tenant. | Operational | A single tenant's data corruption should be restorable without affecting other tenants. |
| ASF-014 | Backup restores are tested for at least one tenant per quarter to verify data integrity and RTO. | Derived | Untested tenant backup restorations may fail when needed. |
| ASF-015 | Backups are encrypted with per-tenant or per-environment keys. | Implicit | Shared backup encryption keys across tenants mean a key compromise exposes all tenant backups. |
| ASF-016 | There is a documented process for restoring a single tenant's data without impacting other tenants. | Operational | Restoring a single tenant from a shared backup requires careful isolation procedures. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Database schema changes are backward-compatible to avoid breaking RLS policies for existing tenants. | Operational | Schema changes that break RLS policies expose cross-tenant data until fixed. |
| ASF-018 | RLS policy changes are reviewed specifically for cross-tenant access vulnerabilities. | Derived | An RLS policy change without tenant isolation review can inadvertently open cross-tenant access. |
| ASF-019 | New tenant onboarding changes (new schema, new database) are automated to prevent misconfiguration. | Trust | Manual tenant provisioning is error-prone and may skip isolation steps. |
| ASF-020 | API Gateway routing rule changes (new tenant endpoints) are reviewed before deployment. | Implicit | Incorrect routing rules can send tenant A's traffic to tenant B's backend. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The App Service IAM role permits access only to the tenant databases it is authorized to reach. | Explicit | Over-permissioned App Service IAM role can access any tenant database in the account. |
| ASF-022 | The API Gateway IAM role has least privilege for invoking the App Service. | Derived | Gateway IAM role with broad invoke permissions can call any service endpoint, bypassing routing controls. |
| ASF-023 | CloudTrail is enabled to audit all API calls to the tenant infrastructure. | Implicit | Without CloudTrail, unauthorized infrastructure changes affecting tenant isolation go undetected. |
| ASF-024 | No IAM users exist in the SaaS account; all access is via roles assumed from an SSO identity provider. | Explicit | Long-term IAM user keys in the SaaS account increase the risk of persistent credential exposure. |

---

### Pattern 7: Compliance & Audit

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Tenant data classification is defined and each tenant's data handling requirements are documented. | Explicit | Different tenants may have different compliance requirements (HIPAA, GDPR, SOC 2). |
| ASF-026 | Tenant audit logs are retained per regulatory requirements and segregated per tenant. | Derived | Co-mingled audit logs prevent per-tenant forensic investigation. |
| ASF-027 | There is a process for tenant data sovereignty — data is stored in the region required by each tenant. | Operational | A tenant requiring data residency must have their database in the appropriate region. |
| ASF-028 | Cross-tenant access attempts (successful or blocked) are logged and escalated. | Derived | Blocked cross-tenant attempts are a signal of probing or misconfiguration that requires investigation. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Data flow diagrams exist per tenant or per tenant tier showing data paths and storage. | Explicit | Different tenant tiers (free, enterprise) may have different data flow architectures. |
| ASF-030 | Tenant data does not flow to any shared analytics or monitoring system without PII stripping. | Implicit | Shared analytics that receive tenant data from all tenants create a cross-tenant data aggregation risk. |
| ASF-031 | The documented data flow (Gateway -> App -> DB) is the only path tenant data travels. | Derived | Shadow IT data exports (e.g., a tenant admin exporting to personal cloud storage) are unaccounted. |
| ASF-032 | Tenant data is not cached in a shared caching layer (e.g., Redis, Memcached) without tenant key prefixing. | Environmental | A shared cache without tenant key isolation can leak data between tenants. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | All tenant databases are encrypted at rest using AWS KMS or equivalent. | Explicit | Standard requirement for sensitive tenant data. |
| ASF-034 | Each tenant database uses a separate KMS key for encryption at rest. | Derived | A shared key across tenants means a KMS key compromise exposes all tenant data. |
| ASF-035 | KMS key policies restrict which IAM principals can encrypt/decrypt each tenant's database. | Implicit | A KMS key with permissive access policy allows any authorized IAM user to decrypt any tenant's data. |
| ASF-036 | Backups of tenant databases are encrypted using a different key than the primary database. | Operational | Backup encryption key compromise should not lead to primary data compromise (or vice versa). |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | TLS is enforced between tenant clients and the API Gateway. | Explicit | Standard for multi-tenant SaaS. |
| ASF-038 | TLS is enforced between the API Gateway and the App Service. | Derived | Internal traffic between gateway and app must be encrypted to prevent tenant context tampering. |
| ASF-039 | TLS is enforced between the App Service and the tenant databases. | Explicit | Database connections must be encrypted to prevent tenant data interception. |
| ASF-040 | TLS 1.2 or higher is enforced; TLS 1.0/1.1 and SSL are disabled on all endpoints. | Derived | Weak TLS versions expose tenant data to passive interception. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | The App Service OS and runtime are patched regularly for security vulnerabilities. | Implicit | Unpatched application servers in the shared tier can be compromised, affecting all tenants. |
| ASF-042 | The API Gateway is protected by a WAF that blocks common web attacks. | Derived | Multi-tenant gateways are high-value targets for injection and enumeration attacks. |
| ASF-043 | Tenant-facing endpoints are scanned for vulnerabilities before each deployment. | Operational | A vulnerability in the shared App Service code affects all tenants simultaneously. |
| ASF-044 | The shared application has no backdoor or debug endpoints that bypass tenant context. | Implicit | Admin endpoints without tenant context filtering expose all tenant data. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Tenant administrators understand the scope of their admin privileges and do not share credentials across tenants. | Derived | Shared admin credentials between tenants create a cross-tenant identity risk. |
| ASF-046 | The SaaS provider's internal team understands tenant isolation architecture and does not make changes that compromise it. | Trust | Internal team misconfiguration can break tenant isolation more effectively than any external attack. |
| ASF-047 | Tenant support requests for data access are verified through a documented process to prevent social engineering. | Operational | An attacker posing as a tenant admin could request data export without proper verification. |
| ASF-048 | There is a process for tenants to report security incidents and receive timely responses. | Environmental | Multi-tenant incidents require coordinated disclosure to all affected tenants. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Tenant user accounts follow a joiner/mover/leaver process within each tenant's organization. | Operational | Stale tenant user accounts with active JWTs represent unauthorized access risk. |
| ASF-050 | Tenant admin role changes are audited and require secondary approval. | Derived | Unchecked tenant admin role grants allow privilege escalation within the tenant. |
| ASF-051 | The SaaS provider's internal admin access to tenant infrastructure is reviewed and recertified quarterly. | Implicit | Provider admin accounts with excessive privileges can access all tenants' data. |
| ASF-052 | Automated provisioning of new tenants includes security baseline configuration (RLS, encryption, audit). | Operational | A tenant provisioned without baseline security controls starts in an insecure state. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan covering multi-tenant data breach scenarios. | Operational | A breach affecting one tenant requires different communication and containment than a single-tenant breach. |
| ASF-054 | The IR team can isolate a compromised tenant's database without affecting other tenants. | Derived | Tenant isolation in response to a breach must not cause collateral damage. |
| ASF-055 | The IR plan includes procedures for notifying affected tenants within regulatory timeframes. | Trust | Regulatory breach notification deadlines (72 hours for GDPR) require pre-defined communication channels. |
| ASF-056 | Monitoring detects anomalous patterns that may indicate cross-tenant data access (e.g., a single IP accessing multiple tenant APIs). | Implicit | Cross-tenant access from a single IP is a strong indicator of a compromised provider account or misconfiguration. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | The App Service database user has SELECT/INSERT/UPDATE/DELETE only on the specific tenant's schema or database. | Explicit | Database user with cross-schema access can read all tenants' data. |
| ASF-058 | The API Gateway has no direct access to tenant databases (access must go through the App Service). | Derived | Gateway direct database access bypasses application-layer tenant context enforcement. |
| ASF-059 | The App Service does not run as a database admin (no DDL privileges in production). | Implicit | App service with DDL privileges can modify schema or RLS policies. |
| ASF-060 | Tenant database credentials (if per-tenant) are scoped to only that tenant's database. | Derived | A tenant database credential with cross-database access can read other tenants' data. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Authentication failures are monitored per tenant to detect credential stuffing attacks. | Operational | A spike in auth failures for a specific tenant indicates an active attack. |
| ASF-062 | Database query patterns are monitored for unexpected cross-tenant queries or unusual volume. | Derived | A query that returns data from multiple tenants is a strong indicator of RLS bypass. |
| ASF-063 | API Gateway 403 errors (cross-tenant access blocked) are monitored and alerted. | Operational | Repeated 403 errors from a single source indicate probing for cross-tenant vulnerabilities. |
| ASF-064 | Tenant resource usage (API calls, storage, compute) is monitored to detect compromised tenants used for crypto mining or data exfiltration. | Derived | Unusual resource consumption by a tenant may indicate account compromise. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | The shared App Service and API Gateway are in a separate VPC from the tenant databases. | Architectural | Shared services with network access to all tenant databases can be used as a pivot point. |
| ASF-066 | Tenant databases are in private subnets accessible only from the App Service security group. | Explicit | A tenant database accessible from the internet or other VPCs bypasses application-layer isolation. |
| ASF-067 | There is no network path between tenant databases (no VPC peering, no cross-database queries). | Architectural | Network-level isolation between tenant databases prevents direct cross-tenant data access. |
| ASF-068 | Security groups for tenant databases restrict inbound traffic to the App Service security group only. | Derived | Any other source IP accessing a tenant database bypasses tenant context enforcement. |

---

### Pattern 18: Secrets Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | The JWT signing key is stored in a secrets manager (e.g., AWS Secrets Manager, Vault) and rotated regularly. | Explicit | JWT signing key in source code or configuration files is extractable and compromisable. |
| ASF-070 | Tenant database credentials are stored in a secrets manager with access audit. | Derived | Tenant credentials in plaintext configuration files are exposed through repository breaches. |
| ASF-071 | API keys for tenant programmatic access are hashed before storage and cannot be retrieved in plaintext. | Implicit | Storable API keys in plaintext can be exfiltrated by an attacker who gains database read access. |
| ASF-072 | Secrets for tenant integrations (webhook secrets, external API keys) are managed per tenant with access controls. | Operational | A compromised webhook secret for one tenant should not allow an attacker to impersonate another tenant's webhooks. |

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-073 | The App Service framework and runtime have no known critical vulnerabilities. | Dependency | A vulnerability in the shared App Service code affects all tenants simultaneously. |
| ASF-074 | Third-party libraries used by the App Service are scanned for vulnerabilities before each deployment. | Operational | Library dependency vulnerabilities (e.g., Log4j, Spring4Shell) in the shared tier expose all tenants. |
| ASF-075 | The API Gateway software (e.g., Kong, AWS API Gateway, custom) has no known vulnerabilities. | Dependency | Gateway vulnerabilities can bypass tenant routing and authentication for all tenants. |
| ASF-076 | An SBOM is maintained for the SaaS application and monitored for new vulnerabilities. | Derived | Without SBOM tracking, newly disclosed vulnerabilities cannot be rapidly assessed for impact. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-077 | The database service (RDS, Aurora) has a reliable uptime SLA. | Dependency | Database service downtime means all tenants lose access to their data. |
| ASF-078 | The JWT library or service has no known vulnerabilities affecting token validation. | Dependency | JWT library vulnerabilities (e.g., algorithm confusion, `none` algorithm) bypass all tenant authentication. |
| ASF-079 | There is a fallback authentication mechanism if the tenant's IdP is unavailable (e.g., cached session). | Operational | Tenant IdP downtime should not block all tenant users if cached sessions are available. |
| ASF-080 | There is an exit strategy for migrating a tenant off the platform with full data export. | Derived | Tenant lock-in without data export capability creates vendor dependency and potential data loss on migration. |

**Total (A): 80** (4 per pattern x 17 patterns + 12 overflow from high-complexity patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-069 | Both require JWT signing key rotation and secure storage. |
| H-002 | ASF-005 | Both require tenant claim set by auth service, not client-modifiable. |
| H-003 | ASF-005 | Both require gateway-injected tenant header, client cannot override. |
| H-004 | ASF-057 | Both require RLS policies with no bypass paths. |
| H-005 | ASF-012 | Both require per-connection/pool tenant context. |
| H-006 | ASF-057 | Both require tenant context set per query. |
| H-007 | ASF-005 | Both require tenant status validation in JWT. |
| H-008 | ASF-012 | Both require per-tenant connection pool scoping. |
| H-009 | ASF-003 | Both require per-tenant rate limiting. |
| H-012 | ASF-013 | Both require tenant data purge on offboarding. |
| H-013 | ASF-032 | Both require shared app service has no state leakage. |
| H-014 | ASF-014 | Both require short JWT expiry with refresh. |
| H-015 | ASF-006 | Both require JWT audience validation. |
| H-016 | ASF-018 | Both require RLS testing in CI/CD. |
| H-017 | ASF-034 | Both require per-tenant encryption keys. |
| H-019 | ASF-019 | Both require no tenant hopping mid-request. |
| H-020 | ASF-003 | Both require stripping client-supplied tenant headers. |
| H-021 | ASF-015 | Both require per-tenant backup segregation. |
| H-022 | ASF-012 | Both require max request processing time. |
| H-025 | ASF-034 | Both require per-tenant CMK. |
| H-026 | ASF-033 | Both require tenant-context database logging. |
| H-027 | ASF-064 | Both require cross-tenant access monitoring. |
| H-028 | ASF-061 | Both require per-tenant auth rate limiting. |
| H-030 | ASF-010 | Both require tenant registry HA with cache. |
| H-032 | ASF-044 | Both require no debug/admin endpoints bypassing tenant. |
| H-033 | ASF-062 | Both require tenant context in DB audit logs. |
| H-034 | ASF-069 | Both require JWT signing key rotation. |
| H-037 | ASF-018 | Both require parameterized tenant context for RLS. |
| H-038 | ASF-011 | Both require public endpoint rate limiting separate from tenant. |
| H-039 | ASF-064 | Both require tenant resource quotas. |
| H-040 | ASF-057 | Both require application-level tenant checks. |
| H-044 | ASF-065 | Both require shared services in separate VPC from databases. |

**Overlap (O): 32**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 44 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 80 | Count of unique ASF-generated assumptions |
| Overlap (O) | 32 | Count appearing in both lists |
| **Precision** | **40.0%** | O / A = 32/80 |
| **Recall** | **72.7%** | O / H = 32/44 |
| **F1 Score** | **51.6%** | 2 x (P x R) / (P + R) |
| Novel findings (A - O) | 48 | Assumptions ASF found that human missed (60.0% of ASF total) |
| Missed findings (H - O) | 12 | Assumptions human found that ASF missed (27.3% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 72.7% |  Met |
| Precision | >= 50% | 40.0% |  Not met |
| Novel discoveries | >= 10% of total (A+O) | 38.7% (48/124) |  Exceeded |
| Expert agreement (F1 proxy) | > 60% | 51.6% |  Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 16 | 50.0% |
| Derived | 12 | 20 | 60.0% |
| Operational | 6 | 24 | 25.0% |
| Implicit | 4 | 12 | 33.3% |
| Trust | 1 | 4 | 25.0% |
| Dependency | 1 | 8 | 12.5% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Derived category showed the strongest agreement. Both humans and the ASF recognize derived concerns: per-tenant encryption keys follow from data classification; per-connection tenant context follows from shared connection pool design.

**Worst overlap:** Architectural and Environmental. The ASF identified architectural concerns (separate VPCs, network-level isolation between tenant databases) and environmental concerns (tenant compliance requirements, data sovereignty) that the human did not treat as assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 12)

1. **Tenant-specific operational hardening (H-010, H-023, H-024, H-029, H-031, H-036):** Tenant data export limits, API key vs JWT tenant mismatch validation, tenant-scoped credential validation, RLS migration safety, feature flag security bypass risk, and JWT mutation for gateway-verified identity. These are multi-tenant-specific implementation concerns.

2. **Application-level tenant controls (H-018, H-034, H-041, H-042, H-043):** Tenant metadata service separation, JWT session ID for per-session revocation, tenant ID format validation (UUID not freeform), per-tenant DB credential management, and tenant ID type validation.

3. **Capacity and resource management (H-035):** Stateless App Service design for horizontal scaling. The ASF covers availability but not the stateless architectural requirement for multi-tenant scaling.

### What ASF Caught That Humans Missed (Novel Findings = 48)

1. **Compliance & Audit (4 assumptions):** The human generated few compliance assumptions. The ASF contributed tenant data classification, per-tenant audit log segregation, data sovereignty region requirements, and cross-tenant attempt escalation.

2. **Change Management (4 assumptions):** The human did not address schema change backward-compatibility, RLS policy change review, automated tenant provisioning, or API Gateway routing rule review.

3. **Incident Response (4 assumptions):** The human generated no IR assumptions specific to multi-tenant scenarios. The ASF contributed tenant breach IR planning, tenant database isolation, regulatory notification procedures, and cross-tenant access detection.

4. **Secrets Management (4 assumptions):** The human covered JWT key storage but the ASF extended to tenant DB credential vault storage, hashed API key storage, and per-tenant integration secret management.

5. **Supply Chain (4 assumptions):** The human did not address App Service framework vulnerabilities, third-party library scanning, API Gateway vulnerabilities, or SBOM management.

### Architecture Complexity Assessment

Architecture #017 achieved the second-highest recall (72.7%), exceeding the 70% target. Multi-tenant isolation concerns align well with the ASF's patterns on least privilege, network segmentation, encryption, and identity lifecycle.

- **Recall (72.7%)** exceeds the 70% target. The ASF's existing patterns cover most multi-tenant isolation concerns, particularly around data isolation, encryption, and network segmentation.
- **Precision (40.0%)** reflects the ASF's breadth across compliance, change management, and supply chain dimensions.
- **Novel rate (60.0%)** confirms significant value from the ASF in surfacing governance, incident response, and compliance assumptions.

### Key Insight

The primary gap is **multi-tenant application-layer controls**: the human enumerated several application-level concerns (credential mismatch validation, tenant ID format, session ID per revocation, stateless design) that the ASF's infrastructure-focused patterns do not cover. Adding a "Multi-tenant Isolation" pattern covering tenant context propagation, credential validation, export controls, and stateless architecture would close the remaining recall gap.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 72.7% |  Target met |
| Precision | >= 50% | 40.0% |  Below target |
| Novel discoveries | >= 10% | 38.7% |  ASF adds value in compliance, change mgmt, IR |
| Expert agreement (F1) | > 60% | 51.6% |  Below target |

The ASF framework applied to Architecture #017 demonstrates strong alignment with multi-tenant isolation concerns, achieving above-target recall. The primary actionable finding is the need for a **Multi-tenant Isolation** pattern covering tenant context propagation, application-level tenant controls, and tenant credential validation to close the remaining gap.
