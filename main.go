package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Mes struct {
	AporteRF         float64 `json:"aporte_rf"`
	AporteFIIs       float64 `json:"aporte_fiis"`
	Saida            float64 `json:"saida"`
	ValorBrutoRF     float64 `json:"valor_bruto_rf"`
	ValorLiquidoRF   float64 `json:"valor_liquido_rf"`
	ValorLiquidoFIIs float64 `json:"valor_liquido_fiis"`
	LucroRetirado    float64 `json:"lucro_retirado"`
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

func menu() {
	dados := carregarDados()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU PRINCIPAL ---")
		fmt.Println("1. Ver resumo completo (tabela horizontal)")
		fmt.Println("2. Ver resumo completo (visualização vertical)")
		fmt.Println("3. Adicionar/editar mês")
		fmt.Println("4. Sair")
		fmt.Print("Escolha uma opção: ")
		scanner.Scan()
		opcao := scanner.Text()

		switch opcao {
		case "1":
			mostrarResumo(dados, true)
		case "2":
			mostrarResumo(dados, false)
		case "3":
			adicionarOuEditarMes(&dados, scanner)
			salvarDados(dados)
		case "4":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}

func mostrarResumo(dados Dados, horizontal bool) {
	anos := ordenarChaves(dados.Anos)
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	valorBrutoFinal := 0.0
	valorLiquidoRFFinal := 0.0
	valorLiquidoFIIsFinal := 0.0
	lucrosRetiradosTotal := 0.0
	lucroLiquidoAcumulado := 0.0
	saldoAnterior := 0.0

	if horizontal {
		fmt.Println("\n📌 Resumo dos aportes e saldos mensais (Tabela Horizontal)")
		fmt.Println("\n| Mês      | Aporte Total | Aporte RF | FIIs | Saída | Lucro Ret. | Bruto RF | Líquido RF | Líquido FIIs | Lucro Mês Bruto | Lucro Mês Líquido |")
		fmt.Println("|----------|--------------|-----------|------|--------|-------------|-----------|-------------|----------------|------------------|---------------------|")
	}

	for _, ano := range anos {
		meses := ordenarChaves(dados.Anos[ano])
		for _, mes := range meses {
			m := dados.Anos[ano][mes]

			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += m.AporteFIIs
			saidaSoFar += m.Saida
			lucrosRetiradosTotal += m.LucroRetirado

			valorBrutoFinal = m.ValorBrutoRF
			valorLiquidoRFFinal = m.ValorLiquidoRF
			valorLiquidoFIIsFinal = m.ValorLiquidoFIIs

			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquido := lucroMesBruto - impostos - m.LucroRetirado
			saldoAnterior = m.ValorBrutoRF

			lucroLiquidoAcumulado += lucroMesLiquido + m.LucroRetirado // soma real do lucro do mês, antes de retirada

			if horizontal {
				fmt.Printf("| %-8s | R$ %10.2f | R$ %7.2f | R$%4.2f | R$%6.2f | R$ %9.2f | R$ %8.2f | R$ %10.2f | R$ %12.2f | R$ %14.2f | R$ %19.2f |\n",
					nomeMes(mes), m.AporteRF+m.AporteFIIs, m.AporteRF, m.AporteFIIs, m.Saida, m.LucroRetirado,
					m.ValorBrutoRF, m.ValorLiquidoRF, m.ValorLiquidoFIIs,
					lucroMesBruto, lucroMesLiquido)
			} else {
				fmt.Printf("\nMês: %s/%s\n", nomeMes(mes), ano)
				fmt.Printf("  Aporte Total:      R$ %.2f\n", m.AporteRF+m.AporteFIIs)
				fmt.Printf("  Aporte RF:         R$ %.2f\n", m.AporteRF)
				fmt.Printf("  FIIs:              R$ %.2f\n", m.AporteFIIs)
				fmt.Printf("  Saída:             R$ %.2f\n", m.Saida)
				fmt.Printf("  Lucro Retirado:    R$ %.2f\n", m.LucroRetirado)
				fmt.Printf("  Bruto RF:          R$ %.2f\n", m.ValorBrutoRF)
				fmt.Printf("  Líquido RF:        R$ %.2f\n", m.ValorLiquidoRF)
				fmt.Printf("  Líquido FIIs:      R$ %.2f\n", m.ValorLiquidoFIIs)
				fmt.Printf("  Lucro Mês Bruto:   R$ %.2f\n", lucroMesBruto)
				fmt.Printf("  Lucro Mês Líquido: R$ %.2f\n", lucroMesLiquido)
				fmt.Println("---------------------------------------")
			}
		}
	}

	totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
	totalAportadoLiquido := totalAportadoBruto - saidaSoFar
	lucroBrutoTotal := valorBrutoFinal - totalAportadoLiquido
	lucroLiquidoTotal := lucroLiquidoAcumulado - lucrosRetiradosTotal

	if horizontal {
		fmt.Println()
	} else {
		fmt.Println("\n--- Mini Resumo Total ---")
	}
	fmt.Printf("Total aportado bruto: R$ %.2f\n", totalAportadoBruto)
	fmt.Printf("Total aportado líquido: R$ %.2f\n", totalAportadoLiquido)
	fmt.Printf("Valor bruto final (RF): R$ %.2f\n", valorBrutoFinal)
	fmt.Printf("Valor líquido final (RF): R$ %.2f\n", valorLiquidoRFFinal)
	fmt.Printf("Valor líquido final (FIIs): R$ %.2f\n", valorLiquidoFIIsFinal)
	fmt.Printf("Lucro bruto total (RF): R$ %.2f\n", lucroBrutoTotal)
	fmt.Printf("Lucro líquido total: R$ %.2f\n", lucroLiquidoTotal)
	fmt.Printf("Lucros retirados: R$ %.2f\n", lucrosRetiradosTotal)
}

func adicionarOuEditarMes(dados *Dados, scanner *bufio.Scanner) {
	fmt.Print("Digite o ano(YYYY): ")
	scanner.Scan()
	ano := scanner.Text()

	fmt.Print("Digite o mês(MM): ")
	scanner.Scan()
	mes := scanner.Text()

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
			fmt.Println("0. Sair da edição")
			fmt.Print("Escolha: ")
			scanner.Scan()
			opcao := scanner.Text()

			switch opcao {
			case "1":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.AporteRF, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "2":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.AporteFIIs, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "3":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.Saida, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "4":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.ValorBrutoRF, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "5":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.ValorLiquidoRF, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "6":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.ValorLiquidoFIIs, _ = strconv.ParseFloat(scanner.Text(), 64)
			case "7":
				fmt.Print("Novo valor: ")
				scanner.Scan()
				m.LucroRetirado, _ = strconv.ParseFloat(scanner.Text(), 64)
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

	fmt.Print("Digite o aporte na Renda Fixa: R$ ")
	scanner.Scan()
	aporteRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o aporte em FIIs: R$ ")
	scanner.Scan()
	aporteFIIs, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite a saída (retirada) do mês: R$ ")
	scanner.Scan()
	saida, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor bruto da Renda Fixa: R$ ")
	scanner.Scan()
	valorBrutoRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor líquido da Renda Fixa: R$ ")
	scanner.Scan()
	valorLiquidoRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor líquido dos FIIs: R$ ")
	scanner.Scan()
	valorLiquidoFIIs, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor de lucro retirado: R$ ")
	scanner.Scan()
	lucroRetirado, _ := strconv.ParseFloat(scanner.Text(), 64)

	dados.Anos[ano][mes] = Mes{
		AporteRF:         aporteRF,
		AporteFIIs:       aporteFIIs,
		Saida:            saida,
		ValorBrutoRF:     valorBrutoRF,
		ValorLiquidoRF:   valorLiquidoRF,
		ValorLiquidoFIIs: valorLiquidoFIIs,
		LucroRetirado:    lucroRetirado,
	}

	fmt.Println("Dados adicionados com sucesso!")
}

func main() {
	menu()
}
