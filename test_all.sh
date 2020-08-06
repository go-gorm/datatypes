#!/bin/bash -e

dialects=("postgres" "mysql" "mssql" "sqlite")

for dialect in "${dialects[@]}" ; do
    GORM_DIALECT=${dialect} go test  --tags "json1"
done
