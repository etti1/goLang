package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sync"
)

type User struct {
	Username string
	Password string
}

var users []User

func register(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	for _, user := range users {
		if user.Username == u.Username {
			return c.String(http.StatusConflict, "Username already exists")
		}
	}

	users = append(users, *u)

	return c.String(http.StatusCreated, "User registered successfully")
}

func login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	for _, user := range users {
		if user.Username == u.Username && user.Password == u.Password {
			return c.String(http.StatusOK, "Login successful")
		}
	}

	return c.String(http.StatusUnauthorized, "Invalid username or password")
}

func printUsername(u *User, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Username: %s\n", u.Username)
}

func registerForm(c echo.Context) error {
	return c.HTML(http.StatusOK, `
		<html>
		<head>
			<title>Registration</title>
		</head>
		<body>
			<form action="/register" method="POST">
				<label>Username:</label><br>
				<input type="text" name="username"><br>
				<label>Password:</label><br>
				<input type="password" name="password"><br>
				<input type="submit" value="Register">
			</form>
		</body>
		</html>
	`)
}

func loginForm(c echo.Context) error {
	return c.HTML(http.StatusOK, `
		<html>
		<head>
			<title>Login</title>
		</head>
		<body>
			<form action="/login" method="POST">
				<label>Username:</label><br>
				<input type="text" name="username"><br>
				<label>Password:</label><br>
				<input type="password" name="password"><br>
				<input type="submit" value="Login">
			</form>
		</body>
		</html>
	`)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	e := echo.New()

	u := User{}
	go printUsername(&u, &wg)
	// Регистрация нового пользователя
	e.POST("/register", register)

	e.POST("/login", login)

	err := e.Start(":8080")
	if err != nil {
		return
	}

}
