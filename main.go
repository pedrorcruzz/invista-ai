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
	Aporte       float64 `json:"aporte"`
	ValorBruto   float64 `json:"valor_bruto"`
	ValorLiquido float64 `json:"valor_liquido"`
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

func menu() {
	dados := carregarDados()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU PRINCIPAL ---")
		fmt.Println("1. Ver resumo completo")
		fmt.Println("2. Adicionar/editar mês")
		fmt.Println("3. Sair")
		fmt.Print("Escolha uma opção: ")
		scanner.Scan()
		opcao := scanner.Text()

		switch opcao {
		case "1":
			mostrarResumo(dados)
		case "2":
			adicionarOuEditarMes(&dados, scanner)
			salvarDados(dados)
		case "3":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
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

func mostrarResumo(dados Dados) {
	fmt.Println("\n📌 Resumo dos aportes e saldos mensais")

	totalAportado := 0.0
	valorBrutoAcumulado := 0.0
	valorLiquidoAcumulado := 0.0
	ultimoAno := ""
	ultimoMes := ""
	var mesAtual Mes

	for ano, meses := range dados.Anos {
		for mes := range meses {
			if ano > ultimoAno || (ano == ultimoAno && mes > ultimoMes) {
				ultimoAno = ano
				ultimoMes = mes
				mesAtual = meses[mes]
			}
		}
	}

	fmt.Printf("\n🗓️  MÊS ATUAL: %s/%s\n", nomeMes(ultimoMes), ultimoAno)
	fmt.Printf("Aporte: R$ %.2f\n", mesAtual.Aporte)
	fmt.Printf("Valor Bruto: R$ %.2f\n", mesAtual.ValorBruto)
	fmt.Printf("Valor Líquido: R$ %.2f\n", mesAtual.ValorLiquido)
	fmt.Printf("Lucro Bruto: R$ %.2f\n", mesAtual.ValorBruto-mesAtual.Aporte)
	fmt.Printf("Lucro Líquido: R$ %.2f\n", mesAtual.ValorLiquido-mesAtual.Aporte)

	fmt.Println("\n| Mês      | Aporte     | Valor Bruto | Valor Líquido | Lucro Bruto | Lucro Líquido |")
	fmt.Println("|----------|------------|-------------|----------------|--------------|----------------|")

	anos := ordenarChaves(dados.Anos)
	for _, ano := range anos {
		meses := ordenarChaves(dados.Anos[ano])
		for _, mes := range meses {
			m := dados.Anos[ano][mes]

			totalAportado += m.Aporte
			valorBrutoAcumulado = m.ValorBruto
			valorLiquidoAcumulado = m.ValorLiquido

			lucroBruto := valorBrutoAcumulado - totalAportado
			lucroLiquido := valorLiquidoAcumulado - totalAportado

			fmt.Printf("| %-8s | R$ %8.2f | R$ %9.2f | R$ %12.2f | R$ %10.2f | R$ %12.2f |\n",
				nomeMes(mes), m.Aporte, m.ValorBruto, m.ValorLiquido, lucroBruto, lucroLiquido)
		}
	}

	fmt.Printf("\nTotal aportado: R$ %.2f\n", totalAportado)
	fmt.Printf("Valor bruto final: R$ %.2f\n", valorBrutoAcumulado)
	fmt.Printf("Valor líquido final: R$ %.2f\n", valorLiquidoAcumulado)
	fmt.Printf("Lucro bruto total: R$ %.2f\n", valorBrutoAcumulado-totalAportado)
	fmt.Printf("Lucro líquido total: R$ %.2f\n", valorLiquidoAcumulado-totalAportado)
}

func ordenarChaves[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func adicionarOuEditarMes(dados *Dados, scanner *bufio.Scanner) {
	fmt.Print("Digite o ano (ex: 2025): ")
	scanner.Scan()
	ano := scanner.Text()

	fmt.Print("Digite o mês (ex: 05): ")
	scanner.Scan()
	mes := scanner.Text()

	fmt.Print("Digite o aporte: R$ ")
	scanner.Scan()
	aporte, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor bruto atual: R$ ")
	scanner.Scan()
	valorBruto, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor líquido atual: R$ ")
	scanner.Scan()
	valorLiquido, _ := strconv.ParseFloat(scanner.Text(), 64)

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}

	dados.Anos[ano][mes] = Mes{
		Aporte:       aporte,
		ValorBruto:   valorBruto,
		ValorLiquido: valorLiquido,
	}

	fmt.Println("Aporte e valores salvos com sucesso!")
}

func main() {
	menu()
}
