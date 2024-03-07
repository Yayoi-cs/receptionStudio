package jsonHelper

type Project struct {
	ProjectName string
	ProjectID   string
	ProjectData ChildNode
}

type ChildNode struct {
	Name        string
	IsReception bool
	Child       []ChildNode
}
