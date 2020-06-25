#!/bin/sh

# ./pull-artifacts.sh "https://artifacts.internal" "v1.2.3,v1.2.4,v1.3.0,v1.3.1,v1.3.2" ~/projects/some-project/collated-dist

set -x

mkdir -p "$3"
curl -X GET --fail "$1/collation?artifacts=$2" | tar -xi -C "$3"
