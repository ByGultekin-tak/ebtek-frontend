package auth

import (
	"html/template"
	"net/http"
	"github.com/gorilla/sessions"
	"time"
)

const (
	sessionName = "session-name"
	sessionDuration = 24 * time.Hour // 1 gün
)

var Store = sessions.NewCookieStore([]byte("gizli-anahtar-buraya"))

var validUser = struct {
	username string
	password string
}{
	username: "admin",
	password: "123456",
}

func init() {
	// Session store ayarları
	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // HTTPS için true yapın
		MaxAge:   int(sessionDuration.Seconds()),
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	var err error
	
	// Session kontrolü
	session, err = Store.Get(r, sessionName)
	if err != nil {
		// Session bozuksa yeni session oluştur
		Store.MaxAge(-1) // Eski session'ı sil
		session, err = Store.New(r, sessionName)
		if err != nil {
			http.Error(w, "Session oluşturulamadı", http.StatusInternalServerError)
			return
		}
	}

	// Eğer zaten giriş yapmışsa dashboard'a yönlendir
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		if username, ok := session.Values["username"].(string); ok && username != "" {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	// Session'ı temizle
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session kaydedilemedi", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/login.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		rememberMe := r.FormValue("remember") == "on"

		if username == validUser.username && password == validUser.password {
			// Yeni temiz session oluştur
			session, _ = Store.New(r, sessionName)
			
			// Session ayarlarını güncelle
			session.Values["authenticated"] = true
			session.Values["username"] = username
			
			if rememberMe {
				session.Options.MaxAge = 7 * 24 * 60 * 60 // 7 gün
			} else {
				session.Options.MaxAge = 0 // Tarayıcı kapanınca silinir
			}
			
			// Session'ı kaydet
			if err := session.Save(r, w); err != nil {
				http.Error(w, "Oturum başlatılamadı", http.StatusInternalServerError)
				return
			}
			
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		data := struct {
			Title    string
			Username string
			Error    string
		}{
			Title:    "Giriş Yap",
			Username: username, // Hatalı girişte kullanıcı adını koru
			Error:    "Geçersiz kullanıcı adı veya şifre",
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Sayfa görüntülenemedi", http.StatusInternalServerError)
		}
		return
	}

	// GET isteği için login formunu göster
	data := struct {
		Title    string
		Username string
		Error    string
	}{
		Title:    "Giriş Yap",
		Username: "",
		Error:    "",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Sayfa görüntülenemedi", http.StatusInternalServerError)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Eski session'ı temizle
	Store.MaxAge(-1)
	session, _ := Store.Get(r, sessionName)
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	session.Save(r, w)

	// Yeni boş session oluştur
	session, _ = Store.New(r, sessionName)
	session.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	session.Save(r, w)

	// Yönlendir
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Şimdilik kayıt işlemi devre dışı, login'e yönlendir
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Session kontrolü
		session, err := Store.Get(r, sessionName)
		if err != nil {
			LogoutHandler(w, r) // Geçersiz session'ı temizle ve login'e yönlendir
			return
		}

		// Oturum ve kullanıcı kontrolü
		auth, authOk := session.Values["authenticated"].(bool)
		username, userOk := session.Values["username"].(string)

		if !authOk || !auth || !userOk || username == "" {
			LogoutHandler(w, r)
			return
		}

		// İsteği işle
		next(w, r)
	}
}
