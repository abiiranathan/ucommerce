// custom eclinic datatypes
package dbtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

const layout = "2006-01-02"

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

// Custom Json encoder
// Called when go types are being converted to json strings
func (date Date) MarshalJSON() ([]byte, error) {
	datetime := time.Time(date)

	if datetime.IsZero() {
		return []byte("null"), nil
	}

	dateBytes, err := datetime.MarshalJSON()
	if err != nil {
		return []byte(""), err
	}

	// Transform the date to format of layout in format yyyy-mm-dd
	dateString := string(dateBytes[1 : len(dateBytes)-1])
	dateString = fmt.Sprintf("\"%s-%02s-%02s\"", dateString[0:4], dateString[5:7], dateString[8:10])
	return []byte(dateString), nil
}

// Custom Json decoder
// Called to convert json strings to go types
func (date *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("date should be a string, got %v", data)
	}

	// by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	// If s is an empty string, assume a zero-value to allow for optional date.
	if strings.TrimSpace(s) == "" {
		s = "0001-01-01"
	}

	// Make sure that the user has provided the standard date format
	_, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("date should be of the format: yyyy-mm-dd")
	}

	// Convert date string to the standard format to RFC 3339 format
	s = fmt.Sprintf("\"%sT00:00:00Z\"", s)
	return (*time.Time)(date).UnmarshalJSON([]byte(s))
}

func (date Date) Year() int {
	return time.Time(date).Year()
}

func (date Date) Month() time.Month {
	return time.Time(date).Month()
}

func (date Date) Day() int {
	return time.Time(date).Day()
}

func (date Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
}
