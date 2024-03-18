package connectors

import (
	"database/sql"
	"encoding/json"
	"example/hello/services/data/serde"
	"fmt"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

const db_file string = "gradientfs.db"

type Sqlite3Connector struct {
    db *sql.DB
    ser serde.Serializer
}

func New(ser serde.Serializer) *Sqlite3Connector {
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        panic(err)
    }

    return &Sqlite3Connector{
    	db: db,
        ser: ser,
    }
}

func (s *Sqlite3Connector) Read(fid int) []byte {
    row := s.db.QueryRow("SELECT column_names, content FROM files WHERE id=?", fid)
    cols := *new(string)
    content := *new([]byte)
    err := row.Scan(&cols, &content)
    if err != nil {
        panic(err)
    }
    fmt.Println(cols)
    return content
}

func (s *Sqlite3Connector) Write(did int, name string, size int64, reader io.Reader) {
    fmt.Println("writing")
    ser, cols, schema := s.ser.Serialize(reader)
    fmt.Printf("%v\n", schema)

    content, err := io.ReadAll(ser)
    if err != nil {
        panic(err)
    }

    mCols, err := json.Marshal(cols)
    if err != nil {
        panic(err)
    }

    _, err = s.db.Exec("INSERT INTO files (directory_id, name, size, column_names, content) VALUES (?, ?, ?, ?, ?)", did, name, size, mCols, content)
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

