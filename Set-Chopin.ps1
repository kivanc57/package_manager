function Set-Chopin {
	[cmdletBinding()]
	Param (
		[Parameter(Mandatory = $true)]
		[ValidateNotNullOrEmpty()]
		[string]$configFileName
	)

	BEGIN {
		$configContent = Get-Content -Path $configFileName
		$chopinPath = "C:\ProgramData\chopin\"
		$chopinVer = "1.0.1"
	}

	PROCESS {
		try {
			Start-Process -FilePath "choco" -ArgumentList "install", "chopin", "--version", $chopinVer, "--source=$chopinURL", "-y", "--no-progress" -Verb RunAs -Wait
			$replacedConfigContent = $configContent -replace "ARTIFACTORY_USER_NAME:ARTIFACTORY_PASSWORD", "${env:ARTIFACTORY_USER_NAME}:${env:ARTIFACTORY_PASSWORD}"
			$replacedConfigContent | Set-Content -Path (Join-Path $chopinPath $configFileName)
		} catch {
			Write-Error "FAILED: '$($MyInvocation.MyCommand.Name)': $_"
		}
	}
}

Set-Chopin -configFileName "appsettings.json"
