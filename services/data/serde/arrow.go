package serde

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
    "strings"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/csv"
)

type ArrowSerializer struct {}

func (s *ArrowSerializer) Serialize(reader io.Reader) (io.Reader, []string, *arrow.Schema) {
    headerW := bytes.Buffer{}
    tee := io.TeeReader(reader, &headerW)
    sc := bufio.NewScanner(tee)
    if !sc.Scan() {
        panic(sc.Err())
    }

    header := strings.TrimSpace(sc.Text())
    columnNames := strings.Split(header, ",")

    multi := io.MultiReader(&headerW, reader)

    arrow_reader := csv.NewInferringReader(multi, csv.WithLazyQuotes(true))
    defer arrow_reader.Release()

    pReader, pWriter := io.Pipe()

    c := make(chan *arrow.Schema)
    go func(c chan *arrow.Schema) {
        defer pWriter.Close()
        sentSchema := true
        first := true
        for arrow_reader.Next() {
            if sentSchema {
                c <- arrow_reader.Schema()
                sentSchema = false
            }
            row, err := arrow_reader.Record().MarshalJSON()
            if err != nil {
                return
            }
            if first {
                fmt.Println(string(row))
                first = false
            }
            pWriter.Write(row)
        }
    }(c)
    return pReader, columnNames, <- c
}
