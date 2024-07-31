package handler

import (
	"auth-service/api/token"
	"auth-service/models"
	"auth-service/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AutheticationHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

type autheticationHandlerImpl struct {
	userService service.UserService
	authService service.AuthService
	logger      *slog.Logger
}

func NewAutheticationHandler(serviceManager service.ServiceManager, logger *slog.Logger) AutheticationHandler {
	return &autheticationHandlerImpl{
		userService: serviceManager.User(),
		authService: serviceManager.Auth(),
		logger:      logger,
	}
}

// @Summary Login user
// @Description Login to exist user acc
// @Tags Auth
// @Accept json
// @Produce json
// @Param Login body models.LoginRequest true "user login"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /api/login [post]
func (h *autheticationHandlerImpl) Login(c *gin.Context) {
	// Implement login logic
	loginRequest := models.LoginRequest{}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.logger.Error("Error binding JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Errors{
			Error: "json binding error",
		})
		return
	}

	user, err := h.authService.Login(c, loginRequest)
	if err != nil {
		h.logger.Error("Error in login user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error invalid username or password",
		})
		return
	}

	is := token.VerifyPassword(loginRequest.Password, user.Password)
	if !is {
		h.logger.Error("Error verifying password")
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "password is incorrect",
		})
		return
	}

	newAccessToken, err := token.GENERATEAccessTokenJWT(user)
	if err != nil {
		h.logger.Error("Error generating access token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error generating access token",
		})
		return
	}

	newRefreshToken, err := token.GENERATERefreshTokenJWT(user)
	if err != nil {
		h.logger.Error("Error generating refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error generating refresh token",
		})
		return
	}

	err = h.authService.SaveRefreshToken(c, models.TokenRequest{
		Username:     loginRequest.Username,
		RefreshToken: newRefreshToken,
		ExpiresAt:    token.GetExpirationTime(),
	})

	if err != nil {
		h.logger.Error("Error saving refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error saving refresh token",
		})
		return
	}

	c.SetCookie("access_token", newAccessToken, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}

// @Summary Register user
// @Description create user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body models.Register true "user register"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /api/register [post]
func (h *autheticationHandlerImpl) Register(c *gin.Context) {
	// Implement registration logic

	user := models.Register{}
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Error binding JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Errors{
			Error: "json binding error",
		})
		return
	}

	resp, err := h.authService.Register(c, user)
	if err != nil {
		h.logger.Error("Error registering user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error in registering user",
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Logout user
// @Description Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Log out user"
// @Success 201 {object} string
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /api/logout [post]
func (h *autheticationHandlerImpl) Logout(c *gin.Context) {
	// Implement logout logic

	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		h.logger.Error("Error getting refresh token from header")
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "unauthorized",
		})
		return
	}

	claims, err := token.ExtractClaimsRefresh(refreshToken)
	if err != nil {
		h.logger.Error("Error extracting claims from refresh token", "error", err.Error())
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "invalid token",
		})
		return
	}
	isValid, err := h.authService.IsRefreshTokenValid(c, claims.Username)
	if !isValid || err != nil {
		h.logger.Error("Error invalidating refresh token", "error", err.Error())
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "invalid refresh token",
		})
		return
	}

	err = h.authService.InvalidateRefreshToken(c, claims.Username)
	if err != nil {
		h.logger.Error("Error invalidating refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error invalidating refresh token",
		})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// @Summary Refresh user
// @Description Get access token user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Acces Token"
// @Success 201 {object} string
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Security ApiKeyAuth
// @Router /api/refresh [post]
func (h *autheticationHandlerImpl) Refresh(c *gin.Context) {
	// Implement token refresh logic
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		h.logger.Error("Error getting refresh token from header")
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "unauthorized",
		})
		return
	}
	claims, err := token.ExtractClaimsRefresh(refreshToken)
	if err != nil {
		h.logger.Error("Error extracting claims from refresh token", "error", err.Error())
		c.JSON(http.StatusUnauthorized, models.Errors{
			Error: "invalid token",
		})
		return
	}

	newAccessToken, err := token.GENERATEAccessTokenJWT(&models.User{
		Username: claims.Username,
		ID:       claims.ID,
		Role:     claims.Role,
	})

	if err != nil {
		h.logger.Error("Error generating access token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Errors{
			Error: "error generating access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"message":      "Successfully refreshed token",
	})
}
