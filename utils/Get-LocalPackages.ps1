function Get-LocalPackages {
	[cmdletBinding()]
	Param ()

	BEGIN {
		$rootDir = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
		$packagesFolder = Join-Path $rootDir "packages"
	}

	PROCESS {
		try {
			$localPackages = Get-ChildItem $packagesFolder -Recurse -Name | Where-Object {
				$_.Split([char]'\').Count -gt 1
			}
			$localPackages = $localPackages | ForEach-Object {$_ -replace "\\", "."}
		} catch {
			Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}

	END {
		return $localPackages
	}
}
