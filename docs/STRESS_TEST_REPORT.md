# Stress Test Report

**Phase 7 — June 2026**

## Summary

The Go native engine was stress-tested with large inputs, many evidence files, and pathological line lengths. All tests completed without errors and produced results matching the Python engine.

## Test Results

### 1. Large TXT File (134KB, 100× Concatenated Policy)

**Setup:** The finance policy text was concatenated 100 times into a single 134KB file.

| Metric | Go | Python |
|--------|----|--------|
| Claims found | 17 | 17 |
| Errors | None | None |
| Completion time | <1s | <1s |

**Notes:** The claim extractor correctly deduplicates identical sentences, so repeated policy text produces the same 17 claims as a single pass. No memory pressure or slowdown observed.

### 2. Many Evidence Files (25 Files)

**Setup:** 25 CSV evidence files were provided via `-e` flag, containing various access control, identity, and network configuration data.

| Metric | Go | Python |
|--------|----|--------|
| Claims found | 17 | 17 |
| Verified | 1 | 1 |
| Contradicted | 5 | 5 |
| Errors | None | None |

**Notes:** Evidence loading scaled linearly. All 25 files were parsed and matched against assumptions correctly. Results identical between engines.

### 3. Single Long Line (20.8K Characters)

**Setup:** A single line of 20,800 characters containing repeated policy sentences with no newlines.

| Metric | Go | Python |
|--------|----|--------|
| Claims found | 1 | 1 |
| Errors | None | None |
| Claim text | Single deduplicated sentence | Single deduplicated sentence |

**Notes:** Both engines correctly identified one unique claim from the repeated text. No buffer overflow, truncation, or performance issues.

## Conclusion

The Go native engine handles stress conditions **identically** to the Python engine. No errors, crashes, or differences observed. The engine is suitable for production workloads with large files and many evidence sources.
