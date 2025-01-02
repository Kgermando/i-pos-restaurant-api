package finance

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedCaisseEntreprise(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var dataList []models.Caisse

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("libelle ILIKE ? OR type_transaction ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("caisses.updated_at DESC").
		Preload("Pos").
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

// Paginate
func GetPaginatedCaisse(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var dataList []models.Caisse

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("libelle ILIKE ? OR type_transaction ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("caisses.updated_at DESC").
		Preload("Pos").
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

// Get All data
func GetAllCaisses(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	var data []models.Caisse
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    data,
	})
}

// Get All data by id
func GetAllCaisseBySearch(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	search := c.Query("search", "")

	var data []models.Caisse
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("libelle ILIKE ? OR type_transaction ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    data,
	})
}

// Get one data
func GetCaisse(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var caisse models.Caisse
	db.Find(&caisse, id)
	if caisse.TypeTransaction == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No caisse name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse found",
			"data":    caisse,
		},
	)
}

// Create data
func CreateCaisse(c *fiber.Ctx) error {
	p := &models.Caisse{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateCaisse(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		TypeTransaction string  `json:"type_transaction"` // Entreé ou Sortie
		Montant         float64 `json:"montant"`          // Montant de la transaction
		Libelle         string  `json:"libelle"`          // Description de la transaction
		Reference       string  `json:"reference"`        // Nombre aleatoire
		Signature       string  `json:"signature"`        // Signature de la transaction
		PosID           uint    `json:"pos_id"`           // ID du point de vente
		CodeEntreprise  uint    `json:"code_entreprise"`
	}

	var updateData UpdateData

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your iunput",
				"data":    nil,
			},
		)
	}

	caisse := new(models.Caisse)

	db.First(&caisse, id)
	caisse.TypeTransaction = updateData.TypeTransaction
	caisse.Montant = updateData.Montant
	caisse.Libelle = updateData.Libelle
	caisse.Reference = updateData.Reference
	caisse.Signature = updateData.Signature
	caisse.PosID = updateData.PosID
	caisse.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse updated success",
			"data":    caisse,
		},
	)

}

// Delete data
func DeleteCaisse(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var caisse models.Caisse
	db.First(&caisse, id)
	if caisse.TypeTransaction == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No caisse name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse deleted success",
			"data":    nil,
		},
	)
}
