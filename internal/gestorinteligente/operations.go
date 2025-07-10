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

func adicionarProduto(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := "ADICIONAR PRODUTO"
	for {
		linhas := []string{titulo, "", "0. Voltar ao Menu"}
		ImprimirCaixa(linhas)
		fmt.Print("Nome do produto (0 para voltar): ")
		nome, _ := reader.ReadString('\n')
		nome = strings.TrimSpace(nome)

		if nome == "0" {
			return
		}

		if _, err := strconv.Atoi(nome); err == nil {
			ImprimirCaixa([]string{"Nome inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("Valor total do produto (R$) (0 para voltar): ")
		valorStr, _ := reader.ReadString('\n')
		valorStr = strings.TrimSpace(valorStr)
		valorStr = strings.ReplaceAll(valorStr, ",", ".")

		if valorStr == "0" {
			return
		}

		valorTotal, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			ImprimirCaixa([]string{"Valor inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("Em quantas vezes será parcelado (0 para voltar): ")
		parcelasStr, _ := reader.ReadString('\n')
		parcelasStr = strings.TrimSpace(parcelasStr)

		if parcelasStr == "0" {
			return
		}

		parcelas, err := strconv.Atoi(parcelasStr)
		if err != nil || parcelas < 1 {
			ImprimirCaixa([]string{"Número de parcelas inválido."})
			time.Sleep(2 * time.Second)
			continue
		}

		parcela := valorTotal / float64(parcelas)

		lista.Produtos = append(lista.Produtos, Produto{
			Nome:       nome,
			Parcela:    parcela,
			ValorTotal: valorTotal,
			Parcelas:   parcelas,
			CriadoEm:   time.Now(),
		})
		lista.Mes = int(time.Now().Month())
		lista.Ano = time.Now().Year()

		ImprimirCaixa([]string{"✅ Produto adicionado! Parcela mensal: R$" + fmt.Sprintf("%.2f", parcela)})
		time.Sleep(2 * time.Second)
		return
	}
}

func removerProduto(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := " REMOVER PRODUTO "
	divisor := strings.Repeat("-", 40)

	fmt.Println("\n" + divisor)
	fmt.Println(titulo)
	fmt.Println(divisor)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divisor)

	if len(lista.Produtos) == 0 {
		fmt.Println("Nenhum produto para remover.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selecionarProdutoPorAnoMesOuNome(reader, lista.Produtos)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("\nTem certeza que deseja remover '%s'? (s/n): ", lista.Produtos[idx].Nome)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "s" && confirm != "sim" {
		fmt.Println("Operação cancelada.")
		time.Sleep(2 * time.Second)
		return
	}

	lista.Produtos = slices.Delete(lista.Produtos, idx, idx+1)

	fmt.Println(divisor)
	fmt.Println("✅ Produto removido!")
	fmt.Println(divisor)

	time.Sleep(2 * time.Second)
}

func editarProduto(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := "EDITAR PRODUTO"

	if len(lista.Produtos) == 0 {
		ImprimirCaixa([]string{"Nenhum produto para editar."})
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selecionarProdutoPorAnoMesOuNome(reader, lista.Produtos)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &lista.Produtos[idx]

	// Nome
	ImprimirCaixa([]string{
		titulo,
		"",
		"0. Voltar ao Menu",
		"",
		"Nome atual: " + p.Nome,
		"Digite o novo nome (ou Enter para manter):",
	})
	fmt.Print(" [1mNovo nome: [0m ")
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)
	if newName == "0" {
		return
	}
	if newName != "" {
		p.Nome = newName
	}

	// Valor total
	ImprimirCaixa([]string{
		titulo,
		"",
		"0. Voltar ao Menu",
		"",
		fmt.Sprintf("Valor total atual: R$%.2f", p.ValorTotal),
		"Digite o novo valor (ou Enter para manter):",
	})
	fmt.Print(" [1mNovo valor: [0m ")
	totalValueStr, _ := reader.ReadString('\n')
	totalValueStr = strings.TrimSpace(totalValueStr)
	if totalValueStr == "0" {
		return
	}
	if totalValueStr != "" {
		totalValueStr = strings.ReplaceAll(totalValueStr, ",", ".")
		totalValue, err := strconv.ParseFloat(totalValueStr, 64)
		if err == nil && totalValue > 0 {
			p.ValorTotal = totalValue
		}
	}

	// Parcelas
	ImprimirCaixa([]string{
		titulo,
		"",
		"0. Voltar ao Menu",
		"",
		fmt.Sprintf("Parcelas atuais: %d", p.Parcelas),
		"Digite o novo número de parcelas (ou Enter para manter):",
	})
	fmt.Print(" [1mNovas parcelas: [0m ")
	installmentsStr, _ := reader.ReadString('\n')
	installmentsStr = strings.TrimSpace(installmentsStr)
	if installmentsStr == "0" {
		return
	}
	if installmentsStr != "" {
		installments, err := strconv.Atoi(installmentsStr)
		if err == nil && installments > 0 {
			p.Parcelas = installments
		}
	}

	p.Parcela = p.ValorTotal / float64(p.Parcelas)

	ImprimirCaixa([]string{"✅ Produto atualizado!"})
	time.Sleep(2 * time.Second)
}

func anteciparParcelas(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := " ANTECIPAR PARCELAS "
	divisor := strings.Repeat("-", 40)

	fmt.Println("\n" + divisor)
	fmt.Println(titulo)
	fmt.Println(divisor)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divisor)

	if len(lista.Produtos) == 0 {
		fmt.Println("Nenhum produto para antecipar.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selecionarProdutoPorAnoMesOuNome(reader, lista.Produtos)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &lista.Produtos[idx]
	now := time.Now()
	currentInstallment := obterNumeroParcela(*p, now.Year(), int(now.Month()))

	if currentInstallment >= p.Parcelas {
		fmt.Println("Este produto já foi totalmente pago.")
		time.Sleep(2 * time.Second)
		return
	}

	remainingInstallments := p.Parcelas - currentInstallment
	remainingValue := p.Parcela * float64(remainingInstallments)

	fmt.Printf("Produto: %s\n", p.Nome)
	fmt.Printf("Parcela atual: %d/%d\n", currentInstallment, p.Parcelas)
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

	anticipateValue := p.Parcela * float64(anticipateCount)
	newRemainingInstallments := remainingInstallments - anticipateCount

	if newRemainingInstallments > 0 {
		newParcel := (remainingValue - anticipateValue) / float64(newRemainingInstallments)
		p.Parcela = newParcel
		p.Parcelas = currentInstallment + newRemainingInstallments
	} else {
		p.Parcelas = currentInstallment
	}

	fmt.Println(divisor)
	fmt.Printf("✅ Parcelas antecipadas! Valor antecipado: R$%.2f\n", anticipateValue)
	fmt.Println(divisor)

	time.Sleep(2 * time.Second)
}

func atualizarLucroMensal(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := " ATUALIZAR LUCRO MENSAL "
	linhas := []string{
		titulo,
		"",
		"0. Voltar ao Menu",
		"",
		fmt.Sprintf("Lucro mensal atual: R$%.2f", lista.LucroMensal),
		"",
	}
	ImprimirCaixa(linhas)
	fmt.Print("Digite o novo lucro mensal (R$): ")
	profitStr, _ := reader.ReadString('\n')
	profitStr = strings.TrimSpace(profitStr)
	profitStr = strings.ReplaceAll(profitStr, ",", ".")

	if profitStr == "0" {
		return
	}

	profit, err := strconv.ParseFloat(profitStr, 64)
	if err != nil || profit < 0 {
		ImprimirCaixa([]string{"Valor inválido."})
		time.Sleep(2 * time.Second)
		return
	}

	lista.LucroMensal = profit

	ImprimirCaixa([]string{fmt.Sprintf("✅ Lucro mensal atualizado para R$%.2f!", profit)})
	time.Sleep(2 * time.Second)
}

func configurarPorcentagemSegura(reader *bufio.Reader, lista *ListaProdutos) {
	titulo := " CONFIGURAR PORCENTAGEM SEGURA "
	divisor := strings.Repeat("-", 40)

	fmt.Println("\n" + divisor)
	fmt.Println(titulo)
	fmt.Println(divisor)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divisor)

	fmt.Printf("Porcentagem segura atual: %.0f%%\n", lista.PorcentagemSegura)
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

	lista.PorcentagemSegura = percentage

	fmt.Println(divisor)
	fmt.Printf("✅ Porcentagem segura atualizada para %.0f%%!\n", percentage)
	fmt.Println(divisor)

	time.Sleep(2 * time.Second)
}

func listarMeses(reader *bufio.Reader, lista ListaProdutos) {
	titulo := "LISTAR MESES"
	linhas := []string{titulo, ""}

	produtosPorAnoMes := mapearProdutosPorAnoMes(lista.Produtos)
	if len(produtosPorAnoMes) == 0 {
		linhas = append(linhas, "Nenhum produto cadastrado.")
		ImprimirCaixa(linhas)
		return
	}

	anos := make([]int, 0, len(produtosPorAnoMes))
	for y := range produtosPorAnoMes {
		anos = append(anos, y)
	}
	sort.Ints(anos)

	linhas = append(linhas, "Selecione o ano (0 para voltar):")
	for i, y := range anos {
		linhas = append(linhas, fmt.Sprintf("%d. %d", i+1, y))
	}
	ImprimirCaixa(linhas)
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return
	}

	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(anos) {
		ImprimirCaixa([]string{"Ano inválido."})
		time.Sleep(2 * time.Second)
		return
	}
	year := anos[yearIdx-1]

	mesesMap := produtosPorAnoMes[year]
	meses := make([]int, 0, len(mesesMap))
	for m := range mesesMap {
		meses = append(meses, m)
	}
	sort.Ints(meses)

	monthLines := []string{fmt.Sprintf("Meses de %d:", year)}
	for i, m := range meses {
		monthLines = append(monthLines, fmt.Sprintf("%d. %s", i+1, nomesMeses[m-1]))
	}
	ImprimirCaixa(monthLines)
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return
	}

	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(meses) {
		ImprimirCaixa([]string{"Mês inválido."})
		time.Sleep(2 * time.Second)
		return
	}
	month := meses[monthIdx-1]

	prodIndexes := mesesMap[month]
	if len(prodIndexes) == 0 {
		ImprimirCaixa([]string{fmt.Sprintf("Nenhum produto encontrado para %s/%d.", nomesMeses[month-1], year)})
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

	productsTitle := fmt.Sprintf("PRODUTOS DE %s/%d", nomesMeses[month-1], year)
	productLines := []string{productsTitle, ""}
	for i, idx := range uniqueIndexes {
		p := lista.Produtos[idx]
		installmentNumber := obterNumeroParcela(p, year, month)
		productLines = append(productLines, fmt.Sprintf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d)",
			i+1, p.Nome, p.ValorTotal, p.Parcela, installmentNumber, p.Parcelas))
	}
	ImprimirCaixa(productLines)

	fmt.Print("\nPressione Enter para voltar...")
	reader.ReadString('\n')
}

func selecionarProdutoPorAnoMesOuNome(reader *bufio.Reader, produtos []Produto) (int, bool) {
	produtosPorAnoMes := mapearProdutosPorAnoMes(produtos)
	if len(produtosPorAnoMes) == 0 {
		ImprimirCaixa([]string{"Nenhum produto cadastrado."})
		return -1, false
	}

	anos := make([]int, 0, len(produtosPorAnoMes))
	for y := range produtosPorAnoMes {
		anos = append(anos, y)
	}
	sort.Ints(anos)

	linhas := []string{"Selecione o ano (0 para voltar):"}
	for i, y := range anos {
		linhas = append(linhas, fmt.Sprintf("%d. %d", i+1, y))
	}
	ImprimirCaixa(linhas)
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return -1, false
	}
	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(anos) {
		ImprimirCaixa([]string{"Ano inválido."})
		return -1, false
	}
	year := anos[yearIdx-1]

	mesesMap := produtosPorAnoMes[year]
	meses := make([]int, 0, len(mesesMap))
	for m := range mesesMap {
		meses = append(meses, m)
	}
	sort.Ints(meses)

	monthLines := []string{fmt.Sprintf("Meses de %d:", year)}
	for i, m := range meses {
		monthLines = append(monthLines, fmt.Sprintf("%d. %s", i+1, nomesMeses[m-1]))
	}
	ImprimirCaixa(monthLines)
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return -1, false
	}
	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(meses) {
		ImprimirCaixa([]string{"Mês inválido."})
		return -1, false
	}
	month := meses[monthIdx-1]

	prodIndexes := mesesMap[month]
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
		p := produtos[idx]
		productLines = append(productLines, fmt.Sprintf("%d - %s | Total: R$%.2f | Parcelas: %d | Adicionado em: %s",
			i+1, p.Nome, p.ValorTotal, p.Parcelas, p.CriadoEm.Format("02/01/2006")))
	}
	ImprimirCaixa(productLines)
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
		if strings.EqualFold(produtos[idx].Nome, prodStr) {
			return idx, true
		}
	}
	ImprimirCaixa([]string{"Produto inválido."})
	return -1, false
}
