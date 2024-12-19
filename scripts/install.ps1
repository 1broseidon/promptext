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

function Get-OSInfo {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { "x86_64" }
        "ARM64" { "arm64" }
        default {
            throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
        }
    }
    return @{
        OS = "Windows"
        Arch = $arch
    }
}

function Get-AssetUrl {
    param($Release)
    $osInfo = Get-OSInfo
    $patterns = @(
        "promptext_$($osInfo.OS)_$($osInfo.Arch).zip",
        "promptext-$($osInfo.OS)-$($osInfo.Arch).zip"
    )
    
    Write-Host "Looking for release assets:" -ForegroundColor Yellow
    foreach ($pattern in $patterns) {
        Write-Host "- $pattern"
    }
    Write-Host "`nAvailable assets:" -ForegroundColor Yellow
    $Release.assets | ForEach-Object { Write-Host "- $($_.name)" }
    
    foreach ($pattern in $patterns) {
        $assets = $Release.assets | Where-Object { $_.name -eq $pattern }
        if ($assets) {
            Write-Host "`nFound matching asset: $($assets[0].name)" -ForegroundColor Green
            return $assets[0].browser_download_url
        }
    }
    
    throw "Could not find compatible Windows binary. Please report this issue."
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
        $FilePath,
        $AssetName
    )
    
    $checksumAsset = $Release.assets | Where-Object { 
        $_.name -eq "checksums.txt" -or 
        $_.name -eq "SHA256SUMS" -or
        $_.name -eq "sha256sums.txt"
    }
    if (-not $checksumAsset) {
        Write-Warning "Skipping checksum verification: checksum file not found in release"
        return
    }
    
    Write-Status "Verifying checksum..."
    Write-Host "Asset name: $AssetName"
    $checksumUrl = $checksumAsset.browser_download_url
    $checksumContent = (Invoke-WebRequest -Uri $checksumUrl -UseBasicParsing).Content
    Write-Host "Checksum content:"
    Write-Host $checksumContent
    $expectedChecksum = ($checksumContent -split "`n" | Where-Object { $_ -like "*$AssetName*" }) -split '\s+' | Select-Object -First 1
    
    if (-not $expectedChecksum) {
        Write-Warning "Skipping checksum verification: checksum not found for $AssetName"
        return
    }
    
    $actualChecksum = Get-FileHash -Path $FilePath -Algorithm SHA256 | Select-Object -ExpandProperty Hash
    if ($actualChecksum -ne $expectedChecksum) {
        throw "Checksum verification failed.`nExpected: $expectedChecksum`nGot: $actualChecksum"
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

    $osInfo = Get-OSInfo
    Write-Status "Installing promptext..."
    Write-Status "OS: $($osInfo.OS)"
    Write-Status "Architecture: $($osInfo.Arch)"
    Write-Status "Installation directory: $defaultInstallDir"

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

    # Download and verify
    Write-Status "Downloading binary..."
    $zipPath = Join-Path $env:TEMP "promptext.zip"
    $assetName = $downloadUrl.Split('/')[-1]
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    Verify-Checksum -Release $release -FilePath $zipPath -AssetName $assetName

    Write-Status "Extracting files..."
    try {
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        $zip = [System.IO.Compression.ZipFile]::OpenRead($zipPath)
        
        # Find the executable entry
        $exeEntry = $zip.Entries | Where-Object { 
            $_.Name -eq "promptext.exe" -and 
            -not $_.FullName.Contains("../") -and
            -not $_.FullName.StartsWith("/")
        } | Select-Object -First 1
        
        if (-not $exeEntry) {
            throw "Could not find valid promptext.exe in the archive"
        }
        
        # Extract the executable
        $exePath = Join-Path $defaultInstallDir "promptext.exe"
        [System.IO.Compression.ZipFileExtensions]::ExtractToFile($exeEntry, $exePath, $true)
        
    } catch {
        throw "Failed to extract executable: $_"
    } finally {
        if ($zip) {
            $zip.Dispose()
        }
        Remove-Item $zipPath -ErrorAction SilentlyContinue
    }

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

    # Verify installation
    $promptextPath = Join-Path $defaultInstallDir "promptext.exe"
    if (-not (Test-Path $promptextPath)) {
        throw "Installation failed: promptext.exe not found at $promptextPath"
    }

    # Test the installation
    try {
        $version = & $promptextPath -v
        Write-Status "Installation verified: $version"
    } catch {
        Write-Warning "Installation completed but verification failed: $_"
    }

    Write-Host "`n✨ Installation complete!" -ForegroundColor Green
    Write-Host "You can use either 'promptext' or 'prx' command after restarting your terminal." -ForegroundColor Yellow
    Write-Host "To uninstall, run this script with -Uninstall flag" -ForegroundColor Yellow

} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}
