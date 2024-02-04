package assert

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

var (
	successful atomic.Int32
	failed     atomic.Int32
	skipped    atomic.Int32
	total      atomic.Int32
)

type Summary struct {
	Failed     int
	Successful int
	Skipped    int
	Total      int
}

const (
	bold    = "\033[1m"
	noBold  = "\033[0m"
	noColor = "\033[0m"
	gray    = "\033[1;30m"
	red     = "\033[0;31m"
	green   = "\033[0;32m"
	yellow  = "\033[0;33m"
	blue    = "\033[0;34m"
	purple  = "\033[0;35m"
	cyan    = "\033[0;36m"
)

func PrintAssertionsSummary(name string, colorsDisabled bool) {
	assertions := AssertionsSummary()
	message := fmt.Sprintf(
		"%s: %d passed, %d failed, %d total\n",
		name,
		assertions.Successful,
		assertions.Failed,
		assertions.Total,
	)
	count := len(message) - 1
	fmt.Printf(strings.Repeat("=", count) + "\n")
	if !colorsDisabled {
		failedColor := noColor
		if assertions.Failed > 0 {
			failedColor = red
		}
		message = fmt.Sprintf(
			"%s%s%s: %s%d passed%s, %s%d failed%s, %s%d total%s\n",
			bold, name, noColor,
			green, assertions.Successful, noColor,
			failedColor, assertions.Failed, noColor,
			noColor, assertions.Total, noColor,
		)
	}
	fmt.Printf(message)
	fmt.Printf(strings.Repeat("=", count) + "\n")
}

func AssertionsSummary() Summary {
	return Summary{
		Failed:     int(failed.Load()),
		Successful: int(successful.Load()),
		Skipped:    int(skipped.Load()),
		Total:      int(total.Load()),
	}
}

type expectLevel1 struct {
	Not *expectLevel7
	To  *expectLevel2
	// ..
}

type expectLevel2 struct {
	Be              *expectLevel3
	Contain         *expectLevel4
	Have            *expectLevel5
	MatchError      func(message string, t TestingInterface)
	MatchExactError func(message string, t TestingInterface)
}

type expectLevel3 struct {
	//Of      *expectLevel6
	Nil     func(t TestingInterface)
	True    func(t TestingInterface)
	False   func(t TestingInterface)
	EqualTo func(expected any, t TestingInterface)
}

type expectLevel4 struct {
	Substring func(sub string, t TestingInterface)
	//Element   func(elem any, t TestingInterface)
}

type expectLevel5 struct {
	LengthOf func(length int, t TestingInterface)
	//Property func(prop any)
}

//type expectLevel6 struct {
//	Type func(expected any)
//}

type expectLevel7 struct {
	To *expectLevel8
}

type expectLevel8 struct {
	Be *expectLevel9
}

type expectLevel9 struct {
	Nil func(t TestingInterface)
}

type TestingInterface interface {
	Helper()
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}

type ExpectInterface func(value any) *expectLevel1

func Expect(value any) *expectLevel1 {
	return expect(value, false)
}

// Require ..
func Require(value any) *expectLevel1 {
	return expect(value, true)
}

func expect(value any, require bool) *expectLevel1 {
	return &expectLevel1{
		Not: &expectLevel7{
			To: &expectLevel8{
				Be: &expectLevel9{
					Nil: func(t TestingInterface) {
						t.Helper()
						total.Add(1)
						if isNil(value) {
							failed.Add(1)
							if require {
								t.Fatalf("expected '%v' to not be nil, but it is", value)
							} else {
								t.Errorf("expected '%v' to not be nil, but it is", value)
							}
							return
						}
						successful.Add(1)
					},
				},
			},
			// ..
		},
		To: &expectLevel2{
			Be: &expectLevel3{
				//Of: &expectLevel6{
				//	Type: func(expected any) {
				//		// TODO: implement me
				//	},
				//},
				Nil: func(t TestingInterface) {
					total.Add(1)
					check(t)
					t.Helper()
					if !isNil(value) {
						t.Errorf("expected '%v' to be nil but it is not", value)
						failed.Add(1)
						return
					}
					successful.Add(1)
				},
				True: func(t TestingInterface) {
					t.Helper()
					total.Add(1)
					v, ok := value.(bool)
					if !ok {
						t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						failed.Add(1)
						return
					}
					if v == false {
						t.Errorf("expected true but got false")
						failed.Add(1)
						return
					}

					successful.Add(1)
				},
				False: func(t TestingInterface) {
					t.Helper()
					total.Add(1)
					v, ok := value.(bool)
					if !ok {
						t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						failed.Add(1)
						return
					}
					if v != false {
						t.Errorf("expected false but got true")
						failed.Add(1)
						return
					}
					successful.Add(1)
				},
				EqualTo: func(expected any, t TestingInterface) {
					t.Helper()
					total.Add(1)
					if !reflect.DeepEqual(expected, value) {
						expectedType := reflect.TypeOf(expected)
						actualType := reflect.TypeOf(value)
						if expectedType != actualType {
							t.Errorf("equality check failed\n\texpected: %#v (type: %s)\n\t  actual: %#v (type: %s)\n", expected, expectedType, value, actualType)
							failed.Add(1)
							return
						}
						t.Errorf("equality check failed\n\texpected: %#v\n\t  actual: %#v\n", expected, value)
						failed.Add(1)
						return
					}

					successful.Add(1)
				},
				// ..
			},
			Have: &expectLevel5{
				// ..
				LengthOf: func(length int, t TestingInterface) {
					t.Helper()
					total.Add(1)

					kind := reflect.TypeOf(value).Kind()

					if kind != reflect.Slice && kind != reflect.Array && kind != reflect.String && kind != reflect.Map {
						t.Errorf("expected target to be slice/array/map/string but it was %s", kind)
						failed.Add(1)
						return
					}

					if kind == reflect.String {
						reflectValue := reflect.ValueOf(value)
						if reflectValue.Len() != length {
							failed.Add(1)
							t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
							return
						}
						successful.Add(1)
						return
					}

					reflectValue := reflect.ValueOf(value)
					if reflectValue.Len() != length {
						t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
						failed.Add(1)
						return
					}

					successful.Add(1)
				},
				//Property: func(prop any) {
				//	// TODO: implement me
				//},
			},
			Contain: &expectLevel4{
				Substring: func(sub string, t TestingInterface) {
					t.Helper()
					total.Add(1)

					s, ok := value.(string)
					if !ok {
						t.Errorf("expected value to be a string but it is not")
						failed.Add(1)
						return
					}

					if !strings.Contains(s, sub) {
						t.Errorf(
							"expected string to contain substring but it does not\n\t   string: %s\n\tsubstring: %s\n",
							s, sub,
						)
						failed.Add(1)
						return
					}

					successful.Add(1)
				},
				//Element: func(elem any, t TestingInterface) {
				//	// TODO: implement me
				//},
			},
			MatchError: func(message string, t TestingInterface) {
				t.Helper()
				total.Add(1)

				if isNil(value) {
					t.Errorf("expected to match error but got nil value")
					failed.Add(1)
					return
				}

				err, ok := value.(error)
				if !ok {
					t.Errorf("expected to match error but value is not an error")
					failed.Add(1)
					return
				}

				if !strings.Contains(err.Error(), message) {
					t.Errorf(
						"expected error to contain message\n\t    actual error: %v (%s)\n\texpected message: %s\n",
						err, reflect.TypeOf(value), message,
					)
					failed.Add(1)
					return
				}

				successful.Add(1)
			},
			MatchExactError: func(message string, t TestingInterface) {
				t.Helper()
				total.Add(1)

				if isNil(value) {
					t.Errorf("expected to match error but got nil value")
					failed.Add(1)
					return
				}

				err, ok := value.(error)
				if !ok {
					t.Errorf("expected to match error but value is not an error")
					failed.Add(1)
					return
				}

				if err.Error() != message {
					t.Errorf(
						"expected error to contain message\n\t    actual error: %v (%s)\n\texpected message: %s\n",
						err, reflect.TypeOf(value), message,
					)
					failed.Add(1)
					return
				}

				successful.Add(1)
			},
		},
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

func check(t TestingInterface) {
	if tt, ok := t.(*testing.T); ok {
		if tt == nil {
			panic("nil *testing.T passed to expect call")
		}
	}
}
