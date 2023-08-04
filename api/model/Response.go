package model

type Error struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Param   interface{} `json:"param"`
	Type    string      `json:"type"`
}
