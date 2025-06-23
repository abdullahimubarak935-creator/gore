#!/bin/sh
# We're not clearing out the state properly on `Run` at the moment, so all the
# tests need to be in a separate process

set -e

if [ -n "$1" ] ; then
	echo "Running tests in $1"
	COUNT=1
	if [ -n "$2" ] ; then
		COUNT=$2
	fi
	echo "Running Test $1 $COUNT times"
	for i in $(seq $COUNT) ; do
		go test -v -count 1 -run "^${1}\$" .
	done
else
	TEST=$(go test -test.list '.*' | grep '^Test')
	echo "All tests: $TEST"
	set -e
	for f in $TEST ; do
		echo "Running $f..."
		go test -v -count 1 -run "^${f}\$" .
	done
fi
