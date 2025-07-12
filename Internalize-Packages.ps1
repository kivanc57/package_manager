function Internalize-Packages{
	[cmdletBinding()]
	Param (
		[Parameter(Mandatory = $true)]
		[ValidateNotNullOrEmpty()]
		[string[]]$packages,

		[Parameter(Mandatory = $false)]
		[string]$resolveConflict
	)

	BEGIN {
		. "$PSScriptRoot\utils\Trace-Conflicts.ps1"
		$chopinConflictPath = "C:\ProgramData\chopin\conflict"
		$chopinConfigPath = "C:\ProgramData\chopin\appsettings.json"
	}

	PROCESS {
		try {
			$packages = $packages -split '\s+'
			foreach ($pckg in $packages) {
				$pckgName, $pckgVer = $pckg -split ':'
				Start-Process -FilePath "chopin" -ArgumentList  "internalize", "-p", $pckgName, "-v", $pckgVer, "--verbose", "Verbose" -Verb RunAs -Wait
			}

			Trace-Conflicts -conflictPath $chopinConflictPath

			if ($resolveConflict -eq 'true') {
				Start-Process -FilePath "chopin" -ArgumentList "resolve", "--all", "--verbose", "Verbose" -Verb RunAs -Wait
				Write-Host "RESOLVED: Conflicts"
			}
		#  FUTURE IMPLICATION --> chopin upload: automatically update the packages
		} catch {
		Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}

	END {
		Remove-Item $chopinConfigPath
		Write-Host "REMOVED: $chopinConfigPath"
	}
}

$packages = $env:PCKG_ID_VERS -split '\s+'
Internalize-Packages -packages $packages -resolveConflict $env:RESOLVE_CONFLICT
