#!/bin/bash

## Extract runtime options from /proc and /sys
# The script formats a list of parameters taken from /proc
# and prints it in the yaml format.

number_check_re='^[0-9]+$'

if [ ! $# == "1" ]; then
	echo "usage: $0 /proc/sys/path/to/subfolder"
	exit
fi

# Filter only files that are writable
for option in `find $1 -maxdepth 1 -type f -perm /222 -printf "%f\n"`; do
	default_val=`cat $1/$option`
	echo "- name:    $option"
	echo "# path:    $1/$option"

	if [[ $default_val == "0" || $default_val == "1" ]]; then
		echo "  type:    string"
    echo "  default: \"$default_val\""
    echo "  choices: [\"0\", \"1\"]"
  elif [[ $default_val =~ $number_check_re ]]; then
		echo "  type:    int"
    echo "  default: $default_val"
    echo "  only:    [$default_val]"
  else
		echo "  type:    string"
    echo "  default: \"$default_val\""
    echo "  only:    [\"$default_val\"]"
	fi

  echo "  when:    test"
	echo ""
done