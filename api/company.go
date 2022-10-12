package api

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	db "github.com/lenimbugua/bot/db/sqlc"
// )

// type createCompanyRequest struct {
// 	Phone string `json:"mobile" binding:"required,e164"`
// 	Email string `json:"email" binding:"required,email"`
// 	Name  string `json:"description" binding:"required"`
// }

// func (server *Server) createCompany(ctx *gin.Context) {
// 	var req createCompanyRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	arg := db.CreateCompanyParams{
// 		Phone: req.Phone,
// 		Email: req.Email,
// 		Name:  req.Name,
// 	}
// 	company, err := server.dbStore.CreateCompany(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, company)
// }

// // type getRoleRequest struct {
// // 	ID int32 `uri:"id" binding:"required,min=1"`
// // }

// // func (server *Server) getRole(ctx *gin.Context) {
// // 	var req getRoleRequest

// // 	if err := ctx.ShouldBindUri(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	role, err := server.dbStore.GetRole(ctx, req.ID)
// // 	if err != nil {
// // 		if err == sql.ErrNoRows {
// // 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// // 			return
// // 		}
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	ctx.JSON(http.StatusOK, role)
// // }

// // type listRolesRequest struct {
// // 	PageID   int32 `form:"page_id" binding:"required,min=1"`
// // 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// // }

// // func (server *Server) listRoles(ctx *gin.Context) {
// // 	var req listRolesRequest
// // 	if err := ctx.ShouldBindQuery(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	arg := db.ListRolesParams{
// // 		Limit:  req.PageSize,
// // 		Offset: (req.PageID - 1) * req.PageSize,
// // 	}

// // 	roles, err := server.dbStore.ListRoles(ctx, arg)
// // 	if err != nil {
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 	}
// // 	ctx.JSON(http.StatusOK, roles)
// // }

// // func (server *Server) updateRole(ctx *gin.Context) {
// // 	var uri getRoleRequest

// // 	var req createRoleRequest

// // 	if err := ctx.ShouldBindUri(&uri); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	id := uri.ID

// // 	if err := ctx.ShouldBindJSON(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	arg := db.UpdateRoleParams{
// // 		Title:       sql.NullString{String: req.Title, Valid: true},
// // 		Slug:        sql.NullString{String: req.Slug, Valid: true},
// // 		Description: sql.NullString{String: req.Description, Valid: true},
// // 		Active:      sql.NullBool{Bool: req.Active, Valid: true},
// // 		ID:          id,
// // 	}

// // 	role, err := server.dbStore.UpdateRole(ctx, arg)

// // 	if err != nil {
// // 		if err == sql.ErrNoRows {
// // 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// // 			return
// // 		}
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	ctx.JSON(http.StatusOK, role)
// // }

// // func (server *Server) deleteRole(ctx *gin.Context) {
// // 	var req getRoleRequest

// // 	if err := ctx.ShouldBindUri(&req); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// // 		return
// // 	}

// // 	err := server.dbStore.DeleteRole(ctx, req.ID)

// // 	if err != nil {
// // 		if err == sql.ErrNoRows {
// // 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// // 			return
// // 		}
// // 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// // 		return
// // 	}
// // 	ctx.JSON(http.StatusOK, req.ID)

// // }
