#!/bin/sh

project_dir="$(cd "$(dirname "$0")/.." && pwd)"

for i in $(seq 100); do
  set -- "$@" $i
  $project_dir/tools/count.sh "$@"
done
