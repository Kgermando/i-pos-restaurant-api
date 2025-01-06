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

	var clientCount int64 = 0

	db.Model(&models.Client{}).
		Where("code_entreprise = ?", codeEntreprise).
		Count(&clientCount)

	var fournisseurCount int64 = 0
	db.Model(&models.Fournisseur{}).
		Where("code_entreprise = ?", codeEntreprise).
		Count(&fournisseurCount)

	var areaCount int64 = 0
	db.Table("livraisons").
		Where("code_entreprise = ?", codeEntreprise).
		Count(&areaCount)

	response := map[string]interface{}{
		"client":      clientCount,
		"fournisseur": fournisseurCount,
		"area":        areaCount,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total ventes",
		"data":    response,
	})
}

// Get Zone de livraison
func GetCourbeZoneLivraison(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	var areaCounts []models.AreaCount
	db.Table("livraisons").
		Select("areas.name as area_name, COUNT(*) as count").
		Joins("JOIN areas ON livraisons.area_id = areas.id").
		Where("livraisons.code_entreprise = ?", codeEntreprise).
		Group("areas.name").
		Scan(&areaCounts)

	response := map[string]interface{}{
		"piechart": areaCounts,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total zone de livraisons",
		"data":    response,
	})
}

// Get Clients with most deliveries
func GetClientsWithMostDeliveries(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	var livraisonAreas []models.LivraisonArea
	query := `
		SELECT
			c.fullname,
			c.telephone,
			c.email, 
			COUNT(a.id)
		FROM
			livraisons l
		JOIN
			clients c ON l.client_id = c.id
		JOIN
			areas a ON l.area_id = a.id
		WHERE
		l.code_entreprise = ?
		GROUP BY
			c.fullname,
			c.telephone,
			c.email
		ORDER BY COUNT(a.id) DESC
		LIMIT 10;
	`
	if err := db.Raw(query, codeEntreprise).Scan(&livraisonAreas).Error; err != nil {
		return err
	}
	


	// db.Table("livraisons").
	// 	Select("client_id, COUNT(*) as count").
	// 	Where("code_entreprise = ?", codeEntreprise).
	// 	Where("created_at BETWEEN ? AND ?", start_date, end_date).
	// 	Group("client_id").
	// 	Order("count DESC").
	// 	Limit(10).
	// 	Scan(&clientDeliveryCounts)

	// var clients []models.Client
	// for _, clientDeliveryCount := range clientDeliveryCounts {
	// 	var client models.Client
	// 	db.First(&client, clientDeliveryCount.ClientID)
	// 	clients = append(clients, client)
	// }

	// response := map[string]interface{}{
	// 	"clients": clients,
	// 	"counts":  clientDeliveryCounts,
	// }

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Clients with most deliveries",
		"data":    livraisonAreas,
	})
}

// Get top 10 Fournisseurs with most stock value
func GetTop10FournisseursWithMostStockValue(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	
	var fournisseurStocks []models.FournisseurStock
	query := `
		SELECT
			f.name,
			f.telephone,
			f.type_fourniture, 
			SUM(prix_achat) as total_value
		FROM
			stocks s
		JOIN
			fournisseurs f ON s.fournisseur_id = f.id
	 WHERE
		s.code_entreprise = ?
		GROUP BY
			f.name,
			f.telephone,
			f.type_fourniture
ORDER BY SUM(prix_achat) DESC
LIMIT 10;
	`
	if err := db.Raw(query, codeEntreprise).Scan(&fournisseurStocks).Error; err != nil {
		return err
	}

	

	// var fournisseurStockValues []struct {
	// 	FournisseurID uint
	// 	TotalValue    float64
	// }

	// db.Table("stocks").
	// 	Select("fournisseur_id, SUM(prix_achat) as total_value").
	// 	Where("code_entreprise = ?", codeEntreprise).
	// 	Where("created_at BETWEEN ? AND ?", start_date, end_date).
	// 	Group("fournisseur_id").
	// 	Order("total_value DESC").
	// 	Limit(10).
	// 	Scan(&fournisseurStockValues)

	// var fournisseurs []models.Fournisseur
	// for _, fournisseurStockValue := range fournisseurStockValues {
	// 	var fournisseur models.Fournisseur
	// 	db.First(&fournisseur, fournisseurStockValue.FournisseurID)
	// 	fournisseurs = append(fournisseurs, fournisseur)
	// }

	// response := map[string]interface{}{
	// 	"fournisseurs": fournisseurs,
	// 	"values":       fournisseurStockValues,
	// }

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Top 10 Fournisseurs with most stock value",
		"data":    fournisseurStocks,
	})
}
