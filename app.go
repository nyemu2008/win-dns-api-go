package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

// DoDNSSet Set
func DoDNSSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zoneName, dnsType, nodeName, Address := vars["zoneName"], vars["dnsType"], vars["nodeName"], vars["Address"]

	// Validate DNS Type
	if dnsType == "A" {
		// Validate Ip Address
		var validIPAddress = regexp.MustCompile(`^(([1-9]?\d|1\d\d|25[0-5]|2[0-4]\d)\.){3}([1-9]?\d|1\d\d|25[0-5]|2[0-4]\d)$`)

		if !validIPAddress.MatchString(Address) {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid IP address ('" + Address + "'). Currently, only IPv4 addresses are accepted."})
			return
		}
	}

	// Validate DNS Type
	var validZoneName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validZoneName.MatchString(zoneName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid zone name ('" + zoneName + "'). Zone names can only contain letters, numbers, dashes (-), and dots (.)."})
		return
	}

	// Validate Node Name
	var validNodeName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validNodeName.MatchString(nodeName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid node name ('" + nodeName + "'). Node names can only contain letters, numbers, dashes (-), and dots (.)."})
		return
	}

	dnsAddRecord := exec.Command("cmd", "/C", "dnscmd /recordadd "+zoneName+" "+nodeName+" "+dnsType+" "+Address)

	if err := dnsAddRecord.Run(); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Add record failed, error was: " + err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "The alias " + dnsType + " record " + nodeName + "." + zoneName + "' was successfully updated to '" + Address + "'."})
}

// EditDNSSet Set
func EditDNSSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zoneName, dnsType, nodeName, Address := vars["zoneName"], vars["dnsType"], vars["nodeName"], vars["Address"]

	// Validate DNS Type
	if dnsType == "A" {
		// Validate Ip Address
		var validIPAddress = regexp.MustCompile(`^(([1-9]?\d|1\d\d|25[0-5]|2[0-4]\d)\.){3}([1-9]?\d|1\d\d|25[0-5]|2[0-4]\d)$`)

		if !validIPAddress.MatchString(Address) {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid IP address ('" + Address + "'). Currently, only IPv4 addresses are accepted."})
			return
		}
	}

	// Validate DNS Type
	var validZoneName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validZoneName.MatchString(zoneName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid zone name ('" + zoneName + "'). Zone names can only contain letters, numbers, dashes (-), and dots (.)."})
		return
	}

	// Validate Node Name
	var validNodeName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validNodeName.MatchString(nodeName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid node name ('" + nodeName + "'). Node names can only contain letters, numbers, dashes (-), and dots (.)."})
		return
	}

	dnsCmdDeleteRecord := exec.Command("cmd", "/C", "dnscmd /recorddelete "+zoneName+" "+nodeName+" "+dnsType+" /f")

	if err := dnsCmdDeleteRecord.Run(); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Edit record failed, error was: " + err.Error()})
		return
	}

	dnsAddDeleteRecord := exec.Command("cmd", "/C", "dnscmd /recordadd "+zoneName+" "+nodeName+" "+dnsType+" "+Address)

	if err := dnsAddDeleteRecord.Run(); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Edit record failed, error was: " + err.Error()})
		return
	}
	retMsg := fmt.Sprintf("The alias " + dnsType + " record '" + nodeName + "." + zoneName + "' was successfully updated to '" + Address + "'.")
	respondWithJSON(w, http.StatusOK, map[string]string{"message": retMsg})
}

// DoDNSRemove Remove
func DoDNSRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zoneName, dnsType, nodeName := vars["zoneName"], vars["dnsType"], vars["nodeName"]

	// Validate DNS Type
	var validZoneName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validZoneName.MatchString(zoneName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Invalid zone name ('" + zoneName + "'). Zone names can only contain letters, numbers, dashes (-), and dots (.).")})
		return
	}

	// Validate Node Name
	var validNodeName = regexp.MustCompile(`[^A-Za-z0-9\.-]+`)

	if validNodeName.MatchString(nodeName) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid node name ('" + nodeName + "'). Node names can only contain letters, numbers, dashes (-), and dots (.)."})
		return
	}

	dnsCmdDeleteRecord := exec.Command("cmd", "/C", "dnscmd /recorddelete "+zoneName+" "+nodeName+" "+dnsType+" /f")

	if err := dnsCmdDeleteRecord.Run(); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	retMsg := fmt.Sprintf("The " + dnsType + " record " + nodeName + "." + zoneName + "' was successfully removed.")
	respondWithJSON(w, http.StatusOK, map[string]string{"message": retMsg})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Could not get the requested route."})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

const (
	serverPort = 3111
)

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Welcome to Win DNS API Go"})
	})

	r.Methods("POST").Path("/dns/{zoneName}/{dnsType}/{nodeName}/set/{Address}").HandlerFunc(DoDNSSet)
	r.Methods("POST").Path("/dns/{zoneName}/{dnsType}/{nodeName}/set/{Address}").HandlerFunc(EditDNSSet)
	r.Methods("POST").Path("/dns/{zoneName}/{dnsType}/{nodeName}/remove").HandlerFunc(DoDNSRemove)
	fmt.Printf("Listening on port %d.\n", serverPort)

	// Start HTTP Server
	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", serverPort),
		r,
	); err != nil {
		log.Fatal(err)
	}
}
