package dashboard

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Total caisses entrees
func GetTotalCaisse(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise") 
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Total des caisses entrees
	var caisseEntrees []models.CaisseItem
	db.Where("caisse_items.code_entreprise = ?", codeEntreprise).
	Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
	Where("type_transaction = ?", "Entrée").
	Find(&caisseEntrees)

	// Calcul du total des caisses sorties
	var caisseSorties []models.CaisseItem
	db.Where("caisse_items.code_entreprise = ?", codeEntreprise).
	Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
	Where("type_transaction = ?", "Sortie").
	Find(&caisseSorties)

	// Calcul du total des caisses entrees
	totalCaisseEntree := 0.0
	for _, caisse := range caisseEntrees {
		totalCaisseEntree += caisse.Montant
	}

	// Calcul du total des caisses sorties
	totalCaisseSorties := 0.0
	for _, caisse := range caisseSorties {
		totalCaisseSorties += caisse.Montant
	}

	// Calcul du total des caisses entrees et sorties
	totalCaisseEntreeSorties := 0.0
	totalCaisseEntreeSorties = totalCaisseEntree + totalCaisseSorties

	response := map[string]interface{} {
		"total_caisse_entree": totalCaisseEntree,
		"total_caisse_sorties": totalCaisseSorties,
		"total_caisse_entree_sorties": totalCaisseEntreeSorties,
	}
 
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "total Caisses",
		"data":    response,
	})
}


// Total des ventes journaileres
func GetTotalVentesParJour(c *fiber.Ctx) error {  
	db := database.DB
	codeEntreprise := c.Params("code_entreprise") 
 
	var commandeLines []models.CommandeLine
	var commandeLineLivraisons []models.CommandeLine
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Second)
	db.Joins("JOIN commandes ON commande_lines.commande_id = commandes.id").
	Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("commande_lines.created_at BETWEEN ? AND ?", startOfDay, endOfDay).
	Where("commandes.status != ? OR commandes.status != ?", "En cours", "Créance").
	Preload("Plat").
	Preload("Product").
	Find(&commandeLines)

	db.Joins("JOIN livraisons ON commande_lines.livraison_id = livraisons.id").
	Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("commande_lines.created_at BETWEEN ? AND ?", startOfDay, endOfDay). 
	Where("livraisons.status != ? OR livraisons.status != ?", "En cours", "Créance").
	Preload("Plat").
	Preload("Product").
	Find(&commandeLineLivraisons)

	commandeLines = append(commandeLines, commandeLineLivraisons...)


	ingredientUsages, err := GetTotalIngredientUsage(db, codeEntreprise)
	if err != nil {
		return nil
	}

	totalSales := 0.0
	totalProfits := 0.0

	profitAmountProduit := 0.0
	profitAmountPlat := 0.0

	for _, line := range commandeLines {
		totalSales += (line.Plat.PrixVente + line.Product.PrixVente)

		stockEntry := getStockEntry(line.ProductID, db)
		if line.ProductID > 0 {
			profitAmountProduit = (line.Product.PrixVente - stockEntry.PrixAchat) * float64(line.Quantity)
		}
		if line.PlatID > 0 {
			priceIngredientUsagePlat := 0.0 // Prix total des ingrédients utilisés pour le plat
			for _, ingredientUsage := range ingredientUsages {
				if ingredientUsage.PlatID == line.PlatID {
					priceIngredientUsagePlat += ingredientUsage.Price
				}
			}
			
			profitAmountPlat = (line.Plat.PrixVente - priceIngredientUsagePlat) * float64(line.Quantity)
		}


		totalProfits += profitAmountProduit + profitAmountPlat
	}

	response := map[string]interface{} {
		"total_ventes": totalSales,
		"total_profits": totalProfits,
	}
 
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "total ventes journalieres",
		"data":    response,
	})
}


func GetCourbeVenteProfit24h(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")


	var commandeLines []models.CommandeLine
	var commandeLineLivraisons []models.CommandeLine
	db.Joins("JOIN commandes ON commande_lines.commande_id = commandes.id").
	Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("commandes.status != ?","En cours").
	Where("date(commande_lines.created_at) = date('now')").
	Preload("Plat").
	Preload("Product").
	Find(&commandeLines)
	db.Joins("JOIN livraisons ON commande_lines.livraison_id = livraisons.id").
	Where("commande_lines.code_entreprise = ?", codeEntreprise).
	Where("livraisons.status != ?", "En cours").
	Where("date(commande_lines.created_at) = date('now')").
	Preload("Plat").
	Preload("Product").
	Find(&commandeLineLivraisons)

	commandeLines = append(commandeLines, commandeLineLivraisons...)

	ingredientUsages, err := GetTotalIngredientUsage(db, codeEntreprise)
	if err != nil {
		return nil
	}

	hourlySales := make(map[int]float64)
	hourlyProfits := make(map[int]float64)
	for i := 0; i < 24; i++ {
		hourlySales[i] = 0
		hourlyProfits[i] = 0
	}

	profitAmountProduit := 0.0
	profitAmountPlat := 0.0

	for _, line := range commandeLines {
		hour := line.CreatedAt.Hour() 
		hourlySales[hour] += (line.Plat.PrixVente + line.Product.PrixVente)

		stockEntry := getStockEntry(line.ProductID, db)
		if line.ProductID > 0 {
			profitAmountProduit = (line.Product.PrixVente - stockEntry.PrixAchat) * float64(line.Quantity)
		}
		if line.PlatID > 0 {
			priceIngredientUsagePlat := 0.0 // Prix total des ingrédients utilisés pour le plat
			for _, ingredientUsage := range ingredientUsages {
				if ingredientUsage.PlatID == line.PlatID {
					priceIngredientUsagePlat += ingredientUsage.Price
				}
			}
			
			profitAmountPlat = (line.Plat.PrixVente - priceIngredientUsagePlat) * float64(line.Quantity)
		}


		hourlyProfits[hour] += profitAmountProduit + profitAmountPlat

	}

	response := map[string]interface{} {
		"hourly_sales": hourlySales,
		"hourly_profits": hourlyProfits,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Courbe de ventes journaliere",
		"data":    response,
	})
}


func GetTableauEntreeSorties(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise") 
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit 

	var dataList []models.CaisseItem

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Offset(offset).
		Limit(limit).
		Order("caisse_items.updated_at DESC").
		Preload("Caisse.Pos").
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
		"message":    "All caisses",
		"data":       dataList,
		"pagination": pagination,
	})
}

func GetTotalParCaisse(c *fiber.Ctx) error { 
	db := database.DB
	codeEntreprise := c.Params("code_entreprise") 
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var dashCaisse models.CaisseDashboard
 
	query := `
		SELECT
    		c.id,
			c.name,
			SUM(CASE WHEN ci.type_transaction = 'Entrée' THEN ci.montant ELSE 0 END) AS total_entrees,
			SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END) AS total_sorties,
			SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END) AS solde
		FROM
			caisses c
		JOIN
			caisse_items ci ON c.id = ci.caisse_id
		WHERE ci.code_entreprise = ? AND ci.created_at BETWEEN ? AND ?
		GROUP BY
			c.id, c.name;
	`

	if err := db.Raw(query, codeEntreprise, startDateStr, endDateStr).
		Scan(&dashCaisse).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Total ParCaisse",
		"data":       "dataList", 
	})
}
