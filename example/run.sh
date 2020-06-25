#!/bin/sh

curl --fail -o /dev/null --silent http://localhost:8080/artifacts || (echo "Please 'make run' before running this script" && exit 1)

rm -f collation-*.tar

set -ex

printf "\nwrite first artifact\n"
tar -C v1.2.0 -c . | curl -X PUT --fail --data-binary @- http://localhost:8080/artifacts/v1.2.0
curl --fail http://localhost:8080/artifacts

printf "\nwrite second artifact\n"
tar -C v1.3.0 -c . | curl -X PUT --fail --data-binary @- http://localhost:8080/artifacts/v1.3.0
curl --fail http://localhost:8080/artifacts

printf "\nwrite third artifact\n"
tar -C v1.3.1 -c . | curl -X PUT --fail --data-binary @- http://localhost:8080/artifacts/v1.3.1
curl --fail http://localhost:8080/artifacts

rm -rf example-collations
mkdir -p example-collations/v1.2.0
mkdir -p example-collations/v1.2.0_v1.3.0_v1.3.1
mkdir -p example-collations/v1.3.1_v1.2.0_v1.3.0

printf "\ncollate first artifact\n"
curl -X GET --silent --fail http://localhost:8080/collation?artifacts=v1.2.0 | tar -xi -C example-collations/v1.2.0

printf "\ncollate first+second+third artifacts\n"
curl -X GET --silent --fail http://localhost:8080/collation?artifacts=v1.2.0,v1.3.0,v1.3.1 | tar -xi -C example-collations/v1.2.0_v1.3.0_v1.3.1

printf "\ncollate third+first+second artifacts\n"
curl -X GET --silent --fail http://localhost:8080/collation?artifacts=v1.3.1,v1.2.0,v1.3.0 | tar -xi -C example-collations/v1.3.1_v1.2.0_v1.3.0


printf "\noverwrite second artifact with first artifact\n"
tar -C v1.2.0 -c . | curl -X PUT --fail --data-binary @- http://localhost:8080/artifacts/v1.3.0
curl --fail http://localhost:8080/artifacts


printf "\ndelete first artifact\n"
curl -X DELETE --fail http://localhost:8080/artifacts/v1.2.0
curl --fail http://localhost:8080/artifacts

printf "\ndelete second artifact\n"
curl -X DELETE --fail http://localhost:8080/artifacts/v1.3.0
curl --fail http://localhost:8080/artifacts

printf "\ndelete third artifact\n"
curl -X DELETE --fail http://localhost:8080/artifacts/v1.3.1
curl --fail http://localhost:8080/artifacts
