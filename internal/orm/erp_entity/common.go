package erp_entity

import (
	"app/internal/global"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
	"time"
	"unicode/utf8"
)

// UTF8String 兼容数据库字符串各种编码
type UTF8String string

func (us *UTF8String) isActuallyUTF8(raw []byte) bool {
	if !utf8.Valid(raw) {
		return false
	}
	// 把 raw 当 UTF-8 解码成 runes，再重新编码成 UTF-8
	reencoded := []byte(string(raw))
	// 只有真 UTF-8，原样 re-encode 才能完全还原
	return bytes.Equal(raw, reencoded)
}

func (us *UTF8String) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:

		// GBK编码
		if global.State.ErpEncoding == 0 {
			str, err := simplifiedchinese.GBK.NewDecoder().Bytes(v)
			if err != nil {
				return err
			}
			*us = UTF8String(str)
		} else if global.State.ErpEncoding == 1 {
			*us = UTF8String(v)
		} else {
			// 自动识别
			if utf8.ValidString(string(v)) && us.isActuallyUTF8(v) {
				*us = UTF8String(v)
			} else {
				str, err := simplifiedchinese.GBK.NewDecoder().Bytes(v)
				if err != nil {
					return err
				}
				*us = UTF8String(str)
			}
		}

		// 去除空格
		*us = UTF8String(strings.TrimSpace(string(*us)))

	case string:
		*us = UTF8String(strings.TrimSpace(v))
	case time.Time:
		*us = UTF8String(v.String())
	case int64, int, int32:
		*us = UTF8String(fmt.Sprintf("%d", v))
	case nil:
		*us = ""
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
func (us *UTF8String) MarshalBinary() (data []byte, err error) {
	return json.Marshal(us.String())
}
func (us *UTF8String) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, us)
}
func (us *UTF8String) String() string {
	return string(*us)
}
