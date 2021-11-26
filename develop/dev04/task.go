package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

OK go vet -c=10 task.go
OK golint task.go
OK go test -run ''  (coverage: 90.9%)

Запуск (слова для поиска анаграмм захардкожены в main()):
go run task.go
*/

// Dictionary - хранилище словаря в удобном для работы виде
type Dictionary struct {
	// мапа [кол-во_букв_в_слове]=слова_через_пробел
	dict map[int]string
}

// AddWords - функция добавления слов в словарь
func (d *Dictionary) AddWords(words []string) {
	for _, word := range words {
		d.dict[utf8.RuneCountInString(word)] += strings.ToLower(strings.TrimSpace(word)) + " "
	}
	//fmt.Printf("%#v", d.dict)
}

// NewDictionary - конструктор словаря
func NewDictionary() *Dictionary {
	return &Dictionary{
		dict: make(map[int]string),
	}
}

// Start - точка входа в программу: на вход - срез слов для поиска анаграмм,
// на выходе - карта[слово]["анаграмма_слова_1", ..., "анаграмма_слова_N"]
func Start(wordList []string, dict *Dictionary) map[string][]string {
	result := make(map[string][]string)
	for _, word := range wordList {
		wordAnagrams := anagrams(word, dict)
		if len(wordAnagrams) > 1 {
			sort.Strings(wordAnagrams)
			result[strings.ToLower(strings.TrimSpace(word))] = wordAnagrams
		}
	}
	return result
}

func anagrams(word string, dict *Dictionary) []string {
	validWord := strings.ToLower(strings.TrimSpace(word))
	validWordLen := utf8.RuneCountInString(validWord)

	// Формируем регулярное выражение вида [М,А,К,А,Р]{5}
	var regExpr strings.Builder
	regExpr.WriteString("[")

	for _, symbol := range validWord {
		regExpr.WriteString(string(symbol) + ",")
	}
	regExpr.WriteString("]{" + strconv.Itoa(validWordLen) + "}")

	// поиск анаграм с помощью регулярных выражений
	re := regexp.MustCompile(regExpr.String())
	//fmt.Printf("%q\n", re.FindAllString(dict.dict[validWordLen], -1))
	return re.FindAllString(dict.dict[validWordLen], -1)
}

func main() {
	myDict := NewDictionary()
	myDict.AddWords([]string{"АМКАР", "КАРМА", "КРАМА", "МАКАР", "МАКРА", "МАРКА", "РАМКА",
		"ПЯТАК", "ПЯТКА", "ТЯПКА", "КОСАЧ", "САЧОК", "ЧАСОК", "АВТОР", "ВАРТО", "ВТОРА", "ОТВАР",
		"РВОТА", "ТАВРО", "ТОВАР", "КАЧУР", "КРАУЧ", "КРУЧА", "КУРЧА", "РУЧКА", "ЧУРКА", "АБНЯ",
		"БАНЯ", "БАЯН", "КОРТ", "КРОТ", "ТРОК", "КОТ", "КТО", "ОТК", "ТОК",
	})

	//anagrams("кот", myDict)
	fmt.Println(Start([]string{"кот"}, myDict))
}
