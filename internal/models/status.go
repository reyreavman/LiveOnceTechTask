package models

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	New        Status = "Новая"
	InProgress Status = "В работе"
	OnHold     Status = "На паузе"
	Completed  Status = "Завершена"
	Cancelled  Status = "Отменена"
)

func (s Status) IsValid() bool {
	switch s {
	case New, InProgress, OnHold, Completed, Cancelled:
		return true
	default:
		return false
	}
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	status := Status(raw)
	if !status.IsValid() {
		return fmt.Errorf("invalid status: %s", status)
	}

	*s = status

	return nil
}
