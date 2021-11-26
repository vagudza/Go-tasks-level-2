package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

OK go vet -c=10 task.go
OK golint task.go
OK go test -run ''  (coverage: 71.7%)

Запуск:
echo "asd1;asd2;asd3;asd4" | go run task.go -d=';' -f=1-3
echo "asd1;asd2;asd3;asd4" | go run task.go  -f=1 -s
*/

// Config - конфигурация программы
type Config struct {
	delim     string
	separated bool
	fields    string
}

// NewConfig - конструктор, парсящий флаги и аргументы
func NewConfig() *Config {
	conf := Config{}

	flagS := flag.Bool("s", false, "Print only-delimited rows")
	flag.StringVar(&conf.delim, "d", "", "Sets custom delimeter")
	flag.StringVar(&conf.fields, "f", "", "List of fields to cut")
	flag.Parse()

	conf.separated = *flagS
	return &conf
}

// Start - Точка входа в программу
func Start(conf *Config) {
	var str strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		str.WriteString(sc.Text())
	}

	result, err := cut(str.String(), conf)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Print(result)
}

func cut(row string, conf *Config) (string, error) {
	var result strings.Builder
	fields := make(map[int]bool)

	delim := "\t"
	if conf.delim != "" {
		if len(conf.delim) == 1 {
			delim = conf.delim
		} else {
			return "", fmt.Errorf("cut: the delimiter must be a single character")
		}
	}

	// разбор параметра fields, допустимые значения: 1 | 1-5 | 1, 4-6, 8
	if conf.fields != "" {
		sequence := strings.Split(conf.fields, ",")

		for _, seqPart := range sequence {
			seqPartRange := strings.Split(strings.TrimSpace(seqPart), "-")
			if len(seqPartRange) == 2 {
				seqPartRangeNumber1, err := strconv.Atoi(seqPartRange[0])
				if err != nil {
					return "", fmt.Errorf("cut: invalid field value: '%s'", seqPartRange[0])
				}

				seqPartRangeNumber2, err := strconv.Atoi(seqPartRange[1])
				if err != nil {
					return "", fmt.Errorf("cut: invalid field value: '%s'", seqPartRange[1])
				}

				if seqPartRangeNumber1 > seqPartRangeNumber2 {
					return "", fmt.Errorf("cut: invalid decreasing range")
				}

				if seqPartRangeNumber1 < 1 {
					return "", fmt.Errorf("cut: fields are numbered from 1")
				}

				for i := seqPartRangeNumber1; i <= seqPartRangeNumber2; i++ {
					fields[i] = true
				}
			} else {
				fieldNum, err := strconv.Atoi(strings.TrimSpace(seqPart))
				if err != nil {
					return "", fmt.Errorf("cut: invalid field value: '%s'", seqPart)
				}

				if fieldNum < 1 {
					return "", fmt.Errorf("cut: fields are numbered from 1")
				}

				fields[fieldNum] = true
			}
		}
	} else {
		return "", fmt.Errorf("cut: you must specify a list of bytes, characters, or fields")
	}

	splittedRow := strings.Split(row, delim)
	if conf.separated && len(splittedRow) == 1 {
		return "", nil
	}

	isNeedDelim := false
	for i, part := range splittedRow {
		_, ok := fields[i+1]
		if ok {
			if isNeedDelim {
				result.WriteString(delim + part)
			} else {
				result.WriteString(part)
				isNeedDelim = true
			}
		}
	}
	//fmt.Printf("%#v", splittedRow)
	return result.String(), nil
}

func main() {
	conf := NewConfig()
	Start(conf)
}
