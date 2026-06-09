from __future__ import annotations

import json
from typing import Any

import networkx as nx

from asf.models import (
    AnalysisResult,
    Assumption,
    Claim,
    Evidence,
    Gap,
    Verification,
    VerificationResult,
)


class GraphModel:
    def __init__(self):
        self.graph = nx.MultiDiGraph()

    def build(self, result: AnalysisResult) -> None:
        self.graph.clear()

        for claim in result.claims:
            self.graph.add_node(
                claim.id,
                type="Claim",
                label=claim.text[:60] + ("..." if len(claim.text) > 60 else ""),
                full_text=claim.text,
                source_document=claim.source_document,
                extraction_confidence=claim.extraction_confidence,
            )

        for assumption in result.assumptions:
            self.graph.add_node(
                assumption.id,
                type="Assumption",
                label=assumption.text[:60] + ("..." if len(assumption.text) > 60 else ""),
                full_text=assumption.text,
                assumption_type=assumption.assumption_type.value,
                verification_status=assumption.verification_status.value,
                confidence=assumption.confidence,
            )
            self.graph.add_edge(assumption.claim_id, assumption.id, relationship="GENERATES")

        for evidence in result.evidence:
            self.graph.add_node(
                evidence.id,
                type="Evidence",
                label=evidence.source.split("/")[-1] if "/" in evidence.source else evidence.source,
                source_type=evidence.source_type.value,
                record_count=len(evidence.records),
                confidence=evidence.confidence,
            )

        for verification in result.verifications:
            self.graph.add_node(
                verification.id,
                type="Verification",
                label=f"Verification: {verification.result.value}",
                result=verification.result.value,
                confidence=verification.confidence,
            )
            self.graph.add_edge(verification.assumption_id, verification.id, relationship="VERIFIES")
            for ev_id in verification.evidence_used:
                self.graph.add_edge(ev_id, verification.id, relationship="SUPPORTS")

            assumption = next(
                (a for a in result.assumptions if a.id == verification.assumption_id),
                None,
            )
            if assumption:
                if verification.result == VerificationResult.VERIFIED:
                    rel = "SUPPORTS"
                elif verification.result == VerificationResult.CONTRADICTED:
                    rel = "CONTRADICTS"
                else:
                    rel = "RELATES_TO"
                self.graph.add_edge(verification.id, assumption.id, relationship=rel)

        for gap in result.gaps:
            self.graph.add_node(
                gap.id,
                type="Gap",
                label=f"Gap: {gap.type.value} ({gap.severity.value})",
                gap_type=gap.type.value,
                severity=gap.severity.value,
                description=gap.description[:80],
            )
            self.graph.add_edge(gap.assumption_id, gap.id, relationship="IDENTIFIES")

    def export_json(self) -> dict[str, Any]:
        nodes_data = []
        for node_id, attrs in self.graph.nodes(data=True):
            node_dict = {"id": node_id, **attrs}
            node_dict["node_type"] = node_dict.pop("type", "Unknown")
            nodes_data.append(node_dict)

        edges_data = []
        for u, v, key, attrs in self.graph.edges(data=True, keys=True):
            edges_data.append({
                "source": u,
                "target": v,
                "key": key,
                **attrs,
            })

        return {
            "nodes": nodes_data,
            "edges": edges_data,
            "node_count": self.graph.number_of_nodes(),
            "edge_count": self.graph.number_of_edges(),
        }

    def export_json_string(self, indent: int = 2) -> str:
        return json.dumps(self.export_json(), indent=indent)

    def summary(self) -> dict[str, int]:
        node_types: dict[str, int] = {}
        for _, attrs in self.graph.nodes(data=True):
            ntype = attrs.get("type", "Unknown")
            node_types[ntype] = node_types.get(ntype, 0) + 1
        return {
            "total_nodes": self.graph.number_of_nodes(),
            "total_edges": self.graph.number_of_edges(),
            "node_types": node_types,
        }
