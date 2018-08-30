package delivery

import (
	"log"
	"strconv"
	"time"

	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/recreation"
	"github.com/atletaid/go-template/util/httputil"
	"github.com/gin-gonic/gin"
)

type RecreationHandler struct {
	ru recreation.Usecase
}

func NewRecreationHandler(router *gin.Engine, ru recreation.Usecase) *gin.Engine {
	handler := &RecreationHandler{ru}

	v1 := router.Group("/api")
	v1.POST("/recreation", handler.CreateRecreationEndpoint())
	v1.GET("/recreation/:recreation_id", handler.GetRecreationEndpoint())
	v1.GET("/recreations", handler.GetAllRecreationsEndpoint())
	v1.POST("/recreation/city", handler.GetRecreationsByCityEndpoint())
	v1.DELETE("/recreation/:recreation_id", handler.DeleteRecreationEndpoint())

	return router
}

type createRecreationRequest struct {
	RecreationName        string  `json:"recreation_name" form:"recreation_name"`
	RecreationTimeMinute  int     `json:"recreation_time_minute" form:"recreation_time_minute"`
	RecreationPrice       int     `json:"recreation_price" form:"recreation_price"`
	PositionLat           float64 `json:"position_lat" form:"position_lat"`
	PositionLong          float64 `json:"position_long" form:"position_long"`
	RecreationCity        string  `json:"recreation_city" form:"recreation_city"`
	RecreationImage       string  `json:"recreation_image" form:"recreation_image"`
	RecreationDescription string  `json:"recreation_description" form:"recreation_description"`
}

type createRecreationResponse struct {
	RecreationID int64 `json:"recreation_id"`
}

func (h *RecreationHandler) CreateRecreationEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := createRecreationRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		recreationID, err := h.ru.CreateRecreation(req.RecreationName, req.RecreationCity, req.RecreationImage, req.RecreationDescription, req.RecreationTimeMinute, req.RecreationPrice, req.PositionLat, req.PositionLong)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := createRecreationResponse{
			RecreationID: recreationID,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success create recreation"}, processTime, resp)
	}
}

type dataRecreationResponse struct {
	Recreation *model.Recreation `json:"recreation"`
}

func (h *RecreationHandler) GetRecreationEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		recreationID, err := strconv.ParseInt(c.Param("recreation_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		recreation, err := h.ru.GetRecreation(recreationID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataRecreationResponse{
			Recreation: recreation,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get recreation"}, processTime, resp)
	}
}

type dataRecreationsResponse struct {
	Recreations model.Recreations `json:"recreations"`
}

func (h *RecreationHandler) GetAllRecreationsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		recreations, err := h.ru.GetAllRecrations()
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataRecreationsResponse{
			Recreations: recreations,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get all recreations"}, processTime, resp)
	}
}

type getRecreationByCityRequest struct {
	City string `json:"recreation_city" form:"recreation_city"`
}

func (h *RecreationHandler) GetRecreationsByCityEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := getRecreationByCityRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		recreations, err := h.ru.GetRecreationsByCity(req.City)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get recreations by city"}, processTime, recreations)
	}
}

func (h *RecreationHandler) DeleteRecreationEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		recreationID, err := strconv.ParseInt(c.Param("recreation_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		err = h.ru.DeleteRecreationByID(recreationID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success delete recreation"}, processTime, nil)
	}
}
