# SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

param (
    [ValidatePattern("^[A-Za-z]$")]
    [string]$DriveLetter = ''
)

$ProgressPreference = 'SilentlyContinue'

if ($DriveLetter -eq "") {
    # No drive letter provided -> dismount all matching VHDX files
    $pattern = "$env:TEMP\~compat_*.vhdx"
    $vhds = Get-ChildItem -Path $pattern -ErrorAction SilentlyContinue
} else {
    # Specific drive letter -> find matching VHD file by name OR by mount
    $pattern = "$env:TEMP\~compat_${DriveLetter}_*.vhdx"
    $vhds = Get-ChildItem -Path $pattern -ErrorAction SilentlyContinue

    if (-not $vhds) {
        # If no file matches, try to find the VHD by mounted drive letter
        try {
            $disk = Get-DiskImage | Where-Object {
                $_.Attached -and ($_.DevicePath -match $DriveLetter + ":")
            }
            if ($disk) {
                $vhds = @($disk)
            }
        } catch {
            Write-Output "No matching VHDX found for $pattern"
        }
    }
}

if (-not $vhds) {
    Write-Output "No matching VHDX found for $pattern"
    exit 0
}

foreach ($vhd in $vhds) {
    try {
        Write-Output "Dismounting VHD: $($vhd.FullName)"
        Dismount-VHD -Path $vhd.FullName -ErrorAction Stop
    } catch {
        Write-Warning "Failed to dismount $($vhd.FullName): $_"
    }

    try {
        Write-Output "Deleting VHD file: $($vhd.FullName)"
        Remove-Item -Path $vhd.FullName -Force -ErrorAction Stop
    } catch {
        Write-Warning "Failed to delete $($vhd.FullName): $_"
    }
}
