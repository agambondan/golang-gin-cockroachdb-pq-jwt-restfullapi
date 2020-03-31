package seed

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func Load(db *sql.DB) {
	var err error
	//_, err = db.Exec("DROP TABLE IF EXISTS accounts")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS users")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS post")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS role")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)");
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id uuid PRIMARY KEY, created_at date, updated_at date, deleted_at date, full_name VARCHAR(55) not null, username VARCHAR(55) unique not null, password VARCHAR(255) not null, email VARCHAR(55) unique not null, role_id int not null )")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS post (id uuid PRIMARY KEY, created_at date, updated_at date, deleted_at date, title VARCHAR(255) not null, content text not null, author_id uuid not null)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS role (id int PRIMARY KEY, created_at date, updated_at date, deleted_at date, name VARCHAR(55) unique not null)")
	if err != nil {
		log.Fatal(err)
	}
	id := 29
	date := time.Now()
	roleName := "admin"
	_, err = db.Exec("INSERT INTO role VALUES ($1, $2, $3, $4, $5)", id, date, date, nil, roleName)
	if err != nil {
		fmt.Println(err.Error())
	}
}
