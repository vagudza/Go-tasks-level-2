package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSortUtil(t *testing.T) {
	err := createTestFiles()
	if err != nil {
		t.Fatalf(err.Error())
	}

	months := [12]string{"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}

	testTable := []struct {
		sc          SortConfig
		name        string
		out         string
		haveError   bool
		errorString string
	}{
		{
			name: "simple sort without parametres",
			sc: SortConfig{
				filename: "testing/sort3.txt",
			},
			out: "1\n11\n2\n4\n4\n5\n7\n8\n9",
		},
		{
			name: "simple reverse sort",
			sc: SortConfig{
				filename:    "testing/sort3.txt",
				reverseSort: true,
			},
			out: "9\n8\n7\n5\n4\n4\n2\n11\n1",
		},
		{
			name: "sort by colon k=2",
			sc: SortConfig{
				filename:   "testing/sort1.txt",
				sortColumn: 2,
			},
			out: "drwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\ndrwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/",
		},
		{
			name: "sort by colon (just reverse) k=2 -r",
			sc: SortConfig{
				filename:    "testing/sort1.txt",
				sortColumn:  2,
				reverseSort: true,
			},
			out: "drwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\ndrwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/",
		},
		{
			name: "sort by colon (as number) k=2 -n",
			sc: SortConfig{
				filename:           "testing/sort1.txt",
				sortByNumericValue: true,
				sortColumn:         2,
			},
			out: "drwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\ndrwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md",
		},
		{
			name: "sort by colon (unique) k=2 -u",
			sc: SortConfig{
				filename:   "testing/sort1.txt",
				sortColumn: 2,
				uniqueRows: true,
			},
			out: "drwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\ndrwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/",
		},
		{
			name: "sort by colon (as number, reverse) k=2 -n -r",
			sc: SortConfig{
				filename:           "testing/sort1.txt",
				sortByNumericValue: true,
				reverseSort:        true,
				sortColumn:         2,
			},
			out: "-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\ndrwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\ndrwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/",
		},
		{
			name: "sort by colon (as Month) k=6 -M",
			sc: SortConfig{
				months:      months,
				filename:    "testing/sort1.txt",
				sortByMonth: true,
				sortColumn:  6,
			},
			out: "drwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\ndrwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod",
		},
		{
			name: "sort by colon (as Month, reverse) k=6 -M -r",
			sc: SortConfig{
				months:      months,
				filename:    "testing/sort1.txt",
				sortByMonth: true,
				sortColumn:  6,
				reverseSort: true,
			},
			out: "-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod\ndrwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/\n-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md\ndrwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/\n-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md\n-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go\n-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum\ndrwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/",
		},
		{
			name: "sort by numeric value -n",
			sc: SortConfig{
				filename:           "testing/sort3.txt",
				sortByNumericValue: true,
			},
			out: "1\n2\n4\n4\n5\n7\n8\n9\n11",
		},
		{
			name: "sort by numeric value (unique) -n -u",
			sc: SortConfig{
				filename:           "testing/sort3.txt",
				sortByNumericValue: true,
				uniqueRows:         true,
			},
			out: "1\n2\n4\n5\n7\n8\n9\n11",
		},
		{
			name: "sort by numeric value (reverse) -n -r",
			sc: SortConfig{
				filename:           "testing/sort3.txt",
				sortByNumericValue: true,
				reverseSort:        true,
			},
			out: "11\n9\n8\n7\n5\n4\n4\n2\n1",
		},
		{
			name: "sort with unique rows -u",
			sc: SortConfig{
				filename:   "testing/sort2.txt",
				uniqueRows: true,
			},
			out: "LAPTOP\nRedHat\ncomputer\ndata\ndebian\nlaptop\nmouse",
		},
		{
			name: "sort with unique rows (reverse) -u -r",
			sc: SortConfig{
				filename:    "testing/sort2.txt",
				uniqueRows:  true,
				reverseSort: true,
			},
			out: "mouse\nlaptop\ndebian\ndata\ncomputer\nRedHat\nLAPTOP",
		},
		{
			name: "sort by Month -M",
			sc: SortConfig{
				months:      months,
				filename:    "testing/sort4.txt",
				sortByMonth: true,
			},
			out: "янв\nфев\nфев\nмар\nмар\nапр\nиюн\nиюл\nноя\nдек",
		},
		{
			name: "sort by Month (reverse) -M -r",
			sc: SortConfig{
				months:      months,
				filename:    "testing/sort4.txt",
				sortByMonth: true,
				reverseSort: true,
			},
			out: "дек\nноя\nиюл\nиюн\nапр\nмар\nмар\nфев\nфев\nянв",
		},
		{
			name: "sort by Month (unique) -M -u",
			sc: SortConfig{
				months:      months,
				filename:    "testing/sort4.txt",
				sortByMonth: true,
				uniqueRows:  true,
			},
			out: "янв\nфев\nмар\nапр\nиюн\nиюл\nноя\nдек",
		},
		{
			name: "check, if already sorted -c",
			sc: SortConfig{
				isRowsAlreadySorted: true,
				filename:            "testing/sort5.txt",
			},
			out: "true",
		},
		{
			name: "check, if already sorted -c",
			sc: SortConfig{
				isRowsAlreadySorted: true,
				filename:            "testing/sort2.txt",
			},
			out: "false",
		},
		{
			name: "error: file not found",
			sc: SortConfig{
				filename: "testing/sort13.txt",
			},
			haveError:   true,
			errorString: "can not read file 'testing/sort13.txt': open testing/sort13.txt: The system cannot find the file specified.",
		},
	}

	for _, testingCase := range testTable {
		t.Run(testingCase.name, func(t *testing.T) {
			result, err := Start(&testingCase.sc)
			if !testingCase.haveError {
				if err != nil {
					t.Errorf("expected err == nil; got '%s'", err.Error())
				}

				if result != testingCase.out {
					t.Errorf("expected result \n'%s';\n\ngot\n'%s'", testingCase.out, result)
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

	sortData1 := `drwxr-xr-x 1 vital 197121    0 янв  8 11:34 develop/
-rw-r--r-- 2 vital 197121  323 дек 18 16:42 go.mod
-rw-r--r-- 4 vital 197121 2311 фев 18 16:44 go.sum
drwxr-xr-x 5 vital 197121    0 ноя 10 17:12 listing/
-rw-r--r-- 5 vital 197121 3591 мар 11 11:05 main.go
drwxr-xr-x 7 vital 197121    0 май  8 11:34 pattern/
-rw-r--r-- 6 vital 197121 1349 апр  8 11:34 README.md
-rw-r--r-- 11 vital 197121 1349 июл  8 11:34 README.md`
	file, err := os.Create(testFilesFolder + "/sort1.txt")

	if err != nil {
		return fmt.Errorf("can not create file (1): '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(sortData1)

	sortData2 := `RedHat
mouse
data
LAPTOP
laptop
laptop
debian
data
computer
data`
	file, err = os.Create(testFilesFolder + "/sort2.txt")

	if err != nil {
		return fmt.Errorf("can not create file (2): '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(sortData2)

	sortData3 := `4
5
7
2
4
1
8
9
11`
	file, err = os.Create(testFilesFolder + "/sort3.txt")

	if err != nil {
		return fmt.Errorf("can not create file (3): '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(sortData3)

	sortData4 := `ноя
фев
янв
фев
апр
дек
мар
мар
июл
июн`
	file, err = os.Create(testFilesFolder + "/sort4.txt")

	if err != nil {
		return fmt.Errorf("can not create file (4): '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(sortData4)

	sortData5 := `1
2
3
4
5`
	file, err = os.Create(testFilesFolder + "/sort5.txt")

	if err != nil {
		return fmt.Errorf("can not create file (5): '%s'", err.Error())
	}
	defer file.Close()
	file.WriteString(sortData5)

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
