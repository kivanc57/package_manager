package main

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"log"
	"strings"
	"os/exec"
	"net/http"
	"path/filepath"

	"github.com/joho/godotenv"
)

func parseEnvArray(key string) []string {
	raw := os.Getenv(key)
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func openFile(inputPath string) *os.File {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Failed to open: %s", err)
	}
	return inputFile
}

func sendGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed GET request: %v", err)
	}
	return resp, nil
}

func checkStatusCode(resp *http.Response) bool {
	return resp.StatusCode == http.StatusOK
}

func createFile(outputPath string) *os.File {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Failed to create file: %s", err)
		return nil
	}
	return outputFile
}

func createFolders(outputFolder string) {
	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directories: %s", err)
	}
}

func copyFile(outputFile *os.File, body io.Reader) error {
	_, err := io.Copy(outputFile, body)
	if err != nil {
		return fmt.Errorf("Failed to write: %s", err)
	}
	return nil
}

func downloadFile(url, outputFolder, outputPath string) error {
	resp, err := sendGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !checkStatusCode(resp) {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	createFolders(outputFolder)

	outFile := createFile(outputPath)
	if outFile == nil {
		return fmt.Errorf("Failed to create file: %s", outputPath)
	}
	defer outFile.Close()

	err = copyFile(outFile, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Download Completed: %s...\n", outputPath)
	log.Printf("-------------------------------------")
	return nil
}

func execPowerhsellCommand(userCommand string) string {
	cmd := exec.Command("powershell", "-Command", userCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing PowerShell command: %s", err)
		log.Printf("PowerShell error output: %s\n", string(output))
	}
	return string(output)
}

func getAbsPath(relPath string) string {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		log.Fatalf("Error converting .nupkg file:%s", err)
	}
	return absPath
}

func executePowershellScripts(commands []string) {
	for _, cmd := range commands {
		log.Printf("Executing PowerShell Command: %s", cmd)
		output := execPowerhsellCommand(cmd)
		log.Printf("PowerShell Command Output: %s", output)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	links := parseEnvArray("LINKS")

	inputPath := filepath.Join("data", "example.txt")
	inputFile := openFile(inputPath)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		pckg := scanner.Text()
		i := strings.Index(pckg, ":")
		pckgName, pckgVersion := pckg[:i], pckg[i+1:]

		for _, link := range links {
			isDownloaded := true
			outputFolder := filepath.Join("output", pckgName, pckgVersion)
			outputFile := fmt.Sprintf("%s.%s%s", pckgName, pckgVersion, ".nupkg")
			outputPath := filepath.Join(outputFolder, outputFile)
			
			pckgLink := fmt.Sprintf("%s&path=%s", link, outputFile)
			err := downloadFile(pckgLink, outputFolder, outputPath)
			if err != nil {
				log.Printf("Error downloading file: %s", outputPath)
				isDownloaded = false
			}

			if isDownloaded {
				absOutputPath := getAbsPath(outputPath)
				newPath := strings.TrimSuffix(absOutputPath, ".nupkg") + ".zip"

				renameCmd := fmt.Sprintf(`if (Test-Path "%s") { Remove-Item "%s" -Force }; Rename-Item -Path "%s" -NewName "%s"`, newPath, newPath, absOutputPath, newPath)
				extractCmd := fmt.Sprintf("Expand-Archive -Path \"%s\" -DestinationPath \"%s\"", newPath, outputFolder)
				detectInstallersCmd := fmt.Sprintf(`Set-Location -Path "%s"; Get-ChildItem -Recurse -File | Where-Object { $_.Extension -eq '.exe' -or $_.Extension -eq '.msi' } | ForEach-Object { Write-Host "Installer found: $($_.FullName)" }`, outputFolder)
				cleanupFilesCmd := fmt.Sprintf(`Set-Location -Path "%s"; Get-ChildItem -Recurse -File | ForEach-Object { if ($_.Extension -ne '.ps1' -and $_.Extension -ne '.nuspec') { Remove-Item $_.FullName -Force } }`, outputFolder)
				removeEmptyDirsCmd := fmt.Sprintf(`Set-Location -Path "%s"; Get-ChildItem -Recurse -Directory | Where-Object { $_.GetFiles().Count -eq 0 -and $_.GetDirectories().Count -eq 0 } | Remove-Item -Force`, outputFolder)
				
				executePowershellScripts([]string{renameCmd, extractCmd, detectInstallersCmd, cleanupFilesCmd, removeEmptyDirsCmd})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
}



