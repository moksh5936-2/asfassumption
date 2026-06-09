from __future__ import annotations

import json
import sqlite3
from pathlib import Path
from threading import Lock
from typing import Any, Optional


class Database:
    def __init__(self, db_path: str | Path = "asf_validator.db"):
        self.db_path = Path(db_path)
        self._lock = Lock()
        self._conn: Optional[sqlite3.Connection] = None
        self._init_db()

    def _get_conn(self) -> sqlite3.Connection:
        if self._conn is None:
            self._conn = sqlite3.connect(str(self.db_path))
            self._conn.row_factory = sqlite3.Row
            self._conn.execute("PRAGMA journal_mode=WAL")
            self._conn.execute("PRAGMA foreign_keys=ON")
        return self._conn

    def _init_db(self) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.executescript("""
                CREATE TABLE IF NOT EXISTS claims (
                    id TEXT PRIMARY KEY,
                    source_document TEXT NOT NULL,
                    source_location TEXT,
                    text TEXT NOT NULL,
                    extraction_confidence REAL NOT NULL DEFAULT 0.5,
                    created_at TEXT NOT NULL,
                    tags TEXT DEFAULT '[]'
                );

                CREATE TABLE IF NOT EXISTS assumptions (
                    id TEXT PRIMARY KEY,
                    claim_id TEXT NOT NULL,
                    text TEXT NOT NULL,
                    assumption_type TEXT NOT NULL,
                    verification_status TEXT NOT NULL DEFAULT 'PENDING',
                    confidence REAL NOT NULL DEFAULT 0.0,
                    created_at TEXT NOT NULL,
                    keywords TEXT DEFAULT '[]'
                );

                CREATE TABLE IF NOT EXISTS evidence (
                    id TEXT PRIMARY KEY,
                    source TEXT NOT NULL,
                    source_type TEXT NOT NULL,
                    timestamp TEXT NOT NULL,
                    content TEXT,
                    confidence REAL NOT NULL DEFAULT 0.8,
                    metadata TEXT DEFAULT '{}',
                    records TEXT DEFAULT '[]'
                );

                CREATE TABLE IF NOT EXISTS verifications (
                    id TEXT PRIMARY KEY,
                    assumption_id TEXT NOT NULL,
                    evidence_used TEXT DEFAULT '[]',
                    result TEXT NOT NULL DEFAULT 'UNKNOWN',
                    confidence REAL NOT NULL DEFAULT 0.0,
                    reasoning TEXT DEFAULT '',
                    created_at TEXT NOT NULL,
                    details TEXT DEFAULT '{}'
                );

                CREATE TABLE IF NOT EXISTS gaps (
                    id TEXT PRIMARY KEY,
                    assumption_id TEXT NOT NULL,
                    severity TEXT NOT NULL,
                    type TEXT NOT NULL,
                    description TEXT NOT NULL,
                    evidence_detail TEXT DEFAULT '',
                    created_at TEXT NOT NULL
                );

                CREATE TABLE IF NOT EXISTS documents (
                    id TEXT PRIMARY KEY,
                    filename TEXT NOT NULL,
                    file_type TEXT NOT NULL,
                    content TEXT,
                    created_at TEXT NOT NULL
                );

                CREATE TABLE IF NOT EXISTS graph_edges (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    source_id TEXT NOT NULL,
                    target_id TEXT NOT NULL,
                    relationship TEXT NOT NULL,
                    UNIQUE(source_id, target_id, relationship)
                );
            """)
            conn.commit()

    def _serialize(self, obj: Any) -> str:
        if isinstance(obj, dict):
            return json.dumps({k: (v.isoformat() if hasattr(v, "isoformat") else v) for k, v in obj.items()})
        if isinstance(obj, list):
            return json.dumps(obj)
        return str(obj)

    def insert_claim(self, claim_dict: dict) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR REPLACE INTO claims
                   (id, source_document, source_location, text, extraction_confidence, created_at, tags)
                   VALUES (?, ?, ?, ?, ?, ?, ?)""",
                (
                    claim_dict["id"],
                    claim_dict["source_document"],
                    claim_dict.get("source_location"),
                    claim_dict["text"],
                    claim_dict["extraction_confidence"],
                    claim_dict["created_at"],
                    self._serialize(claim_dict.get("tags", [])),
                ),
            )
            conn.commit()

    def insert_assumption(self, assumption_dict: dict) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR REPLACE INTO assumptions
                   (id, claim_id, text, assumption_type, verification_status, confidence, created_at, keywords)
                   VALUES (?, ?, ?, ?, ?, ?, ?, ?)""",
                (
                    assumption_dict["id"],
                    assumption_dict["claim_id"],
                    assumption_dict["text"],
                    assumption_dict["assumption_type"],
                    assumption_dict["verification_status"],
                    assumption_dict["confidence"],
                    assumption_dict["created_at"],
                    self._serialize(assumption_dict.get("keywords", [])),
                ),
            )
            conn.commit()

    def insert_evidence(self, evidence_dict: dict) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR REPLACE INTO evidence
                   (id, source, source_type, timestamp, content, confidence, metadata, records)
                   VALUES (?, ?, ?, ?, ?, ?, ?, ?)""",
                (
                    evidence_dict["id"],
                    evidence_dict["source"],
                    evidence_dict["source_type"],
                    evidence_dict["timestamp"],
                    json.dumps(evidence_dict.get("content")) if not isinstance(evidence_dict.get("content"), str) else evidence_dict.get("content"),
                    evidence_dict["confidence"],
                    self._serialize(evidence_dict.get("metadata", {})),
                    self._serialize(evidence_dict.get("records", [])),
                ),
            )
            conn.commit()

    def insert_verification(self, verification_dict: dict) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR REPLACE INTO verifications
                   (id, assumption_id, evidence_used, result, confidence, reasoning, created_at, details)
                   VALUES (?, ?, ?, ?, ?, ?, ?, ?)""",
                (
                    verification_dict["id"],
                    verification_dict["assumption_id"],
                    self._serialize(verification_dict.get("evidence_used", [])),
                    verification_dict["result"],
                    verification_dict["confidence"],
                    verification_dict.get("reasoning", ""),
                    verification_dict["created_at"],
                    self._serialize(verification_dict.get("details", {})),
                ),
            )
            conn.commit()

    def insert_gap(self, gap_dict: dict) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR REPLACE INTO gaps
                   (id, assumption_id, severity, type, description, evidence_detail, created_at)
                   VALUES (?, ?, ?, ?, ?, ?, ?)""",
                (
                    gap_dict["id"],
                    gap_dict["assumption_id"],
                    gap_dict["severity"],
                    gap_dict["type"],
                    gap_dict["description"],
                    gap_dict.get("evidence_detail", ""),
                    gap_dict["created_at"],
                ),
            )
            conn.commit()

    def insert_edge(self, source_id: str, target_id: str, relationship: str) -> None:
        with self._lock:
            conn = self._get_conn()
            conn.execute(
                """INSERT OR IGNORE INTO graph_edges (source_id, target_id, relationship)
                   VALUES (?, ?, ?)""",
                (source_id, target_id, relationship),
            )
            conn.commit()

    def get_claims(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM claims ORDER BY created_at").fetchall()
            return [dict(r) for r in rows]

    def get_assumptions(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM assumptions ORDER BY created_at").fetchall()
            return [dict(r) for r in rows]

    def get_evidence(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM evidence ORDER BY timestamp").fetchall()
            return [dict(r) for r in rows]

    def get_verifications(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM verifications ORDER BY created_at").fetchall()
            return [dict(r) for r in rows]

    def get_gaps(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM gaps ORDER BY created_at").fetchall()
            return [dict(r) for r in rows]

    def get_edges(self) -> list[dict]:
        with self._lock:
            conn = self._get_conn()
            rows = conn.execute("SELECT * FROM graph_edges").fetchall()
            return [dict(r) for r in rows]

    def get_by_id(self, table: str, record_id: str) -> Optional[dict]:
        with self._lock:
            conn = self._get_conn()
            row = conn.execute(f"SELECT * FROM {table} WHERE id = ?", (record_id,)).fetchone()
            return dict(row) if row else None

    def close(self) -> None:
        if self._conn:
            self._conn.close()
            self._conn = None

    def clear_all(self) -> None:
        with self._lock:
            conn = self._get_conn()
            for table in ["claims", "assumptions", "evidence", "verifications", "gaps", "documents", "graph_edges"]:
                conn.execute(f"DELETE FROM {table}")
            conn.commit()
