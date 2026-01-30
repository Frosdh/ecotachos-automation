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
	dropletID = 547986490
)

func main() {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("❌ DO_TOKEN no está configurado")
	}

	client := godo.NewFromToken(token)
	ctx := context.Background()

	// Nombre del snapshot con fecha
	snapshotName := fmt.Sprintf("ecotachos-auto-backup-%s", time.Now().Format("2006-01-02"))

	log.Printf("🔄 Iniciando backup automático: %s\n", snapshotName)

	// Crear snapshot
	action, _, err := client.DropletActions.Snapshot(ctx, dropletID, snapshotName)
	if err != nil {
		log.Fatalf("❌ Error creando snapshot: %v", err)
	}

	log.Printf("✅ Snapshot iniciado (ID: %d)\n", action.ID)

	// Esperar a que se complete (máximo 30 minutos)
	timeout := time.After(30 * time.Minute)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("❌ Timeout: El snapshot tardó más de 30 minutos")
		case <-ticker.C:
			action, _, err := client.DropletActions.Get(ctx, dropletID, action.ID)
			if err != nil {
				log.Printf("⚠️  Error verificando estado: %v", err)
				continue
			}

			if action.Status == "completed" {
				log.Println("✅ Backup completado exitosamente")
				
				// Limpiar snapshots antiguos (mantener últimos 5)
				cleanupOldSnapshots(client, ctx, 5)
				return
			} else if action.Status == "errored" {
				log.Fatal("❌ El backup falló")
			}

			log.Printf("⏳ Estado: %s", action.Status)
		}
	}
}

func cleanupOldSnapshots(client *godo.Client, ctx context.Context, keepLast int) {
	snapshots, _, err := client.Droplets.Snapshots(ctx, dropletID, nil)
	if err != nil {
		log.Printf("⚠️  Error listando snapshots: %v", err)
		return
	}

	if len(snapshots) <= keepLast {
		log.Printf("✅ Solo hay %d snapshots, no se eliminará ninguno\n", len(snapshots))
		return
	}

	toDelete := len(snapshots) - keepLast
	log.Printf("🗑️  Eliminando %d snapshots antiguos...\n", toDelete)

	// Eliminar los más antiguos (los últimos en la lista)
	for i := len(snapshots) - 1; i >= len(snapshots)-toDelete; i-- {
		snapshot := snapshots[i]
		log.Printf("   Eliminando: %s (ID: %d)\n", snapshot.Name, snapshot.ID)
		
		_, err := client.Snapshots.Delete(ctx, snapshot.ID)
		if err != nil {
			log.Printf("⚠️  Error eliminando snapshot %d: %v\n", snapshot.ID, err)
		}
	}

	log.Println("✅ Limpieza de snapshots completada")
}