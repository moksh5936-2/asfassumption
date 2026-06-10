# ASF Release Build Script (PowerShell)
# Usage: .\scripts\build-release.ps1 [-Version "1.0.0"]
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
Get-ChildItem $ReleaseDir -Filter "asf-*" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.tar.gz" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.zip" | Remove-Item -Force
Get-ChildItem $ReleaseDir -Filter "*.exe" | Remove-Item -Force

# Build
Write-Host "`nBuilding for all platforms..."
Push-Location $BuildDir

foreach ($platform in $Platforms) {
    $outName = "asf-$($platform.OS)-$($platform.Arch)"
    if ($platform.OS -eq "windows") {
        $outName += ".exe"
    }
    $outPath = Join-Path $ReleaseDir $outName

    Write-Host "  Building $($platform.OS)/$($platform.Arch)..."
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    go build -ldflags="-X 'main.version=$Version'" -o $outPath .
}

Pop-Location

# Generate checksums
Write-Host "`nGenerating checksums..."
Push-Location $ReleaseDir
$checksums = @()
Get-ChildItem -Filter "asf-*" | ForEach-Object {
    $hash = (Get-FileHash $_.Name -Algorithm SHA256).Hash.ToLower()
    $checksums += "$hash  $($_.Name)"
}
if (Test-Path "install.sh") {
    $hash = (Get-FileHash "install.sh" -Algorithm SHA256).Hash.ToLower()
    $checksums += "$hash  install.sh"
}
if (Test-Path "VERSION") {
    $hash = (Get-FileHash "VERSION" -Algorithm SHA256).Hash.ToLower()
    $checksums += "$hash  VERSION"
}
$checksums -join "`n" | Out-File -FilePath "checksums.txt" -Encoding ASCII
Write-Host "  checksums.txt updated"
Pop-Location

# Create archives
Write-Host "`nCreating archives..."
foreach ($platform in $Platforms) {
    $binaryName = "asf-$($platform.OS)-$($platform.Arch)"
    $binaryPath = Join-Path $ReleaseDir $binaryName
    if ($platform.OS -eq "windows") {
        $binaryPath += ".exe"
    }

    if (Test-Path $binaryPath) {
        if ($platform.OS -eq "windows") {
            $archive = Join-Path $ReleaseDir "$($platform.OS)-$($platform.Arch).zip"
            Compress-Archive -Path @($binaryPath, (Join-Path $ReleaseDir "install.sh"), (Join-Path $ReleaseDir "VERSION"), (Join-Path $ReleaseDir "checksums.txt")) -DestinationPath $archive -Force
            Write-Host "  $archive created"
        } else {
            $archive = Join-Path $ReleaseDir "$($platform.OS)-$($platform.Arch).tar.gz"
            # PowerShell doesn't have native tar.gz support; use tar if available
            $tarAvailable = Get-Command tar -ErrorAction SilentlyContinue
            if ($tarAvailable) {
                Push-Location $ReleaseDir
                tar czf $archive $binaryName install.sh VERSION checksums.txt
                Pop-Location
                Write-Host "  $archive created"
            } else {
                Write-Host "  WARNING: tar not available, skipping archive for $($platform.OS)/$($platform.Arch)" -ForegroundColor Yellow
            }
        }
    }
}

# Update VERSION
Set-Content -Path (Join-Path $ReleaseDir "VERSION") -Value $Version

# Summary
Write-Host "`n=== Build Complete ===" -ForegroundColor Green
Write-Host "Version: v$Version"
Write-Host "Release directory: $ReleaseDir"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Test each binary"
Write-Host "  2. Verify checksums"
Write-Host "  3. Create GitHub release"
Write-Host "  4. Upload archives"
