package models

type Event struct {
	EventId string  `json:"event_id"` // 事件id
	Args    Request `json:"args"`     // 网关请求参数
}
