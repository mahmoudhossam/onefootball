# onefootball
The onefootball backend engineer task

To run the solution file run `go run teams.go`

To run tests run `go test`

TODO
------

1. Remove the first for loop in `main()` and replace it with something cleaner.
2. Proper error handling, possibly logging errors instead of just calling `panic()`
3. Synchronize writes to `roster` for better handling of concurrent writes.
4. Refactor functions to make them more testable and decrease dependency.
5. Write more tests.
