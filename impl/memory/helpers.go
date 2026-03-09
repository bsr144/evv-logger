package memory

import "github.com/bsr144/evv-logger/entity"

func deepCopySchedule(src *entity.Schedule) *entity.Schedule {
	cp := *src
	if src.ClockInTime != nil {
		t := *src.ClockInTime
		cp.ClockInTime = &t
	}
	if src.ClockInLat != nil {
		v := *src.ClockInLat
		cp.ClockInLat = &v
	}
	if src.ClockInLng != nil {
		v := *src.ClockInLng
		cp.ClockInLng = &v
	}
	if src.ClockOutTime != nil {
		t := *src.ClockOutTime
		cp.ClockOutTime = &t
	}
	if src.ClockOutLat != nil {
		v := *src.ClockOutLat
		cp.ClockOutLat = &v
	}
	if src.ClockOutLng != nil {
		v := *src.ClockOutLng
		cp.ClockOutLng = &v
	}
	return &cp
}
