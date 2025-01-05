package dashboard

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get data
func GetTotalPlatProductVendu(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var dataList []models.CommandeLine
	db.Where("commande_lines.code_entreprise = ?", codeEntreprise).
		Where("commande_lines.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Plat").
		Preload("Product").
		Find(&dataList)

	// Par Table
	var totalPlatVenduTable float64 = 0
	var nombreTotalPlatVenduTable float64 = 0
	var totalProduitVenduTable float64 = 0

	// Par Livraison
	var totalPlatVenduLivraison float64 = 0
	var nombreTotalPlatVenduLivraison float64 = 0
	var totalProduitVenduLivraison float64 = 0

	for _, data := range dataList {
		if data.CommandeID > 0 {
			if data.PlatID > 0 {
				totalPlatVenduTable += data.Plat.PrixVente
				nombreTotalPlatVenduTable += float64(data.Quantity)
			}
			if data.ProductID > 0 {
				totalProduitVenduTable += data.Product.PrixVente
			}
		}
		if data.LivraisonID > 0 {
			if data.PlatID > 0 {
				totalPlatVenduLivraison += data.Plat.PrixVente
				nombreTotalPlatVenduLivraison += float64(data.Quantity)
			}
			if data.ProductID > 0 {
				totalProduitVenduLivraison += data.Product.PrixVente
			}
		}
	}

	response := map[string]interface{}{
		"totalPlatVenduTable":           totalPlatVenduTable,
		"nombreTotalPlatVenduTable":     nombreTotalPlatVenduTable,
		"totalProduitVenduTable":        totalProduitVenduTable,
		"totalPlatVenduLivraison":       totalPlatVenduLivraison,
		"nombreTotalPlatVenduLivraison": nombreTotalPlatVenduLivraison,
		"totalProduitVenduLivraison":    totalProduitVenduLivraison,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total ventes",
		"data":    response,
	})
}

// Courbe de ventes et profits des plats et produits
func GetVenteProfitPlatProductMonth(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var commandes []models.Commande

	db.Where("commandes.code_entreprise = ?", codeEntreprise).
		Where("commandes.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("CommandeLines.Plat").
		Preload("CommandeLines.Product").
		Find(&commandes)
 
	ingredientUsages, err := GetTotalIngredientUsage(db, codeEntreprise)
	if err != nil {
		return nil
	} 

	saleAmount := 0.0 // Total des ventes
	saleAmountPlat := 0.0 // Total des ventes des plats
	saleAmountProduit := 0.0 // Total des ventes des produits

	profitAmount := 0.0 // Total des profits
	profitAmountPlat := 0.0 // Total des profits des plats
	profitAmountProduit := 0.0 // Total des profits des produits

	salesByMonth := make(map[string]float64)
	profitsByMonth := make(map[string]float64)

	for _, commande := range commandes {
		month := commande.CreatedAt.Format("2006-01")
		for _, line := range commande.CommandeLines {
			stockEntry := getStockEntry(line.ProductID, db)
			if line.ProductID > 0 {
				saleAmountProduit = line.Product.PrixVente * float64(line.Quantity)
				profitAmountProduit = (line.Product.PrixVente - stockEntry.PrixAchat) * float64(line.Quantity)
			}
			if line.PlatID > 0 {
				priceIngredientUsagePlat := 0.0 // Prix total des ingrédients utilisés pour le plat
				for _, ingredientUsage := range ingredientUsages {
					if ingredientUsage.PlatID == line.PlatID {
						priceIngredientUsagePlat += ingredientUsage.Price
					}
				}
				saleAmountPlat = line.Plat.PrixVente * float64(line.Quantity)
				profitAmountPlat = (line.Plat.PrixVente - priceIngredientUsagePlat) * float64(line.Quantity)
			}

			saleAmount += saleAmountProduit + saleAmountPlat
			profitAmount += profitAmountProduit + profitAmountPlat

	
			salesByMonth[month] += saleAmount
			profitsByMonth[month] += profitAmount
		}
	}
	

	var response []map[string]interface{}

	for month, sales := range salesByMonth {
		response = append(response, map[string]interface{}{
			"mois":    month,
			"ventes":  sales,
			"profits": profitsByMonth[month],
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Ventes et profits",
		"data":    response,
	})
}

// Tabeleau Tous sorties des produits et plats pour Table only
func GetTablePaginatedCmdLineSortieProductPlat(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date") 
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var dataList []models.CommandeLine

	var length int64
	db.Model(dataList).Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("commande_lines.livraison_id = ?", 0).Count(&length)
	db.Where("commande_lines.code_entreprise = ?", codeEntreprise).
		Where("commande_lines.created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("commande_lines.livraison_id = ?", 0).
		Offset(offset).
		Limit(limit).
		Order("commande_lines.updated_at DESC").
		Preload("Product").
		Preload("Plat").
		Preload("Commande.TableBox").
		Preload("Commande.Client").
		Preload("Commande.Pos").
		Find(&dataList)

	if err != nil {
		fmt.Println("error s'est produite: ", err)
		return c.Status(500).SendString(err.Error())
	}

	// Calculate total number of pages
	totalPages := len(dataList) / limit
	if remainder := len(dataList) % limit; remainder > 0 {
		totalPages++
	}
	pagination := map[string]interface{}{
		"total_pages": totalPages,
		"page":        page,
		"page_size":   limit,
		"length":      length,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All dash Table sortie Product Plat",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Tabeleau Tous sorties des produits et plats pour Livraison only
func GetLivraisonPaginatedCmdLineSortieProductPlat(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date") 
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var dataList []models.CommandeLine

	var length int64
	db.Model(dataList).Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("commande_lines.livraison_id != ?", 0).Count(&length)
	db.Where("commande_lines.code_entreprise = ?", codeEntreprise).
		Where("commande_lines.created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("commande_lines.livraison_id != ?", 0).
		Offset(offset).
		Limit(limit).
		Order("commande_lines.updated_at DESC").
		Preload("Product").
		Preload("Plat").
		Preload("Livraison.Client"). 
		Preload("Livraison.Livreur").
		Preload("Livraison.Area").
		Preload("Livraison.Pos").
		Find(&dataList)

	if err != nil {
		fmt.Println("error s'est produite: ", err)
		return c.Status(500).SendString(err.Error())
	}

	// Calculate total number of pages
	totalPages := len(dataList) / limit
	if remainder := len(dataList) % limit; remainder > 0 {
		totalPages++
	}
	pagination := map[string]interface{}{
		"total_pages": totalPages,
		"page":        page,
		"page_size":   limit,
		"length":      length,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All dash Sortie livraison Product Plat",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get percentage of commande_lines with livraison_id != 0 and livraison_id = 0
func GetCommandeLineLivraisonPercentage(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var totalCommandeLines int64
	var livraisonCommandeLines int64
	var tableCommandeLines int64

	db.Model(&models.CommandeLine{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Count(&totalCommandeLines)

	db.Model(&models.CommandeLine{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("livraison_id != ?", 0).
		Count(&livraisonCommandeLines)

	db.Model(&models.CommandeLine{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("livraison_id = ?", 0).
		Count(&tableCommandeLines)

	livraisonPercentage := (float64(livraisonCommandeLines) / float64(totalCommandeLines)) * 100
	tablePercentage := (float64(tableCommandeLines) / float64(totalCommandeLines)) * 100

	response := map[string]interface{}{
		"livraisonPercentage": livraisonPercentage,
		"tablePercentage":     tablePercentage,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pourcentage des commandes par type",
		"data":    response,
	})
}
// Get data for pie chart
func GetCommandeLineLivraisonPieChartData(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var livraisonCommandeLines int64
	var tableCommandeLines int64

	db.Model(&models.CommandeLine{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("livraison_id != ?", 0).
		Count(&livraisonCommandeLines)

	db.Model(&models.CommandeLine{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Where("livraison_id = ?", 0).
		Count(&tableCommandeLines)

	response := map[string]interface{}{
		"livraison": livraisonCommandeLines,
		"table":     tableCommandeLines,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data for pie chart",
		"data":    response,
	})
}

func GetTotalIngredientUsage(db *gorm.DB, code_entreprise string) ([]models.IngredientUsage, error) {
	var ingredientUsages []models.IngredientUsage
	query := `
		SELECT
			cl.plat_id,
			i.name,
			SUM(c.quantity * cl.quantity) AS qty,
			i.unite
		FROM
			commande_lines cl
		JOIN
			compositions c ON cl.plat_id = c.plat_id
		JOIN
			ingredients i ON c.ingredient_id = i.id
		WHERE
            cl.code_entreprise = ?
		GROUP BY
			cl.plat_id, i.name, i.unite;
    `
	if err := db.Raw(query, code_entreprise).Scan(&ingredientUsages).Error; err != nil {
		return nil, err
	}

	return ingredientUsages, nil
}