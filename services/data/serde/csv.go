package serde

import (
	"bufio"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v16/arrow"
)

type CsvSerializer struct {}

func (s *CsvSerializer) Serialize(r io.Reader) (io.Reader, []string, *arrow.Schema) {
    sc := bufio.NewScanner(r)
    if sc.Scan() {
        header := sc.Text()
        fmt.Println(header)
    }

    return r, []string{}, nil
}
