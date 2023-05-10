// Package declaration
package main

//Include code from other packages to use
//Package mysql provides a MySQL driver for Go's database/sql package - both
//log is for logging, Package log implements a simple logging package. It defines a type, Logger, with methods for formatting output.
//It also has a predefined 'standard' Logger accessible through helper functions Print[f|ln]
//net/http is for middleware
//Package template implements data-driven templates for generating textual output.
//To generate HTML output, see package html/template, which has the same interface as this package but automatically secures HTML output
import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// A struct is a type which contains named fields that have to be called throughout in the order they are declared
// Defined outside of functions means it can be used throughout
// Employee is the name of the type
type Employee struct {
	Id      int
	Hours   int
	Name    string
	Project string
}

// Functions are the building blocks of a Go program.
// They have inputs, outputs and a series of steps called statements which are executed in order
// := infers the variables
// db, err := sql.Open uses the MySQL package
// (db *sql.DB) is the functions parameter. name type. db is the name, sql.DB is the type and is the handler
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "admin"
	dbPass := "Too1Shuk!!!"
	dbName := "time"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(ab3.cluster-coad2yefgexn.us-east-1.rds.amazonaws.com:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Must is a helper that wraps a call to a function returning (*Template, error) and panics if the error is non-nil.
// ParseGlob creates a new Template and parses the template definitions from the files identified by the pattern.
// The files are matched according to the semantics of filepath.Match, and the pattern must match at least one file.
// The returned template will have the (base) name and (parsed) contents of the first file matched by the pattern.
// ParseGlob is equivalent to calling ParseFiles with the list of files matched by the pattern.
// When parsing multiple files with the same name in different directories, the last one mentioned will be the one that results.
var tmpl = template.Must(template.ParseGlob("form/*"))

// sub-function is called index, uses the standard middleware handler as a parameter
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id, hours int
		var name, project string
		err = selDB.Scan(&id, &hours, &name, &project)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Hours = hours
		emp.Name = name
		emp.Project = project
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, hours int
		var name, project string
		err = selDB.Scan(&id, &hours, &name, &project)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Hours = hours
		emp.Name = name
		emp.Project = project
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, hours int
		var name, project string
		err = selDB.Scan(&id, &hours, &name, &project)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Hours = hours
		emp.Name = name
		emp.Project = project
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		hours := r.FormValue("hours")
		name := r.FormValue("name")
		project := r.FormValue("project")
		insForm, err := db.Prepare("INSERT INTO employee(hours, name, project) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(hours, name, project)
		log.Println("INSERT: Hours: " + hours + " | Name: " + name + " | Project: " + project)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		hours := r.FormValue("hours")
		name := r.FormValue("name")
		project := r.FormValue("project")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE employee SET hours=?,name=?,project=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(hours, name, project, id)
		log.Println("UPDATE: Hours: " + hours + " | Name: " + name + " | Project: " + project)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
