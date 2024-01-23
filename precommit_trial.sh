# Format all files, and sort go mod
echo "====Format & Sync===="
go work sync

current_dir=$(pwd)
for i in $(go list -m -f '{{.Dir}}')
do 
  cd $i
  go fmt
  go mod tidy
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
staticcheck $(go list -m)

# Run Even stricter linting
gocritic check $(go list -m)

echo "====Error handling===="
# Identifies any non-exhaustive case statements
exhaustive $(go list -m)

# Check for potential Nil panics
# nilaway $(go list -m)

echo "====Performance===="
# Identifies areas where pre-allocating slices could improve performance
prealloc $(go list -m)

echo "====Maintainability===="
# Detects frequently used strings which could be constants
goconst -min-occurrences 5 ./...

# Identifies mixed pointer receivers to an interface
smrcptr $(go list -m)

echo "====Complexity===="
# Limit Cognitive complexity
gocognit -over 15 -top 10 $(go list -m -f '{{.Dir}}')

