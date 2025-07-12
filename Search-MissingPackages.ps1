function Search-MissingPackages {
	[cmdletBinding()]
	Param (
		[Parameter(Mandatory = $false)]
		[string[]]$packages
	)

	BEGIN {
		. "$PSScriptRoot\utils\Get-LocalPackages.ps1"
		. "$PSScriptRoot\utils\Get-LocalPackages.ps1"
		. "$PSScriptRoot\utils\Get-RemotePackages.ps1"
		. "$PSScriptRoot\utils\Get-MissingPackages.ps1"
	}

	PROCESS {
		try {
			$localPackages = Get-LocalPackages
			$remotePackages = Get-RemotePackages -packages $packages
			$missingPackages = Get-MissingPackages -localPackages $localPackages -remotePackages $remotePackages
		} catch {
			Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}

	END {
		if ($missingPackages.Count -ne 0) {
			Write-Host "======MISSING PACKAGES======"
			Write-Host ($missingPackages -join "`n")
			Write-Host "============================"
			Write-Host "TOTAL MISSING PACKAGE AMOUNT => $($missingPackages.Count)"
		}
		Write-Host "TOTAL REMOTE PACKAGE AMOUNT  => $($remotePackages.Count)"
	}
}

$packages = $env:PCKG_ID_VERS -split '\s+'
Search-MissingPackages -packages $packages
