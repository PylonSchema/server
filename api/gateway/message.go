package gateway

type Message struct {
	Op int         `json:"op"`
	D  interface{} `json:"d"`
}
