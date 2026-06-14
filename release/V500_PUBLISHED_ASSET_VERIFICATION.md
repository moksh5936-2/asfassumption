# Published Asset Verification — ASF0 v5.0.0

## Verification Steps

1. **`gh release view v5.0.0`** confirmed:
   - ✅ 8 assets uploaded
   - ✅ Published (not draft)
   - ✅ Not prerelease
   - ✅ All SHA256 digests match local checksums.txt

2. **Downloaded native binary** from release URL:
   ```bash
   curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-darwin-arm64
   ```
   - ✅ `--version` shows `ASF0 v5.0.0`
   - ✅ SHA256 matches local checksums.txt: `66b0102849feb2c78a19422e94cda5258c6a659bbc37150029ae7a5b128f33cf`

## Result
✅ All 8 published assets are verified and match local artifacts.
