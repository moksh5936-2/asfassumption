# ASF Assumption Ontology

**Total assumptions analyzed:** 1697

## Ontology Distribution

| Ontology | Count | Percent |
|----------|-------|--------|
| Explicit | 100 | 5.9% |
| Implicit | 282 | 16.6% |
| Derived | 582 | 34.3% |
| Trust | 142 | 8.4% |
| Operational | 339 | 20.0% |
| Dependency | 225 | 13.3% |
| Architectural | 7 | 0.4% |
| Environmental | 20 | 1.2% |
| **TOTAL** | **1697** | **100%** |

## Missed by Ontology

| Ontology | Missed | Total | % Missed |
|----------|--------|-------|----------|
| Explicit | 59 | 100 | 59.0% |
| Implicit | 249 | 282 | 88.3% |
| Derived | 537 | 582 | 92.3% |
| Trust | 109 | 142 | 76.8% |
| Operational | 325 | 339 | 95.9% |
| Dependency | 203 | 225 | 90.2% |
| Architectural | 4 | 7 | 57.1% |
| Environmental | 12 | 20 | 60.0% |

## Sample Assumptions per Ontology

### Explicit

- **ACCESS**: Database credentials are scoped to least privilege per application.
  _Directly stated security requirement._

- **ACCESS**: API keys are scoped to the minimum permissions required.
  _Directly stated security requirement._

- **ACCESS**: VPN access is granted only to active employees.
  _Directly stated security requirement._

- **ACCESS**: Code repository access is granted based on team membership.
  _Directly stated security requirement._

- **ACCESS**: Partner portal access is limited to authorized vendors.
  _Directly stated security requirement._

- **ACCESS**: API gateway access is controlled by API keys per service.
  _Directly stated security requirement._

- **ACCESS**: File shares are mounted with read-only access for most users.
  _Directly stated security requirement._

- **IDENTITY**: Passwords must be at least 12 characters with complexity requirements.
  _Directly stated security requirement._

- **IDENTITY**: All internal applications integrate with SSO.
  _Directly stated security requirement._

- **IDENTITY**: Service-to-service communication uses certificate-based authentication.
  _Directly stated security requirement._

### Implicit

- **GOVERNANCE**: This policy does not conflict with other security policies.
  _Conflicting policies create unresolvable compliance ambiguity._

- **CONFIGURATION**: Evidence collection does not introduce new security vulnerabilities.
  _Monitoring and logging agents expand the attack surface._

- **GOVERNANCE**: This policy does not conflict with other security policies.
  _Conflicting policies create unresolvable compliance ambiguity._

- **CONFIGURATION**: Evidence collection does not introduce new security vulnerabilities.
  _Monitoring and logging agents expand the attack surface._

- **DEPENDENCY**: The database system exists and is operational.
  _Policy references a system that must be running and maintained._

- **CONFIGURATION**: The database system is properly configured and maintained.
  _System configuration is assumed correct without explicit verification._

- **ACCESS**: Least privilege principle is correctly implemented, not just documented.
  _Documented principle vs. actual implementation gap is a common risk._

- **GOVERNANCE**: This policy does not conflict with other security policies.
  _Conflicting policies create unresolvable compliance ambiguity._

- **CONFIGURATION**: Evidence collection does not introduce new security vulnerabilities.
  _Monitoring and logging agents expand the attack surface._

- **GOVERNANCE**: This policy does not conflict with other security policies.
  _Conflicting policies create unresolvable compliance ambiguity._

### Derived

- **CONFIGURATION**: No compensating controls exist if this policy fails.
  _Defense-in-depth requires compensating controls; policy assumes they exist._

- **GOVERNANCE**: Compliance with this policy is measured and reported.
  _Unmeasured policies are effectively optional._

- **GOVERNANCE**: Exceptions to this policy are tracked and approved.
  _Policy assumes exceptions follow formal waiver process._

- **GOVERNANCE**: This policy covers all instances of payroll without exception.
  _Partial coverage creates blind spots._

- **CONFIGURATION**: Violations of this policy are detected and reported.
  _Undetected violations are indistinguishable from compliance._

- **DOCUMENTATION**: Evidence exists to verify compliance with this policy.
  _Unverifiable policies cannot be audited._

- **CONFIGURATION**: No compensating controls exist if this policy fails.
  _Defense-in-depth requires compensating controls; policy assumes they exist._

- **GOVERNANCE**: Compliance with this policy is measured and reported.
  _Unmeasured policies are effectively optional._

- **GOVERNANCE**: Exceptions to this policy are tracked and approved.
  _Policy assumes exceptions follow formal waiver process._

- **CONFIGURATION**: Violations of this policy are detected and reported.
  _Undetected violations are indistinguishable from compliance._

### Trust

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

- **ACCESS**: Deny rules are enforced for all unauthorized access attempts.
  _Restriction only works if denial is actually enforced at the control plane._

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

- **ACCESS**: Violation of this policy leads to unauthorized data access or exfiltration.
  _The risk consequence justifies the policy's existence._

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

- **ACCESS**: Deny rules are enforced for all unauthorized access attempts.
  _Restriction only works if denial is actually enforced at the control plane._

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

- **ACCESS**: Violation of this policy leads to unauthorized data access or exfiltration.
  _The risk consequence justifies the policy's existence._

- **IDENTITY**: User identities are verified before authorization decisions.
  _Authorization presupposes correct identification._

### Operational

- **CONFIGURATION**: Access attempts are logged for audit purposes.
  _Without logging, policy enforcement cannot be verified._

- **CONFIGURATION**: Audit logs are retained and reviewable for compliance.
  _Log retention is necessary for post-incident investigation._

- **PROCESS**: Supporting processes for this policy are documented and followed.
  _Policy assumes operational processes exist to implement it._

- **PROCESS**: Personnel are trained on this policy consistently.
  _Policy compliance depends on human awareness and training._

- **CONFIGURATION**: Access attempts are logged for audit purposes.
  _Without logging, policy enforcement cannot be verified._

- **CONFIGURATION**: Audit logs are retained and reviewable for compliance.
  _Log retention is necessary for post-incident investigation._

- **PROCESS**: Supporting processes for this policy are documented and followed.
  _Policy assumes operational processes exist to implement it._

- **PROCESS**: Personnel are trained on this policy consistently.
  _Policy compliance depends on human awareness and training._

- **CONFIGURATION**: Access attempts are logged for audit purposes.
  _Without logging, policy enforcement cannot be verified._

- **CONFIGURATION**: Audit logs are retained and reviewable for compliance.
  _Log retention is necessary for post-incident investigation._

### Dependency

- **IDENTITY**: Identity provider is available and correctly configured.
  _IdP downtime breaks all downstream access decisions._

- **GOVERNANCE**: This policy is consistently enforced across all environments.
  _Inconsistent enforcement creates exploitable gaps._

- **IDENTITY**: Identity provider is available and correctly configured.
  _IdP downtime breaks all downstream access decisions._

- **GOVERNANCE**: This policy is consistently enforced across all environments.
  _Inconsistent enforcement creates exploitable gaps._

- **ACCESS**: Authorization decisions are consistent across all enforcement points.
  _Policy assumes no gaps in coverage between different access control systems._

- **IDENTITY**: Identity provider is available and correctly configured.
  _IdP downtime breaks all downstream access decisions._

- **GOVERNANCE**: This policy is consistently enforced across all environments.
  _Inconsistent enforcement creates exploitable gaps._

- **IDENTITY**: Identity provider is available and correctly configured.
  _IdP downtime breaks all downstream access decisions._

- **GOVERNANCE**: This policy is consistently enforced across all environments.
  _Inconsistent enforcement creates exploitable gaps._

- **ACCESS**: Authorization decisions are consistent across all enforcement points.
  _Policy assumes no gaps in coverage between different access control systems._

### Architectural

- **NETWORK**: Network segmentation is enforced at the data link layer, not just documented.
  _Logical segmentation without enforcement is security theater._

- **DEPENDENCY**: The bastion system exists and is operational.
  _Policy references a system that must be running and maintained._

- **CONFIGURATION**: The bastion system is properly configured and maintained.
  _System configuration is assumed correct without explicit verification._

- **GOVERNANCE**: This policy covers all instances of bastion without exception.
  _Partial coverage creates blind spots._

### Environmental

- **GOVERNANCE**: This policy covers all instances of production without exception.
  _Partial coverage creates blind spots._

- **DEPENDENCY**: The production system exists and is operational.
  _Policy references a system that must be running and maintained._

- **CONFIGURATION**: The production system is properly configured and maintained.
  _System configuration is assumed correct without explicit verification._

- **GOVERNANCE**: This policy covers all instances of production without exception.
  _Partial coverage creates blind spots._

- **IDENTITY**: Team membership is accurately reflected in access control groups.
  _Stale group membership creates unauthorized access._

- **DEPENDENCY**: The production system exists and is operational.
  _Policy references a system that must be running and maintained._

- **CONFIGURATION**: The production system is properly configured and maintained.
  _System configuration is assumed correct without explicit verification._

- **GOVERNANCE**: This policy covers all instances of production without exception.
  _Partial coverage creates blind spots._

- **GOVERNANCE**: This policy covers all instances of production without exception.
  _Partial coverage creates blind spots._

- **NETWORK**: No undocumented direct internet paths exist to protected resources.
  _Shadow IT and undocumented connections bypass network controls._


## Cross-tabulation: ASF Type × Ontology

| Type | Explicit | Implicit | Derived | Trust | Operational | Dependency | Architectural | Environmental | Total |
|------|----------|----------|---------|-------|-------------|------------|---------------|---------------|-------|
| ACCESS | 20 | 2 | 0 | 70 | 0 | 50 | 0 | 0 | 142 |
| IDENTITY | 15 | 2 | 0 | 69 | 0 | 52 | 0 | 2 | 140 |
| NETWORK | 15 | 1 | 15 | 0 | 0 | 1 | 4 | 2 | 38 |
| CONFIGURATION | 15 | 140 | 217 | 2 | 101 | 2 | 1 | 6 | 484 |
| PROCESS | 15 | 4 | 4 | 1 | 234 | 2 | 0 | 0 | 260 |
| GOVERNANCE | 10 | 100 | 241 | 0 | 4 | 102 | 1 | 5 | 463 |
| DOCUMENTATION | 5 | 6 | 105 | 0 | 0 | 5 | 0 | 0 | 121 |
| DEPENDENCY | 5 | 27 | 0 | 0 | 0 | 11 | 1 | 5 | 49 |
