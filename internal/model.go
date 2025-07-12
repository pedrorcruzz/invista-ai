package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type FIIAporte struct {
	Quantidade       int      `json:"quantidade"`
	PrecoCota        float64  `json:"preco_cota"`
	ValorTotal       float64  `json:"valor_total"`
	ValorTotalManual *float64 `json:"valor_total_manual,omitempty"` // Valor manual opcional
	Data             string   `json:"data"`
}

type FIIVenda struct {
	Quantidade int     `json:"quantidade"`
	PrecoVenda float64 `json:"preco_venda"`
	ValorTotal float64 `json:"valor_total"`
	LucroVenda float64 `json:"lucro_venda"`
	DARF       float64 `json:"darf"`
	Data       string  `json:"data"`
	Taxas      float64 `json:"taxas"`
	AporteData string  `json:"aporte_data"`
}

type FII struct {
	Codigo     string      `json:"codigo"`
	Aportes    []FIIAporte `json:"aportes"`
	Dividendos float64     `json:"dividendos"`
	Vendas     []FIIVenda  `json:"vendas"`
}

type Mes struct {
	AporteRF       float64 `json:"aporte_rf"`
	Saida          float64 `json:"saida"`
	ValorBrutoRF   float64 `json:"valor_bruto_rf"`
	ValorLiquidoRF float64 `json:"valor_liquido_rf"`
	LucroRetirado  float64 `json:"lucro_retirado"`
	FIIs           []FII   `json:"fiis"`
}

type Ano map[string]Mes

type Dados struct {
	Anos           map[string]Ano `json:"anos"`
	FIIsConhecidos []string       `json:"fiis_conhecidos"`
}

const Arquivo = "dados.json"

func CarregarDados() Dados {
	file, err := os.ReadFile(Arquivo)
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

func SalvarDados(dados Dados) {
	bytes, err := json.MarshalIndent(dados, "", "  ")
	if err != nil {
		fmt.Println("Erro ao salvar dados:", err)
		return
	}
	os.WriteFile(Arquivo, bytes, 0644)
}

// Funções auxiliares para FIIs
func AdicionarFIIConhecido(dados *Dados, codigo string) {
	codigo = strings.ToUpper(strings.TrimSpace(codigo))
	if codigo == "" {
		return
	}

	// Verificar se já existe
	for _, fii := range dados.FIIsConhecidos {
		if fii == codigo {
			return
		}
	}

	dados.FIIsConhecidos = append(dados.FIIsConhecidos, codigo)
	// Ordenar a lista
	sort.Strings(dados.FIIsConhecidos)
}

func EncontrarFIIPorCodigo(fiis []FII, codigo string) *FII {
	for i := range fiis {
		if fiis[i].Codigo == codigo {
			return &fiis[i]
		}
	}
	return nil
}

func CalcularValorTotalFIIs(fiis []FII) float64 {
	total := 0.0
	for _, fii := range fiis {
		for _, aporte := range fii.Aportes {
			if aporte.ValorTotalManual != nil {
				total += *aporte.ValorTotalManual
			} else {
				total += aporte.ValorTotal
			}
		}
	}
	return total
}

func CalcularLucroLiquidoFIIs(fiis []FII) float64 {
	total := 0.0
	for _, fii := range fiis {
		// Dividendos
		total += fii.Dividendos

		// Lucro das vendas (já descontado DARF)
		for _, venda := range fii.Vendas {
			total += venda.LucroVenda - venda.DARF
		}
	}
	return total
}

func CalcularDARFTotal(fiis []FII) float64 {
	total := 0.0
	for _, fii := range fiis {
		for _, venda := range fii.Vendas {
			total += venda.DARF
		}
	}
	return total
}

// Calcular DARF para venda de cotas (15% sobre o lucro)
func CalcularDARFVenda(precoVenda, precoMedio float64, quantidade int) float64 {
	lucroPorCota := precoVenda - precoMedio
	if lucroPorCota <= 0 {
		return 0.0
	}
	lucroTotal := lucroPorCota * float64(quantidade)
	return lucroTotal * 0.15 // 15% de DARF
}

// Calcular preço médio das cotas de um FII
func CalcularPrecoMedioFII(fii FII) float64 {
	totalCotas := 0
	totalValor := 0.0

	for _, aporte := range fii.Aportes {
		totalCotas += aporte.Quantidade
		totalValor += aporte.ValorTotal
	}

	if totalCotas == 0 {
		return 0.0
	}

	return totalValor / float64(totalCotas)
}

// Calcular prazo de pagamento do DARF (até o último dia do mês seguinte)
func CalcularPrazoDARF(mes, ano string) (int, int, int) {
	mesInt, _ := strconv.Atoi(mes)
	anoInt, _ := strconv.Atoi(ano)

	// Prazo: até o último dia do mês seguinte
	mesPagamento := mesInt + 1
	anoPagamento := anoInt
	if mesPagamento > 12 {
		mesPagamento = 1
		anoPagamento++
	}

	// Calcular último dia do mês de pagamento
	ultimoDia := 31
	if mesPagamento == 2 {
		// Verificar se é ano bissexto
		if (anoPagamento%4 == 0 && anoPagamento%100 != 0) || (anoPagamento%400 == 0) {
			ultimoDia = 29
		} else {
			ultimoDia = 28
		}
	} else if mesPagamento == 4 || mesPagamento == 6 || mesPagamento == 9 || mesPagamento == 11 {
		ultimoDia = 30
	}

	return ultimoDia, mesPagamento, anoPagamento
}
