package simulation

import (
	"main/entities"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"sync"
	"time"
)

type SimView struct {
	window fyne.Window
}

func NewSimView(window fyne.Window) *SimView {
	return &SimView{
		window: window,
	}
}

var vehicleContainer = container.NewWithoutLayout()

func (view *SimView) Render() {
	bg := canvas.NewRectangle(color.RGBA{220, 220, 220, 255})
	bg.Resize(fyne.NewSize(620, 720))
	vehicleContainer.Add(bg)

	entryZone := canvas.NewRectangle(color.RGBA{150, 200, 255, 255})
	entryZone.Resize(fyne.NewSize(170, 60))
	entryZone.Move(fyne.NewPos(10, 40))
	vehicleContainer.Add(entryZone)

	parkingZone := canvas.NewRectangle(color.RGBA{210, 210, 210, 255})
	parkingZone.Resize(fyne.NewSize(440, 500))
	parkingZone.Move(fyne.NewPos(10, 110))
	vehicleContainer.Add(parkingZone)

	exitZone := canvas.NewRectangle(color.RGBA{255, 160, 160, 255})
	exitZone.Resize(fyne.NewSize(170, 60))
	exitZone.Move(fyne.NewPos(440, 40))
	vehicleContainer.Add(exitZone)

	view.window.SetContent(vehicleContainer)
}

func (view *SimView) Execute() {
	// El estacionamiento tiene un máximo de 20 espacios disponibles
	lot := entities.NewParkingLot(make(chan int, 20), make(chan struct{}, 1), &sync.Mutex{})
	var waitGroup sync.WaitGroup
	arrivalDist := entities.NewExponentialDist()

	for i := 0; i < 100; i++ { // Generamos 100 vehículos
		waitGroup.Add(1)
		vehicle := entities.NewVehicle(i)
		lot.Register(vehicle)

		go view.manageVehicle(vehicle, lot, &waitGroup)

		// Intervalo aleatorio para la llegada de cada vehículo
		interval := time.Duration(arrivalDist.Generate(2) * float64(time.Second)) // Rate 2 para mayor frecuencia
		time.Sleep(interval)
	}

	waitGroup.Wait()
	fmt.Println("Fin de la simulación.")
}


func (view *SimView) manageVehicle(vehicle *entities.Vehicle, lot *entities.ParkingLot, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	vehicleIcon := vehicle.CreateIcon()
	vehicleIcon.Resize(fyne.NewSize(45, 30))
	vehicleIcon.Move(fyne.NewPos(-60, 80))

	vehicleContainer.Add(vehicleIcon)
	vehicleContainer.Refresh()

	vehicle.Park(lot, vehicleContainer)
}
