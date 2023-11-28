package util

func PStr(s string) *string {
	return &s
}

func PBool(b bool) *bool {
	return &b
}

func PInt(i int) *int {
	return &i
}
