package productplat

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedPlatEntreprise(c *fiber.Ctx) error {
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

	var dataList []models.Plat

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("plats.updated_at DESC").
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
		"message":    "All Plats",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedPlat(c *fiber.Ctx) error {
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

	var dataList []models.Plat

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("plats.updated_at DESC").
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
		"message":    "All Plats",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllPlats(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	var data []models.Plat
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Plats",
		"data":    data,
	})
}

// Get All data by id
func GetAllPlatBySearch(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posId := c.Params("pos_id")

	search := c.Query("search", "")

	var data []models.Plat
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_id = ?", posId).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Plats",
		"data":    data,
	})
}

// Get one data
func GetPlat(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var plat models.Plat
	db.Find(&plat, id)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Plat name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat found",
			"data":    plat,
		},
	)
}

// Create data
func CreatePlat(c *fiber.Ctx) error {
	p := &models.Plat{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Plat created success",
			"data":    p,
		},
	)
}

// Update data
func UpdatePlat(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		Reference      string  `json:"reference"`
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		UniteVente     string  `json:"unite_vente"`
		PrixVente      float64 `json:"prix_vente"`
		Tva            float64 `json:"tva"`
		Signature      string  `json:"signature"`
		PosID          uint    `json:"pos_id"`
		CodeEntreprise uint    `json:"code_entreprise"`
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

	plat := new(models.Plat)

	db.First(&plat, id)
	plat.Reference = updateData.Reference
	plat.Name = updateData.Name
	plat.Description = updateData.Description
	plat.UniteVente = updateData.UniteVente
	plat.PrixVente = updateData.PrixVente
	plat.Tva = updateData.Tva
	plat.Signature = updateData.Signature
	plat.PosID = updateData.PosID
	plat.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&plat)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat updated success",
			"data":    plat,
		},
	)

}

// Delete data
func DeletePlat(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var plat models.Plat
	db.First(&plat, id)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No plat name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&plat)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat deleted success",
			"data":    nil,
		},
	)
}
