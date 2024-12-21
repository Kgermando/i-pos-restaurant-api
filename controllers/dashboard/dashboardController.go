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
func GetStock(c *fiber.Ctx) error {
	db := database.DB
	productId := c.Params("product_id")

	var data models.Stock

	db.Model(data).Where("product_id = ?", productId).Last(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty stocks",
		"data":    data,
	})
}

// Paginate
func GetPaginatedStock(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	fmt.Println("startDateStr", startDateStr)
	fmt.Println("endDateStr", endDateStr)

	// startDate, err := time.Parse("2006-01-02", startDateStr)
	// if err != nil {
	//     return c.Status(400).SendString("Invalid start date format")
	// }

	// endDate, err := time.Parse("2006-01-02", endDateStr)
	// if err != nil {
	//     return c.Status(400).SendString("Invalid end date format")
	// }

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var dataList []models.Stock

	var length int64
	db.Model(dataList).Where("stocks.code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("stocks.code_entreprise = ?", codeEntreprise).
		Where("stocks.created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		// Joins("JOIN products ON stocks.product_id=products.id").
		// Joins("JOIN fournisseurs ON stocks.fournisseur_id=fournisseurs.id").
		// Select(`
		// 	stocks.id AS id,
		// 	products.reference AS reference,
		// 	products.name AS name,
		// 	stocks.quantity AS quantity,
		// 	stocks.prix_achat AS prix_achat,
		// 	stocks.date_expiration AS date_expiration,
		// 	stocks.description AS description,
		// 	stocks.created_at AS created_at,
		// 	stocks.updated_at AS updated_at,
		// 	stocks.signature AS signature
		// `).
		Offset(offset).
		Limit(limit).
		Order("stocks.created_at DESC").
		Preload("Product").
		Preload("Fournisseur").
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
		"message":    "All dashboard stocks",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedCommandeLine(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// startDate, err := time.Parse("2006-01-02", startDateStr)
	// if err != nil {
	//     return c.Status(400).SendString("Invalid start date format")
	// }

	// endDate, err := time.Parse("2006-01-02", endDateStr)
	// if err != nil {
	//     return c.Status(400).SendString("Invalid end date format")
	// }

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
	db.Model(dataList).Where("commande_lines.code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("commande_lines.code_entreprise = ?", codeEntreprise).
		Where("commande_lines.created_at BETWEEN ? AND ?", startDateStr, endDateStr).
		Offset(offset).
		Limit(limit).
		Order("commande_lines.updated_at DESC").
		Preload("Product").
		Preload("Commande.Client").
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
		"message":    "All dashboard commandeLines",
		"data":       dataList,
		"pagination": pagination,
	})
}

func GetEntreeSortie(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var products []models.Product
	db.Where("products.code_entreprise = ?", codeEntreprise).
		Where("products.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Stocks").Preload("CommadeLines").Find(&products)

	response := []map[string]interface{}{}
	for _, product := range products {

		var totalStockIn uint64 = 0
		for _, entry := range product.Stocks {
			totalStockIn += entry.Quantity
		}

		var totalStockOut uint64 = 0
		for _, orderLine := range product.CommadeLines {
			totalStockOut += orderLine.Quantity
		}

		response = append(response, map[string]interface{}{
			"product":   product.Name,
			"stock_in":  totalStockIn,
			"stock_out": totalStockOut,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Entrees et sorties",
		"data":    response,
	})
}

func GetSaleProfit(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var commandes []models.Commande

	db.Where("commandes.code_entreprise = ?", codeEntreprise).
		Where("commandes.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("CommandeLines.Product").Find(&commandes)

	salesByMonth := make(map[string]float64)
	profitsByMonth := make(map[string]float64)

	for _, commande := range commandes {
		month := commande.CreatedAt.Format("2006-01")
		for _, line := range commande.CommandeLines {
			stockEntry := getStockEntry(line.ProductID, db)
			saleAmount := line.Product.PrixVente * float64(line.Quantity)
			profitAmount := (line.Product.PrixVente - stockEntry.PrixAchat) * float64(line.Quantity)

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

func GetStockDisponible(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var products []models.Product
	db.Where("products.code_entreprise = ?", codeEntreprise).
		Where("products.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Stocks").Preload("CommadeLines").Find(&products)

	stockTotal := make(map[string]uint64)

	for _, product := range products {
		var stockEntries []models.Stock
		var commandeLines []models.CommandeLine

		db.Where("product_id = ?", product.ID).Find(&stockEntries)
		db.Where("product_id = ?", product.ID).Find(&commandeLines)

		var totalStockIn uint64 = 0
		for _, entry := range stockEntries {
			totalStockIn += entry.Quantity
		}

		var totalStockOut uint64 = 0
		for _, line := range commandeLines {
			totalStockOut -= line.Quantity
		}

		stockTotal[product.Name] = totalStockIn + totalStockOut
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stocks disponible",
		"data":    stockTotal,
	})
}

func GetTotalProductInStock(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var stocks []models.Stock
	db.Where("stocks.code_entreprise = ?", codeEntreprise).
		Where("stocks.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Product").Find(&stocks)

	var totalProducts uint64 = 0
	for _, entry := range stocks {
		totalProducts += entry.Quantity
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stocks disponible",
		"data":    totalProducts,
	})
}

func GetTotalStockDispoSortie(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var commandes []models.Commande
	db.Where("commandes.code_entreprise = ?", codeEntreprise).
		Where("commandes.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("CommandeLines.Product").Find(&commandes)

	var stocks []models.Stock
	db.Where("stocks.code_entreprise = ?", codeEntreprise).
		Where("stocks.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Product").Find(&stocks)

	totalProfitAvailableStock := 0.0
	totalProfitStockInOut := 0.0

	// Calcul des profits des stocks disponibles
	for _, entry := range stocks {
		profit := (entry.Product.PrixVente - entry.PrixAchat) * float64(entry.Quantity)
		totalProfitAvailableStock += profit
	}

	// Calcul des profits des stocks entrés et sortis
	for _, commande := range commandes {
		for _, line := range commande.CommandeLines {
			stockEntry := getStockEntry(line.ProductID, db)
			profit := (line.Product.PrixVente - stockEntry.PrixAchat) * float64(line.Quantity)
			totalProfitStockInOut += profit
		}
	}

	response := map[string]interface{}{
		"total_profit_available_stock": totalProfitAvailableStock,
		"total_profit_stock_in_out":    totalProfitStockInOut,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total profits dispo et sortie",
		"data":    response,
	})
}

func GetTotalValeurProduct(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	start_date := c.Query("start_date")
	end_date := c.Query("end_date")

	var stocks []models.Stock
	db.Where("stocks.code_entreprise = ?", codeEntreprise).
		Where("stocks.created_at BETWEEN ? AND ?", start_date, end_date).
		Preload("Product").Find(&stocks)

	totalValue := 0.0
	for _, entry := range stocks {
		productValue := entry.Product.PrixVente * float64(entry.Quantity)
		totalValue += productValue
	} 

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "total valeur",
		"data":    totalValue,
	})
}

func getStockEntry(productID uint, db *gorm.DB) models.Stock {
	var stock models.Stock
	db.Where("product_id = ?", productID).First(&stock)
	return stock
}
