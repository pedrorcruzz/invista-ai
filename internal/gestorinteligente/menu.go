package gestorinteligente

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func obterMenuGestorStr() string {
	return `--- MENU GESTOR INTELIGENTE ---
1. Adicionar produto
2. Remover produto
3. Listar meses
4. Atualizar lucro mensal
5. Editar produto
6. Antecipar parcelas
7. Configurar porcentagem segura
8. Voltar ao menu principal`
}

func obterResumoGestorStr(lista ListaProdutos) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("RESUMO DO MÊS (%02d/%d - %s)", time.Now().Month(), time.Now().Year(), nomesMeses[time.Now().Month()-1]))
	b.WriteString("\n------------------------------")
	b.WriteString(fmt.Sprintf("\nLucro mensal: R$%.2f", lista.LucroMensal))
	b.WriteString(fmt.Sprintf("\nTotal de parcelas: R$%.2f", 0.0)) // Placeholder, pode melhorar
	b.WriteString("\n...")                                         // Pode adicionar mais detalhes se quiser
	return b.String()
}

func obterMensagemGestorStr(lista ListaProdutos) string {
	if lista.LucroMensal == 0 {
		return "Por favor, defina seu lucro mensal antes de adicionar produtos."
	}
	return ""
}

func ImprimirMenuGestorCompleto(lista ListaProdutos) {
	titulo := " Gestor Inteligente de Gastos "
	resumo := capturarResumo(lista)
	mensagem := obterMensagemGestorStr(lista)
	menu := obterMenuGestorStr()

	blocos := [][]string{
		{titulo},
		strings.Split(resumo, "\n"),
	}
	if mensagem != "" {
		blocos = append(blocos, strings.Split(mensagem, "\n"))
	}
	blocos = append(blocos, strings.Split(menu, "\n"))

	maxLen := 0
	for _, bloco := range blocos {
		for _, l := range bloco {
			if len(l) > maxLen {
				maxLen = len(l)
			}
		}
	}
	if maxLen < 60 {
		maxLen = 60
	}

	linhaTopo := "╔" + strings.Repeat("═", maxLen+2) + "╗"
	linhaDiv := "╟" + strings.Repeat("─", maxLen+2) + "╢"
	linhaBase := "╚" + strings.Repeat("═", maxLen+2) + "╝"

	fmt.Println(linhaTopo)
	for i, bloco := range blocos {
		for _, l := range bloco {
			fmt.Printf("║ %-*s ║\n", maxLen, l)
		}
		if i < len(blocos)-1 {
			fmt.Println(linhaDiv)
		}
	}
	fmt.Println(linhaBase)
}

func capturarResumo(lista ListaProdutos) string {
	var b strings.Builder
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	mostrarResumo(lista)
	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	b.Write(buf[:n])
	return strings.TrimRight(b.String(), "\n")
}

func MostrarMenuGestor() {
	reader := bufio.NewReader(os.Stdin)
	lista, _ := CarregarProdutos()

	if lista.PorcentagemSegura == 0 {
		lista.PorcentagemSegura = 70
	}

	if lista.LucroMensal == 0 {
		LimparTerminal()
		ImprimirCaixa([]string{
			"🎯 CONFIGURAÇÃO INICIAL",
			"",
			"Bem-vindo ao Gestor Inteligente de Gastos!",
			"",
			"Para começar a usar o sistema, precisamos definir",
			"seu lucro mensal. Este valor será usado para",
			"calcular quanto você pode gastar em parcelas.",
			"",
			"Digite seu lucro mensal (R$):",
		})
		fmt.Print("→ ")
		lucroStr, _ := reader.ReadString('\n')
		lucroStr = strings.TrimSpace(lucroStr)
		lucroStr = strings.ReplaceAll(lucroStr, ",", ".")

		lucro, err := strconv.ParseFloat(lucroStr, 64)
		if err != nil || lucro <= 0 {
			ImprimirCaixa([]string{
				"❌ Valor inválido!",
				"",
				"Por favor, digite um valor válido maior que zero.",
				"Exemplo: 5000 ou 5000,50",
			})
			time.Sleep(3 * time.Second)
			return
		}

		lista.LucroMensal = lucro
		SalvarProdutos(lista)

		ImprimirCaixa([]string{
			"✅ Lucro mensal configurado com sucesso!",
			"",
			fmt.Sprintf("Seu lucro mensal: R$%.2f", lucro),
			"",
			"Agora você pode começar a adicionar produtos.",
		})
		time.Sleep(3 * time.Second)
	}

	for {
		LimparTerminal()
		ImprimirMenuGestorCompleto(lista)
		fmt.Print("Escolha uma opção: ")
		escolha, _ := reader.ReadString('\n')
		escolha = strings.TrimSpace(escolha)

		switch escolha {
		case "1":
			LimparTerminal()
			adicionarProduto(reader, &lista)
		case "2":
			LimparTerminal()
			removerProduto(reader, &lista)
		case "3":
			LimparTerminal()
			listarMeses(reader, lista)
		case "4":
			LimparTerminal()
			atualizarLucroMensal(reader, &lista)
		case "5":
			LimparTerminal()
			editarProduto(reader, &lista)
		case "6":
			LimparTerminal()
			anteciparParcelas(reader, &lista)
		case "7":
			LimparTerminal()
			configurarPorcentagemSegura(reader, &lista)
		case "8":
			SalvarProdutos(lista)
			return
		default:
			fmt.Println("Opção inválida.")
			time.Sleep(1 * time.Second)
		}
		SalvarProdutos(lista)
	}
}
