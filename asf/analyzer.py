from __future__ import annotations

from pathlib import Path
from typing import Optional, Sequence

from asf.assumption import AssumptionEngine
from asf.config import ASFConfig
from asf.confidence import ConfidenceEngine
from asf.db.database import Database
from asf.evidence import EvidenceLoader, EvidenceMapper
from asf.extraction import ClaimExtractor
from asf.gaps import GapEngine
from asf.graph import GraphModel
from asf.ingestion import IngestionPipeline
from asf.models import (
    AnalysisResult,
    Assumption,
    Claim,
    Evidence,
    VerificationResult,
)
from asf.verification import VerificationEngine


class Analyzer:
    def __init__(self, config: ASFConfig | None = None):
        self.config = config or ASFConfig.default()
        self.pipeline = IngestionPipeline()
        self.claim_extractor = ClaimExtractor()
        self.assumption_engine = AssumptionEngine()
        self.evidence_loader = EvidenceLoader(
            field_mappings=self.config.field_mappings or None
        )
        self.evidence_mapper = EvidenceMapper()
        self.verification_engine = VerificationEngine()
        self.confidence_engine = ConfidenceEngine()
        self.gap_engine = GapEngine()
        self.graph_model = GraphModel()
        self.result = AnalysisResult()
        self.db: Optional[Database] = None

        if self.config.db_path:
            self.db = Database(self.config.db_path)

    def analyze(
        self,
        document_paths: list[str | Path],
        evidence_paths: list[str | Path] | None = None,
        persist: bool = False,
    ) -> AnalysisResult:
        self.result = AnalysisResult()

        claims = self._process_documents(document_paths)
        assumptions = self.assumption_engine.convert_many(claims)

        evidence_records = []
        if evidence_paths:
            evidence_records = self._load_evidence(evidence_paths)

        evidence_list = list(evidence_records)

        verifications = []
        for assumption in assumptions:
            matching_evidence = self._find_matching_evidence(assumption, evidence_list)
            verification = self.verification_engine.verify(assumption, matching_evidence)
            verification.confidence = self.confidence_engine.compute_verification_confidence(
                verification, matching_evidence
            )
            verifications.append(verification)

            assumption.confidence = self.confidence_engine.compute_assumption_confidence([verification])
            if verification.result == VerificationResult.VERIFIED:
                from asf.models import VerificationStatus
                assumption.verification_status = VerificationStatus.VERIFIED
            elif verification.result == VerificationResult.CONTRADICTED:
                from asf.models import VerificationStatus
                assumption.verification_status = VerificationStatus.CONTRADICTED
            elif verification.result == VerificationResult.PARTIALLY_VERIFIED:
                from asf.models import VerificationStatus
                assumption.verification_status = VerificationStatus.IN_REVIEW

        gaps = self.gap_engine.generate_gaps(assumptions, verifications)

        self.result = AnalysisResult(
            claims=claims,
            assumptions=assumptions,
            evidence=evidence_list,
            verifications=verifications,
            gaps=gaps,
        )

        self.graph_model.build(self.result)

        if persist and self.db:
            self._persist()

        return self.result

    def _process_documents(self, paths: list[str | Path]) -> list[Claim]:
        all_claims: list[Claim] = []
        for path in paths:
            path = Path(path)
            if not path.exists():
                continue
            text = self.pipeline.parse_text(path)
            doc_meta = self.pipeline.get_document_metadata(path)
            claims = self.claim_extractor.extract(
                text,
                source_document=doc_meta["filename"],
                source_location=str(path),
            )
            all_claims.extend(claims)
        return all_claims

    def _load_evidence(self, paths: list[str | Path]) -> list[Evidence]:
        evidence_list: list[Evidence] = []
        for path in paths:
            path = Path(path)
            if not path.exists():
                continue
            evidence = self.evidence_loader.load(path, auto_map=True)
            evidence_list.append(evidence)
        return evidence_list

    def _find_matching_evidence(self, assumption: Assumption, evidence_list: list[Evidence]) -> list[Evidence]:
        compatible_types = self.evidence_mapper.get_compatible_source_types(assumption.assumption_type)
        if not compatible_types:
            return evidence_list
        matched = [e for e in evidence_list if e.source_type in compatible_types]
        if not matched:
            return evidence_list
        return matched

    def _persist(self) -> None:
        if not self.db:
            return
        for claim in self.result.claims:
            self.db.insert_claim(claim.model_dump())
        for assumption in self.result.assumptions:
            self.db.insert_assumption(assumption.model_dump())
        for ev in self.result.evidence:
            self.db.insert_evidence(ev.model_dump())
        for v in self.result.verifications:
            self.db.insert_verification(v.model_dump())
        for gap in self.result.gaps:
            self.db.insert_gap(gap.model_dump())

        for a in self.result.assumptions:
            self.db.insert_edge(a.claim_id, a.id, "GENERATES")
        for v in self.result.verifications:
            self.db.insert_edge(v.assumption_id, v.id, "VERIFIES")
            for ev_id in v.evidence_used:
                self.db.insert_edge(ev_id, v.id, "SUPPORTS")
        for g in self.result.gaps:
            self.db.insert_edge(g.assumption_id, g.id, "IDENTIFIES")

    def close(self) -> None:
        if self.db:
            self.db.close()
