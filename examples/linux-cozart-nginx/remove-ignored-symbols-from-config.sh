#!/bin/bash

ignored_config="$2"

# For each line in the specialized config
while IFS= read -r line; do
  # Check if the symbol is in the ignored config
  if grep -q "$line" "$ignored_config"; then
    # If it is, don't write it to the output
    continue
  fi
  echo "$line"
done < "$1"
