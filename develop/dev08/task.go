package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.

pwd - утилита, которая позволяет вывести в терминал путь к текущей папке.
ps - перечень всех процессов всех пользователей.
Netcat — утилита Unix, позволяющая устанавливать соединения TCP и UDP, принимать оттуда данные и передавать их.
fork - команда запускающая процесс потомок.
exec - команда, запускающая внешний исполняемый файл/скрипт
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	quit = "\\quit"
	cd   = "cd"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			_, err = fmt.Fprintln(os.Stderr, err)
			if err != nil {
				return
			}
		}

		err = runCommand(cmdString)
		if err != nil {
			_, err = fmt.Fprintln(os.Stderr, err)
			if err != nil {
				return
			}
		}
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSpace(commandStr)
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)

	if len(arrCommandStr) == 0 {
		return nil
	}

	switch arrCommandStr[0] {
	case cd:
		if len(arrCommandStr) < 2 {
			return nil
		}

		err := os.Chdir(arrCommandStr[1])
		if err != nil {
			return err
		}
		return nil
	case quit:
		os.Exit(0)
	}

	cmd := exec.Command("bash", "-c", commandStr)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
