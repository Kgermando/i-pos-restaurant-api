package livraison

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/database"
	"kgermando/i-pos-restaurant-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedLivreur(c *fiber.Ctx) error {
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

	var dataList []models.Livreur

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name_society ILIKE ? OR livreur_name ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("livreurs.updated_at DESC"). 
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
		"message":    "All livreurs",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllLivreurs(c *fiber.Ctx) error {
	codeEntreprise := c.Params("code_entreprise")
	db := database.DB

	var data []models.Livreur
	db.Where("code_entreprise = ?", codeEntreprise).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All livreurs",
		"data":    data,
	})
}

// Get one data
func GetLivreur(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var livreur models.Livreur
	db.Find(&livreur, id)
	if livreur.NameSociety == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livreur found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur found",
			"data":    livreur,
		},
	)
}

// Create data
func CreateLivreur(c *fiber.Ctx) error {
	p := &models.Livreur{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateLivreur(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		NameSociety    string `json:"name_society"`
		LivreurName    string `json:"livreur_name"`
		Telephone      string `json:"telephone"`
		Email          string `json:"email"`
		Rccm           string `json:"rccm"`
		IdNat          string `json:"idnat"`
		Signature      string `json:"signature"`
		CodeEntreprise uint   `json:"code_entreprise"` 
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

	livreur := new(models.Livreur)

	db.First(&livreur, id)
	livreur.NameSociety = updateData.NameSociety
	livreur.LivreurName = updateData.LivreurName
	livreur.Telephone = updateData.Telephone
	livreur.Email = updateData.Email
	livreur.Rccm = updateData.Rccm
	livreur.IdNat = updateData.IdNat 
	livreur.Signature = updateData.Signature
	livreur.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&livreur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur updated success",
			"data":    livreur,
		},
	)

}

// Delete data
func DeleteLivreur(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var livreur models.Livreur
	db.First(&livreur, id)
	if livreur.NameSociety == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livreur found",
				"data":    nil,
			},
		)
	}

	db.Delete(&livreur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur deleted success",
			"data":    nil,
		},
	)
}
