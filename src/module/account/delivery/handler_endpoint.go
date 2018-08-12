package delivery

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/atletaid/go-template/src/common/auth"
	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/account"
	"github.com/atletaid/go-template/util/httputil"
)

type AccountHandler struct {
	au account.Usecase
}

func NewAccountHandler(router *gin.Engine, m *auth.Middleware, au account.Usecase) *gin.Engine {
	handler := &AccountHandler{au}

	v1 := router.Group("/api/v1")
	v1.POST("/account", handler.CreateAccountEndpoint())
	v1.GET("/account/:account_id", handler.GetAccountEndpoint())

	v1.Use(m.AuthUserToken())
	{
		v1.GET("/accounts", handler.GetAccountsEndpoint())
		v1.PUT("/account/:account_id", handler.UpdateAccountEndpoint())
	}

	return router
}

type createAccountRequest struct {
	Email    string `json:"user_email" form:"user_email"`
	Fullname string `json:"user_fullname" form:"user_fullname"`
}

type createAccountResponse struct {
	AccountID int64 `json:"account_id"`
}

func (h *AccountHandler) CreateAccountEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := createAccountRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		accountID, err := h.au.CreateAccount(req.Email, req.Fullname)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := createAccountResponse{
			AccountID: accountID,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success create account"}, processTime, resp)
	}
}

type dataAccountResponse struct {
	Account *model.Account `json:"account"`
}

func (h *AccountHandler) GetAccountEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		account, err := h.au.GetAccount(accountID)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataAccountResponse{
			Account: account,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get account"}, processTime, resp)
	}
}

type dataAccountsResponse struct {
	Accounts model.Accounts `json:"accounts"`
}

func (h *AccountHandler) GetAccountsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accounts, err := h.au.GetAccounts()
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		resp := dataAccountsResponse{
			Accounts: accounts,
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success get accounts"}, processTime, resp)
	}
}

type updateAccountRequest struct {
	Email    string `json:"user_email" form:"user_email"`
	Fullname string `json:"user_fullname" form:"user_fullname"`
}

func (h *AccountHandler) UpdateAccountEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		req := updateAccountRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		err = h.au.UpdateAccount(accountID, req.Email, req.Fullname)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		processTime := time.Now().Sub(startTime).Seconds()
		httputil.WriteResponse(c, []string{"Success update account"}, processTime, nil)
	}
}
