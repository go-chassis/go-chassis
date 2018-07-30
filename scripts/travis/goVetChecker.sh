diff -u <(echo -n) <(find . -type d -not -path "./vendor/*" -not -path "./third_party/*"| xargs go vet )
if [ $? == 0 ]; then
	echo "Hurray....all OKAY..."
	exit 0
else
	echo "There are some static issues in the project...please run go vet"
	exit 1
fi
