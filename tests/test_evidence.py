from asf.evidence import EvidenceMapper, EvidenceLoader
from asf.models import AssumptionType, SourceType


class TestEvidenceMapper:
    def setup_method(self):
        self.mapper = EvidenceMapper()

    def test_get_required_evidence_access(self):
        req = self.mapper.get_required_evidence(AssumptionType.ACCESS)
        assert len(req) >= 2
        types = [r["source_type"] for r in req]
        assert SourceType.ACL_LIST in types or SourceType.IAM_EXPORT in types

    def test_get_required_evidence_network(self):
        req = self.mapper.get_required_evidence(AssumptionType.NETWORK)
        assert len(req) >= 2
        types = [r["source_type"] for r in req]
        assert SourceType.FIREWALL_RULES in types or SourceType.SECURITY_GROUPS in types

    def test_get_compatible_types(self):
        types = self.mapper.get_compatible_source_types(AssumptionType.IDENTITY)
        assert SourceType.IAM_EXPORT in types

    def test_get_required_identity(self):
        req = self.mapper.get_required_evidence(AssumptionType.IDENTITY)
        assert len(req) >= 1


class TestEvidenceLoader:
    def setup_method(self):
        self.loader = EvidenceLoader()

    def test_load_csv(self):
        import tempfile
        import csv
        tmp = tempfile.NamedTemporaryFile(mode="w", suffix=".csv", delete=False)
        w = csv.writer(tmp)
        w.writerow(["user", "role", "access"])
        w.writerow(["alice", "admin", "true"])
        tmp.close()

        evidence = self.loader.load(tmp.name)
        assert evidence.source_type == SourceType.CSV
        assert len(evidence.records) == 1
        import os
        os.unlink(tmp.name)

    def test_load_json(self):
        import tempfile
        import json
        tmp = tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False)
        json.dump([{"user": "alice", "access": "admin"}], tmp)
        tmp.close()

        evidence = self.loader.load(tmp.name)
        assert evidence.source_type == SourceType.JSON
        import os
        os.unlink(tmp.name)
