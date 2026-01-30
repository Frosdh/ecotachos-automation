package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	url := getenv("HEALTH_URL", "")
	if url == "" {
		fmt.Println("OK")
		writeJSON("OK", 200)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		writeJSON(err.Error(), 0)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("OK: %s\n", resp.Status)
		writeJSON("OK", resp.StatusCode)
	} else {
		fmt.Printf("UNHEALTHY: %s\n", resp.Status)
		writeJSON("UNHEALTHY", resp.StatusCode)
		os.Exit(2)
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

type HealthLog struct {
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Timestamp string `json:"timestamp"`
}

func writeJSON(status string, code int) {
	dir := getenv("HEALTH_LOG_DIR", "logs")
	_ = os.MkdirAll(dir, 0o755)
	f := filepath.Join(dir, "health.json")
	h := HealthLog{Status: status, Code: code, Timestamp: time.Now().UTC().Format(time.RFC3339)}
	b, _ := json.MarshalIndent(h, "", "  ")
	_ = os.WriteFile(f, b, 0o644)
}
