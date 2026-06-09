# ASF v1 Diagnostic Test Results

**Date:** 2026-06-09 10:33:06 UTC
**ASF Version:** 0.1.0

---

## Test 1: Happy Path

**Purpose:** See `run_asf_tests.py` — test 1 docstring.

### Input

**Policy:**
```
Only Finance employees may access payroll processing system. All payroll access requires MFA. Production databases are not internet accessible. Backups are encrypted. Quarterly access reviews are performed.
```

### ASF Output

**Claims:**
```json
[
  {
    "id": "clm_bc549726f1e6",
    "text": "Only Finance employees may access payroll processing system.",
    "tags": [
      "access",
      "process"
    ]
  },
  {
    "id": "clm_e8ccdb7e1049",
    "text": "All payroll access requires MFA.",
    "tags": [
      "access",
      "identity"
    ]
  },
  {
    "id": "clm_54592e1caa22",
    "text": "Production databases are not internet accessible.",
    "tags": [
      "access",
      "network"
    ]
  },
  {
    "id": "clm_aa5720031543",
    "text": "Backups are encrypted.",
    "tags": [
      "configuration"
    ]
  },
  {
    "id": "clm_3e648f7e697d",
    "text": "Quarterly access reviews are performed.",
    "tags": [
      "access",
      "governance"
    ]
  }
]
```

**Assumptions:**
```json
[
  {
    "id": "asm_1310d9c6c148",
    "type": "ACCESS",
    "text": "System assumes access control: Only Finance employees may access payroll processing system."
  },
  {
    "id": "asm_7d654a35c6b8",
    "type": "IDENTITY",
    "text": "System assumes identity posture: All payroll access requires MFA."
  },
  {
    "id": "asm_3a44763c8518",
    "type": "NETWORK",
    "text": "System assumes network posture: Production databases are not internet accessible."
  },
  {
    "id": "asm_507d40be6b8e",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: Backups are encrypted."
  },
  {
    "id": "asm_37b47a1f05b3",
    "type": "ACCESS",
    "text": "System assumes access control: Quarterly access reviews are performed."
  }
]
```

**Verifications:**
```json
[
  {
    "id": "vrf_f0130293e1b9",
    "assumption_id": "asm_1310d9c6c148",
    "result": "VERIFIED",
    "confidence": 0.78,
    "reasoning": "Only users in 'finance employees' found with access (4 users); Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "inline-acl": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [
          "bob",
          "carol",
          "alice",
          "frank"
        ],
        "resources_found": [
          "payroll"
        ],
        "total_records": 4
      },
      "inline-mfa": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 4
      },
      "inline-network": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 3
      },
      "inline-backup": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-db",
          "backup-server",
          "finance-fs"
        ],
        "total_records": 3
      },
      "inline-governance": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "q3-2025",
          "q1-2025",
          "q2-2025",
          "q4-2025"
        ],
        "total_records": 4
      }
    }
  },
  {
    "id": "vrf_c75e4d2612fc",
    "assumption_id": "asm_7d654a35c6b8",
    "result": "CONTRADICTED",
    "confidence": 0.95,
    "reasoning": "MFA not enabled for 4 user(s); MFA enabled for all 4 user(s); No identity evidence matched assumption; No identity evidence matched assumption; No identity evidence matched assumption",
    "details": {
      "inline-acl": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": [
          "alice",
          "bob",
          "carol",
          "frank"
        ]
      },
      "inline-mfa": {
        "checks": [],
        "mfa_enabled_users": [
          "alice",
          "bob",
          "carol",
          "frank"
        ],
        "mfa_disabled_users": []
      },
      "inline-network": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": []
      },
      "inline-backup": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": []
      },
      "inline-governance": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": []
      }
    }
  },
  {
    "id": "vrf_d22b91949623",
    "assumption_id": "asm_3a44763c8518",
    "result": "VERIFIED",
    "confidence": 0.85,
    "reasoning": "Confirmed: no exposure found across 4 asset(s); Could not verify network posture; Confirmed: no exposure found across 3 asset(s); Confirmed: no exposure found across 3 asset(s); Could not verify network posture",
    "details": {
      "inline-acl": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll",
          "payroll",
          "payroll",
          "payroll"
        ]
      },
      "inline-mfa": {
        "exposed_assets": [],
        "isolated_assets": []
      },
      "inline-network": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-db",
          "finance-db",
          "analytics-db"
        ]
      },
      "inline-backup": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-db",
          "finance-fs",
          "backup-server"
        ]
      },
      "inline-governance": {
        "exposed_assets": [],
        "isolated_assets": []
      }
    }
  },
  {
    "id": "vrf_471d30af6bc3",
    "assumption_id": "asm_507d40be6b8e",
    "result": "VERIFIED",
    "confidence": 0.85,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Could not verify configuration from evidence; All 3 resource(s) compliant with configuration; Could not verify configuration from evidence",
    "details": {
      "inline-acl": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "inline-mfa": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "inline-network": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "inline-backup": {
        "compliant": 3,
        "non_compliant": 0,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "backup-server"
        ],
        "examples_non_compliant": []
      },
      "inline-governance": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_c63cc0d905ca",
    "assumption_id": "asm_37b47a1f05b3",
    "result": "UNKNOWN",
    "confidence": 0.0,
    "reasoning": "Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "inline-acl": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll"
        ],
        "total_records": 4
      },
      "inline-mfa": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 4
      },
      "inline-network": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 3
      },
      "inline-backup": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-db",
          "backup-server",
          "finance-fs"
        ],
        "total_records": 3
      },
      "inline-governance": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "q3-2025",
          "q1-2025",
          "q2-2025",
          "q4-2025"
        ],
        "total_records": 4
      }
    }
  }
]
```

**Gaps:**
```json
[
  {
    "id": "gap_988d0a6d9534",
    "severity": "CRITICAL",
    "type": "IDENTITY_GAP",
    "description": "Assumption contradicted: System assumes identity posture: All payroll access requires MFA."
  },
  {
    "id": "gap_029574ac2b03",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: Quarterly access reviews are performed."
  }
]
```

### Analysis

_Extracted 5 claims, 5 assumptions, 5 verifications, 2 gaps (1 critical). All verified: False. Verdict: FAIL_

**Verdict:** `FAIL`

---

## Test 2: Direct Contradiction

**Purpose:** See `run_asf_tests.py` — test 2 docstring.

### Input

**Policy:**
```
Only Finance employees may access payroll processing system.
```

### ASF Output

**Claims:**
```json
[
  {
    "id": "clm_ccfb463ed2e6",
    "text": "Only Finance employees may access payroll processing system."
  }
]
```

**Assumptions:**
```json
[
  {
    "id": "asm_ee7956d20c47",
    "type": "ACCESS",
    "text": "System assumes access control: Only Finance employees may access payroll processing system."
  }
]
```

**Verifications:**
```json
[
  {
    "id": "vrf_6f9b672f4a3c",
    "assumption_id": "asm_ee7956d20c47",
    "result": "CONTRADICTED",
    "confidence": 0.92,
    "reasoning": "Found 1 user(s) outside 'finance employees' with access: sarah",
    "details": {
      "inline-acl": {
        "expected_group": "finance employees",
        "users_outside_group": [
          "sarah"
        ],
        "users_inside_group": [
          "john"
        ],
        "resources_found": [
          "payroll"
        ],
        "total_records": 2
      }
    }
  }
]
```

**Gaps:**
```json
[
  {
    "id": "gap_cae6f40e8c49",
    "severity": "CRITICAL",
    "type": "ACCESS_GAP",
    "description": "Assumption contradicted: System assumes access control: Only Finance employees may access payroll processing system."
  }
]
```

### Analysis

_Contradicted: True. Reasoning explains issue: True (found 'sarah' or 'outside' or 'Engineering' in reasoning). Verdict: PASS_

**Verdict:** `PASS`

---

## Test 3: Missing Evidence

**Purpose:** See `run_asf_tests.py` — test 3 docstring.

### Input

**Policy:**
```
All payroll access requires MFA.
```

### ASF Output

**Claims:**
```json
[
  {
    "id": "clm_d59f42761956",
    "text": "All payroll access requires MFA."
  }
]
```

**Assumptions:**
```json
[
  {
    "id": "asm_757288eb7b92",
    "type": "IDENTITY",
    "text": "System assumes identity posture: All payroll access requires MFA."
  }
]
```

**Verifications:**
```json
[
  {
    "id": "vrf_14719812e899",
    "assumption_id": "asm_757288eb7b92",
    "result": "UNKNOWN",
    "confidence": 0.0,
    "reasoning": "No matching evidence available for verification",
    "details": {}
  }
]
```

**Gaps:**
```json
[
  {
    "id": "gap_eda5b076bab7",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes identity posture: All payroll access requires MFA."
  }
]
```

### Analysis

_All verifications UNKNOWN: True. ASF correctly returned UNKNOWN instead of assuming false. Gap type: EVIDENCE_GAP (LOW). Verdict: PASS_

**Verdict:** `PASS`

---

## Test 4: Garbage Evidence

**Purpose:** See `run_asf_tests.py` — test 4 docstring.

### Input

**Policy:**
```
Only Finance employees may access payroll.
```

### ASF Output

**Claims:**
```json
[
  {
    "id": "clm_0f89659d6ba5",
    "text": "Only Finance employees may access payroll."
  }
]
```

**Assumptions:**
```json
[
  {
    "id": "asm_df4216e2a3d3",
    "type": "ACCESS",
    "text": "System assumes access control: Only Finance employees may access payroll."
  }
]
```

**Verifications:**
```json
[
  {
    "id": "vrf_07e7dbbaaf11",
    "assumption_id": "asm_df4216e2a3d3",
    "result": "UNKNOWN",
    "confidence": 0.0,
    "reasoning": "Could not determine access patterns from evidence",
    "details": {
      "garbage.csv": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 2
      }
    }
  }
]
```

**Gaps:**
```json
[
  {
    "id": "gap_9065012c77eb",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: Only Finance employees may access payroll."
  }
]
```

### Analysis

_Crashed: False. Output UNKNOWN: True. ASF remained stable on garbage input. Verdict: PASS_

**Verdict:** `PASS`

---

## Test 5: Real Policy Analysis

**Purpose:** See `run_asf_tests.py` — test 5 docstring.

### Input

**Policy File:** `/Users/moksh/Project/cybersec/sample_data/finance_policy.txt`
**Evidence Files:**
- `/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv`
- `/Users/moksh/Project/cybersec/sample_data/mfa_status.csv`
- `/Users/moksh/Project/cybersec/sample_data/backup_config.csv`
- `/Users/moksh/Project/cybersec/sample_data/network_exposure.csv`

### ASF Output

**Claims:**
```json
[
  {
    "id": "clm_d8704c572fb7",
    "text": "ACCESS CONTROL\n\nOnly Finance employees may access the payroll processing system.",
    "tags": [
      "access",
      "process"
    ]
  },
  {
    "id": "clm_35d048ac366d",
    "text": "All payroll data is encrypted at rest and in transit.",
    "tags": [
      "configuration"
    ]
  },
  {
    "id": "clm_e69c0c20fe08",
    "text": "Only the VP of Finance can approve payroll runs.",
    "tags": [
      "governance"
    ]
  },
  {
    "id": "clm_2a7f2d1b87c2",
    "text": "SYSTEM ACCESS\n\nProduction databases are not internet accessible.",
    "tags": [
      "access",
      "network"
    ]
  },
  {
    "id": "clm_22f812083f65",
    "text": "All administrative access requires multi-factor authentication.",
    "tags": [
      "access",
      "identity"
    ]
  },
  {
    "id": "clm_446a4e2d8761",
    "text": "Database access is restricted to database administrators only.",
    "tags": [
      "access"
    ]
  },
  {
    "id": "clm_e72b44b7c60a",
    "text": "SSH access to production servers is restricted to the infrastructure team.",
    "tags": [
      "access"
    ]
  },
  {
    "id": "clm_4238857d5dae",
    "text": "AUDIT AND COMPLIANCE\n\nAll access to financial systems is logged and monitored.",
    "tags": [
      "access",
      "configuration"
    ]
  },
  {
    "id": "clm_9ff853d4d651",
    "text": "Security reviews are conducted quarterly.",
    "tags": [
      "governance"
    ]
  },
  {
    "id": "clm_515eeefcd3d5",
    "text": "All configuration changes must be approved by the security team.",
    "tags": [
      "governance"
    ]
  },
  {
    "id": "clm_cc421d96b37b",
    "text": "BACKUP AND RECOVERY\n\nAll financial data is backed up daily.",
    "tags": [
      "configuration"
    ]
  },
  {
    "id": "clm_176f08d969e2",
    "text": "Backup data is encrypted using AES-256.",
    "tags": [
      "configuration"
    ]
  },
  {
    "id": "clm_d4c193d23492",
    "text": "Backups are tested monthly to ensure recoverability.",
    "tags": [
      "configuration",
      "process"
    ]
  },
  {
    "id": "clm_2ee1f53ce81d",
    "text": "Offsite backup storage is maintained at a separate geographic location.",
    "tags": [
      "configuration"
    ]
  },
  {
    "id": "clm_25c3a4cc9fab",
    "text": "NETWORK SECURITY\n\nThe finance network segment is isolated from other network segments.",
    "tags": [
      "network"
    ]
  },
  {
    "id": "clm_a669603f3776",
    "text": "All network traffic is inspected by the intrusion detection system.",
    "tags": [
      "network"
    ]
  },
  {
    "id": "clm_854877590a8b",
    "text": "Security groups restrict access to only required ports and protocols.",
    "tags": [
      "access"
    ]
  }
]
```

**Assumptions:**
```json
[
  {
    "id": "asm_e0b3b469bd25",
    "type": "ACCESS",
    "text": "System assumes access control: ACCESS CONTROL\n\nOnly Finance employees may access the payroll processing system.",
    "confidence": 0.0968
  },
  {
    "id": "asm_261ec92ea111",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: All payroll data is encrypted at rest and in transit.",
    "confidence": 0.35
  },
  {
    "id": "asm_2ec9fbff1177",
    "type": "PROCESS",
    "text": "System assumes process compliance: Only the VP of Finance can approve payroll runs.",
    "confidence": 0.0952
  },
  {
    "id": "asm_acf471c6dfb0",
    "type": "ACCESS",
    "text": "System assumes access control: SYSTEM ACCESS\n\nProduction databases are not internet accessible.",
    "confidence": 0.0
  },
  {
    "id": "asm_a76b59349d74",
    "type": "IDENTITY",
    "text": "System assumes identity posture: All administrative access requires multi-factor authentication.",
    "confidence": 0.098
  },
  {
    "id": "asm_5239942f3c34",
    "type": "ACCESS",
    "text": "System assumes access control: Database access is restricted to database administrators only.",
    "confidence": 0.0
  },
  {
    "id": "asm_aa6f1b93c5fb",
    "type": "ACCESS",
    "text": "System assumes access control: SSH access to production servers is restricted to the infrastructure team.",
    "confidence": 0.0
  },
  {
    "id": "asm_188afee9fc91",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: AUDIT AND COMPLIANCE\n\nAll access to financial systems is logged and monitored.",
    "confidence": 0.35
  },
  {
    "id": "asm_fa16d63ee584",
    "type": "GOVERNANCE",
    "text": "System assumes governance compliance: Security reviews are conducted quarterly.",
    "confidence": 0.0952
  },
  {
    "id": "asm_9ba11e1aea41",
    "type": "PROCESS",
    "text": "System assumes process compliance: All configuration changes must be approved by the security team.",
    "confidence": 0.0952
  },
  {
    "id": "asm_a0f821f66b8a",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: BACKUP AND RECOVERY\n\nAll financial data is backed up daily.",
    "confidence": 0.35
  },
  {
    "id": "asm_6cda6401d6d9",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: Backup data is encrypted using AES-256.",
    "confidence": 0.35
  },
  {
    "id": "asm_3a7dcba090ed",
    "type": "GOVERNANCE",
    "text": "System assumes governance compliance: Backups are tested monthly to ensure recoverability.",
    "confidence": 0.0952
  },
  {
    "id": "asm_52cffdb5e0a6",
    "type": "CONFIGURATION",
    "text": "System assumes configuration state: Offsite backup storage is maintained at a separate geographic location.",
    "confidence": 0.35
  },
  {
    "id": "asm_186ad9c69708",
    "type": "NETWORK",
    "text": "System assumes network posture: NETWORK SECURITY\n\nThe finance network segment is isolated from other network segments.",
    "confidence": 0.096
  },
  {
    "id": "asm_839472ccc099",
    "type": "NETWORK",
    "text": "System assumes network posture: All network traffic is inspected by the intrusion detection system.",
    "confidence": 0.0
  },
  {
    "id": "asm_a5672ca9c37f",
    "type": "ACCESS",
    "text": "System assumes access control: Security groups restrict access to only required ports and protocols.",
    "confidence": 0.0
  }
]
```

**Verifications:**
```json
[
  {
    "id": "vrf_64e6576079a2",
    "assumption_id": "asm_e0b3b469bd25",
    "result": "CONTRADICTED",
    "confidence": 0.968,
    "reasoning": "Found 4 user(s) outside 'finance employees' with access: henry.taylor, dave.wilson, eve.brown, jack.thompson; Found 4 user(s) outside 'finance employees' with access: henry.taylor, dave.wilson, eve.brown, jack.thompson; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "expected_group": "finance employees",
        "users_outside_group": [
          "henry.taylor",
          "dave.wilson",
          "eve.brown",
          "jack.thompson"
        ],
        "users_inside_group": [
          "alice.jones",
          "iris.anderson",
          "frank.miller",
          "carol.davis",
          "grace.lee",
          "bob.smith"
        ],
        "resources_found": [
          "payroll-system"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "expected_group": "finance employees",
        "users_outside_group": [
          "henry.taylor",
          "dave.wilson",
          "eve.brown",
          "jack.thompson"
        ],
        "users_inside_group": [
          "alice.jones",
          "iris.anderson",
          "frank.miller",
          "carol.davis",
          "grace.lee",
          "bob.smith"
        ],
        "resources_found": [],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "backup-server",
          "dev-db",
          "analytics-db",
          "monitoring-data",
          "hr-system",
          "payroll-db",
          "staging-db",
          "logs-archive",
          "finance-fs",
          "customer-db"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "expected_group": "finance employees",
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      }
    }
  },
  {
    "id": "vrf_8cee73436b07",
    "assumption_id": "asm_261ec92ea111",
    "result": "PARTIALLY_VERIFIED",
    "confidence": 0.7,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Partially compliant: 7 OK, 3 non-compliant; Could not verify configuration from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "compliant": 7,
        "non_compliant": 3,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "customer-db"
        ],
        "examples_non_compliant": [
          "analytics-db",
          "dev-db",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_d5a85eb78cb8",
    "assumption_id": "asm_2ec9fbff1177",
    "result": "CONTRADICTED",
    "confidence": 0.952,
    "reasoning": "No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending)",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      }
    }
  },
  {
    "id": "vrf_82d43d57e666",
    "assumption_id": "asm_acf471c6dfb0",
    "result": "UNKNOWN",
    "confidence": 0.4,
    "reasoning": "Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-system"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "backup-server",
          "dev-db",
          "analytics-db",
          "monitoring-data",
          "hr-system",
          "payroll-db",
          "staging-db",
          "logs-archive",
          "finance-fs",
          "customer-db"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      }
    }
  },
  {
    "id": "vrf_cccf1a428e61",
    "assumption_id": "asm_a76b59349d74",
    "result": "CONTRADICTED",
    "confidence": 0.98,
    "reasoning": "MFA not enabled for 10 user(s); MFA enabled for 7 user(s) but missing for 3; No identity evidence matched assumption; No identity evidence matched assumption",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": [
          "alice.jones",
          "bob.smith",
          "carol.davis",
          "dave.wilson",
          "eve.brown",
          "frank.miller",
          "grace.lee",
          "henry.taylor",
          "iris.anderson",
          "jack.thompson"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "checks": [],
        "mfa_enabled_users": [
          "alice.jones",
          "bob.smith",
          "carol.davis",
          "dave.wilson",
          "eve.brown",
          "grace.lee",
          "henry.taylor"
        ],
        "mfa_disabled_users": [
          "frank.miller",
          "iris.anderson",
          "jack.thompson"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": []
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "checks": [],
        "mfa_enabled_users": [],
        "mfa_disabled_users": []
      }
    }
  },
  {
    "id": "vrf_431a52f9156d",
    "assumption_id": "asm_5239942f3c34",
    "result": "UNKNOWN",
    "confidence": 0.4,
    "reasoning": "Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-system"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "backup-server",
          "dev-db",
          "analytics-db",
          "monitoring-data",
          "hr-system",
          "payroll-db",
          "staging-db",
          "logs-archive",
          "finance-fs",
          "customer-db"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      }
    }
  },
  {
    "id": "vrf_7ad8e7f2e609",
    "assumption_id": "asm_aa6f1b93c5fb",
    "result": "UNKNOWN",
    "confidence": 0.4,
    "reasoning": "Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-system"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "backup-server",
          "dev-db",
          "analytics-db",
          "monitoring-data",
          "hr-system",
          "payroll-db",
          "staging-db",
          "logs-archive",
          "finance-fs",
          "customer-db"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      }
    }
  },
  {
    "id": "vrf_da82a02b6fc5",
    "assumption_id": "asm_188afee9fc91",
    "result": "PARTIALLY_VERIFIED",
    "confidence": 0.7,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Partially compliant: 7 OK, 3 non-compliant; Could not verify configuration from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "compliant": 7,
        "non_compliant": 3,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "customer-db"
        ],
        "examples_non_compliant": [
          "analytics-db",
          "dev-db",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_274c9681efdb",
    "assumption_id": "asm_fa16d63ee584",
    "result": "CONTRADICTED",
    "confidence": 0.952,
    "reasoning": "No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending)",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      }
    }
  },
  {
    "id": "vrf_15ae48d133a7",
    "assumption_id": "asm_9ba11e1aea41",
    "result": "CONTRADICTED",
    "confidence": 0.952,
    "reasoning": "No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending)",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      }
    }
  },
  {
    "id": "vrf_81f8d0424b58",
    "assumption_id": "asm_a0f821f66b8a",
    "result": "PARTIALLY_VERIFIED",
    "confidence": 0.7,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Partially compliant: 7 OK, 3 non-compliant; Could not verify configuration from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "compliant": 7,
        "non_compliant": 3,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "customer-db"
        ],
        "examples_non_compliant": [
          "analytics-db",
          "dev-db",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_648d0cf8373b",
    "assumption_id": "asm_6cda6401d6d9",
    "result": "PARTIALLY_VERIFIED",
    "confidence": 0.7,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Partially compliant: 7 OK, 3 non-compliant; Could not verify configuration from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "compliant": 7,
        "non_compliant": 3,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "customer-db"
        ],
        "examples_non_compliant": [
          "analytics-db",
          "dev-db",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_5ef7b5fdf4d3",
    "assumption_id": "asm_3a7dcba090ed",
    "result": "CONTRADICTED",
    "confidence": 0.952,
    "reasoning": "No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending); No governance reviews completed (10 pending)",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "reviews_completed": 0,
        "reviews_pending": 10
      }
    }
  },
  {
    "id": "vrf_68b2da551769",
    "assumption_id": "asm_52cffdb5e0a6",
    "result": "PARTIALLY_VERIFIED",
    "confidence": 0.7,
    "reasoning": "Could not verify configuration from evidence; Could not verify configuration from evidence; Partially compliant: 7 OK, 3 non-compliant; Could not verify configuration from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "compliant": 7,
        "non_compliant": 3,
        "examples_compliant": [
          "payroll-db",
          "finance-fs",
          "customer-db"
        ],
        "examples_non_compliant": [
          "analytics-db",
          "dev-db",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "compliant": 0,
        "non_compliant": 0,
        "examples_compliant": [],
        "examples_non_compliant": []
      }
    }
  },
  {
    "id": "vrf_65955e0ebeb6",
    "assumption_id": "asm_186ad9c69708",
    "result": "CONTRADICTED",
    "confidence": 0.96,
    "reasoning": "All 10 asset(s) appear isolated; All 0 asset(s) appear isolated; All 10 asset(s) appear isolated; Claimed isolated but found 4 exposed asset(s): customer-portal, api-gateway, dev-server, admin-panel",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "exposed_assets": [],
        "isolated_assets": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-db",
          "finance-fs",
          "customer-db",
          "analytics-db",
          "backup-server",
          "logs-archive",
          "dev-db",
          "staging-db",
          "monitoring-data",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "exposed_assets": [
          "customer-portal",
          "api-gateway",
          "dev-server",
          "admin-panel"
        ],
        "isolated_assets": [
          "payroll-app",
          "finance-db",
          "staging-app",
          "analytics-db",
          "backup-server",
          "monitoring"
        ]
      }
    }
  },
  {
    "id": "vrf_6275fbf1e3a7",
    "assumption_id": "asm_839472ccc099",
    "result": "UNKNOWN",
    "confidence": 0.4,
    "reasoning": "Could not verify network posture; Could not verify network posture; Could not verify network posture; Could not verify network posture",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system",
          "payroll-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "exposed_assets": [],
        "isolated_assets": []
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "exposed_assets": [],
        "isolated_assets": [
          "payroll-db",
          "finance-fs",
          "customer-db",
          "analytics-db",
          "backup-server",
          "logs-archive",
          "dev-db",
          "staging-db",
          "monitoring-data",
          "hr-system"
        ]
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "exposed_assets": [
          "customer-portal",
          "api-gateway",
          "dev-server",
          "admin-panel"
        ],
        "isolated_assets": [
          "payroll-app",
          "finance-db",
          "staging-app",
          "analytics-db",
          "backup-server",
          "monitoring"
        ]
      }
    }
  },
  {
    "id": "vrf_4565528c4161",
    "assumption_id": "asm_a5672ca9c37f",
    "result": "UNKNOWN",
    "confidence": 0.4,
    "reasoning": "Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence; Could not determine access patterns from evidence",
    "details": {
      "/Users/moksh/Project/cybersec/sample_data/payroll_acl.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "payroll-system"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/mfa_status.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/backup_config.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [
          "backup-server",
          "dev-db",
          "analytics-db",
          "monitoring-data",
          "hr-system",
          "payroll-db",
          "staging-db",
          "logs-archive",
          "finance-fs",
          "customer-db"
        ],
        "total_records": 10
      },
      "/Users/moksh/Project/cybersec/sample_data/network_exposure.csv": {
        "expected_group": null,
        "users_outside_group": [],
        "users_inside_group": [],
        "resources_found": [],
        "total_records": 10
      }
    }
  }
]
```

**Gaps:**
```json
[
  {
    "id": "gap_23a45559e167",
    "severity": "CRITICAL",
    "type": "ACCESS_GAP",
    "description": "Assumption contradicted: System assumes access control: ACCESS CONTROL\n\nOnly Finance employees may access the payroll processing system."
  },
  {
    "id": "gap_823d2376d4d3",
    "severity": "MEDIUM",
    "type": "CONFIGURATION_GAP",
    "description": "Assumption only partially verified: System assumes configuration state: All payroll data is encrypted at rest and in transit."
  },
  {
    "id": "gap_bc9d405e59b5",
    "severity": "HIGH",
    "type": "PROCESS_GAP",
    "description": "Assumption contradicted: System assumes process compliance: Only the VP of Finance can approve payroll runs."
  },
  {
    "id": "gap_22f5902e6eb2",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: SYSTEM ACCESS\n\nProduction databases are not internet accessible."
  },
  {
    "id": "gap_2f7b7abce567",
    "severity": "CRITICAL",
    "type": "IDENTITY_GAP",
    "description": "Assumption contradicted: System assumes identity posture: All administrative access requires multi-factor authentication."
  },
  {
    "id": "gap_542a8f59999f",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: Database access is restricted to database administrators only."
  },
  {
    "id": "gap_3a868dcc7e3a",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: SSH access to production servers is restricted to the infrastructure team."
  },
  {
    "id": "gap_4e0bafe024cf",
    "severity": "MEDIUM",
    "type": "CONFIGURATION_GAP",
    "description": "Assumption only partially verified: System assumes configuration state: AUDIT AND COMPLIANCE\n\nAll access to financial systems is logged and monitored."
  },
  {
    "id": "gap_8c41d947bb94",
    "severity": "HIGH",
    "type": "GOVERNANCE_GAP",
    "description": "Assumption contradicted: System assumes governance compliance: Security reviews are conducted quarterly."
  },
  {
    "id": "gap_292648bfb69a",
    "severity": "HIGH",
    "type": "PROCESS_GAP",
    "description": "Assumption contradicted: System assumes process compliance: All configuration changes must be approved by the security team."
  },
  {
    "id": "gap_a2acfdb10b02",
    "severity": "MEDIUM",
    "type": "CONFIGURATION_GAP",
    "description": "Assumption only partially verified: System assumes configuration state: BACKUP AND RECOVERY\n\nAll financial data is backed up daily."
  },
  {
    "id": "gap_9bac80c57ae1",
    "severity": "MEDIUM",
    "type": "CONFIGURATION_GAP",
    "description": "Assumption only partially verified: System assumes configuration state: Backup data is encrypted using AES-256."
  },
  {
    "id": "gap_193c3afd6bad",
    "severity": "HIGH",
    "type": "GOVERNANCE_GAP",
    "description": "Assumption contradicted: System assumes governance compliance: Backups are tested monthly to ensure recoverability."
  },
  {
    "id": "gap_0dcd0e17e1b2",
    "severity": "MEDIUM",
    "type": "CONFIGURATION_GAP",
    "description": "Assumption only partially verified: System assumes configuration state: Offsite backup storage is maintained at a separate geographic location."
  },
  {
    "id": "gap_e394f3a6997a",
    "severity": "CRITICAL",
    "type": "NETWORK_GAP",
    "description": "Assumption contradicted: System assumes network posture: NETWORK SECURITY\n\nThe finance network segment is isolated from other network segments."
  },
  {
    "id": "gap_376bbb583934",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes network posture: All network traffic is inspected by the intrusion detection system."
  },
  {
    "id": "gap_477c54dc30f6",
    "severity": "LOW",
    "type": "EVIDENCE_GAP",
    "description": "Insufficient evidence to verify: System assumes access control: Security groups restrict access to only required ports and protocols."
  }
]
```

**Summary:**
```json
{
  "claims_found": 17,
  "assumptions_found": 17,
  "verified": 0,
  "contradicted": 7,
  "unknown": 5,
  "critical_gaps": 3
}
```

### Analysis

_Real policy: 17 claims, 0 verified, 7 contradicted, 5 unknown, 3 critical gaps. Report is actionable: YES. Verdict: PASS_

**Verdict:** `PASS`

---

## Final Summary

| Test | Name | Verdict | Key Finding |
|------|------|---------|-------------|
| 1 | Happy Path | FAIL | Extracted 5 claims, 5 assumptions, 5 verifications, 2 gaps ( |
| 2 | Direct Contradiction | PASS | Contradicted: True. Reasoning explains issue: True (found 's |
| 3 | Missing Evidence | PASS | All verifications UNKNOWN: True. ASF correctly returned UNKN |
| 4 | Garbage Evidence | PASS | Crashed: False. Output UNKNOWN: True. ASF remained stable on |
| 5 | Real Policy Analysis | PASS | Real policy: 17 claims, 0 verified, 7 contradicted, 5 unknow |
