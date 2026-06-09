"""
Forensic Analysis of Missed Assumptions.
Classifies all 1,697 ground truth assumptions into 8-type ontology.
This is analysis, not code. Output is knowledge.
"""
from __future__ import annotations
import json
import re
import random
from pathlib import Path
from collections import Counter


# ── Ontology classification rules ──────────────────────────────

def classify_ontology(text: str, current_category: str, asf_type: str, commentary: str) -> str:
    tl = text.lower()
    cl = commentary.lower()

    # Explicit — directly stated
    if current_category == "explicit":
        return "Explicit"

    # Dependency — external systems, vendors, third parties
    dep_keywords = [
        "vendor", "third.party", "supplier", "external", "sla",
        "depends? on", "reli(es|ed|es)\s+on", "integration",
        "connected", "communicate", "upstream", "downstream",
        "partner", "federation", "federated",
        "certificate authority", "identity provider", "idp",
        "open source", "oss", "library",
        "supply chain", "container base image", "image registry",
    ]
    if any(re.search(k, tl) for k in dep_keywords):
        return "Dependency"

    # Trust — identity, authorization, relationship
    trust_keywords = [
        "trust", "trusted", "federat",
        "identity", "authentication", "authorization", "auth",
        "group memberships? are accurate",
        "permissions? are correct",
        "role assignment", "rbac",
        "access decisions are consistent",
        "enrollment is complete",
        "hr data is accurate", "hr system",
        "employee status", "termination notification",
        "consent ", "consent record",
        "vendor access is promptly revoked",
        "vendor security posture",
        "vendor breach notification",
    ]
    if any(re.search(k, tl) for k in trust_keywords):
        return "Trust"

    # Operational — process, procedure, review, testing
    operational_keywords = [
        "process", "procedure",
        "review", "audit", "inspect",
        "train", "test", "tested", "testing",
        "performed", "conduct", "schedule",
        "compliance is validated",
        "change management", "change approval",
        "incident response", "ir plan",
        "playbook", "runbook",
        "policy review", "policy updat",
        "access review", "recertification",
        "remediation", "remediate",
        "backup test", "restore test",
        "disaster recovery test",
        "security training", "awareness training",
        "forensic", "evidence chain",
        "lessons learned",
    ]
    if any(re.search(k, tl) for k in operational_keywords):
        return "Operational"

    # Architectural — design, structure, boundaries
    architectural_keywords = [
        "architect", "design",
        "boundar", "trust boundary",
        "segment", "isolat", "network boundary",
        "tier", "layer",
        "single point of failure",
        "failover mechanism",
        "circuit breaker", "backpressure",
        "graceful degradation",
        "no single point",
        "route table", "security group",
        "vpc peering", "transit gateway",
        "service mesh", "sidecar",
        "tenant isolation", "multi.tenant",
        "dmz", "bastion",
    ]
    if any(re.search(k, tl) for k in architectural_keywords):
        return "Architectural"

    # Environmental — infrastructure, environment, capacity
    environmental_keywords = [
        "environment", "production", "staging", "development",
        "region", "availability zone", "multi.region",
        "capacity", "scaling", "auto.scal", "load bal",
        "resource", "compute", "storage capacity",
        "network flow log", "vpc flow log",
        "monitoring cover", "all critical",
        "log storage capacity", "retention period",
        "backup job complet", "backup window",
        "rpo", "rto", "recovery point", "recovery time",
        "ups", "generator", "power backup",
        "cooling", "hvac", "temperature",
        "disk encryption",
        "geo.redundan", "cross.region",
    ]
    if any(re.search(k, tl) for k in environmental_keywords):
        return "Environmental"

    # Derived — logical consequence, risk statement
    derived_keywords = [
        "violation of this policy leads to",
        "if this policy fails",
        "consequence", "therefore",
        "leads to", "results in",
        "risk", "compromise", "exposure",
        "would", "could", "may",
        "unverifiable", "cannot be verified",
    ]
    if any(re.search(k, tl) for k in derived_keywords):
        return "Derived"

    # Implicit — implied but not stated, default fallback for "implicit" category
    if current_category == "implicit":
        return "Implicit"

    return "Derived"


def run_analysis():
    base = Path(__file__).parent.parent / "report"
    gt_path = base / "ground_truth.json"

    with open(gt_path) as f:
        gt = json.load(f)

    assumptions = gt["assumptions"]
    print(f"Total assumptions loaded: {len(assumptions)}")

    # Classify every assumption
    classifications = []
    for a in assumptions:
        onto = classify_ontology(a["text"], a["category"], a["type"], a.get("commentary", ""))
        classifications.append({
            "id": a["id"],
            "text": a["text"],
            "current_category": a["category"],
            "asf_type": a["type"],
            "ontology": onto,
            "keywords": a.get("keywords", []),
            "commentary": a.get("commentary", ""),
        })

    # Count frequencies
    onto_counts = Counter(c["ontology"] for c in classifications)
    type_x_onto: dict[str, Counter] = {}
    for c in classifications:
        t = c["asf_type"]
        if t not in type_x_onto:
            type_x_onto[t] = Counter()
        type_x_onto[t][c["ontology"]] += 1

    cat_x_onto: dict[str, Counter] = {}
    for c in classifications:
        cat = c["current_category"]
        if cat not in cat_x_onto:
            cat_x_onto[cat] = Counter()
        cat_x_onto[cat][c["ontology"]] += 1

    total = len(assumptions)

    # Print results
    print("\n" + "=" * 60)
    print("ASF ASSUMPTION ONTOLOGY — FORENSIC ANALYSIS")
    print("=" * 60)

    print(f"\nTotal assumptions: {total}")
    print(f"\n## Ontology Distribution")
    print(f"\n{'Ontology':25s} {'Count':>6s} {'Percent':>8s}")
    print("-" * 42)
    for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
        cnt = onto_counts.get(onto, 0)
        pct = cnt / total * 100
        bar = "█" * int(cnt / total * 40)
        print(f"{onto:25s} {cnt:6d} {pct:7.1f}%  {bar}")
    print("-" * 42)
    print(f"{'TOTAL':25s} {total:6d} {100.0:7.1f}%")

    print(f"\n## Cross-tabulation: ASF Type × Ontology")
    print(f"\n{'Type':20s}", end="")
    for onto in ["Exp", "Imp", "Der", "Tru", "Ops", "Dep", "Arc", "Env"]:
        print(f"{onto:>7s}", end="")
    print(f"{'Total':>7s}")
    print("-" * 62)
    for asf_type in ["ACCESS", "IDENTITY", "NETWORK", "CONFIGURATION", "PROCESS", "GOVERNANCE", "DOCUMENTATION", "DEPENDENCY"]:
        if asf_type not in type_x_onto:
            continue
        c = type_x_onto[asf_type]
        row_total = sum(c.values())
        print(f"{asf_type:20s}", end="")
        for onto_short, onto_full in [("Exp", "Explicit"), ("Imp", "Implicit"), ("Der", "Derived"), ("Tru", "Trust"), ("Ops", "Operational"), ("Dep", "Dependency"), ("Arc", "Architectural"), ("Env", "Environmental")]:
            print(f"{c.get(onto_full, 0):7d}", end="")
        print(f"{row_total:7d}")

    print(f"\n## Cross-tabulation: Current Category × Ontology")
    for cat in ["explicit", "implicit", "derived"]:
        c = cat_x_onto.get(cat, Counter())
        print(f"\n  {cat} ({sum(c.values())} total):")
        for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
            if c.get(onto, 0) > 0:
                print(f"    {onto:20s} {c[onto]}")

    # ── Sample analysis: 200 missed assumptions ──
    # Determine which assumptions were "missed" (ASF didn't match them)
    print("\n" + "=" * 60)
    print("MISSED ASSUMPTIONS ANALYSIS")
    print("=" * 60)

    asf_path = base / "asf_assumptions.json"
    with open(asf_path) as f:
        asf_data = json.load(f)

    from benchmark.runner import compute_match, GroundTruthAssumption  # noqa

    # Reconstruct ground truth objects
    gt_objects = []
    for a in assumptions:
        gt_objects.append(GroundTruthAssumption(
            id=a["id"], policy_id="", text=a["text"],
            type=a["type"], category=a["category"],
            commentary=a.get("commentary", ""),
            keywords=a.get("keywords", []),
        ))

    # For each GT assumption, check if it matched any ASF assumption
    missed_by_onto: Counter = Counter()
    missed_examples: dict[str, list] = {"Explicit": [], "Implicit": [], "Derived": [], "Trust": [],
                                         "Operational": [], "Dependency": [], "Architectural": [], "Environmental": []}

    # We need to match each GT to its policy. Since we lost the policy_id in the export,
    # let's re-load from the ground truth builder
    from benchmark.ground_truth import build_ground_truth
    full_gt = build_ground_truth()

    # Build a lookup: text -> ontology classification
    text_to_onto = {c["text"]: c["ontology"] for c in classifications}

    missed = 0
    total_checked = 0
    for policy in full_gt.policies:
        pid = policy.id
        gt_list = full_gt.get_assumptions_for(pid)
        asm_list = asf_data.get(pid, [])
        for gt_a in gt_list:
            score = compute_match(gt_a, asm_list)
            total_checked += 1
            if score < 0.35:
                missed += 1
                onto = text_to_onto.get(gt_a.text, "Derived")
                missed_by_onto[onto] += 1
                if len(missed_examples.get(onto, [])) < 25:
                    missed_examples.setdefault(onto, []).append({
                        "text": gt_a.text[:100],
                        "type": str(gt_a.type),
                        "commentary": gt_a.commentary[:80],
                    })

    print(f"\nMissed assumptions (score < 0.35): {missed} / {total_checked}")
    print(f"\n## Missed by Ontology")
    print(f"\n{'Ontology':25s} {'Missed':>6s} {'% Missed':>9s}")
    print("-" * 42)
    for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
        cnt = missed_by_onto.get(onto, 0)
        total_onto = onto_counts.get(onto, 0)
        pct = cnt / total_onto * 100 if total_onto > 0 else 0
        print(f"{onto:25s} {cnt:6d}  {pct:6.1f}%")

    print(f"\n## Sample Missed Assumptions by Ontology")
    for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
        examples = missed_examples.get(onto, [])
        if examples:
            print(f"\n  {onto} ({len(examples)} examples):")
            for ex in examples[:8]:
                print(f"    • {ex['text'][:90]}")
                print(f"      ({ex['type']}) {ex['commentary'][:70]}")

    # ── Save results ──
    output_dir = base / "ontology"
    output_dir.mkdir(exist_ok=True)

    # Save full classification
    with open(output_dir / "full_classification.json", "w") as f:
        json.dump(classifications, f, indent=2)

    # Save ontology summary
    summary = {
        "total_assumptions": total,
        "ontology_distribution": {k: v for k, v in sorted(onto_counts.items(), key=lambda x: -x[1])},
        "missed_by_ontology": {k: v for k, v in sorted(missed_by_onto.items(), key=lambda x: -x[1])},
        "cross_type_ontology": {t: dict(c) for t, c in sorted(type_x_onto.items())},
        "methodology": "Rule-based classification using keyword/pattern matching against assumption text and commentary. 8-type ontology: Explicit (directly written), Implicit (implied), Derived (logical consequence), Trust (relationship-based), Operational (process/procedure), Dependency (external), Architectural (design constraint), Environmental (infrastructure).",
    }
    with open(output_dir / "ontology_summary.json", "w") as f:
        json.dump(summary, f, indent=2)

    # Save markdown report
    with open(output_dir / "ASF_Assumption_Ontology.md", "w") as f:
        f.write("# ASF Assumption Ontology\n\n")
        f.write(f"**Total assumptions analyzed:** {total}\n\n")
        f.write("## Ontology Distribution\n\n")
        f.write("| Ontology | Count | Percent |\n")
        f.write("|----------|-------|--------|\n")
        for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
            cnt = onto_counts.get(onto, 0)
            pct = cnt / total * 100
            f.write(f"| {onto} | {cnt} | {pct:.1f}% |\n")
        f.write(f"| **TOTAL** | **{total}** | **100%** |\n")

        f.write("\n## Missed by Ontology\n\n")
        f.write("| Ontology | Missed | Total | % Missed |\n")
        f.write("|----------|--------|-------|----------|\n")
        for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
            cnt = missed_by_onto.get(onto, 0)
            total_onto = onto_counts.get(onto, 0)
            pct = cnt / total_onto * 100 if total_onto > 0 else 0
            f.write(f"| {onto} | {cnt} | {total_onto} | {pct:.1f}% |\n")

        f.write("\n## Sample Assumptions per Ontology\n\n")
        for onto in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]:
            examples = missed_examples.get(onto, [])
            if examples:
                f.write(f"### {onto}\n\n")
                for ex in examples[:10]:
                    f.write(f"- **{ex['type']}**: {ex['text']}\n")
                    f.write(f"  _{ex['commentary']}_\n\n")

        f.write("\n## Cross-tabulation: ASF Type × Ontology\n\n")
        f.write("| Type | Explicit | Implicit | Derived | Trust | Operational | Dependency | Architectural | Environmental | Total |\n")
        f.write("|------|----------|----------|---------|-------|-------------|------------|---------------|---------------|-------|\n")
        for asf_type in ["ACCESS", "IDENTITY", "NETWORK", "CONFIGURATION", "PROCESS", "GOVERNANCE", "DOCUMENTATION", "DEPENDENCY"]:
            if asf_type not in type_x_onto:
                continue
            c = type_x_onto[asf_type]
            row_total = sum(c.values())
            vals = [str(c.get(o, 0)) for o in ["Explicit", "Implicit", "Derived", "Trust", "Operational", "Dependency", "Architectural", "Environmental"]]
            f.write(f"| {asf_type} | {' | '.join(vals)} | {row_total} |\n")

    print(f"\n\nResults saved to {output_dir}/")


if __name__ == "__main__":
    run_analysis()
