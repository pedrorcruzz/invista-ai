package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pedrorcruzz/invista-ai/internal"
	"github.com/pedrorcruzz/invista-ai/internal/gestorinteligente"
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

		if inMenuInicial && opcao == "5" {
			fmt.Println("Saindo...")
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
				internal.MostrarResumoAno(dados, ano, false)
				fmt.Print("\nPressione Enter para voltar ao menu...")
				scanner.Scan()
			}
			internal.PrintMenuPrincipalSozinho()
			inMenuInicial = false
		case "2":
			ano := internal.SelecionarAno(dados, scanner)
			if ano != "" {
				internal.MostrarResumoAno(dados, ano, true)
				fmt.Print("\nPressione Enter para voltar ao menu...")
				scanner.Scan()
			}
			internal.PrintMenuPrincipalSozinho()
			inMenuInicial = false
		case "3":
			internal.AdicionarOuEditarMes(&dados, scanner)
			internal.SalvarDados(dados)
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		case "4":
			gestorinteligente.ShowGestorMenu()
			internal.PrintTelaUnificada(dados)
			inMenuInicial = true
		default:
			fmt.Println("Opção inválida!")
			internal.PrintMenuPrincipalSozinho()
			inMenuInicial = false
		}
	}
}
