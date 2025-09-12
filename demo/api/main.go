package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type Case struct {
	ID            int64     `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Details       string    `db:"details" json:"details"`
	Sender        string    `db:"sender" json:"sender"`
	ETA           time.Time `db:"eta" json:"eta"`
	SLADays       int       `db:"sla_days" json:"sla_days"`
	Hypercare     bool      `db:"hypercare" json:"hypercare"`
	Label         string    `db:"label" json:"label"`
	HSCode        string    `db:"hs_code" json:"hs_code"`
	Preference    string    `db:"preference" json:"preference"`
	SuppUnits     string    `db:"supp_units" json:"supp_units"`
	AssignedTo    string    `db:"assigned_to" json:"assigned_to"`
	PriorityScore int       `db:"priority_score" json:"priority_score"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

func main() {
	dsn := "postgres://postgres:postgres@db:5432/logistics?sslmode=disable"
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	api.POST("/cases", createCase)
	api.GET("/cases", listCases)
	api.POST("/duplicate_check", duplicateCheck)
	api.POST("/hs_suggest", hsSuggest)
	api.POST("/bulk_update", bulkUpdate)
	api.POST("/assign", assignCase)
	
	// Advanced ML Integration Endpoints
	api.POST("/ml/classify", mlClassifyCase)
	api.POST("/ml/extract_metadata", mlExtractMetadata)
	api.POST("/ml/analyze_document", mlAnalyzeDocument)
	api.POST("/ml/forecast_demand", mlForecastDemand)
	api.POST("/ml/optimize_route", mlOptimizeRoute)
	api.POST("/ml/capacity_planning", mlCapacityPlanning)
	api.POST("/ml/similarity_check", mlSimilarityCheck)

	port := "5000"
	r.Run(":" + port)
}

func createCase(c *gin.Context) {
	var in struct {
		Title     string `json:"title"`
		Details   string `json:"details"`
		Sender    string `json:"sender"`
		ETA       string `json:"eta"`
		SLADays   int    `json:"sla_days"`
		Hypercare bool   `json:"hypercare"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	eta := time.Now().Add(48 * time.Hour)
	if in.ETA != "" {
		if t, err := time.Parse(time.RFC3339, in.ETA); err == nil {
			eta = t
		}
	}
	score := computePriorityScore(eta, in.SLADays, in.Hypercare, in.Sender)

	var id int64
	err := db.QueryRow("INSERT INTO cases (title, details, sender, eta, sla_days, hypercare, priority_score, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7, now()) RETURNING id",
		in.Title, in.Details, in.Sender, eta, in.SLADays, in.Hypercare, score).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cs Case
	if err := db.Get(&cs, "SELECT * FROM cases WHERE id=$1", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"case": cs})
}

func listCases(c *gin.Context) {
	cases := []Case{}
	err := db.Select(&cases, "SELECT * FROM cases ORDER BY priority_score DESC, created_at DESC LIMIT 200")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cases)
}

func duplicateCheck(c *gin.Context) {
	var in struct {
		Title  string `json:"title"`
		Sender string `json:"sender"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tokens := tokenize(in.Title)
	rows := []Case{}
	err := db.Select(&rows, "SELECT * FROM cases ORDER BY created_at DESC LIMIT 500")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var candidates []Case
	for _, r := range rows {
		score := tokenOverlapScore(tokens, tokenize(r.Title))
		if strings.EqualFold(strings.TrimSpace(in.Sender), strings.TrimSpace(r.Sender)) {
			score += 2
		}
		if score >= 2 {
			candidates = append(candidates, r)
		}
	}
	c.JSON(200, gin.H{"candidates": candidates})
}

func hsSuggest(c *gin.Context) {
	var in struct {
		Title   string `json:"title"`
		Details string `json:"details"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	keymap := map[string][]string{
		"electronics": {"8542.31", "8517.62"},
		"battery":     {"8507.80"},
		"clothing":    {"6203.42", "6204.69"},
		"toy":         {"9503.00"},
		"food":        {"1602.49"},
	}
	text := strings.ToLower(in.Title + " " + in.Details)
	suggestions := []map[string]interface{}{}
	for k, v := range keymap {
		if strings.Contains(text, k) {
			for _, code := range v {
				suggestions = append(suggestions, map[string]interface{}{"hs": code, "matched": k})
			}
		}
	}
	if len(suggestions) == 0 {
		suggestions = append(suggestions, map[string]interface{}{"hs": "UNKNOWN", "matched": ""})
	}
	c.JSON(200, gin.H{"suggestions": suggestions})
}

func bulkUpdate(c *gin.Context) {
	var in struct {
		IDs        []int64 `json:"ids"`
		HSCode     string  `json:"hs_code"`
		Preference string  `json:"preference"`
		SuppUnits  string  `json:"supp_units"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tx := db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "transaction failed"})
		}
	}()
	for _, id := range in.IDs {
		_, err := tx.Exec("UPDATE cases SET hs_code=$1, preference=$2, supp_units=$3 WHERE id=$4", in.HSCode, in.Preference, in.SuppUnits, id)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "updated": len(in.IDs)})
}

func assignCase(c *gin.Context) {
	var in struct {
		ID       int64  `json:"id"`
		Assignee string `json:"assignee"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("UPDATE cases SET assigned_to=$1 WHERE id=$2", in.Assignee, in.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

// helpers

func computePriorityScore(eta time.Time, slaDays int, hypercare bool, sender string) int {
	score := 0
	hours := time.Until(eta).Hours()
	if hours <= 24 {
		score += 50
	} else if hours <= 72 {
		score += 20
	}
	if slaDays <= 1 {
		score += 30
	} else if slaDays <= 3 {
		score += 10
	}
	if hypercare {
		score += 40
	}
	if strings.Contains(strings.ToLower(sender), "premium") || strings.Contains(strings.ToLower(sender), "vip") {
		score += 30
	}
	return score
}

func tokenize(s string) []string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", " ")
	s = strings.ReplaceAll(s, ".", " ")
	parts := strings.Fields(s)
	var tokens []string
	for _, p := range parts {
		if len(p) > 2 {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

func tokenOverlapScore(a, b []string) int {
	set := map[string]bool{}
	for _, t := range a {
		set[t] = true
	}
	score := 0
	for _, t := range b {
		if set[t] {
			score++
		}
	}
	return score
}

// ML Integration Handler Functions

func mlClassifyCase(c *gin.Context) {
	var input struct {
		Text string `json:"text"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call NLP service
	payload := map[string]interface{}{
		"text": input.Text,
	}
	result, err := callMLService("http://nlp:8000/classify", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func mlExtractMetadata(c *gin.Context) {
	var input struct {
		Text string `json:"text"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call NLP service for metadata extraction
	payload := map[string]interface{}{
		"text": input.Text,
	}
	result, err := callMLService("http://nlp:8000/extract_metadata", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func mlAnalyzeDocument(c *gin.Context) {
	var input struct {
		Image string `json:"image"`
		Type  string `json:"type"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call CV service for document analysis
	payload := map[string]interface{}{
		"image": input.Image,
		"type":  input.Type,
	}
	
	// First extract text
	textResult, err := callMLService("http://cv:8001/extract_text", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Then quality check
	qualityResult, err := callMLService("http://cv:8001/quality_check", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Combine results
	result := map[string]interface{}{
		"text_extraction": textResult,
		"quality_check":   qualityResult,
	}

	c.JSON(http.StatusOK, result)
}

func mlForecastDemand(c *gin.Context) {
	var input struct {
		Data            []map[string]interface{} `json:"data"`
		Days            int                      `json:"days"`
		Model           string                   `json:"model"`
		ConfidenceLevel float64                  `json:"confidence_level"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if input.Days == 0 {
		input.Days = 30
	}
	if input.Model == "" {
		input.Model = "random_forest"
	}
	if input.ConfidenceLevel == 0 {
		input.ConfidenceLevel = 0.95
	}

	// Call time series service
	payload := map[string]interface{}{
		"data":             input.Data,
		"days":             input.Days,
		"model":            input.Model,
		"confidence_level": input.ConfidenceLevel,
	}
	result, err := callMLService("http://ts:8002/forecast", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func mlOptimizeRoute(c *gin.Context) {
	var input struct {
		Locations       []map[string]interface{} `json:"locations"`
		StartLocation   map[string]interface{}   `json:"start_location"`
		VehicleCapacity int                      `json:"vehicle_capacity"`
		Algorithm       string                   `json:"algorithm"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if input.VehicleCapacity == 0 {
		input.VehicleCapacity = 1000
	}
	if input.Algorithm == "" {
		input.Algorithm = "nearest_neighbor"
	}

	// Call optimization service
	payload := map[string]interface{}{
		"locations":        input.Locations,
		"start_location":   input.StartLocation,
		"vehicle_capacity": input.VehicleCapacity,
		"algorithm":        input.Algorithm,
	}
	result, err := callMLService("http://opt:8003/optimize_route", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func mlCapacityPlanning(c *gin.Context) {
	var input struct {
		CurrentCapacity       int                      `json:"current_capacity"`
		DemandData           []map[string]interface{} `json:"demand_data"`
		ExpectedGrowthRate   float64                  `json:"expected_growth_rate"`
		PlanningHorizonMonths int                     `json:"planning_horizon_months"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if input.CurrentCapacity == 0 {
		input.CurrentCapacity = 100
	}
	if input.ExpectedGrowthRate == 0 {
		input.ExpectedGrowthRate = 0.1
	}
	if input.PlanningHorizonMonths == 0 {
		input.PlanningHorizonMonths = 12
	}

	// Call time series service
	payload := map[string]interface{}{
		"current_capacity":         input.CurrentCapacity,
		"demand_data":             input.DemandData,
		"expected_growth_rate":    input.ExpectedGrowthRate,
		"planning_horizon_months": input.PlanningHorizonMonths,
	}
	result, err := callMLService("http://ts:8002/capacity_planning", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func mlSimilarityCheck(c *gin.Context) {
	var input struct {
		Text1 string `json:"text1"`
		Text2 string `json:"text2"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call NLP service for similarity
	payload := map[string]interface{}{
		"text1": input.Text1,
		"text2": input.Text2,
	}
	result, err := callMLService("http://nlp:8000/similarity", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper function to call ML services
func callMLService(url string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service error: %v", result)
	}

	return result, nil
}
