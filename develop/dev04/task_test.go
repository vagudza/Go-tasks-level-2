package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	myDict := NewDictionary()
	myDict.AddWords([]string{"АМКАР", "КАРМА", "КРАМА", "МАКАР", "МАКРА", "МАРКА", "РАМКА",
		"ПЯТАК", "ПЯТКА", "ТЯПКА", "КОСАЧ", "САЧОК", "ЧАСОК", "АВТОР", "ВАРТО", "ВТОРА", "ОТВАР",
		"РВОТА", "ТАВРО", "ТОВАР", "КАЧУР", "КРАУЧ", "КРУЧА", "КУРЧА", "РУЧКА", "ЧУРКА", "АБНЯ",
		"БАНЯ", "БАЯН", "КОРТ", "КРОТ", "ТРОК", "КОТ", "КТО", "ОТК", "ТОК",
	})

	testTable := []struct {
		name string
		in   []string
		out  map[string][]string
	}{
		{name: "one key word", in: []string{"МАКАР"}, out: map[string][]string{"макар": {"амкар", "карма", "крама", "макар", "макра", "марка", "рамка"}}},
		{name: "several key words", in: []string{"кот", "Макар", "сачок"}, out: map[string][]string{
			"кот":   {"кот", "кто", "отк", "ток"},
			"макар": {"амкар", "карма", "крама", "макар", "макра", "марка", "рамка"},
			"сачок": {"косач", "сачок", "часок"},
		}},
		{name: "empty key words", in: []string{}, out: map[string][]string{}},
		{name: "key word not in dictionary", in: []string{"собака"}, out: map[string][]string{}},
		{name: "bad-symbols", in: []string{"]"}, out: map[string][]string{}},
	}

	for _, testingCase := range testTable {
		t.Run(testingCase.name, func(t *testing.T) {
			result := Start(testingCase.in, myDict)

			for word, anagrams := range testingCase.out {
				resAnagrams, ok := result[word]
				if !ok {
					t.Errorf("expected key word '%s' in result, but it have not", word)
				}

				for i, anagram := range anagrams {
					if anagram != resAnagrams[i] {
						t.Errorf("expected anagram '%s' in result, but it have not", anagram)
					}
				}
			}
		})
	}
}
