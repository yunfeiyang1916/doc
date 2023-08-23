package comm

import "time"

const TimeLayout = "2006-01-02 15:04:05"

func Now() *time.Time {
	now := time.Now()
	return &now
}

// 时间转字符串
func TimeFormat(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(TimeLayout)
}

// 字符串转时间
func TimeParse(str string) *time.Time {
	if t, err := time.Parse(TimeLayout, str); err != nil {
		return nil
	} else {
		return &t
	}
}
