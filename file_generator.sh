# !/bin/bash

location=_generated_files/

mkdir $location
cd $location

for i in $(seq -w 1 200); do
    folder_name="folder_${i}"
    mkdir -p "$folder_name"

    for j in $(seq -w 1 200); do
      file_name="${folder_name}/file_${j}.txt"
      head /dev/urandom | tr -dc A-Za-z0-9 | head -c 1000 > "$file_name"
    done

done

echo "Folders and files created successfully."
