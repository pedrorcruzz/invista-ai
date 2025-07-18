package internal

import (
	"fmt"
	"sort"
	"strconv"
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

	// Para detalhes dos FIIs
	todosFIIs := make(map[string]FII)

	// Para verificar DARF a pagar
	totalDARF := 0.0

	// Para calcular lucro bruto total
	lucroBrutoTotalAcumulado := 0.0

	// Verificar se estamos no mês atual
	hoje := time.Now()
	anoAtual := fmt.Sprintf("%04d", hoje.Year())
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))

	saldoAnterior := 0.0
	for _, ano := range anos {
		mesesMap := dados.Anos[ano]
		meses := OrdenarChaves(mesesMap)
		for _, mes := range meses {
			m := mesesMap[mes]
			isMesAtual := (ano == anoAtual && mes == mesAtual)

			// SEM filtro: acumula aportes, FIIs, saídas, retiradas
			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += CalcularValorTotalFIIs(m.FIIs)
			saidaSoFar += m.Saida
			lucrosRetiradosTotal += m.LucroRetirado

			// Sempre pega o saldo do último mês
			ultimoBrutoFinal = m.ValorBrutoRF
			ultimoLiquidoFinal = m.ValorLiquidoRF

			// Acumular FIIs
			for _, fii := range m.FIIs {
				if fiiExistente, existe := todosFIIs[fii.Codigo]; existe {
					// Merge dos aportes
					fiiExistente.Aportes = append(fiiExistente.Aportes, fii.Aportes...)
					fiiExistente.Dividendos += fii.Dividendos
					fiiExistente.Vendas = append(fiiExistente.Vendas, fii.Vendas...)
					todosFIIs[fii.Codigo] = fiiExistente
				} else {
					todosFIIs[fii.Codigo] = fii
				}
			}

			// Acumular DARF
			totalDARF += CalcularDARFTotal(m.FIIs)

			// Cálculo do lucro líquido (COM filtro)
			lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
			impostos := m.ValorBrutoRF - m.ValorLiquidoRF
			lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
			lucroLiquidoFIIs := CalcularLucroLiquidoFIIs(m.FIIs)
			lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs

			// FIIs profit should always be accumulated, regardless of RF profit status
			lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs

			lucroValido := lucroMesBruto > impostos

			// Lucro bruto sempre acumula (todos os meses)
			lucroBrutoTotalAcumulado += lucroMesBruto

			// Se for o mês atual e não for válido, não acumular lucros líquidos
			if isMesAtual && !lucroValido {
				// Não acumular lucros líquidos do mês atual se não for válido
				// Mas continuar acumulando FIIs (já foi feito acima)
			} else if lucroValido {
				lucroLiquidoAcumulado += lucroMesLiquidoRF
				lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
			}
			saldoAnterior = m.ValorBrutoRF
		}
	}
	// Totais - usar valores finais (bruto e líquido) da RF
	// FIIs bruto: aportes + dividendos + lucro/prejuízo vendas + ajuste manual
	totalDividendos := 0.0
	totalLucroVendas := 0.0
	for _, fii := range todosFIIs {
		totalDividendos += fii.Dividendos
		for _, venda := range fii.Vendas {
			totalLucroVendas += venda.LucroVenda
		}
	}
	fiisBruto := aporteFIIsSoFar + totalDividendos + totalLucroVendas + dados.ValorAjusteFIIs
	totalAportadoBruto := ultimoBrutoFinal + fiisBruto
	totalAportadoLiquido := ultimoLiquidoFinal + fiisBruto
	// Lucro bruto total = valor final - total aportado bruto (sem considerar saídas no cálculo)
	lucroBrutoTotal := ultimoBrutoFinal - totalAportadoBruto
	// Corrigir: usar o acumulado dos lucros brutos dos meses válidos
	lucroBrutoTotal = lucroBrutoTotalAcumulado

	// Porcentagens e valores de RF e FIIs (bruto)
	percRFBruto := 0.0
	percFIIsBruto := 0.0
	if totalAportadoBruto > 0 {
		percRFBruto = (ultimoBrutoFinal / totalAportadoBruto) * 100
		percFIIsBruto = (fiisBruto / totalAportadoBruto) * 100
	}

	// Porcentagens e valores de RF e FIIs (líquido, saídas só afetam RF)
	rfLiquido := ultimoLiquidoFinal
	// FIIs líquido: aportes + dividendos + lucro/prejuízo vendas + ajuste manual
	totalDividendos = 0.0
	totalLucroVendas = 0.0
	for _, fii := range todosFIIs {
		totalDividendos += fii.Dividendos
		for _, venda := range fii.Vendas {
			totalLucroVendas += venda.LucroVenda
		}
	}
	fiisLiquido := aporteFIIsSoFar + totalDividendos + totalLucroVendas + dados.ValorAjusteFIIs
	percRFLiquido := 0.0
	percFIIsLiquido := 0.0
	if totalAportadoLiquido > 0 {
		percRFLiquido = (rfLiquido / totalAportadoLiquido) * 100
		percFIIsLiquido = (fiisLiquido / totalAportadoLiquido) * 100
	}

	// Preparar detalhes dos FIIs com porcentagem do lucro
	fiisDetalhes := ""
	if len(todosFIIs) > 0 {
		fiisDetalhes = "\n[FIIs Detalhados]\n"
		for codigo, fii := range todosFIIs {
			totalQtd := 0
			totalValor := 0.0
			lucroFII := fii.Dividendos
			for _, aporte := range fii.Aportes {
				totalQtd += aporte.Quantidade
				if aporte.ValorTotalManual != nil {
					totalValor += *aporte.ValorTotalManual
				} else {
					totalValor += aporte.ValorTotal
				}
			}
			// Adicionar lucro das vendas
			for _, venda := range fii.Vendas {
				lucroFII += venda.LucroVenda - venda.DARF
			}

			// Calcular porcentagem do lucro total
			porcentagem := 0.0
			if lucroLiquidoFIIsAcumulado > 0 {
				porcentagem = (lucroFII / lucroLiquidoFIIsAcumulado) * 100
			}

			fiisDetalhes += fmt.Sprintf("  - %s (%.1f%%): %d cotas (R$ %s) | Preço médio: R$ %s\n", codigo, porcentagem, totalQtd, FormatFloatBR(totalValor), FormatFloatBR(CalcularPrecoMedioFII(fii)))
		}
	}

	// Alerta de DARF
	alertaDARF := ""
	if totalDARF > 0 {
		// Coletar detalhes por mês/ano para prazo
		prazo := ""
		for ano, mesesMap := range dados.Anos {
			for mes, m := range mesesMap {
				darfMes := CalcularDARFTotal(m.FIIs)
				if darfMes > 0 {
					// Calcular prazo: último dia do mês seguinte
					mesInt, _ := strconv.Atoi(mes)
					anoInt, _ := strconv.Atoi(ano)
					mesPrazo := mesInt + 1
					anoPrazo := anoInt
					if mesPrazo > 12 {
						mesPrazo = 1
						anoPrazo++
					}
					t := time.Date(anoPrazo, time.Month(mesPrazo)+1, 0, 0, 0, 0, 0, time.UTC)
					prazo = t.Format("02/01/2006")
				}
			}
		}
		alertaDARF = "\n╔════════════════════════════════════════════════════╗\n" +
			"║  ⚠️  DARF a pagar: R$ " + FormatFloatBR(totalDARF) + " até " + prazo + "         ║\n" +
			"╚════════════════════════════════════════════════════╝\n"
	} else {
		alertaDARF = "\n╔════════════════════════════════════════════════════╗\n" +
			"║  ✅ Nenhum DARF a pagar!                           ║\n" +
			"╚════════════════════════════════════════════════════╝\n"
	}

	// Cálculo do bloco [FIIs] global
	fiisTotalInvestido := 0.0
	fiisDividendos := 0.0
	fiisLucroVendas := 0.0
	for _, fii := range todosFIIs {
		for _, aporte := range fii.Aportes {
			if aporte.ValorTotalManual != nil {
				fiisTotalInvestido += *aporte.ValorTotalManual
			} else {
				fiisTotalInvestido += aporte.ValorTotal
			}
		}
		fiisDividendos += fii.Dividendos
		for _, venda := range fii.Vendas {
			fiisLucroVendas += venda.LucroVenda
		}
	}
	fiisCarteira := fiisTotalInvestido + fiisDividendos + fiisLucroVendas + dados.ValorAjusteFIIs
	rendimentoFIIs := fiisDividendos + fiisLucroVendas
	if abs(rendimentoFIIs) < 0.005 {
		rendimentoFIIs = 0.0
	}
	if abs(fiisTotalInvestido) < 0.005 {
		fiisTotalInvestido = 0.0
	}
	if abs(fiisCarteira) < 0.005 {
		fiisCarteira = 0.0
	}

	// Preparar detalhes dos FIIs fora da caixinha
	fiisDetalhes = ""
	if len(todosFIIs) > 0 {
		fiisDetalhes = "\n[FIIs Detalhados]\n"
		for codigo, fii := range todosFIIs {
			totalQtd := 0
			totalValor := 0.0
			for _, aporte := range fii.Aportes {
				totalQtd += aporte.Quantidade
				if aporte.ValorTotalManual != nil {
					totalValor += *aporte.ValorTotalManual
				} else {
					totalValor += aporte.ValorTotal
				}
			}
			precoMedio := CalcularPrecoMedioFII(fii)
			fiisDetalhes += fmt.Sprintf("  - %s: %d cotas (R$ %s) | Preço médio: R$ %s\n", codigo, totalQtd, FormatFloatBR(totalValor), FormatFloatBR(precoMedio))
		}
	}

	// Montar o resumo principal sem bug de formatação
	var resumo string
	resumo = fmt.Sprintf(`================== InvistAI ==================

--- Total Investido ---

[VALOR BRUTO]
Total valor bruto: R$ %s
  - Renda Fixa: %.2f%% (R$ %s)
  - FIIs: %.2f%% (R$ %s)

--------------------

[VALOR LÍQUIDO]
Total valor líquido: R$ %s
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

`,
		FormatFloatBR(totalAportadoBruto), percRFBruto, FormatFloatBR(ultimoBrutoFinal), percFIIsBruto, FormatFloatBR(fiisBruto),
		FormatFloatBR(totalAportadoLiquido), percRFLiquido, FormatFloatBR(rfLiquido), percFIIsLiquido, FormatFloatBR(fiisLiquido),
		FormatFloatBR(ultimoBrutoFinal), FormatFloatBR(ultimoLiquidoFinal), FormatFloatBR(lucrosRetiradosTotal), FormatFloatBR(lucroBrutoTotal),
		FormatFloatBR(lucroLiquidoAcumulado))

	// [FIIs] bloco global
	resumo += "[FIIs]\n"
	resumo += fmt.Sprintf("Total Investido: R$ %s\n", FormatFloatBR(fiisTotalInvestido))
	carteiraFIIs := fiisTotalInvestido + dados.ValorAjusteFIIs
	resumo += fmt.Sprintf("Carteira: R$ %s\n", FormatFloatBR(carteiraFIIs))
	sinalAjuste := "+"
	if dados.ValorAjusteFIIs < 0 {
		sinalAjuste = "-"
	}
	resumo += fmt.Sprintf("Lucro/Prejuízo: R$ %s%s\n", sinalAjuste, FormatFloatBR(abs(dados.ValorAjusteFIIs)))
	// Linha de rendimento FIIs sozinha
	linhaRendimento := fmt.Sprintf("[Rendimento FIIs: R$ %s]", FormatFloatBR(rendimentoFIIs))
	resumo += linhaRendimento + "\n"
	// Linha de resumo dos FIIs (ex: ' - VGIR11 (100%) | R$ 10,00')
	fiisResumo := ""
	if len(todosFIIs) > 0 && rendimentoFIIs > 0.0 {
		for codigo, fii := range todosFIIs {
			lucroFII := fii.Dividendos
			for _, venda := range fii.Vendas {
				lucroFII += venda.LucroVenda - venda.DARF
			}
			porcentagem := 0.0
			if rendimentoFIIs > 0 {
				porcentagem = (lucroFII / rendimentoFIIs) * 100
			}
			fiisResumo += fmt.Sprintf(" - %s (%.0f%%) | R$ %s\n", codigo, porcentagem, FormatFloatBR(lucroFII))
		}
	}
	if fiisResumo != "" {
		resumo += fiisResumo
	}
	resumo += "\n"
	fiisDetalhes = ""
	if len(todosFIIs) > 0 {
		fiisDetalhes = "[FIIs Detalhados]\n"
		for codigo, fii := range todosFIIs {
			totalQtd := 0
			totalValor := 0.0
			for _, aporte := range fii.Aportes {
				totalQtd += aporte.Quantidade
				if aporte.ValorTotalManual != nil {
					totalValor += *aporte.ValorTotalManual
				} else {
					totalValor += aporte.ValorTotal
				}
			}
			precoMedio := CalcularPrecoMedioFII(fii)
			fiisDetalhes += fmt.Sprintf("  - %s: %d cotas (R$ %s) | Preço médio: R$ %s\n", codigo, totalQtd, FormatFloatBR(totalValor), FormatFloatBR(precoMedio))
		}
	}
	if fiisDetalhes != "" {
		resumo += fiisDetalhes
	}

	resumo += alertaDARF
	resumo += fmt.Sprintf(`
------------------------------------------------------

╔════════════════════════════════════════════════════╗
║  Lucro Total Bruto (RF + FIIs): R$ %s           ║
║  Lucro Total Líquido (RF + FIIs): R$ %s           ║
╚════════════════════════════════════════════════════╝
`,
		FormatFloatBR(lucroBrutoTotalAcumulado),
		FormatFloatBR(lucroMesLiquidoTotalAcumulado))

	return resumo
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
				lucroLiquidoFIIs := CalcularLucroLiquidoFIIs(m.FIIs)
				lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
				aporteFIIs := CalcularValorTotalFIIs(m.FIIs)
				titulo := fmt.Sprintf("Mês: %s/%s", NomeMes(mes), ano)
				resumo := fmt.Sprintf(`%s
  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)
---------------------------------------
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
  %-22s R$ %s
---------------------------------------`,
					titulo,
					"Aporte Total:", FormatFloatBR(m.AporteRF+aporteFIIs),
					"Aporte RF:", FormatFloatBR(m.AporteRF),
					"Aporte FIIs:", FormatFloatBR(aporteFIIs),
					"Saída:", FormatFloatBR(m.Saida),
					"Lucro Retirado:", FormatFloatBR(m.LucroRetirado),
					"Bruto RF:", FormatFloatBR(m.ValorBrutoRF),
					"Líquido RF:", FormatFloatBR(m.ValorLiquidoRF),
					"Lucro Mês Bruto RF:", FormatFloatBR(lucroMesBruto),
					"Lucro Líquido RF:", FormatFloatBR(lucroMesLiquidoRF),
					"Lucro FIIs:", FormatFloatBR(lucroLiquidoFIIs),
					"Lucro Mês Líquido:", FormatFloatBR(lucroMesLiquidoTotal))
				if len(m.FIIs) > 0 {
					fiisStr := "\n  FIIs do mês:\n"
					for _, fii := range m.FIIs {
						totalQtd := 0
						totalValor := 0.0
						for _, aporte := range fii.Aportes {
							totalQtd += aporte.Quantidade
							if aporte.ValorTotalManual != nil {
								totalValor += *aporte.ValorTotalManual
							} else {
								totalValor += aporte.ValorTotal
							}
						}
						cotasVendidas := 0
						for _, venda := range fii.Vendas {
							cotasVendidas += venda.Quantidade
						}
						totalQtdOriginal := 0
						for _, aporte := range fii.Aportes {
							// Soma a quantidade original do aporte (quantidade atual + cotas vendidas daquele aporte)
							qtdVendidaAporte := 0
							for _, venda := range fii.Vendas {
								if venda.AporteData == aporte.Data {
									qtdVendidaAporte += venda.Quantidade
								}
							}
							totalQtdOriginal += aporte.Quantidade + qtdVendidaAporte
						}
						if cotasVendidas > 0 {
							fiisStr += fmt.Sprintf("    %s: %d cotas atuais | %d cotas vendidas | %d cotas | R$ %s\n", fii.Codigo, totalQtd, cotasVendidas, totalQtdOriginal, FormatFloatBR(totalValor))
						} else {
							fiisStr += fmt.Sprintf("    %s: %d cotas | R$ %s\n", fii.Codigo, totalQtdOriginal, FormatFloatBR(totalValor))
						}
						for _, aporte := range fii.Aportes {
							// Exibir data completa (dd/mm/aaaa)
							data := aporte.Data
							// Calcular quantidade original do aporte
							qtdVendidaAporte := 0
							for _, venda := range fii.Vendas {
								if venda.AporteData == aporte.Data {
									qtdVendidaAporte += venda.Quantidade
								}
							}
							quantidadeOriginal := aporte.Quantidade + qtdVendidaAporte
							fiisStr += fmt.Sprintf("      Aporte (%s): | %d cotas | R$ %s/cota | R$ %s\n",
								data,
								quantidadeOriginal,
								FormatFloatBR(aporte.PrecoCota),
								FormatFloatBR(aporte.ValorTotal),
							)
						}
						for _, venda := range fii.Vendas {
							msg := fmt.Sprintf("      Venda (%s): | %d cotas | Preço médio: R$ %s | Preço total da venda: R$ %s | Taxas: R$ %s",
								venda.Data,
								venda.Quantidade,
								FormatFloatBR(venda.PrecoVenda),
								FormatFloatBR(venda.ValorTotal),
								FormatFloatBR(venda.Taxas),
							)
							if venda.DARF > 0 {
								msg += fmt.Sprintf(" | DARF: R$ %s", FormatFloatBR(venda.DARF))
							}
							fiisStr += msg + "\n"
						}
					}
					return resumo + fiisStr + "---------------------------------------"
				}
				return resumo + "---------------------------------------"
			}
			saldoAnterior = mesesMap[mes].ValorBrutoRF
		}
	}
	return "Mês atual não possui dados."
}

func MostrarResumoAno(dados Dados, ano string) {
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
	valorLiquidoFinal := 0.0
	lucroLiquidoAcumulado := 0.0
	lucroLiquidoFIIsAcumulado := 0.0
	lucroMesLiquidoTotalAcumulado := 0.0
	lucroBrutoTotalAcumulado := 0.0
	saldoAnterior := 0.0
	hoje := time.Now()
	mesAtual := fmt.Sprintf("%02d", int(hoje.Month()))
	anoAtual := fmt.Sprintf("%04d", hoje.Year())

	// Para detalhes dos FIIs do ano
	fiisAno := make(map[string]FII)

	fmt.Printf("\n📌 Resumo dos aportes e saldos mensais - Ano %s\n", ano)

	for _, mes := range meses {
		m := mesesMap[mes]
		lucroMesBruto := m.ValorBrutoRF - (saldoAnterior + m.AporteRF - m.Saida)
		impostos := m.ValorBrutoRF - m.ValorLiquidoRF
		lucroMesLiquidoRF := lucroMesBruto - impostos - m.LucroRetirado
		lucroLiquidoFIIs := CalcularLucroLiquidoFIIs(m.FIIs)
		lucroMesLiquidoTotal := lucroMesLiquidoRF + lucroLiquidoFIIs
		isMesAtual := (ano == anoAtual && mes == mesAtual)

		// Acumular FIIs do ano
		for _, fii := range m.FIIs {
			if fiiExistente, existe := fiisAno[fii.Codigo]; existe {
				fiiExistente.Aportes = append(fiiExistente.Aportes, fii.Aportes...)
				fiisAno[fii.Codigo] = fiiExistente
			} else {
				fiisAno[fii.Codigo] = fii
			}
		}

		fmt.Printf("\nMês: %s/%s\n", NomeMes(mes), ano)
		if isMesAtual {
			fmt.Println("  ⚠️ Mês atual em andamento — valores podem parecer distorcidos (lucro líquido ainda parcial)")
		}
		impostoValido := impostos > 0
		if lucroMesBruto > impostos && impostoValido {
			fmt.Println("  ✅ Agora os lucros já cobrem os impostos!")
		}
		fmt.Println("---------------------------------------")
		aporteFIIs := CalcularValorTotalFIIs(m.FIIs)
		fmt.Printf("  %-22s R$ %s\n", "Aporte Total:", FormatFloatBR(m.AporteRF+aporteFIIs))
		fmt.Printf("  %-22s R$ %s\n", "Aporte RF:", FormatFloatBR(m.AporteRF))
		fmt.Printf("  %-22s R$ %s\n", "Aporte FIIs:", FormatFloatBR(aporteFIIs))
		fmt.Printf("  %-22s R$ %s\n", "Saída:", FormatFloatBR(m.Saida))
		fmt.Printf("  %-22s R$ %s\n", "Lucro Retirado:", FormatFloatBR(m.LucroRetirado))
		fmt.Printf("  %-22s R$ %s\n", "Bruto RF:", FormatFloatBR(m.ValorBrutoRF))
		fmt.Printf("  %-22s R$ %s\n", "Líquido RF:", FormatFloatBR(m.ValorLiquidoRF))
		fmt.Printf("  %-22s R$ %s\n", "Lucro Mês Bruto RF:", FormatFloatBR(lucroMesBruto))
		fmt.Printf("  %-22s R$ %s\n", "Lucro Líquido RF:", FormatFloatBR(lucroMesLiquidoRF))
		fmt.Printf("  %-22s R$ %s\n", "Lucro FIIs:", FormatFloatBR(lucroLiquidoFIIs))
		fmt.Printf("  %-22s R$ %s\n", "Lucro Mês Líquido:", FormatFloatBR(lucroMesLiquidoTotal))

		// Mostrar detalhes dos FIIs do mês se houver
		if len(m.FIIs) > 0 {
			fmt.Println("  FIIs do mês:")
			for _, fii := range m.FIIs {
				totalQtd := 0
				totalValor := 0.0
				for _, aporte := range fii.Aportes {
					totalQtd += aporte.Quantidade
					if aporte.ValorTotalManual != nil {
						totalValor += *aporte.ValorTotalManual
					} else {
						totalValor += aporte.ValorTotal
					}
				}
				cotasVendidas := 0
				for _, venda := range fii.Vendas {
					cotasVendidas += venda.Quantidade
				}
				totalQtdOriginal := 0
				for _, aporte := range fii.Aportes {
					// Soma a quantidade original do aporte (quantidade atual + cotas vendidas daquele aporte)
					qtdVendidaAporte := 0
					for _, venda := range fii.Vendas {
						if venda.AporteData == aporte.Data {
							qtdVendidaAporte += venda.Quantidade
						}
					}
					totalQtdOriginal += aporte.Quantidade + qtdVendidaAporte
				}
				if cotasVendidas > 0 {
					fmt.Printf("    %s: %d cotas atuais | %d cotas vendidas | %d cotas | R$ %s\n", fii.Codigo, totalQtd, cotasVendidas, totalQtdOriginal, FormatFloatBR(totalValor))
				} else {
					fmt.Printf("    %s: %d cotas | R$ %s\n", fii.Codigo, totalQtdOriginal, FormatFloatBR(totalValor))
				}
				for _, aporte := range fii.Aportes {
					// Exibir data completa (dd/mm/aaaa)
					data := aporte.Data
					// Calcular quantidade original do aporte
					qtdVendidaAporte := 0
					for _, venda := range fii.Vendas {
						if venda.AporteData == aporte.Data {
							qtdVendidaAporte += venda.Quantidade
						}
					}
					quantidadeOriginal := aporte.Quantidade + qtdVendidaAporte
					fmt.Printf("      Aporte (%s): | %d cotas | R$ %s/cota | R$ %s\n",
						data,
						quantidadeOriginal,
						FormatFloatBR(aporte.PrecoCota),
						FormatFloatBR(aporte.ValorTotal),
					)
				}
				for _, venda := range fii.Vendas {
					msg := fmt.Sprintf("      Venda (%s): | %d cotas | Preço médio: R$ %s | Preço total da venda: R$ %s | Taxas: R$ %s",
						venda.Data,
						venda.Quantidade,
						FormatFloatBR(venda.PrecoVenda),
						FormatFloatBR(venda.ValorTotal),
						FormatFloatBR(venda.Taxas),
					)
					if venda.DARF > 0 {
						msg += fmt.Sprintf(" | DARF: R$ %s", FormatFloatBR(venda.DARF))
					}
					fmt.Printf("      %s\n", msg)
				}
			}
		}
		fmt.Println("---------------------------------------")

		// Acumular valores (sem filtro para o ano selecionado)
		aporteRFSoFar += m.AporteRF
		aporteFIIsSoFar += CalcularValorTotalFIIs(m.FIIs)
		saidaSoFar += m.Saida
		lucrosRetiradosTotal += m.LucroRetirado
		valorBrutoFinal = m.ValorBrutoRF
		valorLiquidoFinal = m.ValorLiquidoRF

		// Lucro bruto sempre acumula (todos os meses)
		lucroBrutoTotalAcumulado += lucroMesBruto

		// Lucro líquido só acumula se for válido (lucro cobre imposto)
		lucroValido := lucroMesBruto > impostos
		if isMesAtual && !lucroValido {
			// Não acumular lucros líquidos do mês atual se não for válido
			// Mas continuar acumulando FIIs
			lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
		} else if lucroValido {
			lucroLiquidoAcumulado += lucroMesLiquidoRF
			lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
			lucroMesLiquidoTotalAcumulado += lucroMesLiquidoTotal
		} else {
			// Meses passados sempre acumulam FIIs
			lucroLiquidoFIIsAcumulado += lucroLiquidoFIIs
		}

		saldoAnterior = m.ValorBrutoRF
	}

	// Totais - usar valores finais (bruto e líquido) da RF
	totalAportadoBruto := valorBrutoFinal + aporteFIIsSoFar
	totalAportadoLiquido := valorLiquidoFinal + aporteFIIsSoFar
	// Lucro bruto total = valor final - total aportado bruto (sem considerar saídas no cálculo)
	lucroBrutoTotal := valorBrutoFinal - totalAportadoBruto
	// Corrigir: usar o acumulado dos lucros brutos
	lucroBrutoTotal = lucroBrutoTotalAcumulado
	// Usar o acumulado correto dos lucros líquidos (não sobrescrever)
	// lucroLiquidoTotal := valorLiquidoFinal - (aporteRFSoFar - saidaSoFar)
	lucroLiquidoTotal := lucroLiquidoAcumulado

	// Calcular porcentagens

	// Preparar detalhes dos FIIs com porcentagem do lucro
	fiisDetalhes := ""
	if len(fiisAno) > 0 {
		fiisDetalhes = "\n[FIIs Detalhados do Ano]\n"
		for codigo, fii := range fiisAno {
			totalQtd := 0
			totalValor := 0.0
			for _, aporte := range fii.Aportes {
				totalQtd += aporte.Quantidade
				if aporte.ValorTotalManual != nil {
					totalValor += *aporte.ValorTotalManual
				} else {
					totalValor += aporte.ValorTotal
				}
			}
			// Exibir apenas: - CÓDIGO: N cotas (R$ X,XX)
			fiisDetalhes += fmt.Sprintf("  - %s: %d cotas (R$ %s)\n", codigo, totalQtd, FormatFloatBR(totalValor))
		}
	}

	// Cálculo do bloco [FIIs] do ANO (igual ao global, mas só com FIIs do ano)
	fiisTotalInvestido := 0.0
	fiisDividendos := 0.0
	fiisLucroVendas := 0.0
	for _, fii := range fiisAno {
		for _, aporte := range fii.Aportes {
			if aporte.ValorTotalManual != nil {
				fiisTotalInvestido += *aporte.ValorTotalManual
			} else {
				fiisTotalInvestido += aporte.ValorTotal
			}
		}
		fiisDividendos += fii.Dividendos
		for _, venda := range fii.Vendas {
			fiisLucroVendas += venda.LucroVenda
		}
	}
	// O ajuste manual é global, mas entra no cálculo da carteira do ano
	fiisCarteira := fiisTotalInvestido + fiisDividendos + fiisLucroVendas + dados.ValorAjusteFIIs
	rendimentoFIIs := fiisDividendos + fiisLucroVendas
	if rendimentoFIIs < 0.005 && rendimentoFIIs > -0.005 {
		rendimentoFIIs = 0.0
	}

	fmt.Println()
	fmt.Println("================== InvistAI ==================")
	fmt.Println()
	fmt.Println("--- Total Investido do Ano ---")
	fmt.Println()
	fmt.Println("[VALOR BRUTO (valor atual da carteira no ano)]")
	fmt.Printf("Total valor bruto: R$ %s\n", FormatFloatBR(totalAportadoBruto))
	totalBrutoAno := valorBrutoFinal + fiisCarteira
	percRFBrutoAno := 0.0
	percFIIsBrutoAno := 0.0
	if totalBrutoAno > 0 {
		percRFBrutoAno = (valorBrutoFinal / totalBrutoAno) * 100
		percFIIsBrutoAno = (fiisCarteira / totalBrutoAno) * 100
	}
	fmt.Printf("  - Renda Fixa: %.2f%% (R$ %s)\n", percRFBrutoAno, FormatFloatBR(valorBrutoFinal))
	fmt.Printf("  - FIIs: %.2f%% (R$ %s)\n", percFIIsBrutoAno, FormatFloatBR(fiisCarteira))
	fmt.Println()
	fmt.Println("--------------------")
	fmt.Println()
	fmt.Println("[VALOR LÍQUIDO]")
	fmt.Printf("Total valor líquido: R$ %s\n", FormatFloatBR(totalAportadoLiquido))
	totalLiquidoAno := valorLiquidoFinal + fiisCarteira
	percRFLiquidoAno := 0.0
	percFIIsLiquidoAno := 0.0
	if totalLiquidoAno > 0 {
		percRFLiquidoAno = (valorLiquidoFinal / totalLiquidoAno) * 100
		percFIIsLiquidoAno = (fiisCarteira / totalLiquidoAno) * 100
	}
	fmt.Printf("  - Renda Fixa: %.2f%% (R$ %s)\n", percRFLiquidoAno, FormatFloatBR(valorLiquidoFinal))
	fmt.Printf("  - FIIs: %.2f%% (R$ %s)\n", percFIIsLiquidoAno, FormatFloatBR(fiisCarteira))
	fmt.Println()
	fmt.Println("---------------------------------------")
	fmt.Println()
	fmt.Println("[RENDA FIXA]")
	fmt.Printf("Valor Bruto Final (RF): R$ %s\n", FormatFloatBR(valorBrutoFinal))
	fmt.Printf("Valor Líquido Final (RF): R$ %s\n", FormatFloatBR(valorLiquidoFinal))
	fmt.Printf("Lucros Retirados: R$ %s\n", FormatFloatBR(lucrosRetiradosTotal))
	fmt.Printf("Lucro Bruto Total (RF): R$ %s\n", FormatFloatBR(lucroBrutoTotal))
	fmt.Printf("Lucro Líquido RF: R$ %s\n", FormatFloatBR(lucroLiquidoTotal))
	fmt.Println()
	fmt.Println("---------------------------------------")
	fmt.Println()
	fmt.Println("[FIIs]")
	fmt.Printf("Total Investido: R$ %s\n", FormatFloatBR(fiisTotalInvestido))
	fmt.Printf("Carteira: R$ %s\n", FormatFloatBR(fiisCarteira))
	// Calcular variacao do ano (carteira FIIs)
	variacaoAno := fiisCarteira - fiisTotalInvestido
	fmt.Printf("Variação (carteira FIIs): R$ %s\n", FormatFloatBR(variacaoAno))
	// Linha de rendimento FIIs sozinha
	linhaRendimento := fmt.Sprintf("[Rendimento FIIs: R$ %s]", FormatFloatBR(rendimentoFIIs))
	fmt.Println(linhaRendimento)
	// Linha de resumo dos FIIs (ex: ' - VGIR11 (100%) | R$ 10,00')
	fiisResumo := ""
	if len(fiisAno) > 0 && rendimentoFIIs > 0.0 {
		for codigo, fii := range fiisAno {
			lucroFII := fii.Dividendos
			for _, venda := range fii.Vendas {
				lucroFII += venda.LucroVenda - venda.DARF
			}
			porcentagem := 0.0
			if rendimentoFIIs > 0 {
				porcentagem = (lucroFII / rendimentoFIIs) * 100
			}
			fiisResumo += fmt.Sprintf(" - %s (%.0f%%) | R$ %s\n", codigo, porcentagem, FormatFloatBR(lucroFII))
		}
	}
	if fiisResumo != "" {
		fmt.Print(fiisResumo)
	}
	if fiisDetalhes != "" {
		fmt.Print(fiisDetalhes)
	}
	fmt.Println()
	fmt.Println("---------------------------------------")
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Printf("║  Lucro Total Bruto (RF + FIIs): R$ %s           ║\n", FormatFloatBR(lucroBrutoTotalAcumulado))
	fmt.Printf("║  Lucro Total Líquido (RF + FIIs): R$ %s           ║\n", FormatFloatBR(lucroMesLiquidoTotalAcumulado))
	fmt.Println("╚════════════════════════════════════════════════════╝")
	return
}
