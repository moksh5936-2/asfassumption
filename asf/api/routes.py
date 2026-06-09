from __future__ import annotations

import json
import tempfile
from pathlib import Path
from typing import Optional

from fastapi import APIRouter, File, Form, HTTPException, UploadFile
from fastapi.responses import HTMLResponse, JSONResponse

from asf.analyzer import Analyzer
from asf.config import ASFConfig
from asf.db.database import Database
from asf.models import AnalysisResult, VerificationResult, GapSeverity

router = APIRouter(prefix="/api/v1", tags=["asf"])
_db: Optional[Database] = None
_analyzer: Optional[Analyzer] = None
_config: Optional[ASFConfig] = None


def init(db: Database, analyzer: Analyzer, config: Optional[ASFConfig] = None) -> None:
    global _db, _analyzer, _config
    _db = db
    _analyzer = analyzer
    _config = config


@router.post("/documents")
async def upload_document(file: UploadFile = File(...)):
    if _analyzer is None:
        raise HTTPException(503, "Analyzer not initialized")

    suffix = Path(file.filename or "upload").suffix
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        content = await file.read()
        tmp.write(content)
        tmp_path = tmp.name

    try:
        text = _analyzer.pipeline.parse_text(tmp_path)
        doc_meta = _analyzer.pipeline.get_document_metadata(tmp_path)
        claims = _analyzer.claim_extractor.extract(text, source_document=doc_meta["filename"])
        return {
            "document": doc_meta,
            "claims_found": len(claims),
            "claims": [c.model_dump() for c in claims],
        }
    finally:
        Path(tmp_path).unlink(missing_ok=True)


@router.post("/evidence")
async def upload_evidence(file: UploadFile = File(...)):
    if _analyzer is None:
        raise HTTPException(503, "Analyzer not initialized")

    suffix = Path(file.filename or "upload").suffix
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        content = await file.read()
        tmp.write(content)
        tmp_path = tmp.name

    try:
        evidence = _analyzer.evidence_loader.load(tmp_path, auto_map=True)
        return {"evidence": evidence.model_dump()}
    finally:
        Path(tmp_path).unlink(missing_ok=True)


@router.post("/analyze")
async def run_analysis(documents: list[str] = Form(...), evidence: Optional[list[str]] = Form(None)):
    if _analyzer is None:
        raise HTTPException(503, "Analyzer not initialized")

    doc_paths = [p for p in documents if Path(p).exists()]
    ev_paths = [p for p in (evidence or []) if Path(p).exists()]

    result = _analyzer.analyze(doc_paths, ev_paths, persist=True)
    return _serialize_result(result)


@router.get("/assumptions")
async def get_assumptions():
    if _db is None:
        raise HTTPException(503, "Database not initialized")
    return {"assumptions": _db.get_assumptions()}


@router.get("/gaps")
async def get_gaps():
    if _db is None:
        raise HTTPException(503, "Database not initialized")
    return {"gaps": _db.get_gaps()}


@router.get("/claims")
async def get_claims():
    if _db is None:
        raise HTTPException(503, "Database not initialized")
    return {"claims": _db.get_claims()}


@router.get("/graph")
async def get_graph():
    if _analyzer is None:
        raise HTTPException(503, "Analyzer not initialized")
    return JSONResponse(content=_analyzer.graph_model.export_json())


@router.get("/summary")
async def get_summary():
    if _db is None:
        raise HTTPException(503, "Database not initialized")
    claims = _db.get_claims()
    assumptions = _db.get_assumptions()
    verifications = _db.get_verifications()
    gaps = _db.get_gaps()

    return {
        "summary": {
            "claims_found": len(claims),
            "assumptions": len(assumptions),
            "verified": sum(1 for v in verifications if v["result"] == VerificationResult.VERIFIED),
            "contradicted": sum(1 for v in verifications if v["result"] == VerificationResult.CONTRADICTED),
            "unknown": sum(1 for v in verifications if v["result"] == VerificationResult.UNKNOWN),
            "critical_gaps": sum(1 for g in gaps if g["severity"] == GapSeverity.CRITICAL),
        }
    }


def _serialize_result(result: AnalysisResult) -> dict:
    return {
        "summary": {
            "claims_found": result.claims_found,
            "assumptions": result.assumptions_found,
            "verified": result.verified_count,
            "contradicted": result.contradicted_count,
            "unknown": result.unknown_count,
            "critical_gaps": result.critical_gaps,
        },
        "claims": [c.model_dump() for c in result.claims],
        "assumptions": [a.model_dump() for a in result.assumptions],
        "verifications": [
            {
                "assumption_id": v.assumption_id,
                "result": v.result.value,
                "confidence": v.confidence,
                "evidence_used": v.evidence_used,
                "reasoning": v.reasoning,
            }
            for v in result.verifications
        ],
        "gaps": [
            {
                "assumption_id": g.assumption_id,
                "type": g.type.value,
                "severity": g.severity.value,
                "description": g.description,
                "evidence_detail": g.evidence_detail,
            }
            for g in result.gaps
        ],
    }
