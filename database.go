package main

import (
	"database/sql"
	"fmt"
)

type Database struct {
	mySQLconn *sql.DB
	Connected bool
	Status string

}

var DbConn = &Database{}

func (d *Database) ping() error {
	return d.mySQLconn.Ping()
}

func (d *Database) GetConnection() interface{} {
	return d.mySQLconn
}

func (d *Database) Connect(host, port, uname, pass, dbname string) {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		uname,
		pass,
		host,
		port,
		dbname,
	)

	c, err := sql.Open("mysql", dbSource)
	if err != nil {
		d.Connected = false
		d.Status = fmt.Sprintf("Status: %s %s@%s; %s", "Unable to connect", host, dbname, err.Error())
		return
	}

	d.mySQLconn = c

	err = d.mySQLconn.Ping()
	if err != nil {
		d.Connected = false
		d.Status = fmt.Sprintf("Status: %s %s@%s; %s", "Unable to connect", host, dbname, err.Error())
	}

	d.Connected = true
}
