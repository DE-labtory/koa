#!/bin/bash

# goimports test
diff <(goimports -d $(find . -type f -name '*.go' -not -path "*/vendor/*")) <(printf "")

if [ $? -ne 0 ]; then
echo "goimports format error" >&2
exit 1
fi

# run go test
go test -v -mod=vendor ./...

if [ $? -ne 0 ]; then
echo "go test fail" >&2
exit 1
fi

exit 0
