// cntr/auth.go
package cntr

import (
	"net/http"
	"net/url"
	"xi/app/lib"
	"xi/app/model"
	"xi/app/cfg"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthCntr struct {
	db        *gorm.DB
	rdb       *redis.Client
	jwtSecret []byte
}

// Singleton controller
var Auth = &AuthCntr{
	db:        lib.DB.GetCli(),
	jwtSecret: []byte("supersecretkey"),
}

// signup/signout
func (a *AuthCntr) Signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user model.User

	if err := a.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Compare passwords (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := lib.Auth.GenToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthCntr) ShowSignup(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // auth ID or slug
	refKey := "/auth/" + url.QueryEscape(rawUID) + "/" + url.QueryEscape(rawID)

	// Try cache
	lib.View.RenderCache(c, "page/auth", refKey)

	// Prep data
	P := cfg.View.Pages["auth"]
	P.Data["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, "page/auth", refKey, P)
}

func (a *AuthCntr) ShowSignout(c *gin.Context) {}

func (a *AuthCntr) Signout(c *gin.Context) {}
 
// login/logout
func (a *AuthCntr) Logins(c *gin.Context) {
	refKey := "/auth"

	// Try cache
	lib.View.RenderCache(c, "page/auth", refKey)

	// Build data
	P := cfg.View.Pages["auths"]
	P.Data["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, "page/auths", refKey, P)
}

func (a *AuthCntr) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user model.User

	if err := a.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Compare passwords (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := lib.Auth.GenToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthCntr) ShowLogin(c *gin.Context) {
	rawUID := c.Param("uid") // @username or UID
	rawID := c.Param("id")   // auth ID or slug
	refKey := "/auth/" + url.QueryEscape(rawUID) + "/" + url.QueryEscape(rawID)

	// Try cache
	lib.View.RenderCache(c, "page/auth", refKey)

	// Prep data
	P := cfg.View.Pages["auth"]
	P.Data["url"] = c.Request.URL.String()

	// Cache renderer
	lib.View.RenderAndCache(c, "page/auth", refKey, P)
}

func (a *AuthCntr) ShowLogout(c *gin.Context) {}

func (a *AuthCntr) Logout(c *gin.Context) {}