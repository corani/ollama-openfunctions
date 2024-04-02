#!/bin/bash

prog=$(realpath "$0")
rc=0

function print_usage {
  echo "Usage: $prog [options...]"
  echo
  echo "Available options:"
  echo "  -h            prints this help"
  echo "  -b            build app"
}

function do_echo {
    echo "[CMD] $*"
    TIMEFORMAT="[INFO] took %3lR"
    time "$@"
    code=$?
    if [ $code -ne 0 ]; then
        rc=$code
        echo "[ERROR] return code $rc"
    fi
}

if [ "$#" -eq "0" ]; then
    print_usage
    exit 0
fi

mkdir -p bin .cache models

while [ "$#" -gt "0" ]; do
  arg=$1
  shift

  case $arg in
    -h)
        print_usage
        exit 0
        ;;
    -b)
        do_echo go build -o bin/weather .
        ;;
    *)
        echo "[ERROR] unrecognized argument '$arg'"
        print_usage
        exit 1
        ;;
  esac
done

exit "${rc}"
