package lang

import "fmt"

type FilterVisitor func(Filter)
type Filter interface {
    FType() string
    AcceptF(FilterVisitor)
}

type KeepFilter struct {
    Condition Condition `json:"c"`
    QType string `json:"type"`
    SubType string `json:"ftype"`
}

type RemoveFilter struct {
    Condition Condition `json:"c"`
    QType string `json:"type"`
    SubType string `json:"ftype"`
}

func (kf *KeepFilter) Type() string {
    return kf.QType
}

func (rf *RemoveFilter) Type() string {
    return rf.QType
}

func (kf *KeepFilter) Accept(visitor Visitor) {
    visitor(kf)
}

func (rf *RemoveFilter) Accept(visitor Visitor) {
    visitor(rf)
}

func (kf *KeepFilter) AcceptF(visitor FilterVisitor) {
    visitor(kf)
}

func (rf *RemoveFilter) AcceptF(visitor FilterVisitor) {
    visitor(rf)
}

func (kf *KeepFilter) FType() string {
    return kf.SubType
}

func (rf *RemoveFilter) FType() string {
    return rf.SubType
}

func NewKeepFilter(condition Condition) (*KeepFilter) {
    return &KeepFilter{
    	Condition: condition,
        QType: "f",
        SubType: "k",
    }
}

func NewRemoveFilter(condition Condition) (*RemoveFilter) {
    return &RemoveFilter{
    	Condition: condition,
        QType: "f",
        SubType: "r",
    }
}

func DesKf(m map[string]interface{}) *KeepFilter {
    return NewKeepFilter(DesC(m["c"].(map[string]interface{})))
}

func DesRf(m map[string]interface{}) *RemoveFilter {
    return NewRemoveFilter(DesC(m["c"].(map[string]interface{})))
}

func DesF(m map[string]interface{}) Filter {
    switch m["ftype"] {
    case "k":
        return DesKf(m)
    case "r":
        return DesRf(m)
    default:
        panic(fmt.Sprintf("Unknown filter type: %v", m))
    }
}
