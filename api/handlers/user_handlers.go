package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"

	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		db: db,
		q:  queries.New(db),
	}
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.q.AddUser(c.Request.Context(), queries.AddUserParams{
		Username:       req.Username,
		HashedPassword: string(hashedPassword[:]),
		Email:          req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagesize, err := strconv.ParseInt(c.DefaultQuery("page_size", "50"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.q.GetAllUsers(c.Request.Context(), queries.GetAllUsersParams{
		Page:     int32(page - 1),
		PageSize: int32(pagesize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) GetUserByUuid(c *gin.Context) {
	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.q.GetUserByUuid(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		NewUsername *string `json:"username"`
		NewName     *string `json:"full_name"`
		NewBio      *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := queries.UpdateUserDetailsParams{
		UserUuid: pgUuid,
	}

	// If no username is provided, we need to fetch the current username so it doesn't get overwritten
	if req.NewUsername != nil {
		params.NewUsername = *req.NewUsername
	} else {
		currentUser, err := h.q.GetUserByUuid(c.Request.Context(), pgUuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch current user details"})
			return
		}
		params.NewUsername = currentUser.Username
	}

	// Ensure that any missing fields don't overwrite the current field
	if req.NewName != nil {
		params.NewName = pgtype.Text{String: *req.NewName, Valid: true}
	} else {
		params.NewName = pgtype.Text{Valid: false}
	}

	if req.NewBio != nil {
		params.NewBio = pgtype.Text{String: *req.NewBio, Valid: true}
	} else {
		params.NewBio = pgtype.Text{Valid: false}
	}

	// Update the user details
	user, err := h.q.UpdateUserDetails(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.q.DeleteUser(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "user deleted"})
}
