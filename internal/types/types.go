package types

import "time"

const (
	timeFormart = "2006-01-02 15:04:05"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var _t time.Time
	var err error

	_t, err = time.ParseInLocation(`"`+time.RFC3339+`"`, string(data), time.Local)
	if err == nil {
		goto end
	}

	_t, err = time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
end:
	*t = Time(_t)

	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')

	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
