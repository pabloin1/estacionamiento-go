package entities

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"math/rand"
	"time"
)

type Vehicle struct {
	id   int
	icon *canvas.Rectangle
}

func NewVehicle(id int) *Vehicle {
	return &Vehicle{
		id:   id,
		icon: canvas.NewRectangle(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}),
	}
}

func (v *Vehicle) Update(status string) {
	fmt.Printf("Vehicle %d recibió notificación: %s\n", v.id, status)
}

func (v *Vehicle) GetID() int {
	return v.id
}

func (v *Vehicle) Park(lot *ParkingLot, vehicleContainer *fyne.Container) {
	v.Enter(lot, vehicleContainer)

	// Movimiento hacia la zona de estacionamiento
	for y := 0; y < 210; y += 20 {
		v.icon.Move(fyne.NewPos(130, float32(140+y)))
		vehicleContainer.Refresh()
		time.Sleep(100 * time.Millisecond)
	}

	// Permanecer estacionado 3 segundos
	time.Sleep(3 * time.Second)

	// Luego, salir del estacionamiento
	v.Exit(lot, vehicleContainer)
}


func (v *Vehicle) Enter(lot *ParkingLot, vehicleContainer *fyne.Container) {
	// Espera hasta que haya espacio disponible en el estacionamiento
	lot.spots <- v.id // Espera en el canal de spots cuando el estacionamiento esté lleno
	fmt.Printf("Vehicle %d ha ingresado\n", v.id)
	lot.NotifyObservers(v.id, fmt.Sprintf("Vehicle %d ha ingresado", v.id))

	// Animación de entrada al estacionamiento
	for x := 0; x < 210; x += 20 {
		v.icon.Move(fyne.NewPos(float32(20+x), 80))
		vehicleContainer.Refresh()
		time.Sleep(100 * time.Millisecond)
	}
}


func (v *Vehicle) Exit(lot *ParkingLot, vehicleContainer *fyne.Container) {
	<-lot.spots
	lot.entryGate <- struct{}{}
	fmt.Printf("Vehicle %d ha salido\n", v.id)
	lot.NotifyObservers(v.id, fmt.Sprintf("Vehicle %d ha salido", v.id))

	for x := 0; x < 310; x += 20 {
		v.icon.Move(fyne.NewPos(float32(420-x), 650))
		vehicleContainer.Refresh()
		time.Sleep(100 * time.Millisecond)
	}
	vehicleContainer.Remove(v.icon)
	vehicleContainer.Refresh()
	<-lot.entryGate
}

func (v *Vehicle) CreateIcon() *canvas.Rectangle {
	return v.icon
}
