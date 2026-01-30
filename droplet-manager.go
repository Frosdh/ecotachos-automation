package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
)

const (
	dropletID = 547986490 // Tu droplet ID
)

// DropletManager gestiona operaciones del droplet
type DropletManager struct {
	client *godo.Client
	ctx    context.Context
}

// NewDropletManager crea una nueva instancia del manager
func NewDropletManager(token string) *DropletManager {
	client := godo.NewFromToken(token)
	return &DropletManager{
		client: client,
		ctx:    context.Background(),
	}
}

// GetDropletInfo obtiene información del droplet
func (dm *DropletManager) GetDropletInfo() error {
	droplet, _, err := dm.client.Droplets.Get(dm.ctx, dropletID)
	if err != nil {
		return fmt.Errorf("error obteniendo info del droplet: %v", err)
	}

	fmt.Println("\n=== 📊 Información del Droplet ===")
	fmt.Printf("Nombre: %s\n", droplet.Name)
	fmt.Printf("Estado: %s\n", droplet.Status)
	fmt.Printf("Región: %s\n", droplet.Region.Name)
	fmt.Printf("Memoria: %d MB\n", droplet.Memory)
	fmt.Printf("CPUs: %d\n", droplet.Vcpus)
	fmt.Printf("Disco: %d GB\n", droplet.Disk)
	fmt.Printf("IP Pública: %s\n", droplet.Networks.V4[0].IPAddress)
	fmt.Printf("Creado: %s\n", droplet.Created)
	
	return nil
}

// CreateSnapshot crea un snapshot del droplet
func (dm *DropletManager) CreateSnapshot(name string) error {
	if name == "" {
		name = fmt.Sprintf("ecotachos-backup-%s", time.Now().Format("2006-01-02-15-04"))
	}

	fmt.Printf("\n🔄 Creando snapshot: %s\n", name)
	
	action, _, err := dm.client.DropletActions.Snapshot(dm.ctx, dropletID, name)
	if err != nil {
		return fmt.Errorf("error creando snapshot: %v", err)
	}

	fmt.Printf("✅ Snapshot iniciado (ID: %d)\n", action.ID)
	fmt.Println("⏳ El proceso puede tardar varios minutos...")
	
	// Esperar a que se complete
	return dm.waitForAction(action.ID)
}

// waitForAction espera a que una acción se complete
func (dm *DropletManager) waitForAction(actionID int) error {
	for {
		action, _, err := dm.client.DropletActions.Get(dm.ctx, dropletID, actionID)
		if err != nil {
			return err
		}

		if action.Status == "completed" {
			fmt.Println("✅ Acción completada exitosamente")
			return nil
		} else if action.Status == "errored" {
			return fmt.Errorf("❌ acción falló")
		}

		fmt.Printf("⏳ Estado: %s...\n", action.Status)
		time.Sleep(10 * time.Second)
	}
}

// ListSnapshots lista todos los snapshots del droplet
func (dm *DropletManager) ListSnapshots() error {
	snapshots, _, err := dm.client.Droplets.Snapshots(dm.ctx, dropletID, nil)
	if err != nil {
		return fmt.Errorf("error listando snapshots: %v", err)
	}

	fmt.Println("\n=== 📸 Snapshots Disponibles ===")
	if len(snapshots) == 0 {
		fmt.Println("No hay snapshots disponibles")
		return nil
	}

	for i, snapshot := range snapshots {
		fmt.Printf("\n[%d] %s\n", i+1, snapshot.Name)
		fmt.Printf("    ID: %d\n", snapshot.ID)
		fmt.Printf("    Tamaño: %.2f GB\n", snapshot.SizeGigaBytes)
		fmt.Printf("    Creado: %s\n", snapshot.Created)
		fmt.Printf("    Regiones: %v\n", snapshot.Regions)
	}

	return nil
}

// GetMetrics obtiene métricas del droplet (simulado - DO no expone métricas detalladas via API)
func (dm *DropletManager) GetMetrics() error {
	fmt.Println("\n=== 📈 Métricas del Droplet ===")
	fmt.Println("ℹ️  DigitalOcean no expone métricas detalladas vía API")
	fmt.Println("    Para métricas en tiempo real, usa el dashboard de DO o instala Prometheus")
	
	// Obtener información básica
	droplet, _, err := dm.client.Droplets.Get(dm.ctx, dropletID)
	if err != nil {
		return err
	}

	fmt.Printf("\n✅ Estado del servidor: %s\n", droplet.Status)
	fmt.Printf("📊 Recursos asignados:\n")
	fmt.Printf("   - vCPUs: %d\n", droplet.Vcpus)
	fmt.Printf("   - RAM: %d MB\n", droplet.Memory)
	fmt.Printf("   - Disco: %d GB\n", droplet.Disk)
	
	return nil
}

// DeleteOldSnapshots elimina snapshots antiguos (mantiene los últimos N)
func (dm *DropletManager) DeleteOldSnapshots(keepLast int) error {
	snapshots, _, err := dm.client.Droplets.Snapshots(dm.ctx, dropletID, nil)
	if err != nil {
		return err
	}

	if len(snapshots) <= keepLast {
		fmt.Printf("✅ Solo hay %d snapshots, no se eliminará ninguno\n", len(snapshots))
		return nil
	}

	toDelete := len(snapshots) - keepLast
	fmt.Printf("🗑️  Eliminando %d snapshots antiguos...\n", toDelete)

	for i := 0; i < toDelete; i++ {
		snapshot := snapshots[len(snapshots)-1-i]
		fmt.Printf("   Eliminando: %s (ID: %d)\n", snapshot.Name, snapshot.ID)
		
		_, err := dm.client.Snapshots.Delete(dm.ctx, fmt.Sprintf("%d", snapshot.ID))
		if err != nil {
			log.Printf("⚠️  Error eliminando snapshot %d: %v\n", snapshot.ID, err)
		}
	}

	fmt.Println("✅ Limpieza completada")
	return nil
}

func main() {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("❌ Error: Debes configurar la variable de entorno DO_TOKEN")
	}

	manager := NewDropletManager(token)

	// Mostrar menú
	if len(os.Args) < 2 {
		fmt.Println("\n=== 🚀 EcoTachosTec - Gestor de Droplet ===")
		fmt.Println("\nUso: go run droplet-manager.go [comando]")
		fmt.Println("\nComandos disponibles:")
		fmt.Println("  info       - Mostrar información del droplet")
		fmt.Println("  snapshot   - Crear un nuevo snapshot")
		fmt.Println("  list       - Listar todos los snapshots")
		fmt.Println("  metrics    - Mostrar métricas básicas")
		fmt.Println("  cleanup    - Eliminar snapshots antiguos (mantener últimos 3)")
		fmt.Println("\nEjemplo:")
		fmt.Println("  export DO_TOKEN=tu_token")
		fmt.Println("  go run droplet-manager.go info")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "info":
		if err := manager.GetDropletInfo(); err != nil {
			log.Fatal(err)
		}
	case "snapshot":
		if err := manager.CreateSnapshot(""); err != nil {
			log.Fatal(err)
		}
	case "list":
		if err := manager.ListSnapshots(); err != nil {
			log.Fatal(err)
		}
	case "metrics":
		if err := manager.GetMetrics(); err != nil {
			log.Fatal(err)
		}
	case "cleanup":
		if err := manager.DeleteOldSnapshots(3); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("❌ Comando desconocido: %s", command)
	}
}