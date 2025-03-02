#!/bin/bash

vanilla_config="$2"

# For each line in the specialized config
while IFS= read -r line; do
  # Check if the symbol is already in the vanilla config
  if grep -q "$line" "$vanilla_config"; then
    # If it is, write it to the output
    echo "$line"
  fi
done < "$1"
