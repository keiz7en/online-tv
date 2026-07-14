# Setup script for Online TV
# Downloads VLC portable and prepares it for embedding

$vlcVersion = "3.0.21"
$vlcUrl = "https://get.videolan.org/vlc/$vlcVersion/win64/vlc-$vlcVersion-win64.zip"
$vlcDir = "vlc-files"
$tempDir = "$env:TEMP\vlc-download"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Online TV - VLC Setup Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if VLC files already exist
if (Test-Path "$vlcDir\vlc.exe") {
    Write-Host "VLC files already present in $vlcDir" -ForegroundColor Green
    Write-Host "Skipping download." -ForegroundColor Green
    exit 0
}

Write-Host "Downloading VLC portable..." -ForegroundColor Yellow
Write-Host "URL: $vlcUrl" -ForegroundColor Gray

# Create temp directory
if (-not (Test-Path $tempDir)) {
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
}

$zipFile = "$tempDir\vlc.zip"

try {
    # Download VLC
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $vlcUrl -OutFile $zipFile -UseBasicParsing
    
    Write-Host "Extracting VLC..." -ForegroundColor Yellow
    
    # Extract
    Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force
    
    # Find the extracted directory
    $extractedDir = Get-ChildItem -Path $tempDir -Directory | Where-Object { $_.Name -like "vlc-*" } | Select-Object -First 1
    
    if ($extractedDir) {
        # Copy files to vlc-files directory
        if (-not (Test-Path $vlcDir)) {
            New-Item -ItemType Directory -Path $vlcDir -Force | Out-Null
        }
        
        Copy-Item -Path "$($extractedDir.FullName)\*" -Destination $vlcDir -Recurse -Force
        
        Write-Host ""
        Write-Host "VLC downloaded and extracted successfully!" -ForegroundColor Green
        Write-Host "Location: $vlcDir" -ForegroundColor Green
        
        # Clean up
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    } else {
        Write-Host "Error: Could not find extracted VLC directory" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error downloading VLC: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please download VLC manually:" -ForegroundColor Yellow
    Write-Host "1. Go to https://get.videolan.org/vlc/" -ForegroundColor Yellow
    Write-Host "2. Download VLC for Windows (64-bit)" -ForegroundColor Yellow
    Write-Host "3. Extract the zip file" -ForegroundColor Yellow
    Write-Host "4. Copy all contents to the '$vlcDir' directory" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "Setup complete! You can now build the app:" -ForegroundColor Cyan
Write-Host "  wails build" -ForegroundColor White
