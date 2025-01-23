package usecases

import (
	"fmt"
	"inventory/pkg/models"
	"mime"
	"os"
	"path/filepath"
)

func (app *Application) CreatePurchase(factors []models.Factor) (models.Purchase, error) {
	for _, v := range factors {
		if v.FileName.Valid {
			if err := ValidateImageFile(v.FileName.String); err != nil {
				return models.Purchase{}, err
			}
		}
	}
	return app.db.CreatePurchase(factors)
}

func (app *Application) GetAllPurchases() ([]models.Purchase, error) {
	return app.db.GetAllPurchases()
}

func (app *Application) GetAllFactors() ([]models.Factor, error) {

	return app.db.GetAllFactors()
}

func ValidateImageFile(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	if err != nil {
		return fmt.Errorf("error accessing file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return fmt.Errorf("invalid file type: %s (only JPG and PNG are allowed)", mimeType)
	}

	return nil
}
