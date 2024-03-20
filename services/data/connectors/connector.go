package connectors

import "io"

type FileLocator struct {
    Id string 
    Name string
    Size int
}

type DataConnector interface {
    Read(fid string) ([]string, [][]any)
    Write(did int, name string, size int64, content io.Reader) 
    List(did int) []FileLocator
    Delete(fid string)
}
