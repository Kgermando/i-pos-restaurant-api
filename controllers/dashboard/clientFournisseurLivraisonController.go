package dashboard

import (
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"

	"github.com/gofiber/fiber/v2"
)

// Get total Client et Fournisseurs
func GetTotalClientFournisseur(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var clientCount int64 = 0

	db.Model(&models.Client{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", start_date, end_date). // Filter by date
		Count(&clientCount)

	var fournisseurCount int64 = 0
	db.Model(&models.Fournisseur{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", start_date, end_date). // Filter by date
		Count(&fournisseurCount)

	 

	response := map[string]interface{}{
		"client":     clientCount,
		"fournisseur": fournisseurCount,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total ventes",
		"data":    response,
	})
}


