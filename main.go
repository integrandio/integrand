package main

import (
	"fmt"
	"integrand/persistence"
	"integrand/utils"
	"integrand/web"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	utils.GetEnvrionmentVariableString("DEV_MODE", "true")
	utils.GetEnvrionmentVariableString("DB_FILE_LOCATION", "integrand.db")
	utils.GetEnvrionmentVariableString("ROOT_EMAIL", "admin")
	utils.GetEnvrionmentVariableString("ROOT_PASSWORD", "admin")
	utils.GetEnvrionmentVariableString("INITIAL_API_KEY", "11111")

	// We need to generate new keys....
	utils.GetEnvrionmentVariableString("GOOGLE_CALLBACK_URL", "")
	utils.GetEnvrionmentVariableString("GOOGLE_CLIENT_KEY", "")
	utils.GetEnvrionmentVariableString("GOOGLE_CLIENT_SECRET", "")

	utils.GetEnvrionmentVariableString("GITHUB_CALLBACK_URL", "")
	utils.GetEnvrionmentVariableString("GITHUB_CLIENT_ID", "")
	utils.GetEnvrionmentVariableString("GITHUB_CLIENT_SECRET", "")
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	persistence.Initialize()
	router := web.NewNewWebRouter()
	port := ":8000"
	slog.Info(fmt.Sprintf("Server started on http//:localhost%s\n", port))
	log.Fatal(http.ListenAndServe(port, router))
}
