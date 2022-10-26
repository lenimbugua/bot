package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
)

type createBotRequest struct {
	Title     string `json:"title" binding:"required"`
	CompanyID int64  `json:"company_id" binding:"required,min=0"`
}

func (server *Server) createBot(ctx *gin.Context) {
	var req createBotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	company, err := server.dbStore.GetCompanyByID(ctx, req.CompanyID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateBotParams{
		Title:     req.Title,
		CompanyID: company.ID,
	}
	bot, err := server.dbStore.CreateBot(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bot)
}

type updateBotRequestURI struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateBotRequestParams struct {
	Title     string `json:"title" binding:"required"`
	CompanyID int64  `json:"company_id" binding:"required,min=1"`
}

func (server *Server) updateBot(ctx *gin.Context) {
	var uri updateBotRequestURI

	var req updateBotRequestParams

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if req.CompanyID != authPayload.CompanyID {
		err := errors.New("bot doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateBotParams{
		Title:     sql.NullString{String: req.Title, Valid: true},
		CompanyID: sql.NullInt64{Int64: req.CompanyID, Valid: true},
		ID:        uri.ID,
	}

	bot, err := server.dbStore.UpdateBot(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, bot)
}

type deleteBotRequestURI struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBot(ctx *gin.Context) {
	var req deleteBotRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bot, err := server.dbStore.GetBot(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if bot.CompanyID != authPayload.CompanyID {
		err := errors.New("You cannot delete a bot that does not belong to your company")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.dbStore.DeleteBot(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)

}

type listBotsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listBots(ctx *gin.Context) {
	var req listBotsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAllBotsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	bots, err := server.dbStore.ListAllBots(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, bots)
}

type listCompanyBotsRequest struct {
	CompanyID int64 `form:"company_id" binding:"required,min=1"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCompanyBots(ctx *gin.Context) {
	var req listCompanyBotsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompanyBotsParams{
		CompanyID: req.CompanyID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	bots, err := server.dbStore.ListCompanyBots(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, bots)
}

type getBotRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBot(ctx *gin.Context) {
	var req getBotRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bot, err := server.dbStore.GetBot(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if bot.CompanyID != authPayload.CompanyID {
		err := errors.New("No access rights for that company")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, bot)
}
