package internal

import (
	"encoding/json"
	"fmt"
	"os"
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

