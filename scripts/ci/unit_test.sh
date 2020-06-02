#!/bin/sh
set -e
# Make the Coverage File
echo "mode: atomic" > coverage.txt
# Make Necessary directories needed by Test (Ideally it should get created automatically but Travis is not allowing to create it using os.MkdriAll)
# I know this is insane but nothing can be done

go test $(go list ./... |  grep -v third_party | grep -v examples) -cover -covermode atomic -coverprofile coverage.out -timeout=30m

sed '1d;$d' coverage.out >> coverage.txt
