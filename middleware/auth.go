package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/project/pkg/models"
)

var jwtKey = []byte("sdjf2uie1wwk2hi1jqen2")

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	fmt.Println("Парсим токен:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Неверный метод подписи")
			return nil, errors.New("неверный метод подписи")
		}
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println("Ошибка парсинга токена:", err)
		return nil, errors.New("неверный или просроченный токен")
	}

	if !token.Valid {
		fmt.Println("Токен не валиден")
		return nil, errors.New("токен не валиден")
	}

	fmt.Println("Токен успешно распарсен")

	return token, nil
}

func IsValidEmployee(tokenString string) (bool, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		fmt.Println("Ошибка при разборе токена:", err)
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Ошибка: неверный формат claims")
		return false, errors.New("неверный формат данных в токене")
	}

	fmt.Println("Claims:", claims)

	role, ok := claims["role"].(string)
	if !ok {
		fmt.Println("Ошибка: не удалось получить роль из claims")
		return false, errors.New("не удалось получить роль")
	}

	if role != "employee" {
		fmt.Println("Ошибка: роль не employee:", role)
		return false, errors.New("неверная роль в токене")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		fmt.Println("Ошибка: не удалось получить время жизни токена")
		return false, errors.New("ошибка времени действия токена")
	}

	if time.Now().Unix() > int64(exp) {
		fmt.Println("Токен просрочен")
		return false, errors.New("токен просрочен")
	}

	fmt.Println("Роль employee подтверждена, токен действителен")
	return true, nil
}

func EmployeeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		fmt.Printf("Получен заголовок Authorization: '%s'\n", c.GetHeader("Authorization"))
		fmt.Println("Получен токен:", token)

		if token == "" {
			fmt.Println("Ошибка: токен отсутствует")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Доступ запрещен.",
			})
			c.Abort()
			return
		}

		valid, err := IsValidEmployee(token)
		if err != nil {
			fmt.Println("Ошибка при валидации токена:", err)
		}
		if err != nil || !valid {
			fmt.Println("Ошибка: токен недействителен или роль не employee")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Доступ запрещен.",
			})
			c.Abort()
			return
		}

		fmt.Println("Доступ разрешён")
		c.Next()
	}
}
