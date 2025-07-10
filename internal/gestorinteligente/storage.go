package gestorinteligente

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const pastaDados = "data"
const arquivoDados = "produtos.json"

func garantirPastaDados() error {
	if _, err := os.Stat(pastaDados); os.IsNotExist(err) {
		return os.Mkdir(pastaDados, 0755)
	}
	return nil
}

func CarregarProdutos() (ListaProdutos, error) {
	var lista ListaProdutos
	lista.PorcentagemSegura = 70

	if err := garantirPastaDados(); err != nil {
		return lista, err
	}

	caminhoArquivo := filepath.Join(pastaDados, arquivoDados)
	if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
		return lista, nil
	}

	dados, err := os.ReadFile(caminhoArquivo)
	if err != nil {
		return lista, err
	}

	if len(dados) == 0 {
		return lista, nil
	}

	err = json.Unmarshal(dados, &lista)
	return lista, err
}

func SalvarProdutos(lista ListaProdutos) error {
	if err := garantirPastaDados(); err != nil {
		return err
	}

	dados, err := json.MarshalIndent(lista, "", "  ")
	if err != nil {
		return err
	}

	caminhoArquivo := filepath.Join(pastaDados, arquivoDados)
	return os.WriteFile(caminhoArquivo, dados, 0644)
}
