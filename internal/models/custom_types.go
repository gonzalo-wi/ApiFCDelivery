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

// Value implementa driver.Valuer para guardar en la base de datos
func (cd CustomDate) Value() (driver.Value, error) {
	return cd.Time, nil
}

// Scan implementa sql.Scanner para leer desde la base de datos
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
