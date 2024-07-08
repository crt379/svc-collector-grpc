package util

import (
	"net"
	"reflect"
	"unicode"
)

// 获取本机ip
func GetIP() string {
	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		return ""
	}

	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)

	return addr.IP.String()
}

// 根据 newvalueobj 修改 obj 和 newvalueobj 相同字段名相同类型的值为newvalueobj的, 如果 newvalueobj 某个字段为 Zero 则忽略该字段
func UpdateValue(newvalueobj any, obj any) (isupdate bool) {
	var k1 reflect.Kind
	var k2 reflect.Kind

	t1 := reflect.TypeOf(newvalueobj)
	v1 := reflect.ValueOf(newvalueobj)
	k1 = t1.Kind()
	if k1 == reflect.Ptr {
		t1 = t1.Elem()
		v1 = v1.Elem()
	}
	if t1.Kind() != reflect.Struct {
		return isupdate
	}

	t2 := reflect.TypeOf(obj)
	v2 := reflect.ValueOf(obj)
	k2 = t2.Kind()
	if k2 != reflect.Ptr {
		return isupdate
	}

	// 同一指针
	if k1 == k2 && reflect.DeepEqual(newvalueobj, obj) {
		return isupdate
	}

	t2 = t2.Elem()
	v2 = v2.Elem()
	if t2.Kind() != reflect.Struct {
		return isupdate
	}

	fieldnum1 := t1.NumField()
	for i := 0; i < fieldnum1; i++ {
		v2v := v2.FieldByName(t1.Field(i).Name)
		if v2v.IsValid() && v2v.CanSet() {
			v1v := v1.Field(i)
			if !v1v.IsZero() && v1v.Type() == v2v.Type() && !v1v.Equal(v2v) {
				v2v.Set(v1v)
				isupdate = true
			}
		}
	}
	return isupdate
}

// 根据 newobj 修改 obj 和 newobj 相同字段名相同类型的值为newobj的, 要求obj 和 newobj 类型要一致, 如果 newobj 某个字段为 Zero 则忽略该字段
func UpdateValueSame[T any](newobj T, obj T) (isupdate bool) {
	t1 := reflect.TypeOf(newobj)
	v1 := reflect.ValueOf(newobj)
	if t1.Kind() != reflect.Ptr {
		return false
	}

	t1 = t1.Elem()
	v1 = v1.Elem()
	if t1.Kind() != reflect.Struct {
		return false
	}

	// 同一指针
	if reflect.DeepEqual(newobj, obj) {
		return false
	}

	v2 := reflect.ValueOf(obj).Elem()
	for i := 0; i < t1.NumField(); i++ {
		v2v := v2.Field(i)
		if v2v.IsValid() && v2v.CanSet() {
			v1v := v1.Field(i)
			if !v1v.IsZero() && !v1v.Equal(v2v) {
				v2v.Set(v1v)
				isupdate = true
			}
		}
	}

	return
}

// 判断 s1 是不是 s2 的前缀
func StrInStrBegin(s1, s2 string) bool {
	if len(s2) < len(s1) {
		return false
	}

	if s2[:len(s1)] == s1 {
		return true
	}

	return false
}

// 判断s中的符号是不是 legitimate, 如果不是则返回 false
func StrPunctIllegal(s string, legitimate rune) bool {
	for _, c := range s {
		if unicode.IsPunct(c) && c != legitimate {
			return true
		}
	}

	return false
}
