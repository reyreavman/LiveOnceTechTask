package models

type Status string

const (
	New        Status = "Новая"
	InProgress Status = "В работе"
	OnHold     Status = "На паузе"
	InReview   Status = "На проверке"
	Completed  Status = "Завершена"
	Cancelled  Status = "Отменена"
)
