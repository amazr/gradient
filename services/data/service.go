package data

import (
	"example/hello/services/data/connectors"
	"example/hello/services/data/serde"
	"fmt"
	"io"
	"strings"
)

type DataService interface {
    List(did int) []connectors.FileLocator
    Read(fid string) ([]string, [][]any) 
    Write(did int, name string, size int64, content io.Reader) 
    Delete(fid string)
}

type DataResource struct {
    connector connectors.DataConnector
}

func NewSqlite() DataService {
    return &DataResource{
        connector: connectors.New(&serde.CsvSerializer{}),
    }
}

func (r *DataResource) Write(did int, name string, size int64, content io.Reader) {
    r.connector.Write(did, strings.Split(name, ".")[0], size, content)
}

func (r *DataResource) List(did int) []connectors.FileLocator {
    return r.connector.List(0)
}

func (r *DataResource) Read(fid string) ([]string, [][]any) {
    return r.connector.Read(fid)
}

func (r *DataResource) Delete(fid string) {
    fmt.Printf("Deleting file: %s\n", fid)
    r.connector.Delete(fid)
}
