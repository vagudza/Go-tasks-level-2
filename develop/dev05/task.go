package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

OK go vet -c=10 task.go
OK golint task.go
OK go test -run ''  (coverage: 65.0%)

Запуск:
go run task.go -A=2 -i vital grep1.txt
go run task.go -c -i -v vital grep1.txt
*/

// Config - конфигурация программы
type Config struct {
	after       int
	before      int
	contextRows int
	count       bool
	ignoreCase  bool
	invert      bool
	fixed       bool
	strNum      bool
	regExp      string
	filename    string
}

// NewConfig - конструктор, парсящий флаги и аргументы
func NewConfig() *Config {
	conf := Config{}
	flag.IntVar(&conf.after, "A", 0, "Print +N rows after match")
	flag.IntVar(&conf.before, "B", 0, "Print +N rows before match")
	flag.IntVar(&conf.contextRows, "C", 0, "Print +N rows after and before match")
	flagC := flag.Bool("c", false, "Print count of match rows")
	flagI := flag.Bool("i", false, "Ignore case")
	flagV := flag.Bool("v", false, "Instead of a match, exclude")
	flagF := flag.Bool("F", false, "Exact match with a string, not a pattern")
	flagN := flag.Bool("n", false, "Print line number of match rows")

	flag.Parse()

	args := flag.Args()
	conf.count = *flagC
	conf.ignoreCase = *flagI
	conf.invert = *flagV
	conf.fixed = *flagF
	conf.strNum = *flagN

	if len(args) == 2 {
		conf.regExp = args[0]
		conf.filename = args[1]
	} else {
		log.Fatalf("The argument (path to the file name) must be one")
	}

	return &conf
}

// Start - Точка входа в программу
func Start(conf *Config) (interface{}, error) {
	rows, err := readFile(conf.filename)
	if err != nil {
		return "", fmt.Errorf("can not read file '%s': %s", conf.filename, err.Error())
	}

	return grep(rows, conf)
}

func grep(rows []string, conf *Config) (interface{}, error) {
	prefix := ""
	postfix := ""
	if conf.ignoreCase {
		prefix = "(?i)"
	}

	if conf.fixed {
		prefix += "^"
		postfix = "$"
	}

	re, err := regexp.Compile(prefix + conf.regExp + postfix)
	if err != nil {
		return "error", fmt.Errorf("invalid regular expression")
	}

	switch {
	case conf.after != 0:
		{
			for i, row := range rows {
				if re.MatchString(row) {
					if conf.after <= len(rows)-i {
						return rows[i : i+conf.after+1], nil
					}

					return rows[i:], nil
				}
			}
			return "not found", nil
		}

	case conf.before != 0:
		{
			for i, row := range rows {
				if re.MatchString(row) {
					if conf.before-1 <= i {
						return rows[i-conf.before : i+1], nil
					}
					return rows[:i+1], nil
				}
			}
			return "not found", nil
		}

	case conf.contextRows != 0:
		{
			for i, row := range rows {
				if re.MatchString(row) {
					startIndex := 0
					endIndex := len(rows)

					if conf.contextRows-1 <= i {
						startIndex = i - conf.contextRows
					}

					if conf.contextRows <= len(rows)-i {
						endIndex = i + conf.contextRows + 1
					}

					return rows[startIndex:endIndex], nil
				}
			}
			return "not found", nil
		}

	case conf.count:
		{
			totalCount := 0
			for _, row := range rows {
				totalCount += len(re.FindAllString(row, -1))
			}

			if conf.invert {
				return len(rows) - totalCount, nil
			}

			return totalCount, nil
		}

	case conf.strNum:
		{
			numberOfRows := []int{}
			for i, row := range rows {
				if re.MatchString(row) {
					numberOfRows = append(numberOfRows, i)
				}
			}
			return numberOfRows, nil
		}

	default:
		{
			result := []string{}
			for _, row := range rows {
				if re.MatchString(row) {
					result = append(result, row)
				}
			}
			return result, nil
		}
	}
}

func readFile(filename string) ([]string, error) {
	rows := []string{}
	file, err := os.Open(filename)
	if err != nil {
		return rows, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}
	return rows, nil
}

func main() {
	conf := NewConfig()
	res, err := Start(conf)
	if err != nil {
		log.Fatalf(err.Error())
	}

	switch result := res.(type) {
	case []string:
		{
			for _, row := range result {
				fmt.Println(row)
			}
		}
	case []int:
		{
			for _, row := range result {
				fmt.Printf("%d ", row)
			}
		}
	case int:
		fmt.Println(result)
	case string:
		fmt.Println(result)
	default:
		fmt.Printf("unknown type %T\n", res)
	}
}
