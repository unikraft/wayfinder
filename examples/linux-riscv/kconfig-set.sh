#!/bin/bash

# kconfig-set.sh $kconfig_file $variable $value

kconfig_file=$1
kconfig_variable=$2
kconfig_value=$3

is_present=$(cat $kconfig_file | grep -P "$kconfig_variable" | wc -l)

if [ "$is_present" -eq "0" ]; then
  echo "[W] Variable $kconfig_variable not present in $kconfig_file, appending at the end"
  echo "$kconfig_variable=$kconfig_value" >> $kconfig_file
else
  if grep -q "$kconfig_variable=$kconfig_value" $kconfig_file; then
    echo "[I] Variable $kconfig_variable already set to $kconfig_value, doing nothing"
  else
    echo "[I] Setting variable $kconfig_variable to $kconfig_value"
    sed -i -e "s/# $kconfig_variable is not set/$kconfig_variable=$kconfig_value/" $kconfig_file
    sed -i -e "s/$kconfig_variable=.*/$kconfig_variable=$kconfig_value/" $kconfig_file
  fi
fi
