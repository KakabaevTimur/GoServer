package main

import (
	"database/sql"
	"hash/fnv"
	"strconv"

	_ "github.com/lib/pq"
)

type DB struct {
	postgres *sql.DB
	urls     map[string]string
	get      func(string) string
	put      func(string) string
}

func (db *DB) init(usePostgres bool) {
	if usePostgres {
		connStr := "user=postgres password=12345 dbname=URLs sslmode=disable"
		var err error
		db.postgres, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		db.get = db.getFromPostgres
		db.put = db.putToPostgres

	} else {
		db.get = db.getFromMap
		db.put = db.putToMap
	}
	db.urls = make(map[string]string)
}

func (db *DB) deinit() {
	if db.postgres != nil {
		db.postgres.Close()
	}
}

func (db *DB) getFromMap(req string) string {
	return db.urls[req]
}

func (db *DB) getFromPostgres(req string) string {
	rows, err := db.postgres.Query("select url from urls where hash = '" + req + "'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var url string
	for rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			continue
		}
	}
	return url
}

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 10)
}

func (db *DB) putToMap(req string) string {
	db.urls[hash(req)] = req
	return hash(req)
}

func (db *DB) putToPostgres(req string) string {
	_, err := db.postgres.Exec("insert into URLs (hash, url) values ($1, $2)",
		hash(req), req)
	if err != nil {
		panic(err)
	}
	return hash(req)
}
