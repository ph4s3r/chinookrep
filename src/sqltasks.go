package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

func openDB(dbpath string) *sql.DB {

	if _, err := os.Stat(dbpath); err != nil {
		fmt.Printf("DB does not exist, exiting\n")
		os.Exit(1)
	}
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)

	return db
}

func sqlprep() {

	db := openDB("./chinook.db")
	defer db.Close()

	//2a. adds a new column “ReleaseDate” to the tracks table, moves on if the col exists

	_, err := db.Exec("ALTER TABLE tracks ADD COLUMN ReleaseDate DATETIME;")
	moveonErr(err)
	fmt.Println("2a done")

	//2b. fills it with the attached data

	csvfile, err := os.Open("track_release_dates.csv")
	checkErr(err)
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		_, err = db.Exec("UPDATE tracks SET (ReleaseDate) = (?) WHERE TrackId = (?) ;", record[1], record[0])
		checkErr(err)
	}

	fmt.Println("2b done")

	//2c. creates a table named “reports” with columns: FromDate, ToDate and Result. The column type can be text for all the columns.

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS reports ( [FromDate] TEXT, [ToDate] TEXT, [Result] TEXT	)")
	checkErr(err)
}

func sqltasker_a(fromDate, toDate string) string {

	db := openDB("./chinook.db")
	defer db.Close()

	result := "top 5 artists with the highest sum of UnitPrice produced: \r\n"
	rows, err := db.Query("select SUM(invoice_items.UnitPrice), artists.name from invoice_items,tracks,albums,artists where invoice_items.trackId = tracks.trackid AND tracks.albumid = albums.albumid and artists.artistid=albums.artistid AND tracks.ReleaseDate BETWEEN (?) AND (?) group by artists.artistid ORDER BY 1 DESC LIMIT 5", fromDate, toDate)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var moneymade float64
		var artistname string
		err = rows.Scan(&moneymade, &artistname)
		checkErr(err)
		result += artistname + " : " + fmt.Sprintf("%.0f", moneymade) + "$ \r\n"

	}
	err = rows.Err()
	checkErr(err)
	return result
}

func sqltasker_b(fromDate, toDate string) string {

	db := openDB("./chinook.db")
	defer db.Close()

	result := "top 5 genres with the highest sum of UnitPrice produced: \r\n"

	rows, err := db.Query("select SUM(invoice_items.UnitPrice), genres.name from invoice_items,tracks,genres where invoice_items.trackId = tracks.trackid AND tracks.genreid = genres.genreid AND tracks.ReleaseDate BETWEEN (?) AND (?) group by genres.genreID ORDER BY 1 DESC LIMIT 15;", fromDate, toDate)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var moneymade float64
		var genrename string
		err = rows.Scan(&moneymade, &genrename)
		checkErr(err)
		result += genrename + " : " + fmt.Sprintf("%.0f", moneymade) + "$ \r\n"

	}
	err = rows.Err()
	checkErr(err)
	return result
}

func sqltasker_c(fromDate, toDate string) string {

	db := openDB("./chinook.db")
	defer db.Close()

	result := "top 3 artists with the highest average of UnitPrice divided by the length of the produced tracks in seconds: \r\n"

	rows, err := db.Query("select (AVG(tracks.unitprice)/(sum(tracks.MILLISECONDS)))*1000, artists.name from artists,albums,tracks where artists.artistid = albums.artistid and albums.albumid = tracks.albumid AND tracks.ReleaseDate BETWEEN (?) AND (?) GROUP BY artists.artistid ORDER BY 1 LIMIT 5;", fromDate, toDate)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var eff float64
		var artistname string
		err = rows.Scan(&eff, &artistname)
		checkErr(err)
		result += artistname + " : " + fmt.Sprintf("%.8f", eff) + "$ \r\n"

	}
	err = rows.Err()
	checkErr(err)
	return result
}
