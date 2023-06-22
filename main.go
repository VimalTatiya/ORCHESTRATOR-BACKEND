package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type RecordCluster struct {
	Cluster string `json:"cluster"`
}
type RecordNamespace struct {
	Namespace string `json:"namespace"`
}
type RecordDeployment struct {
	Deployment string `json:"deployment"`
}
type DetailsCaseOne struct {
	Cluster1                string `json:"cluster"`
	Namespace1              string `json:"namespace"`
	Deployment1             string `json:"deployment"`
	Container1              string `json:"container"`
	Updated_enabled1        int    `json:"updated_enabled"`
	Optimization_available1 int    `json:"optimization_available"`
	Current_cpu_request1    int64  `json:"current_cpu_request"`
	Current_mem_request1    int64  `json:"current_mem_request"`
	Recommend_cpu_request1  int64  `json:"recommend_cpu_request"`
	Recommend_mem_request1  int64  `json:"recommend_mem_request"`
}
type DetailsCaseTwo struct {
	Cluster2                string `json:"cluster"`
	Namespace2              string `json:"namespace"`
	Deployment2             string `json:"deployment"`
	Container2              string `json:"container"`
	Updated_enabled2        int    `json:"updated_enabled"`
	Optimization_available2 int    `json:"optimization_available"`
	Current_cpu_request2    int64  `json:"current_cpu_request"`
	Current_mem_request2    int64  `json:"current_mem_request"`
	Recommend_cpu_request2  int64  `json:"recommend_cpu_request"`
	Recommend_mem_request2  int64  `json:"recommend_mem_request"`
}
type DetailsCaseThree struct {
	Cluster3                string `json:"cluster"`
	Namespace3              string `json:"namespace"`
	Deployment3             string `json:"deployment"`
	Container3              string `json:"container"`
	Updated_enabled3        int    `json:"updated_enabled"`
	Optimization_available3 int    `json:"optimization_available"`
	Current_cpu_request3    int64  `json:"current_cpu_request"`
	Current_mem_request3    int64  `json:"current_mem_request"`
	Recommend_cpu_request3  int64  `json:"recommend_cpu_request"`
	Recommend_mem_request3  int64  `json:"recommend_mem_request"`
}

// To establish connection a connection with database:
func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := mux.NewRouter()

	router.HandleFunc("/api/clusters", getClusterHandler).Methods("GET")
	log.Println("Server started on port 8000")

	router.HandleFunc("/api/namespaces", getNamespaceHandler).Queries("cluster", "{cluster}")
	log.Println("Server started on port 8000")

	router.HandleFunc("/api/deployments", getDeploymentHandler).Queries("cluster", "{cluster}", "namespace", "{namespace}")
	log.Println("Server started on port 8000")

	router.HandleFunc("/api/container1", getCaseOneHandler).Queries("cluster", "{cluster}")
	log.Println("Server started on port 8000")

	router.HandleFunc("/api/container2", getCaseTwoHandler).Queries("cluster", "{cluster}", "namespace", "{namespace}")
	log.Println("Server started on port 8000")

	router.HandleFunc("/api/container3", getCaseThreeHandler).Queries("cluster", "{cluster}", "namespace", "{namespace}", "deployment", "{deployment}")
	log.Println("Server started on port 8000")

	c := cors.Default()
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}

// To fetch Cluster Data
func getClusterHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT DISTINCT cluster_name FROM metadata")
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch cluster data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataCluster []RecordCluster
	for rows.Next() {
		var data RecordCluster
		err := rows.Scan(&data.Cluster)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataCluster = append(dataCluster, data)
	}

	response, err := json.Marshal(dataCluster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Fetching namespace data from data base
func getNamespaceHandler(w http.ResponseWriter, r *http.Request) {

	selectedCluster := r.URL.Query().Get("cluster")

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT DISTINCT namespace FROM metadata WHERE cluster_name = '%s'", selectedCluster)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch namespaces data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataNamespace []RecordNamespace
	for rows.Next() {
		var data RecordNamespace
		err := rows.Scan(&data.Namespace)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataNamespace = append(dataNamespace, data)
	}

	response, err := json.Marshal(dataNamespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Fetching deployment data from data base
func getDeploymentHandler(w http.ResponseWriter, r *http.Request) {

	selectedCluster := r.URL.Query().Get("cluster")
	selectedNamespace := r.URL.Query().Get("namespace")

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT DISTINCT deployment FROM metadata WHERE cluster_name = '%s' AND namespace='%s'", selectedCluster, selectedNamespace)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch deployment data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataDeployment []RecordDeployment
	for rows.Next() {
		var data RecordDeployment
		err := rows.Scan(&data.Deployment)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataDeployment = append(dataDeployment, data)
	}

	response, err := json.Marshal(dataDeployment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// To fetch the data when only first dropdown is selected
func getCaseOneHandler(w http.ResponseWriter, r *http.Request) {

	selectedCluster := r.URL.Query().Get("cluster")

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT cluster_name,namespace,deployment,container,managed_by_orchestrator,optimization_available,cpu_request,mem_request,cpu_recommendation,mem_recommendation FROM metadata WHERE cluster_name = '%s'", selectedCluster)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch deployment data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataCaseOne []DetailsCaseOne
	for rows.Next() {
		var data DetailsCaseOne
		err := rows.Scan(&data.Cluster1, &data.Namespace1, &data.Deployment1, &data.Container1, &data.Updated_enabled1,
			&data.Optimization_available1, &data.Current_cpu_request1, &data.Current_mem_request1, &data.Recommend_cpu_request1, &data.Recommend_mem_request1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataCaseOne = append(dataCaseOne, data)
	}

	response, err := json.Marshal(dataCaseOne)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// To fetch the data when only first and second dropdown is selected
func getCaseTwoHandler(w http.ResponseWriter, r *http.Request) {

	selectedCluster := r.URL.Query().Get("cluster")
	selectedNamespace := r.URL.Query().Get("namespace")

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT cluster_name,namespace,deployment,container,managed_by_orchestrator,optimization_available,cpu_request,mem_request,cpu_recommendation,mem_recommendation FROM metadata WHERE cluster_name = '%s' AND namespace = '%s'", selectedCluster, selectedNamespace)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch deployment data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataCaseTwo []DetailsCaseTwo
	for rows.Next() {
		var data DetailsCaseTwo
		err := rows.Scan(&data.Cluster2, &data.Namespace2, &data.Deployment2, &data.Container2, &data.Updated_enabled2,
			&data.Optimization_available2, &data.Current_cpu_request2, &data.Current_mem_request2, &data.Recommend_cpu_request2, &data.Recommend_mem_request2)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataCaseTwo = append(dataCaseTwo, data)
	}

	response, err := json.Marshal(dataCaseTwo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// To fetch the data when only first second and third dropdown is selected
func getCaseThreeHandler(w http.ResponseWriter, r *http.Request) {

	selectedCluster := r.URL.Query().Get("cluster")
	selectedNamespace := r.URL.Query().Get("namespace")
	selectedDeployment := r.URL.Query().Get("deployment")

	db, err := sql.Open("mysql", "stage_orchestrator:123456@tcp(localhost:3306)/stage_orchestrator")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		log.Println(err)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT cluster_name,namespace,deployment,container,managed_by_orchestrator,optimization_available,cpu_request,mem_request,cpu_recommendation,mem_recommendation FROM metadata WHERE cluster_name = '%s' AND namespace = '%s' AND deployment = '%s'", selectedCluster, selectedNamespace, selectedDeployment)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch deployment data from the database")
		log.Println(err)
		return
	}
	defer rows.Close()

	var dataCaseThree []DetailsCaseThree
	for rows.Next() {
		var data DetailsCaseThree
		err := rows.Scan(&data.Cluster3, &data.Namespace3, &data.Deployment3, &data.Container3, &data.Updated_enabled3,
			&data.Optimization_available3, &data.Current_cpu_request3, &data.Current_mem_request3, &data.Recommend_cpu_request3, &data.Recommend_mem_request3)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to scan data from the database")
			log.Println(err)
			return
		}
		dataCaseThree = append(dataCaseThree, data)
	}

	response, err := json.Marshal(dataCaseThree)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate JSON response")
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
