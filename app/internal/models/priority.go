package models

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigth
)

func (p Priority) PriorityString() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigth:
		return "hight"
	default:
		return "unknown"
	}
}
