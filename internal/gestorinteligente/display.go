package gestorinteligente

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func mostrarResumo(lista ListaProdutos) {
	var totalParcela float64

	agora := time.Now()
	anoAlvo := agora.Year()
	mesAlvo := int(agora.Month())

	var produtosAtivos []Produto

	for _, p := range lista.Produtos {
		anoInicio, mesInicio := p.CriadoEm.Year(), int(p.CriadoEm.Month())
		dataFim := p.CriadoEm.AddDate(0, p.Parcelas-1, 0)
		anoFim, mesFim := dataFim.Year(), int(dataFim.Month())

		if (anoAlvo > anoInicio || (anoAlvo == anoInicio && mesAlvo >= mesInicio)) &&
			(anoAlvo < anoFim || (anoAlvo == anoFim && mesAlvo <= mesFim)) {
			produtosAtivos = append(produtosAtivos, p)
			totalParcela += p.Parcela
		}
	}

	percentualUsado := 0.0
	percentualRestante := 100.0
	valorReinvestir := 0.0

	if lista.LucroMensal > 0 {
		percentualUsado = (totalParcela / lista.LucroMensal) * 100
		percentualRestante = 100 - percentualUsado
		valorReinvestir = (percentualRestante / 100) * lista.LucroMensal
	}

	percentualGastar := 100.0 - lista.PorcentagemSegura
	valorGastar := (percentualGastar / 100) * lista.LucroMensal
	valorRestanteGastar := valorGastar - totalParcela
	if valorRestanteGastar < 0 {
		valorRestanteGastar = 0
	}

	// nomeMes removido pois não é mais usado

	divisorResumo := strings.Repeat("-", 60)
	tituloPrograma := "================== InvistAI =================="
	fmt.Println("\n" + tituloPrograma)
	fmt.Println(divisorResumo)

	fmt.Printf("Lucro mensal: R$%.2f\n", lista.LucroMensal)
	fmt.Printf("Total de parcelas: R$%.2f\n", totalParcela)
	fmt.Printf("Usado: %.2f%% | Para reinvestir: %.2f%% (R$%.2f)\n", percentualUsado, percentualRestante, valorReinvestir)
	fmt.Printf("Porcentagem segura configurada: %.0f%%\n", lista.PorcentagemSegura)
	fmt.Printf("Disponível para gastos: %.0f%% (R$%.2f) | Restante: R$%.2f\n",
		percentualGastar, valorGastar, valorRestanteGastar)

	fmt.Println("")
	fmt.Println(divisorResumo)

	if percentualRestante >= lista.PorcentagemSegura {
		fmt.Println("✅ Você pode usar parte do seu lucro para pagar as parcelas!")
	} else {
		fmt.Println("❌ Não recomendado. Crie uma caixinha separada para alguns produtos!")
		sugerirProdutosParaSeparar(produtosAtivos, lista.LucroMensal, lista.PorcentagemSegura)
	}

	if len(produtosAtivos) > 0 {
		tituloProdutos := " PRODUTOS ATIVOS NESTE MÊS "
		fmt.Println("\n" + divisorResumo)
		fmt.Println(tituloProdutos)
		fmt.Println(divisorResumo)

		for i, p := range produtosAtivos {
			numeroParcela := obterNumeroParcela(p, anoAlvo, mesAlvo)
			fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d)\n",
				i+1, p.Nome, p.ValorTotal, p.Parcela, numeroParcela, p.Parcelas)
		}
		fmt.Println(divisorResumo)
	}
}

func sugerirProdutosParaSeparar(produtos []Produto, lucroMensal float64, porcentagemSegura float64) {
	if len(produtos) == 0 {
		return
	}
	type ProdutoComIndice struct {
		Indice  int
		Produto Produto
	}

	produtosComIndice := make([]ProdutoComIndice, len(produtos))
	for i, p := range produtos {
		produtosComIndice[i] = ProdutoComIndice{i, p}
	}

	sort.Slice(produtosComIndice, func(i, j int) bool {
		return produtosComIndice[i].Produto.Parcela > produtosComIndice[j].Produto.Parcela
	})

	totalParcela := 0.0
	for _, p := range produtos {
		totalParcela += p.Parcela
	}
	targetParcela := totalParcela - (lucroMensal * (porcentagemSegura / 100))
	if targetParcela <= 0 {
		return
	}

	var produtosSugeridos []Produto
	var somaParcelasSugeridas float64

	for _, pci := range produtosComIndice {
		if somaParcelasSugeridas >= targetParcela {
			break
		}
		produtosSugeridos = append(produtosSugeridos, pci.Produto)
		somaParcelasSugeridas += pci.Produto.Parcela
	}

	divisorSugestao := strings.Repeat("-", 50)

	if len(produtosSugeridos) == 1 {
		fmt.Println(divisorSugestao)
		fmt.Printf("💡 Sugestão: Separe o produto '%s' (Parcela: R$%.2f) em uma caixinha separada.\n",
			produtosSugeridos[0].Nome, produtosSugeridos[0].Parcela)
		fmt.Println(divisorSugestao)
	} else if len(produtosSugeridos) > 1 {
		fmt.Println(divisorSugestao)
		fmt.Println("💡 Sugestão: Separe os seguintes produtos em uma caixinha:")
		for i, p := range produtosSugeridos {
			fmt.Printf("  %d. %s (Parcela: R$%.2f)\n", i+1, p.Nome, p.Parcela)
		}
		fmt.Printf("  Total a separar: R$%.2f\n", somaParcelasSugeridas)
		fmt.Println(divisorSugestao)
	}
}
