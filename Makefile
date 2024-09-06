run_mysql:
	docker run --name mysql -p 3306:3306 mysql

start_mysql:
	docker container start mysql

go_get:
	go get github.com/go-sql-driver/mysql
	go get github.com/justinas/alice
	go get github.com/go-playground/form/v4
	go get github.com/alexedwards/scs/v2
	go get github.com/alexedwards/scs/mysqlstore
	go get golang.org/x/crypto/bcrypt
	go get github.com/justinas/nosurf

mysql_root:
	docker exec -it mysql mysql -uroot -prootpwd

.PHONY: run_mysql start_mysql go_get mysql_root