#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

#######################################################
# This script is used to package ui/dist in go source #
#######################################################

rm -f cli/commands/init/rice-box.go
rice embed-go -i ./cli/commands/init
