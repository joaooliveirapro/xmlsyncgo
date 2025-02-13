package models

type Job struct {
	Content map[string]interface{} `json:"content"`
	Hash    string                 `json:"hash"`
	Edits   []Edit                 `json:"edits"`
	Deleted bool                   `json:"deleted"`
}

type Edit struct {
	Type            string `json:"type"`
	Ts              int64  `json:"ts"`
	RemoteFileModTs string `json:"remote_file_mod_ts"`
	Key             string `json:"key,omitempty"`
	Value           string `json:"value,omitempty"`
	NewValue        string `json:"new_value,omitempty"`
	OldValue        string `json:"old_value,omitempty"`
}
