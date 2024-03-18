package serde

import (
	"io"

	"github.com/apache/arrow/go/v16/arrow"
)

type Serializer interface {
    Serialize(io.Reader) (io.Reader, []string, *arrow.Schema)
}
