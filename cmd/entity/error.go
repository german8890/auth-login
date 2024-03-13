package entity

type Error struct {
	Error  interface{} `json:"error,omitempty"`
	Detail string      `json:"detail,omitempty"`
}
