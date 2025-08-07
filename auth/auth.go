package auth

import (
	"html/template"
	"net/http"
)

type LoginData struct {
	Error string
}

// Basit kullanıcı doğrulama için örnek kullanıcı bilgileri
var validUser = struct {
	username string
	password string
}{
	username: "admin",
	password: "123456",
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == validUser.username && password == validUser.password {
			http.SetCookie(w, &http.Cookie{
				Name:  "authenticated",
				Value: "true",
				Path:  "/",
			})
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		tmpl.Execute(w, LoginData{Error: "Geçersiz kullanıcı adı veya şifre"})
		return
	}

	tmpl.Execute(w, nil)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authenticated")
		if err != nil || cookie.Value != "true" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
