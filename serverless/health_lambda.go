package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// HealthResponse representa el payload de salud
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Handle es un handler simple para ambientes serverless
// En AWS, se usaría: lambda.Start(Handle)
// En GCP (Cloud Functions), se expondría como función HTTP
func Handle(ctx context.Context) (string, error) {
	h := HealthResponse{Status: "ok", Timestamp: time.Now().UTC().Format(time.RFC3339)}
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// main permite pruebas locales sin dependencias externas
func main() {
	s, err := Handle(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
