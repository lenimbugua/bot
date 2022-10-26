package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lib/pq"
)

type createQuestionRequest struct {
	Question       string    `json:"question" binding:"required"`
	CompanyID      int64     `json:"company_id" binding:"required",min=1`
	Type           string    `json:"type" binding:"required"`
	ParentID       int64     `json:"parent_id" binding:"min=1"`
	BotID          int64     `json:"bot_id" binding:"required",min=1`
	NextQuestionID int64     `json:"next_question_id" binding:"min=1"`
}

func (server *Server) createQuestion(ctx *gin.Context) {
	var req createQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateQuestionParams{
		Question:       req.Question,
		Type:           req.Type,
		ParentID:       req.ParentID,
		BotID:       req.BotID,
		NextQuestionID: req.NextQuestionID,
	}
	company, err := server.dbStore.CreateQuestion(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

// type getCompanyByEmailRequest struct {
// 	Email string `form:"email" binding:"required,email"`
// }

// func (server *Server) getCompanyByEmail(ctx *gin.Context) {
// 	var req getCompanyByEmailRequest

// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	company, err := server.dbStore.GetCompanyByEmail(ctx, req.Email)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, company)
// }

// type getCompanyByIDRequest struct {
// 	ID int64 `uri:"id" binding:"required,min=1"`
// }

// func (server *Server) getCompanyByID(ctx *gin.Context) {
// 	var req getCompanyByIDRequest

// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	company, err := server.dbStore.GetCompanyByID(ctx, req.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, company)
// }

// type listCompaniesRequest struct {
// 	PageID   int32 `form:"page_id" binding:"required,min=1"`
// 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// }

// func (server *Server) listCompanies(ctx *gin.Context) {
// 	var req listCompaniesRequest
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := db.ListCompaniesParams{
// 		Limit:  req.PageSize,
// 		Offset: (req.PageID - 1) * req.PageSize,
// 	}

// 	companies, err := server.dbStore.ListCompanies(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 	}
// 	ctx.JSON(http.StatusOK, companies)
// }

// type updateCompanyUri struct {
// 	ID int64 `uri:"id" binding:"required"`
// }

// type updateCompanyRequest struct {
// 	Phone string `json:"phone" binding:"required,e164"`
// 	Name  string `json:"name" binding:"required"`
// 	Email string `json:"email" binding:"required,email"`
// }

// func (server *Server) updateCompany(ctx *gin.Context) {
// 	var uri updateCompanyUri

// 	var req updateCompanyRequest

// 	if err := ctx.ShouldBindUri(&uri); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	id := uri.ID

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := db.UpdateCompanyParams{
// 		Name:      sql.NullString{String: req.Name, Valid: true},
// 		Phone:     sql.NullString{String: req.Phone, Valid: true},
// 		Email:     sql.NullString{String: req.Email, Valid: true},
// 		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: false},
// 		ID:        id,
// 	}

// 	company, err := server.dbStore.UpdateCompany(ctx, arg)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, company)
// }

// func (server *Server) deleteCompany(ctx *gin.Context) {
// 	var req getCompanyByIDRequest

// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	err := server.dbStore.DeleteCompany(ctx, req.ID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, req.ID)

// }
