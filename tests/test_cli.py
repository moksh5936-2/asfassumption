from pathlib import Path
from click.testing import CliRunner

from asf.cli.main import cli

SAMPLE_DIR = Path(__file__).parent.parent / "sample_data"


class TestCLI:
    def setup_method(self):
        self.runner = CliRunner()

    def test_analyze_help(self):
        result = self.runner.invoke(cli, ["analyze", "--help"])
        assert result.exit_code == 0
        assert "PATHS" in result.output

    def test_analyze_with_data(self):
        doc = str(SAMPLE_DIR / "finance_policy.txt")
        ev = str(SAMPLE_DIR / "payroll_acl.csv")
        result = self.runner.invoke(cli, ["analyze", doc, "-e", ev])
        assert result.exit_code == 0
        assert "Documents Processed" in result.output
        assert "Assumptions" in result.output

    def test_analyze_json_output(self):
        doc = str(SAMPLE_DIR / "finance_policy.txt")
        result = self.runner.invoke(cli, ["analyze", doc, "--json"])
        assert result.exit_code == 0
        import json
        data = json.loads(result.output)
        assert "summary" in data
        assert data["summary"]["claims_found"] > 0
