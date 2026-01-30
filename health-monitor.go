package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/digitalocean/godo"
)

const (
	dropletID = 547986490
)

type HealthStatus struct {
	Timestamp   string `json:"timestamp"`
	DropletName string `json:"droplet_name"`
	Status      string `json:"status"`
	IPAddress   string `json:"ip_address"`
	Memory      int    `json:"memory_mb"`
	CPUs        int    `json:"cpus"`
	Disk        int    `json:"disk_gb"`
	Region      string `json:"region"`
	Uptime      string `json:"uptime"`
	Healthy     bool   `json:"healthy"`
}

type ServiceHealth struct {
	Service string `json:"service"`
	URL     string `json:"url"`
	Status  int    `json:"status"`
	Healthy bool   `json:"healthy"`
}

func main() {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("❌ DO_TOKEN no está configurado")
	}

	client := godo.NewFromToken(token)
	ctx := context.Background()

	// Obtener información del droplet
	droplet, _, err := client.Droplets.Get(ctx, dropletID)
	if err != nil {
		log.Fatalf("❌ Error obteniendo info del droplet: %v", err)
	}

	// Calcular uptime aproximado
	created, _ := time.Parse(time.RFC3339, droplet.Created)
	uptime := time.Since(created)

	// Estado del droplet
	health := HealthStatus{
		Timestamp:   time.Now().Format(time.RFC3339),
		DropletName: droplet.Name,
		Status:      droplet.Status,
		IPAddress:   getPublicIP(droplet),
		Memory:      droplet.Memory,
		CPUs:        droplet.Vcpus,
		Disk:        droplet.Disk,
		Region:      droplet.Region.Name,
		Uptime:      formatDuration(uptime),
		Healthy:     droplet.Status == "active",
	}

	// Verificar servicios
	services := []ServiceHealth{
		checkService("Frontend", fmt.Sprintf("http://%s", health.IPAddress)),
		checkService("Backend API", fmt.Sprintf("http://%s/api/ia/health/", health.IPAddress)),
	}

	// Imprimir reporte
	fmt.Println("\n╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║       🚀 ECOTACHOSTEC - HEALTH CHECK REPORT          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	
	fmt.Printf("\n📅 Timestamp: %s\n", health.Timestamp)
	fmt.Println("\n=== 💻 Droplet Status ===")
	fmt.Printf("Nombre:     %s\n", health.DropletName)
	fmt.Printf("Estado:     %s %s\n", getStatusEmoji(health.Healthy), health.Status)
	fmt.Printf("IP Pública: %s\n", health.IPAddress)
	fmt.Printf("Región:     %s\n", health.Region)
	fmt.Printf("Uptime:     %s\n", health.Uptime)
	
	fmt.Println("\n=== 📊 Recursos ===")
	fmt.Printf("CPUs:   %d vCPUs\n", health.CPUs)
	fmt.Printf("RAM:    %d MB\n", health.Memory)
	fmt.Printf("Disco:  %d GB\n", health.Disk)
	
	fmt.Println("\n=== 🌐 Servicios ===")
	allHealthy := health.Healthy
	for _, svc := range services {
		fmt.Printf("%s %s: %s (HTTP %d)\n", 
			getStatusEmoji(svc.Healthy), 
			svc.Service, 
			getHealthText(svc.Healthy),
			svc.Status)
		if !svc.Healthy {
			allHealthy = false
		}
	}
	
	fmt.Println("\n=== 🎯 Resumen ===")
	if allHealthy {
		fmt.Println("✅ Todos los sistemas operando normalmente")
	} else {
		fmt.Println("⚠️  Algunos servicios presentan problemas")
	}
	
	// Guardar reporte en JSON
	saveReport(health, services)
	
	// Exit code basado en salud
	if !allHealthy {
		os.Exit(1)
	}
}

func getPublicIP(droplet *godo.Droplet) string {
	if len(droplet.Networks.V4) > 0 {
		for _, network := range droplet.Networks.V4 {
			if network.Type == "public" {
				return network.IPAddress
			}
		}
		return droplet.Networks.V4[0].IPAddress
	}
	return "N/A"
}

func checkService(name, url string) ServiceHealth {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return ServiceHealth{
			Service: name,
			URL:     url,
			Status:  0,
			Healthy: false,
		}
	}
	defer resp.Body.Close()
	
	return ServiceHealth{
		Service: name,
		URL:     url,
		Status:  resp.StatusCode,
		Healthy: resp.StatusCode >= 200 && resp.StatusCode < 400,
	}
}

func getStatusEmoji(healthy bool) string {
	if healthy {
		return "✅"
	}
	return "❌"
}

func getHealthText(healthy bool) string {
	if healthy {
		return "HEALTHY"
	}
	return "UNHEALTHY"
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func saveReport(health HealthStatus, services []ServiceHealth) {
	report := map[string]interface{}{
		"droplet":  health,
		"services": services,
	}
	
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Printf("⚠️  Error generando JSON: %v", err)
		return
	}
	
	filename := fmt.Sprintf("health-report-%s.json", time.Now().Format("2006-01-02-15-04"))
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("⚠️  Error guardando reporte: %v", err)
		return
	}
	
	fmt.Printf("\n📄 Reporte guardado: %s\n", filename)
}