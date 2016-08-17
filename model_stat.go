package ui

// Stat describes a statistic/metric
type Stat struct {
	Metric string `json:"metric"`
	Value  int64  `json:"value"`
}
