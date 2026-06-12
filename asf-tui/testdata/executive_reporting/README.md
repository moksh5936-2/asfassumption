# Executive Risk Narratives — Benchmark Testdata

This directory contains benchmark test data for validating the ERN engine
across different domain knowledge packs.

## Structure

Each subdirectory represents a domain with expected ERN outputs:

- `healthcare/` — Healthcare domain (HIPAA, PHI, EHR)
- `fintech/` — Fintech domain (PCI DSS, payments, cardholder data)
- `saas/` — SaaS domain (multi-tenant, SOC 2)
- `government/` — Government domain (FedRAMP, FISMA)
- `kubernetes/` — Kubernetes domain (container security, CIS)

## Expected Outputs

Each domain subdirectory contains:

- `input.json` — ERNInput fixture
- `expected_board_summary.txt` — Expected board-level summary
- `expected_ciso_briefing.json` — Expected CISO briefing structure
- `expected_roadmap.json` — Expected remediation roadmap

## Usage

Testdata is consumed by `TestERNBenchmark*` functions in `intelligence/ern_test.go`.
