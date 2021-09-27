package webhook

type WebHookMessageReceived struct {
	SchemaID  int32 `json:"ttn_schema_id"`
	MessageID int64 `json:"ttn_message_id"`
}

type ProcessSchema struct {
	SchemaID int32 `json:"schema_id"`
}
