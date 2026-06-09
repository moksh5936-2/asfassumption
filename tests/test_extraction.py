from asf.extraction import ClaimExtractor


class TestClaimExtraction:
    def setup_method(self):
        self.extractor = ClaimExtractor()

    def test_extract_access_claim(self):
        text = "Only Finance employees may access the payroll processing system."
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) >= 1
        assert "Finance" in claims[0].text

    def test_extract_encryption_claim(self):
        text = "All payroll data is encrypted at rest and in transit."
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) >= 1
        assert "encrypted" in claims[0].text.lower()

    def test_extract_multi_claims(self):
        text = (
            "Only Finance employees may access the payroll system. "
            "All data is encrypted. "
            "Production databases are not internet accessible."
        )
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) >= 2

    def test_no_claims_for_narrative(self):
        text = "We had a meeting about security. The team discussed various options."
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) == 0

    def test_deduplication(self):
        text = "Only Finance can access payroll. Only Finance can access payroll."
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) == 1

    def test_tags_extracted(self):
        text = "MFA is required for all administrative access."
        claims = self.extractor.extract(text, source_document="test.txt")
        assert len(claims) >= 1
        assert "identity" in claims[0].tags or "access" in claims[0].tags

    def test_full_policy_extraction(self):
        text = (
            "ONLY FINANCE ACCESS. "
            "ENCRYPTED AT REST. "
            "LOGGED AND MONITORED. "
            "QUARTERLY REVIEWS. "
            "BACKED UP DAILY."
        )
        claims = self.extractor.extract(text, source_document="policy.txt")
        assert len(claims) >= 3
