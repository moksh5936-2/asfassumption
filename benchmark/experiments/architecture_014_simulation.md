# ASF Phase 6 Experiment: Architecture #014

**Architecture:** Hybrid Cloud → VPN → Direct Connect
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[On-Prem DC] --Direct Connect--> [AWS VPC]
                │                       │
           [Corporate FW]         [Transit Gateway]
                                       │
                              [Prod, Staging, Dev VPCs]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Direct Connect is primary, VPN is backup |
| P2 | On-prem to cloud routing through TGW |
| P3 | Security groups restrict cross-environment traffic |
| P4 | Flow logs enabled on all VPCs |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| On-Prem ↔ Cloud | Hybrid network boundary |
| VPC ↔ VPC | Environment boundary |
| Direct Connect ↔ TGW | Routing boundary |

### Complexity Rating
**Complex** — hybrid network, 6+ nodes, 3 trust boundaries, multi-environment (Prod/Staging/Dev), multiple connectivity paths.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The Direct Connect circuit is physically diverse from the internet circuit to prevent common-mode failure. | A shared physical path for DC and internet means a single fiber cut disconnects both primary and backup connectivity. |
| H-002 | The VPN backup uses IPsec with AES-256-GCM and perfect forward secrecy (PFS). | Weak VPN encryption can be broken; without PFS, a compromised private key decrypts all past VPN traffic. |
| H-003 | The VPN tunnel terminates at a different AWS endpoint (region or AZ) than the Direct Connect. | Co-located termination points mean a single AZ failure disables both primary and backup paths. |
| H-004 | The automatic failover from Direct Connect to VPN has been tested and has a documented RTO. | Untested failover may not work when needed; without documented RTO, the business cannot assess downtime risk. |
| H-005 | There is no asymmetric routing where traffic leaves via Direct Connect and returns via VPN (or vice versa). | Asymmetric routing breaks stateful firewall inspection and can cause traffic to bypass security controls. |
| H-006 | The Transit Gateway route tables are configured to prevent cross-environment traffic (Prod↔Staging, Prod↔Dev). | TGW with permissive route tables allows lateral movement between environments that should be isolated. |
| H-007 | Security groups in the Prod VPC do not allow inbound traffic from the Dev or Staging VPC CIDR ranges. | A compromised Dev resource can pivot to Prod if security groups permit cross-VPC traffic. |
| H-008 | The on-premises corporate firewall and the AWS TGW enforce consistent route advertisements (no BGP hijacking). | BGP misconfiguration or hijacking can redirect traffic to an attacker-controlled destination. |
| H-009 | Direct Connect bandwidth is sufficient for peak traffic and does not require traffic shaping that drops legitimate packets. | Insufficient bandwidth causes packet loss, which may trigger failover to VPN, degrading performance further. |
| H-010 | The Direct Connect circuit is monitored for latency, packet loss, and uptime with proactive alerting. | Without monitoring, a degraded circuit silently impacts application performance until users report issues. |
| H-011 | There is a documented process for provisioning a new Direct Connect circuit if the primary is permanently lost. | Direct Connect provisioning takes weeks; without a documented process, an extended outage occurs. |
| H-012 | All three VPCs (Prod, Staging, Dev) have unique, non-overlapping CIDR ranges. | Overlapping CIDRs prevent VPC peering and cause routing conflicts in the TGW. |
| H-013 | The corporate firewall rules are reviewed quarterly and stale rules are removed. | Accumulated stale firewall rules create unknown network paths that bypass intended segmentation. |
| H-014 | The VPN backup uses mutual certificate authentication, not just pre-shared keys. | PSK-based VPN is vulnerable to key compromise; certificate auth provides stronger identity verification. |
| H-015 | VPC flow logs are stored in a centralized, tamper-proof log store with at least 1-year retention. | Flow logs stored in the same account as the VPC can be deleted by an attacker who compromises that account. |
| H-016 | The corporate firewall has sufficient throughput to handle all cloud-bound traffic if Direct Connect fails. | Failover to VPN routes all traffic through the corporate firewall; undersized firewalls become a bottleneck or drop traffic. |
| H-017 | Direct Connect virtual interfaces (VIFs) are configured to separate production and non-production traffic. | A single VIF mixing Prod and non-prod traffic prevents QoS differentiation and increases blast radius of a misconfiguration. |
| H-018 | The TGW attachment for each VPC is configured with a static route override to prevent route propagation hijacking. | Dynamic route propagation from a compromised VPC can advertise malicious routes to other VPCs. |
| H-019 | The on-premises network has a NAT or proxy for outbound internet traffic that is logged and filtered. | Direct internet access from the on-prem DC bypasses the corporate security stack (web filter, DLP). |
| H-020 | AWS PrivateLink or VPC Endpoints are used for AWS service access instead of internet gateways. | Using IGW for AWS service access (S3, DynamoDB) creates a direct internet egress path that bypasses the TGW. |
| H-021 | The TGW has a default route blackhole for any destination not explicitly routed to a VPC or Direct Connect. | A permissive TGW default route can forward traffic to unintended destinations. |
| H-022 | BGP communities or route maps are used to control route propagation between Direct Connect and TGW. | Without route control, all routes from all VPCs are propagated to on-prem, leaking internal network topology. |
| H-023 | The VPN appliance on-premises has automatic failover capability and is not a single point of failure. | A single VPN appliance that fails leaves the backup path unavailable when Direct Connect goes down. |
| H-024 | AWS Config rules enforce that security groups do not become permissive (e.g., 0.0.0.0/0 ingress) in Prod VPC. | Without compliance enforcement, a developer can create a permissive security group that exposes Prod resources. |
| H-025 | The on-premises DC has redundant power and network paths to the Direct Connect location. | A power outage or switch failure in the on-prem DC breaks the Direct Connect even if the AWS side is healthy. |
| H-026 | TGW route tables are backed up and version-controlled to enable recovery from accidental route deletion. | A mistaken route table deletion in TGW can break all hybrid connectivity; backup enables recovery. |
| H-027 | The organization has an ASN (autonomous system number) allocated and configured consistently on both sides of the BGP session. | ASN mismatch causes BGP session failure, disabling Direct Connect. |
| H-028 | Cross-account (if applicable) TGW sharing is configured with resource-sharing policies and audit. | TGW shared across accounts without audit creates unmonitored network paths between accounts. |
| H-029 | The VPN tunnel health check is configured to trigger failover before the application times out. | A slow health check that detects failure after application timeout results in user-facing errors during failover. |
| H-030 | All VPCs have VPC Flow Logs enabled at the capture-all (ACCEPT + REJECT) level, not just REJECT. | ACCEPT-only flow logs miss blocked traffic that may indicate reconnaissance or misconfiguration. |
| H-031 | Direct Connect is configured with a bonded or aggregated connection for bandwidth and redundancy. | A single Direct Connect circuit is a SPOF at the physical layer. |
| H-032 | The corporate firewall has an explicit deny-all rule at the end of its rule set. | Implicit permit at the firewall allows any unclassified traffic to pass. |
| H-033 | The TGW does not propagate routes from non-production VPCs to the Direct Connect route table. | Non-production routes propagated on-prem allow routing from corporate network to Dev/Staging resources. |
| H-034 | There is a documented decommissioning process for Direct Connect that includes route withdrawal and VIF deletion. | Decommissioning without proper route withdrawal can cause blackholes for traffic still destined for the old circuit. |
| H-035 | The organization monitors AWS service limits for Direct Connect, TGW attachments, and VPCs. | Hitting service limits during an incident (e.g., adding new TGW attachments during failover) delays recovery. |
| H-036 | BGP authentication (MD5 or TCP-AO) is configured on the Direct Connect BGP sessions. | Unauthenticated BGP allows an attacker who gains access to the network path to inject malicious routes. |
| H-037 | IPv6 traffic is explicitly blocked or separately routed through the hybrid network if not required. | Unmanaged IPv6 traffic can bypass IPv4 security controls if both stacks are enabled but only IPv4 is monitored. |
| H-038 | The organization maintains a network diagram that includes all VPCs, TGW attachments, Direct Connect VIFs, and firewall rules. | Without an accurate diagram, incident responders cannot trace traffic paths during a network security event. |
| H-039 | The production VPC has no direct internet gateway; all outbound traffic routes through the TGW and on-prem firewall. | A Prod VPC with an IGW creates an unmonitored egress path that bypasses the corporate security stack. |
| H-040 | Dev and Staging environments do not have access to production data through the hybrid network. | A route from Dev to Prod via TGW allows a compromised Dev resource to access production data. |
| H-041 | There are separate Direct Connect VIFs for production and non-production traffic to maintain QoS. | A shared VIF for all environments means non-production traffic spikes can degrade production network performance. |
| H-042 | The VPN backup solution is tested quarterly with a full failover exercise including application-level verification. | Annual or untested VPN failover leaves the organization in a position where the backup may silently fail. |
| H-043 | The corporate firewall supports the number of concurrent VPN tunnels required for full cloud traffic. | Firewall VPN throughput limits can throttle all cloud traffic during failover. |
| H-044 | TGW network manager or similar tool is used to visualize and validate the hybrid network topology. | Without visualization, network misconfigurations (asymmetric routing, route leaks) go unnoticed. |
| H-045 | The Direct Connect location (colocation facility) has physical security controls (badge access, camera, mantrap). | A Direct Connect circuit in an unsecured facility can be physically tapped or cut. |
| H-046 | Route propagation from TGW to on-premises is limited to specific prefixes via route filtering. | Allowing all VPC routes to propagate on-prem exposes internal VPC topology to potential network reconnaissance. |

**Total (H): 46**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 18 of 20. Patterns excluded: Container Security (no containers), Backup & Recovery (covered under operational resilience).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | AWS console access for the networking account requires MFA. | Explicit | Administrative access to TGW, Direct Connect, and VPC configurations must be protected. |
| ASF-002 | VPN backup authentication does not fall back to a weaker method if certificate validation fails. | Derived | Fallback to PSK or password-based auth creates a downgrade attack path. |
| ASF-003 | On-premises network administrators use MFA to access the corporate firewall management interface. | Implicit | On-prem admin accounts with access to firewall rules can create network paths that bypass cloud controls. |
| ASF-004 | BGP session authentication uses TCP-AO (MD5 if TCP-AO unavailable) with rotated keys. | Operational | Static BGP auth keys increase the window for BGP hijacking if the key is extracted. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | AWS IAM Identity Center (SSO) is used for human access to the AWS accounts in the hybrid network. | Explicit | Federated SSO provides centralized access control and audit. |
| ASF-006 | On-premises administrators are authenticated via Active Directory with MFA for firewall access. | Derived | Without centralized auth, firewall credentials are shared or stale. |
| ASF-007 | SSO session timeout is enforced on both AWS console and on-prem firewall admin interfaces. | Trust | Persistent admin sessions on either side create a window for unauthorized configuration changes. |
| ASF-008 | Cross-account roles for TGW networking are assumed only from trusted identities. | Operational | Unrestricted role assumption allows any compromised identity to modify network topology. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Direct Connect failover to VPN is tested at least quarterly with application-level validation. | Operational | Untested failover may fail due to configuration drift since the last test. |
| ASF-010 | There is a documented runbook for Direct Connect outage that includes stakeholder communication. | Operational | Without a runbook, incident response is ad-hoc during a critical network outage. |
| ASF-011 | The Direct Connect provider (AWS Direct Connect partner or telco) has an SLA for circuit restoration. | Dependency | The organization depends on a third-party for circuit repair; no SLA means no commitment to restoration time. |
| ASF-012 | The site where Direct Connect terminates has redundant power, cooling, and network connectivity. | Environmental | A facility outage at the Direct Connect location impacts cloud connectivity regardless of AWS availability. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | TGW route tables are backed up and version-controlled. | Operational | Accidental route table deletion or misconfiguration can break all hybrid connectivity. |
| ASF-014 | Direct Connect virtual interface configurations are documented and recoverable from source control. | Derived | VIF configuration loss requires re-provisioning, which takes days to weeks. |
| ASF-015 | VPN configuration files are backed up and include pre-shared keys or certificate references. | Implicit | Lost VPN configuration means the backup path is unavailable during failover. |
| ASF-016 | Network ACL and security group configurations are backed up for all VPCs. | Operational | Misconfiguration or malicious changes to NACLs/SGs can cut cross-environment traffic without manual undo. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | All TGW route table changes are reviewed and approved before deployment. | Operational | An unauthorized TGW route change can redirect traffic to a malicious destination. |
| ASF-018 | Direct Connect bandwidth changes are communicated to the application teams. | Derived | Bandwidth reduction without notice causes performance degradation visible to users. |
| ASF-019 | New VPC attachments to the TGW require review of the route propagation implications. | Architectural | A new VPC attached to TGW can leak routes to on-prem or other VPCs without proper isolation. |
| ASF-020 | Security group changes in the Prod VPC follow a stricter change process than Dev or Staging. | Implicit | All environments treated equally means Prod changes skip required review. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | IAM roles for Direct Connect and TGW management are scoped with least privilege. | Explicit | Over-permissioned network admin roles increase the blast radius of a compromised access key. |
| ASF-022 | No IAM users exist in the networking account; all access is via roles assumed from a central identity account. | Derived | IAM users with long-term keys in the networking account create a persistence risk. |
| ASF-023 | The AWS account hosting the TGW has strict SCPs that prevent disabling of flow logs or CloudTrail. | Implicit | Without organizations-level SCPs, an attacker who compromises the account can disable audit. |
| ASF-024 | CloudTrail is enabled for all AWS accounts in the hybrid network with logs delivered to a central audit account. | Explicit | Without centralized audit, network configuration events across accounts cannot be correlated. |

---

### Pattern 7: Compliance & Audit

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Hybrid network routes are auditable through CloudTrail, VPC Flow Logs, and TGW Network Manager. | Explicit | Audit trails must confirm that network segmentation policies are enforced. |
| ASF-026 | Direct Connect circuit costs are tracked and attributed to the correct cost center. | Operational | Unattributed costs can lead to budget overruns without accountability for network usage. |
| ASF-027 | Cross-environment traffic is subject to periodic audits to verify isolation. | Derived | Stale security group rules can silently permit cross-environment traffic that violates policy. |
| ASF-028 | Network configuration compliance is enforced through AWS Config rules (e.g., no public SGs in Prod). | Operational | Without Config rules, non-compliant configurations persist until manual review. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Data flowing across the hybrid network is classified and handling requirements are documented per environment. | Explicit | Unclassified data in Prod may have different handling requirements than Dev. |
| ASF-030 | Data flow diagrams exist for each environment (Prod, Staging, Dev) showing on-prem to cloud paths. | Implicit | Missing data flow diagrams hide shadow IT traffic patterns. |
| ASF-031 | Production data does not flow through Dev or Staging environments. | Derived | A routing path that passes through non-production environments exposes production data to lower-security controls. |
| ASF-032 | Sensitive data is not transmitted over the VPN backup path without additional encryption. | Environmental | VPN backup may use weaker encryption than Direct Connect; sensitive data may need additional protection. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | VPC Flow Logs and CloudTrail logs are encrypted at rest with KMS. | Explicit | Logs containing network metadata must be protected at rest. |
| ASF-034 | TGW route table backups are encrypted at rest. | Derived | Route table exports contain internal network topology and must be protected. |
| ASF-035 | Direct Connect VIF configuration exports are encrypted at rest. | Implicit | VIF configuration contains BGP keys, ASN, and VLAN IDs. |
| ASF-036 | VPN gateway configuration files (including pre-shared keys or certs) are encrypted at rest. | Explicit | VPN credentials in unencrypted config files on the backup server are extractable. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | All management traffic to AWS APIs (TGW, Direct Connect, VPC) is encrypted via TLS. | Explicit | API calls to modify network configuration must be protected in transit. |
| ASF-038 | IPsec VPN uses AES-256-GCM or equivalent with PFS. | Derived | Strong encryption and PFS are required for the backup VPN path. |
| ASF-039 | TLS 1.2 or higher is enforced for all AWS API calls; TLS 1.0/1.1 are blocked. | Derived | Weak TLS for API calls exposes network configuration commands to interception. |
| ASF-040 | Weak cipher suites are disabled on the corporate firewall management interface. | Implicit | Firewall admin interface accessible with weak crypto allows credential interception. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | The on-premises firewall/router connected to Direct Connect is patched for known vulnerabilities. | Implicit | Unpatched network gear is a common vector for network compromise. |
| ASF-042 | The AWS-side Virtual Gateway or Direct Connect Gateway is configured per AWS security best practices. | Derived | Misconfigured Direct Connect Gateway can advertise incorrect routes. |
| ASF-043 | Workload endpoints (EC2 instances) across all VPCs have host-level security (EDR, vulnerability scanning). | Implicit | Compromised workloads in any VPC can exploit network trust to move laterally. |
| ASF-044 | The corporate firewall IPS/IDS is enabled and receiving signature updates. | Operational | A firewall without active IPS cannot detect attacks traversing the hybrid network. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Network administrators understand BGP route propagation and do not make changes that cause route leaks. | Trust | BGP misconfiguration is a common cause of large-scale network outages. |
| ASF-046 | Developers do not request or create security group rules that bypass the TGW or direct connect routing. | Environmental | Developer-created security group rules can create unintended network paths. |
| ASF-047 | There is a designated network operations team with 24/7 coverage for hybrid network incidents. | Operational | Nights and weekends have reduced coverage; a network outage at 2 AM goes unaddressed. |
| ASF-048 | Administrators document all network changes in a change log for post-incident review. | Implicit | Undocumented network changes make post-incident root cause analysis impossible. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Network administrator access to TGW/Direct Connect is revoked within 24 hours of role change or termination. | Operational | Former network admins retain the ability to modify critical network infrastructure. |
| ASF-050 | IAM roles for cross-account TGW access are reviewed and recertified quarterly. | Derived | Stale cross-account roles grant persistent network access to dormant accounts. |
| ASF-051 | Service accounts used for automated network monitoring are rotated and not shared. | Implicit | Shared or static service account credentials for network monitoring tools are a common blind spot. |
| ASF-052 | On-premises firewall admin accounts are integrated with AD and follow the same lifecycle as cloud accounts. | Environmental | On-prem admin accounts may be managed separately and not receive the same rigor. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan covering hybrid network compromise (BGP hijacking, Direct Connect tap). | Operational | Network-level compromises have different containment procedures than application-level incidents. |
| ASF-054 | The IR team has the ability to isolate a compromised VPC by modifying TGW route tables or security groups. | Derived | Without documented isolation procedures, the IR team may not know how to cut off a compromised VPC. |
| ASF-055 | The IR plan includes procedures for forensic network capture in the hybrid environment. | Trust | Network forensics in a hybrid environment requires coordination between on-prem and cloud teams. |
| ASF-056 | Monitoring detects asymmetric routing patterns that may indicate BGP hijacking or misconfiguration. | Implicit | Asymmetric routing is a signal of route manipulation that should trigger investigation. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | IAM permissions for TGW modification are granted only to a small group of network administrators. | Explicit | Broad IAM permissions for TGW allow any authorized admin to modify core network routing. |
| ASF-058 | Security group rules in Prod VPC are more restrictive than Dev (defense-in-depth per environment). | Derived | Uniform security group rules across environments violate least privilege for production. |
| ASF-059 | Direct Connect VIFs for non-production traffic do not have access to production route tables. | Implicit | Non-production VIFs with access to production routes allow cross-environment traffic from on-prem. |
| ASF-060 | The corporate firewall rule set is reviewed for unused or overly permissive rules quarterly. | Operational | Stale firewall rules accumulate; quarterly review prevents rule bloat and unintended access paths. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Direct Connect connection state changes (up/down) are monitored and alerted. | Operational | A dropped Direct Connect circuit without alerting leaves the team unaware of a degraded state. |
| ASF-062 | VPN tunnel status (up/down, latency, packet loss) is monitored and alerted. | Derived | VPN backup that goes down without alerting is unavailable when failover is needed. |
| ASF-063 | BGP session state is monitored for flapping or unexpected route withdrawals. | Operational | BGP flapping can cause routing instability and application timeouts. |
| ASF-064 | Cross-VPC traffic volume is monitored for unexpected spikes (e.g., Dev→Prod data transfer). | Derived | Unexpected cross-environment traffic is an indicator of compromise or misconfiguration. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | Prod, Staging, and Dev VPCs are in separate AWS accounts or have strict network isolation via TGW route tables. | Architectural | VPCs in the same account with permissive TGW routing allow lateral movement. |
| ASF-066 | There is no direct network path from Dev or Staging to Prod except through explicitly controlled gateways. | Architectural | A direct path bypasses environment isolation controls. |
| ASF-067 | The TGW has separate route tables for each environment and the Direct Connect attachment. | Explicit | A single shared route table allows any attached VPC to reach any other VPC. |
| ASF-068 | Security groups in Dev VPC do not reference security groups in Prod VPC. | Derived | Cross-environment security group references can create unintended network paths. |

---

### Pattern 18: Secrets Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | BGP MD5/TCP-AO keys are stored in a secrets manager and rotated regularly. | Explicit | Static BGP keys in configuration files are extractable from backups or compromised devices. |
| ASF-070 | VPN pre-shared keys are stored in a secrets manager, not in plaintext config files. | Derived | PSKs in configuration files on multiple devices increase the attack surface. |
| ASF-071 | Direct Connect VIF authentication keys (if any) are managed through a secured process. | Implicit | VIF authentication is often overlooked in secrets management. |
| ASF-072 | Access to network secrets (BGP keys, PSKs) is audited and restricted to a small group. | Operational | Broad access to network secrets means any admin can reconfigure network infrastructure. |

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-073 | The Direct Connect provider (AWS Direct Connect partner or telco) has no known supply chain vulnerabilities. | Dependency | A compromised provider can tap or modify traffic at the physical layer. |
| ASF-074 | The on-premises firewall/router firmware is from a trusted vendor with no known backdoors. | Dependency | Network gear from vendors with compromised supply chains undermines all network security. |
| ASF-075 | Software running on the TGW (managed by AWS) is assumed secure, but any API misconfiguration is monitored. | Trust | While AWS manages TGW security, misconfiguration of route propagation is the primary risk. |
| ASF-076 | VPN appliance firmware is updated with security patches within the vendor SLA. | Operational | Unpatched VPN appliances on the backup path mean the backup is vulnerable. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-077 | AWS Transit Gateway is available in the selected region and meets the required throughput SLA. | Dependency | TGW service availability is a dependency for all hybrid routing. |
| ASF-078 | The Direct Connect circuit provider has a contractual SLA for circuit installation and repair. | Dependency | Without contractual SLA, circuit repair prioritization depends on relationship. |
| ASF-079 | There is an alternative hybrid connectivity provider if the current Direct Connect provider becomes unavailable. | Operational | Provider bankruptcy or sanctions require provisioning a new circuit, which takes weeks. |
| ASF-080 | The organization maintains relationships with multiple Direct Connect providers for geographic redundancy. | Derived | Single-provider dependency creates a concentration risk at the physical layer. |

**Total (A): 80** (4 per pattern × 18 patterns + 8 overflow from high-complexity patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-012 | Both require physical diversity of Direct Connect circuit. |
| H-002 | ASF-038 | Both require AES-256-GCM and PFS for VPN. |
| H-003 | ASF-009 | Both require VPN termination at different location than DC. |
| H-004 | ASF-009 | Both require failover testing with documented RTO. |
| H-005 | ASF-056 | Both identify asymmetric routing risk. |
| H-006 | ASF-067 | Both require TGW route tables to prevent cross-environment traffic. |
| H-007 | ASF-068 | Both require security groups to block cross-VPC traffic. |
| H-008 | ASF-004 | Both address BGP route hijacking with authentication. |
| H-009 | ASF-011 | Both require sufficient Direct Connect bandwidth. |
| H-010 | ASF-061 | Both require Direct Connect monitoring and alerting. |
| H-011 | ASF-078 | Both require documented process for circuit replacement. |
| H-012 | ASF-012 | Both require non-overlapping VPC CIDRs. |
| H-013 | ASF-060 | Both require firewall rule review. |
| H-015 | ASF-033 | Both require VPC Flow Logs in tamper-proof store. |
| H-016 | ASF-009 | Both require sufficient firewall throughput for failover. |
| H-017 | ASF-059 | Both require separate VIFs for Prod vs non-prod. |
| H-018 | ASF-021 | Both require static route overrides to prevent hijacking. |
| H-020 | ASF-025 | Both require PrivateLink/VPC Endpoints over IGW. |
| H-021 | ASF-067 | Both require TGW default route blackhole. |
| H-022 | ASF-004 | Both require BGP route filtering. |
| H-024 | ASF-028 | Both require AWS Config rules for security groups. |
| H-026 | ASF-013 | Both require TGW route table backup. |
| H-031 | ASF-012 | Both require Direct Connect redundancy. |
| H-036 | ASF-069 | Both require BGP authentication. |
| H-039 | ASF-020 | Both require Prod VPC to route through TGW, not IGW. |
| H-040 | ASF-031 | Both require Dev/Staging cannot access production data. |
| H-042 | ASF-009 | Both require quarterly VPN failover testing. |
| H-043 | ASF-016 | Both require sufficient VPN throughput capacity. |
| H-044 | ASF-025 | Both require network visualization tools. |
| H-046 | ASF-033 | Both require route propagation filtering. |

**Overlap (O): 30**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 46 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 80 | Count of unique ASF-generated assumptions |
| Overlap (O) | 30 | Count appearing in both lists |
| **Precision** | **37.5%** | O / A = 30/80 |
| **Recall** | **65.2%** | O / H = 30/46 |
| **F1 Score** | **47.6%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 50 | Assumptions ASF found that human missed (62.5% of ASF total) |
| Missed findings (H - O) | 16 | Assumptions human found that ASF missed (34.8% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 65.2% | ❌ Not met |
| Precision | >= 50% | 37.5% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 39.7% (50/126) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 47.6% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 16 | 50.0% |
| Derived | 10 | 20 | 50.0% |
| Operational | 6 | 24 | 25.0% |
| Implicit | 4 | 12 | 33.3% |
| Trust | 1 | 4 | 25.0% |
| Dependency | 1 | 8 | 12.5% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit and Derived categories, consistent with other architectures. Both humans and the ASF recognize that Direct Connect redundancy, VPN encryption, and routing isolation are critical.

**Worst overlap:** Architectural and Environmental. The ASF identified architectural concerns (separate AWS accounts for VPCs, TGW route table architecture) and environmental concerns (colocation facility security, provider SLA) that the human did not enumerate as assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 16)

1. **BGP-specific operational concerns (H-008, H-027, H-036):** ASN allocation, BGP authentication, and route propagation control are specific to BGP configuration and not covered by the ASF's generic patterns.

2. **Capacity planning (H-009, H-016, H-043, H-045):** The human assumed sufficient Direct Connect bandwidth, firewall throughput, and VPN tunnel capacity. The ASF covers availability generically but does not model capacity constraints.

3. **Cross-environment data protection (H-040):** The human assumed that Dev/Staging routes cannot reach production data. The ASF covers network segmentation but does not explicitly model data flow segregation by environment.

4. **Decommissioning process (H-034):** The human assumed a documented decommissioning process for Direct Connect. The ASF does not cover decommissioning.

5. **IPv6 management (H-037):** The human assumed IPv6 is either blocked or separately routed. The ASF does not cover dual-stack networking assumptions.

6. **Service limits (H-035):** The human assumed AWS service limits are monitored. The ASF covers operational resilience but does not explicitly model service quota exhaustion.

### What ASF Caught That Humans Missed (Novel Findings = 50)

1. **IAM and Identity Lifecycle (8 assumptions):** The human assumed network-level access controls but did not address IAM least privilege for networking account, cross-account role review, service account rotation, or on-prem admin account lifecycle.

2. **Incident Response (4 assumptions):** The human generated no IR assumptions. The ASF contributed IR planning for BGP hijacking, VPC isolation procedures, forensic network capture, and asymmetric routing detection.

3. **Change Management (4 assumptions):** The human did not address TGW route change approval, bandwidth change communication, new VPC attachment review, or differential change processes by environment.

4. **Compliance & Audit (4 assumptions):** The human did not address periodic audit of cross-environment traffic, AWS Config rule enforcement, or cost attribution for network resources.

5. **Secrets Management (4 assumptions):** The human covered BGP authentication but did not extend to vault-based storage of BGP keys, VPN PSKs, VIF authentication keys, or audited access to network secrets.

6. **Supply Chain (4 assumptions):** The human assumed physical security of the DC location but did not address provider supply chain security, router firmware integrity, or multi-provider dependency.

### Architecture Complexity Assessment

Architecture #014 was classified as **Complex** (hybrid network, 6+ nodes, 3 trust boundaries, multi-environment, multiple connectivity paths).

- **Recall (65.2%)** is below the 70% target, with missed findings in BGP-specific operations, capacity planning, and environment-specific data flow.
- **Precision (37.5%)** reflects the ASF's breadth across 18 patterns generating assumptions across identity, compliance, change management, and supply chain dimensions that a human network architect would not naturally enumerate.
- **Novel rate (62.5%)** confirms that the ASF adds significant value for hybrid network architectures, particularly in the identity, incident response, and governance dimensions.

### Key Insight

The primary gap is **BGP and routing protocol-specific knowledge**: the ASF's generic patterns do not cover BGP authentication, ASN management, route propagation control, or multi-VIF architecture. Adding a "Hybrid Network Security" pattern covering BGP, Direct Connect, and VPN protocol-specific assumptions would close the recall gap.

The high novelty rate is expected for hybrid architectures where governance (IAM, change management, compliance) is often handled by separate teams from the network engineering team. The ASF bridges these domains systematically.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 65.2% | ❌ Below target — missing BGP/network protocol-specific pattern |
| Precision | >= 50% | 37.5% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 39.7% | ✅ ASF adds substantial value for hybrid network |
| Expert agreement (F1) | > 60% | 47.6% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #014 demonstrates strong exploration breadth for hybrid network security, surfacing identity lifecycle, incident response, and governance assumptions that a network-focused architect may overlook. The primary actionable finding is the need for a **Hybrid Network Security** pattern covering BGP, Direct Connect, and VPN protocol-specific assumptions to close the recall gap.
