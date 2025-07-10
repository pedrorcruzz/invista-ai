package gestorinteligente

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dataDir = "data"
const dataFile = "products.json"

func ensureDataDir() error {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return os.Mkdir(dataDir, 0755)
	}
	return nil
}

func LoadProducts() (ProductList, error) {
	var list ProductList
	list.SafePercentage = 70

	if err := ensureDataDir(); err != nil {
		return list, err
	}

	filePath := filepath.Join(dataDir, dataFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return list, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return list, err
	}

	if len(data) == 0 {
		return list, nil
	}

	err = json.Unmarshal(data, &list)
	return list, err
}

func SaveProducts(list ProductList) error {
	if err := ensureDataDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(dataDir, dataFile)
	return os.WriteFile(filePath, data, 0644)
}
