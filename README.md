# Graves

go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


migrate -path migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/your_database" up
