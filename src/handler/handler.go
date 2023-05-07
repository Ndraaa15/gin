package handler

import (
	"gin/src/middleware"
	"gin/src/usecase"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var once = sync.Once{}

type Handler struct {
	http *gin.Engine
	uc   *usecase.Usecase
}

func InitHandler(uc *usecase.Usecase) *Handler {
	r := Handler{}
	once.Do(func() {
		r.http = gin.Default()
		r.uc = uc
		r.RoutesAndMiddleware()

	})
	return &r
}

func (h *Handler) RoutesAndMiddleware() {
	/*
		Endpoint testing
	*/

	h.http.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	v1 := h.http.Group("/api/v1")

	/*
		CORS
	*/
	h.http.Use(middleware.CORS())

	/*
		Endpoint for not user
	*/

	h.http.POST(v1.BasePath()+"/tourism/create", h.PostTourism)  //post ecotourism by admin
	h.http.POST(v1.BasePath()+"/trash/validate", h.ValidateCode) //validate exchange trash by admin

	/*
		Endpoint for user
	*/
	user := h.http.Group(v1.BasePath() + "/user")
	user.POST("/signup", h.UserRegister) //new user signup
	user.POST("/login", h.UserLogin)     //user login
	user.Use(middleware.IsUserLoggedIn)
	user.GET("/profile", h.GetProfile)
	user.DELETE("/delete", h.DeleteAccount)          //user delete account
	user.PUT("/update", h.UserUpdate)                //user update profile without photo profile
	user.POST("/upload/photo", h.UploadPhotoProfile) //user upload photo
	user.GET("/cart", h.GetCart)                     //Get user cart
	user.GET("/history")                             //Get user booking history

	h.http.Use(middleware.CORS())

	eco := h.http.Group(v1.BasePath() + "/tourism")
	eco.Use(middleware.IsUserLoggedIn).
		GET("/", h.GetAllTourism).                                       //get all eco tourism
		GET("/:id", h.GetTourismByID).                                   //get eco tourism by id
		GET("/filter/category/:category", h.GetTourismByCategory).       //get filtered eco tourism by category
		GET("/filter/price/:startPrice/:endPrice", h.GetTourismByPrice). //get filtered eco tourism by price
		GET("/filter/region/:region", h.GetTourismByRegion).             //get filtered eco tourism by region
		POST("/:id/comment", h.UploadComment).                           //comment on current eco tourism, id in here is id ecotourism
		POST("/add/:id/cart", h.AddCart).                                //Add to cart
		DELETE("/del/:id/cart", h.DeleteCartContent)                     //Delete content cart

	trash := h.http.Group(v1.BasePath() + "/trash")
	trash.Use(middleware.IsUserLoggedIn).
		POST("/exchange", h.ExchangeTrash).              //exchange trash into coin
		GET("/exchange/history", h.ExchangeTrashHistory) //get exchange trash history

	article := h.http.Group(v1.BasePath() + "/article")
	article.Use(middleware.IsUserLoggedIn).
		GET("/", h.GetAllArticles).      //get all articles
		GET(":id", h.GetArticleByID).    //get article by id
		POST("/create", h.CreateArticle) //create a article

	payment := h.http.Group(v1.BasePath() + "/payment")
	payment.Use(middleware.IsUserLoggedIn).
		POST("/online"). //Payment gateway using online payment
		POST("/coin")    //Payment gateway using traveleen coin
}

func (h *Handler) Run() {
	h.http.Run()
}
