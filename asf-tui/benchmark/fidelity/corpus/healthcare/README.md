# Healthcare Architecture Benchmark

## Architecture Description

A healthcare application that handles patient health information (PHI).

## Facts

### Security Controls
- MFA is enabled for all users
- Encryption is enabled for data at rest and in transit
- HIPAA compliance is required
- Backups are automated daily
- WAF is enabled
- Audit logging is enabled
- VPN is used for admin access
- Role-based access control is enforced

### Components
- Patient Database
- API Gateway
- Auth0 Service
- Audit Log
- Load Balancer
- CDN
- Message Queue

### Relationships
- API Gateway -> Patient Database (queries)
- Auth0 Service -> API Gateway (authenticates)
- Load Balancer -> API Gateway (routes)
- CDN -> API Gateway (caches)
- Admin -> Patient Database (manages)

## Expected Assumptions

### Hidden Assumptions (should be generated)
1. Certificates are rotated before expiry
2. Restore testing is performed periodically
3. Backup data is encrypted
4. WAF rules are reviewed regularly
5. Auth0 administrators are restricted
6. Session tokens have short expiration
7. PHI access is logged with immutable audit trails
8. Break-glass procedures are documented
9. Patient safety mechanisms are in place
10. Database replication and failover are configured
11. API has rate limiting and throttling
12. Load balancer has health checks
13. CDN cache invalidation is configured
14. Message queue has dead letter queue
15. Log retention policies meet compliance requirements

## Expected Contradictions

### None (all facts are consistent)

## Metrics

- Expected fidelity score: 90%+
- Expected assumption count: 15+
- Expected contradiction count: 0
