diff -u <(echo -n) <(find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" | xargs gofmt -s -d)
if [ $? == 0 ]; then
	echo "Hurray....all code is formatted properly..."
	exit 0
else
	echo "There is issues's with the code formatting....please run go fmt on your code"
	exit 1
fi
