package internal

import (
	"bufio"
	"os"
	"os/exec"
	"time"
)

func ClearTerminal() {
	cmd := exec.Command("clear")
	if _, ok := os.LookupEnv("OS"); ok {
		cmd = exec.Command("cls") // para Windows
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Pause(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func WaitEnter(scanner *bufio.Scanner) {
	_ = scanner.Scan()
}

