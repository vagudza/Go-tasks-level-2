package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	testTable := []struct {
		name        string
		conf        Config
		in          string
		out         string
		haveError   bool
		errorString string
	}{
		{
			name: "simple cut",
			in: "apple1	juice2	so3	tasty4",
			conf: Config{
				fields: "1",
			},
			out: "apple1",
		},
		{
			name: "cut with delimeter",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "2",
				delim:  ";",
			},
			out: "juice2",
		},
		{
			name: "cut with delimeter several fields",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "2-4",
				delim:  ";",
			},
			out: "juice2;so3;tasty4",
		},
		{
			name: "cut with delimeter several fields 2",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "1,3-4",
				delim:  ";",
			},
			out: "apple1;so3;tasty4",
		},
		{
			name: "cut with delimeter several fields 3",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "1,2,4",
				delim:  ";",
			},
			out: "apple1;juice2;tasty4",
		},
		{
			name: "cut with delimeter overlimited fields",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "10-13",
				delim:  ";",
			},
			out: "",
		},
		{
			name: "cut with delimeter overlimited fields",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "10-13",
				delim:  ";",
			},
			out: "",
		},
		{
			name: "cut with decreasing range",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "3-1",
			},
			haveError:   true,
			errorString: "cut: invalid decreasing range",
			out:         "",
		},
		{
			name: "cut with field 0",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "0",
			},
			haveError:   true,
			errorString: "cut: fields are numbered from 1",
			out:         "",
		},
		{
			name: "cut with range from 0",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "0-3",
			},
			haveError:   true,
			errorString: "cut: fields are numbered from 1",
			out:         "",
		},
		{
			name: "cut with invalid delimeter",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "1",
				delim:  ";;",
			},
			haveError:   true,
			errorString: "cut: the delimiter must be a single character",
			out:         "",
		},
		{
			name: "cut with invalid field value",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "abc",
				delim:  ";",
			},
			haveError:   true,
			errorString: "cut: invalid field value: 'abc'",
			out:         "",
		},
		{
			name: "cut with invalid field range value",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "s-2",
				delim:  ";",
			},
			haveError:   true,
			errorString: "cut: invalid field value: 's'",
			out:         "",
		},
		{
			name: "cut with invalid field range value",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields: "1-s",
				delim:  ";",
			},
			haveError:   true,
			errorString: "cut: invalid field value: 's'",
			out:         "",
		},
		{
			name:        "cut with empty field flag",
			in:          "apple1;juice2;so3;tasty4",
			conf:        Config{},
			haveError:   true,
			errorString: "cut: you must specify a list of bytes, characters, or fields",
			out:         "",
		},
		{
			name: "cut with flag s",
			in:   "apple1;juice2;so3;tasty4",
			conf: Config{
				fields:    "3,4",
				separated: true,
			},
			out: "",
		},
	}

	for _, testingCase := range testTable {
		t.Run(testingCase.name, func(t *testing.T) {
			result, err := cut(testingCase.in, &testingCase.conf)
			if !testingCase.haveError {
				if err != nil {
					t.Errorf("expected err == nil; got '%s'", err.Error())
				}

				if result != testingCase.out {
					t.Errorf("expected result '%s'; got '%s'", testingCase.out, result)
				}
			} else {
				if err != nil {
					if err.Error() != testingCase.errorString {
						t.Errorf("expected err.Error() == '%s'; got '%s'", testingCase.errorString, err.Error())
					}
				} else {
					t.Errorf("expected err.Error() == '%s', but err is nil", testingCase.errorString)
				}
			}
		})
	}
}
