package handlers

import (
	"fmt"
	"net/http"

	"github.com/himtar/go-boilerplate/pkg/router"
)

func RegisterGroupsHandlers (server *router.RouterMux) {
	fmt.Println("Registering Groups Routes")

	server.AddCustomGetHandler("/group/get", getGroupsHandler)
	server.AddCustomPostHandler("/group/create", createGroupsHandler)
}

func getGroupsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Group Get Route called")
}

func createGroupsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Group Create Route called")
}