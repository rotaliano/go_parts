package sqlx_x

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sqlx.DB

func Test() {
	db, err := sqlx.Open("sqlite3", ":memory")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	schema := `CREATE TABLE place (
		country text,
		city text NULL,
		telcode integer);`

	result, _ := db.Exec(schema)
	log.Println("result:", result)

	cityState := `INSERT INTO place (country, telcode) VALUES (?, ?)`
	countryCity := `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
	db.MustExec(cityState, "Hong Kong", 852)
	db.MustExec(cityState, "Singapore", 65)
	db.MustExec(countryCity, "Sounth Africa", "Johannesburg", 27)

	//rows, err := db.Query("SELECT country, city, telcode FROM place")
	//for rows.Next() {
	//	var country string
	//	var city sql.NullString
	//	var telcode int
	//	err = rows.Scan(&country, &city, &telcode)
	//	log.Println(country, city, telcode)
	//}

	type Place struct {
		Country       string
		City          sql.NullString
		TelephoneCode int `db:"telcode"`
	}

	rows, err := db.Queryx("SELECT * FROM place")
	for rows.Next() {
		var p Place
		err = rows.StructScan(&p)
		log.Println(p)
	}

	row := db.QueryRow("SELECT * FROM place WHERE telcode=? LIMIT 1", 852)
	var telcode int
	err = row.Scan(&telcode)
	log.Println("telcode:", telcode)

	var p Place
	err = db.QueryRowx("SELECT city, telcode FROM place LIMIT 1").StructScan(&p)
	log.Println("db.QueryRowx:", p)

	p = Place{}
	pp := []Place{}

	err = db.Get(&p, "SELECT * FROM place LIMIT 1")
	log.Println("db.Get: ", p)

	err = db.Select(&pp, "SELECT * FROM place WHERE telcode > ?", 50)
	log.Println("db.Select:", pp)

	// count
	var id int
	err = db.Get(&id, "SELECT count(*) FROM place")
	log.Println("db.Get:", id)

	// in
	var levels = []int{852, 6, 7}
	qeury, args, err := sqlx.In("SELECT * FROM place WHERE telcode IN (?);", levels)
	if err != nil {
		log.Fatalln(err)
	}
	qeury = db.Rebind(qeury)
	rowss, err := db.Queryx(qeury, args...)
	for rowss.Next() {
		var p Place
		rowss.StructScan(&p)
		log.Println("in:", p)
	}

	tx := db.MustBegin()
	tx.MustExec("UPDATE place SET telcode = 123")
	err = tx.Commit()

}
