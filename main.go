package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Mes struct {
	AporteRF         float64 `json:"aporte_rf"`
	AporteFIIs       float64 `json:"aporte_fiis"`
	Saida            float64 `json:"saida"`
	ValorBrutoRF     float64 `json:"valor_bruto_rf"`
	ValorLiquidoRF   float64 `json:"valor_liquido_rf"`
	ValorLiquidoFIIs float64 `json:"valor_liquido_fiis"`
	LucroRetirado    float64 `json:"lucro_retirado"`
	LucroLiquidoFIIs float64 `json:"lucro_liquido_fiis"`
}

type Ano map[string]Mes

type Dados struct {
	Anos map[string]Ano `json:"anos"`
}

const arquivo = "dados.json"

func carregarDados() Dados {
	file, err := os.ReadFile(arquivo)
	if err != nil {
		return Dados{Anos: make(map[string]Ano)}
	}

	var dados Dados
	err = json.Unmarshal(file, &dados)
	if err != nil {
		fmt.Println("Erro ao carregar dados:", err)
		return Dados{Anos: make(map[string]Ano)}
	}
	return dados
}

func salvarDados(dados Dados) {
	bytes, err := json.MarshalIndent(dados, "", "  ")
	if err != nil {
		fmt.Println("Erro ao salvar dados:", err)
		return
	}
	os.WriteFile(arquivo, bytes, 0644)
}

func nomeMes(m string) string {
	nomes := map[string]string{
		"01": "Janeiro", "02": "Fevereiro", "03": "Março",
		"04": "Abril", "05": "Maio", "06": "Junho",
		"07": "Julho", "08": "Agosto", "09": "Setembro",
		"10": "Outubro", "11": "Novembro", "12": "Dezembro",
	}
	return nomes[m]
}

func ordenarChaves[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printMenuBox(options []string) {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                  MENU PRINCIPAL                     ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	for i, opt := range options {
		fmt.Printf("║  %d. %-46s║\n", i+1, opt)
	}
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func mostrarResumoTotalAcumulado(dados Dados) {
	anos := ordenarChaves(dados.Anos)
	if len(anos) == 0 {
		fmt.Println("Nenhum dado disponível ainda.")
		return
	}
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	valorBrutoFinal := 0.0
	valorLiquidoRFFinal := 0.0
	valorLiquidoFIIsFinal := 0.0
	lucrosRetiradosTotal := 0.0
	lucroLiquidoAcumulado := 0.0
	lucroLiquidoFIIsAcumulado := 0.0
	lucroMesLiquidoTotalAcumulado := 0.0
	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := ordenarChaves(mesesMap)
		for _, mes := range meses {
			m := mesesMap[mes]
			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
			lucroLiquidoFIIs := m.LucroLiquidoFIIs
			lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
			lucroValido := lucroMesBruto > impostos
			if lucroValido {
				aporteRFSoFar += m.AporteRF
				aporteFIIsSoFar += m.AporteFIIs
				saidaSoFar += m.Saida
				lucrosRetiradosTotal += m.LucroRetirado
				valorBrutoFinal = m.ValorBrutoRF
				valorLiquidoRFFinal = m.ValorLiquidoRF
				valorLiquidoFIIsFinal = m.ValorLiquidoFIIs
				lucroLiquidoAcumulado += lucroMesLiquidoRF
				lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
				lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
				saldoAnterior = m.ValorBrutoRF
			}
		}
	}
	totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
	totalAportadoLiquido := totalAportadoBruto - saidaSoFar
	lucroBrutoTotal := valorBrutoFinal - totalAportadoLiquido
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Printf("║           RESUMO TOTAL ACUMULADO                    ║\n")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Printf("║ Total aportado bruto:      R$ %10.2f               ║\n", totalAportadoBruto)
	fmt.Printf("║ Total aportado líquido:    R$ %10.2f               ║\n", totalAportadoLiquido)
	fmt.Printf("║ Valor bruto final (RF):    R$ %10.2f               ║\n", valorBrutoFinal)
	fmt.Printf("║ Valor líquido final (RF):  R$ %10.2f               ║\n", valorLiquidoRFFinal)
	fmt.Printf("║ Valor líquido final (FIIs):R$ %10.2f               ║\n", valorLiquidoFIIsFinal)
	fmt.Printf("║ Lucro bruto total (RF):    R$ %10.2f               ║\n", lucroBrutoTotal)
	fmt.Printf("║ Lucro Líquido RF:          R$ %10.2f               ║\n", lucroLiquidoAcumulado)
	fmt.Printf("║ Lucro Líquido FIIs:        R$ %10.2f               ║\n", lucroLiquidoFIIsAcumulado)
	fmt.Printf("║ Lucro Total Líquido:       R$ %10.2f               ║\n", lucroMesLiquidoTotalAcumulado)
	fmt.Printf("║ Lucros retirados:          R$ %10.2f               ║\n", lucrosRetiradosTotal)
	fmt.Println("╚══════════════════════════════════════════════════════╝\n")
}

func mostrarResumoMesAtual(dados Dados) {
	hoje := time.Now()
	anoAtual := fmt.Sprintf("%04d", hoje.Year())
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))
	anos := ordenarChaves(dados.Anos)
	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := ordenarChaves(mesesMap)
		for _, mes := range meses {
			if ano == anoAtual && mes == mesAtual {
				m := mesesMap[mes]
				lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
				impostos := m.ValorBrutoRF - m.ValorLiquidoRF
				lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
				lucroLiquidoFIIs := m.LucroLiquidoFIIs
				lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
				fmt.Println("╔══════════════════════════════════════════════════════╗")
				fmt.Printf("║ Mês: %s/%s\n", nomeMes(mes), ano)
				fmt.Println("║  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)")
				fmt.Println("╠══════════════════════════════════════════════════════╣")
				fmt.Printf("║  Aporte Total:         R$ %10.2f                 ║\n", m.AporteRF+m.AporteFIIs)
				fmt.Printf("║  Aporte RF:            R$ %10.2f                 ║\n", m.AporteRF)
				fmt.Printf("║  FIIs:                 R$ %10.2f                 ║\n", m.AporteFIIs)
				fmt.Printf("║  Saída:                R$ %10.2f                 ║\n", m.Saida)
				fmt.Printf("║  Lucro Retirado:       R$ %10.2f                 ║\n", m.LucroRetirado)
				fmt.Printf("║  Bruto RF:             R$ %10.2f                 ║\n", m.ValorBrutoRF)
				fmt.Printf("║  Líquido RF:           R$ %10.2f                 ║\n", m.ValorLiquidoRF)
				fmt.Printf("║  Líquido FIIs:         R$ %10.2f                 ║\n", m.ValorLiquidoFIIs)
				fmt.Printf("║  Lucro Mês Bruto:      R$ %10.2f                 ║\n", lucroMesBruto)
				fmt.Printf("║  Lucro Líquido RF:     R$ %10.2f                 ║\n", lucroMesLiquidoRF)
				fmt.Printf("║  Lucro Líquido FIIs:   R$ %10.2f                 ║\n", lucroLiquidoFIIs)
				fmt.Printf("║  Lucro Mês Líquido:    R$ %10.2f                 ║\n", lucroMesLiquidoTotal)
				fmt.Println("╚══════════════════════════════════════════════════════╝\n")
				return
			}
			saldoAnterior = mesesMap[mes].ValorBrutoRF
		}
	}
}

func printTelaUnificada(dados Dados) {
	clearTerminal()
	time.Sleep(300 * time.Millisecond)
	// Preparar strings de cada seção
	resumoTotal := getResumoTotalAcumuladoStr(dados)
	resumoMes := getResumoMesAtualStr(dados)
	menu := getMenuPrincipalStr()

	// Descobrir o maior comprimento de linha
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

	// Bordas
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

func getResumoTotalAcumuladoStr(dados Dados) string {
	anos := ordenarChaves(dados.Anos)
	if len(anos) == 0 {
		return "Nenhum dado disponível ainda."
	}
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	valorBrutoFinal := 0.0
	valorLiquidoRFFinal := 0.0
	valorLiquidoFIIsFinal := 0.0
	lucrosRetiradosTotal := 0.0
	lucroLiquidoAcumulado := 0.0
	lucroLiquidoFIIsAcumulado := 0.0
	lucroMesLiquidoTotalAcumulado := 0.0
	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := ordenarChaves(mesesMap)
		for _, mes := range meses {
			m := mesesMap[mes]
			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
			lucroLiquidoFIIs := m.LucroLiquidoFIIs
			lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
			lucroValido := lucroMesBruto > impostos
			if lucroValido {
				aporteRFSoFar += m.AporteRF
				aporteFIIsSoFar += m.AporteFIIs
				saidaSoFar += m.Saida
				lucrosRetiradosTotal += m.LucroRetirado
				valorBrutoFinal = m.ValorBrutoRF
				valorLiquidoRFFinal = m.ValorLiquidoRF
				valorLiquidoFIIsFinal = m.ValorLiquidoFIIs
				lucroLiquidoAcumulado += lucroMesLiquidoRF
				lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
				lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
				saldoAnterior = m.ValorBrutoRF
			}
		}
	}
	totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
	totalAportadoLiquido := totalAportadoBruto - saidaSoFar
	lucroBrutoTotal := valorBrutoFinal - totalAportadoLiquido
	return fmt.Sprintf(`--- Resumo Total Acumulado ---
Total aportado bruto: R$ %.2f
Total aportado líquido: R$ %.2f
Valor bruto final (RF): R$ %.2f
Valor líquido final (RF): R$ %.2f
Valor líquido final (FIIs): R$ %.2f
Lucro bruto total (RF): R$ %.2f
Lucro Líquido RF: R$ %.2f
Lucro Líquido FIIs: R$ %.2f
Lucro Total Líquido (RF + FIIs): R$ %.2f
Lucros retirados: R$ %.2f`,
		totalAportadoBruto, totalAportadoLiquido, valorBrutoFinal, valorLiquidoRFFinal, valorLiquidoFIIsFinal, lucroBrutoTotal, lucroLiquidoAcumulado, lucroLiquidoFIIsAcumulado, lucroMesLiquidoTotalAcumulado, lucrosRetiradosTotal)
}

func getResumoMesAtualStr(dados Dados) string {
	hoje := time.Now()
	anoAtual := fmt.Sprintf("%04d", hoje.Year())
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))
	anos := ordenarChaves(dados.Anos)
	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := ordenarChaves(mesesMap)
		for _, mes := range meses {
			if ano == anoAtual && mes == mesAtual {
				m := mesesMap[mes]
				lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
				impostos := m.ValorBrutoRF - m.ValorLiquidoRF
				lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
				lucroLiquidoFIIs := m.LucroLiquidoFIIs
				lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
				titulo := fmt.Sprintf("Mês: %s/%s", nomeMes(mes), ano)
				return fmt.Sprintf(`%s
  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)
---------------------------------------
  Aporte Total:         R$ %.2f
  Aporte RF:            R$ %.2f
  FIIs:                 R$ %.2f
  Saída:                R$ %.2f
  Lucro Retirado:       R$ %.2f
  Bruto RF:             R$ %.2f
  Líquido RF:           R$ %.2f
  Líquido FIIs:         R$ %.2f
  Lucro Mês Bruto:      R$ %.2f
  Lucro Líquido RF:     R$ %.2f
  Lucro Líquido FIIs:   R$ %.2f
  Lucro Mês Líquido:    R$ %.2f
---------------------------------------`,
					titulo,
					m.AporteRF+m.AporteFIIs, m.AporteRF, m.AporteFIIs, m.Saida, m.LucroRetirado, m.ValorBrutoRF, m.ValorLiquidoRF, m.ValorLiquidoFIIs, lucroMesBruto, lucroMesLiquidoRF, lucroLiquidoFIIs, lucroMesLiquidoTotal)
			}
			saldoAnterior = mesesMap[mes].ValorBrutoRF
		}
	}
	return "Mês atual não possui dados."
}

func getMenuPrincipalStr() string {
	return `--- MENU PRINCIPAL ---
1. Ver resumo completo (visualização vertical)
2. Ver resumo completo (tabela horizontal)
3. Adicionar/editar mês
4. Sair do programa`
}

func printMenuPrincipalSozinho() {
	clearTerminal()
	time.Sleep(300 * time.Millisecond)
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║ --- MENU PRINCIPAL ---                             ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Ver resumo completo (visualização vertical)      ║")
	fmt.Println("║ 2. Ver resumo completo (tabela horizontal)          ║")
	fmt.Println("║ 3. Adicionar/editar mês                             ║")
	fmt.Println("║ 4. Voltar para o menu inicial                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func menu() {
	dados := carregarDados()
	scanner := bufio.NewScanner(os.Stdin)

	// Mostrar tudo em uma caixa só na tela inicial
	printTelaUnificada(dados)

	inMenuInicial := true

	for {
		fmt.Print("Escolha uma opção: ")
		scanner.Scan()
		opcao := scanner.Text()

		if inMenuInicial && opcao == "4" {
			fmt.Println("Saindo...")
			return
		}

		if !inMenuInicial && opcao == "4" {
			// Voltar ao menu inicial (com resumos)
			printTelaUnificada(dados)
			inMenuInicial = true
			continue
		}

		switch opcao {
		case "1":
			ano := selecionarAno(dados, scanner)
			if ano != "" {
				mostrarResumoAno(dados, ano, false)
				fmt.Print("\nPressione Enter para voltar ao menu...")
				scanner.Scan()
			}
			printMenuPrincipalSozinho()
			inMenuInicial = false
		case "2":
			ano := selecionarAno(dados, scanner)
			if ano != "" {
				mostrarResumoAno(dados, ano, true)
				fmt.Print("\nPressione Enter para voltar ao menu...")
				scanner.Scan()
			}
			printMenuPrincipalSozinho()
			inMenuInicial = false
		case "3":
			adicionarOuEditarMes(&dados, scanner)
			salvarDados(dados)
			// Atualizar tela unificada após edição
			printTelaUnificada(dados)
			inMenuInicial = true
		default:
			fmt.Println("Opção inválida!")
			printMenuPrincipalSozinho()
			inMenuInicial = false
		}
	}
}

func selecionarAno(dados Dados, scanner *bufio.Scanner) string {
	if len(dados.Anos) == 0 {
		fmt.Println("Nenhum dado disponível ainda.")
		return ""
	}

	anos := ordenarChaves(dados.Anos)

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

func mostrarResumoAno(dados Dados, ano string, horizontal bool) {
	clearTerminal()
	time.Sleep(300 * time.Millisecond)
	mesesMap, ok := dados.Anos[ano]
	if !ok || len(mesesMap) == 0 {
		fmt.Printf("Não há dados para o ano %s.\n", ano)
		return
	}

	meses := ordenarChaves(mesesMap)

	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	valorBrutoFinal := 0.0
	valorLiquidoRFFinal := 0.0
	valorLiquidoFIIsFinal := 0.0
	lucrosRetiradosTotal := 0.0
	lucroLiquidoAcumulado := 0.0
	lucroLiquidoFIIsAcumulado := 0.0
	lucroMesLiquidoTotalAcumulado := 0.0
	saldoAnterior := 0.0

	hoje := time.Now()
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))
	anoAtual := fmt.Sprintf("%04d", hoje.Year())

	if horizontal {
		fmt.Printf("\n📌 Resumo dos aportes e saldos mensais - Ano %s (Tabela Horizontal)\n", ano)
		fmt.Println("\n| Mês      | Aporte Total | Aporte RF | FIIs | Saída | Lucro Ret. | Bruto RF | Líquido RF | Líquido FIIs | Lucro Mês Bruto | Lucro Líquido RF | Lucro Líquido FIIs | Lucro Mês Líquido |")
		fmt.Println("|----------|--------------|-----------|------|--------|------------|----------|------------|--------------|-----------------|------------------|--------------------|-------------------|")
	} else {
		fmt.Printf("\n📌 Resumo dos aportes e saldos mensais - Ano %s (Visualização Vertical)\n", ano)
	}

	for _, mes := range meses {
		m := mesesMap[mes]

		lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
		impostos := m.ValorBrutoRF - m.ValorLiquidoRF
		lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
		lucroLiquidoFIIs := m.LucroLiquidoFIIs
		lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs

		isMesAtual := (ano == anoAtual && mes == mesAtual)

		if horizontal {
			fmt.Printf("| %-8s | R$ %10.2f | R$ %7.2f | R$%4.2f | R$%6.2f | R$ %9.2f | R$ %8.2f | R$ %10.2f | R$ %12.2f | R$ %14.2f | R$ %16.2f | R$ %18.2f | R$ %17.2f |\n",
				nomeMes(mes), m.AporteRF+m.AporteFIIs, m.AporteRF, m.AporteFIIs, m.Saida, m.LucroRetirado,
				m.ValorBrutoRF, m.ValorLiquidoRF, m.ValorLiquidoFIIs,
				lucroMesBruto, lucroMesLiquidoRF, lucroLiquidoFIIs, lucroMesLiquidoTotal)
		} else {
			fmt.Printf("\nMês: %s/%s\n", nomeMes(mes), ano)
			if isMesAtual {
				fmt.Println("  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)")
			}

			impostoValido := impostos > 0
			if lucroMesBruto > impostos && impostoValido {
				fmt.Println("  ✅ Agora os lucros já cobrem os impostos!")
			}

			fmt.Println("---------------------------------------")

			fmt.Printf("  Aporte Total:         R$ %.2f\n", m.AporteRF+m.AporteFIIs)
			fmt.Printf("  Aporte RF:            R$ %.2f\n", m.AporteRF)
			fmt.Printf("  FIIs:                 R$ %.2f\n", m.AporteFIIs)
			fmt.Printf("  Saída:                R$ %.2f\n", m.Saida)
			fmt.Printf("  Lucro Retirado:       R$ %.2f\n", m.LucroRetirado)
			fmt.Printf("  Bruto RF:             R$ %.2f\n", m.ValorBrutoRF)
			fmt.Printf("  Líquido RF:           R$ %.2f\n", m.ValorLiquidoRF)
			fmt.Printf("  Líquido FIIs:         R$ %.2f\n", m.ValorLiquidoFIIs)
			fmt.Printf("  Lucro Mês Bruto:      R$ %.2f\n", lucroMesBruto)
			fmt.Printf("  Lucro Líquido RF:     R$ %.2f\n", lucroMesLiquidoRF)
			fmt.Printf("  Lucro Líquido FIIs:   R$ %.2f\n", lucroLiquidoFIIs)
			fmt.Printf("  Lucro Mês Líquido:    R$ %.2f\n", lucroMesLiquidoTotal)

			fmt.Println("---------------------------------------")
		}

		lucroValido := lucroMesBruto > impostos

		if lucroValido {
			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += m.AporteFIIs
			saidaSoFar += m.Saida
			lucrosRetiradosTotal += m.LucroRetirado

			valorBrutoFinal = m.ValorBrutoRF
			valorLiquidoRFFinal = m.ValorLiquidoRF
			valorLiquidoFIIsFinal = m.ValorLiquidoFIIs

			lucroLiquidoAcumulado += lucroMesLiquidoRF
			lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
			lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
			saldoAnterior = m.ValorBrutoRF
		}
	}

	totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
	totalAportadoLiquido := totalAportadoBruto - saidaSoFar
	lucroBrutoTotal := valorBrutoFinal - totalAportadoLiquido
	lucroLiquidoTotal := lucroLiquidoAcumulado
	lucroLiquidoFIIsTotal := lucroLiquidoFIIsAcumulado
	lucroMesLiquidoTotalAno := lucroMesLiquidoTotalAcumulado

	fmt.Println()
	fmt.Println("--- Resumo Total do Ano ---")
	fmt.Printf("Total aportado bruto: R$ %.2f\n", totalAportadoBruto)
	fmt.Printf("Total aportado líquido: R$ %.2f\n", totalAportadoLiquido)
	fmt.Printf("Valor bruto final (RF): R$ %.2f\n", valorBrutoFinal)
	fmt.Printf("Valor líquido final (RF): R$ %.2f\n", valorLiquidoRFFinal)
	fmt.Printf("Valor líquido final (FIIs): R$ %.2f\n", valorLiquidoFIIsFinal)
	fmt.Printf("Lucro bruto total (RF): R$ %.2f\n", lucroBrutoTotal)
	fmt.Printf("Lucro Líquido RF: R$ %.2f\n", lucroLiquidoTotal)
	fmt.Printf("Lucro Líquido FIIs: R$ %.2f\n", lucroLiquidoFIIsTotal)
	fmt.Printf("Lucro Total Líquido (RF + FIIs): R$ %.2f\n", lucroMesLiquidoTotalAno)
	fmt.Printf("Lucros retirados: R$ %.2f\n", lucrosRetiradosTotal)
}

// Caixa para inputs
func inputBox(prompt string, scanner *bufio.Scanner) string {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Printf("║ %-48s ║\n", prompt)
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Print("→ ")
	scanner.Scan()
	return scanner.Text()
}

func mostrarResumoTodosAnos(dados Dados) {
	anos := ordenarChaves(dados.Anos)
	if len(anos) == 0 {
		fmt.Println("Nenhum dado disponível ainda.")
		return
	}
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := ordenarChaves(mesesMap)
		aporteRFSoFar := 0.0
		aporteFIIsSoFar := 0.0
		saidaSoFar := 0.0
		valorBrutoFinal := 0.0
		valorLiquidoRFFinal := 0.0
		valorLiquidoFIIsFinal := 0.0
		lucrosRetiradosTotal := 0.0
		lucroLiquidoAcumulado := 0.0
		lucroLiquidoFIIsAcumulado := 0.0
		lucroMesLiquidoTotalAcumulado := 0.0
		saldoAnterior := 0.0
		for _, mes := range meses {
			m := mesesMap[mes]
			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
			lucroLiquidoFIIs := m.LucroLiquidoFIIs
			lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
			lucroValido := lucroMesBruto > impostos
			if lucroValido {
				aporteRFSoFar += m.AporteRF
				aporteFIIsSoFar += m.AporteFIIs
				saidaSoFar += m.Saida
				lucrosRetiradosTotal += m.LucroRetirado
				valorBrutoFinal = m.ValorBrutoRF
				valorLiquidoRFFinal = m.ValorLiquidoRF
				valorLiquidoFIIsFinal = m.ValorLiquidoFIIs
				lucroLiquidoAcumulado += lucroMesLiquidoRF
				lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
				lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
				saldoAnterior = m.ValorBrutoRF
			}
		}
		totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
		totalAportadoLiquido := totalAportadoBruto - saidaSoFar
		lucroBrutoTotal := valorBrutoFinal - totalAportadoLiquido
		// Caixa bonita para cada ano
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Printf("║           RESUMO TOTAL DO ANO %s                        ║\n", ano)
		fmt.Println("╠══════════════════════════════════════════════════════╣")
		fmt.Printf("║ Total aportado bruto:      R$ %10.2f               ║\n", totalAportadoBruto)
		fmt.Printf("║ Total aportado líquido:    R$ %10.2f               ║\n", totalAportadoLiquido)
		fmt.Printf("║ Valor bruto final (RF):    R$ %10.2f               ║\n", valorBrutoFinal)
		fmt.Printf("║ Valor líquido final (RF):  R$ %10.2f               ║\n", valorLiquidoRFFinal)
		fmt.Printf("║ Valor líquido final (FIIs):R$ %10.2f               ║\n", valorLiquidoFIIsFinal)
		fmt.Printf("║ Lucro bruto total (RF):    R$ %10.2f               ║\n", lucroBrutoTotal)
		fmt.Printf("║ Lucro Líquido RF:          R$ %10.2f               ║\n", lucroLiquidoAcumulado)
		fmt.Printf("║ Lucro Líquido FIIs:        R$ %10.2f               ║\n", lucroLiquidoFIIsAcumulado)
		fmt.Printf("║ Lucro Total Líquido:       R$ %10.2f               ║\n", lucroMesLiquidoTotalAcumulado)
		fmt.Printf("║ Lucros retirados:          R$ %10.2f               ║\n", lucrosRetiradosTotal)
		fmt.Println("╚══════════════════════════════════════════════════════╝\n")
	}
}

func adicionarOuEditarMes(dados *Dados, scanner *bufio.Scanner) {
	ano := inputBox("Digite o ano(YYYY):", scanner)
	mes := inputBox("Digite o mês(MM):", scanner)

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}

	m := dados.Anos[ano][mes]
	if m != (Mes{}) {
		for {
			fmt.Println("\n--- EDITAR CAMPOS ---")
			fmt.Printf("1. Aporte RF (atual: %.2f)\n", m.AporteRF)
			fmt.Printf("2. Aporte FIIs (atual: %.2f)\n", m.AporteFIIs)
			fmt.Printf("3. Saída (atual: %.2f)\n", m.Saida)
			fmt.Printf("4. Valor Bruto RF (atual: %.2f)\n", m.ValorBrutoRF)
			fmt.Printf("5. Valor Líquido RF (atual: %.2f)\n", m.ValorLiquidoRF)
			fmt.Printf("6. Valor Líquido FIIs (atual: %.2f)\n", m.ValorLiquidoFIIs)
			fmt.Printf("7. Lucro Retirado (atual: %.2f)\n", m.LucroRetirado)
			fmt.Printf("8. Lucro Líquido FIIs (atual: %.2f)\n", m.LucroLiquidoFIIs)
			fmt.Println("0. Sair da edição")
			opcao := inputBox("Escolha:", scanner)

			switch opcao {
			case "1":
				valor := inputBox("Novo valor:", scanner)
				m.AporteRF, _ = strconv.ParseFloat(valor, 64)
			case "2":
				valor := inputBox("Novo valor:", scanner)
				m.AporteFIIs, _ = strconv.ParseFloat(valor, 64)
			case "3":
				valor := inputBox("Novo valor:", scanner)
				m.Saida, _ = strconv.ParseFloat(valor, 64)
			case "4":
				valor := inputBox("Novo valor:", scanner)
				m.ValorBrutoRF, _ = strconv.ParseFloat(valor, 64)
			case "5":
				valor := inputBox("Novo valor:", scanner)
				m.ValorLiquidoRF, _ = strconv.ParseFloat(valor, 64)
			case "6":
				valor := inputBox("Novo valor:", scanner)
				m.ValorLiquidoFIIs, _ = strconv.ParseFloat(valor, 64)
			case "7":
				valor := inputBox("Novo valor:", scanner)
				m.LucroRetirado, _ = strconv.ParseFloat(valor, 64)
			case "8":
				valor := inputBox("Novo valor:", scanner)
				m.LucroLiquidoFIIs, _ = strconv.ParseFloat(valor, 64)
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

	aporteRF, _ := strconv.ParseFloat(inputBox("Digite o aporte na Renda Fixa: R$", scanner), 64)
	aporteFIIs, _ := strconv.ParseFloat(inputBox("Digite o aporte em FIIs: R$", scanner), 64)
	saida, _ := strconv.ParseFloat(inputBox("Digite a saída (retirada) do mês: R$", scanner), 64)
	valorBrutoRF, _ := strconv.ParseFloat(inputBox("Digite o valor bruto da Renda Fixa: R$", scanner), 64)
	valorLiquidoRF, _ := strconv.ParseFloat(inputBox("Digite o valor líquido da Renda Fixa: R$", scanner), 64)
	valorLiquidoFIIs, _ := strconv.ParseFloat(inputBox("Digite o valor líquido dos FIIs: R$", scanner), 64)
	lucroRetirado, _ := strconv.ParseFloat(inputBox("Digite o valor de lucro retirado: R$", scanner), 64)
	lucroLiquidoFIIs, _ := strconv.ParseFloat(inputBox("Digite o lucro líquido dos FIIs: R$", scanner), 64)

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

func clearTerminal() {
	cmd := exec.Command("clear")
	if _, ok := os.LookupEnv("OS"); ok {
		cmd = exec.Command("cls") // para Windows
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	menu()
}
