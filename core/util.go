package core

import (
	log "github.com/sirupsen/logrus"
	"io"
	"reflect"
	"strings"
)

func PStr(s string) *string {
	return &s
}

func PBool(b bool) *bool {
	return &b
}

func PInt(i int) *int {
	return &i
}

func ReadValue(value reflect.Value) any {
	if value.Kind() == reflect.Ptr && !value.Elem().IsValid() {
		return nil
	}
	typeStr := value.Type().String()
	switch typeStr {
	case "bool", "*bool":
		return reflect.Indirect(value).Bool()
	case "int", "*int":
		return reflect.Indirect(value).Int()
	case "string", "*string":
		return reflect.Indirect(value).String()
	default:
		log.Warn("Type not support: ", typeStr)
		return nil
	}
}

func UnCapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func ReadError(err error) *string {
	if err == nil {
		return nil
	}
	return PStr(err.Error())
}

func NoError(err error) bool {
	if err != nil {
		log.Error("Error occurred! ", err)
	}
	return err == nil
}

func Close(db io.Closer) {
	NoError(db.Close())
}
