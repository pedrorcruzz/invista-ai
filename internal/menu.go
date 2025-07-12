package internal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func PrintMenuPrincipalSozinho() {
	ClearTerminal()
	Pause(300)
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║ --- MENU PRINCIPAL ---                             ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Ver resumo completo                              ║")
	fmt.Println("║ 2. Adicionar/editar mês                             ║")
	fmt.Println("║ 3. Gestor Inteligente de Gastos                     ║")
	fmt.Println("║ 4. Voltar para o menu inicial                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func PrintTelaUnificada(dados Dados) {
	ClearTerminal()
	Pause(300)
	resumoTotal := GetResumoTotalAcumuladoStr(dados)
	resumoMes := GetResumoMesAtualStr(dados)
	menu := GetMenuPrincipalStr()

	maxLen := 0
	for _, s := range []string{resumoTotal, resumoMes, menu} {
		for _, l := range splitLines(s) {
			if len(l) > maxLen {
				maxLen = len(l)
			}
		}
	}
	if maxLen < 60 {
		maxLen = 60
	}

	linhaTopo := "╔" + repeatStr("═", maxLen+2) + "╗"
	linhaDiv := "╟" + repeatStr("─", maxLen+2) + "╢"
	linhaBase := "╚" + repeatStr("═", maxLen+2) + "╝"

	fmt.Println(linhaTopo)
	for i, bloco := range []string{resumoTotal, resumoMes, menu} {
		for _, l := range splitLines(bloco) {
			fmt.Printf("║ %-*s ║\n", maxLen, l)
		}
		if i < 2 {
			fmt.Println(linhaDiv)
		}
	}
	fmt.Println(linhaBase)
}

func splitLines(s string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func repeatStr(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}

func GetMenuPrincipalStr() string {
	return `--- MENU PRINCIPAL ---
1. Ver resumo completo
2. Adicionar/editar mês
3. Gestor Inteligente de Gastos
4. Retirar Lucro
5. Sair do programa`
}

func SelecionarAno(dados Dados, scanner *bufio.Scanner) string {
	if len(dados.Anos) == 0 {
		fmt.Println("Nenhum dado disponível ainda.")
		return ""
	}
	anos := OrdenarChaves(dados.Anos)
	fmt.Println("\nAnos disponíveis:")
	for i, a := range anos {
		fmt.Printf("%d - %s\n", i+1, a)
	}
	fmt.Print("Digite o número ou o ano desejado (YYYY): ")
	scanner.Scan()
	input := scanner.Text()
	if idx, err := strconv.Atoi(input); err == nil {
		if idx >= 1 && idx <= len(anos) {
			return anos[idx-1]
		}
	}
	for _, a := range anos {
		if a == input {
			return a
		}
	}
	fmt.Printf("Não há dados para o ano ou opção '%s'.\n", input)
	fmt.Println("Anos disponíveis:")
	for _, a := range anos {
		fmt.Println(" -", a)
	}
	return ""
}

func AdicionarOuEditarMes(dados *Dados, scanner *bufio.Scanner) {
	// Seleção de ano
	anos := OrdenarChaves(dados.Anos)
	ano := ""
	if len(anos) > 0 {
		fmt.Println("Anos disponíveis:")
		for i, a := range anos {
			fmt.Printf("%d - %s\n", i+1, a)
		}
		fmt.Print("Digite o número ou o ano desejado (YYYY): ")
		scanner.Scan()
		input := scanner.Text()
		if idx, err := strconv.Atoi(input); err == nil {
			if idx >= 1 && idx <= len(anos) {
				ano = anos[idx-1]
			}
		}
		if ano == "" {
			for _, a := range anos {
				if a == input {
					ano = a
				}
			}
		}
	}
	if ano == "" {
		ano = InputBox("Digite o ano(YYYY):", scanner)
	}

	// Seleção de mês
	mes := ""
	mesesExistentes := []string{}
	if dados.Anos[ano] != nil {
		mesesExistentes = OrdenarChaves(dados.Anos[ano])
	}
	if len(mesesExistentes) > 0 {
		fmt.Println("Meses disponíveis:")
		for i, m := range mesesExistentes {
			fmt.Printf("%d - %s\n", i+1, NomeMes(m))
		}
		fmt.Print("Digite o número ou o mês desejado (MM): ")
		scanner.Scan()
		input := scanner.Text()
		if idx, err := strconv.Atoi(input); err == nil {
			if idx >= 1 && idx <= len(mesesExistentes) {
				mes = mesesExistentes[idx-1]
			}
		}
		if mes == "" {
			for _, m := range mesesExistentes {
				if m == input {
					mes = m
				}
			}
		}
	}
	if mes == "" {
		mes = InputBox("Digite o mês(MM):", scanner)
	}

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}
	m := dados.Anos[ano][mes]
	if m != (Mes{}) {
		for {
			// Menu de edição em caixinha
			lines := []string{
				"EDITAR CAMPOS",
				"",
				"0. Sair da edição",
				"",
				fmt.Sprintf("1. Aporte RF (atual: %s)", FormatFloatBR(m.AporteRF)),
				fmt.Sprintf("2. Aporte FIIs (atual: %s)", FormatFloatBR(m.AporteFIIs)),
				fmt.Sprintf("3. Saída (atual: %s)", FormatFloatBR(m.Saida)),
				fmt.Sprintf("4. Valor Bruto RF (atual: %s)", FormatFloatBR(m.ValorBrutoRF)),
				fmt.Sprintf("5. Valor Líquido RF (atual: %s)", FormatFloatBR(m.ValorLiquidoRF)),
				fmt.Sprintf("6. Valor Líquido FIIs (atual: %s)", FormatFloatBR(m.ValorLiquidoFIIs)),
				fmt.Sprintf("7. Lucro Retirado (atual: %s)", FormatFloatBR(m.LucroRetirado)),
				fmt.Sprintf("8. Lucro Líquido FIIs (atual: %s)", FormatFloatBR(m.LucroLiquidoFIIs)),
			}
			PrintCaixa(lines)
			opcao := InputBox("Escolha:", scanner)
			switch opcao {
			case "1":
				valor := InputBox("Novo valor:", scanner)
				m.AporteRF, _ = ParseFloatBR(valor)
			case "2":
				valor := InputBox("Novo valor:", scanner)
				m.AporteFIIs, _ = ParseFloatBR(valor)
			case "3":
				valor := InputBox("Novo valor:", scanner)
				m.Saida, _ = ParseFloatBR(valor)
			case "4":
				valor := InputBox("Novo valor:", scanner)
				m.ValorBrutoRF, _ = ParseFloatBR(valor)
			case "5":
				valor := InputBox("Novo valor:", scanner)
				m.ValorLiquidoRF, _ = ParseFloatBR(valor)
			case "6":
				valor := InputBox("Novo valor:", scanner)
				m.ValorLiquidoFIIs, _ = ParseFloatBR(valor)
			case "7":
				valor := InputBox("Novo valor:", scanner)
				m.LucroRetirado, _ = ParseFloatBR(valor)
			case "8":
				valor := InputBox("Novo valor:", scanner)
				m.LucroLiquidoFIIs, _ = ParseFloatBR(valor)
			case "0":
				dados.Anos[ano][mes] = m
				PrintCaixa([]string{"✅ Edição concluída!"})
				return
			default:
				PrintCaixa([]string{"❌ Opção inválida."})
			}
			dados.Anos[ano][mes] = m
		}
	}
	aporteRF, _ := ParseFloatBR(InputBox("Digite o aporte na Renda Fixa: R$", scanner))
	aporteFIIs, _ := ParseFloatBR(InputBox("Digite o aporte em FIIs: R$", scanner))
	saida, _ := ParseFloatBR(InputBox("Digite a saída (retirada) do mês: R$", scanner))
	valorBrutoRF, _ := ParseFloatBR(InputBox("Digite o valor bruto da Renda Fixa: R$", scanner))
	valorLiquidoRF, _ := ParseFloatBR(InputBox("Digite o valor líquido da Renda Fixa: R$", scanner))
	valorLiquidoFIIs, _ := ParseFloatBR(InputBox("Digite o valor líquido dos FIIs: R$", scanner))
	lucroRetirado, _ := ParseFloatBR(InputBox("Digite o valor de lucro retirado: R$", scanner))
	lucroLiquidoFIIs, _ := ParseFloatBR(InputBox("Digite o lucro líquido dos FIIs: R$", scanner))
	dados.Anos[ano][mes] = Mes{
		AporteRF:         aporteRF,
		AporteFIIs:       aporteFIIs,
		Saida:            saida,
		ValorBrutoRF:     valorBrutoRF,
		ValorLiquidoRF:   valorLiquidoRF,
		ValorLiquidoFIIs: valorLiquidoFIIs,
		LucroRetirado:    lucroRetirado,
		LucroLiquidoFIIs: lucroLiquidoFIIs,
	}
	PrintCaixa([]string{"✅ Dados adicionados com sucesso!"})
}

// PrintCaixa exibe uma caixinha com as linhas fornecidas
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
	linhaTopo := "╔" + repeatStr("═", maxLen+2) + "╗"
	linhaBase := "╚" + repeatStr("═", maxLen+2) + "╝"
	fmt.Println(linhaTopo)
	for _, l := range lines {
		fmt.Printf("║ %-*s ║\n", maxLen, l)
	}
	fmt.Println(linhaBase)
}

// ParseFloatBR converte string com vírgula ou ponto para float64
func ParseFloatBR(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

// Função para retirar lucro do mês atual
func RetirarLucro(dados *Dados, scanner *bufio.Scanner) {
	hoje := time.Now()
	anoAtual := fmt.Sprintf("%04d", hoje.Year())
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))

	if dados.Anos[anoAtual] == nil {
		dados.Anos[anoAtual] = make(Ano)
	}
	m := dados.Anos[anoAtual][mesAtual]

	PrintCaixa([]string{"Digite o valor de lucro a retirar (será descontado do Lucro Líquido RF + FIIs):"})
	fmt.Print("Valor: R$ ")
	scanner.Scan()
	valorStr := scanner.Text()
	valor, err := ParseFloatBR(valorStr)
	if err != nil || valor <= 0 {
		PrintCaixa([]string{"❌ Valor inválido!"})
		return
	}
	m.LucroRetirado += valor
	dados.Anos[anoAtual][mesAtual] = m
	PrintCaixa([]string{fmt.Sprintf("✅ Lucro de R$ %s retirado com sucesso!", FormatFloatBR(valor))})
}
