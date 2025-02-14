package models

type Job struct {
	Id      int64
	Content map[string]interface{}
	Hash    string
	Edits   []Edit
	Deleted bool
}

type Edit struct {
	Id              int64
	Type            string
	Ts              int64
	RemoteFileModTs string
	Key             string
	Value           string
	NewValue        string
	OldValue        string
}
