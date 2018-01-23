diff -u <(echo -n) <(goconst ./... | grep -v vendor | grep -v third_party)
if [ $? == 0 ]; then
	echo "No goConst problem"
	exit 0
else
	echo "Has goConst Problem"
	exit 1
fi
