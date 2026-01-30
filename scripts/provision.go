package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

type Result struct {
	Provider string `json:"provider"`
	Action   string `json:"action"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
	Time     string `json:"time"`
}

func main() {
	provider := flag.String("provider", "aws", "Proveedor cloud: aws|gcp|do")
	action := flag.String("action", "list", "Acción: create-vm|list|status")
	name := flag.String("name", "eco-vm", "Nombre del recurso")
	region := flag.String("region", "us-east-1", "Región")
	flag.Parse()

	res := execute(*provider, *action, *name, *region)
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

func execute(provider, action, name, region string) Result {
	// Nota: Implementación simulada. Integra SDKs:
	// - AWS: github.com/aws/aws-sdk-go-v2
	// - GCP: cloud.google.com/go/compute/apiv1
	// - DO: github.com/digitalocean/godo
	// para producción.
	res := Result{Provider: provider, Action: action, Status: "ok", Time: time.Now().UTC().Format(time.RFC3339)}
	switch action {
	case "create-vm":
		res.Detail = fmt.Sprintf("VM '%s' solicitada en %s (%s)", name, provider, region)
	case "list":
		res.Detail = fmt.Sprintf("Listado de instancias en %s (%s)", provider, region)
	case "status":
		res.Detail = fmt.Sprintf("Estado de '%s' en %s (%s): running", name, provider, region)
	default:
		res.Status = "error"
		res.Detail = "Acción no soportada"
	}
	return res
}
