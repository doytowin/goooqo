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
	typeStr := value.Type().String()
	log.Debug("Read value for type: ", typeStr)
	if value.Kind() == reflect.Ptr && !value.Elem().IsValid() {
		return nil
	}
	return reflect.Indirect(value).Interface()
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
