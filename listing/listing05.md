Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Программа выведет: error
В main определяется переменная err типа интерфейса error, который реализуется в customError. 
После вызова метода test() в err содержится интерфейс (Type=*customError, Value=nil). 
Аналогично примеру из listing3.md, функция возвращает кастомную ошибку, которая имеет свой указатель на тип.
Поэтому такое значение реализации интерфейса кастомной ошибки будет не равно nil, даже если значение указателя V внутри равно nil.

```
