package main

import (
	"fmt"
	"log"
	"strings"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.

OK go vet -c=10 task.go
OK golint task.go
OK go test -run ''		(coverage: 83.3%)

% покрытия кода тестами в VS Code можно узнать так:
ПКМ --> палитра команд (ctrl+shift+p)
и ввести cove... и выбрать "Go toggle test coverage"

запуск из консоли теста: go test -run ''

Запуск (строка для распаковки захардкожена в main):
go run task.go
*/

func main() {
	str := "a4bc2d5e"
	unpacked, err := unpackString(str)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Распакованная строка: '%s'", unpacked)
}

func unpackString(str string) (string, error) {
	var unpackedStr strings.Builder
	var lastRune rune
	var isEscapeSymbol bool

	for _, symbol := range str {
		switch {
		case isEscapeSymbol:
			{
				isEscapeSymbol = false
				lastRune = symbol
			}
		case symbol <= '9' && symbol >= '0':
			{
				if lastRune != 0 {
					// поскольку в талблице ASCI/UTF-8 каждый символ имеет код, то для перевода из числовой руны в int
					// достаточно вычесть из кода этой руны код руны '0'
					iterationCount := int(symbol - '0')
					for i := 0; i < iterationCount; i++ {
						unpackedStr.WriteRune(lastRune)
					}
					lastRune = 0
				} else {
					return "", fmt.Errorf("некорректная строка")
				}
			}
		case symbol == '\\':
			{
				isEscapeSymbol = true
				if lastRune != 0 {
					unpackedStr.WriteRune(lastRune)
				}
			}
		default:
			if lastRune != 0 {
				unpackedStr.WriteRune(lastRune)
			}
			lastRune = symbol
		}
	}

	if lastRune != 0 {
		unpackedStr.WriteRune(lastRune)
	}
	return unpackedStr.String(), nil
}
