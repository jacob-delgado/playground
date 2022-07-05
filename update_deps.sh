grep -v indirect go.mod | grep -v require | grep -v \) | grep -v "go 1.18" | grep -v "github.com/jacob-delgado/playground" | awk {print } | awk NF | xargs -L1 go get -u
