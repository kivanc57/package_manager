
# Package Manager

## Overview

This project provides a modular Go/Python and Powershell based framework for a very efficient package manager, where the user can filter, exclude, download and upload necessary packages within a given artifactory.

- **Downloader**: Downloads files from specified URLs, manages extraction, and processes packages.
- **Package Manager**: Processes JSON package lists, checks availability in Chocolatey repositories, and logs missing packages.


## Installation

Clone this repository and navigate to the project directory:

```bash
git clone https://github.com/your-username/your-project.git
cd your-project
```

Initialize Go modules:

```bash
go mod tidy
```

## Environment Configuration

Create a `.env` file at the root with necessary environment variables:

```env
LINKS=https://example.com/download?token=12345
REPOSITORIES=https://repo.example.com
EXCLUDED_PACKAGES=package1,package2
EXTENSIONS=.zip,.exe,.nupkg
```

## Running the Scripts

### Downloader

To execute the downloader tool:

```bash
go run cmd/downloader/main.go
```

### Package Manager

To execute the package manager:

```bash
go run cmd/packagemanager/main.go
```

## Python Modular Downloader Script

This Python script downloads and extracts package files from given URLs and processes them in a modular way.

### Usage

Ensure you have the `.env` file with `LINKS` variable set, and the input file `data/final_list.txt` ready.

Run the script with:

```bash
python main.py
```

### Requirements

- Python 3.x
- `requests`
- `python-dotenv`

Install dependencies with:

```bash
pip install requests python-dotenv
```
---

## Output

Results will be stored in the `output` directory. Ensure it exists or will be created automatically by the scripts.

## Dependencies

- [godotenv](https://github.com/joho/godotenv) for loading `.env` variables

Install dependencies using:

```bash
go get github.com/joho/godotenv
```
