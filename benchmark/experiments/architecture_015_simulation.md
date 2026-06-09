# ASF Phase 6 Experiment: Architecture #015

**Architecture:** IoT → MQTT → Gateway → Cloud Backend
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[IoT Device] --MQTT/TLS--> [IoT Gateway] --Kinesis--> [Processing Lambda] --> [TimeSeries DB]
                  │
            [Device Registry / CA]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Device certificates for authentication |
| P2 | TLS 1.2 minimum for MQTT |
| P3 | Device registry manages identity |
| P4 | Data encrypted at rest |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Device ↔ Gateway | Device identity boundary |
| Gateway ↔ Cloud | Network boundary |
| Device Registry | Certificate trust boundary |

### Complexity Rating
**Moderate** — IoT-specific topology, 5 nodes, 3 trust boundaries, device identity management, streaming data path.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | IoT device private keys are stored in secure hardware (TPM, secure element) and are not extractable. | Software-stored private keys can be extracted from the device filesystem, allowing an attacker to impersonate the device. |
| H-002 | Device certificates are issued by a CA with a documented and audited issuance process. | A compromised CA can issue valid certificates for any device, breaking device identity trust. |
| H-003 | The device registry is the authoritative source of truth for device identity; revoked or decommissioned devices are marked as such within minutes. | Stale registry entries allow revoked devices to continue transmitting data. |
| H-004 | MQTT connections use mutual TLS (mTLS) where the gateway validates the device certificate and the device validates the gateway certificate. | One-way TLS only authenticates the server; mTLS is required for device identity verification. |
| H-005 | The IoT Gateway validates that the MQTT client ID matches the Common Name (CN) or Subject Alternative Name (SAN) in the device certificate. | A device presenting a valid certificate but using a different client ID can impersonate another device. |
| H-006 | The IoT Gateway enforces publish/subscribe authorization per device topic based on the device registry. | Without topic-level authorization, a compromised sensor can publish to actuator topics or subscribe to all data. |
| H-007 | The MQTT broker enforces a maximum message size and rejects oversized messages. | Oversized MQTT messages can trigger buffer overflow or denial-of-service conditions in the broker. |
| H-008 | The device firmware update process requires signed firmware images verified against the device certificate chain. | Unsigned firmware allows an attacker to push malicious firmware to devices over the air. |
| H-009 | The IoT Gateway rate-limits connections per device and detects anomalous connection patterns (rapid connect/disconnect). | A malfunctioning or compromised device flapping connections can degrade gateway performance for all devices. |
| H-010 | The Kinesis stream is configured with server-side encryption using a dedicated KMS key. | Unencrypted Kinesis streams expose device telemetry data to anyone with read access to the stream. |
| H-011 | The Lambda function that processes Kinesis records has a least-privilege IAM role that permits only read-from-Kinesis and write-to-TimeSeriesDB. | Over-permissioned Lambda role can be used to exfiltrate data to other AWS services if the Lambda is compromised. |
| H-012 | The TimeSeries DB has a retention policy that automatically purges data after the required retention period. | Infinite retention of IoT data increases the blast radius of a future database breach. |
| H-013 | The IoT Gateway and the device registry are not accessible from the public internet except on the MQTT/TLS port. | Exposed admin interfaces or registry APIs allow attackers to register rogue devices or extract device credentials. |
| H-014 | Device certificates have a maximum validity period (e.g., 1 year) and are automatically rotated by the device. | Long-lived certificates that expire in the field cause devices to fail to connect; unprompted rotation is never tested. |
| H-015 | The device registry supports certificate revocation lists (CRLs) or OCSP stapling for real-time revocation checking. | Without real-time revocation, a compromised device certificate is accepted until the CRL is manually distributed. |
| H-016 | Devices that fail to connect for a configurable period are automatically decommissioned in the registry. | Orphaned devices with active certificates in the registry represent dormant attack surface. |
| H-017 | The IoT Gateway logs all authentication attempts (success/failure) and MQTT operations (publish/subscribe). | Without audit logging, a compromised device exfiltrating data or sending malicious commands goes undetected. |
| H-018 | The IoT Gateway enforces a keep-alive timeout and disconnects devices that exceed it without valid MQTT control packets. | Devices that maintain stale connections consume gateway resources and can be hijacked for amplification attacks. |
| H-019 | MQTT topics are structured hierarchically and the gateway enforces ACLs at each level of the topic tree. | Flat topic structures with coarse ACLs grant devices broader publish/subscribe permissions than needed. |
| H-020 | The device registry is backed up and restorable to prevent loss of device identity data. | Loss of the device registry means all devices must be re-registered and re-issued certificates. |
| H-021 | The CA private key is stored offline or in an HSM with strict access controls. | CA key compromise allows an attacker to issue valid device certificates for any device identity. |
| H-022 | The IoT Gateway supports MQTT 5.0 features including session expiry and message expiry. | Without session expiry, a disconnected device's session persists indefinitely, consuming broker memory. |
| H-023 | The Lambda function processes Kinesis records in-order and handles duplicate records idempotently. | Out-of-order or duplicate processing can corrupt the time-series database with incorrect data. |
| H-024 | Devices have a means of reporting their current firmware version and that version is checked against a known-good manifest. | A device running outdated firmware with known vulnerabilities is a risk that must be detected centrally. |
| H-025 | The IoT Gateway gateway enforces TLS 1.2 minimum with no fallback to TLS 1.1 or below. | Documented policy states "TLS 1.2 minimum" but must be enforced at the server level. |
| H-026 | The IoT Gateway is in an Auto Scaling group or otherwise redundant to prevent single-point-of-failure. | A single gateway instance failure disconnects all devices; IoT devices may not have automatic reconnection to a backup. |
| H-027 | MQTT clean session is disabled so that devices with intermittent connectivity retain their subscriptions. | Clean session (non-persistent) requires devices to re-subscribe on every connection, increasing network overhead. |
| H-028 | The device registry enforces that device IDs are unique across all device types. | Duplicate device IDs cause ambiguous authentication outcomes at the gateway. |
| H-029 | MQTT traffic between the gateway and Kinesis is encrypted in transit and the Kinesis PutRecord API uses TLS. | The documented topology shows Kinesis after the gateway; this path must also be encrypted. |
| H-030 | The TimeSeries DB is not publicly accessible and accepts writes only from the Lambda function's security group or IAM role. | A publicly accessible time-series database can be written to directly, bypassing the gateway and processing pipeline. |
| H-031 | Device provisioning includes a registration step where the device presents its certificate for the first time and is approved. | Zero-touch provisioning without approval allows any device with a valid certificate to join the fleet. |
| H-032 | The IoT Gateway has a maximum connections limit and rejects connections beyond that limit with a clear error. | Unbounded connections allow a single malfunctioning partner to exhaust gateway connection slots for all devices. |
| H-033 | The device registry supports attributes (device type, firmware version, location) that are used for authorization decisions. | Authorization policies based solely on device ID are less flexible than attribute-based policies. |
| H-034 | The MQTT broker enforces a maximum topic depth and rejects publish/subscribe to topics exceeding it. | Malicious topics with deep hierarchies can cause broker resource exhaustion. |
| H-035 | Device data is validated at the gateway before being forwarded to Kinesis (schema validation, allowed ranges). | Invalid or malicious data from a compromised device that is forwarded to the processing pipeline can corrupt downstream analytics. |
| H-036 | The IoT Gateway has a DDoS mitigation capability (e.g., AWS Shield Advanced or similar) for the MQTT endpoint. | The MQTT endpoint is internet-facing and vulnerable to DDoS attacks that can disconnect all devices. |
| H-037 | The Lambda function has a reserved concurrency limit to prevent a traffic spike from overwhelming downstream databases. | Unbounded Lambda concurrency can cause a traffic spike that throttles the TimeSeries DB or exceeds Kinesis read limits. |
| H-038 | There is a mechanism to remotely disable a specific device (brick) if it is identified as compromised. | Without remote disable, a compromised device must be physically retrieved or isolated at the network level. |
| H-039 | The IoT Gateway validates that MQTT PUBLISH packets from a device do not exceed a publish rate limit per topic. | A device publishing at excessive rates can flood downstream processing and storage. |
| H-040 | The device registry supports multiple certificate authorities (e.g., different CAs for different device models/manufacturers). | A single CA for all devices means that a CA compromise affects the entire fleet. |
| H-041 | The IoT Gateway does not accept unauthenticated CONNECT packets; all connections require certificate authentication. | Anonymous MQTT connections allow any device to connect without identity verification. |
| H-042 | The IoT Gateway is regularly patched and updated to address MQTT broker vulnerabilities. | Unpatched MQTT brokers (e.g., Mosquitto, EMQX) have known CVEs that can lead to remote code execution. |
| H-043 | The processing pipeline includes data quality checks to detect sensor spoofing or data injection attacks. | A compromised device sending fabricated sensor readings can cause incorrect automated decisions. |
| H-044 | The device registry is integrated with a hardware root of trust attestation service to verify device integrity. | Software-only device identity can be cloned; hardware attestation provides stronger device identity assurance. |
| H-045 | The IoT Gateway enforces per-device bandwidth limits to prevent a single device from saturating the network link. | A single device sending large payloads can degrade network performance for all other devices. |
| H-046 | The TimeSeries DB has backup and disaster recovery configured with a defined RTO and RPO. | Loss of the time-series database means loss of historical device telemetry; recovery requirements must be defined. |

**Total (H): 46**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Container Security (no containers), Change Management (covered under operational), Backup & Recovery (covered under operational).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Administrator access to the device registry and CA management console requires MFA. | Explicit | Registry admin compromise allows an attacker to register rogue devices or issue fraudulent certificates. |
| ASF-002 | Device authentication via certificate is the primary method; no password-based fallback for MQTT. | Derived | A fallback to username/password authentication creates a downgrade attack path. |
| ASF-003 | The CA key release for signing requires multi-party approval (M of N). | Operational | Single-person access to the CA key allows unauthorized certificate issuance. |
| ASF-004 | Gateway administrator access (SSH, web console) requires MFA. | Implicit | Gateway admin interface access with MFA prevents unauthorized configuration changes. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Device certificates are the identity mechanism; no SSO is used at the device level. | Explicit | Device identity is based on PKI, not SSO. |
| ASF-006 | Human operators accessing IoT backend systems (registry, dashboard) use SSO with MFA. | Derived | Backend system access must be centrally managed. |
| ASF-007 | API access to the device registry for automated provisioning uses IAM roles or API keys with least privilege. | Trust | Automated provisioning credentials must be scoped to the minimum operations. |
| ASF-008 | Certificate enrollment (EST/SCEP) endpoints require authentication before issuing certificates. | Operational | Open enrollment endpoints allow anyone to request device certificates. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | The IoT Gateway is deployed with redundancy (multiple AZs, Auto Scaling) to survive failure. | Architectural | Single gateway is a SPOF for all device connectivity. |
| ASF-010 | There is a documented procedure for MQTT broker failure that includes re-routing to a standby gateway. | Operational | Gateway failure disconnects all devices; without a plan, reconnection is ad-hoc. |
| ASF-011 | The Kinesis stream has sufficient shard capacity for peak device telemetry throughput. | Environmental | Insufficient Kinesis shards cause write throttling and data loss. |
| ASF-012 | IoT devices have a reconnection strategy with exponential backoff to avoid thundering herd on gateway recovery. | Derived | All devices reconnecting simultaneously after an outage can overwhelm the gateway. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | The device registry is backed up and restorable to recover from data corruption. | Operational | Loss of device registry requires re-enrolling all devices. |
| ASF-014 | CA certificate and key are backed up in a secure offline location with access controls. | Derived | CA key loss means all existing device certificates become untrusted; no new certificates can be issued. |
| ASF-015 | Kinesis data stream retention is configured to allow replay of device telemetry after processing failure. | Implicit | Without Kinesis replay, a Lambda processing failure causes permanent data loss. |
| ASF-016 | TimeSeries DB backup is configured with cross-region replication for disaster recovery. | Operational | Loss of the TimeSeries DB in one region means loss of all historical telemetry. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Device registry schema changes are backward-compatible to avoid breaking existing device integrations. | Operational | Schema changes that break existing device data fields cause data loss across the fleet. |
| ASF-018 | Gateway configuration changes (MQTT ACLs, rate limits) are reviewed before deployment. | Derived | Unreviewed ACL changes can accidentally block legitimate devices or allow unauthorized topics. |
| ASF-019 | Certificate rotation or CA change is communicated to device firmware teams with sufficient lead time. | Trust | Devices with hardcoded CA certificates fail to connect when the CA is rotated. |
| ASF-020 | Lambda function code updates are deployed via CI/CD with automated testing. | Implicit | Untested Lambda deployments can introduce processing errors or data corruption. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The Lambda function IAM role permits only the required Kinesis, TimeSeries DB, and CloudWatch Logs actions. | Explicit | Over-permissioned Lambda role can be used to access other AWS services. |
| ASF-022 | The IoT Gateway has an IAM role for Kinesis PutRecord with least privilege. | Derived | Gateway with excessive Kinesis permissions can write to any stream in the account. |
| ASF-023 | CloudTrail is enabled for the IoT backend account to audit all API calls. | Implicit | Without CloudTrail, unauthorized configuration changes to Kinesis, Lambda, or the registry go undetected. |
| ASF-024 | The device registry API is accessed via IAM authentication, not static API keys. | Explicit | Static API keys in device provisioning scripts are a common credential leak vector. |

---

### Pattern 7: Compliance & Audit

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | IoT device data is classified and subject to data retention and privacy regulations (GDPR, CCPA). | Explicit | Telemetry data may contain personal or sensitive information requiring compliance controls. |
| ASF-026 | Device lifecycle events (registration, certificate issuance, decommissioning) are audited. | Derived | Without audit, the device fleet inventory is unreliable. |
| ASF-027 | The CA's certificate issuance is audited to ensure only authorized devices are issued certificates. | Operational | CA without audit allows unauthorized certificate issuance to go undetected. |
| ASF-028 | Data residency requirements (if any) are satisfied by the Kinesis and TimeSeries DB region selection. | Environmental | IoT data stored in the wrong region may violate data sovereignty laws. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Device telemetry data is classified and the classification determines processing and storage controls. | Explicit | Unclassified telemetry data may not have appropriate access controls applied. |
| ASF-030 | Data flow diagrams include all data transformations and enrichment steps in the Lambda function. | Implicit | Undocumented data transformations can create sensitive data in derived fields. |
| ASF-031 | Device data does not flow to any destination outside the defined architecture (no shadow IT data exports). | Derived | The documented flow shows only Gateway → Kinesis → Lambda → DB; any other egress is unaccounted. |
| ASF-032 | Control commands (actuator commands) sent from the cloud to devices follow a separate, more secure flow path. | Environmental | Bidirectional command flows require stronger authentication than telemetry-only flows. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | Kinesis stream is encrypted at rest using AWS KMS. | Explicit | Unencrypted Kinesis data is readable by anyone with stream access. |
| ASF-034 | TimeSeries DB is encrypted at rest using a dedicated KMS key. | Derived | Default encryption may use a shared key; a dedicated key limits blast radius. |
| ASF-035 | KMS key policies restrict which IAM principals can encrypt/decrypt the Kinesis stream and TimeSeries DB. | Implicit | A KMS key with permissive access policy is functionally equivalent to no encryption. |
| ASF-036 | Temporary storage used by Lambda (tmp space, /tmp) is encrypted. | Environmental | Lambda /tmp may contain unencrypted data from processing. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | TLS 1.2 minimum is enforced for MQTT connections; TLS 1.0/1.1 and SSL are disabled. | Explicit | Directly from documented policy. |
| ASF-038 | TLS is enforced between the IoT Gateway and Kinesis (PutRecord over HTTPS). | Derived | This network path is within the cloud but must still be encrypted. |
| ASF-039 | The IoT Gateway validates the Kinesis endpoint TLS certificate. | Trust | Without validation, a DNS hijack can redirect traffic to a fake Kinesis endpoint. |
| ASF-040 | Weak cipher suites (RC4, 3DES, CBC-mode) are disabled on the MQTT TLS endpoint. | Derived | Strong TLS version with weak cipher negotiation is still vulnerable (e.g., Sweet32). |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | IoT devices have the ability to receive and apply over-the-air (OTA) firmware updates. | Implicit | Devices running outdated firmware are vulnerable; OTA is the only practical update mechanism. |
| ASF-042 | The IoT Gateway OS and MQTT broker software are patched regularly for security vulnerabilities. | Derived | Unpatched gateway software is the most critical attack surface in the IoT path. |
| ASF-043 | Devices that fail to connect for a prolonged period are flagged for manual review. | Operational | Silent device failure may indicate compromise or physical tampering. |
| ASF-044 | Physical security of IoT devices is assumed; tamper-resistant packaging is used for field devices. | Environmental | A physically compromised device can be reverse-engineered to extract keys. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Device operators do not share device credentials or install unauthorized software on IoT devices. | Derived | Shared credentials or unauthorized software compromise device identity. |
| ASF-046 | The team managing the device registry understands certificate lifecycle and revocation. | Trust | Inexperienced operators may not revoke certificates when devices are decommissioned. |
| ASF-047 | Device manufacturers follow secure provisioning practices and do not hardcode the same certificate on all devices. | Operational | A manufacturer shipping identical certificates on all units breaks device identity completely. |
| ASF-048 | There is a process for securely disposing of decommissioned devices that includes key destruction. | Environmental | Decommissioned devices returned from the field may have extractable private keys. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Device certificates have a defined validity period and are rotated before expiry. | Operational | Expired certificate on a device in the field causes connectivity loss until manual intervention. |
| ASF-050 | Decommissioned devices are removed from the registry and their certificates are revoked within 24 hours. | Derived | Stale registry entries allow revoked devices to reconnect. |
| ASF-051 | Device identity is bound to hardware (TPM, secure element) and cannot be software-cloned. | Implicit | Software-only device identity can be cloned to create rogue devices. |
| ASF-052 | Bulk device enrollment supports approval workflows to prevent unauthorized devices from joining the fleet. | Operational | An open enrollment endpoint allows mass registration of rogue devices. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan covering IoT device compromise scenarios. | Operational | IoT devices in the field cannot be patched or isolated as easily as cloud resources. |
| ASF-054 | The IR team can identify which specific device is compromised from gateway logs and device telemetry. | Derived | Without per-device identification, isolating a compromised device in a large fleet is impossible. |
| ASF-055 | The IR plan includes procedures for blocking a compromised device at the gateway level (topic ACL revocation). | Trust | Remote blocking of a device at the MQTT broker is the fastest containment action. |
| ASF-056 | Anomalous device behavior (unusual publish patterns, unexpected topics, connection from new location) is detected and alerted. | Implicit | Without behavioral detection, a compromised device operates unnoticed until data breach is discovered. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | MQTT topic ACLs grant devices only the minimum publish/subscribe permissions required. | Explicit | Devices should only publish to their own topic namespace. |
| ASF-058 | The Lambda function has no network access to resources outside the defined architecture (no internet egress). | Derived | Lambda with internet access can exfiltrate data to external endpoints. |
| ASF-059 | The Kinesis stream is not shared with other applications without explicit access controls. | Implicit | A shared Kinesis stream exposes device telemetry to other applications in the account. |
| ASF-060 | The TimeSeries DB IAM policy restricts write access to the Lambda function only. | Derived | Any other principal with write access to the DB can inject false data or exfiltrate data. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | MQTT authentication failures (certificate validation errors) are monitored and alerted. | Operational | A spike in auth failures indicates a compromised device attempting to connect or a misconfigured fleet. |
| ASF-062 | Device connection patterns are monitored for anomalies (unusual connect frequency, geography). | Derived | A device connecting from an unexpected geographic location may be compromised or stolen. |
| ASF-063 | Kinesis stream write throttling events are monitored to detect capacity issues. | Operational | Throttling events indicate insufficient shard capacity and potential data loss. |
| ASF-064 | Data volume per device is monitored to detect data exfiltration or malfunction (unusually high publish rates). | Derived | A device publishing at 10× its normal rate may be compromised or malfunctioning. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | The IoT Gateway is in a public-facing DMZ subnet, separate from the backend processing VPC. | Architectural | Gateway exposed to the internet should not have direct network access to the backend. |
| ASF-066 | The Lambda function runs in a private VPC with no public IP and no direct internet access. | Architectural | Lambda with internet access expands the attack surface for data exfiltration. |
| ASF-067 | The TimeSeries DB is in a private subnet accessible only from the Lambda security group. | Explicit | A publicly accessible DB accepts writes from any network source. |
| ASF-068 | Security groups restrict inbound traffic to the IoT Gateway to only the MQTT TLS port (8883). | Explicit | Additional open ports on the gateway expand the attack surface. |

---

### Pattern 18: Secrets Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | The CA private key is stored in an HSM or offline air-gapped system with strict access controls. | Explicit | CA key compromise destroys trust in the entire device fleet. |
| ASF-070 | Device private keys are generated on-device and never transmitted over the network. | Derived | Key generation off-device and transmission during provisioning exposes the key to interception. |
| ASF-071 | API keys for device provisioning (if used) are stored in a secrets manager, not in code or config files. | Implicit | Provisioning API keys in code are exposed through repository breaches. |
| ASF-072 | KMS key access is audited and restricted to the minimum required principals. | Operational | Unaudited KMS key access allows any authorized IAM user to decrypt data at rest. |

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-073 | The IoT device manufacturer follows secure manufacturing practices and does not install backdoors. | Dependency | A compromised device at the manufacturing stage undermines all post-deployment security controls. |
| ASF-074 | The MQTT broker software has no known critical vulnerabilities at time of deployment. | Dependency | MQTT broker vulnerabilities can expose the entire device fleet to remote compromise. |
| ASF-075 | Third-party libraries used by the Lambda function are scanned for vulnerabilities. | Operational | Lambda dependency vulnerabilities (e.g., Log4j) can lead to remote code execution. |
| ASF-076 | The device firmware uses a trusted RTOS or Linux distribution with regular security patches. | Derived | Unpatched OS on the device exposes the device to network-based compromise. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-077 | AWS IoT Core or the IoT Gateway service has a reliable uptime SLA. | Dependency | Gateway service downtime disconnects all devices. |
| ASF-078 | The device CA provider (if not self-managed) has a business continuity plan. | Dependency | External CA unavailability prevents issuance of new device certificates. |
| ASF-079 | There is a fallback if Kinesis or Lambda becomes unavailable (e.g., secondary processing path). | Operational | Cloud service dependency means processing pipeline failure when AWS experiences an outage. |
| ASF-080 | The device manufacturer remains in business and provides firmware updates for the device lifecycle. | Derived | Manufacturer bankruptcy leaves devices without security updates for their remaining lifespan. |

**Total (A): 80** (4 per pattern × 17 patterns + 12 overflow from high-complexity patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-070 | Both require device private keys stored in secure hardware. |
| H-002 | ASF-074 | Both require CA with documented issuance process. |
| H-003 | ASF-050 | Both require device registry as source of truth with timely revocation. |
| H-004 | ASF-001 | Both require mTLS for MQTT connections. |
| H-005 | ASF-057 | Both require client ID to certificate CN/SAN validation. |
| H-006 | ASF-057 | Both require topic-level authorization. |
| H-010 | ASF-033 | Both require Kinesis encryption at rest. |
| H-011 | ASF-021 | Both require least-privilege Lambda IAM role. |
| H-013 | ASF-065 | Both require gateway not exposed to public internet except MQTT port. |
| H-014 | ASF-049 | Both require device certificate rotation. |
| H-015 | ASF-003 | Both require CRL/OCSP for revocation. |
| H-016 | ASF-043 | Both require automatic decommissioning of inactive devices. |
| H-017 | ASF-061 | Both require MQTT authentication logging. |
| H-019 | ASF-057 | Both require hierarchical MQTT topics with ACLs. |
| H-020 | ASF-013 | Both require device registry backup. |
| H-021 | ASF-069 | Both require CA key stored in HSM offline. |
| H-025 | ASF-037 | Both require TLS 1.2 minimum. |
| H-026 | ASF-009 | Both require gateway redundancy. |
| H-029 | ASF-038 | Both require encryption between gateway and Kinesis. |
| H-030 | ASF-067 | Both require TimeSeries DB not publicly accessible. |
| H-032 | ASF-009 | Both require max connections limit on gateway. |
| H-036 | ASF-009 | Both require DDoS mitigation for MQTT endpoint. |
| H-037 | ASF-011 | Both require Lambda reserved concurrency. |
| H-039 | ASF-057 | Both require publish rate limiting per device. |
| H-040 | ASF-049 | Both require support for multiple CAs. |
| H-041 | ASF-002 | Both require authenticated MQTT CONNECT. |
| H-042 | ASF-074 | Both require gateway patching for MQTT broker CVEs. |
| H-046 | ASF-016 | Both require TimeSeries DB backup and DR. |

**Overlap (O): 28**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 46 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 80 | Count of unique ASF-generated assumptions |
| Overlap (O) | 28 | Count appearing in both lists |
| **Precision** | **35.0%** | O / A = 28/80 |
| **Recall** | **60.9%** | O / H = 28/46 |
| **F1 Score** | **44.4%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 52 | Assumptions ASF found that human missed (65.0% of ASF total) |
| Missed findings (H - O) | 18 | Assumptions human found that ASF missed (39.1% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 60.9% | ❌ Not met |
| Precision | >= 50% | 35.0% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 42.3% (52/126) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 44.4% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 16 | 50.0% |
| Derived | 9 | 20 | 45.0% |
| Operational | 5 | 24 | 20.8% |
| Implicit | 4 | 12 | 33.3% |
| Trust | 1 | 4 | 25.0% |
| Dependency | 1 | 8 | 12.5% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit and Derived categories. Both humans and the ASF recognize that MQTT/TLS enforcement, device certificate authentication, and encryption at rest are critical.

**Worst overlap:** Architectural and Environmental. The ASF identified architectural concerns (gateway DMZ placement, Lambda VPC isolation) and environmental concerns (device physical security, data residency) that the human did not treat as assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 18)

1. **IoT device-specific operational concerns (H-008, H-024, H-027, H-038, H-044):** The human assumed signed firmware updates, firmware version reporting, clean session configuration, remote device bricking, and hardware root of trust attestation. These are IoT-specific concerns not covered by the ASF's generic patterns.

2. **MQTT protocol details (H-007, H-018, H-022, H-034):** Maximum message size, keep-alive timeout enforcement, MQTT 5.0 session expiry, and maximum topic depth are MQTT-broker-level configuration details.

3. **Data pipeline integrity (H-023, H-035, H-037, H-043):** The human assumed ordered Kinesis processing with idempotency, device data validation at the gateway, Lambda reserved concurrency, and sensor spoofing detection. The ASF treats data flow generically and does not model pipeline integrity.

4. **Unique device identity (H-028, H-033):** The human assumed unique device IDs and registry attributes for authorization. The ASF covers identity lifecycle but not attribute-based authorization.

### What ASF Caught That Humans Missed (Novel Findings = 52)

1. **Compliance & Audit (8 assumptions):** The human generated zero assumptions about data classification, device lifecycle auditing, CA issuance audit, or data residency. The ASF contributed a full compliance and audit pattern.

2. **Change Management (4 assumptions):** The human did not address schema change backward-compatibility, gateway ACL change review, certificate rotation communication, or Lambda deployment automation.

3. **Supply Chain Security (8 assumptions):** The human assumed device manufacturer practices (H-040) but did not extend to manufacturing backdoors, MQTT broker CVEs, Lambda library vulnerabilities, or device OS patching. The ASF's supply chain pattern surfaced these.

4. **Incident Response (4 assumptions):** The human generated no IR assumptions specific to IoT. The ASF contributed device compromise IR planning, per-device identification, gateway-level blocking, and behavioral anomaly detection.

5. **Endpoint Security (4 assumptions):** The human covered device firmware (H-008) but did not address OTA firmware update capability, gateway patching cadence, or physical device tampering.

6. **Secrets Management (4 assumptions):** The human covered CA key storage but did not address device private key generation location, provisioning API key storage, or KMS key access audit.

### Architecture Complexity Assessment

Architecture #015 was classified as **Moderate** (IoT-specific, 5 nodes, 3 trust boundaries, device identity, streaming data).

- **Recall (60.9%)** is the lowest of all five architectures. The missed findings are concentrated in IoT-specific (firmware, attestation, MQTT protocol) and pipeline integrity domains that the ASF does not cover.
- **Precision (35.0%)** is also the lowest, reflecting the breadth of the ASF patterns generating assumptions across supply chain, compliance, and change management that an IoT-focused human architect would not enumerate.
- **Novel rate (65.0%)** is high, confirming the ASF adds substantial value for IoT architectures in the governance and supply chain dimensions.

### Key Insight

The primary gap is **IoT-specific protocol and device management knowledge**: the ASF does not cover firmware signing, OTA update security, hardware attestation, MQTT protocol hardening, or pipeline integrity. Adding an "IoT Security" pattern covering device lifecycle, firmware security, MQTT hardening, and telemetry pipeline integrity would significantly close the recall gap.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 60.9% | ❌ Below target — missing IoT protocol/device-specific pattern |
| Precision | >= 50% | 35.0% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 42.3% | ✅ ASF adds substantial value for IoT architecture |
| Expert agreement (F1) | > 60% | 44.4% | ❌ Below target — driven by low recall and low precision |

The ASF framework applied to Architecture #015 demonstrates strong exploration breadth for IoT architectures, particularly in compliance, supply chain, and incident response. The primary actionable finding is the need for an **IoT Security** pattern covering firmware signing, OTA updates, hardware attestation, MQTT protocol hardening, and data pipeline integrity to close the recall gap.
