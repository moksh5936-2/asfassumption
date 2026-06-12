# Healthcare PHI Architecture

## Metadata

- **Name:** Healthcare PHI Architecture
- **Version:** 1.0
- **Purpose:** Security assessment of cloud-native healthcare data platform
- **Compliance:** HIPAA, SOC2, ISO27001

## System

**HealthData Platform** — Cloud-native platform for managing protected health information (PHI).

## Components

| Component | Type | Description |
|-----------|------|-------------|
| Auth0 | identity_provider | OAuth2/OIDC identity provider for authentication |
| WebApp | web_application | React-based healthcare provider portal |
| APIGateway | api_gateway | AWS API Gateway for request routing and throttling |
| PHIDatabase | database | PostgreSQL database storing patient health records |
| KMS | encryption_service | AWS KMS for encryption key management |
| AuditLog | logging_service | Immutable audit log for PHI access events |
| BackupService | storage_service | Encrypted backup service for PHI data |
| ThirdPartyAnalytics | external_service | Third-party analytics provider with PHI access |
| AdminConsole | admin_tool | Internal admin console for system management |

## Relationships

```
WebApp --OAuth2--> Auth0
WebApp --HTTPS--> APIGateway
APIGateway --TLS--> PHIDatabase
APIGateway --TLS--> KMS
PHIDatabase --TLS--> BackupService
PHIDatabase --TLS--> AuditLog
AdminConsole --TLS--> PHIDatabase
AdminConsole --HTTPS--> Auth0
ThirdPartyAnalytics --TLS--> PHIDatabase
```

## Assumptions

1. MFA is enforced for all Auth0 user authentication.
2. Auth0 administrative access is restricted to authorized admins only.
3. All API requests pass through APIGateway for authentication validation.
4. PHI data is encrypted at rest in PHIDatabase.
5. PHI data is encrypted in transit between all components.
6. Encryption keys are stored and managed in KMS with automatic rotation.
7. KMS access is restricted to authorized services only.
8. Audit logging is immutable and tamper-proof.
9. All PHI access events are logged with user and timestamp.
10. Backup data is encrypted at rest and in transit.
11. Backup restore procedures are tested regularly.
12. ThirdPartyAnalytics has access only to de-identified PHI.
13. Third-party provider maintains equivalent security controls.
14. AdminConsole requires MFA for all administrative access.
15. Object-level authorization is enforced for PHI record access.
16. Database connection pooling does not leak data between sessions.
17. TLS certificates are monitored and renewed before expiry.
18. API rate limiting prevents abuse of PHI endpoints.
19. Network segmentation isolates the PHI database in a private subnet.
20. Unauthorized PHI export is detected and alerted.
21. Session tokens expire and are rotated periodically.
22. Auth0 tenant is configured with breach detection and anomaly alerts.
23. Audit log storage has sufficient retention for compliance.
24. PHI data minimization policies are enforced at application layer.
25. Incident response plan includes PHI breach notification.
26. Vendor risk assessments are conducted for ThirdPartyAnalytics.
27. Database backups are stored in a separate geographic region.
28. KMS key deletion is protected with multi-factor authorization.
29. API Gateway logs are monitored for anomalous access patterns.
30. System health and availability monitoring covers all components.

## Security Controls

### Authentication
- MFA
- Password_Policy
- Session_Management
- SSO_Integration
- Breach_Detection

### Authorization
- RBAC
- Object_Level_Access_Control
- API_Gateway_Validation
- Least_Privilege

### Encryption
- KMS_Key_Management
- TLS_1.3
- AES256_At_Rest
- Key_Rotation
- Key_Deletion_Protection

### Logging
- Immutable_Audit_Log
- PHI_Access_Logging
- API_Request_Logging
- Anomaly_Detection

### Backup
- Encrypted_Backups
- Cross_Region_Replication
- Restore_Testing

### Network
- Private_Subnet
- Network_Segmentation
- TLS_Termination
- Rate_Limiting

### Monitoring
- Health_Checks
- Availability_Monitoring
- Unauthorized_Access_Detection

### Third Party
- Vendor_Risk_Assessment
- Contractual_Security_Controls
- De_identified_Access_Only

## Expected Results

- **Minimum Assumptions:** 25
- **Minimum Critical:** 3
- **Minimum High:** 8
- **Expected STRIDE Categories:** Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, Elevation of Privilege
- **Expected High-Risk Themes:**
  - PHI data exposure
  - Unauthorized access
  - Key compromise
  - Third-party risk
  - Audit integrity

## Validation Criteria

- All PHI access must be authenticated and authorized.
- All data in transit must use TLS 1.3.
- All data at rest must use AES-256 encryption.
- Encryption keys must be managed by KMS.
- Audit logs must be immutable with tamper detection.
- Backup must be encrypted and geographically redundant.
- Third-party access must be restricted to de-identified data.
- Administrative access must require MFA.

## Notes

- This architecture supports regulated healthcare workloads.
- Formal HIPAA risk assessment is required before production.
- SOC2 Type II audit is in progress.
- ISO27001 certification is targeted for Q3 2026.
