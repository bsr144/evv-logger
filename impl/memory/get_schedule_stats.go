package memory

import (
	"context"
	"time"

	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/evv/schedule"
)

func (s *Store) GetScheduleStats(_ context.Context) (*schedule.ScheduleStatsResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	today := time.Now().Format("2006-01-02")
	stats := &schedule.ScheduleStatsResult{Total: len(s.schedules)}

	for _, sched := range s.schedules {
		switch {
		case sched.Status == entity.ScheduleStatusMissed:
			stats.Missed++
		case sched.Status == entity.ScheduleStatusUpcoming && sched.Date == today:
			stats.Upcoming++
		case sched.Status == entity.ScheduleStatusCompleted && sched.Date == today:
			stats.Completed++
		}
	}

	return stats, nil
}
