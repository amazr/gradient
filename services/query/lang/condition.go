package lang

import "fmt"

type ConditionVisitor func(Condition)
type Condition interface {
    Type() string
    Accept(ConditionVisitor)
}

type IsEqual struct {
    L any `json:"l"`
    R any `json:"r"`
    CType string `json:"type"`
}

func (ie *IsEqual) Type() string {
    return ie.CType
}

func (ie *IsEqual) Accept(visitor ConditionVisitor) {
    visitor(ie)
}

func NewIsEqual(l any, r any) (*IsEqual) {
    return &IsEqual{
        L: l,
        R: r,
        CType: "ie",
    }
}

func DesIe(m map[string]interface{}) *IsEqual {
    return NewIsEqual(m["l"], m["r"])
}

func DesC(m map[string]interface{}) Condition {
    switch m["type"] {
    case "ie":
        return DesIe(m)
    default:
        panic(fmt.Sprintf("Unknown condition type: %v", m))
    }
}
