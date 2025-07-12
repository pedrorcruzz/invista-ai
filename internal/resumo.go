package internal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func NomeMes(m string) string {
	nomes := map[string]string{
		"01": "Janeiro", "02": "Fevereiro", "03": "Março",
		"04": "Abril", "05": "Maio", "06": "Junho",
		"07": "Julho", "08": "Agosto", "09": "Setembro",
		"10": "Outubro", "11": "Novembro", "12": "Dezembro",
	}
	return nomes[m]
}

func OrdenarChaves[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FormatFloatBR formata float64 para string com vírgula como separador decimal
func FormatFloatBR(f float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", f), ".", ",")
}

func GetResumoTotalAcumuladoStr(dados Dados) string {
	anos := OrdenarChaves(dados.Anos)
	if len(anos) == 0 {
		return "Nenhum dado disponível ainda."
	}
	// Acumuladores SEM filtro (para aportes, FIIs, saídas, retiradas)
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	lucrosRetiradosTotal := 0.0

	// Acumuladores COM filtro (para lucros líquidos e saldo final)
	// valorBrutoFinal removido pois não é mais usado
	lucroLiquidoAcumulado := 0.0
	lucroLiquidoFIIsAcumulado := 0.0
	lucroMesLiquidoTotalAcumulado := 0.0

	// Novo: pegar sempre o saldo do último mês para os valores finais
	ultimoBrutoFinal := 0.0
	ultimoLiquidoFinal := 0.0

	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := OrdenarChaves(mesesMap)
		for _, mes := range meses {
			m := mesesMap[mes]
			// SEM filtro: acumula aportes, FIIs, saídas, retiradas
			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += m.AporteFIIs
			saidaSoFar += m.Saida
			lucrosRetiradosTotal += m.LucroRetirado

			// Sempre pega o saldo do último mês
			ultimoBrutoFinal = m.ValorBrutoRF
			ultimoLiquidoFinal = m.ValorLiquidoRF

			// Cálculo do lucro líquido (COM filtro)
			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
			lucroLiquidoFIIs := m.LucroLiquidoFIIs
			lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs

			lucroValido := lucroMesBruto > impostos
			if lucroValido {
				lucroLiquidoAcumulado += lucroMesLiquidoRF
				lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
				lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
			}
			saldoAnterior = m.ValorBrutoRF
		}
	}
	// Totais
	totalAportadoBruto := aporteRFSoFar + aporteFIIsSoFar
	totalAportadoLiquido := totalAportadoBruto - saidaSoFar
	// Novo: calcular aportes e saídas só de RF para o lucro bruto total
	totalAportadoBrutoRF := aporteRFSoFar
	totalAportadoLiquidoRF := totalAportadoBrutoRF - saidaSoFar
	lucroBrutoTotal := ultimoBrutoFinal - totalAportadoLiquidoRF

	// Porcentagens e valores de RF e FIIs (bruto)
	percRFBruto := 0.0
	percFIIsBruto := 0.0
	if totalAportadoBruto > 0 {
		percRFBruto = (aporteRFSoFar / totalAportadoBruto) * 100
		percFIIsBruto = (aporteFIIsSoFar / totalAportadoBruto) * 100
	}

	// Porcentagens e valores de RF e FIIs (líquido, saídas só afetam RF)
	rfLiquido := aporteRFSoFar - saidaSoFar
	fiisLiquido := aporteFIIsSoFar // FIIs não tem saída
	percRFLiquido := 0.0
	percFIIsLiquido := 0.0
	if totalAportadoLiquido > 0 {
		percRFLiquido = (rfLiquido / totalAportadoLiquido) * 100
		percFIIsLiquido = (fiisLiquido / totalAportadoLiquido) * 100
	}

	return fmt.Sprintf(`================== InvistAI ==================

--- Total Investido ---

[APORTES BRUTOS (tudo que já foi investido)]
Total aportado bruto: R$ %s
  - Renda Fixa: %.2f%% (R$ %s)
  - FIIs: %.2f%% (R$ %s)

--------------------

[APORTES LÍQUIDOS (após retiradas)]
Total aportado líquido: R$ %s
  - Renda Fixa: %.2f%% (R$ %s)
  - FIIs: %.2f%% (R$ %s)

---------------------------------------

[RENDA FIXA]
Valor Bruto Final (RF): R$ %s
Valor Líquido Final (RF): R$ %s
Lucros Retirados: R$ %s
Lucro Bruto Total (RF): R$ %s
Lucro Líquido RF: R$ %s

---------------------------------------

[FIIs]
Lucro Líquido FIIs: R$ %s

---------------------------------------

╔════════════════════════════════════════════════════╗
║  Lucro Total Líquido (RF + FIIs): R$ %s           ║
╚════════════════════════════════════════════════════╝
`,
		FormatFloatBR(totalAportadoBruto), percRFBruto, FormatFloatBR(aporteRFSoFar), percFIIsBruto, FormatFloatBR(aporteFIIsSoFar), FormatFloatBR(totalAportadoLiquido), percRFLiquido, FormatFloatBR(rfLiquido), percFIIsLiquido, FormatFloatBR(fiisLiquido), FormatFloatBR(ultimoBrutoFinal), FormatFloatBR(ultimoLiquidoFinal), FormatFloatBR(lucrosRetiradosTotal), FormatFloatBR(lucroBrutoTotal), FormatFloatBR(lucroLiquidoAcumulado), FormatFloatBR(lucroLiquidoFIIsAcumulado), FormatFloatBR(lucroMesLiquidoTotalAcumulado))
}

func GetResumoMesAtualStr(dados Dados) string {
	hoje := time.Now()
	anoAtual := fmt.Sprintf("%04d", hoje.Year())
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))
	anos := OrdenarChaves(dados.Anos)
	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := OrdenarChaves(mesesMap)
		for _, mes := range meses {
			if ano == anoAtual && mes == mesAtual {
				m := mesesMap[mes]
				lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
				impostos := m.ValorBrutoRF - m.ValorLiquidoRF
				lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
				lucroLiquidoFIIs := m.LucroLiquidoFIIs
				lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
				titulo := fmt.Sprintf("Mês: %s/%s", NomeMes(mes), ano)
				return fmt.Sprintf(`%s
  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)
---------------------------------------
  Aporte Total:         R$ %s
  Aporte RF:            R$ %s
  Aporte FIIs:          R$ %s
  Saída:                R$ %s
  Bruto RF:             R$ %s
  Líquido RF:           R$ %s
  Lucro Retirado:       R$ %s
  Lucro Líquido RF:     R$ %s
  Lucro Líquido FIIs:   R$ %s
  Lucro Mês Bruto:      R$ %s
  Lucro Mês Líquido:    R$ %s
---------------------------------------`,
					titulo,
					FormatFloatBR(m.AporteRF+m.AporteFIIs), FormatFloatBR(m.AporteRF), FormatFloatBR(m.AporteFIIs), FormatFloatBR(m.Saida), FormatFloatBR(m.ValorBrutoRF), FormatFloatBR(m.ValorLiquidoRF), FormatFloatBR(m.LucroRetirado), FormatFloatBR(lucroMesLiquidoRF), FormatFloatBR(lucroLiquidoFIIs), FormatFloatBR(lucroMesBruto), FormatFloatBR(lucroMesLiquidoTotal))
			}
			saldoAnterior = mesesMap[mes].ValorBrutoRF
		}
	}
	return "Mês atual não possui dados."
}

func MostrarResumoAno(dados Dados, ano string, horizontal bool) {
	mesesMap, ok := dados.Anos[ano]
	if !ok || len(mesesMap) == 0 {
		fmt.Printf("Não há dados para o ano %s.\n", ano)
		return
	}
	meses := OrdenarChaves(mesesMap)
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0
	saidaSoFar := 0.0
	lucrosRetiradosTotal := 0.0
	valorBrutoFinal := 0.0
	// valorLiquidoRFFinal removido pois não é mais usado
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
			fmt.Printf("| %-8s | R$ %10s | R$ %7s | R$%4s | R$%6s | R$ %9s | R$ %8s | R$ %10s | R$ %12s | R$ %14s | R$ %16s | R$ %18s | R$ %17s |\n",
				NomeMes(mes), FormatFloatBR(m.AporteRF+m.AporteFIIs), FormatFloatBR(m.AporteRF), FormatFloatBR(m.AporteFIIs), FormatFloatBR(m.Saida), FormatFloatBR(m.LucroRetirado),
				FormatFloatBR(m.ValorBrutoRF), FormatFloatBR(m.ValorLiquidoRF), FormatFloatBR(m.ValorLiquidoFIIs),
				FormatFloatBR(lucroMesBruto), FormatFloatBR(lucroMesLiquidoRF), FormatFloatBR(lucroLiquidoFIIs), FormatFloatBR(lucroMesLiquidoTotal))
		} else {
			fmt.Printf("\nMês: %s/%s\n", NomeMes(mes), ano)
			if isMesAtual {
				fmt.Println("  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)")
			}
			impostoValido := impostos > 0
			if lucroMesBruto > impostos && impostoValido {
				fmt.Println("  ✅ Agora os lucros já cobrem os impostos!")
			}
			fmt.Println("---------------------------------------")
			fmt.Printf("  Aporte Total:         R$ %s\n", FormatFloatBR(m.AporteRF+m.AporteFIIs))
			fmt.Printf("  Aporte RF:            R$ %s\n", FormatFloatBR(m.AporteRF))
			fmt.Printf("  Aporte FIIs:          R$ %s\n", FormatFloatBR(m.AporteFIIs))
			fmt.Printf("  Saída:                R$ %s\n", FormatFloatBR(m.Saida))
			fmt.Printf("  Lucro Retirado:       R$ %s\n", FormatFloatBR(m.LucroRetirado))
			fmt.Printf("  Bruto RF:             R$ %s\n", FormatFloatBR(m.ValorBrutoRF))
			fmt.Printf("  Líquido RF:           R$ %s\n", FormatFloatBR(m.ValorLiquidoRF))
			fmt.Printf("  Lucro Mês Bruto:      R$ %s\n", FormatFloatBR(lucroMesBruto))
			fmt.Printf("  Lucro Líquido RF:     R$ %s\n", FormatFloatBR(lucroMesLiquidoRF))
			fmt.Printf("  Lucro Líquido FIIs:   R$ %s\n", FormatFloatBR(lucroLiquidoFIIs))
			fmt.Printf("  Lucro Mês Líquido:    R$ %s\n", FormatFloatBR(lucroMesLiquidoTotal))
			fmt.Println("---------------------------------------")
		}
		lucroValido := lucroMesBruto > impostos
		if lucroValido {
			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += m.AporteFIIs
			saidaSoFar += m.Saida
			lucrosRetiradosTotal += m.LucroRetirado
			valorBrutoFinal = m.ValorBrutoRF
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
	fmt.Printf("Total aportado bruto: R$ %s\n", FormatFloatBR(totalAportadoBruto))
	fmt.Printf("Total aportado líquido: R$ %s\n", FormatFloatBR(totalAportadoLiquido))
	fmt.Printf("Valor bruto final (RF): R$ %s\n", FormatFloatBR(valorBrutoFinal))
	fmt.Printf("Valor líquido final (RF): R$ %s\n", FormatFloatBR(valorBrutoFinal))
	fmt.Printf("Lucro bruto total (RF): R$ %s\n", FormatFloatBR(lucroBrutoTotal))
	fmt.Printf("Lucro Líquido RF: R$ %s\n", FormatFloatBR(lucroLiquidoTotal))
	fmt.Printf("Lucro Líquido FIIs: R$ %s\n", FormatFloatBR(lucroLiquidoFIIsTotal))
	fmt.Printf("Lucro Total Líquido (RF + FIIs): R$ %s\n", FormatFloatBR(lucroMesLiquidoTotalAno))
	fmt.Printf("Lucros retirados: R$ %s\n", FormatFloatBR(lucrosRetiradosTotal))
}
