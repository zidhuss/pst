package db

import (
	"database/sql"
	"fmt"
	"log"
)

// Paste stores a pastes ID and its contents
type Paste struct {
	ID   string
	Data []byte
}

type PasteDatabase struct {
	db *sql.DB
}

func CreatePasteDatabase(name string) *PasteDatabase {
	pasteDB := &PasteDatabase{}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("opening db: %s\n", err)
	}

	db.Exec(`
		CREATE TABLE IF NOT EXISTS Paste (
			id INTEGER NOT NULL,
			data BLOB NOT NULL, PRIMARY KEY(id))
	`)

	// TODO: Check err
	pasteDB.db = db

	return pasteDB
}

// RetrievePaste queries the database for the paste and returns it.
func (pasteDB *PasteDatabase) RetrievePaste(id string) (*Paste, error) {
	paste := &Paste{ID: id}

	number, err := idToInt(id)
	if err != nil {
		return paste, err
	}

	stmt, err := pasteDB.db.Prepare("SELECT data FROM Paste WHERE id = $1")
	if err != nil {
		return paste, err
	}

	err = stmt.QueryRow(number).Scan(&paste.Data)

	return paste, err
}

// StorePaste stores the paste in the database.
func (pasteDB *PasteDatabase) StorePaste(data []byte) (*Paste, error) {
	n := pasteDB.nextID()
	paste := &Paste{ID: intToID(n), Data: data}

	stmt, err := pasteDB.db.Prepare("INSERT INTO Paste (id, data) VALUES ($1, $2)")
	if err != nil {
		return paste, fmt.Errorf("StorePaste preparing insert statement: %s", err)
	}
	_, err = stmt.Exec(n, data)
	if err != nil {
		return paste, fmt.Errorf("StorePaste executing insert: %s", err)
	}

	return paste, nil
}

func (pasteDB *PasteDatabase) nextID() uint64 {
	var i uint64

	// Reuse IDs
	// rows, _ := pasteDB.db.Query("SELECT id FROM DeletedId ORDER BY ASC")
	// if rows.Next() {
	// 	db.Exec()
	// 	rows.Scan(&i)
	// 	return intToID()
	// }

	rows, _ := pasteDB.db.Query("SELECT id FROM Paste ORDER BY id DESC")
	if rows.Next() {
		rows.Scan(&i)
		rows.Close()
		return i + 1
	}

	return 0
}
