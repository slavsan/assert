package assert

import (
	"reflect"
	"testing"
)

type expectLevel1 struct {
	Not *expectLevel7
	To  *expectLevel2
	// ..
}

type expectLevel2 struct {
	Be         *expectLevel3
	Contain    *expectLevel4
	Have       *expectLevel5
	MatchError func(message string)
}

type expectLevel3 struct {
	Of      *expectLevel6
	Nil     func()
	True    func()
	False   func()
	EqualTo func(expected any)
}

type expectLevel4 struct {
	Substring func(sub string)
	Element   func(elem any)
}

type expectLevel5 struct {
	LengthOf func(length int)
	Property func(prop any)
}

type expectLevel6 struct {
	Type func(expected any)
}

type expectLevel7 struct {
	To *expectLevel8
}

type expectLevel8 struct {
	Be *expectLevel9
}

type expectLevel9 struct {
	Nil func()
}

type TestingInterface interface {
	Helper()
	Errorf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}

// Expect ..
func Expect(t TestingInterface) func(any) *expectLevel1 {
	t.Helper()
	return func(value any) *expectLevel1 {
		t.Helper()
		return &expectLevel1{
			Not: &expectLevel7{
				To: &expectLevel8{
					Be: &expectLevel9{
						Nil: func() {
							if isNil(value) {
								t.Errorf("expected '%v' to not be nil, but it is", value)
							}
						},
					},
				},
				// ..
			},
			To: &expectLevel2{
				Be: &expectLevel3{
					Of: &expectLevel6{
						Type: func(expected any) {
							// TODO: implement me
						},
					},
					Nil: func() {
						t.Helper()
						if !isNil(value) {
							t.Errorf("expected '%v' to be nil but it is not", value)
						}
					},
					True: func() {
						t.Helper()
						v, ok := value.(bool)
						if !ok {
							t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
							return
						}
						if v == false {
							t.Errorf("expected true but got false")
						}
					},
					False: func() {
						t.Helper()
						v, ok := value.(bool)
						if !ok {
							t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
							return
						}
						if v != false {
							t.Errorf("expected false but got true")
						}
					},
					EqualTo: func(expected any) {
						t.Helper()
						if !reflect.DeepEqual(expected, value) {
							expectedType := reflect.TypeOf(expected)
							actualType := reflect.TypeOf(value)
							if expectedType != actualType {
								t.Errorf("equality check failed\n\texpected: %#v (type: %s)\n\t  actual: %#v (type: %s)\n", expected, expectedType, value, actualType)
								return
							}
							t.Errorf("equality check failed\n\texpected: %#v\n\t  actual: %#v\n", expected, value)
						}
					},
					// ..
				},
				Have: &expectLevel5{
					// ..
					LengthOf: func(length int) {
						t.Helper()

						kind := reflect.TypeOf(value).Kind()

						if kind != reflect.Slice && kind != reflect.Array && kind != reflect.String && kind != reflect.Map {
							t.Errorf("expected target to be slice/array/map/string but it was %s", kind)
							return
						}

						if kind == reflect.String {
							reflectValue := reflect.ValueOf(value)
							if reflectValue.Len() != length {
								t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
							}
							return
						}

						reflectValue := reflect.ValueOf(value)
						if reflectValue.Len() != length {
							t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
						}
					},
					Property: func(prop any) {
						// TODO: implement me
					},
				},
				Contain: &expectLevel4{
					Substring: func(sub string) {
						// TODO: implement me
					},
					Element: func(elem any) {
						// TODO: implement me
					},
				},
				MatchError: func(message string) {
					// TODO: check if value is an error
					// TODO: implement me
				},
				// ..
			},
			// ..
		}
	}
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	valueOf := reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Chan, reflect.UnsafePointer, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		if valueOf.IsNil() {
			return true
		}
	}
	return false
}
