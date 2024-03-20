package exec

import "example/hello/services/query/lang"

type Executor interface {
    Execute(r *lang.QueryRoot) ([]string, [][]any) 
}
