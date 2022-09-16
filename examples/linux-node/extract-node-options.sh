#!/bin/zsh

## Extract node options from the internet
# The script formats a list of parameters taken from the node website
# and prints it in the yaml format.

node_link="https://nodejs.org/docs/latest-v8.x/api/cli.html"
v8_link="https://gist.githubusercontent.com/sarupbanskota/a68e8148aa4cdc95e66a1b0e93df48ef/raw/701463e3a70b4cb2cbd679d0fe5b07df73adcee7/node_8_v8_options"

node_options_list=$(curl -s $node_link | grep "<li><code>--" | grep -o -e "--[a-z1-9-]*")
v8_options_list=$(curl -s $v8_link | tail -n +18)

FILE="$(mktemp)"
echo "$node_options_list" > $FILE

while read -r option; do
  option=$(echo ${option:2} | tr '[:lower:]' '[:upper:]' | tr '-' '_')

  echo "- name:    $option"
  echo "  type:    string"
  echo "  default: \"n\""
  echo "  only:    [\"n\", \"y\"]"
  echo "  when:    test"
  echo ""
done < $FILE

new_line=0
skip_first=0
name=""
desc=""

echo "$v8_options_list" > $FILE

while read -r option; do
  if [[ $new_line == 0 ]]; then
    name=$(echo ${option:2} | grep -o -e "[a-z_-]*" | head -n 1 | tr '[:lower:]' '[:upper:]' | tr '-' '_')
    desc="$(echo ${option})"
    new_line=1
  else
    type=$(echo ${option} | grep -o "bool")
    if [[ $type == "bool" ]]; then
      echo "- name:    $name"
      echo "  description: |"
      echo "  $desc"
      echo "  type:    string"

      default=$(echo ${option} | grep -o -e "default: .*" | grep -o "false")
      if [[ $default == "false" ]]; then
        echo "  default: \"n\""
        echo "  only:    [\"n\", \"y\"]"
      else
        echo "  default: \"y\""
        echo "  only:    [\"y\", \"n\"]"
      fi

      echo "  when:    test"
      echo ""
    fi
    new_line=0
  fi
done < $FILE

rm $FILE
