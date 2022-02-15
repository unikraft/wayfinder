#!/bin/bash

# kconfig-set.sh $kconfig_file $variable $value

kconfig_file=$1
kconfig_variable=$2
kconfig_value=$3

# quick early check...
is_present=$(cat $kconfig_file | grep -P "$kconfig_variable" | wc -l)
cp $kconfig_file /tmp/kconfig_file.tmp

if [ "$is_present" -eq "1" ]; then
  echo "[W] Variable $kconfig_variable not present in $kconfig_file appending at the end"
  echo "$kconfig_variable=$kconfig_value" >> $kconfig_file
else
  echo "[I] Replacing..."
  sed -i -e "s/# $kconfig_variable is not set/$kconfig_variable=$kconfig_value/" $kconfig_file
  sed -i -e "s/$kconfig_variable=.*/$kconfig_variable=$kconfig_value/" $kconfig_file
fi
