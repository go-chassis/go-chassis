diff -u <(echo -n) <(find . -type d -not -path "./vendor/*" | xargs deadcode)
if [ $? == 0 ]; then
	echo "Hurray....all code's are reachable and utilised..."
	exit 0
else
	echo "There are some deadcode in the project...please remove the unused code"
	exit 1
fi
