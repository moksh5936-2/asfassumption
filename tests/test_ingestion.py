from pathlib import Path

from asf.ingestion import IngestionPipeline
from asf.models import SourceType


SAMPLE_DIR = Path(__file__).parent.parent / "sample_data"


class TestIngestion:
    def setup_method(self):
        self.pipeline = IngestionPipeline()

    def test_detect_txt(self):
        ftype = self.pipeline.detect_type(SAMPLE_DIR / "finance_policy.txt")
        assert ftype == SourceType.TXT

    def test_detect_csv(self):
        ftype = self.pipeline.detect_type(SAMPLE_DIR / "payroll_acl.csv")
        assert ftype == SourceType.CSV

    def test_detect_json(self):
        ftype = self.pipeline.detect_type(SAMPLE_DIR / "iam_export.json")
        assert ftype == SourceType.JSON

    def test_parse_txt(self):
        text = self.pipeline.parse_text(SAMPLE_DIR / "finance_policy.txt")
        assert "Only Finance employees may access the payroll processing system" in text
        assert len(text) > 100

    def test_parse_csv_to_text(self):
        text = self.pipeline.parse_text(SAMPLE_DIR / "payroll_acl.csv")
        assert len(text) > 0

    def test_parse_csv_records(self):
        records = self.pipeline.parse_to_records(SAMPLE_DIR / "payroll_acl.csv")
        assert len(records) == 10
        assert records[0]["user"] == "alice.jones"

    def test_parse_json_records(self):
        records = self.pipeline.parse_to_records(SAMPLE_DIR / "iam_export.json")
        assert len(records) == 1  # single dict becomes list of 1
        assert "users" in records[0]

    def test_document_metadata(self):
        meta = self.pipeline.get_document_metadata(SAMPLE_DIR / "finance_policy.txt")
        assert meta["filename"] == "finance_policy.txt"
        assert meta["file_type"] == "TXT"
        assert meta["size_bytes"] > 0
