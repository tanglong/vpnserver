package redis

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

const (
	dollar byte = 36
	colon  byte = 58
	minus  byte = 45
	plus   byte = 43
	star   byte = 42
)

var delim []byte = []byte{13, 10}

// formatArgWithBuff formats the given argument to a Redis-styled argument byte slice.
func formatArgWithBuff(buff *bytes.Buffer, v interface{}) error {
	var bs []byte

	switch vt := v.(type) {
	case []byte:
		bs = vt
	case string:
		bs = []byte(vt)
	case int:
		bs = []byte(strconv.Itoa(vt))
	case int8:
		bs = []byte(strconv.FormatInt(int64(vt), 10))
	case int16:
		bs = []byte(strconv.FormatInt(int64(vt), 10))
	case int32:
		bs = []byte(strconv.FormatInt(int64(vt), 10))
	case int64:
		bs = []byte(strconv.FormatInt(vt, 10))
	case uint:
		bs = []byte(strconv.FormatUint(uint64(vt), 10))
	case uint8:
		bs = []byte(strconv.FormatUint(uint64(vt), 10))
	case uint16:
		bs = []byte(strconv.FormatUint(uint64(vt), 10))
	case uint32:
		bs = []byte(strconv.FormatUint(uint64(vt), 10))
	case uint64:
		bs = []byte(strconv.FormatUint(vt, 10))
	case bool:
		if vt {
			bs = []byte{49}
		} else {
			bs = []byte{48}
		}
	case nil:
		// empty byte slice
	default:
		// Fallback to reflect-based.
		switch reflect.TypeOf(vt).Kind() {
		case reflect.Slice:
			rv := reflect.ValueOf(vt)
			for i := 0; i < rv.Len(); i++ {
				if err := formatArgWithBuff(buff, rv.Index(i).Interface()); err != nil {
					return err
				}
			}
			return nil
		case reflect.Map:
			rv := reflect.ValueOf(vt)
			keys := rv.MapKeys()
			for _, k := range keys {
				if err := formatArgWithBuff(buff, k.Interface()); err != nil {
					return err
				}
				if err := formatArgWithBuff(buff, rv.MapIndex(k).Interface()); err != nil {
					return err
				}
			}
			return nil
		default:
			var buf bytes.Buffer

			fmt.Fprint(&buf, v)
			bs = buf.Bytes()
		}
	}

	if err := buff.WriteByte(dollar); err != nil {
		return err
	}
	if _, err := buff.WriteString(strconv.Itoa(len(bs))); err != nil {
		return err
	}
	if _, err := buff.Write(delim); err != nil {
		return err
	}
	if _, err := buff.Write(bs); err != nil {
		return err
	}
	if _, err := buff.Write(delim); err != nil {
		return err
	}
	return nil
}

// createRequest creates a Redis request for the given call and its arguments.
func createRequestWithBuff(buff *bytes.Buffer, calls ...call) error {
	for _, call := range calls {
		// Calculate number of arguments.
		argsLen := 1
		for _, arg := range call.args {
			switch arg.(type) {
			case []byte:
				argsLen++
			case nil:
				argsLen++
			default:
				// Fallback to reflect-based.
				kind := reflect.TypeOf(arg).Kind()
				switch kind {
				case reflect.Slice:
					argsLen += reflect.ValueOf(arg).Len()
				case reflect.Map:
					argsLen += reflect.ValueOf(arg).Len() * 2
				default:
					argsLen++
				}
			}
		}

		// number of arguments
		if err := buff.WriteByte(star); err != nil {
			return err
		}
		if _, err := buff.WriteString(strconv.Itoa(argsLen)); err != nil {
			return err
		}
		if _, err := buff.Write(delim); err != nil {
			return err
		}

		// command
		if err := buff.WriteByte(dollar); err != nil {
			return err
		}
		if _, err := buff.WriteString(strconv.Itoa(len(call.cmd))); err != nil {
			return err
		}
		if _, err := buff.Write(delim); err != nil {
			return err
		}
		if _, err := buff.WriteString(string(call.cmd)); err != nil {
			return err
		}
		if _, err := buff.Write(delim); err != nil {
			return err
		}

		// arguments
		for _, arg := range call.args {
			if err := formatArgWithBuff(buff, arg); err != nil {
				return err
			}
		}
	}
	return nil
}
