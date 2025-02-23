#!/bin/bash

while IFS= read -r line; do
  echo "- name:    $line"
  echo "  type:    string"
  echo "  only:    [\"y\", \"n\"]"
  echo "  default: \"y\""
  echo "  when:    build"
  echo ""
done < "$1"
