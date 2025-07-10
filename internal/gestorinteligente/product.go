package gestorinteligente

import "time"

type Produto struct {
	Nome       string    `json:"nome"`
	Parcela    float64   `json:"parcela"`
	ValorTotal float64   `json:"valor_total"`
	Parcelas   int       `json:"parcelas"`
	CriadoEm   time.Time `json:"criado_em"`
}

type ListaProdutos struct {
	Produtos          []Produto `json:"produtos"`
	LucroMensal       float64   `json:"lucro_mensal"`
	Mes               int       `json:"mes"`
	Ano               int       `json:"ano"`
	PorcentagemSegura float64   `json:"porcentagem_segura"`
}
