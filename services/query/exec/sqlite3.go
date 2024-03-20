package exec

import (
	"database/sql"
	"encoding/json"
	"example/hello/services/query/lang"
	"fmt"
)

const db_file = "gradientfs.db"

type Sqlite3Executor struct {
    db *sql.DB
    root *lang.QueryRoot
    stmt string
}

func NewSqlite3Executor() *Sqlite3Executor {
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        panic(err)
    }

    return &Sqlite3Executor{
    	db: db,
        root: nil,
        stmt: "",
    }
}

func (ctx *Sqlite3Executor) Execute(r *lang.QueryRoot) ([]string, [][]any) {
    ctx.root = r
    r.Query.Accept(ctx.visit)

    fmt.Printf("Executing: %s\n", ctx.stmt)
    rows, err := ctx.db.Query(ctx.stmt)
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

func (ctx *Sqlite3Executor) visit(q lang.Query) {
    switch q.(type) {
    case lang.Filter:
        q.(lang.Filter).AcceptF(ctx.visitFilter)
    default:
        panic(fmt.Sprintf("Unknown query in exec: %v", q))
    }
}

func (ctx *Sqlite3Executor) visitFilter(f lang.Filter) {
    switch f.(type) {
    case *lang.KeepFilter:
        f.(*lang.KeepFilter).Condition.Accept(func (c lang.Condition) {
            switch c.(type) {
            case *lang.IsEqual:
                ie := c.(*lang.IsEqual)
                fmt.Println(ctx)
                ctx.stmt = fmt.Sprintf("SELECT * FROM '%s' WHERE \"%s\"=\"%s\"", ctx.root.Fid, ie.L, ie.R)
            default:
                panic(fmt.Sprintf("Unknown condition in exec: %v", c))
            }
        })
    case *lang.RemoveFilter:
        f.(*lang.RemoveFilter).Condition.Accept(func (c lang.Condition) {
            switch c.(type) {
            case *lang.IsEqual:
                ie := c.(*lang.IsEqual)
                ctx.stmt = fmt.Sprintf("SELECT * FROM '%s' WHERE \"%s\"!=\"%s\"", ctx.root.Fid, ie.L, ie.R)
            default:
                panic(fmt.Sprintf("Unknown condition in exec: %v", c))
            }
        })
    default:
        panic(fmt.Sprintf("Unknown filter in exec: %v", f))
    }
}
