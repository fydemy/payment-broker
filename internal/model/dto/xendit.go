package dto

type XenditResponse struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}
