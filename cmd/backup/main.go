package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Directorio de snapshots local
	outDir := getenv("SNAPSHOT_DIR", "snapshots")
	_ = os.MkdirAll(outDir, 0o755)

	name := fmt.Sprintf("snapshot-%s.txt", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(outDir, name)

	content := fmt.Sprintf("Backup creado: %s (UTC)\n", time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Printf("ERROR creando snapshot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Snapshot generado: %s\n", path)
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
