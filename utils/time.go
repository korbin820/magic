package util

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type DateTime time.Time

func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *DateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = DateTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (c *DateTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = DateTime(t) //set result using the pointer
	return nil
}

func (c *DateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*c)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02"))), nil
}
