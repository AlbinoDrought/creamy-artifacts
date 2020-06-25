#!/bin/sh

# ./remove-artifact.sh "https://artifacts.internal" v1.2.3

set -x

curl -X DELETE --fail "$1/artifacts/$2"
