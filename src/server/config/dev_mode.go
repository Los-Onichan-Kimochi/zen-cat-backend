package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var IsDevMode bool = false

func SetDevMode(devMode bool) {
	IsDevMode = devMode
}

func GetDevMode() bool {
	return IsDevMode
}

func InitDevMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("¿Ejecutar en modo desarrollo? (omite autenticación JWT) [y/N]: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo entrada, ejecutando en modo producción...")
		SetDevMode(false)
		return
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		SetDevMode(true)
		fmt.Println("🔓 Modo dev activado - Autenticación JWT deshabilitada")
	} else {
		SetDevMode(false)
		fmt.Println("🔒 Modo producción activado - Autenticación JWT habilitada")
	}
}
