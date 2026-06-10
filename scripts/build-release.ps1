# ASF Release Build Script (PowerShell)
# Usage: .\scripts\build-release.ps1 -Version "1.0.0"
# Requires: Go 1.24+

param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$ReleaseDir = Join-Path $RootDir "release"
$BuildDir = Join-Path $RootDir "asf-tui"

$Platforms = @(
    @{OS="linux"; Arch="amd64"}
    @{OS="linux"; Arch="arm64"}
    @{OS="darwin"; Arch="amd64"}
    @{OS="darwin"; Arch="arm64"}
    @{OS="windows"; Arch="amd64"}
)

Write-Host "=== ASF Release Build v$Version ===" -ForegroundColor Cyan

# Check Go
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "ERROR: Go is not installed. Install Go 1.24+ to build." -ForegroundColor Red
    exit 1
}
Write-Host "Go: $goVersion"

# Clean
Write-Host "`nCleaning release directory..."
Get-ChildItem $ReleaseDir -Filter "ASF-v*" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "asf-*" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.tar.gz" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.zip" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.exe" | Remove-Item -Force

# Build
Write-Host "`nBuilding for all platforms..."
Push-Location $BuildDir

foreach ($platform in $Platforms) {
    $outName = "ASF-v$Version-$($platform.OS)-$($platform.Arch)"
    if ($platform.OS -eq "windows") {
        $outName += ".exe"
    }
    $outPath = Join-Path $ReleaseDir $outName

    Write-Host "  Building $($platform.OS)/$($platform.Arch)..."
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    $env:CGO_ENABLED = "0"
    go build -ldflags="-s -w" -o $outPath .
    
    $size = (Get-Item $outPath).Length / 1MB
    Write-Host "  -> $("{0:N1}MB" -f $size)"
}

Pop-Location

# Generate checksums
Write-Host "`nGenerating checksums.txt..."
Push-Location $ReleaseDir
$checksums = @()
Get-ChildItem -Filter "ASF-v*" | ForEach-Object {
    $hash = (Get-FileHash $_.Name -Algorithm SHA256).Hash.ToLower()
    $checksums += "$hash  $($_.Name)"
}
$checksums -join "`n" | Out-File -FilePath "checksums.txt" -Encoding ASCII
Write-Host "  checksums.txt updated"
Pop-Location

# Update VERSION
Set-Content -Path (Join-Path $ReleaseDir "VERSION") -Value $Version

# Summary
Write-Host "`n=== Build Complete ===" -ForegroundColor Green
Write-Host "Version: v$Version"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  git tag v$Version && git push origin v$Version"
Write-Host "  # GitHub Actions will build and publish automatically"
