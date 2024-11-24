package models

type Message struct {
	Type string `json:"type"` // Тип сообщения
	Data string `json:"data"` // Данные сообщения
}
