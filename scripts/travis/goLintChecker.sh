diff -u <(echo -n) <(golint ./... | grep -v vendor | grep -v third_party | grep -v stutters | grep -v _test)
if [ $? == 0 ]; then
	echo "No GoLint warnings found"
	exit 0
else
	echo "GoLint Warnings found"
	exit 1
fi
