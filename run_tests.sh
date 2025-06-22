#!/bin/sh
# We're not clearing out the state properly on `Run` at the moment, so all the
# tests need to be in a separate process
TEST=$(go test -test.list '.*' | grep '^Test')
echo "All tests: $TEST"
set -e
for f in $TEST ; do
	echo "Running $f..."
	go test -v -count 1 -run "^${f}\$" .
done
