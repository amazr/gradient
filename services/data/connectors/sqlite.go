package connectors

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"example/hello/services/data/serde"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
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

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS file_refs (file_uuid primary key, name)")
    if err != nil {
        panic(err)
    }

    return &Sqlite3Connector{
    	db: db,
        ser: ser,
    }
}

func (s *Sqlite3Connector) Read(fid string) ([]string, [][]any) {
    query := fmt.Sprintf("SELECT * FROM \"%s\"", fid)
    rows, err := s.db.Query(query)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    cols, err := rows.Columns()
    if err != nil {
        panic(err)
    }

    var content [][]any
    for rows.Next() {
        rawv := make([]interface{}, len(cols))
        p := make([]interface{}, len(cols))
		for i := range rawv {
			p[i] = &rawv[i]
		}

        err := rows.Scan(p...)
        if err != nil {
            panic(err)
        }

        bytev, err := json.Marshal(rawv)
        if err != nil {
            panic(err)
        }

        val := []any{}
        err = json.Unmarshal(bytev, &val)
        if err != nil {
            panic(err)
        }

        content = append(content, val)
    }

    return cols, content
}

func (s *Sqlite3Connector) Write(did int, name string, size int64, reader io.Reader) {
    sc := bufio.NewScanner(reader)
    if !sc.Scan() {
        panic(sc.Err())
    }
    header := strings.TrimSpace(sc.Text())
    cols := strings.Split(header, ",")

    tableId := uuid.New()
    fmt.Printf("Adding new file: %s %s\n", name, tableId.String())
    _, err := s.db.Exec(fmt.Sprintf("CREATE TABLE \"%s\" (id integer primary key autoincrement)", tableId.String()))
    if err != nil {
        panic(err)
    }

    for _, col_name := range cols {
        noDoubleQuotes := strings.Replace(col_name, "\"", "", -1)
        add_col_sql := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN '%s'", tableId.String(), noDoubleQuotes)
        _, err := s.db.Exec(add_col_sql)
        if err != nil {
            panic(err)
        }
    }

    p1 := strings.Repeat("?,", len(cols)+1)
    placeholders := p1[:len(p1)-1]

    counter := 0
    for sc.Scan() {
        line := strings.TrimSpace(sc.Text())
        args := strings.Split(line, ",")

        var gargs []interface{}
        gargs = append(gargs, nil)
        for _, arg := range args {
            gargs = append(gargs, arg)
        }
        for i := 0; i < len(cols)+1-len(gargs); i++ {
            gargs = append(gargs, nil)
            fmt.Println("adding an extra nil")
        }

        insert_sql := fmt.Sprintf("INSERT INTO \"%s\" VALUES (%s)", tableId.String(), placeholders)
        _, err = s.db.Exec(insert_sql, gargs...)
        if err != nil {
            panic(err)
        }
        counter++
    }

    _, err = s.db.Exec("INSERT INTO file_refs (file_uuid, name) VALUES (?, ?)", tableId.String(), name)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Wrote %d lines to table\n", counter)
}

func (s *Sqlite3Connector) List(did int) []FileLocator {
    files, err := s.db.Query("SELECT file_uuid, name FROM file_refs")
    defer files.Close()
    if err != nil {
        panic(err)
    }

    locators := []FileLocator{}
    for files.Next() {
        file_uuid := *new(string)
        name := *new(string)
        err := files.Scan(&file_uuid, &name)
        if err != nil {
            fmt.Printf("Error: %v", err)
            continue
        }
        locators = append(locators, FileLocator{file_uuid, name, 0})
    }
    return locators
}

func (s *Sqlite3Connector) Delete(fid string) {
    _, err := s.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS \"%s\"", fid))
    if err != nil {
        panic(err)
    }

    _, err = s.db.Exec(fmt.Sprintf("DELETE FROM file_refs WHERE file_uuid==\"%s\"", fid))
    if err != nil {
        panic(err)
    }
}
