package main

import (
	"fmt"
	"integrand/persistence"
	"integrand/services"
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
	utils.GetEnvrionmentVariableString("SINK_URL", "http://localhost:5000")

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	persistence.Initialize()

	// Enable our Workflower
	go services.Workflower()

	router := web.NewNewWebRouter()
	port := ":8000"
	slog.Info(fmt.Sprintf("Server started on http//:localhost%s\n", port))
	log.Fatal(http.ListenAndServe(port, router))
}
