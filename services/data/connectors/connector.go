package connectors

import "io"

type FileLocator struct {
    Id int
    Name string
    Size int
}

type DataConnector interface {
    Read(fid int) []byte 
    Write(did int, name string, size int64, content io.Reader) 
    List(did int) []FileLocator
    Delete(fid int)
    GetRootDir() int
}
