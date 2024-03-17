package data

import (
	"example/hello/services/data/connectors"
	"io"
)

type DataService interface {
    List(did int) []connectors.FileLocator
    Read(fid int) []byte
    Write(did int, name string, size int64, content io.Reader) 
    Delete(fid int)
}

type DataResource struct {
    connector connectors.DataConnector
}

func NewSqlite() DataService {
    return &DataResource{
        connector: connectors.NewSqlite(),
    }
}

func (r *DataResource) Write(did int, name string, size int64, content io.Reader) {
    r.connector.Write(did, name, size, content)
}

func (r *DataResource) List(did int) []connectors.FileLocator {
    return r.connector.List(r.connector.GetRootDir())
}

func (r *DataResource) Read(fid int) []byte {
    return r.connector.Read(fid)
}

func (r *DataResource) Delete(fid int) {
    r.connector.Delete(fid)
}
