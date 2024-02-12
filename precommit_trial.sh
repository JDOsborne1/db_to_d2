# Format all files, and sort go mod
echo "====Format & Sync===="
go work sync

current_dir=$(pwd)
for i in $(go list -m -f '{{.Dir}}')
do 
  cd $i
  go fmt
  go mod tidy
  # go install golang.org/x/tools/cmd/goimports@latest
  goimports -w .
  cd $current_dir
done

# Test all packages
echo "====Test===="
go test $(go list -m) -cover

echo "====Health & Standards===="
# Run standard health checks
go vet $(go list -m)

# Stricter linting
# go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck $(go list -m)

# Run Even stricter linting
# go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
gocritic check $(go list -m)

echo "====Error handling===="
# Identifies any non-exhaustive case statements
# go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
exhaustive $(go list -m)

# Check for potential Nil panics
# nilaway $(go list -m)

echo "====Performance===="
# Identifies areas where pre-allocating slices could improve performance
# go install github.com/alexkohler/prealloc@latest
prealloc $(go list -m)

echo "====Maintainability===="
# Detects frequently used strings which could be constants
# go get github.com/jgautheron/goconst/cmd/goconst
goconst -min-occurrences 5 ./...

# Identifies mixed pointer receivers to an interface
# go install github.com/nikolaydubina/smrcptr@latest
smrcptr $(go list -m)

echo "====Complexity===="
# Limit Cognitive complexity
# go install github.com/jgautheron/goconst/cmd/goconst@latest
gocognit -over 15 -top 10 $(go list -m -f '{{.Dir}}')

