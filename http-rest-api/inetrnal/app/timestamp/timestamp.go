package timestamp

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	//"github.com/globalsign/mgo/bson"
)

// TimestampFormatLayout formatlayout, save in mysql
const TimestampFormatLayout = "2006-01-02 15:04:05"

type Timestamp struct {
	time.Time
}

// New return timestamp
func New() Timestamp {
	return Timestamp{}
}

// NewWithTime return a specified time's Timestamp
func NewWithTime(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// Now returns the current local time.
func Now() Timestamp {
	return Timestamp{Time: time.Now()}
}

// Unix returns the local Time corresponding to the given Unix time,
func Unix(sec int64, nsec int64) Timestamp {
	return Timestamp{time.Unix(sec, nsec)}
}

// Date return time with date
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Timestamp {
	return Timestamp{time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// AutoSet auto set type
func AutoSet(data interface{}, fields ...string) {
	if m, ok := data.(map[string]interface{}); ok {
		for _, field := range fields {
			v, exists := m[field]
			if !exists {
				continue
			}
			if t, ok := v.(float64); ok {
				m[field] = Unix(int64(t), 0)
			}
			if t, ok := v.(int64); ok {
				m[field] = Unix(t, 0)
			}
		}
	}
}

func (t Timestamp) Pointer() *Timestamp {
	return &t
}

// MarshalJSON implements the json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {

	if t.IsZero() || t.Time.UnixNano() == time.Unix(0, 0).UnixNano() {
		return []byte("null"), nil
	}

	ts := t.Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var (
		err error
		str = strings.Replace(string(b), "\"", "", -1)
	)

	if str == "null" || str == "0" {
		return nil
	}

	t.Time, err = dateparse.ParseAny(str)
	return err
}

func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(time.Time); ok {
		*t = Timestamp{v}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", value)
}

// Value insert timestamp into mysql need this function.
func (t Timestamp) Value() (driver.Value, error) {
	if t.IsZero() || t.Time.UnixNano() == time.Unix(0, 0).UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
