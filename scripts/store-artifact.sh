#!/bin/sh

# ./store-artifact.sh "https://artifacts.internal" v1.2.3 ~/projects/some-project/dist

set -x

tar -C "$3" -c . | curl -X PUT --fail --data-binary @- "$1/artifacts/$2"
