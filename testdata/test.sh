#!/bin/bash
# test.sh - A simple test script for goshellcheck

echo "Hello, World!"

# Potential issue: unquoted variable
name=World
echo Hello $name

# Another potential issue: using backticks instead of $()
files=`ls -la`
