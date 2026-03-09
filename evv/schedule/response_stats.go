package schedule

type StatsResponse struct {
	Total     int `json:"total"`
	Missed    int `json:"missed"`
	Upcoming  int `json:"upcoming"`
	Completed int `json:"completed"`
}
