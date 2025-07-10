package gestorinteligente

import (
	"bufio"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func addProduct(reader *bufio.Reader, list *ProductList) {
	title := "ADICIONAR PRODUTO"
	for {
		lines := []string{title, "", "0. Voltar ao Menu"}
		PrintCaixa(lines)
		fmt.Print("Nome do produto (0 para voltar): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "0" {
			return
		}

		if _, err := strconv.Atoi(name); err == nil {
			PrintCaixa([]string{"Nome inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("Valor total do produto (R$) (0 para voltar): ")
		valueStr, _ := reader.ReadString('\n')
		valueStr = strings.TrimSpace(valueStr)
		valueStr = strings.ReplaceAll(valueStr, ",", ".")

		if valueStr == "0" {
			return
		}

		totalValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			PrintCaixa([]string{"Valor inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("Em quantas vezes será parcelado (0 para voltar): ")
		installmentsStr, _ := reader.ReadString('\n')
		installmentsStr = strings.TrimSpace(installmentsStr)

		if installmentsStr == "0" {
			return
		}

		installments, err := strconv.Atoi(installmentsStr)
		if err != nil || installments < 1 {
			PrintCaixa([]string{"Número de parcelas inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		parcel := totalValue / float64(installments)

		list.Products = append(list.Products, Product{
			Name:         name,
			Parcel:       parcel,
			TotalValue:   totalValue,
			Installments: installments,
			CreatedAt:    time.Now(),
		})
		list.Month = int(time.Now().Month())
		list.Year = time.Now().Year()

		PrintCaixa([]string{"✅ Produto adicionado! Parcela mensal: R$" + fmt.Sprintf("%.2f", parcel)})
		time.Sleep(2 * time.Second)
		return
	}
}

func removeProduct(reader *bufio.Reader, list *ProductList) {
	title := " REMOVER PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para remover.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonthOrName(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("\nTem certeza que deseja remover '%s'? (s/n): ", list.Products[idx].Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "s" && confirm != "sim" {
		fmt.Println("Operação cancelada.")
		time.Sleep(2 * time.Second)
		return
	}

	list.Products = slices.Delete(list.Products, idx, idx+1)

	fmt.Println(divider)
	fmt.Println("✅ Produto removido!")
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func editProduct(reader *bufio.Reader, list *ProductList) {
	title := " EDITAR PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para editar.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonthOrName(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &list.Products[idx]

	fmt.Printf("Nome atual: %s. Novo nome (ou Enter para manter, 0 para voltar): ", p.Name)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	if newName == "0" {
		return
	}

	if newName != "" {
		p.Name = newName
	}

	fmt.Printf("Valor total atual: R$%.2f. Novo valor (ou Enter para manter, 0 para voltar): ", p.TotalValue)
	totalValueStr, _ := reader.ReadString('\n')
	totalValueStr = strings.TrimSpace(totalValueStr)

	if totalValueStr == "0" {
		return
	}

	if totalValueStr != "" {
		totalValueStr = strings.ReplaceAll(totalValueStr, ",", ".")
		totalValue, err := strconv.ParseFloat(totalValueStr, 64)
		if err == nil && totalValue > 0 {
			p.TotalValue = totalValue
		}
	}

	fmt.Printf("Parcelas atuais: %d. Novo número de parcelas (ou Enter para manter, 0 para voltar): ", p.Installments)
	installmentsStr, _ := reader.ReadString('\n')
	installmentsStr = strings.TrimSpace(installmentsStr)

	if installmentsStr == "0" {
		return
	}

	if installmentsStr != "" {
		installments, err := strconv.Atoi(installmentsStr)
		if err == nil && installments > 0 {
			p.Installments = installments
		}
	}

	p.Parcel = p.TotalValue / float64(p.Installments)

	fmt.Println(divider)
	fmt.Println("✅ Produto atualizado!")
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func anticipateInstallments(reader *bufio.Reader, list *ProductList) {
	title := " ANTECIPAR PARCELAS "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para antecipar.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonthOrName(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &list.Products[idx]
	now := time.Now()
	currentInstallment := getInstallmentNumber(*p, now.Year(), int(now.Month()))

	if currentInstallment >= p.Installments {
		fmt.Println("Este produto já foi totalmente pago.")
		time.Sleep(2 * time.Second)
		return
	}

	remainingInstallments := p.Installments - currentInstallment
	remainingValue := p.Parcel * float64(remainingInstallments)

	fmt.Printf("Produto: %s\n", p.Name)
	fmt.Printf("Parcela atual: %d/%d\n", currentInstallment, p.Installments)
	fmt.Printf("Valor restante: R$%.2f\n", remainingValue)
	fmt.Printf("Parcelas restantes: %d\n", remainingInstallments)

	fmt.Print("Quantas parcelas deseja antecipar? (0 para voltar): ")
	anticipateStr, _ := reader.ReadString('\n')
	anticipateStr = strings.TrimSpace(anticipateStr)

	if anticipateStr == "0" {
		return
	}

	anticipateCount, err := strconv.Atoi(anticipateStr)
	if err != nil || anticipateCount < 1 || anticipateCount > remainingInstallments {
		fmt.Println("Número de parcelas inválido.")
		time.Sleep(2 * time.Second)
		return
	}

	anticipateValue := p.Parcel * float64(anticipateCount)
	newRemainingInstallments := remainingInstallments - anticipateCount

	if newRemainingInstallments > 0 {
		newParcel := (remainingValue - anticipateValue) / float64(newRemainingInstallments)
		p.Parcel = newParcel
		p.Installments = currentInstallment + newRemainingInstallments
	} else {
		p.Installments = currentInstallment
	}

	fmt.Println(divider)
	fmt.Printf("✅ Parcelas antecipadas! Valor antecipado: R$%.2f\n", anticipateValue)
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func updateMonthlyProfit(reader *bufio.Reader, list *ProductList) {
	title := " ATUALIZAR LUCRO MENSAL "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	fmt.Printf("Lucro mensal atual: R$%.2f\n", list.MonthlyProfit)
	fmt.Print("Novo lucro mensal (R$): ")
	profitStr, _ := reader.ReadString('\n')
	profitStr = strings.TrimSpace(profitStr)
	profitStr = strings.ReplaceAll(profitStr, ",", ".")

	if profitStr == "0" {
		return
	}

	profit, err := strconv.ParseFloat(profitStr, 64)
	if err != nil || profit < 0 {
		fmt.Println("Valor inválido.")
		time.Sleep(2 * time.Second)
		return
	}

	list.MonthlyProfit = profit

	fmt.Println(divider)
	fmt.Printf("✅ Lucro mensal atualizado para R$%.2f!\n", profit)
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func configureSafePercentage(reader *bufio.Reader, list *ProductList) {
	title := " CONFIGURAR PORCENTAGEM SEGURA "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	fmt.Printf("Porcentagem segura atual: %.0f%%\n", list.SafePercentage)
	fmt.Print("Nova porcentagem segura (%): ")
	percentageStr, _ := reader.ReadString('\n')
	percentageStr = strings.TrimSpace(percentageStr)
	percentageStr = strings.ReplaceAll(percentageStr, ",", ".")

	if percentageStr == "0" {
		return
	}

	percentage, err := strconv.ParseFloat(percentageStr, 64)
	if err != nil || percentage < 0 || percentage > 100 {
		fmt.Println("Porcentagem inválida.")
		time.Sleep(2 * time.Second)
		return
	}

	list.SafePercentage = percentage

	fmt.Println(divider)
	fmt.Printf("✅ Porcentagem segura atualizada para %.0f%%!\n", percentage)
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func listMonths(reader *bufio.Reader, list ProductList) {
	title := "LISTAR MESES"
	lines := []string{title, ""}

	byYearMonth := mapProductsByYearMonth(list.Products)
	if len(byYearMonth) == 0 {
		lines = append(lines, "Nenhum produto cadastrado.")
		PrintCaixa(lines)
		return
	}

	years := make([]int, 0, len(byYearMonth))
	for y := range byYearMonth {
		years = append(years, y)
	}
	sort.Ints(years)

	lines = append(lines, "Selecione o ano (0 para voltar):")
	for i, y := range years {
		lines = append(lines, fmt.Sprintf("%d. %d", i+1, y))
	}
	PrintCaixa(lines)
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return
	}

	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(years) {
		PrintCaixa([]string{"Ano inválido."})
		time.Sleep(2 * time.Second)
		return
	}
	year := years[yearIdx-1]

	monthsMap := byYearMonth[year]
	months := make([]int, 0, len(monthsMap))
	for m := range monthsMap {
		months = append(months, m)
	}
	sort.Ints(months)

	monthLines := []string{fmt.Sprintf("Meses de %d:", year)}
	for i, m := range months {
		monthLines = append(monthLines, fmt.Sprintf("%d. %s", i+1, monthNames[m-1]))
	}
	PrintCaixa(monthLines)
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return
	}

	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(months) {
		PrintCaixa([]string{"Mês inválido."})
		time.Sleep(2 * time.Second)
		return
	}
	month := months[monthIdx-1]

	prodIndexes := monthsMap[month]
	if len(prodIndexes) == 0 {
		PrintCaixa([]string{fmt.Sprintf("Nenhum produto encontrado para %s/%d.", monthNames[month-1], year)})
		time.Sleep(2 * time.Second)
		return
	}

	uniqueIndexes := make([]int, 0)
	seen := make(map[int]bool)
	for _, idx := range prodIndexes {
		if !seen[idx] {
			seen[idx] = true
			uniqueIndexes = append(uniqueIndexes, idx)
		}
	}

	productsTitle := fmt.Sprintf("PRODUTOS DE %s/%d", monthNames[month-1], year)
	productLines := []string{productsTitle, ""}
	for i, idx := range uniqueIndexes {
		p := list.Products[idx]
		installmentNumber := getInstallmentNumber(p, year, month)
		productLines = append(productLines, fmt.Sprintf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d)",
			i+1, p.Name, p.TotalValue, p.Parcel, installmentNumber, p.Installments))
	}
	PrintCaixa(productLines)

	fmt.Print("\nPressione Enter para voltar...")
	reader.ReadString('\n')
}

func selectProductByYearMonthOrName(reader *bufio.Reader, products []Product) (int, bool) {
	byYearMonth := mapProductsByYearMonth(products)
	if len(byYearMonth) == 0 {
		PrintCaixa([]string{"Nenhum produto cadastrado."})
		return -1, false
	}

	years := make([]int, 0, len(byYearMonth))
	for y := range byYearMonth {
		years = append(years, y)
	}
	sort.Ints(years)

	lines := []string{"Selecione o ano (0 para voltar):"}
	for i, y := range years {
		lines = append(lines, fmt.Sprintf("%d. %d", i+1, y))
	}
	PrintCaixa(lines)
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return -1, false
	}
	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(years) {
		PrintCaixa([]string{"Ano inválido."})
		return -1, false
	}
	year := years[yearIdx-1]

	monthsMap := byYearMonth[year]
	months := make([]int, 0, len(monthsMap))
	for m := range monthsMap {
		months = append(months, m)
	}
	sort.Ints(months)

	monthLines := []string{fmt.Sprintf("Meses de %d:", year)}
	for i, m := range months {
		monthLines = append(monthLines, fmt.Sprintf("%d. %s", i+1, monthNames[m-1]))
	}
	PrintCaixa(monthLines)
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return -1, false
	}
	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(months) {
		PrintCaixa([]string{"Mês inválido."})
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

	productLines := []string{"Selecione o produto (número ou nome, 0 para voltar):"}
	for i, idx := range uniqueIndexes {
		p := products[idx]
		productLines = append(productLines, fmt.Sprintf("%d - %s | Total: R$%.2f | Parcelas: %d | Adicionado em: %s",
			i+1, p.Name, p.TotalValue, p.Installments, p.CreatedAt.Format("02/01/2006")))
	}
	PrintCaixa(productLines)
	fmt.Print("Produto: ")
	prodStr, _ := reader.ReadString('\n')
	prodStr = strings.TrimSpace(prodStr)

	if prodStr == "0" {
		return -1, false
	}

	// Tenta por número
	prodIdx, err := strconv.Atoi(prodStr)
	if err == nil && prodIdx >= 1 && prodIdx <= len(uniqueIndexes) {
		return uniqueIndexes[prodIdx-1], true
	}
	// Tenta por nome
	for _, idx := range uniqueIndexes {
		if strings.EqualFold(products[idx].Name, prodStr) {
			return idx, true
		}
	}
	PrintCaixa([]string{"Produto inválido."})
	return -1, false
}
