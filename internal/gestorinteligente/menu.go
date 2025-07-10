package gestorinteligente

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getGestorMenuStr() string {
	return `--- MENU GESTOR INTELIGENTE ---
1. Adicionar produto
2. Remover produto
3. Listar meses
4. Atualizar lucro mensal
5. Editar produto
6. Antecipar parcelas
7. Configurar porcentagem segura
8. Voltar ao menu principal`
}

func getGestorResumoStr(list ProductList) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("RESUMO DO MÊS (%02d/%d - %s)", time.Now().Month(), time.Now().Year(), monthNames[time.Now().Month()-1]))
	b.WriteString("\n------------------------------")
	b.WriteString(fmt.Sprintf("\nLucro mensal: R$%.2f", list.MonthlyProfit))
	b.WriteString(fmt.Sprintf("\nTotal de parcelas: R$%.2f", 0.0)) // Placeholder, pode melhorar
	b.WriteString("\n...")                                         // Pode adicionar mais detalhes se quiser
	return b.String()
}

func getGestorMensagemStr(list ProductList) string {
	if list.MonthlyProfit == 0 {
		return "Por favor, defina seu lucro mensal antes de adicionar produtos."
	}
	return ""
}

func PrintGestorMenuCompleto(list ProductList) {
	titulo := " Gestor Inteligente de Gastos "
	resumo := captureShowSummary(list)
	mensagem := getGestorMensagemStr(list)
	menu := getGestorMenuStr()

	// Quebrar em linhas
	blocos := [][]string{
		{titulo},
		strings.Split(resumo, "\n"),
	}
	if mensagem != "" {
		blocos = append(blocos, strings.Split(mensagem, "\n"))
	}
	blocos = append(blocos, strings.Split(menu, "\n"))

	// Calcular o maior comprimento
	maxLen := 0
	for _, bloco := range blocos {
		for _, l := range bloco {
			if len(l) > maxLen {
				maxLen = len(l)
			}
		}
	}
	if maxLen < 60 {
		maxLen = 60
	}

	linhaTopo := "╔" + strings.Repeat("═", maxLen+2) + "╗"
	linhaDiv := "╟" + strings.Repeat("─", maxLen+2) + "╢"
	linhaBase := "╚" + strings.Repeat("═", maxLen+2) + "╝"

	fmt.Println(linhaTopo)
	for i, bloco := range blocos {
		for _, l := range bloco {
			fmt.Printf("║ %-*s ║\n", maxLen, l)
		}
		if i < len(blocos)-1 {
			fmt.Println(linhaDiv)
		}
	}
	fmt.Println(linhaBase)
}

// Função para capturar a saída de showSummary como string
func captureShowSummary(list ProductList) string {
	var b strings.Builder
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	showSummary(list)
	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	b.Write(buf[:n])
	return strings.TrimRight(b.String(), "\n")
}

func ShowGestorMenu() {
	reader := bufio.NewReader(os.Stdin)
	list, _ := LoadProducts()

	if list.SafePercentage == 0 {
		list.SafePercentage = 70
	}

	// Verificar se é a primeira vez e o lucro mensal não foi definido
	if list.MonthlyProfit == 0 {
		ClearTerminal()
		PrintCaixa([]string{
			"🎯 CONFIGURAÇÃO INICIAL",
			"",
			"Bem-vindo ao Gestor Inteligente de Gastos!",
			"",
			"Para começar a usar o sistema, precisamos definir",
			"seu lucro mensal. Este valor será usado para",
			"calcular quanto você pode gastar em parcelas.",
			"",
			"Digite seu lucro mensal (R$):",
		})
		fmt.Print("→ ")
		profitStr, _ := reader.ReadString('\n')
		profitStr = strings.TrimSpace(profitStr)
		profitStr = strings.ReplaceAll(profitStr, ",", ".")

		profit, err := strconv.ParseFloat(profitStr, 64)
		if err != nil || profit <= 0 {
			PrintCaixa([]string{
				"❌ Valor inválido!",
				"",
				"Por favor, digite um valor válido maior que zero.",
				"Exemplo: 5000 ou 5000,50",
			})
			time.Sleep(3 * time.Second)
			return
		}

		list.MonthlyProfit = profit
		SaveProducts(list)

		PrintCaixa([]string{
			"✅ Lucro mensal configurado com sucesso!",
			"",
			fmt.Sprintf("Seu lucro mensal: R$%.2f", profit),
			"",
			"Agora você pode começar a adicionar produtos.",
		})
		time.Sleep(3 * time.Second)
	}

	for {
		ClearTerminal()
		PrintGestorMenuCompleto(list)
		fmt.Print("Escolha uma opcão: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			ClearTerminal()
			addProduct(reader, &list)
		case "2":
			ClearTerminal()
			removeProduct(reader, &list)
		case "3":
			ClearTerminal()
			listMonths(reader, list)
		case "4":
			ClearTerminal()
			updateMonthlyProfit(reader, &list)
		case "5":
			ClearTerminal()
			editProduct(reader, &list)
		case "6":
			ClearTerminal()
			anticipateInstallments(reader, &list)
		case "7":
			ClearTerminal()
			configureSafePercentage(reader, &list)
		case "8":
			SaveProducts(list)
			return
		default:
			fmt.Println("Opcão inválida.")
			time.Sleep(1 * time.Second)
		}
		SaveProducts(list)
	}
}
