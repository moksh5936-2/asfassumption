# ASF — Architecture Security Framework
# Windows Installer — https://github.com/moksh5936-2/asfassumption
#
# Usage:
#   powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"
#
# Environment:
#   $env:ASF_VERSION="2.0.0"    — pin a specific version

param(
    [switch]$Upgrade,
    [switch]$Help,
    [string]$Token = ""
)

# Fall back to env var if no param
if (-not $Token) { $Token = [Environment]::GetEnvironmentVariable("GITHUB_TOKEN", "User") }
if (-not $Token) { $Token = [Environment]::GetEnvironmentVariable("GITHUB_TOKEN", "Process") }

# Auth header for private repos
$AuthHeader = @{}
if ($Token) { $AuthHeader["Authorization"] = "token $Token" }

$Repo = "moksh5936-2/asfassumption"
$Version = [Environment]::GetEnvironmentVariable("ASF_VERSION", "User")
$InstallDir = "$env:LOCALAPPDATA\ASF"
$BinDir = "$env:LOCALAPPDATA\ASF\bin"
$ConfigDir = "$env:APPDATA\ASF"

# ─── Help ──────────────────────────────────────────────────
if ($Help) {
    Write-Host "ASF Windows Installer" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  powershell -c `"irm https://raw.githubusercontent.com/${Repo}/main/install.ps1 | iex`""
    Write-Host "  powershell -c `"... | iex`" -Upgrade"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Upgrade         Upgrade an existing installation"
    Write-Host "  -Help            Show this help"
    Write-Host ""
    Write-Host "Environment:"
    Write-Host "  ASF_VERSION=x.y.z   Pin a specific version (default: latest)"
    exit 0
}

# ─── Header ────────────────────────────────────────────────
Write-Host "  /\   /\" -ForegroundColor Cyan
Write-Host " (  o.o  )" -ForegroundColor Cyan
Write-Host "  >  ^  < " -ForegroundColor Cyan
Write-Host " ASF Security Framework" -ForegroundColor Cyan
Write-Host ""

# ─── Detect OS/Arch ───────────────────────────────────────
$OsArch = "windows-amd64"

# ─── Detect version ───────────────────────────────────────
if (-not $Version) {
    Write-Host "  Detecting latest version..." -ForegroundColor Cyan
    try {
        $apiUrl = "https://api.github.com/repos/${Repo}/releases/latest"
        $headers = @{"Accept"="application/vnd.github.v3+json"} + $AuthHeader
        $response = Invoke-RestMethod -Uri $apiUrl -Headers $headers -ErrorAction Stop
        $Version = $response.tag_name -replace "^v", ""
        Write-Host "  ✓ Latest: v${Version}" -ForegroundColor Green
    } catch {
        $Version = "2.0.0"
        Write-Host "  ⚠  Could not detect version, defaulting to v${Version}" -ForegroundColor Yellow
        if (-not $Token) { Write-Host "  ⚠  Set GITHUB_TOKEN env var for private repos" -ForegroundColor Yellow }
    }
}

$BinaryName = "ASF-v${Version}-${OsArch}"
$DownloadUrl = "https://github.com/${Repo}/releases/download/v${Version}/${BinaryName}.exe"
$ChecksumsUrl = "https://github.com/${Repo}/releases/download/v${Version}/checksums.txt"

# ─── Check existing ───────────────────────────────────────
$InstalledBin = "${InstallDir}\asf.exe"
if (Test-Path $InstalledBin -PathType Leaf) {
    if (-not $Upgrade) {
        $installedVer = "unknown"
        try { $installedVer = & $InstalledBin --version } catch {}
        Write-Host "  ✓ ASF is already installed (${installedVer})" -ForegroundColor Green
        Write-Host ""
        Write-Host "  Run: asf" -ForegroundColor Cyan
        Write-Host "  To upgrade: add -Upgrade flag" -ForegroundColor Yellow
        exit 0
    }
}

# ─── Download ──────────────────────────────────────────────
$TmpDir = "$env:TEMP\asf-install"
if (Test-Path $TmpDir) { Remove-Item -Recurse -Force $TmpDir }
New-Item -ItemType Directory -Force -Path $TmpDir | Out-Null

Write-Host ""
Write-Host "  Downloading ASF v${Version} for Windows..." -ForegroundColor Cyan
Write-Host "  ${DownloadUrl}" -ForegroundColor Cyan
Write-Host ""

try {
    # For private repos, download via API if token is available
    if ($Token) {
        $releaseUrl = "https://api.github.com/repos/${Repo}/releases/tags/v${Version}"
        $release = Invoke-RestMethod -Uri $releaseUrl -Headers $AuthHeader -ErrorAction Stop
        $asset = $release.assets | Where-Object { $_.name -eq $BinaryName }
        if ($asset) {
            $dlHeaders = $AuthHeader.Clone()
            $dlHeaders["Accept"] = "application/octet-stream"
            Invoke-WebRequest -Uri $asset.url -Headers $dlHeaders -OutFile "${TmpDir}\asf.exe" -ErrorAction Stop
        } else {
            throw "Binary ${BinaryName} not found in release assets"
        }
    } else {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile "${TmpDir}\asf.exe" -ErrorAction Stop
    }
} catch {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
    Write-Host "  ✗ Download failed: $_" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path "${TmpDir}\asf.exe") -or ((Get-Item "${TmpDir}\asf.exe").Length -eq 0)) {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
    Write-Host "  ✗ Downloaded file is empty or missing" -ForegroundColor Red
    exit 1
}

# ─── Checksum ──────────────────────────────────────────────
Write-Host "  Verifying checksum..." -ForegroundColor Cyan
try {
    if ($Token) {
        $releaseUrl = "https://api.github.com/repos/${Repo}/releases/tags/v${Version}"
        $release = Invoke-RestMethod -Uri $releaseUrl -Headers $AuthHeader -ErrorAction Stop
        $csAsset = $release.assets | Where-Object { $_.name -eq "checksums.txt" }
        if ($csAsset) {
            $csHeaders = $AuthHeader.Clone()
            $csHeaders["Accept"] = "application/octet-stream"
            $checksums = (Invoke-WebRequest -Uri $csAsset.url -Headers $csHeaders -ErrorAction Stop).Content
        } else {
            $checksums = ""
        }
    } else {
        $checksums = (Invoke-WebRequest -Uri $ChecksumsUrl -ErrorAction Stop).Content
    }
    $expectedHash = ($checksums -split "`n" | Where-Object { $_ -match $BinaryName } | ForEach-Object { ($_ -split "\s+")[0] })
    if ($expectedHash) {
        $computedHash = (Get-FileHash "${TmpDir}\asf.exe" -Algorithm SHA256).Hash.ToLower()
        if ($computedHash -ne $expectedHash) {
            Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
            Write-Host "  ✗ Checksum mismatch! Expected ${expectedHash}, got ${computedHash}" -ForegroundColor Red
            exit 1
        }
        Write-Host "  ✓ Checksum verified" -ForegroundColor Green
    } else {
        Write-Host "  ⚠  No checksum found for ${BinaryName}" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ⚠  Could not verify checksum" -ForegroundColor Yellow
}

# ─── Verify binary ─────────────────────────────────────────
try {
    $binVer = & "${TmpDir}\asf.exe" --version
    Write-Host "  ✓ Binary verified: ${binVer}" -ForegroundColor Green
} catch {
    Write-Host "  ⚠  Could not verify binary version" -ForegroundColor Yellow
}

# ─── Install ───────────────────────────────────────────────
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

Copy-Item "${TmpDir}\asf.exe" -Destination "${InstalledBin}" -Force

# Add to PATH if needed
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*${BinDir}*") {
    $newPath = "${BinDir};${userPath}"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "${BinDir};${env:Path}"
    Write-Host "  ✓ Added ${BinDir} to PATH" -ForegroundColor Green
}

# ─── Config ────────────────────────────────────────────────
New-Item -ItemType Directory -Force -Path $ConfigDir | Out-Null
$ConfigPath = "${ConfigDir}\config.yaml"
if (-not (Test-Path $ConfigPath)) {
    @"
general:
  theme: Dark
  fox_style: Classic
analysis:
  depth: deep
  stride: true
  controls: true
ai:
  enabled: false
  active_model: ""
  installed_models: []
output:
  default: markdown
  directory: ./reports
appearance:
  theme: Dark
  fox_style: Classic
engine:
  use_native_engine: true
"@ | Out-File -FilePath $ConfigPath -Encoding ascii
    Write-Host "  ✓ Created default config" -ForegroundColor Green
}

# ─── Cleanup ──────────────────────────────────────────────
Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue

# ─── Create alias convenience ─────────────────────────────
$AsfExe = "${InstalledBin}\asf.exe"

# ─── Success ──────────────────────────────────────────────
$BinSize = "{0:N1}MB" -f ((Get-Item $AsfExe).Length / 1MB)
Write-Host ""
Write-Host "  ✓ ASF v${Version} installed  (${BinSize})" -ForegroundColor Green
Write-Host ""
Write-Host "  Run: asf" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Prerequisites (full functionality):" -ForegroundColor Cyan
Write-Host "    Tesseract (OCR): choco install tesseract" -ForegroundColor Cyan
Write-Host "    Ollama (AI):      https://ollama.com/download/windows" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Documentation: https://github.com/${Repo}" -ForegroundColor Cyan
