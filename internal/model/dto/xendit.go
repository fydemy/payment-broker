package dto

type XenditResponse struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

type XenditNonEventWebhook struct {
	ExternalID string `json:"external_id"`
}

type XenditEventWebhook struct {
	ReferenceID string `json:"reference_id"`
}
