package models

import "fmt"

type Error struct {
	Type    string
	Message string
}

func ErrorToJSON(e Error) string {
	return fmt.Sprintf(`{"type:"%s", "message":"%s"}`, e.Type, e.Message)
}