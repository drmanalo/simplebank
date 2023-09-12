package api

import (
	"net/http"
	"time"

	db "github.com/drmanalo/simplebank/db/sqlc"
	"github.com/drmanalo/simplebank/util"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,alphanum"`
}

type createUserResponse struct {
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Username:       req.Username,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Email:     user.Email,
		FullName:  user.FullName,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time,
	}
	ctx.JSON(http.StatusCreated, resp)
}
