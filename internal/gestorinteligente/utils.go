package gestorinteligente

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var monthNames = []string{
	"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
	"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
}

func PrintCaixa(lines []string) {
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	if maxLen < 60 {
		maxLen = 60
	}
	linhaTopo := "╔" + strings.Repeat("═", maxLen+2) + "╗"
	linhaBase := "╚" + strings.Repeat("═", maxLen+2) + "╝"
	fmt.Println(linhaTopo)
	for _, l := range lines {
		fmt.Printf("║ %-*s ║\n", maxLen, l)
	}
	fmt.Println(linhaBase)
}

func ClearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func mapProductsByYearMonth(products []Product) map[int]map[int][]int {
	result := make(map[int]map[int][]int)
	for idx, p := range products {
		endDate := p.CreatedAt.AddDate(0, p.Installments-1, 0)
		currentDate := p.CreatedAt

		for !currentDate.After(endDate) {
			currentYear := currentDate.Year()
			currentMonth := int(currentDate.Month())

			if _, ok := result[currentYear]; !ok {
				result[currentYear] = make(map[int][]int)
			}
			result[currentYear][currentMonth] = append(result[currentYear][currentMonth], idx)

			currentDate = currentDate.AddDate(0, 1, 0)
		}
	}
	return result
}

func selectProductByYearMonth(reader *bufio.Reader, products []Product) (int, bool) {
	byYearMonth := mapProductsByYearMonth(products)
	if len(byYearMonth) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return -1, false
	}

	years := make([]int, 0, len(byYearMonth))
	for y := range byYearMonth {
		years = append(years, y)
	}
	sort.Ints(years)

	fmt.Println("\nSelecione o ano (0 para voltar):")
	for i, y := range years {
		fmt.Printf("%d. %d\n", i+1, y)
	}
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return -1, false
	}

	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(years) {
		fmt.Println("Ano inválido.")
		return -1, false
	}
	year := years[yearIdx-1]

	monthsMap := byYearMonth[year]
	months := make([]int, 0, len(monthsMap))
	for m := range monthsMap {
		months = append(months, m)
	}
	sort.Ints(months)

	fmt.Println("\nSelecione o mês (0 para voltar):")
	for i, m := range months {
		fmt.Printf("%d. %s\n", i+1, monthNames[m-1])
	}
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return -1, false
	}

	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(months) {
		fmt.Println("Mês inválido.")
		return -1, false
	}
	month := months[monthIdx-1]

	prodIndexes := monthsMap[month]

	uniqueIndexes := make([]int, 0)
	seen := make(map[int]bool)
	for _, idx := range prodIndexes {
		if !seen[idx] {
			seen[idx] = true
			uniqueIndexes = append(uniqueIndexes, idx)
		}
	}

	fmt.Println("\nSelecione o produto (0 para voltar):")
	for i, idx := range uniqueIndexes {
		p := products[idx]
		fmt.Printf("%d. %s | Total: R$%.2f | Parcelas: %d | Adicionado em: %s\n",
			i+1, p.Name, p.TotalValue, p.Installments, p.CreatedAt.Format("02/01/2006"))
	}
	fmt.Print("Produto: ")
	prodStr, _ := reader.ReadString('\n')
	prodStr = strings.TrimSpace(prodStr)

	if prodStr == "0" {
		return -1, false
	}

	prodIdx, err := strconv.Atoi(prodStr)
	if err != nil || prodIdx < 1 || prodIdx > len(uniqueIndexes) {
		fmt.Println("Produto inválido.")
		return -1, false
	}
	return uniqueIndexes[prodIdx-1], true
}

func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Print(prompt)
	valueStr, _ := reader.ReadString('\n')
	valueStr = strings.TrimSpace(valueStr)
	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	return strconv.ParseFloat(valueStr, 64)
}

func isProductActiveInMonth(p Product, targetYear, targetMonth int) bool {
	startDate := p.CreatedAt
	endDate := startDate.AddDate(0, p.Installments-1, 0)

	targetDate := time.Date(targetYear, time.Month(targetMonth), 1, 0, 0, 0, 0, time.UTC)
	return !targetDate.Before(startDate) && !targetDate.After(endDate)
}

func getInstallmentNumber(p Product, targetYear, targetMonth int) int {
	startDate := p.CreatedAt
	yearDiff := targetYear - startDate.Year()
	monthDiff := targetMonth - int(startDate.Month())
	totalMonthDiff := yearDiff*12 + monthDiff + 1

	if totalMonthDiff < 1 {
		return 1
	}
	if totalMonthDiff > p.Installments {
		return p.Installments
	}
	return totalMonthDiff
}
