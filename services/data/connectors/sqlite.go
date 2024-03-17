package connectors

import (
	"database/sql"
	"fmt"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

const db_file string = "gradientfs.db"

type Sqlite3Connector struct {
    db *sql.DB
}

func NewSqlite() *Sqlite3Connector {
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        panic(err)
    }

    return &Sqlite3Connector{
    	db: db,
    }
}

func (s *Sqlite3Connector) Read(fid int) []byte {
    row := s.db.QueryRow("SELECT content FROM files WHERE id=?", fid)
    content := *new([]byte)
    err := row.Scan(&content)
    if err != nil {
        panic(err)
    }
    return content
}

func (s *Sqlite3Connector) Write(did int, name string, size int64, reader io.Reader) {
    content, err := io.ReadAll(reader)
    if err != nil {
        panic(err)
    }

    _, err = s.db.Exec("INSERT INTO files (directory_id, name, size, content) VALUES (?, ?, ?, ?)", did, name, size, content)
    if err != nil {
        panic(err)
    }
}

func (s *Sqlite3Connector) List(did int) []FileLocator {
    rows, err := s.db.Query("SELECT id, name, size FROM files WHERE directory_id=?", did)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    locators := []FileLocator{}
    for rows.Next() {
        id := *new(int)
        name := *new(string)
        size := *new(int)
        err := rows.Scan(&id, &name, &size)
        if err != nil {
            fmt.Printf("Error: %v", err)
            continue
        }
        locators = append(locators, FileLocator{id, name, size})
    }

    return locators
}

func (s *Sqlite3Connector) Delete(fid int) {
    _, err := s.db.Exec("DELETE FROM files WHERE id=?", fid)
    if err != nil {
        panic(err)
    }
}

func (s *Sqlite3Connector) GetRootDir() int {
    row := s.db.QueryRow("SELECT id FROM directories WHERE parent_id IS NULL")
    did := *new(int)
    err := row.Scan(&did)
    if err != nil {
        panic(err)
    }
    return did
}

