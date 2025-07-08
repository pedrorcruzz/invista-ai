package internal

import (
	"bufio"
	"fmt"
)

func InputBox(prompt string, scanner *bufio.Scanner) string {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Printf("║ %-48s ║\n", prompt)
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Print("→ ")
	scanner.Scan()
	return scanner.Text()
}

