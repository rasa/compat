param (
    [ValidatePattern("^[A-Za-z]$")]
    [char]$DriveLetter = 'Z'
)
$ProgressPreference = 'SilentlyContinue'
# Locate the temp VHDX file based on naming pattern
$pattern = "$env:TEMP\~tmp_${DriveLetter}_*.vhdx"
$vhds = Get-ChildItem -Path $pattern -ErrorAction SilentlyContinue

if (-not $vhds) {
    Write-Error "No matching VHDX found for $pattern"
    exit 0
}

foreach ($vhd in $vhds) {
    try {
        Write-Host "Dismounting VHD: $($vhd.FullName)"
        Dismount-VHD -Path $vhd.FullName -ErrorAction Stop
    } catch {
        Write-Warning "Failed to dismount $($vhd.FullName): $_"
    }

    try {
        Write-Host "Deleting VHD file: $($vhd.FullName)"
        Remove-Item -Path $vhd.FullName -Force -ErrorAction Stop
    } catch {
        Write-Warning "Failed to delete $($vhd.FullName): $_"
    }
}
