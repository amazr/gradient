package lang

import (
	"encoding/json"
	"fmt"
)

type Visitor func(Query)
type Query interface {
    Type() string
    Accept(Visitor)
}

type QueryRoot struct {
    Fid string `json:"fid"`
    Query Query `json:"query"`
}

func New(fid string, query Query) (*QueryRoot) {
    return &QueryRoot{
    	Fid:   fid,
    	Query: query,
    }
}

func Ser(r *QueryRoot) []byte {
    b, e := json.Marshal(r)
    if e != nil {
        panic(e)
    }
    return b
}

func Des(b []byte) *QueryRoot {
    m := map[string]interface{}{}
    e := json.Unmarshal(b, &m)
    if e != nil {
        panic(e)
    }
    return New(m["fid"].(string), DesQ(m["query"].(map[string]interface{})))
}

func DesQ(m map[string]interface{}) Query {
    switch m["type"].(string) {
    case "f":
       return DesF(m).(Query) 
    default:
        panic(fmt.Sprintf("Unknown query type: %v", m))
    }
}
