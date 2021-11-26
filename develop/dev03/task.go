package main

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы							(игнорируются по умолчанию)
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов  	(?)

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

Отличия флага от аргумента:
	Флаг: 		-f
	Аргумент:	sort.txt

Сортировка по умолчанию: по возрастанию

OK go vet -c=10 task.go
OK golint task.go
OK go test -run ''  (coverage: 77.2%)

Запуск:
go run task.go -M -r -u sort1.txt
*/
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SortConfig - Конфигурация сортировки
type SortConfig struct {
	sortColumn          int
	sortByNumericValue  bool
	reverseSort         bool
	uniqueRows          bool
	sortByMonth         bool
	isRowsAlreadySorted bool
	months              [12]string
	filename            string
}

// NewSortConfig - Конструктор конфига
func NewSortConfig() *SortConfig {
	s := SortConfig{}
	s.months = [12]string{"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}
	flag.IntVar(&s.sortColumn, "k", 0, "Sets column for sort")
	flagN := flag.Bool("n", false, "Makes sort by numeric value")
	flagR := flag.Bool("r", false, "Makes reverse sort")
	flagU := flag.Bool("u", false, "Ignore duplicate lines")
	flagM := flag.Bool("M", false, "Makes sort by month")
	flagC := flag.Bool("c", false, "Check if rows already sorted")

	flag.Parse()

	args := flag.Args()
	s.sortByNumericValue = *flagN
	s.reverseSort = *flagR
	s.uniqueRows = *flagU
	s.sortByMonth = *flagM
	s.isRowsAlreadySorted = *flagC

	if len(args) == 1 {
		s.filename = args[0]
	} else {
		log.Fatalf("The argument (path to the file name) must be one")
	}

	return &s
}

// Start - Точка входа в программу сортировки
func Start(s *SortConfig) (string, error) {
	rows, err := readFile(s.filename)
	if err != nil {
		return "", fmt.Errorf("can not read file '%s': %s", s.filename, err.Error())
	}

	return sortRows(rows, s)
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

func uniqualizer(rows []string) []string {
	tempBuf := make(map[string]bool)
	for _, row := range rows {
		tempBuf[row] = true
	}

	// "очищаем" исходный слайс и заполняем его снова
	rows = rows[:0]
	for key := range tempBuf {
		rows = append(rows, key)
	}
	return rows
}

func getColumnValue(row string, s *SortConfig) (string, error) {
	// поиск разделителя: один или более пробельных символов
	re := regexp.MustCompile(`\s+`)
	// -b тег по умолчанию - при работе с колонками отрезаем пробелы вначале и в конце.
	listOfColumns := re.Split(strings.TrimSpace(row), -1)
	if len(listOfColumns) >= s.sortColumn {
		return listOfColumns[s.sortColumn-1], nil
	}
	return "", fmt.Errorf("can not find column")
}

func sortRows(rows []string, s *SortConfig) (string, error) {
	var sourceRows []string
	if s.isRowsAlreadySorted {
		sourceRows = make([]string, len(rows))
		_ = copy(sourceRows, rows)

		// принудительно игнорируем все остальные флаги, если есть флаг -c
		s.sortColumn = 0
		s.sortByMonth = false
		s.sortByNumericValue = false
		s.reverseSort = false
		s.uniqueRows = false
	}

	switch {
	case s.sortColumn > 0:
		{
			// уникализируем, если надо. Уникализация по колонке - не работает, т.к. по заданию - уникализация строки.
			if s.uniqueRows {
				rows = uniqualizer(rows)
			}

			sort.SliceStable(rows, func(i, j int) bool {
				// выбираем значение в выбранной колонке
				ith, err := getColumnValue(rows[i], s)
				if err != nil {
					return false
				}
				jth, err := getColumnValue(rows[j], s)
				if err != nil {
					return false
				}

				// сортировка по числам
				if s.sortByNumericValue {
					if s.reverseSort {
						return !ithNumLessThanJth(ith, jth)
					}
					return ithNumLessThanJth(ith, jth)
				}

				// сортировка по месяцам
				if s.sortByMonth {
					//fmt.Printf("%#v\n\n", s)
					if s.reverseSort {
						return !ithMonthLessThanJth(ith, jth, s)
					}
					return ithMonthLessThanJth(ith, jth, s)
				}

				if s.reverseSort {
					return ith < jth
				}
				return ith > jth
			})
		}
	case s.sortByNumericValue:
		{
			// уникализируем, если надо
			if s.uniqueRows {
				rows = uniqualizer(rows)
			}

			sort.SliceStable(rows, func(i, j int) bool {
				if s.reverseSort {
					return !ithNumLessThanJth(rows[i], rows[j])
				}
				return ithNumLessThanJth(rows[i], rows[j])
			})
		}
	case s.uniqueRows:
		{
			rows = uniqualizer(rows)
			sort.SliceStable(rows, func(i, j int) bool {
				// сортировка по месяцам
				if s.sortByMonth {
					if s.reverseSort {
						return !ithMonthLessThanJth(rows[i], rows[j], s)
					}
					return ithMonthLessThanJth(rows[i], rows[j], s)
				}

				if s.reverseSort {
					return rows[i] > rows[j]
				}
				return rows[i] < rows[j]
			})
		}
	case s.sortByMonth:
		{
			sort.SliceStable(rows, func(i, j int) bool {
				if s.reverseSort {
					return !ithMonthLessThanJth(rows[i], rows[j], s)
				}
				return ithMonthLessThanJth(rows[i], rows[j], s)
			})
		}
	default:
		sort.SliceStable(rows, func(i, j int) bool {
			if s.reverseSort {
				return rows[i] > rows[j]
			}
			return rows[i] < rows[j]
		})
	}

	if s.isRowsAlreadySorted {
		for i, row := range rows {
			if row != sourceRows[i] {
				return "false", nil
			}
		}
		return "true", nil
	}

	var result strings.Builder
	lenRows := len(rows)
	for i, row := range rows {
		if i < lenRows-1 {
			_, _ = result.WriteString(row + "\n")
		} else {
			_, _ = result.WriteString(row)
		}
	}

	return result.String(), nil
}

func ithNumLessThanJth(strI, strJ string) bool {
	ith, _ := strconv.Atoi(strI)
	jth, _ := strconv.Atoi(strJ)
	return ith < jth
}

func ithMonthLessThanJth(strI, strJ string, s *SortConfig) bool {
	for _, month := range s.months {
		switch {
		case month == strI:
			return true
		case month == strJ:
			return false
		default:
			continue
		}
	}
	return true
}

func main() {
	s := NewSortConfig()
	res, err := Start(s)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(res)
}
