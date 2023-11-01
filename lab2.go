package main

import (
	"fmt"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"bufio"
	"strings"
	"strconv"
)

var DB *sql.DB

type Row struct{
	id int 
	name string
	count int
}

func rowToJson(row Row) string {
	var json string = "[" + strconv.Itoa(row.id) + ", " + row.name + "]"
	return json
}

func rowToFullJson(row Row) string {
	var json string = "[" + strconv.Itoa(row.id) + ", " + row.name + ", " + strconv.Itoa(row.count) + "]"
	return json
}

func printHelp(){
	fmt.Println("Usage: ./lab2 [opts]")
	fmt.Println("	opts:")
	fmt.Println("	--start - Starts http server on 8080 port")
	fmt.Println("	--createdb - Creates new database and fills it from monitors file")
	fmt.Println("	--help - Print this help message")
}

func main() {
	if len(os.Args) > 1 {
		command := strings.ToLower(os.Args[1])
		if command == "--help" {
			printHelp()
		} else if command == "--createdb" {
			createDatabase()
			fillDatabase("monitors")
		} else if command == "--start" {
				http.HandleFunc("/category/monitors", getAll)
				http.HandleFunc("/category/monitor/", getById)
				err := http.ListenAndServe(":8080", nil)
				if err != nil {
					log.Fatal("Failed to start the server!\n", err)
				} else {
					fmt.Println("Server started!")
				}
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}
}

func getAllFromDB() []Row{
	openDatabase()
	rows, err := DB.Query("select id, name, count from monitors")
	if err != nil {
		log.Fatal("Error while queriyng data!\n", err)
		return []Row{}
	}
	var objRows []Row
	for rows.Next() {
		var row Row
		rows.Scan(&row.id, &row.name, &row.count)
		objRows = append(objRows, row)
	}
	DB.Close()
	return objRows
}

func getAll(w http.ResponseWriter, r *http.Request){
	openDatabase()

	key := r.URL.Query().Get("dev")

	var objRows = getAllFromDB()	
	var json = "{ \"monitors\": ["

	for i := range objRows[:len(objRows)-1]{
		if key == "true"{
			json += rowToFullJson(objRows[i])
		} else {
			json += rowToJson(objRows[i])
		}

		json += ", "
	}
	if key == "true"{
		json += rowToFullJson(objRows[len(objRows)-1])
	} else {
		json += rowToJson(objRows[len(objRows)-1])
	}
	json += "] }"
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, json)
}

func getById(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()

	if err == nil {
		monitorId, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/category/monitor/"))
		monitor := getMonitor(monitorId)
		var json string
		if monitor != nil {
			json = "{ \"monitor\": " + rowToJson(*monitor) + " }"
		} else {
			json = "{}"
		} 
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, json)
	} else {
		fmt.Println("Failed to parse request!\n", err)
	}
}

func getMonitor(id int) *Row {
	openDatabase()
	var row Row
	err := DB.QueryRow("select id, name from monitors where id = ?" , id).Scan(&row.id, &row.name)
	if err != nil {
		fmt.Println("Monitor not found!")
		return nil
	}
	updateMonitorCount(id)
	DB.Close()
	return &row 
}

func updateMonitorCount(id int){
	_, err := DB.Exec("update monitors set count = (select count+1 from monitors where id = $1) where id = $1", id)	
	if err != nil {
		fmt.Println("Monitor not found!")
		return 
	}
}

func openDatabase(){
	db, err := sql.Open("sqlite3", "lab2.db")
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func createDatabase(){
    fmt.Println("Creating database...")
	openDatabase()

	_, err := DB.Exec("create table if not exists monitors(id int, name varchar(255), count int); delete from monitors")

	if err != nil {
		fmt.Println("Failed to create database")
		log.Fatal(err)
	} else {
		fmt.Println("Database created")	
	}

	DB.Close()
}

func fillDatabase(filename string){
	var file *os.File
	var err error
	file, err = os.Open(filename)
	
	if err != nil {
		log.Fatal("Failed to open file!\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	openDatabase()

	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ",")
		id := arr[0]
		name := arr[1]
		_, err := DB.Exec("insert into monitors(id, name, count) values ($1, $2, 0)", id, name)

		if err != nil {
			log.Fatal("Failed to insert into database!\n", err)
		}
	}

	defer file.Close()
	DB.Close()
}
