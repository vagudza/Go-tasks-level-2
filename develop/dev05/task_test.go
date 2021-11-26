package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGrep(t *testing.T) {
	err := createTestFiles()
	if err != nil {
		t.Fatalf(err.Error())
	}

	testTable := []struct {
		conf        Config
		name        string
		out         interface{}
		haveError   bool
		errorString string
	}{
		{
			name: "simple grep without parametres",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "listing",
			},
			out: []string{"vital ноя 10 17:12 listing/", "vital ноя 17:12 listing_1.txt", "listing ноя 17:01 newfile.csv"},
		},
		{
			name: "grep: Print +N rows after match (-A=2)",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "README",
				after:    2,
			},
			out: []string{"vital апр 11:34 README.md", "vital дек 16:42 gopher1.go", "vital фев 16:44 gopher2.go"},
		},
		{
			name: "grep: Print +N rows after match (-A=100) - too much rows",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "develop.txt",
				after:    100,
			},
			out: []string{"ViTaL янв 11:34 develop.txt", "LiStinG июл 11:34 README.md", "Vital июл  8 11:34 README.md", "UsErNAME"},
		},
		{
			name: "grep: case ignore",
			conf: Config{
				filename:   "testing/grep1.txt",
				regExp:     "listing",
				ignoreCase: true,
			},
			out: []string{"vital ноя 10 17:12 listing/", "vital ноя 17:12 listing_1.txt", "listing ноя 17:01 newfile.csv", "Listing июл 11:34 README.md", "LiStinG июл 11:34 README.md"},
		},
		{
			name: "grep: Print +N rows before match (-B=2)",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "gopher",
				before:   2,
			},
			out: []string{"username", "vital апр 11:34 README.md", "vital дек 16:42 gopher1.go"},
		},
		{
			name: "grep: Print +N rows before match (-B=100) - too much rows",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "go.sum",
				before:   100,
			},
			out: []string{"vital июл 11:34 develop/", "vital дек 18 16:42 go.mod", "vital фев 18 16:44 go.sum"},
		},
		{
			name: "grep: Print +N rows after and before match (-C=1)",
			conf: Config{
				filename:    "testing/grep1.txt",
				regExp:      "gopher",
				contextRows: 1,
			},
			out: []string{"vital апр 11:34 README.md", "vital дек 16:42 gopher1.go", "vital фев 16:44 gopher2.go"},
		},
		{
			name: "grep: Print count of match rows (-c)",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "vital",
				count:    true,
			},
			out: 13,
		},
		{
			name: "grep: Print count of match rows, ignore case (-c -i)",
			conf: Config{
				filename:   "testing/grep1.txt",
				regExp:     "vital",
				count:      true,
				ignoreCase: true,
			},
			out: 15,
		},
		{
			name: "grep: Print count of match rows, ignore case (-c -i -v)",
			conf: Config{
				filename:   "testing/grep1.txt",
				regExp:     "vital",
				count:      true,
				ignoreCase: true,
				invert:     true,
			},
			out: 5,
		},
		{
			name: "grep: Exact match with a string, not a pattern (-F)",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "username",
				fixed:    true,
			},
			out: []string{"username"},
		},
		{
			name: "grep: Exact match with a string, not a pattern, case ignore (-F -i)",
			conf: Config{
				filename:   "testing/grep1.txt",
				regExp:     "username",
				fixed:      true,
				ignoreCase: true,
			},
			out: []string{"username", "UsErNAME"},
		},
		{
			name: "grep: Print line number of match rows (-n)",
			conf: Config{
				filename: "testing/grep1.txt",
				regExp:   "listing",
				strNum:   true,
			},
			out: []int{3, 10, 14},
		},
		{
			name: "grep: Print line number of match rows, ignore case (-n -i)",
			conf: Config{
				filename:   "testing/grep1.txt",
				regExp:     "listing",
				strNum:     true,
				ignoreCase: true,
			},
			out: []int{3, 10, 14, 15, 17},
		},
		{
			name: "grep: Print line number of match rows, ignore case (-n -i)",
			conf: Config{
				filename: "testing/grep13.txt",
			},
			haveError:   true,
			errorString: "can not read file 'testing/grep13.txt': open testing/grep13.txt: The system cannot find the file specified.",
		},
		{
			name: "grep: Print +N rows after match: with search regExp, which does not find (-A=1)",
			conf: Config{
				filename: "testing/grep1.txt",
				after:    1,
				regExp:   "regexp, which does not find result",
			},
			out: "not found",
		},
		{
			name: "grep: Print +N rows before match: with search regExp, which does not find (-B=1)",
			conf: Config{
				filename: "testing/grep1.txt",
				before:   1,
				regExp:   "regexp, which does not find result",
			},
			out: "not found",
		},
		{
			name: "grep: Print +N rows after and before match: with search regExp, which does not find (-C=1)",
			conf: Config{
				filename:    "testing/grep1.txt",
				contextRows: 1,
				regExp:      "regexp, which does not find result",
			},
			out: "not found",
		},
	}

	for _, testingCase := range testTable {
		t.Run(testingCase.name, func(t *testing.T) {
			res, err := Start(&testingCase.conf)

			if !testingCase.haveError {
				switch result := res.(type) {
				case []string:
					{
						rightResults := testingCase.out.([]string)
						if len(rightResults) == len(result) {
							for i, row := range result {
								if row != rightResults[i] {
									t.Errorf("expected result \n'%s';\n\ngot\n'%s'", rightResults[i], result)
								}
							}
						} else {
							t.Errorf("length of result and test answer is not equal")
						}
					}
				case []int:
					{
						rightResults := testingCase.out.([]int)
						if len(rightResults) == len(result) {
							for i, row := range result {
								if row != rightResults[i] {
									t.Errorf("expected result \n'%d';\n\ngot\n'%d'", rightResults[i], result)
								}
							}
						} else {
							t.Errorf("length of result and test answer is not equal")
						}
					}
				case int:
					{
						rightResult := testingCase.out.(int)
						if result != rightResult {
							t.Errorf("expected result \n'%d';\n\ngot\n'%d'", rightResult, result)
						}
					}
				case string:
					{
						rightResult := testingCase.out.(string)
						if result != rightResult {
							t.Errorf("expected result \n'%s';\n\ngot\n'%s'", rightResult, result)
						}
					}
				default:
					t.Errorf("expected result: int, []int, []string, but result have unknown type %T\n", res)
				}
			} else {
				if err != nil {
					if err.Error() != testingCase.errorString {
						t.Errorf("expected err.Error() == '%s'; got '%s'", testingCase.errorString, err.Error())
					}
				} else {
					t.Errorf("expected err != nil, but errs is nil")
				}
			}

		})
	}
}

func createTestFiles() error {
	testFilesFolder := "testing"
	isExists, err := exists(testFilesFolder)
	if err != nil {
		return err
	}

	if isExists {
		return nil
	}

	err = os.Mkdir(testFilesFolder, 0777)
	if err != nil {
		return fmt.Errorf("can not create folder for create testing files: '%s'", err.Error())
	}

	grepData1 := `vital июл 11:34 develop/
vital дек 18 16:42 go.mod
vital фев 18 16:44 go.sum
vital ноя 10 17:12 listing/
vital мар 11 11:05 main.go
vital май 11:34 pattern/
username
vital апр 11:34 README.md
vital дек 16:42 gopher1.go
vital фев 16:44 gopher2.go
vital ноя 17:12 listing_1.txt
vital мар 11:05 mainy.go
vital май 11:34 patterns/
vital апр 11:34 README.md
listing ноя 17:01 newfile.csv
Listing июл 11:34 README.md
ViTaL янв 11:34 develop.txt
LiStinG июл 11:34 README.md
Vital июл  8 11:34 README.md
UsErNAME`

	file, err := os.Create(testFilesFolder + "/grep1.txt")
	if err != nil {
		return fmt.Errorf("can not create file: '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(grepData1)

	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
