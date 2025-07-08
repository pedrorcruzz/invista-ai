package internal

import (
	"fmt"
	"sort"
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

func GetResumoTotalAcumuladoStr(dados Dados) string {
	anos := OrdenarChaves(dados.Anos)
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
		meses := OrdenarChaves(mesesMap)
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
				NomeMes(mes), m.AporteRF+m.AporteFIIs, m.AporteRF, m.AporteFIIs, m.Saida, m.LucroRetirado,
				m.ValorBrutoRF, m.ValorLiquidoRF, m.ValorLiquidoFIIs,
				lucroMesBruto, lucroMesLiquidoRF, lucroLiquidoFIIs, lucroMesLiquidoTotal)
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

