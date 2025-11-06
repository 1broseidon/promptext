#Requires -Version 5.1
[CmdletBinding()]
param (
    [string]$Version = "latest",
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Configuration
$owner = "1broseidon"
$repo = "promptext"
$binaryName = "promptext.exe"
$defaultInstallDir = Join-Path $env:LOCALAPPDATA "promptext"

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
    # GoReleaser uses lowercase "windows" and formats like: promptext_0.5.3_windows_amd64.zip
    # But for "latest" download URLs, version is omitted
    $patterns = @(
        "promptext_windows_$($osInfo.Arch).zip",
        "promptext_Windows_$($osInfo.Arch).zip",
        "promptext-windows-$($osInfo.Arch).zip"
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
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -like "*$defaultInstallDir*") {
        $newPath = ($currentPath.Split(';') | Where-Object { $_ -ne $defaultInstallDir }) -join ';'
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = ($env:Path.Split(';') | Where-Object { $_ -ne $defaultInstallDir }) -join ';'
    }
    
    # Remove installation directory
    if (Test-Path $defaultInstallDir) {
        Remove-Item -Path $defaultInstallDir -Recurse -Force
    }
    
    # Remove alias from profile
    if (Test-Path $PROFILE.CurrentUserCurrentHost) {
        $content = Get-Content $PROFILE.CurrentUserCurrentHost | Where-Object { $_ -notmatch "Set-Alias.*prx.*promptext" }
        Set-Content -Path $PROFILE.CurrentUserCurrentHost -Value $content
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
        
        # Find the executable entry (GoReleaser builds it as "prx.exe")
        $exeEntry = $zip.Entries | Where-Object {
            $_.Name -eq "prx.exe" -and
            -not $_.FullName.Contains("../") -and
            -not $_.FullName.StartsWith("/")
        } | Select-Object -First 1

        if (-not $exeEntry) {
            throw "Could not find valid prx.exe in the archive"
        }

        # Extract the executable and rename to promptext.exe
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

    # Update PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not ($currentPath -like "*$defaultInstallDir*")) {
        $newPath = "$currentPath;$defaultInstallDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = "$env:Path;$defaultInstallDir"
        Write-Status "Added $defaultInstallDir to PATH"
    }

    # Add alias to PowerShell profile
    $profileDir = Split-Path $PROFILE.CurrentUserCurrentHost -Parent
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $PROFILE.CurrentUserCurrentHost)) {
        New-Item -ItemType File -Path $PROFILE.CurrentUserCurrentHost -Force | Out-Null
    }

    $aliasLine = "Set-Alias prx '$defaultInstallDir\promptext.exe'"
    if (-not (Select-String -Path $PROFILE.CurrentUserCurrentHost -Pattern "Set-Alias.*prx.*promptext" -Quiet)) {
        Add-Content -Path $PROFILE.CurrentUserCurrentHost -Value $aliasLine
        Write-Status "Added 'prx' alias to PowerShell profile"
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
