package internal

import (
	"bufio"
	"fmt"
	"sort"
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
	fmt.Println("║ 2. Renda Fixa                                       ║")
	fmt.Println("║ 3. FIIs                                             ║")
	fmt.Println("║ 4. Gestor Inteligente de Gastos                     ║")
	fmt.Println("║ 5. Ajustar valor da carteira                        ║")
	fmt.Println("║ 6. Ajuste Preço Médio                               ║")
	fmt.Println("║ 7. Voltar ao menu principal                         ║")
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
2. Renda Fixa
3. FIIs
4. Gestor Inteligente de Gastos
5. Retirar Lucro
6. Sair do programa`
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
	if m.AporteRF != 0 || m.Saida != 0 || m.ValorBrutoRF != 0 || m.ValorLiquidoRF != 0 || m.LucroRetirado != 0 || len(m.FIIs) > 0 {
		for {
			// Menu de edição em caixinha
			lines := []string{
				"EDITAR CAMPOS",
				"",
				"0. Sair da edição",
				"",
				fmt.Sprintf("1. Aporte RF (atual: %s)", FormatFloatBR(m.AporteRF)),
				fmt.Sprintf("2. Saída (atual: %s)", FormatFloatBR(m.Saida)),
				fmt.Sprintf("3. Valor Bruto RF (atual: %s)", FormatFloatBR(m.ValorBrutoRF)),
				fmt.Sprintf("4. Valor Líquido RF (atual: %s)", FormatFloatBR(m.ValorLiquidoRF)),
				fmt.Sprintf("5. Lucro Retirado (atual: %s)", FormatFloatBR(m.LucroRetirado)),
			}
			PrintCaixa(lines)
			opcao := InputBox("Escolha:", scanner)
			switch opcao {
			case "1":
				valor := InputBox("Novo valor:", scanner)
				m.AporteRF, _ = ParseFloatBR(valor)
			case "2":
				valor := InputBox("Novo valor:", scanner)
				m.Saida, _ = ParseFloatBR(valor)
			case "3":
				valor := InputBox("Novo valor:", scanner)
				m.ValorBrutoRF, _ = ParseFloatBR(valor)
			case "4":
				valor := InputBox("Novo valor:", scanner)
				m.ValorLiquidoRF, _ = ParseFloatBR(valor)
			case "5":
				valor := InputBox("Novo valor:", scanner)
				m.LucroRetirado, _ = ParseFloatBR(valor)
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
	saida, _ := ParseFloatBR(InputBox("Digite a saída (retirada) do mês: R$", scanner))
	valorBrutoRF, _ := ParseFloatBR(InputBox("Digite o valor bruto da Renda Fixa: R$", scanner))
	valorLiquidoRF, _ := ParseFloatBR(InputBox("Digite o valor líquido da Renda Fixa: R$", scanner))
	lucroRetirado, _ := ParseFloatBR(InputBox("Digite o valor de lucro retirado: R$", scanner))
	dados.Anos[ano][mes] = Mes{
		AporteRF:       aporteRF,
		Saida:          saida,
		ValorBrutoRF:   valorBrutoRF,
		ValorLiquidoRF: valorLiquidoRF,
		LucroRetirado:  lucroRetirado,
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

func GerenciarRendaFixa(dados *Dados, scanner *bufio.Scanner) {
	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                 RENDA FIXA                          ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Adicionar/editar mês                             ║")
	fmt.Println("║ 2. Voltar ao menu principal                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	opcao := InputBox("Escolha uma opção:", scanner)
	switch opcao {
	case "1":
		AdicionarOuEditarMes(dados, scanner)
	case "2":
		return
	default:
		fmt.Println("Opção inválida.")
		Pause(2000)
	}
}

func GerenciarFIIs(dados *Dados, scanner *bufio.Scanner) {
	for {
		// Calcular o valor total investido em FIIs (soma de todos os meses/anos)
		totalInvestido := 0.0
		for _, meses := range dados.Anos {
			for _, mes := range meses {
				for _, fii := range mes.FIIs {
					for _, aporte := range fii.Aportes {
						if aporte.ValorTotalManual != nil {
							totalInvestido += *aporte.ValorTotalManual
						} else {
							totalInvestido += aporte.ValorTotal
						}
					}
				}
			}
		}
		ajuste := dados.ValorAjusteFIIs

		ClearTerminal()
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Println("║                     FIIs                            ║")
		fmt.Println("╟──────────────────────────────────────────────────────╢")
		fmt.Printf("║ %-54s║\n", fmt.Sprintf("Carteira: R$ %s", FormatFloatBR(totalInvestido+ajuste)))
		sinal := "+"
		if ajuste < 0 {
			sinal = "-"
		}
		fmt.Printf("║ %-54s║\n", fmt.Sprintf("Lucro/Prejuízo: R$ %s%s", sinal, FormatFloatBR(abs(ajuste))))
		fmt.Println("╠══════════════════════════════════════════════════════╣")
		fmt.Println("║ 1. Adicionar/editar FIIs do mês                     ║")
		fmt.Println("║ 2. Gerenciar dividendos e vendas                    ║")
		fmt.Println("║ 3. Ver DARF a pagar                                 ║")
		fmt.Println("║ 4. Ver FIIs conhecidos                              ║")
		fmt.Println("║ 5. Ajustar valor da carteira                        ║")
		fmt.Println("║ 6. Ajuste Preço Médio                               ║")
		fmt.Println("║ 7. Voltar ao menu principal                         ║")
		fmt.Println("╚══════════════════════════════════════════════════════╝")

		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			GerenciarFIIsMes(dados, scanner)
		case "2":
			GerenciarDividendosEVendas(dados, scanner)
		case "3":
			MostrarDARFAPagar(dados, scanner)
		case "4":
			MostrarFIIsConhecidos(dados, scanner)
		case "5":
			MostrarEEditarTotalInvestidoFIIs(dados, scanner)
		case "6":
			AjustarPrecoMedioFIIs(dados, scanner)
		case "7":
			return
		default:
			fmt.Println("Opção inválida.")
			Pause(2000)
		}
	}
}

func GerenciarFIIsMes(dados *Dados, scanner *bufio.Scanner) {
	// Seleção de ano e mês (reutilizar lógica existente)
	ano, mes := SelecionarAnoMes(dados, scanner)
	if ano == "" || mes == "" {
		return
	}

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}

	mesData := dados.Anos[ano][mes]
	m := &mesData

	for {
		ClearTerminal()
		fmt.Printf("╔══════════════════════════════════════════════════════╗\n")
		fmt.Printf("║                FIIs - %s/%s                        ║\n", NomeMes(mes), ano)
		fmt.Printf("╠══════════════════════════════════════════════════════╣\n")
		fmt.Printf("║ 1. Adicionar FII                                    ║\n")
		fmt.Printf("║ 2. Editar FII                                       ║\n")
		fmt.Printf("║ 3. Remover FII                                      ║\n")
		fmt.Printf("║ 4. Ver FIIs do mês                                  ║\n")
		fmt.Printf("║ 5. Voltar                                           ║\n")
		fmt.Printf("╚══════════════════════════════════════════════════════╝\n")

		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			AdicionarFII(dados, m, scanner)
		case "2":
			EditarFII(m, scanner)
		case "3":
			RemoverFII(m, scanner)
		case "4":
			MostrarFIIsMes(m, mes, ano, scanner)
		case "5":
			// Salvar as mudanças no mês
			dados.Anos[ano][mes] = *m
			return
		default:
			fmt.Println("Opção inválida.")
			Pause(2000)
		}
	}
}

func SelecionarAnoMes(dados *Dados, scanner *bufio.Scanner) (string, string) {
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

	return ano, mes
}

func SelecionarAnoMesComFIIs(dados *Dados, scanner *bufio.Scanner) (string, string) {
	// Filtrar anos que têm FIIs
	anosComFIIs := []string{}
	for ano, meses := range dados.Anos {
		for _, mes := range meses {
			if len(mes.FIIs) > 0 {
				anosComFIIs = append(anosComFIIs, ano)
				break
			}
		}
	}

	if len(anosComFIIs) == 0 {
		fmt.Println("Nenhum ano com FIIs encontrado.")
		Pause(2000)
		return "", ""
	}

	// Ordenar anos
	sort.Strings(anosComFIIs)

	// Seleção de ano
	ano := ""
	fmt.Println("Anos disponíveis:")
	for i, a := range anosComFIIs {
		fmt.Printf("%d - %s\n", i+1, a)
	}
	fmt.Print("Digite o número ou o ano desejado (YYYY): ")
	scanner.Scan()
	input := scanner.Text()
	if idx, err := strconv.Atoi(input); err == nil {
		if idx >= 1 && idx <= len(anosComFIIs) {
			ano = anosComFIIs[idx-1]
		}
	}
	if ano == "" {
		for _, a := range anosComFIIs {
			if a == input {
				ano = a
			}
		}
	}
	if ano == "" {
		fmt.Println("Ano inválido.")
		Pause(2000)
		return "", ""
	}

	// Filtrar meses que têm FIIs no ano selecionado
	mesesComFIIs := []string{}
	if dados.Anos[ano] != nil {
		for mes, mesData := range dados.Anos[ano] {
			if len(mesData.FIIs) > 0 {
				mesesComFIIs = append(mesesComFIIs, mes)
			}
		}
	}

	if len(mesesComFIIs) == 0 {
		fmt.Println("Nenhum mês com FIIs encontrado neste ano.")
		Pause(2000)
		return "", ""
	}

	// Ordenar meses
	sort.Strings(mesesComFIIs)

	// Seleção de mês
	mes := ""
	fmt.Println("Meses disponíveis:")
	for i, m := range mesesComFIIs {
		fmt.Printf("%d - %s\n", i+1, NomeMes(m))
	}
	fmt.Print("Digite o número ou o mês desejado (MM): ")
	scanner.Scan()
	input = scanner.Text()
	if idx, err := strconv.Atoi(input); err == nil {
		if idx >= 1 && idx <= len(mesesComFIIs) {
			mes = mesesComFIIs[idx-1]
		}
	}
	if mes == "" {
		for _, m := range mesesComFIIs {
			if m == input {
				mes = m
			}
		}
	}
	if mes == "" {
		fmt.Println("Mês inválido.")
		Pause(2000)
		return "", ""
	}

	return ano, mes
}

func AdicionarFII(dados *Dados, m *Mes, scanner *bufio.Scanner) {
	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                ADICIONAR FII                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	// Mostrar FIIs conhecidos se existirem
	if len(dados.FIIsConhecidos) > 0 {
		fmt.Println("\nFIIs conhecidos:")
		for i, codigo := range dados.FIIsConhecidos {
			fmt.Printf("%d - %s\n", i+1, codigo)
		}
		fmt.Println()
	}

	codigo := InputBox("Código do FII (ex: VGIR11):", scanner)
	codigo = strings.ToUpper(strings.TrimSpace(codigo))
	if codigo == "" {
		return
	}

	quantidadeStr := InputBox("Quantidade de cotas:", scanner)
	quantidade, err := strconv.Atoi(quantidadeStr)
	if err != nil || quantidade <= 0 {
		fmt.Println("Quantidade inválida.")
		Pause(2000)
		return
	}

	precoStr := InputBox("Preço por cota (R$):", scanner)
	precoStr = strings.ReplaceAll(precoStr, ",", ".")
	preco, err := strconv.ParseFloat(precoStr, 64)
	if err != nil || preco <= 0 {
		fmt.Println("Preço inválido.")
		Pause(2000)
		return
	}

	// Definir data padrão como hoje
	hoje := time.Now()
	dataStr := fmt.Sprintf("%02d/%02d/%04d", hoje.Day(), hoje.Month(), hoje.Year())
	dataInput := InputBox("Data do aporte (DD/MM/AAAA) [Enter para hoje]:", scanner)
	if dataInput != "" {
		dataStr = dataInput
	}

	valorTotal := float64(quantidade) * preco

	// Perguntar valor manual
	valorManualStr := InputBox("Colocar manual o valor final do aporte? (Enter para manter automático):", scanner)
	valorManualStr = strings.ReplaceAll(valorManualStr, ",", ".")
	var valorTotalManual *float64
	if valorManualStr != "" {
		valorManual, err := strconv.ParseFloat(valorManualStr, 64)
		if err == nil && valorManual > 0 {
			valorTotalManual = &valorManual
			valorTotal = valorManual
		}
	}

	novoAporte := FIIAporte{
		Quantidade:       quantidade,
		PrecoCota:        preco,
		ValorTotal:       valorTotal,
		ValorTotalManual: valorTotalManual,
		Data:             dataStr,
	}

	// Verificar se o FII já existe no mês
	fiiExistente := EncontrarFIIPorCodigo(m.FIIs, codigo)
	if fiiExistente != nil {
		// Adicionar aporte ao FII existente
		fiiExistente.Aportes = append(fiiExistente.Aportes, novoAporte)
		fmt.Printf("\n✅ Aporte adicionado ao FII %s!\n", codigo)
	} else {
		// Criar novo FII
		novoFII := FII{
			Codigo:  codigo,
			Aportes: []FIIAporte{novoAporte},
		}
		m.FIIs = append(m.FIIs, novoFII)
		fmt.Printf("\n✅ Novo FII %s criado!\n", codigo)
	}

	// AporteFIIs não existe mais no modelo, será calculado dinamicamente

	// Adicionar à lista de FIIs conhecidos
	AdicionarFIIConhecido(dados, codigo)

	fmt.Printf("Quantidade: %d cotas\n", quantidade)
	fmt.Printf("Preço por cota: R$ %.2f\n", preco)
	fmt.Printf("Valor total: R$ %.2f\n", valorTotal)
	Pause(3000)
}

func EditarFII(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII para editar.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                EDITAR FII                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		totalQtd := 0
		totalValor := 0.0
		for _, aporte := range fii.Aportes {
			totalQtd += aporte.Quantidade
			totalValor += aporte.ValorTotal
		}
		fmt.Printf("%d - %s (%d cotas, R$ %.2f total, %d aportes)\n", i+1, fii.Codigo, totalQtd, totalValor, len(fii.Aportes))
	}

	opcaoStr := InputBox("Escolha o FII para editar:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]

	if len(fii.Aportes) == 0 {
		fmt.Println("Nenhum aporte para editar neste FII.")
		Pause(2000)
		return
	}

	fmt.Printf("\nAportes de %s:\n", fii.Codigo)
	for i, aporte := range fii.Aportes {
		data := aporte.Data
		if data == "" {
			data = "Data não informada"
		}
		fmt.Printf("%d - %d cotas, R$ %.2f cada, R$ %.2f total (%s)\n", i+1, aporte.Quantidade, aporte.PrecoCota, aporte.ValorTotal, data)
	}

	apStr := InputBox("Escolha o aporte para editar:", scanner)
	apIdx, err := strconv.Atoi(apStr)
	if err != nil || apIdx < 1 || apIdx > len(fii.Aportes) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	ap := &fii.Aportes[apIdx-1]
	fmt.Printf("Quantidade atual: %d\n", ap.Quantidade)
	qtdStr := InputBox("Nova quantidade (Enter para manter):", scanner)
	if qtdStr != "" {
		qtd, err := strconv.Atoi(qtdStr)
		if err == nil && qtd > 0 {
			ap.Quantidade = qtd
		}
	}

	fmt.Printf("Preço atual: R$ %.2f\n", ap.PrecoCota)
	precoStr := InputBox("Novo preço por cota (Enter para manter):", scanner)
	if precoStr != "" {
		precoStr = strings.ReplaceAll(precoStr, ",", ".")
		preco, err := strconv.ParseFloat(precoStr, 64)
		if err == nil && preco > 0 {
			ap.PrecoCota = preco
		}
	}

	// Valor manual
	if ap.ValorTotalManual != nil {
		fmt.Printf("Valor final manual atual: R$ %.2f\n", *ap.ValorTotalManual)
	}
	manualStr := InputBox("Valor Total [MANUAL] (Enter para manter/limpar):", scanner)
	manualStr = strings.ReplaceAll(manualStr, ",", ".")
	if manualStr != "" {
		manual, err := strconv.ParseFloat(manualStr, 64)
		if err == nil && manual > 0 {
			ap.ValorTotalManual = &manual
			ap.ValorTotal = manual
		} else {
			ap.ValorTotalManual = nil
			ap.ValorTotal = float64(ap.Quantidade) * ap.PrecoCota
		}
	} else {
		// Se deixou vazio, limpar valor manual e recalcular
		ap.ValorTotalManual = nil
		ap.ValorTotal = float64(ap.Quantidade) * ap.PrecoCota
	}

	// Ao editar aporte:
	fmt.Printf("Data atual do aporte: %s\n", ap.Data)
	dataEdit := InputBox("Nova data (DD/MM/AAAA) [Enter para manter]:", scanner)
	if dataEdit != "" {
		ap.Data = dataEdit
	}

	// Recalcular valor total apenas se não tiver valor manual
	if ap.ValorTotalManual == nil {
		ap.ValorTotal = float64(ap.Quantidade) * ap.PrecoCota
	}
	fmt.Println("✅ Aporte atualizado!")
	Pause(2000)
}

func RemoverFII(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII para remover.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║               REMOVER FII/APORTE                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		totalQtd := 0
		totalValor := 0.0
		for _, aporte := range fii.Aportes {
			totalQtd += aporte.Quantidade
			totalValor += aporte.ValorTotal
		}
		fmt.Printf("%d - %s (%d cotas, R$ %.2f total, %d aportes)\n", i+1, fii.Codigo, totalQtd, totalValor, len(fii.Aportes))
	}

	opcaoStr := InputBox("Escolha o FII para remover aporte ou o FII inteiro:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]

	if len(fii.Aportes) == 0 {
		fmt.Println("Nenhum aporte para remover neste FII.")
		Pause(2000)
		return
	}

	fmt.Printf("\nAportes de %s:\n", fii.Codigo)
	for i, aporte := range fii.Aportes {
		data := aporte.Data
		if data == "" {
			data = "Data não informada"
		}
		fmt.Printf("%d - %d cotas, R$ %.2f cada, R$ %.2f total (%s)\n", i+1, aporte.Quantidade, aporte.PrecoCota, aporte.ValorTotal, data)
	}
	fmt.Printf("%d - Remover FII inteiro\n", len(fii.Aportes)+1)

	apStr := InputBox("Escolha o aporte para remover ou o FII inteiro:", scanner)
	apIdx, err := strconv.Atoi(apStr)
	if err != nil || apIdx < 1 || apIdx > len(fii.Aportes)+1 {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	if apIdx == len(fii.Aportes)+1 {
		// Remover FII inteiro
		m.FIIs = append(m.FIIs[:opcao-1], m.FIIs[opcao:]...)
		fmt.Println("✅ FII removido!")
		Pause(2000)
		return
	}

	// Remover aporte específico
	fii.Aportes = append(fii.Aportes[:apIdx-1], fii.Aportes[apIdx:]...)
	// AporteFIIs não existe mais no modelo, será calculado dinamicamente
	fmt.Println("✅ Aporte removido!")
	Pause(2000)
}

func MostrarFIIsMes(m *Mes, mes, ano string, scanner *bufio.Scanner) {
	ClearTerminal()
	fmt.Printf("╔══════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                FIIs - %s/%s                        ║\n", NomeMes(mes), ano)
	fmt.Printf("╚══════════════════════════════════════════════════════╝\n")

	if len(m.FIIs) == 0 {
		fmt.Println("\nNenhum FII registrado neste mês.")
	} else {
		fmt.Printf("\nTotal de FIIs: %d\n", len(m.FIIs))
		fmt.Printf("Valor total investido: R$ %s\n", FormatFloatBR(CalcularValorTotalFIIs(m.FIIs)))
		fmt.Println("\nDetalhes:")
		fmt.Println("---------------------------------------")

		for i, fii := range m.FIIs {
			totalQtd := 0
			totalValor := 0.0
			for _, aporte := range fii.Aportes {
				totalQtd += aporte.Quantidade
				totalValor += aporte.ValorTotal
			}
			fmt.Printf("%d. %s\n", i+1, fii.Codigo)
			fmt.Printf("   Total: %d cotas, R$ %.2f\n", totalQtd, totalValor)
			for _, aporte := range fii.Aportes {
				data := aporte.Data
				if data == "" {
					data = "Data não informada"
				}
				fmt.Printf("   Aporte (%s): %d cotas, R$ %.2f cada, R$ %.2f total\n", data, aporte.Quantidade, aporte.PrecoCota, aporte.ValorTotal)
			}
			fmt.Println("---------------------------------------")
		}
	}

	InputBox("Pressione Enter para continuar...", scanner)
}

func MostrarFIIsConhecidos(dados *Dados, scanner *bufio.Scanner) {
	for {
		ClearTerminal()
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Println("║              FIIs CONHECIDOS                        ║")
		fmt.Println("╠══════════════════════════════════════════════════════╣")
		fmt.Println("║ 1. Ver lista de FIIs conhecidos                     ║")
		fmt.Println("║ 2. Remover FII conhecido                           ║")
		fmt.Println("║ 3. Voltar                                          ║")
		fmt.Println("║ 4. Ajuste Preço Médio                               ║")
		fmt.Println("╚══════════════════════════════════════════════════════╝")

		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			MostrarListaFIIsConhecidos(dados, scanner)
		case "2":
			RemoverFIIConhecido(dados, scanner)
		case "3":
			return
		case "4":
			AjustarPrecoMedioFIIs(dados, scanner)
		default:
			fmt.Println("Opção inválida.")
			Pause(2000)
		}
	}
}

func MostrarListaFIIsConhecidos(dados *Dados, scanner *bufio.Scanner) {
	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║              FIIs CONHECIDOS                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	if len(dados.FIIsConhecidos) == 0 {
		fmt.Println("\nNenhum FII conhecido ainda.")
	} else {
		fmt.Printf("\nTotal de FIIs conhecidos: %d\n", len(dados.FIIsConhecidos))
		fmt.Println("\nLista:")
		for i, codigo := range dados.FIIsConhecidos {
			fmt.Printf("%d - %s\n", i+1, codigo)
		}
	}

	InputBox("Pressione Enter para continuar...", scanner)
}

func RemoverFIIConhecido(dados *Dados, scanner *bufio.Scanner) {
	if len(dados.FIIsConhecidos) == 0 {
		ClearTerminal()
		fmt.Println("Nenhum FII conhecido para remover.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║            REMOVER FII CONHECIDO                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs conhecidos:")
	for i, codigo := range dados.FIIsConhecidos {
		fmt.Printf("%d - %s\n", i+1, codigo)
	}

	opcaoStr := InputBox("Escolha o FII para remover:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(dados.FIIsConhecidos) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fiiParaRemover := dados.FIIsConhecidos[opcao-1]
	confirm := InputBox(fmt.Sprintf("Tem certeza que deseja remover '%s' da lista de FIIs conhecidos? (s/n): ", fiiParaRemover), scanner)
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm == "s" || confirm == "sim" {
		// Remover o FII da lista
		dados.FIIsConhecidos = append(dados.FIIsConhecidos[:opcao-1], dados.FIIsConhecidos[opcao:]...)
		fmt.Printf("✅ FII '%s' removido da lista de conhecidos!\n", fiiParaRemover)
	} else {
		fmt.Println("Operação cancelada.")
	}
	Pause(2000)
}

func GerenciarDividendosEVendas(dados *Dados, scanner *bufio.Scanner) {
	// Seleção de ano e mês apenas para períodos que têm FIIs
	ano, mes := SelecionarAnoMesComFIIs(dados, scanner)
	if ano == "" || mes == "" {
		return
	}

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}

	mesData := dados.Anos[ano][mes]
	m := &mesData

	for {
		ClearTerminal()
		fmt.Printf("╔══════════════════════════════════════════════════════╗\n")
		fmt.Printf("║         DIVIDENDOS E VENDAS - %s/%s                ║\n", NomeMes(mes), ano)
		fmt.Printf("╠══════════════════════════════════════════════════════╣\n")
		fmt.Printf("║ 1. Adicionar dividendos                             ║\n")
		fmt.Printf("║ 2. Registrar venda de cotas                         ║\n")
		fmt.Printf("║ 3. Ver resumo de dividendos e vendas                ║\n")
		fmt.Printf("║ 4. Editar dividendos                                ║\n")
		fmt.Printf("║ 5. Remover dividendos                              ║\n")
		fmt.Printf("║ 6. Editar venda de cotas                            ║\n")
		fmt.Printf("║ 7. Remover venda de cotas                           ║\n")
		fmt.Printf("║ 8. Voltar                                           ║\n")
		fmt.Printf("╚══════════════════════════════════════════════════════╝\n")

		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			AdicionarDividendos(m, scanner)
		case "2":
			RegistrarVendaCotas(m, scanner)
		case "3":
			MostrarResumoDividendosEVendas(m, mes, ano, scanner)
		case "4":
			EditarDividendos(m, scanner) // a implementar
		case "5":
			RemoverDividendos(m, scanner) // a implementar
		case "6":
			EditarVendaCotas(m, scanner) // a implementar
		case "7":
			RemoverVendaCotas(m, scanner) // a implementar
		case "8":
			dados.Anos[ano][mes] = *m
			return
		default:
			fmt.Println("Opção inválida.")
			Pause(2000)
		}
	}
}

func AdicionarDividendos(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║              ADICIONAR DIVIDENDOS                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		fmt.Printf("%d - %s (Dividendos atuais: R$ %.2f)\n", i+1, fii.Codigo, fii.Dividendos)
	}

	opcaoStr := InputBox("Escolha o FII para adicionar dividendos:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]
	fmt.Printf("Dividendos atuais de %s: R$ %.2f\n", fii.Codigo, fii.Dividendos)

	valorStr := InputBox("Digite o valor dos dividendos (R$):", scanner)
	valorStr = strings.ReplaceAll(valorStr, ",", ".")
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil || valor < 0 {
		fmt.Println("Valor inválido.")
		Pause(2000)
		return
	}

	fii.Dividendos = valor
	fmt.Printf("✅ Dividendos de R$ %.2f definidos para %s!\n", valor, fii.Codigo)
	Pause(2000)
}

func RegistrarVendaCotas(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║            REGISTRAR VENDA DE COTAS                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		totalQtd := 0
		for _, aporte := range fii.Aportes {
			totalQtd += aporte.Quantidade
		}
		fmt.Printf("%d - %s (%d cotas disponíveis)\n", i+1, fii.Codigo, totalQtd)
	}

	opcaoStr := InputBox("Escolha o FII para vender cotas:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]

	if len(fii.Aportes) == 0 {
		fmt.Println("Nenhum aporte disponível para venda neste FII.")
		Pause(2000)
		return
	}

	fmt.Println("\nAportes disponíveis:")
	for i, aporte := range fii.Aportes {
		data := aporte.Data
		if data == "" {
			data = "Data não informada"
		}
		fmt.Printf("%d - %d cotas, R$ %.2f cada, R$ %.2f total (%s)\n", i+1, aporte.Quantidade, aporte.PrecoCota, aporte.ValorTotal, data)
	}

	apStr := InputBox("Escolha o aporte do qual deseja vender cotas:", scanner)
	apIdx, err := strconv.Atoi(apStr)
	if err != nil || apIdx < 1 || apIdx > len(fii.Aportes) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	aporte := &fii.Aportes[apIdx-1]
	cotasDisponiveis := aporte.Quantidade
	if cotasDisponiveis <= 0 {
		fmt.Println("Não há cotas disponíveis para venda neste aporte.")
		Pause(2000)
		return
	}

	fmt.Printf("Cotas disponíveis para venda deste aporte: %d\n", cotasDisponiveis)
	qtdStr := InputBox("Quantidade de cotas a vender:", scanner)
	qtd, err := strconv.Atoi(qtdStr)
	if err != nil || qtd <= 0 || qtd > cotasDisponiveis {
		fmt.Println("Quantidade inválida.")
		Pause(2000)
		return
	}

	precoStr := InputBox("Preço por cota na venda (R$):", scanner)
	precoStr = strings.ReplaceAll(precoStr, ",", ".")
	precoVenda, err := strconv.ParseFloat(precoStr, 64)
	if err != nil || precoVenda <= 0 {
		fmt.Println("Preço inválido.")
		Pause(2000)
		return
	}

	taxasStr := InputBox("Taxas pagas na venda (R$):", scanner)
	taxasStr = strings.ReplaceAll(taxasStr, ",", ".")
	taxas, err := strconv.ParseFloat(taxasStr, 64)
	if err != nil || taxas < 0 {
		fmt.Println("Taxas inválidas.")
		Pause(2000)
		return
	}

	dataVenda := InputBox("Data da venda (DD/MM/AAAA) [Enter para hoje]:", scanner)
	if dataVenda == "" {
		hoje := time.Now()
		dataVenda = fmt.Sprintf("%02d/%02d/%04d", hoje.Day(), hoje.Month(), hoje.Year())
	}

	// Calcular valores
	valorTotalVenda := float64(qtd) * precoVenda
	precoMedioGlobal := CalcularPrecoMedioFII(*fii)
	custoTotal := float64(qtd) * precoMedioGlobal
	lucro := (valorTotalVenda - taxas) - custoTotal
	darf := 0.0
	if lucro > 0 {
		darf = lucro * 0.20 // 20%
	}

	// Atualizar quantidade de cotas do aporte
	aporte.Quantidade -= qtd

	// Registrar venda
	venda := FIIVenda{
		Quantidade: qtd,
		PrecoVenda: precoVenda,
		ValorTotal: valorTotalVenda,
		LucroVenda: lucro,
		DARF:       darf,
		Data:       dataVenda,
		Taxas:      taxas,
		AporteData: aporte.Data,
	}
	fii.Vendas = append(fii.Vendas, venda)

	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                RESUMO DA VENDA                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Printf("Aporte de origem: %s\n", aporte.Data)
	fmt.Printf("Quantidade vendida: %d cotas\n", qtd)
	fmt.Printf("Preço de venda: R$ %.2f\n", precoVenda)
	fmt.Printf("Valor total da venda: R$ %.2f\n", valorTotalVenda)
	fmt.Printf("Preço médio de compra (global): R$ %.2f\n", precoMedioGlobal)
	fmt.Printf("Custo total de compra: R$ %.2f\n", custoTotal)
	fmt.Printf("Taxas pagas: R$ %.2f\n", taxas)
	fmt.Printf("Lucro líquido: R$ %.2f\n", lucro)
	if darf > 0 {
		fmt.Printf("DARF a pagar (20%%): R$ %.2f\n", darf)
	} else {
		fmt.Println("Não há DARF a pagar (sem lucro na venda).")
	}
	Pause(4000)
}

func MostrarResumoDividendosEVendas(m *Mes, mes, ano string, scanner *bufio.Scanner) {
	ClearTerminal()
	fmt.Printf("╔══════════════════════════════════════════════════════╗\n")
	fmt.Printf("║         DIVIDENDOS E VENDAS - %s/%s                ║\n", NomeMes(mes), ano)
	fmt.Printf("╚══════════════════════════════════════════════════════╝\n")

	if len(m.FIIs) == 0 {
		fmt.Println("\nNenhum FII registrado neste mês.")
	} else {
		totalDividendos := 0.0
		totalVendas := 0.0
		totalLucroVendas := 0.0
		totalDARF := 0.0

		fmt.Println("\nResumo por FII:")
		fmt.Println("---------------------------------------")

		for i, fii := range m.FIIs {
			fmt.Printf("%d. %s\n", i+1, fii.Codigo)
			fmt.Printf("   Dividendos: R$ %.2f\n", fii.Dividendos)

			if len(fii.Vendas) > 0 {
				fmt.Printf("   Vendas:\n")
				for idx, venda := range fii.Vendas {
					fmt.Printf("     Venda %d: %d cotas, R$ %.2f cada, R$ %.2f total\n",
						idx+1, venda.Quantidade, venda.PrecoVenda, venda.ValorTotal)
					fmt.Printf("     Lucro: R$ %.2f, DARF: R$ %.2f\n", venda.LucroVenda, venda.DARF)
				}
			} else {
				fmt.Printf("   Nenhuma venda registrada\n")
			}

			totalDividendos += fii.Dividendos
			for _, venda := range fii.Vendas {
				totalVendas += venda.ValorTotal
				totalLucroVendas += venda.LucroVenda
				totalDARF += venda.DARF
			}

			fmt.Println("---------------------------------------")
		}

		fmt.Printf("\nTOTAIS DO MÊS:\n")
		fmt.Printf("Total de dividendos: R$ %.2f\n", totalDividendos)
		fmt.Printf("Total de vendas: R$ %.2f\n", totalVendas)
		fmt.Printf("Total de lucro das vendas: R$ %.2f\n", totalLucroVendas)
		fmt.Printf("Total de DARF a pagar: R$ %.2f\n", totalDARF)
		fmt.Printf("Lucro FIIs: R$ %.2f\n", totalDividendos+totalLucroVendas-totalDARF)
	}

	InputBox("Pressione Enter para continuar...", scanner)
}

func MostrarDARFAPagar(dados *Dados, scanner *bufio.Scanner) {
	ClearTerminal()
	linhas := []string{" DARF A PAGAR ", "---"}

	// Coletar todos os DARFs por mês/ano
	darfPorMes := make(map[string]map[string]float64) // ano -> mes -> valor
	totalDARF := 0.0

	for ano, mesesMap := range dados.Anos {
		for mes, m := range mesesMap {
			darfMes := CalcularDARFTotal(m.FIIs)
			if darfMes > 0 {
				if darfPorMes[ano] == nil {
					darfPorMes[ano] = make(map[string]float64)
				}
				darfPorMes[ano][mes] = darfMes
				totalDARF += darfMes
			}
		}
	}

	if totalDARF == 0 {
		linhas = append(linhas, "✅ Nenhum DARF a pagar!")
	} else {
		linhas = append(linhas, "⚠️  ATENÇÃO: Você tem DARF a pagar!", fmt.Sprintf("Total de DARF: R$ %s", FormatFloatBR(totalDARF)), "---")
		anos := OrdenarChaves(darfPorMes)
		for _, ano := range anos {
			meses := OrdenarChaves(darfPorMes[ano])
			linhas = append(linhas, fmt.Sprintf("Ano %s:", ano))
			for _, mes := range meses {
				darf := darfPorMes[ano][mes]
				ultimoDia, mesPagamento, anoPagamento := CalcularPrazoDARF(mes, ano)
				linhas = append(linhas, fmt.Sprintf("  %s: R$ %s", NomeMes(mes), FormatFloatBR(darf)))
				linhas = append(linhas, fmt.Sprintf("    Prazo: até %02d/%02d/%04d", ultimoDia, mesPagamento, anoPagamento))
			}
			linhas = append(linhas, "---")
		}
		linhas = append(linhas, "💡 Dica: DARF pode ser pago até o último dia do mês seguinte ao mês em que ocorreu a venda.")
	}
	PrintCaixa(linhas)
	InputBox("Pressione Enter para continuar...", scanner)
}

// Funções stub para edição/remoção de dividendos e vendas
func EditarDividendos(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                EDITAR DIVIDENDOS                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		fmt.Printf("%d - %s (Dividendos atuais: R$ %.2f)\n", i+1, fii.Codigo, fii.Dividendos)
	}

	opcaoStr := InputBox("Escolha o FII para editar dividendos:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]
	fmt.Printf("Dividendos atuais de %s: R$ %.2f\n", fii.Codigo, fii.Dividendos)

	valorStr := InputBox("Digite o novo valor dos dividendos (R$):", scanner)
	valorStr = strings.ReplaceAll(valorStr, ",", ".")
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil || valor < 0 {
		fmt.Println("Valor inválido.")
		Pause(2000)
		return
	}

	fii.Dividendos = valor
	fmt.Printf("✅ Dividendos de %s atualizados para R$ %.2f!\n", fii.Codigo, valor)
	Pause(2000)
}

func RemoverDividendos(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                REMOVER DIVIDENDOS                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		fmt.Printf("%d - %s (Dividendos atuais: R$ %.2f)\n", i+1, fii.Codigo, fii.Dividendos)
	}

	opcaoStr := InputBox("Escolha o FII para remover dividendos:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]
	fmt.Printf("Dividendos atuais de %s: R$ %.2f\n", fii.Codigo, fii.Dividendos)
	confirma := InputBox("Tem certeza que deseja remover (zerar) os dividendos deste FII? (s/n):", scanner)
	if strings.TrimSpace(confirma) == "" || strings.ToLower(strings.TrimSpace(confirma)) == "s" {
		fii.Dividendos = 0
		fmt.Printf("✅ Dividendos de %s removidos (zerados)!\n", fii.Codigo)
	} else {
		fmt.Println("Operação cancelada.")
	}
	Pause(2000)
}

func EditarVendaCotas(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                EDITAR VENDA DE COTAS                ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		fmt.Printf("%d - %s (%d vendas)\n", i+1, fii.Codigo, len(fii.Vendas))
	}

	opcaoStr := InputBox("Escolha o FII para editar venda:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]
	if len(fii.Vendas) == 0 {
		fmt.Println("Nenhuma venda registrada neste FII.")
		Pause(2000)
		return
	}

	fmt.Println("\nVendas registradas:")
	for i, venda := range fii.Vendas {
		fmt.Printf("%d - %d cotas, R$ %.2f cada, R$ %.2f total, Data: %s\n", i+1, venda.Quantidade, venda.PrecoVenda, venda.ValorTotal, venda.Data)
	}

	vendaStr := InputBox("Escolha a venda para editar:", scanner)
	vendaIdx, err := strconv.Atoi(vendaStr)
	if err != nil || vendaIdx < 1 || vendaIdx > len(fii.Vendas) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	venda := &fii.Vendas[vendaIdx-1]

	// Editar campos
	qtdStr := InputBox(fmt.Sprintf("Nova quantidade de cotas (atual: %d, Enter para manter):", venda.Quantidade), scanner)
	if qtdStr != "" {
		qtd, err := strconv.Atoi(qtdStr)
		if err == nil && qtd > 0 {
			venda.Quantidade = qtd
		}
	}
	precoStr := InputBox(fmt.Sprintf("Novo preço por cota (atual: %.2f, Enter para manter):", venda.PrecoVenda), scanner)
	if precoStr != "" {
		precoStr = strings.ReplaceAll(precoStr, ",", ".")
		preco, err := strconv.ParseFloat(precoStr, 64)
		if err == nil && preco > 0 {
			venda.PrecoVenda = preco
			venda.ValorTotal = float64(venda.Quantidade) * preco
		}
	}
	taxasStr := InputBox(fmt.Sprintf("Novas taxas (atual: %.2f, Enter para manter):", venda.Taxas), scanner)
	if taxasStr != "" {
		taxasStr = strings.ReplaceAll(taxasStr, ",", ".")
		taxas, err := strconv.ParseFloat(taxasStr, 64)
		if err == nil && taxas >= 0 {
			venda.Taxas = taxas
		}
	}
	dataStr := InputBox(fmt.Sprintf("Nova data (atual: %s, Enter para manter):", venda.Data), scanner)
	if dataStr != "" {
		venda.Data = dataStr
	}
	fmt.Println("✅ Venda editada!")
	Pause(2000)
}

func RemoverVendaCotas(m *Mes, scanner *bufio.Scanner) {
	if len(m.FIIs) == 0 {
		fmt.Println("Nenhum FII registrado neste mês.")
		Pause(2000)
		return
	}

	ClearTerminal()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║                REMOVER VENDA DE COTAS               ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	fmt.Println("\nFIIs do mês:")
	for i, fii := range m.FIIs {
		fmt.Printf("%d - %s (%d vendas)\n", i+1, fii.Codigo, len(fii.Vendas))
	}

	opcaoStr := InputBox("Escolha o FII para remover venda:", scanner)
	opcao, err := strconv.Atoi(opcaoStr)
	if err != nil || opcao < 1 || opcao > len(m.FIIs) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	fii := &m.FIIs[opcao-1]
	if len(fii.Vendas) == 0 {
		fmt.Println("Nenhuma venda registrada neste FII.")
		Pause(2000)
		return
	}

	fmt.Println("\nVendas registradas:")
	for i, venda := range fii.Vendas {
		fmt.Printf("%d - %d cotas, R$ %.2f cada, R$ %.2f total, Data: %s\n", i+1, venda.Quantidade, venda.PrecoVenda, venda.ValorTotal, venda.Data)
	}

	vendaStr := InputBox("Escolha a venda para remover:", scanner)
	vendaIdx, err := strconv.Atoi(vendaStr)
	if err != nil || vendaIdx < 1 || vendaIdx > len(fii.Vendas) {
		fmt.Println("Opção inválida.")
		Pause(2000)
		return
	}

	venda := fii.Vendas[vendaIdx-1]
	confirma := InputBox("Tem certeza que deseja remover esta venda? (s/n):", scanner)
	if strings.TrimSpace(confirma) == "" || strings.ToLower(strings.TrimSpace(confirma)) == "s" {
		// Devolver cotas ao aporte de origem
		for i := range fii.Aportes {
			if fii.Aportes[i].Data == venda.AporteData {
				fii.Aportes[i].Quantidade += venda.Quantidade
				break
			}
		}
		// Remover venda
		fii.Vendas = append(fii.Vendas[:vendaIdx-1], fii.Vendas[vendaIdx:]...)
		fmt.Println("✅ Venda removida e cotas devolvidas ao aporte de origem!")
	} else {
		fmt.Println("Operação cancelada.")
	}
	Pause(2000)
}

// Função para mostrar e editar o total investido em FIIs
func MostrarEEditarTotalInvestidoFIIs(dados *Dados, scanner *bufio.Scanner) {
	// Calcular o valor total investido em FIIs (soma de todos os meses/anos)
	totalInvestido := 0.0
	for _, meses := range dados.Anos {
		for _, mes := range meses {
			for _, fii := range mes.FIIs {
				for _, aporte := range fii.Aportes {
					if aporte.ValorTotalManual != nil {
						totalInvestido += *aporte.ValorTotalManual
					} else {
						totalInvestido += aporte.ValorTotal
					}
				}
			}
		}
	}
	ajuste := dados.ValorAjusteFIIs

	for {
		caixa := []string{
			fmt.Sprintf("Valor total investido: R$ %s", FormatFloatBR(totalInvestido+ajuste)),
		}
		sinal := "+"
		if ajuste < 0 {
			sinal = "-"
		}
		caixa = append(caixa, fmt.Sprintf("Lucro/Prejuízo: R$ %s%s", sinal, FormatFloatBR(abs(ajuste))))
		PrintCaixa(caixa)
		fmt.Println("1 - Lucro")
		fmt.Println("2 - Prejuízo")
		fmt.Println("3 - Manual")
		fmt.Println("4 - Voltar")
		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			valorStr := InputBox("Digite o valor do lucro:", scanner)
			valor, err := ParseFloatBR(valorStr)
			if err != nil || valor < 0 {
				PrintCaixa([]string{"Valor inválido!"})
				Pause(2000)
				continue
			}
			dados.ValorAjusteFIIs += valor
			SalvarDados(*dados)
			ajuste = dados.ValorAjusteFIIs
			PrintCaixa([]string{fmt.Sprintf("Novo valor total investido (com lucro): R$ %s", FormatFloatBR(totalInvestido+ajuste))})
			Pause(2000)
		case "2":
			valorStr := InputBox("Digite o valor do prejuízo:", scanner)
			valor, err := ParseFloatBR(valorStr)
			if err != nil || valor < 0 {
				PrintCaixa([]string{"Valor inválido!"})
				Pause(2000)
				continue
			}
			dados.ValorAjusteFIIs -= valor
			if (totalInvestido + dados.ValorAjusteFIIs) < 0 {
				dados.ValorAjusteFIIs = -totalInvestido
			}
			SalvarDados(*dados)
			ajuste = dados.ValorAjusteFIIs
			PrintCaixa([]string{fmt.Sprintf("Novo valor total investido (com prejuízo): R$ %s", FormatFloatBR(totalInvestido+ajuste))})
			Pause(2000)
		case "3":
			valorStr := InputBox("Digite o valor manual:", scanner)
			valor, err := ParseFloatBR(valorStr)
			if err != nil || valor < 0 {
				PrintCaixa([]string{"Valor inválido!"})
				Pause(2000)
				continue
			}
			dados.ValorAjusteFIIs = valor - totalInvestido
			SalvarDados(*dados)
			ajuste = dados.ValorAjusteFIIs
			PrintCaixa([]string{fmt.Sprintf("Novo valor total investido (manual): R$ %s", FormatFloatBR(totalInvestido+ajuste))})
			Pause(2000)
		case "4":
			return
		default:
			PrintCaixa([]string{"Opção inválida."})
			Pause(2000)
		}
	}
}

// Função auxiliar para valor absoluto
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// Função para ajuste de preço médio dos FIIs
func AjustarPrecoMedioFIIs(dados *Dados, scanner *bufio.Scanner) {
	for {
		ClearTerminal()
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Println("║              AJUSTE PREÇO MÉDIO FII                 ║")
		fmt.Println("╠══════════════════════════════════════════════════════╣")
		fmt.Println("║ 1. Automático (padrão do sistema)                   ║")
		fmt.Println("║ 2. Manual (definir preço médio inicial)             ║")
		fmt.Println("║ 3. Voltar                                          ║")
		fmt.Println("╚══════════════════════════════════════════════════════╝")

		opcao := InputBox("Escolha uma opção:", scanner)
		switch opcao {
		case "1":
			// Remover qualquer ajuste manual global
			if dados.ValorAjusteFIIs != 0 {
				dados.ValorAjusteFIIs = 0
				fmt.Println("Ajuste manual removido. Sistema volta ao cálculo automático.")
				Pause(2000)
			}
			return
		case "2":
			// Listar FIIs disponíveis
			fiisUnicos := make(map[string]*FII)
			for _, ano := range dados.Anos {
				for _, mes := range ano {
					for i := range mes.FIIs {
						codigo := mes.FIIs[i].Codigo
						if _, existe := fiisUnicos[codigo]; !existe {
							fiisUnicos[codigo] = &mes.FIIs[i]
						}
					}
				}
			}
			if len(fiisUnicos) == 0 {
				fmt.Println("Nenhum FII disponível para ajuste.")
				Pause(2000)
				continue
			}
			// Exibir lista em uma caixinha
			codigos := make([]string, 0, len(fiisUnicos))
			for codigo := range fiisUnicos {
				codigos = append(codigos, codigo)
			}
			sort.Strings(codigos)
			caixa := []string{"FIIs disponíveis:"}
			for i, codigo := range codigos {
				precoMedio := CalcularPrecoMedioFII(*fiisUnicos[codigo])
				caixa = append(caixa, fmt.Sprintf("%d - %s (Preço Médio: R$ %.2f)", i+1, codigo, precoMedio))
			}
			PrintCaixa(caixa)
			input := InputBox("Digite o número ou código do FII para ajustar:", scanner)
			input = strings.TrimSpace(input)
			var fiiPtr *FII
			if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= len(codigos) {
				fiiPtr = fiisUnicos[codigos[idx-1]]
			} else {
				codigo := strings.ToUpper(input)
				if ptr, ok := fiisUnicos[codigo]; ok {
					fiiPtr = ptr
				}
			}
			if fiiPtr == nil {
				fmt.Println("FII não encontrado.")
				Pause(2000)
				continue
			}
			precoStr := InputBox("Digite o preço médio desejado:", scanner)
			precoStr = strings.ReplaceAll(precoStr, ",", ".")
			preco, err := strconv.ParseFloat(precoStr, 64)
			if err != nil || preco <= 0 {
				fmt.Println("Preço inválido.")
				Pause(2000)
				continue
			}
			// Calcular total de cotas
			totalCotas := 0
			for _, ap := range fiiPtr.Aportes {
				totalCotas += ap.Quantidade
			}
			if totalCotas == 0 {
				fmt.Println("FII sem cotas para ajustar.")
				Pause(2000)
				continue
			}
			// Ajustar valor do primeiro aporte para refletir o novo preço médio
			if len(fiiPtr.Aportes) > 0 {
				valorManual := preco * float64(totalCotas)
				fiiPtr.Aportes[0].ValorTotalManual = &valorManual
				fiiPtr.Aportes[0].ValorTotal = valorManual
				fmt.Printf("Preço médio ajustado para R$ %.2f em %s.\n", preco, fiiPtr.Codigo)
				Pause(2000)
			}
			return
		case "3":
			return
		default:
			fmt.Println("Opção inválida.")
			Pause(2000)
		}
	}
}
