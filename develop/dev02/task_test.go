package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	testTable := []struct {
		name      string
		in        string
		out       string
		haveError bool
	}{
		{name: "simple unpack", in: "a4bc2d5e", out: "aaaabccddddde"},
		{name: "no unpack operations", in: "abcd", out: "abcd"},
		{name: "incorrect string", in: "45", out: "", haveError: true},
		{name: "empty string", in: "", out: ""},
		{name: "escape - последовательность 1", in: "qwe\\4\\5", out: "qwe45"},
		{name: "escape - последовательность 2", in: "qwe\\45", out: "qwe44444"},
		{name: "escape - последовательность 3", in: "qwe\\\\5", out: "qwe\\\\\\\\\\"},
	}

	for _, testingCase := range testTable {
		t.Run(testingCase.name, func(t *testing.T) {
			result, err := unpackString(testingCase.in)
			if !testingCase.haveError {
				if err != nil {
					t.Errorf("expected err == nil; got '%s'", err.Error())
				}

				if result != testingCase.out {
					t.Errorf("expected result '%s'; got '%s'", testingCase.out, result)
				}
			} else {
				if err != nil {
					if err.Error() != "некорректная строка" {
						t.Errorf("expected err.Error() == 'некорректная строка'; got '%s'", err.Error())
					}
				} else {
					t.Errorf("expected err.Error() == 'некорректная строка', but err is nil")
				}
			}
		})
	}
}
