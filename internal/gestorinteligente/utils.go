package gestorinteligente

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var nomesMeses = []string{
	"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
	"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
}

func ImprimirCaixa(linhas []string) {
	maxLen := 0
	for _, l := range linhas {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	if maxLen < 60 {
		maxLen = 60
	}
	linhaTopo := "╔" + strings.Repeat("═", maxLen+2) + "╗"
	linhaBase := "╚" + strings.Repeat("═", maxLen+2) + "╝"
	fmt.Println(linhaTopo)
	for _, l := range linhas {
		fmt.Printf("║ %-*s ║\n", maxLen, l)
	}
	fmt.Println(linhaBase)
}

func LimparTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func mapearProdutosPorAnoMes(produtos []Produto) map[int]map[int][]int {
	resultado := make(map[int]map[int][]int)
	for idx, p := range produtos {
		dataFim := p.CriadoEm.AddDate(0, p.Parcelas-1, 0)
		dataAtual := p.CriadoEm

		for !dataAtual.After(dataFim) {
			anoAtual := dataAtual.Year()
			mesAtual := int(dataAtual.Month())

			if _, ok := resultado[anoAtual]; !ok {
				resultado[anoAtual] = make(map[int][]int)
			}
			resultado[anoAtual][mesAtual] = append(resultado[anoAtual][mesAtual], idx)

			dataAtual = dataAtual.AddDate(0, 1, 0)
		}
	}
	return resultado
}

func selecionarProdutoPorAnoMes(reader *bufio.Reader, produtos []Produto) (int, bool) {
	porAnoMes := mapearProdutosPorAnoMes(produtos)
	if len(porAnoMes) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return -1, false
	}

	anos := make([]int, 0, len(porAnoMes))
	for y := range porAnoMes {
		anos = append(anos, y)
	}
	sort.Ints(anos)

	fmt.Println("\nSelecione o ano (0 para voltar):")
	for i, y := range anos {
		fmt.Printf("%d. %d\n", i+1, y)
	}
	fmt.Print("Ano: ")
	anoStr, _ := reader.ReadString('\n')
	anoStr = strings.TrimSpace(anoStr)

	if anoStr == "0" {
		return -1, false
	}

	anoIdx, err := strconv.Atoi(anoStr)
	if err != nil || anoIdx < 1 || anoIdx > len(anos) {
		fmt.Println("Ano inválido.")
		return -1, false
	}
	ano := anos[anoIdx-1]

	mesesMap := porAnoMes[ano]
	meses := make([]int, 0, len(mesesMap))
	for m := range mesesMap {
		meses = append(meses, m)
	}
	sort.Ints(meses)

	fmt.Println("\nSelecione o mês (0 para voltar):")
	for i, m := range meses {
		fmt.Printf("%d. %s\n", i+1, nomesMeses[m-1])
	}
	fmt.Print("Mês: ")
	mesStr, _ := reader.ReadString('\n')
	mesStr = strings.TrimSpace(mesStr)

	if mesStr == "0" {
		return -1, false
	}

	mesIdx, err := strconv.Atoi(mesStr)
	if err != nil || mesIdx < 1 || mesIdx > len(meses) {
		fmt.Println("Mês inválido.")
		return -1, false
	}
	mes := meses[mesIdx-1]

	idxProdutos := mesesMap[mes]

	idxUnicos := make([]int, 0)
	vistos := make(map[int]bool)
	for _, idx := range idxProdutos {
		if !vistos[idx] {
			vistos[idx] = true
			idxUnicos = append(idxUnicos, idx)
		}
	}

	fmt.Println("\nSelecione o produto (0 para voltar):")
	for i, idx := range idxUnicos {
		p := produtos[idx]
		fmt.Printf("%d. %s | Total: R$%.2f | Parcelas: %d | Adicionado em: %s\n",
			i+1, p.Nome, p.ValorTotal, p.Parcelas, p.CriadoEm.Format("02/01/2006"))
	}
	fmt.Print("Produto: ")
	prodStr, _ := reader.ReadString('\n')
	prodStr = strings.TrimSpace(prodStr)

	if prodStr == "0" {
		return -1, false
	}

	prodIdx, err := strconv.Atoi(prodStr)
	if err != nil || prodIdx < 1 || prodIdx > len(idxUnicos) {
		fmt.Println("Produto inválido.")
		return -1, false
	}
	return idxUnicos[prodIdx-1], true
}

func lerFloat(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Print(prompt)
	valorStr, _ := reader.ReadString('\n')
	valorStr = strings.TrimSpace(valorStr)
	valorStr = strings.ReplaceAll(valorStr, ",", ".")
	return strconv.ParseFloat(valorStr, 64)
}

func produtoAtivoNoMes(p Produto, anoAlvo, mesAlvo int) bool {
	dataInicio := p.CriadoEm
	dataFim := dataInicio.AddDate(0, p.Parcelas-1, 0)

	dataAlvo := time.Date(anoAlvo, time.Month(mesAlvo), 1, 0, 0, 0, 0, time.UTC)
	return !dataAlvo.Before(dataInicio) && !dataAlvo.After(dataFim)
}

func obterNumeroParcela(p Produto, anoAlvo, mesAlvo int) int {
	dataInicio := p.CriadoEm
	anoDif := anoAlvo - dataInicio.Year()
	mesDif := mesAlvo - int(dataInicio.Month())
	difTotalMes := anoDif*12 + mesDif + 1

	if difTotalMes < 1 {
		return 1
	}
	if difTotalMes > p.Parcelas {
		return p.Parcelas
	}
	return difTotalMes
}
