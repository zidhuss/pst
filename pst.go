package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zidhuss/pst/db"
)

const programURL = "pst.zidhuss.tech"

type app struct {
	db *db.PasteDatabase
}

func main() {
	pst := &app{db.CreatePasteDatabase("pst.db")}
	err := http.ListenAndServe("127.0.0.1:8081", handler(pst))
	if err != nil {
		log.Fatal(err)
	}
}

func handler(app *app) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", app.home).
		Methods("GET")

	r.HandleFunc("/", app.postPaste).
		Methods("POST")

	r.HandleFunc("/{pasteID}", app.retrievePaste)
	return r
}

func (pst *app) retrievePaste(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pasteID := vars["pasteID"]

	paste, err := pst.db.RetrievePaste(pasteID)
	if err != nil {
		// TODO: 404?
		log.Printf("%s\n", err)
		return
	}
	w.Write(paste.Data)
}

func (pst *app) postPaste(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024) // 2 Mb

	for i := 1; ; i++ {
		key := fmt.Sprintf("f:%d", i)
		contents := []byte(r.FormValue(key))

		if len(contents) == 0 {

			f, _, err := r.FormFile(key)
			if err != nil {
				break
			}

			contents, err = ioutil.ReadAll(f)
			if err != nil {
				break
			}

		}

		paste, err := pst.db.StorePaste(contents)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		fmt.Fprintf(w, "https://%s/%s\n", programURL, paste.ID)

		log.Printf("Storing paste %s from %s\n", paste.ID, ip)
	}
}

func (pst *app) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Help)
}

var Help = fmt.Sprintf(`
pst(1)                               PST                                  pst(1)

NAME

	pst: command line pastebin.


TL;DR

	~$ echo Hello world. | curl -F 'f:1=<-' %[1]s
	https://%[1]s/fpW


GET

	%[1]s/ID
		raw


POST

	%[1]s/

		f:N    contents or attached file.

	where N is a unique number within request. (This allows you to post
	multiple files at once.)

	returns: https://%[1]s/id for N in request


EXAMPLES

	Anonymous, unnamed paste, two ways:

		cat file.ext | curl -F 'f:1=<-' %[1]s
		curl -F 'f:1=@file.ext' %[1]s


SEE ALSO

	https://github.com/zidhuss/pst
`, programURL)
