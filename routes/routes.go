package routes

import (
	"kgermando/i-pos-restaurant-api/controllers/area"
	"kgermando/i-pos-restaurant-api/controllers/auth"
	"kgermando/i-pos-restaurant-api/controllers/commande"
	"kgermando/i-pos-restaurant-api/controllers/contact"
	"kgermando/i-pos-restaurant-api/controllers/dashboard"
	"kgermando/i-pos-restaurant-api/controllers/entreprise"
	"kgermando/i-pos-restaurant-api/controllers/finance"
	"kgermando/i-pos-restaurant-api/controllers/fournisseurclient"
	"kgermando/i-pos-restaurant-api/controllers/ingredient"
	"kgermando/i-pos-restaurant-api/controllers/livraison"
	"kgermando/i-pos-restaurant-api/controllers/pos"
	"kgermando/i-pos-restaurant-api/controllers/productplat"
	"kgermando/i-pos-restaurant-api/controllers/stock"
	"kgermando/i-pos-restaurant-api/controllers/tablebox"
	"kgermando/i-pos-restaurant-api/controllers/users"
	"kgermando/i-pos-restaurant-api/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {

	api := app.Group("/api", logger.New())

	// Authentification controller
	au := api.Group("/auth")
	au.Post("/register", auth.Register)
	au.Post("/login", auth.Login)
	au.Post("/forgot-password", auth.Forgot)
	au.Post("/reset/:token", auth.ResetPassword)

	au.Post("/entreprise", entreprise.CreateEntreprise)

	app.Use(middlewares.IsAuthenticated)

	au.Get("/user", auth.AuthUser)
	au.Put("/profil/info", auth.UpdateInfo)
	au.Put("/change-password", auth.ChangePassword)
	au.Post("/logout", auth.Logout)

	// User controller
	u := api.Group("/users")
	u.Get("/all", users.GetAllUsers)
	u.Get("/all/paginate", users.GetPaginatedUsers)
	u.Get("/all/paginate/:entreprise_id", users.GetPaginatedUserByID)
	u.Get("/all/:id", users.GetUserByID)
	u.Get("/get/:id", users.GetUser)
	u.Post("/create", users.CreateUser)
	u.Put("/update/:id", users.UpdateUser)
	u.Delete("/delete/:id", users.DeleteUser)

	// Entreprise controller
	e := api.Group("/entreprises")
	e.Get("/all", entreprise.GetAllEntreprises)
	e.Get("/all/paginate", entreprise.GetPaginatedEntreprise)
	e.Get("/get/:id", entreprise.GetEntreprise)
	e.Post("/create", entreprise.CreateEntreprise)
	e.Put("/update/:id", entreprise.UpdateEntreprise)
	e.Delete("/delete/:id", entreprise.DeleteEntreprise)

	// POS controller
	p := api.Group("/pos")
	p.Get("/all", pos.GetAllPoss)
	p.Get("/all/:entreprise_id", pos.GetAllPosById)
	p.Get("/all/paginate", pos.GetPaginatedPos)
	p.Get("/all/paginate/:entreprise_id", pos.GetPaginatedPosByID)
	p.Get("/get/:id", pos.GetPos)
	p.Post("/create", pos.CreatePos)
	p.Put("/update/:id", pos.UpdatePos)
	p.Delete("/delete/:id", pos.DeletePos)

	// Table Box controller
	tb := api.Group("/table-box")
	tb.Get("/:code_entreprise/all/paginate", tablebox.GetPaginatedTableBoxEntreprise)
	tb.Get("/:code_entreprise/:pos_id/all/paginate", tablebox.GetPaginatedTableBox)
	tb.Get("/:code_entreprise/:pos_id/all", tablebox.GetAllTableBox)
	tb.Get("/get/:id", tablebox.GetTableBox)
	tb.Post("/create", tablebox.CreateTableBox)
	tb.Put("/update/:id", tablebox.UpdateTableBox)
	tb.Delete("/delete/:id", tablebox.DeleteTableBox)

	// Product controller
	pr := api.Group("/products")
	pr.Get("/:code_entreprise/all/paginate", productplat.GetPaginatedProductEntreprise)
	pr.Get("/:code_entreprise/:pos_id/all", productplat.GetAllProducts)
	pr.Get("/:code_entreprise/:pos_id/all/paginate", productplat.GetPaginatedProduct)
	pr.Get("/:code_entreprise/:pos_id/all/search", productplat.GetAllProductBySearch)
	pr.Get("/get/:id", productplat.GetProduct)
	pr.Post("/create", productplat.CreateProduct)
	pr.Put("/update/:id", productplat.UpdateProduct)
	pr.Delete("/delete/:id", productplat.DeleteProduct)

	// Stock controller
	s := api.Group("/stocks")
	s.Get("/all", stock.GetAllStocks)
	s.Get("/all/paginate/:product_id", stock.GetPaginatedStock)
	s.Get("/all/total/:product_id", stock.GetTotalStock)
	s.Get("/all/get/:product_id", stock.GetStockMargeBeneficiaire)
	s.Get("/get/:id", stock.GetStock)
	s.Post("/create", stock.CreateStock)
	s.Put("/update/:id", stock.UpdateStock)
	s.Delete("/delete/:id", stock.DeleteStock)

	// Plat controller
	pl := api.Group("/plats")
	pl.Get("/:code_entreprise/all/paginate", productplat.GetPaginatedPlatEntreprise)
	pl.Get("/:code_entreprise/:pos_id/all", productplat.GetAllPlats)
	pl.Get("/:code_entreprise/:pos_id/all/paginate", productplat.GetPaginatedPlat)
	pl.Get("/:code_entreprise/:pos_id/all/search", productplat.GetAllPlatBySearch)
	pl.Get("/get/:id", productplat.GetPlat)
	pl.Post("/create", productplat.CreatePlat)
	pl.Put("/update/:id", productplat.UpdatePlat)
	pl.Delete("/delete/:id", productplat.DeletePlat)

	// Ingredient controller
	in := api.Group("/ingredients")
	in.Get("/:code_entreprise/all/paginate", ingredient.GetPaginatedIngredientEntreprise)
	in.Get("/:code_entreprise/:pos_id/all", ingredient.GetAllIngredients)
	in.Get("/:code_entreprise/:pos_id/all/paginate", ingredient.GetPaginatedIngredient)
	in.Get("/:code_entreprise/:pos_id/all/search", ingredient.GetAllIngredientBySearch)
	in.Get("/get/:id", ingredient.GetIngredient)
	in.Post("/create", ingredient.CreateIngredient)
	in.Put("/update/:id", ingredient.UpdateIngredient)
	in.Delete("/delete/:id", ingredient.DeleteIngredient)

	// ingredients Stock controller
	is := api.Group("/ingredients-stocks")
	is.Get("/all", ingredient.GetAllIngredientStocks)
	is.Get("/all/paginate/:ingredient_id", ingredient.GetPaginatedIngredientStock)
	is.Get("/:code_entreprise/total/get-all/:ingredient_id", ingredient.GetStatsIngredientStock)
	is.Get("/:code_entreprise/total/get/:ingredient_id", ingredient.GetStatsParIngredientStock)
	is.Get("/get/:id", ingredient.GetIngredientStock)
	is.Post("/create", ingredient.CreateIngredientStock)
	is.Put("/update/:id", ingredient.UpdateIngredientStock)
	is.Delete("/delete/:id", ingredient.DeleteIngredientStock)

	// Commande controller
	cmd := api.Group("/commandes")
	cmd.Get("/:code_entreprise/all/paginate", commande.GetPaginatedCommandeEntreprise)
	cmd.Get("/:code_entreprise/:pos_id/all", commande.GetAllCommandes)
	cmd.Get("/:code_entreprise/:pos_id/:table_box_id/all/paginate", commande.GetPaginatedCommandeByTableBox)
	cmd.Get("/:code_entreprise/:table_box_id/total", commande.GetTotalCommande)
	cmd.Get("/get/:id", commande.GetCommande)
	cmd.Post("/create", commande.CreateCommande)
	cmd.Put("/update/:id", commande.UpdateCommande)
	cmd.Delete("/delete/:id", commande.DeleteCommande)

	// Commande line controller
	cmdl := api.Group("/commandes-lines")
	cmdl.Get("/all", commande.GetAllCommandeLines)
	cmdl.Get("/all/:commande_id", commande.GetAllCommandeLineById)
	cmdl.Get("/all/livraison/:livraison_id", commande.GetAllCommandeLineByIdLivraison)
	cmdl.Get("/all/paginate/:commande_id", commande.GetPaginatedCommandeLineByID)
	cmdl.Get("/all/total/:product_id", commande.GetTotalCommandeLine)
	cmdl.Get("/get/:id", commande.GetCommandeLine)
	cmdl.Post("/create", commande.CreateCommandeLine)
	cmdl.Put("/update/:id", commande.UpdateCommandeLine)
	cmdl.Delete("/delete/:id", commande.DeleteCommandeLine)

	// Compositions controller
	comp := api.Group("/compositions")
	comp.Get("/:code_entreprise/:pos_id/all", commande.GetAllCompositions)
	comp.Get("/all/paginate/:plat_id", commande.GetPaginatedComposition)
	comp.Get("/all/total/:plat_id", commande.GetTotalComposition)
	comp.Get("/all/get/:plat_id", commande.GetCompositionMargeBeneficiaire)
	comp.Get("/get/:id", commande.GetComposition)
	comp.Post("/create", commande.CreateComposition)
	comp.Put("/update/:id", commande.UpdateComposition)
	comp.Delete("/delete/:id", commande.DeleteComposition)

	// Client controller
	cl := api.Group("/clients")
	cl.Get("/:code_entreprise/all", fournisseurclient.GetAllClients)
	cl.Get("/:code_entreprise/all/paginate", fournisseurclient.GetPaginatedClient)
	cl.Get("/get/:id", fournisseurclient.GetClient)
	cl.Post("/create", fournisseurclient.CreateClient)
	cl.Post("/uploads", fournisseurclient.UploadCsvDataClient)
	cl.Put("/update/:id", fournisseurclient.UpdateClient)
	cl.Delete("/delete/:id", fournisseurclient.DeleteClient)

	// Fournisseur controller
	fs := api.Group("/fournisseurs")
	fs.Get("/:code_entreprise/all", fournisseurclient.GetAllFournisseurs)
	fs.Get("/:code_entreprise/all/paginate", fournisseurclient.GetPaginatedFournisseur)
	fs.Get("/get/:id", fournisseurclient.GetFournisseur)
	fs.Post("/create", fournisseurclient.CreateFournisseur)
	fs.Put("/update/:id", fournisseurclient.UpdateFournisseur)
	fs.Delete("/delete/:id", fournisseurclient.DeleteFournisseur)

	// Livreur controller
	lv := api.Group("/livreurs")
	lv.Get("/:code_entreprise/all", livraison.GetAllLivreurs)
	lv.Get("/:code_entreprise/all/paginate", livraison.GetPaginatedLivreur)
	lv.Get("/get/:id", livraison.GetLivreur)
	lv.Post("/create", livraison.CreateLivreur)
	lv.Put("/update/:id", livraison.UpdateLivreur)
	lv.Delete("/delete/:id", livraison.DeleteLivreur)

	// Livraison controller
	lvs := api.Group("/livraisons")
	lvs.Get("/:code_entreprise/all/paginate", livraison.GetPaginatedLivraisonEntreprise)
	lvs.Get("/:code_entreprise/:pos_id/all", livraison.GetAllLivraisons)
	lvs.Get("/:code_entreprise/:pos_id/all/paginate", livraison.GetPaginatedLivraison)
	lvs.Get("/:code_entreprise/:pos_id/all/search", livraison.GetAllLivraisonBySearch)
	lvs.Get("/get/:id", livraison.GetLivraison)
	lvs.Post("/create", livraison.CreateLivraison)
	lvs.Put("/update/:id", livraison.UpdateLivraison)
	lvs.Delete("/delete/:id", livraison.DeleteLivraison)

	// Area controller
	are := api.Group("/areas")
	are.Get("/:code_entreprise/all", area.GetAllAreas)
	are.Get("/:code_entreprise/all/paginate", area.GetPaginatedArea)
	are.Get("/get/:id", area.GetArea)
	are.Post("/create", area.CreateArea)
	are.Put("/update/:id", area.UpdateArea)
	are.Delete("/delete/:id", area.DeleteArea)

	// Contact controller
	ctc := api.Group("/contacts")
	ctc.Get("/:code_entreprise/all", contact.GetAllContacts)
	ctc.Get("/:code_entreprise/all/paginate", contact.GetPaginatedContact)
	ctc.Get("/get/:id", contact.GetContact)
	ctc.Post("/create", contact.CreateContact)
	ctc.Put("/update/:id", contact.UpdateContact)
	ctc.Delete("/delete/:id", contact.DeleteContact)

	// Finance controller
	cais := api.Group("/caisses") 
	cais.Get("/:code_entreprise/all/total", finance.GetTotalAllCaisses)
	cais.Get("/:code_entreprise/all", finance.GetAllCaisses)
	cais.Get("/:code_entreprise/:pos_id/all", finance.GetAllCaisseByPos)
	cais.Get("/get/:id", finance.GetCaisse)
	cais.Post("/create", finance.CreateCaisse)
	cais.Put("/update/:id", finance.UpdateCaisse)
	cais.Delete("/delete/:id", finance.DeleteCaisse)

	// Caisse item Controller
	caisseItem := api.Group("/caisse-items") 
	caisseItem.Get("/:code_entreprise/:caisse_id/all/paginate", finance.GetPaginatedCaisseItems)
	caisseItem.Get("/:code_entreprise/:caisse_id/all", finance.GetAllCaisseItems)
	caisseItem.Get("/get/:id", finance.GetCaisseItem)
	caisseItem.Post("/create", finance.CreateCaisseItem)
	caisseItem.Put("/update/:id", finance.UpdateCaisseItem)
	caisseItem.Delete("/delete/:id", finance.DeleteCaisseItem)

	// Dashboard controller
	dash := api.Group("/dashboard")
	dash.Get("/:code_entreprise/all/stocks", dashboard.GetPaginatedStock)
	dash.Get("/:code_entreprise/all/commandeline", dashboard.GetPaginatedCommandeLine)
	dash.Get("/:code_entreprise/all/entree-sortie", dashboard.GetEntreeSortie)
	dash.Get("/:code_entreprise/all/sales-profits", dashboard.GetSaleProfit)
	
	dash.Get("/:code_entreprise/all/total-product-in-stock", dashboard.GetTotalProductInStock)
	dash.Get("/:code_entreprise/all/total-stock-dispo-sortie", dashboard.GetTotalStockDispoSortie)
	dash.Get("/:code_entreprise/all/total-valeur-products", dashboard.GetTotalValeurProduct)
	dash.Get("/:code_entreprise/all/courbe-ventes-jour", dashboard.GetCourbeVente24h)
	dash.Get("/:code_entreprise/all/total-ventes-jour", dashboard.GetTotalVente24h)



	// Dash Plat Product controller
	dash.Get("/:code_entreprise/plats-products/ventes", dashboard.GetTotalPlatProductVendu)
	dash.Get("/:code_entreprise/plats-products/courbe-ventes-profits", dashboard.GetVenteProfitPlatProductMonth)
	dash.Get("/:code_entreprise/plats-products/table/tableau-sortie-products-plats", dashboard.GetTablePaginatedCmdLineSortieProductPlat)
	dash.Get("/:code_entreprise/plats-products/livraison/tableau-sortie-products-plats", dashboard.GetLivraisonPaginatedCmdLineSortieProductPlat)
	dash.Get("/:code_entreprise/all/stocks-disponible", dashboard.GetStockDisponible)
	dash.Get("/:code_entreprise/plats-products/percentage", dashboard.GetCommandeLineLivraisonPercentage)
	dash.Get("/:code_entreprise/plats-products/livraison-count", dashboard.GetCommandeLineLivraisonPieChartData)

	// Dash Caisses controller
	dash.Get("/:code_entreprise/caisses/total", dashboard.GetTotalCaisse)
	dash.Get("/:code_entreprise/caisses/total-ventes-journalieres", dashboard.GetTotalVentesParJour)
	dash.Get("/:code_entreprise/caisses/courbe-ventes-profits", dashboard.GetCourbeVenteProfit24h)
	dash.Get("/:code_entreprise/caisses/tableau-entrees-sorties", dashboard.GetTableauEntreeSorties)
	dash.Get("/:code_entreprise/caisses/total-par-caisse", dashboard.GetTotalParCaisse)

	// Dash Client Fournisseur et Livraison
	dash.Get("/:code_entreprise/cl-fseur-liv/total", dashboard.GetTotalClientFournisseur)
	dash.Get("/:code_entreprise/cl-fseur-liv/courbe-areas", dashboard.GetCourbeZoneLivraison)
	dash.Get("/:code_entreprise/cl-fseur-liv/best-clients", dashboard.GetClientsWithMostDeliveries)
	dash.Get("/:code_entreprise/cl-fseur-liv/best-fournisseurs", dashboard.GetTop10FournisseursWithMostStockValue)
}
