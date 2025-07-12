function Trace-Conflicts {
	[cmdletBinding()]
	Param (
		[Parameter(Mandatory = $true)]
		[ValidateNotNullOrEmpty()]
		[string]$conflictPath
	)

	BEGIN {
		$conflictFolders = Get-ChildItem -Path $conflictPath -Directory | ForEach-Object { $_.FullName }
	}

	PROCESS {
		if ($conflictFolders.Count -gt 0) {
			Write-Host "Conflict found...`n"

			foreach ($conflict in $conflictFolders) {
				$refFolder = Get-ChildItem -Path $conflict -Filter "*-fromLocalRepo" | Select-Object -First 1
				$difFolder = Get-ChildItem -Path $conflict -Filter "*-toBeInternalized" | Select-Object -First 1

				git --no-pager diff --no-index $refFolder.FullName $difFolder.FullName
			}
		} else {
			Write-Host "Conflict not found..."
		}
	}
}


