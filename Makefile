run_mysql:
	docker run --name mysql -p 3306:3306 mysql

start_mysql:
	docker container start mysql

go_get:
	go get github.com/go-sql-driver/mysql
	go get github.com/justinas/alice
	go get github.com/go-playground/form/v4

.PHONY: run_mysql start_mysql go_get