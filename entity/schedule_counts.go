package entity

type ScheduleCounts struct {
	Total            int
	StatusCounts     map[string]int
	DateStatusCounts map[string]int
}
