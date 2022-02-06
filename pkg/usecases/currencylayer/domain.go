package currencylayer

type Response struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int                `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}
