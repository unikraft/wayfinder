#!/bin/bash

# Set the paths to your kernel config files
ALLDEFCONFIG=.config-riscv64-def
ALLNOCONFIG=.config-riscv64-allno

# Parse the two configs into arrays of enabled options
alldefconfig=($(grep -v "^#" $ALLDEFCONFIG | grep -o "CONFIG_[^=]*" | sort))
allnoconfig=($(grep -v "^#" $ALLNOCONFIG | grep -o "CONFIG_[^=]*" | sort))

# Find the options that are enabled in the default config but not in the minimal one
diff=($(comm -23 <(echo "${alldefconfig[*]}") <(echo "${allnoconfig[*]}")))

echo "Options enabled in the default config but not in the minimal one:"
echo "${diff[@]}"
echo "Total: ${#diff[@]}"

# Check if each option is a top-level variable that is not a dependency
diff_clean=()

for option in "${diff[@]}"
do
    # Check if the option exists in the minimal config but simply disabled. If so, add it to diff_clean
    if grep -q "# $option is not set" $ALLNOCONFIG; then
        diff_clean+=($option)
    fi
done

# Print the results
echo "Top-level options enabled in the default config but not in the minimal one:"
echo "${diff_clean[@]}"
echo "Total: ${#diff_clean[@]}"

# Save the results to a file
echo "${diff_clean[@]}" | tr ' ' '\n' > diff_clean.txt

# Create a yaml file for the options
for option in "${diff_clean[@]}"
do
    echo "- name: $option" >> diff_clean.yaml
    echo "  type: string" >> diff_clean.yaml
    echo "  default: \"y\"" >> diff_clean.yaml
    echo "  only: [\"y\", \"n\"]" >> diff_clean.yaml
    echo "" >> diff_clean.yaml
done
