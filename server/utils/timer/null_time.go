package timer

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// NullTime 自定义时间类型，处理空值和零值
type NullTime struct {
	Time  time.Time
	Valid bool // Valid 为 true 表示有有效值
}

// UnmarshalJSON 实现JSON反序列化接口
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// 处理空字符串
	if s == "" {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	// 解析时间，支持多种格式
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

// MarshalJSON 实现JSON序列化接口
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time.Format("2006-01-02T15:04:05Z07:00"))
}

// UnmarshalText 处理form表单解析
func (nt *NullTime) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

// Scan 实现数据库扫描接口
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
		nt.Valid = true
	case []byte:
		return nt.UnmarshalText(v)
	case string:
		return nt.UnmarshalText([]byte(v))
	default:
		return errors.New(fmt.Sprintf("无法将 %T 转换为 NullTime", v))
	}
	return nil
}

// Value 实现数据库值接口
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
