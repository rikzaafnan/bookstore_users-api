package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ClientDB *sql.DB
	// username = os.Getenv("mysql_users_username")
	// password = os.Getenv("mysql_users_password")
	// host     = os.Getenv("mysql_users_host")
	// schema   = os.Getenv("mysql_users_schema")
	username = mysql_users_username
	password = mysql_users_password
	host     = mysql_users_host
	schema   = mysql_users_schema
)

const (
	mysql_users_username = "root"
	mysql_users_password = "root_password"
	mysql_users_host     = "127.0.0.1:35200"
	mysql_users_schema   = "users_db"
)

func init() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host, schema)

	log.Println(fmt.Sprintf("about to connect in host %s ", host))

	var err error
	ClientDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = ClientDB.Ping(); err != nil {
		panic(err)
	}

	log.Println("dataabse successfully configuration")

}
