function Get-MissingPackages {
	[cmdletBinding()]
	Param(
			[Parameter(Mandatory = $true)]
			[ValidateNotNullOrEmpty()]
			[string[]]$localPackages,

			[Parameter(Mandatory = $true)]
			[ValidateNotNullOrEmpty()]
			[string[]]$remotePackages
			)

	BEGIN {
		$missingPackages = @()
		$rootDir = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
		$excludedPackages = Get-Content (Join-Path $rootDir "packages_excluded.txt")
	}

	PROCESS {
		try {
			foreach ($pckg in $remotePackages) {
				if ($localPackages -contains $pckg) {
					Write-Host "FOUND: $pckg"
					continue
				}
				$isExcluded = $false
				foreach ($pckgExcluded in $excludedPackages) {
					if ($pckg -match $pckgExcluded) {
						Write-Host "EXCLUDED: $pckg"
						$isExcluded = $true
						break
					}
				}
				if (-not $isExcluded) {
					Write-Host "MISSING: $pckg"
					$missingPackages += $pckg
				}
			}
		} catch {
			Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}

	END {
		return $missingPackages
	}
}
