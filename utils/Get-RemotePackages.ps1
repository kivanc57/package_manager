function Get-RemotePackages {
	[cmdletBinding()]
	Param (
		[Parameter(Mandatory = $false)]
		[string[]]$packages
	)

	BEGIN {
		$remotePackages = @()
	}

	PROCESS {
		try {
			foreach ($repo in $artifactoryList) {
				$req = Invoke-WebRequest -Uri $repo -UseBasicParsing
				$fetchedPackages = $req.links | ForEach-Object { $_.href -replace "\.nupkg$", "" }
				$remotePackages =  ($remotePackages + $fetchedPackages) | Select-Object -Unique
			}
		} catch {
			Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}

	END {
		return $remotePackages
	}
}
