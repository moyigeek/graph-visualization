package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Node struct {
	FromPackage string `json:"from_package"`
	ToPackage   string `json:"to_package"`
	FromDepends int    `json:"from_depends"`
	ToDepends   int    `json:"to_depends"`
}

func main() {
	http.HandleFunc("/nodes", corsMiddleware(nodesHandler))
	log.Println("Server started on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	minCount, err := strconv.Atoi(r.URL.Query().Get("min_count"))

	if err != nil {
		http.Error(w, "Invalid min_count parameter", http.StatusBadRequest)
		return
	}

	view, err := strconv.Atoi(r.URL.Query().Get("view"))
	if err != nil || view < 1 || view > 5 {
		http.Error(w, "Invalid view parameter", http.StatusBadRequest)
		return
	}
	log.Println("[info] minCount: ", minCount, "view: ", getviewname(view))

	nodes, err := getNodesWithMinCount(minCount, view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[info] Nodes sent")
}

func getNodesWithMinCount(minCount int, view int) ([]Node, error) {
	dsn := "postgres://postgres:TSE9%2FdM78kyOsioH@222.20.126.219:5432/criticality_score?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	var viewName string
	switch view {
	case 1:
		viewName = "draw_arch"
	case 2:
		viewName = "draw_debian"
	case 3:
		viewName = "draw_gentoo"
	case 4:
		viewName = "draw_homebrew"
	case 5:
		viewName = "draw_nix"
	default:
		return nil, fmt.Errorf("invalid view number")
	}

	query := fmt.Sprintf(`
       SELECT FromPackage, ToPackage, FromDepends, ToDepends
       FROM %s
       WHERE FromDepends > $1 AND ToDepends > $1
   `, viewName)

	rows, err := db.Query(query, minCount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var nodes []Node
	for rows.Next() {
		var node Node
		if err := rows.Scan(&node.FromPackage, &node.ToPackage, &node.FromDepends, &node.ToDepends); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func getviewname(num int) string {
	switch num {
	case 1:
		return "draw_arch"
	case 2:
		return "draw_debian"
	case 3:
		return "draw_gentoo"
	case 4:
		return "draw_homebrew"
	case 5:
		return "draw_nix"
	default:
		return "draw_arch"
	}
}
