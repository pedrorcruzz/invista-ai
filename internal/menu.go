package internal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func PrintMenuPrincipalSozinho() {
	ClearTerminal()
	Pause(300)
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║ --- MENU PRINCIPAL ---                             ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Ver resumo completo (visualização vertical)      ║")
	fmt.Println("║ 2. Ver resumo completo (tabela horizontal)          ║")
	fmt.Println("║ 3. Adicionar/editar mês                             ║")
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
1. Ver resumo completo (visualização vertical)
2. Ver resumo completo (tabela horizontal)
3. Adicionar/editar mês
4. Sair do programa`
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
	ano := InputBox("Digite o ano(YYYY):", scanner)
	mes := InputBox("Digite o mês(MM):", scanner)
	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}
	m := dados.Anos[ano][mes]
	if m != (Mes{}) {
		for {
			fmt.Println("\n--- EDITAR CAMPOS ---")
			fmt.Printf("1. Aporte RF (atual: %s)\n", FormatFloatBR(m.AporteRF))
			fmt.Printf("2. Aporte FIIs (atual: %s)\n", FormatFloatBR(m.AporteFIIs))
			fmt.Printf("3. Saída (atual: %s)\n", FormatFloatBR(m.Saida))
			fmt.Printf("4. Valor Bruto RF (atual: %s)\n", FormatFloatBR(m.ValorBrutoRF))
			fmt.Printf("5. Valor Líquido RF (atual: %s)\n", FormatFloatBR(m.ValorLiquidoRF))
			fmt.Printf("6. Valor Líquido FIIs (atual: %s)\n", FormatFloatBR(m.ValorLiquidoFIIs))
			fmt.Printf("7. Lucro Retirado (atual: %s)\n", FormatFloatBR(m.LucroRetirado))
			fmt.Printf("8. Lucro Líquido FIIs (atual: %s)\n", FormatFloatBR(m.LucroLiquidoFIIs))
			fmt.Println("0. Sair da edição")
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
				fmt.Println("Edição concluída.")
				return
			default:
				fmt.Println("Opção inválida.")
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
	fmt.Println("Dados adicionados com sucesso!")
}

// ParseFloatBR converte string com vírgula ou ponto para float64
func ParseFloatBR(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}
