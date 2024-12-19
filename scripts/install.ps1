#Requires -Version 5.1
[CmdletBinding()]
param (
    [switch]$UserInstall,
    [string]$Version = "latest",
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Configuration
$owner = "1broseidon"
$repo = "promptext"
$binaryName = "promptext.exe"
$defaultInstallDir = if ($UserInstall) {
    Join-Path $env:LOCALAPPDATA "promptext"
} else {
    Join-Path ${env:ProgramFiles} "promptext"
}

function Write-Status {
    param([string]$Message)
    Write-Host "→ $Message" -ForegroundColor Blue
}

function Get-LatestRelease {
    $url = "https://api.github.com/repos/$owner/$repo/releases/latest"
    $release = Invoke-RestMethod -Uri $url -UseBasicParsing
    return $release
}

function Get-AssetUrl {
    param($Release)
    $assets = $Release.assets | Where-Object { 
        $_.name -like "*windows*64*.zip" -or 
        $_.name -like "*win*64*.zip" -or
        $_.name -like "*Windows*64*.zip" -or
        $_.name -like "*Win*64*.zip"
    }
    if (-not $assets) {
        Write-Host "Available assets:" -ForegroundColor Yellow
        $Release.assets | ForEach-Object { Write-Host "- $($_.name)" }
        throw "No compatible Windows binary found in release. Please report this issue."
    }
    return $assets[0].browser_download_url
}

function Uninstall-Promptext {
    Write-Status "Uninstalling promptext..."
    
    # Remove from PATH
    $envTarget = if ($UserInstall) { "User" } else { "Machine" }
    $currentPath = [Environment]::GetEnvironmentVariable("Path", $envTarget)
    $binPath = $defaultInstallDir
    
    if ($currentPath -like "*$binPath*") {
        $newPath = ($currentPath.Split(';') | Where-Object { $_ -ne $binPath }) -join ';'
        [Environment]::SetEnvironmentVariable("Path", $newPath, $envTarget)
        $env:Path = ($env:Path.Split(';') | Where-Object { $_ -ne $binPath }) -join ';'
    }
    
    # Remove installation directory
    if (Test-Path $defaultInstallDir) {
        Remove-Item -Path $defaultInstallDir -Recurse -Force
    }
    
    # Remove alias from profile
    $aliasPath = if ($UserInstall) {
        Join-Path $env:USERPROFILE "Documents\WindowsPowerShell"
    } else {
        "$env:SystemRoot\System32\WindowsPowerShell\v1.0"
    }
    $profilePath = Join-Path $aliasPath "Microsoft.PowerShell_profile.ps1"
    
    if (Test-Path $profilePath) {
        $content = Get-Content $profilePath | Where-Object { $_ -notmatch "Set-Alias.*prx.*promptext" }
        Set-Content -Path $profilePath -Value $content
    }
    
    Write-Host "Promptext has been uninstalled successfully." -ForegroundColor Green
    exit 0
}

function Verify-Checksum {
    param(
        $Release,
        $FilePath
    )
    
    $checksumAsset = $Release.assets | Where-Object { $_.name -like "*checksums.txt" }
    if (-not $checksumAsset) {
        throw "Checksum file not found in release"
    }
    
    $checksumUrl = $checksumAsset.browser_download_url
    $checksumContent = (Invoke-WebRequest -Uri $checksumUrl -UseBasicParsing).Content
    $expectedChecksum = ($checksumContent -split "`n" | Where-Object { $_ -like "*windows*amd64*.zip" }) -split '\s+' | Select-Object -First 1
    
    if (-not $expectedChecksum) {
        throw "Checksum not found for Windows binary"
    }
    
    $actualChecksum = Get-FileHash -Path $FilePath -Algorithm SHA256 | Select-Object -ExpandProperty Hash
    if ($actualChecksum -ne $expectedChecksum) {
        throw "Checksum verification failed. Expected: $expectedChecksum, Got: $actualChecksum"
    }
    
    Write-Status "Checksum verification successful"
}

try {
    # Handle execution policy
    $policy = Get-ExecutionPolicy
    if ($policy -eq "Restricted") {
        Write-Status "Setting execution policy to RemoteSigned for current process..."
        Set-ExecutionPolicy -Scope Process -ExecutionPolicy RemoteSigned -Force
    }

    if ($Uninstall) {
        Uninstall-Promptext
        exit 0
    }

    Write-Status "Installing promptext..."

    # Check for admin rights if not user install
    if (-not $UserInstall) {
        $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
        if (-not $isAdmin) {
            throw "Administrator rights required. Run as administrator or use -UserInstall for current user installation."
        }
    }

    # Create install directory
    if (-not (Test-Path $defaultInstallDir)) {
        New-Item -ItemType Directory -Path $defaultInstallDir | Out-Null
    }

    # Get latest release info
    Write-Status "Fetching release information..."
    $release = Get-LatestRelease
    $downloadUrl = Get-AssetUrl $release

    # Download and extract
    Write-Status "Downloading binary..."
    $zipPath = Join-Path $env:TEMP "promptext.zip"
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing

    Write-Status "Extracting files..."
    Expand-Archive -Path $zipPath -DestinationPath $defaultInstallDir -Force
    Remove-Item $zipPath

    # Add to PATH
    $binPath = $defaultInstallDir
    $envTarget = if ($UserInstall) { "User" } else { "Machine" }
    $currentPath = [Environment]::GetEnvironmentVariable("Path", $envTarget)
    
    if ($currentPath -notlike "*$binPath*") {
        Write-Status "Adding to PATH..."
        $newPath = "$currentPath;$binPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, $envTarget)
        $env:Path = "$env:Path;$binPath"
    }

    # Create alias
    Write-Status "Creating alias..."
    $aliasPath = if ($UserInstall) {
        Join-Path $env:USERPROFILE "Documents\WindowsPowerShell"
    } else {
        "$env:SystemRoot\System32\WindowsPowerShell\v1.0"
    }
    
    if (-not (Test-Path $aliasPath)) {
        New-Item -ItemType Directory -Path $aliasPath -Force | Out-Null
    }

    $profilePath = Join-Path $aliasPath "Microsoft.PowerShell_profile.ps1"
    $aliasContent = "Set-Alias -Name prx -Value promptext.exe"
    
    if (Test-Path $profilePath) {
        if (-not (Get-Content $profilePath | Select-String "Set-Alias.*prx.*promptext")) {
            Add-Content -Path $profilePath -Value "`n$aliasContent"
        }
    } else {
        Set-Content -Path $profilePath -Value $aliasContent
    }

    Write-Status "Installation complete! Run 'promptext -v' to verify."
    Write-Host "Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
    Write-Host "To uninstall, run this script with -Uninstall flag" -ForegroundColor Yellow

} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}
