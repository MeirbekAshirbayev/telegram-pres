package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	BotToken    string  `json:"bot_token"`
	BotUsername string  `json:"bot_username"`
	AdminIDs    []int64 `json:"admin_ids"`
	Port        string  `json:"port"`
	DBPath      string  `json:"db_path"`
	BaseURL     string  `json:"base_url"`
}

var appConfig Config
var templates *template.Template

func loadConfig() {
	// Try loading from file first
	file, err := os.Open("config.json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		_ = decoder.Decode(&appConfig)
	}

	// Override with Environment Variables (for Render/Cloud)
	if envPort := os.Getenv("PORT"); envPort != "" {
		appConfig.Port = ":" + envPort
	}
	if envToken := os.Getenv("BOT_TOKEN"); envToken != "" {
		appConfig.BotToken = envToken
	}
	if envUsername := os.Getenv("BOT_USERNAME"); envUsername != "" {
		appConfig.BotUsername = envUsername
	}
	if envURL := os.Getenv("BASE_URL"); envURL != "" {
		appConfig.BaseURL = envURL
	}
	// Handle AdminIDs from Env (comma separated) if needed, but for now config.json is fine or we rely on seed
}

func main() {
	loadConfig()

	// check for flags
	if len(os.Args) > 1 && os.Args[1] == "-seed" {
		seedDB()
		return
	}

	initDB(appConfig.DBPath)

	// Auto-seed if specified (useful for ephemeral cloud storage)
	if os.Getenv("AUTO_SEED") == "true" {
		seedDB()
	}

	initBot(appConfig.BotToken)

	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", dashboardHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/auth/telegram", telegramAuthHandler).Methods("GET")
	r.HandleFunc("/auth/telegram", telegramAuthHandler).Methods("GET")
	r.HandleFunc("/view/{id}", viewHandler).Methods("GET")
	r.HandleFunc("/post/{id}", postHandler).Methods("POST") // Use POST for actions

	log.Printf("Server starting on %s", appConfig.Port)
	err = http.ListenAndServe(appConfig.Port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	redirectTo := r.URL.Query().Get("redirect_to")
	authURL := appConfig.BaseURL + "/auth/telegram"
	if redirectTo != "" {
		authURL += "?redirect_to=" + url.QueryEscape(redirectTo)
	}

	data := struct {
		BotUsername string
		AuthURL     string
	}{
		BotUsername: appConfig.BotUsername,
		AuthURL:     authURL,
	}
	templates.ExecuteTemplate(w, "login.html", data)
}

func checkTelegramAuth(query url.Values, token string) (map[string]string, error) {
	// Check if hash is present
	hash := query.Get("hash")
	if hash == "" {
		return nil, errors.New("hash is missing")
	}

	// Create data-check-string
	var args []string
	for k, v := range query {
		if k == "hash" || k == "redirect_to" {
			continue
		}
		args = append(args, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(args)
	dataCheckString := strings.Join(args, "\n")

	// Compute secret key
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(token))
	secretKey := sha256Hash.Sum(nil)

	// Compute HMAC-SHA256 signature
	hmacHash := hmac.New(sha256.New, secretKey)
	hmacHash.Write([]byte(dataCheckString))
	signature := hex.EncodeToString(hmacHash.Sum(nil))

	// DEBUG LOGGING
	if signature != hash {
		log.Printf("Auth FAILED. Expected: %s, Got: %s", signature, hash)
		log.Printf("Data Check String was:\n%s", dataCheckString)
	}

	// Compare signatures
	if signature != hash {
		return nil, errors.New("signature mismatch")
	}

	// Check auth_date
	authDateStr := query.Get("auth_date")
	var authDate int64
	fmt.Sscanf(authDateStr, "%d", &authDate)
	if time.Now().Unix()-authDate > 86400 {
		return nil, errors.New("auth data is outdated")
	}

	user := make(map[string]string)
	user["id"] = query.Get("id")
	user["first_name"] = query.Get("first_name")
	user["username"] = query.Get("username")
	return user, nil
}

func telegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	user, err := checkTelegramAuth(r.URL.Query(), appConfig.BotToken)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Set a cookie or session here for the user
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user["id"],
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // Required for SameSite=None
		SameSite: http.SameSiteNoneMode, // Allows cross-site cookie for Telegram redirect
		MaxAge:   3600 * 24 * 30,        // 30 days
	})

	// Redirect to home or dashboard (or the requested page)
	redirectTo := r.URL.Query().Get("redirect_to")
	if redirectTo != "" {
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Explicitly ignore cookie usage error by logging (or just using it if needed later)
	_ = cookie

	presentations, err := getAllPresentations()
	if err != nil {
		http.Error(w, "Error fetching presentations", http.StatusInternalServerError)
		return
	}

	// Group presentations
	grouped := make(map[string][]Presentation)
	for _, p := range presentations {
		grouped[p.GroupName] = append(grouped[p.GroupName], p)
	}

	data := struct {
		User                 struct{ FirstName string }
		GroupedPresentations map[string][]Presentation
	}{
		User:                 struct{ FirstName string }{FirstName: "User"},
		GroupedPresentations: grouped,
	}
	templates.ExecuteTemplate(w, "dashboard.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get user from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Deep link redirection: Send them to login, then bring them back here
		returnURL := "/view/" + id
		http.Redirect(w, r, "/login?redirect_to="+url.QueryEscape(returnURL), http.StatusSeeOther)
		return
	}

	userIDStr := cookie.Value
	var userID int64
	fmt.Sscanf(userIDStr, "%d", &userID)

	presentation, err := getPresentation(id)
	if err != nil {
		http.Error(w, "Presentation not found", http.StatusNotFound)
		return
	}

	allowed, err := isUserMember(userID, presentation.AllowedChannelID)
	if err != nil {
		log.Printf("Error checking membership: %v", err)
		http.Error(w, "Error checking permissions", http.StatusInternalServerError)
		return
	}

	// Prevent caching so permissions are checked every time
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if !allowed {
		http.Error(w, "You are not a member of the required channel to view this presentation.", http.StatusForbidden)
		return
	}

	templates.ExecuteTemplate(w, "viewer.html", map[string]string{
		"Title":    presentation.Title,
		"CanvaURL": presentation.CanvaEmbedURL,
	})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Security: In a real app, check if the requester is an Admin
	// Here we assume only admins access the dashboard.

	presentation, err := getPresentation(id)
	if err != nil {
		http.Error(w, "Presentation not found", http.StatusNotFound)
		return
	}

	// Construct the View URL
	viewURL := appConfig.BaseURL + "/view/" + id

	// Send to Telegram
	err = PostPresentationToChannel(presentation.AllowedChannelID, presentation.Title, viewURL)
	if err != nil {
		log.Printf("Error posting to channel: %v", err)
		http.Error(w, "Failed to post to Telegram: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to dashboard with success message (or just back)
	http.Redirect(w, r, "/?posted=true", http.StatusSeeOther)
}
