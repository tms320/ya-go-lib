@echo ff
cls
go test -coverprofile=coverage.out
go tool cover -html=coverage.out