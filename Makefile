run-accounts: accounts-build accounts-start

accounts-build:
	@echo " >> building binaries"
	@go build -o bin/accounts cmd/accounts/app.go

accounts-start:
	@echo " >> starting binaries"
	@./bin/accounts