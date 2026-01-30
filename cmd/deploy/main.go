package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
)

func main() {
	env := getenv("ENV", "dev")
	app := getenv("APP", "ecotachos")
	version := getenv("VERSION", "v0.1.0")
	host := getenv("HOST", "")
	user := getenv("USER", "root")
	keyPath := getenv("KEY_PATH", "")
	port := getenv("PORT", "22")

	fmt.Printf("Desplegando %s %s al entorno %s...\n", app, version, env)
	if host == "" || keyPath == "" {
		fmt.Println("Modo simulación: faltan HOST o KEY_PATH. Ejecutando pasos locales.")
		simulate()
		return
	}

	if err := deployRemote(host, user, port, keyPath); err != nil {
		log.Fatalf("Error en despliegue remoto: %v", err)
	}
	fmt.Println("Despliegue remoto completado.")
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func simulate() {
	steps := []string{
		"Conectar al servidor",
		"Actualizar repos",
		"Reiniciar servicios",
		"Validar estado",
	}
	for _, s := range steps {
		fmt.Println("-", s)
		time.Sleep(300 * time.Millisecond)
	}
}

func deployRemote(host, user, port, keyPath string) error {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("leyendo llave: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("parseando llave: %w", err)
	}
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 10 * time.Second,
	}
	addr := host + ":" + port
	c, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return fmt.Errorf("conectando: %w", err)
	}
	defer c.Close()

	script := `
set -e
cd /var/www/ecotachostec-backend
if [ -d "/var/www/ecotachostec-backend" ]; then
  sudo cp -r /var/www/ecotachostec-backend /var/www/ecotachostec-backend.backup.$(date +%Y%m%d_%H%M%S)
fi
git fetch origin
git reset --hard origin/main || git reset --hard origin/master

cd /var/www/ecotachostec-backend/docker
sudo docker-compose down
sudo docker-compose build --no-cache
sudo docker-compose up -d
sleep 10
sudo docker-compose ps
curl -f http://localhost:8000/api/ia/health/ || true
curl -f http://localhost/api/ia/health/ || true
`

	if err := runSSH(c, script); err != nil {
		return err
	}
	return nil
}

func runSSH(client *ssh.Client, cmd string) error {
	s, err := client.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()
	s.Stdout = os.Stdout
	s.Stderr = os.Stderr
	return s.Run(cmd)
}
