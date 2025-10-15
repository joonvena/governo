package auditlog

// Simplified Cloud Audit Log entry
type AuditLogEntry struct {
	ProtoPayload struct {
		MethodName   string         `json:"methodName"`
		ServiceName  string         `json:"serviceName"`
		ResourceName string         `json:"resourceName"`
		Request      map[string]any `json:"request"`
		Response     map[string]any `json:"response"`
	} `json:"protoPayload"`
}
