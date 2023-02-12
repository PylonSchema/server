package gateway

type Message struct {
	Op int                    `json:"op"`
	D  map[string]interface{} `json:"d"`
}
