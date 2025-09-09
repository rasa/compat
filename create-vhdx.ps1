# SPDX-FileCopyrightText: Copyright (c) 2025 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

param (
    [ValidatePattern("^[A-Za-z]?$")]
    [string]$DriveLetter,   # optional, auto-picked if not supplied

    [ValidateSet("exFAT", "FAT", "FAT32", "NTFS", "ReFS")]
    [string]$FileSystem = 'NTFS',

    [ValidatePattern("^\d+(MB|GB|TB)$")]
    [string]$Size = "2GB"
)

$ProgressPreference = 'SilentlyContinue'
$Size       = $Size.ToUpper()
$FileSystem = $FileSystem.ToUpper()

Add-Type -AssemblyName System.Numerics

function Convert-ToBase36 {
    param([System.Numerics.BigInteger]$Value)
    $chars = "0123456789abcdefghijklmnopqrstuvwxyz"
    if ($Value -eq 0) { return "0" }
    $result = ""
    while ($Value -gt 0) {
        $result = $chars[$Value % 36] + $result
        $Value  = [System.Numerics.BigInteger]::Divide($Value, 36)
    }
    return $result
}

function Convert-ToBytes {
    param([Parameter(Mandatory)][string]$Value)
    if ($Value -as [UInt64]) { return [UInt64]$Value }
    switch -regex ($Value.ToUpperInvariant()) {
        '^(\d+)MB$' { return [UInt64]$Matches[1] * 1MB }
        '^(\d+)GB$' { return [UInt64]$Matches[1] * 1GB }
        '^(\d+)TB$' { return [UInt64]$Matches[1] * 1TB }
        default     { throw "Unsupported size format: $Value" }
    }
}

# Gather all currently used drive letters (local + network)
$usedLetters = [System.IO.DriveInfo]::GetDrives().Name.TrimEnd('\') |
    ForEach-Object { $_[0].ToString().ToUpper() }

# If no drive letter supplied, auto-pick from Z down to A
if (-not $DriveLetter) {
    foreach ($letter in [char[]]([byte][char]'Z'..[byte][char]'A')) {
        if ($usedLetters -notcontains $letter) {
            $DriveLetter = $letter
            Write-Output "Auto-selected drive letter $DriveLetter"
            break
        }
    }
    if (-not $DriveLetter) {
        Write-Error "No available drive letters found (Z..D all in use)"
        exit 1
    }
} else {
    $DriveLetter = $DriveLetter.ToUpper()
    if ($usedLetters -contains $DriveLetter) {
        Write-Error "Drive letter $DriveLetter is already in use"
        exit 1
    }
}

# Convert size string to bytes
$sizeBytes = Convert-ToBytes $Size

# Generate random base36 suffix
$rand36 = (Convert-ToBase36 (Get-Random -Minimum 0 -Maximum (
    [System.Numerics.BigInteger]::Pow(36,8)))).PadLeft(8,'0')

# Temp VHD path
$vhdfile = "$env:TEMP\~compat_${DriveLetter}_${FileSystem}_${Size}_${rand36}.vhdx"

# Create and mount
New-VHD -Path $vhdfile -Dynamic -SizeBytes $sizeBytes | Out-Null
Mount-VHD -Path $vhdfile | Out-Null
Start-Sleep -Seconds 1

# Find the newly mounted RAW disk
$disk = Get-Disk | Where-Object { $_.PartitionStyle -eq 'RAW' }
if (-not $disk) {
    Write-Error "Unable to locate newly mounted RAW disk for $vhdfile"
    exit 1
}

# Decide partition style
$partitionStyle = if ($FileSystem -in @("exFAT","NTFS","ReFS")) { "GPT" } else { "MBR" }

Write-Output "Initializing disk $($disk.Number) as $partitionStyle with drive letter $DriveLetter"

Initialize-Disk -Number $disk.Number -PartitionStyle $partitionStyle -ErrorAction Stop
$part = New-Partition -DiskNumber $disk.Number -DriveLetter $DriveLetter -UseMaximumSize -ErrorAction Stop

# Label & format
$label = "${DriveLetter}${FileSystem}${Size}"
Format-Volume -Partition $part -FileSystem $FileSystem `
    -NewFileSystemLabel $label -Confirm:$false -ErrorAction Stop | Out-Null

# Verify FS type
$vol = Get-Volume -DriveLetter $part.DriveLetter
if ($vol.FileSystem -ne $FileSystem) {
    Write-Error "Requested $FileSystem but got $($vol.FileSystem) on $($part.DriveLetter):"
    Dismount-VHD -Path $vhdfile -ErrorAction SilentlyContinue
    Remove-Item -Path $vhdfile -Force -ErrorAction SilentlyContinue
    exit 1
}

Write-Output "Created a $Size $FileSystem drive $($part.DriveLetter): ($label) at $vhdfile"
