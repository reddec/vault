#!/usr/bin/env bash
set -e -o pipefail

ARGS="-u redis://127.0.0.1:6379 -u redis://127.0.0.1:6380 -u redis://127.0.0.1:6381"

docker-compose stop
docker-compose rm --force -v
docker-compose up -d
sleep 2

dd if=/dev/urandom of=sample.txt bs=10M count=1

cat sample.txt | go run ../cmd/vault-push/main.go $ARGS[@] sample

go run ../cmd/vault-pull/main.go $ARGS[@] sample > sample2.txt

docker-compose stop
docker-compose rm --force -v

if cmp -s sample.txt sample2.txt; then
    echo "Success!"
else
    echo "Failed: file are different"
    exit 1
fi

