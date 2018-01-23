diff -u <(echo -n) <(find . -type f -not -path "./vendor/*" -not -path "./third_party/*" -print0 | xargs -0 misspell)
if [ $? == 0 ]; then
	echo "No Misspell found"
	exit 0
else
	echo "Misspell found"
	exit 1
fi
