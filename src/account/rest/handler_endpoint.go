package rest

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/model"
	"github.com/sportivaid/go-template/util/httputil"
)

func NewAccountHandler(au account.Usecase) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/account", CreateAccountEndpoint(au))
		v1.GET("/account/:account_id", GetAccountEndpoint(au))
		v1.GET("/accounts", GetAccountsEndpoint(au))
		v1.PUT("/account/:account_id", UpdateAccountEndpoint(au))
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

func CreateAccountEndpoint(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := createAccountRequest{}
		if err := httputil.DecodeFormRequest(c.Request, &req); err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteDecodeErrorResponse(c, processTime, &req)
			return
		}

		accountID, err := au.CreateAccount(req.Email, req.Fullname)
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

func GetAccountEndpoint(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		account, err := au.GetAccount(accountID)
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

func GetAccountsEndpoint(au account.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		accounts, err := au.GetAccounts()
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

func UpdateAccountEndpoint(au account.Usecase) gin.HandlerFunc {
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

		err = au.UpdateAccount(accountID, req.Email, req.Fullname)
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
