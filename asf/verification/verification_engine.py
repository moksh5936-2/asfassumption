from __future__ import annotations

import json
import re
from datetime import datetime, timezone
from typing import Any

from asf.models import (
    Assumption,
    AssumptionType,
    Evidence,
    Verification,
    VerificationResult,
)


class VerificationEngine:
    def verify(self, assumption: Assumption, evidence_records: list[Evidence]) -> Verification:
        matched_evidence_ids: list[str] = []
        result = VerificationResult.UNKNOWN
        confidence = 0.0
        reasoning_parts: list[str] = []
        details: dict[str, Any] = {}

        for evidence in evidence_records:
            if not evidence.records:
                reasoning_parts.append(f"No structured records in {evidence.source}")
                continue

            matched_evidence_ids.append(evidence.id)
            check_result = self._check_assumption_against_evidence(assumption, evidence)

            if check_result:
                vres, vconf, vreason, vdetails = check_result
                reasoning_parts.append(vreason)
                details[evidence.source] = vdetails

                if vres == VerificationResult.VERIFIED:
                    if result != VerificationResult.CONTRADICTED:
                        result = VerificationResult.VERIFIED
                    confidence = max(confidence, vconf)
                elif vres == VerificationResult.CONTRADICTED:
                    result = VerificationResult.CONTRADICTED
                    confidence = max(confidence, vconf)
                elif vres == VerificationResult.PARTIALLY_VERIFIED:
                    if result not in (VerificationResult.CONTRADICTED, VerificationResult.VERIFIED):
                        result = VerificationResult.PARTIALLY_VERIFIED
                    confidence = max(confidence, vconf)

        if not matched_evidence_ids:
            reasoning_parts.append("No matching evidence available for verification")

        verification = Verification(
            assumption_id=assumption.id,
            evidence_used=matched_evidence_ids,
            result=result,
            confidence=confidence,
            reasoning="; ".join(reasoning_parts) if reasoning_parts else "No evidence processed",
            details=details,
        )
        return verification

    def _check_assumption_against_evidence(
        self, assumption: Assumption, evidence: Evidence
    ) -> tuple[VerificationResult, float, str, dict] | None:
        atype = assumption.assumption_type

        if atype == AssumptionType.ACCESS:
            return self._check_access(assumption, evidence)
        elif atype == AssumptionType.IDENTITY:
            return self._check_identity(assumption, evidence)
        elif atype == AssumptionType.NETWORK:
            return self._check_network(assumption, evidence)
        elif atype == AssumptionType.CONFIGURATION:
            return self._check_configuration(assumption, evidence)
        elif atype == AssumptionType.GOVERNANCE:
            return self._check_governance(assumption, evidence)
        elif atype == AssumptionType.PROCESS:
            return self._check_process(assumption, evidence)

        return None

    def _check_access(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        text_lower = assumption.text.lower()
        records = evidence.records
        resource_keywords = self._extract_resource_keywords(assumption.text)

        only_pattern = re.search(r"only\s+(.+?)\s+(?:can|may|has|should)", text_lower)
        restricted_group = only_pattern.group(1).strip() if only_pattern else None

        users_outside = []
        users_inside = []
        resources_found = set()

        for rec in records:
            rec_lower = {k.lower(): str(v).lower() for k, v in rec.items()}

            user = self._find_field(rec_lower, ["user", "username", "identity", "principal", "name", "email"])
            resource = self._find_field(rec_lower, ["resource", "application", "system", "target", "service", "scope"])
            permission = self._find_field(rec_lower, ["permission", "access", "role", "right", "privilege", "action"])
            group = self._find_field(rec_lower, ["group", "department", "team", "unit", "division", "org"])

            if resource:
                resources_found.add(resource)
            if user and permission:
                if restricted_group:
                    user_in_group = restricted_group and group and (group in restricted_group or restricted_group in group)
                    if user_in_group:
                        users_inside.append(user)
                    else:
                        users_outside.append(user)

        details = {
            "expected_group": restricted_group,
            "users_outside_group": list(set(users_outside)),
            "users_inside_group": list(set(users_inside)),
            "resources_found": list(resources_found),
            "total_records": len(records),
        }

        if users_outside and restricted_group:
            return (
                VerificationResult.CONTRADICTED,
                0.92,
                f"Found {len(set(users_outside))} user(s) outside '{restricted_group}' with access: {', '.join(set(users_outside[:5]))}",
                details,
            )

        if restricted_group and users_inside:
            return (
                VerificationResult.VERIFIED,
                0.78,
                f"Only users in '{restricted_group}' found with access ({len(set(users_inside))} users)",
                details,
            )

        return (
            VerificationResult.UNKNOWN,
            0.3,
            "Could not determine access patterns from evidence",
            details,
        )

    def _check_identity(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        text_lower = assumption.text.lower()
        records = evidence.records
        details: dict[str, Any] = {"checks": []}

        has_mfa_claim = "mfa" in text_lower or "multi-factor" in text_lower or "multifactor" in text_lower

        if has_mfa_claim:
            mfa_users = []
            no_mfa_users = []
            for rec in records:
                rec_lower = {k.lower(): str(v).lower() for k, v in rec.items()}
                user = self._find_field(rec_lower, ["user", "username", "identity", "name", "email"])
                mfa_status = self._find_field(rec_lower, ["mfa", "mfa_enabled", "multi_factor", "2fa", "totp"])

                if user:
                    if mfa_status and mfa_status in ("true", "yes", "enabled", "1"):
                        mfa_users.append(user)
                    else:
                        no_mfa_users.append(user)

            details["mfa_enabled_users"] = mfa_users
            details["mfa_disabled_users"] = no_mfa_users

            if no_mfa_users and not mfa_users:
                return (
                    VerificationResult.CONTRADICTED,
                    0.95,
                    f"MFA not enabled for {len(no_mfa_users)} user(s)",
                    details,
                )
            if no_mfa_users and mfa_users:
                return (
                    VerificationResult.PARTIALLY_VERIFIED,
                    0.6,
                    f"MFA enabled for {len(mfa_users)} user(s) but missing for {len(no_mfa_users)}",
                    details,
                )
            if mfa_users and not no_mfa_users:
                return (
                    VerificationResult.VERIFIED,
                    0.85,
                    f"MFA enabled for all {len(mfa_users)} user(s)",
                    details,
                )

        return (
            VerificationResult.UNKNOWN,
            0.3,
            "No identity evidence matched assumption",
            details,
        )

    def _check_network(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        text_lower = assumption.text.lower()
        records = evidence.records
        details: dict[str, Any] = {}

        is_exposed = "internet" in text_lower or "public" in text_lower or "exposed" in text_lower or "external" in text_lower
        is_isolated = "isolat" in text_lower or "segment" in text_lower or "private" in text_lower

        exposures = []
        isolations = []
        for rec in records:
            rec_lower = {k.lower(): str(v).lower() for k, v in rec.items()}
            asset = self._find_field(rec_lower, ["asset", "resource", "host", "server", "system", "name", "service"])
            public_val = self._find_field(rec_lower, ["public", "exposed", "internet_facing", "is_public", "exposure"])

            if asset:
                if public_val and public_val in ("true", "yes", "exposed", "1"):
                    exposures.append(asset)
                else:
                    isolations.append(asset)

        details["exposed_assets"] = exposures
        details["isolated_assets"] = isolations

        if is_isolated and exposures:
            return (
                VerificationResult.CONTRADICTED,
                0.9,
                f"Claimed isolated but found {len(exposures)} exposed asset(s): {', '.join(exposures[:5])}",
                details,
            )

        if is_isolated and not exposures:
            return (
                VerificationResult.VERIFIED,
                0.8,
                f"All {len(isolations)} asset(s) appear isolated",
                details,
            )

        if is_exposed and not exposures and isolations:
            negation = "no" in text_lower or "not" in text_lower or "never" in text_lower
            if negation:
                return (
                    VerificationResult.VERIFIED,
                    0.85,
                    f"Confirmed: no exposure found across {len(isolations)} asset(s)",
                    details,
                )

        if is_exposed and exposures:
            return (
                VerificationResult.VERIFIED,
                0.85,
                f"Found {len(exposures)} exposed asset(s) as expected",
                details,
            )

        return (
            VerificationResult.UNKNOWN,
            0.3,
            "Could not verify network posture",
            details,
        )

    def _check_configuration(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        text_lower = assumption.text.lower()
        records = evidence.records
        details: dict[str, Any] = {}

        is_encrypted = "encrypt" in text_lower
        is_backed_up = "backup" in text_lower
        is_logged = "log" in text_lower or "audit" in text_lower

        compliant = 0
        non_compliant = 0
        examples_compliant = []
        examples_non_compliant = []

        for rec in records:
            rec_lower = {k.lower(): str(v).lower() for k, v in rec.items()}
            resource = self._find_field(rec_lower, ["resource", "system", "asset", "service", "name", "component"])
            enabled = self._find_field(rec_lower, ["enabled", "status", "state", "active", "value", "configuration"])

            if resource:
                if enabled and enabled in ("true", "yes", "enabled", "1", "active", "on"):
                    compliant += 1
                    if len(examples_compliant) < 3:
                        examples_compliant.append(resource)
                elif enabled and enabled in ("false", "no", "disabled", "0", "inactive", "off"):
                    non_compliant += 1
                    if len(examples_non_compliant) < 3:
                        examples_non_compliant.append(resource)

        details["compliant"] = compliant
        details["non_compliant"] = non_compliant
        details["examples_compliant"] = examples_compliant
        details["examples_non_compliant"] = examples_non_compliant

        if non_compliant > 0 and compliant == 0:
            return (
                VerificationResult.CONTRADICTED,
                0.9,
                f"Configuration not applied: {non_compliant} non-compliant resource(s)",
                details,
            )

        if non_compliant > 0 and compliant > 0:
            return (
                VerificationResult.PARTIALLY_VERIFIED,
                0.5,
                f"Partially compliant: {compliant} OK, {non_compliant} non-compliant",
                details,
            )

        if compliant > 0 and non_compliant == 0:
            return (
                VerificationResult.VERIFIED,
                0.85,
                f"All {compliant} resource(s) compliant with configuration",
                details,
            )

        return (
            VerificationResult.UNKNOWN,
            0.3,
            "Could not verify configuration from evidence",
            details,
        )

    def _check_governance(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        records = evidence.records
        details: dict[str, Any] = {}

        reviewed = 0
        not_reviewed = 0
        for rec in records:
            rec_lower = {k.lower(): str(v).lower() for k, v in rec.items()}
            status = self._find_field(rec_lower, ["status", "reviewed", "approved", "state", "completed"])
            if status and status in ("true", "yes", "completed", "approved", "reviewed", "1"):
                reviewed += 1
            else:
                not_reviewed += 1

        details["reviews_completed"] = reviewed
        details["reviews_pending"] = not_reviewed

        if not_reviewed > 0 and reviewed == 0:
            return (
                VerificationResult.CONTRADICTED,
                0.88,
                f"No governance reviews completed ({not_reviewed} pending)",
                details,
            )
        if not_reviewed > 0:
            return (
                VerificationResult.PARTIALLY_VERIFIED,
                0.55,
                f"{reviewed} reviews done, {not_reviewed} pending",
                details,
            )
        if reviewed > 0:
            return (
                VerificationResult.VERIFIED,
                0.85,
                f"All {reviewed} governance reviews completed",
                details,
            )

        return (VerificationResult.UNKNOWN, 0.3, "No governance evidence", details)

    def _check_process(self, assumption: Assumption, evidence: Evidence) -> tuple[VerificationResult, float, str, dict]:
        return self._check_governance(assumption, evidence)

    def _extract_resource_keywords(self, text: str) -> list[str]:
        words = re.findall(r"\b[A-Z][a-z]+\b", text)
        return words[:5]

    def _find_field(self, record: dict[str, str], candidates: list[str]) -> str | None:
        for candidate in candidates:
            if candidate in record:
                return record[candidate]
        return None
