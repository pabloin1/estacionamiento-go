package main

import (
	"fmt"
	"main/ui"
)

func main() {
	fmt.Println("Iniciando Parking Simulator...")

	appUI := ui.NewAppUI()
	fmt.Println("Interfaz principal iniciada.")

	appUI.Start()
	fmt.Println("Aplicación en ejecución.")
}
