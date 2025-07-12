package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"os/exec"
	"encoding/json"

	"github.com/joho/godotenv"
)

type Package struct {
	Path  string              `json:"path"`
	Props map[string][]string `json:"props"`
}

func listFiles(src_dir, ext_file string) []string {
	root := os.DirFS(src_dir)
	mdFiles, err := fs.Glob(root, ext_file)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, v := range mdFiles {
		files = append(files, path.Join(src_dir, v))
	}
	return files
}

func execPowerhsellCommand(userCommand string) string {
	cmd := exec.Command("powershell", "-Command", userCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing PowerShell command for package: %s", err)
		log.Printf("PowerShell error output: %s\n", string(output))
	}
	return string(output)
}

func parseEnvArray(key string) []string {
	raw := os.Getenv(key)
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	repositories := parseEnvArray("REPOSITORIES")
	excluded_packages := parseEnvArray("EXCLUDED_PACKAGES")
	extensions := parseEnvArray("EXTENSIONS")

	// Create the output file to store filtered packages
	newFile, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("Error creating output list: %v", err)
	}
	defer newFile.Close()

	// Iterate through each file found in the directory
	data_dir := "data"
	files := listFiles(data_dir, "*.json")
	for _, filePath := range files {
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Unable to read file: %v", err)
			return
		}

		var packages []Package
		err = json.Unmarshal(jsonData, &packages)
		if err != nil {
			log.Fatalf("JSON unmarshal error: %v", err)
			return
		}

		// Iterate over each package in the JSON file
		for _, jsonPckg := range packages {
			skip := false
			// Check the repository, excluded_packages and extension filters
			for _, repo := range repositories {
				if strings.HasPrefix(jsonPckg.Path, repo) {
					for _, ext := range extensions {
						if strings.HasSuffix(jsonPckg.Path, ext) {
							for _, pack := range excluded_packages {
								if strings.Contains(jsonPckg.Path, pack) {
									skip = true
									break
								}
							}
						}
					}
				}
			}
			// Skip the current package if it matches the exclusion filter
			if skip {
				continue
			}

			id := ""
			ver := ""
			// Check if the package has 'nuget.id' and 'nuget.version' properties
			for key, value := range jsonPckg.Props {
				if key == "nuget.id" && len(value) > 0 {
					id = value[0]
				}
				if key == "nuget.version" && len(value) > 0 {
					ver = value[0]
				}
			}

			// Proceed only if both id and version are found
			if id != "" && ver != "" {
				// nameVersion := fmt.Sprintf("%s:%s", id, ver)
				chocoCommand := fmt.Sprintf("choco find %s --exact --version %s", id, ver)

				// Execute the Chocolatey command and check if the package is found
				execOutput := execPowerhsellCommand(chocoCommand)

				if strings.Contains(execOutput, "0 packages found") {
					fmt.Printf("Package: %s:%s not found in public repository, adding to your list...\n", id, ver)

					// Write the package name to the output file
					if _, err := fmt.Fprintln(newFile, jsonPckg.Path); err != nil {
						log.Printf("Error writing to output file: %v", err)
						return
					}
				}
			}
		}
	}
}
