package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time)
}

func (cd CustomDate) Value() (driver.Value, error) {
	return cd.Time, nil
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	if t, ok := value.(time.Time); ok {
		cd.Time = t
		return nil
	}

	return nil
}
