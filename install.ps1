# ASF — Architecture Security Framework
# Windows Installer — https://github.com/moksh5936-2/asfassumption
#
# Usage:
#   powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"
#
# Environment:
#   $env:ASF_VERSION="2.0.1"    — pin a specific version

param(
    [switch]$Upgrade,
    [switch]$Repair,
    [switch]$Clean,
    [switch]$Purge,
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

# Validate --purge requires --clean
if ($Purge -and -not $Clean) {
    Write-Host "Error: -Purge must be used with -Clean" -ForegroundColor Red
    exit 1
}

# ─── Help ──────────────────────────────────────────────────
if ($Help) {
    Write-Host "ASF Windows Installer" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  powershell -c `"irm https://raw.githubusercontent.com/${Repo}/main/install.ps1 | iex`""
    Write-Host "  powershell -c `"... | iex`" -Upgrade"
    Write-Host "  powershell -c `"... | iex`" -Repair"
    Write-Host "  powershell -c `"... | iex`" -Clean"
    Write-Host "  powershell -c `"... | iex`" -Clean -Purge"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Upgrade         Upgrade existing installation (backs up config)"
    Write-Host "  -Repair          Fix PATH without re-downloading"
    Write-Host "  -Clean           Remove binary, keep config, reinstall"
    Write-Host "  -Purge           Only with -Clean: removes config/cache/data too"
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
        $apiUrl = "https://api.github.com/repos/${Repo}/releases?per_page=10"
        $headers = @{"Accept"="application/vnd.github.v3+json"} + $AuthHeader
        $releases = Invoke-RestMethod -Uri $apiUrl -Headers $headers -ErrorAction Stop
        $release = $releases | Where-Object { -not $_.draft } | Select-Object -First 1
        if ($release) {
            $LatestTag = $release.tag_name
            $Version = $LatestTag -replace "^v", ""
        } else {
            throw "no non-draft release found"
        }
        Write-Host "  ✓ Latest: ${LatestTag}" -ForegroundColor Green
    } catch {
        $LatestTag = "v5.0.1"
        $Version = "5.0.1"
        Write-Host "  ⚠  Could not detect version, defaulting to ${LatestTag}" -ForegroundColor Yellow
        if (-not $Token) { Write-Host "  ⚠  Set GITHUB_TOKEN env var for private repos" -ForegroundColor Yellow }
    }
} else {
    $LatestTag = "v${Version}"
}

# Normalize version values
$LatestTag = $LatestTag.Trim("`r`n").Trim()
$Version = $Version.Trim("`r`n").Trim()

$AssetVersion = $Version
$BinaryName = "ASF-${LatestTag}-${OsArch}"
$DownloadUrl = "https://github.com/${Repo}/releases/download/${LatestTag}/${BinaryName}.exe"
$ChecksumsUrl = "https://github.com/${Repo}/releases/download/${LatestTag}/checksums.txt"

# Validate URL before download
if ($DownloadUrl -match "\s") {
    Write-Host "  ✗ Download URL contains whitespace or newlines: ${DownloadUrl}" -ForegroundColor Red
    exit 1
}
if ($BinaryName -notmatch "^ASF-") {
    Write-Host "  ✗ Invalid binary name: ${BinaryName}" -ForegroundColor Red
    exit 1
}

# ─── Check existing ───────────────────────────────────────
$InstalledBin = "${InstallDir}\asf.exe"
$HasExisting = Test-Path $InstalledBin -PathType Leaf

# ─── Clean mode ───────────────────────────────────────────
if ($Clean) {
    Write-Host "  Cleaning old ASF installation..." -ForegroundColor Cyan
    if (Test-Path $InstalledBin) { Remove-Item $InstalledBin -Force -ErrorAction SilentlyContinue }
    if (Test-Path "${BinDir}\asf.exe") { Remove-Item "${BinDir}\asf.exe" -Force -ErrorAction SilentlyContinue }
    if ($Purge) {
        if (Test-Path $ConfigDir) { Remove-Item $ConfigDir -Recurse -Force -ErrorAction SilentlyContinue }
        Write-Host "  ✓ Config, cache, data removed" -ForegroundColor Green
    } else {
        Write-Host "  ✓ Old binaries removed (config kept)" -ForegroundColor Green
    }
}

# ─── Repair mode ───────────────────────────────────────────
if ($Repair) {
    if (-not $HasExisting) {
        Write-Host "  ✗ No ASF binary found. Run installer without -Repair." -ForegroundColor Red
        exit 1
    }
    Write-Host "  Repairing ASF installation..." -ForegroundColor Cyan

    # Fix PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*${BinDir}*") {
        $newPath = "${BinDir};${userPath}"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "  ✓ Added ${BinDir} to PATH" -ForegroundColor Green
    } else {
        Write-Host "  ✓ ${BinDir} already in PATH" -ForegroundColor Green
    }

    # Verify
    $verOut = & $InstalledBin --version 2>$null
    Write-Host "  ✓ ${verOut}" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Repair complete. If 'asf' is not recognized, open a new PowerShell window." -ForegroundColor Cyan
    exit 0
}

# ─── Existing install detection ───────────────────────────
if ($HasExisting -and -not $Upgrade -and -not $Clean) {
    $installedVer = "unknown"
    try { $installedVer = & $InstalledBin --version } catch {}
    Write-Host "  ✓ ASF v${Version} is already installed (${installedVer})" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Run: asf" -ForegroundColor Cyan
    Write-Host "  To upgrade: add -Upgrade flag" -ForegroundColor Yellow
    exit 0
}

# ─── Backup existing config on upgrade ────────────────────
if ($Upgrade -and $HasExisting) {
    $BackupDir = "${InstallDir}\backups"
    New-Item -ItemType Directory -Force -Path $BackupDir | Out-Null
    $ConfigFile = "${ConfigDir}\config.yaml"
    if (Test-Path $ConfigFile) {
        $stamp = Get-Date -Format "yyyyMMdd-HHmmss"
        Copy-Item $ConfigFile -Destination "${BackupDir}\config.yaml.bak.${stamp}" -Force
        Write-Host "  ✓ Config backed up" -ForegroundColor Green
    }
    $LicenseFile = "${ConfigDir}\license.key"
    if (Test-Path $LicenseFile) {
        $stamp = Get-Date -Format "yyyyMMdd-HHmmss"
        Copy-Item $LicenseFile -Destination "${BackupDir}\license.key.bak.${stamp}" -Force
        Write-Host "  ✓ License backed up" -ForegroundColor Green
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
} else {
    Write-Host "  ✓ ${BinDir} already in PATH" -ForegroundColor Green
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

# ─── Verify ──────────────────────────────────────────────
Write-Host ""
Write-Host "  Verifying installation..." -ForegroundColor Cyan
$allOk = $true
if (Test-Path $InstalledBin) {
    Write-Host "  ✓ Binary: ${InstalledBin}" -ForegroundColor Green
} else {
    Write-Host "  ⚠  Binary not found" -ForegroundColor Yellow
    $allOk = $false
}
if (Test-Path "${BinDir}\asf.exe") {
    Write-Host "  ✓ Command: ${BinDir}\asf.exe" -ForegroundColor Green
} else {
    Write-Host "  ⚠  Command not in bin dir" -ForegroundColor Yellow
    $allOk = $false
}
$verOut = & $InstalledBin --version 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ ${verOut}" -ForegroundColor Green
} else {
    Write-Host "  ⚠  Version check failed" -ForegroundColor Yellow
    $allOk = $false
}
if ($allOk) {
    Write-Host "  ✓ All checks passed." -ForegroundColor Green
} else {
    Write-Host "  ⚠  Some checks failed" -ForegroundColor Yellow
}

# ─── Success ──────────────────────────────────────────────
$BinSize = "{0:N1}MB" -f ((Get-Item $InstalledBin).Length / 1MB)
Write-Host ""
Write-Host "  ✓ ASF v${Version} installed  (${BinSize})" -ForegroundColor Green
Write-Host ""
Write-Host "  Run: asf" -ForegroundColor Cyan
Write-Host ""
Write-Host "  If 'asf' is not recognized, open a new PowerShell window." -ForegroundColor Cyan
Write-Host ""
Write-Host "  Prerequisites (full functionality):" -ForegroundColor Cyan
Write-Host "    Tesseract (OCR): choco install tesseract" -ForegroundColor Cyan
Write-Host "    Ollama (AI):      https://ollama.com/download/windows" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Documentation: https://github.com/${Repo}" -ForegroundColor Cyan
