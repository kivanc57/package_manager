import os
import requests
from zipfile import ZipFile
from dotenv import load_dotenv

import time # Library for counting time

load_dotenv()

def join_paths(*args):
  return os.path.join(*args)

def remove_file(path):
  try:
    if os.path.exists(path):
      os.remove(path)
  except Exception as e:
    print(f"Error removing file: {e}")

def create_directory(target_dir):
  try:
    os.makedirs(target_dir, 0o755, exist_ok=True)
  except OSError as e:
    print(f"Directory {target_dir} can not be created: {e}")

def write_file(url, output_path):
  try:
    response = requests.get(url)
    if response.status_code == 200:
      with open(output_path, "wb") as file:
        file.write(response.content)
      print(f"Download succeeded from: {url}")
      print(f"File is written to: {output_path}")
      return True
    else:
      print(f"Download failed from: {url}")
      print("------------------------------------------------------------------")
      return False
  except Exception as e:
    print(f"Error downloading file: {e}")
    print("------------------------------------------------------------------")
    return False

def convert_package(src_pckg, target_pck):
  try:

    os.rename(src_pckg, target_pck)
    return target_pck
  except FileNotFoundError:
    print(f"Error: File not found: {src_pckg}")
    return None
  except Exception as e:
    print(f"An error occurred: {e}")
    return None

def extract_zip(zip_path, output_folder):
  is_placeholder_present = False

  try:
    with ZipFile(zip_path, 'r') as zObject:
      for file in zObject.namelist():

        if file.endswith(".nuspec") or file.endswith(".ps1"):
          zObject.extract(file, output_folder)

        elif (file.endswith(".exe") or file.endswith(".msi")) and not is_placeholder_present:
          is_placeholder_present = True
          placeholder_path = join_paths(output_folder, "placeholder")
          with open(placeholder_path, "w") as pHolder:
            pHolder.write("Installer found")
    remove_file(zip_path)
    print(f"Extraction completed: {zip_path}")
  except FileNotFoundError:
    print(f"Error: File not found: {zip_path}")
  except Exception as e:
    print(f"An error occurred: {e}")

def main():
  input_file = join_paths("data", "final_list.txt")
  links = os.getenv("LINKS").split(",")
  count = 0

  with open(input_file, "r") as f:
    for pckg in f:
      pckg = pckg.strip()
      if not pckg or any(c.isspace() for c in pckg):
        continue
      i = pckg.index(":")
      pckg_name = pckg[:i]
      pckg_ver = pckg[i+1:].strip()

      for link in links:
        output_file = f"{pckg_name}.{pckg_ver}.nupkg"
        output_folders = join_paths("output", pckg_name, pckg_ver)
        output_path = join_paths(output_folders, output_file)
        zip_path = output_path.replace(".nupkg", ".zip")

        if os.path.exists(zip_path):
          print(f"Already processed: {zip_path}")
          break

        create_directory(output_folders)
        pckg_link = f"{link}&path={output_file}"
        if write_file(pckg_link, output_path):
          new_path = convert_package(output_path, zip_path)
          if new_path:
            extract_zip(new_path, output_folders)
            count += 1
            print("------------------------------------------------------------------")
          break
  print(f"PACKAGE COUNT: {count}")

if __name__ == "__main__":
  start = time.time()
  main()
  end = time.time()
  length = end - start
  print(f"Program took {length} seconds...")

