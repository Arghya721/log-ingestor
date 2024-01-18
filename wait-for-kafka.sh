#!/bin/bash
# wait-for-kafka.sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until (echo > /dev/tcp/$host/$port) >/dev/null 2>&1; do
  >&2 echo "Kafka is unavailable - sleeping"
  sleep 1
done

>&2 echo "Kafka is up - executing command"
exec $cmd
