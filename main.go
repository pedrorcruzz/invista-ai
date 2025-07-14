package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pedrorcruzz/invista-ai-cli/internal"
	"github.com/pedrorcruzz/invista-ai-cli/internal/gestorinteligente"
)

func main() {
	dados := internal.CarregarDados()
	scanner := bufio.NewScanner(os.Stdin)

	internal.PrintTelaUnificada(dados)
	inMenuInicial := true

	for {
		fmt.Print("Escolha uma opção: ")
		scanner.Scan()
		opcao := scanner.Text()

		if inMenuInicial && opcao == "6" {
			internal.ClearTerminal()
			return
		}

		if !inMenuInicial && opcao == "5" {
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
			continue
		}

		switch opcao {
		case "1":
			ano := internal.SelecionarAno(dados, scanner)
			if ano != "" {
				internal.MostrarResumoAno(dados, ano)
				fmt.Print("\nPressione Enter para voltar ao menu...")
				scanner.Scan()
			}
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		case "2":
			internal.GerenciarRendaFixa(&dados, scanner)
			internal.SalvarDados(dados)
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		case "3":
			internal.GerenciarFIIs(&dados, scanner)
			internal.SalvarDados(dados)
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		case "4":
			gestorinteligente.MostrarMenuGestor()
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		case "5":
			internal.RetirarLucro(&dados, scanner)
			internal.SalvarDados(dados)
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		default:
			fmt.Println("Opção inválida!")
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		}
	}
}
