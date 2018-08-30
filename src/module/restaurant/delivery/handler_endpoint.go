package delivery

import (
	"log"
	"strconv"
	"time"

	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/restaurant"
	"github.com/atletaid/go-template/util/httputil"
	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	rtu restaurant.Usecase
}

func NewRestaurantHandler(router *gin.Engine, rtu restaurant.Usecase) *gin.Engine {
	handler := &RestaurantHandler{rtu}

	v1 := router.Group("/api")
	v1.POST("/restaurant", handler.CreateRestaurantEndpoint())
	v1.GET("/restaurant/:restaurant_id", handler.GetRestaurantEndpoint())
	v1.GET("/restaurants", handler.GetAllRestaurantsEndpoint())
	v1.POST("/restaurant/city", handler.GetRestaurantsByCityEndpoint())
	v1.DELETE("/restaurant/:restaurant_id", handler.DeleteRestaurantEndpoint())

	return router
}

type createRestaurantRequest struct {
	RestaurantName        string  `json:"restaurant_name" form:"restaurant_name"`
	RestaurantTimeMinute  int     `json:"restaurant_time_minute" form:"restaurant_time_minute"`
	RestaurantPrice       int     `json:"restaurant_price" form:"restaurant_price"`
	PositionLat           float64 `json:"position_lat" form:"position_lat"`
	PositionLong          float64 `json:"position_long" form:"position_long"`
	RestaurantCity        string  `json:"restaurant_city" form:"restaurant_city"`
	RestaurantImage       string  `json:"restaurant_image" form:"restaurant_image"`
	RestaurantDescription string  `json:"restaurant_description" form:"restaurant_description"`
}

type createRestaurantResponse struct {
	RestaurantID int64 `json:"restaurant_id"`
}

func (h *RestaurantHandler) CreateRestaurantEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := createRestaurantRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		restaurantID, err := h.rtu.CreateRestaurant(req.RestaurantName, req.RestaurantCity, req.RestaurantImage, req.RestaurantDescription, req.RestaurantTimeMinute, req.RestaurantPrice, req.PositionLat, req.PositionLong)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := createRestaurantResponse{
			RestaurantID: restaurantID,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success create restaurant"}, processTime, resp)
	}
}

type dataRestaurantResponse struct {
	Restaurant *model.Restaurant `json:"restaurant"`
}

func (h *RestaurantHandler) GetRestaurantEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		restaurantID, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		restaurant, err := h.rtu.GetRestaurant(restaurantID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataRestaurantResponse{
			Restaurant: restaurant,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get restaurant"}, processTime, resp)
	}
}

type dataRestaurantsResponse struct {
	Restaurants model.Restaurants `json:"restaurants"`
}

func (h *RestaurantHandler) GetAllRestaurantsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		restaurants, err := h.rtu.GetAllRestaurants()
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataRestaurantsResponse{
			Restaurants: restaurants,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get all restaurants"}, processTime, resp)
	}
}

type getRestaurantByCityRequest struct {
	City string `json:"restaurant_city" form:"restaurant_city"`
}

func (h *RestaurantHandler) GetRestaurantsByCityEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := getRestaurantByCityRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		restaurants, err := h.rtu.GetRestaurantsByCity(req.City)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get restaurants by city"}, processTime, restaurants)
	}
}

func (h *RestaurantHandler) DeleteRestaurantEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		restaurantID, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		err = h.rtu.DeleteRestaurantByID(restaurantID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success delete restaurant"}, processTime, nil)
	}
}
