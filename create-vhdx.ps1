param (
    [ValidatePattern("^[A-Za-z]$")]
    [string]$DriveLetter = 'Z',

    [ValidatePattern("^\d+(MB|GB|TB)$")]
    [string]$Size = "2GB",

    [ValidateSet("exFAT", "FAT", "FAT32", "NTFS", "ReFS")]
    [string]$FileSystem = 'NTFS'
)
$ProgressPreference = 'SilentlyContinue'
$DriveLetter = $DriveLetter.ToUpper()
$Size = $Size.ToUpper()
$FileSystem = $FileSystem.ToUpper()

function Convert-ToBytes {
    param([Parameter(Mandatory)][string]$Value)
    if ($Value -as [UInt64]) {
        return [UInt64]$Value   # already numeric, e.g. "2147483648"
    }
    switch -regex ($Value.ToUpperInvariant()) {
        '^(\d+)MB$' { return [UInt64]$Matches[1] * 1MB }
        '^(\d+)GB$' { return [UInt64]$Matches[1] * 1GB }
        '^(\d+)TB$' { return [UInt64]$Matches[1] * 1TB }
        default     { throw "Unsupported size format: $Value" }
    }
}

$sizeBytes = Convert-ToBytes $Size

$randHex = '{0:x8}' -f (Get-Random -Minimum 0 -Maximum 4294967295)
$vhdfile = "$env:TEMP\~tmp_${DriveLetter}_${FileSystem}_${Size}_${randHex}.vhdx"

New-VHD -Path $vhdfile -Dynamic -SizeBytes $sizeBytes | Out-Null
Mount-VHD -Path $vhdfile | Out-Null
Start-Sleep -Seconds 1

$disk = Get-Disk | Where-Object { $_.PartitionStyle -eq 'RAW' }
if (-not $disk) {
    Write-Error "Unable to locate newly mounted RAW disk for $vhdfile"
    exit 1
}

$partitionStyle = if ($FileSystem -in @("exFAT", "NTFS", "ReFS")) { "GPT" } else { "MBR" }

Initialize-Disk -Number $disk.Number -PartitionStyle GPT -ErrorAction Stop
$part = New-Partition -DiskNumber $disk.Number -DriveLetter $DriveLetter -UseMaximumSize -ErrorAction Stop

$label = "${DriveLetter}${FileSystem}${Size}"
# Attempt to format
$null = Format-Volume -Partition $part -FileSystem $FileSystem `
    -NewFileSystemLabel $label -Confirm:$false -ErrorAction Stop

# Verify FS type
$vol = Get-Volume -DriveLetter $part.DriveLetter
if ($vol.FileSystem -ne $FileSystem) {
    Write-Error "Requested $FileSystem but got $($vol.FileSystem) on $($part.DriveLetter):"
    Dismount-VHD -Path $vhdfile -ErrorAction SilentlyContinue
    Remove-Item -Path $vhdfile -Force -ErrorAction SilentlyContinue
    exit 1  # ensures %ERRORLEVEL% / $LASTEXITCODE is set
}

Write-Output "Created a $Size $FileSystem drive $($part.DriveLetter): ($label) at $vhdfile"
