package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	"todo-app"
	"todo-app/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

const (

	// Константа, которую приложение добавляет при хешировании
	salt = "lhkljGE65E&^$JHGEfgiouero23hg3hj4j"

	// Константа, которую приложение добавляет при генерации JWT токенов
	signingKey = "dsl*@(H#*@H#*Dr%4hr"

	// Время жизни JWT токена
	tokenTTL = 12 * time.Hour
)

// Обертка над jwt.StandartClaims с дополнительным полем
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// Структура сервиса аутентификации
type AuthService struct {
	repo repository.Authorization
}

// Конструктор
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// Метод передачи структуры пользователя в слой repository
func (s *AuthService) CreateUser(user todo.User) (int, error) {

	// Хэширует пароль пользователя
	user.Password = generatePasswordHash(user.Password)

	// Передаёт структуру в слой repository
	return s.repo.CreateUser(user)
}

// Метод генерации JWT токена
func (s *AuthService) GenerateToken(username, password string) (string, error) {

	// Получение пользователя из БД
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	// Проврека на получение
	if err != nil {
		return "", err
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{

			// Задаёт время жизни токена
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),

			// Записывает время генерации токена
			IssuedAt: time.Now().Unix(),
		},

		// Дополнительное поле с id пользователя
		user.Id,
	})

	// Возвращает токен и ошибку
	return token.SignedString([]byte(signingKey))
}

// Метод парсинга JWT токена
func (s *AuthService) ParseToken(accessToken string) (int, error) {

	// Парсинг токена
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		// Если метод подписи токена не HMAC, возвращает ошибку
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// Возвращает ключ подписи
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	// Приводит полученные данные к обёртке tokenClaims
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	// Возвращает id из токена
	return claims.UserId, nil
}

// Функция хеширования пароля
func generatePasswordHash(password string) string {

	// Инициализация структуры для хеширования
	hash := sha1.New()

	// Хеширование пароля
	hash.Write([]byte(password))

	// Возвращает строку с хешом
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
