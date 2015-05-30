Fast Rename
======================

Install
--------------

    go get -u -x github.com/c9s/fsrename


Usage
---------------


    fsrename -match "_stmt.go" -replace "_stmt.go" src/c6

    fsrename -match "_[a-z]*.go" -replace ".go" src/c6
