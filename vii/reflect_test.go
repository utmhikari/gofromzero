package vii

import (
	"fmt"
	"reflect"
	"testing"
)

type UserInfo struct {
	Desc string `json:"desc"`
	City string `json:"city" desc:"the city of user"`
}

type User struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	Info *UserInfo `json:"info"`
}

func (u *User) ToString() string {
	return fmt.Sprintf("[%d]%s", u.ID, u.Name)
}

func (u *User) Response(from string, msg string) string {
	return fmt.Sprintf("User %s received msg %s from <%s>", u.ToString(), msg, from)
}

func TestReflect_Integer(t *testing.T) {
	i := 1
	v := reflect.ValueOf(i)
	vKind := v.Kind()
	vType := v.Type()
	t.Logf("i kind: %+v\n", vKind)
	t.Logf("i type: %+v\n", vType)

	itf := v.Interface()
	j, ok := itf.(int)
	t.Logf("j val: %+v\n", j)
	if !ok || j != i {
		t.Fatalf("i != j")
	}
}

func TestReflect_Pointer(t *testing.T) {
	u := &User{1, "jack", nil}
	vPtr := reflect.ValueOf(u)
	vPtrKind := vPtr.Kind()
	vPtrType := vPtr.Type()
	t.Logf("ptr kind: %+v\n", vPtrKind)
	t.Logf("ptr type: %+v\n", vPtrType)
	meth, ok := vPtrType.MethodByName("ToString")
	t.Logf("ptr meth ToString: %+v (%+v)\n", meth, ok)

	vVal := reflect.Indirect(vPtr)
	// vVal := vPtr.Elem()
	vValKind := vVal.Kind()
	vValType := vVal.Type()
	t.Logf("val kind: %+v\n", vValKind)
	t.Logf("val type: %+v\n", vValType)
	meth, ok = vValType.MethodByName("ToString")
	t.Logf("val meth ToString: %+v (%+v)\n", meth, ok)
}

func TestReflect_Method(t *testing.T) {
	u := &User{1, "jack", nil}
	uPtr := reflect.ValueOf(u)
	meth, ok := uPtr.Type().MethodByName("Response")
	if !ok {
		t.Fatalf("no method named Response")
	}
	t.Logf("meth Response: %+v\n", meth)

	methType := meth.Type
	if methType.NumIn() != 3 {
		t.Fatalf("invalid NumIn %d, expected %d", methType.NumIn(), 3)
	}
	if methType.NumOut() != 1 {
		t.Fatalf("invalid NumOut %d, expected %d", methType.NumOut(), 1)
	}
	from, msg := reflect.ValueOf("client"), reflect.ValueOf("ping")
	rets := meth.Func.Call([]reflect.Value{uPtr, from, msg})
	if len(rets) != 1 {
		t.Fatalf("invalid num rets %d, expected %d", len(rets), 1)
	}
	respVal := rets[0]
	if respVal.Type() != reflect.TypeOf("") {
		t.Fatalf("invalid ret type %v, expected %s", respVal.Type(), "STRING")
	}
	resp, ok := respVal.Interface().(string)
	if !ok {
		t.Fatalf("ret value cannot be converted to string")
	}
	t.Logf("resp: %s\n", resp)
}

func TestReflect_CopySliceAndMap(t *testing.T) {
	mp := map[string]int{
		"jack": 1,
		"tom":  2,
	}
	sl := []int{1, 1, 2, 3, 5, 8}
	vals := []reflect.Value{reflect.ValueOf(mp), reflect.ValueOf(sl)}
	var copyVals []reflect.Value
	for _, val := range vals {
		var copyVal reflect.Value
		switch val.Kind() {
		case reflect.Map:
			copyVal = reflect.MakeMap(val.Type())
			iter := val.MapRange()
			for iter.Next() {
				copyVal.SetMapIndex(iter.Key(), iter.Value())
			}
		case reflect.Slice:
			copyVal = reflect.AppendSlice(
				reflect.MakeSlice(val.Type(), 0, val.Len()),
				val)
		}
		copyVals = append(copyVals, copyVal)
	}

	for _, val := range copyVals {
		switch val.Kind() {
		case reflect.Map:
			if val.Len() != len(mp) {
				t.Fatalf("invalid map length %d, expected %d", val.Len(), len(mp))
			}
			copyVal, ok := val.Interface().(map[string]int)
			if !ok {
				t.Fatalf("map convert failed")
			}
			t.Logf("copied map: %+v", copyVal)
			for k, v := range mp {
				copyV, ok := copyVal[k]
				if !ok || !reflect.DeepEqual(v, copyV) {
					t.Fatalf("copy value of key %s failed, expected %d, actual %d", k, v, copyV)
				}
			}
		case reflect.Slice:
			if val.Len() != len(sl) {
				t.Fatalf("invalid slice length %d, expected %d", val.Len(), len(sl))
			}
			copyVal, ok := val.Interface().([]int)
			if !ok {
				t.Fatalf("slice convert failed")
			}
			t.Logf("copied slice: %+v", copyVal)
			if !reflect.DeepEqual(copyVal, sl) {
				t.Fatalf("slice not equal")
			}
		}
	}
}

func TestReflect_Struct(t *testing.T) {
	st := reflect.ValueOf(&User{
		ID:   9,
		Name: "Ronaldo",
		Info: &UserInfo{
			Desc: "SC",
			City: "Madrid",
		},
	}).Elem().Type()
	numField := st.NumField()
	t.Logf("num fields: %d", numField)
	for i := 0; i < numField; i++ {
		field := st.Field(i)
		t.Logf("field %d -> name: %s, type: %v, json: %s",
			i+1,
			field.Name,
			field.Type,
			field.Tag.Get("json"))
	}
	cityField := st.FieldByIndex([]int{2, 1})
	cityFieldDesc, ok := cityField.Tag.Lookup("desc")
	if !ok {
		t.Fatalf("cannot find city field desc")
	}
	t.Logf("CityField -> name: %s, type: %v, desc: %s",
		cityField.Name,
		cityField.Type,
		cityFieldDesc)
}
