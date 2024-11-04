// gui.go
package interfaces

import (
    "context"
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
    "main/aplicacion"
    "math/rand"
    "time"
)

var (
    cancelFunc context.CancelFunc
)

func IniciarGUI() {
    a := app.New()

    ventana := a.NewWindow("Simulador de Estacionamiento")
    rand.Seed(time.Now().UnixNano())

    // Variables de binding
    estadoBinding := binding.NewString()
    estadoBinding.Set("Vehículos en el estacionamiento: 0")
    estadoLabel := widget.NewLabelWithData(estadoBinding)

    contadorBinding := binding.NewString()
    contadorBinding.Set("Contador de Vehículos: 0")
    contador := widget.NewLabelWithData(contadorBinding)

    eventLogBinding := binding.NewStringList()
    logList := widget.NewListWithData(
        eventLogBinding,
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(i binding.DataItem, o fyne.CanvasObject) {
            str, _ := i.(binding.String).Get()
            o.(*widget.Label).SetText(str)
        },
    )
    logList.Resize(fyne.NewSize(800, 200)) // Aumentar el largo del registro de eventos

    // Nuevas variables de binding para las listas de vehículos
    vehiculosEstacionadosBinding := binding.NewStringList()
    vehiculosEnEsperaBinding := binding.NewStringList()

    listaEstacionados := widget.NewListWithData(
        vehiculosEstacionadosBinding,
        func() fyne.CanvasObject {
            carImage := canvas.NewImageFromFile("interfaces/recursos/vehiculo_estacionado.jpg")
            carImage.SetMinSize(fyne.NewSize(50, 50))
            carImage.FillMode = canvas.ImageFillContain
            label := widget.NewLabel("")
            label.Alignment = fyne.TextAlignCenter
            return container.NewHBox(
                container.NewCenter(carImage),
                label,
            )
        },
        func(i binding.DataItem, o fyne.CanvasObject) {
            container := o.(*fyne.Container)
            str, _ := i.(binding.String).Get()
            label := container.Objects[1].(*widget.Label)
            label.SetText(str)
        },
    )
    listaEstacionados.Resize(fyne.NewSize(800, 600)) // Aumentar el largo de la lista de vehículos estacionados

    listaEnEspera := widget.NewListWithData(
        vehiculosEnEsperaBinding,
        func() fyne.CanvasObject {
            carImage := canvas.NewImageFromFile("interfaces/recursos/vehiculo_espera.png")
            carImage.SetMinSize(fyne.NewSize(50, 50))
            carImage.FillMode = canvas.ImageFillContain
            label := widget.NewLabel("")
            label.Alignment = fyne.TextAlignCenter
            return container.NewHBox(
                container.NewCenter(carImage),
                label,
            )
        },
        func(i binding.DataItem, o fyne.CanvasObject) {
            container := o.(*fyne.Container)
            str, _ := i.(binding.String).Get()
            label := container.Objects[1].(*widget.Label)
            label.SetText(str)
        },
    )
    listaEnEspera.Resize(fyne.NewSize(800, 600)) // Aumentar el largo de la lista de vehículos en espera

    capacidadEntry := widget.NewEntry()
    capacidadEntry.SetPlaceHolder("Capacidad del estacionamiento")

    var botonIniciar *widget.Button
    botonIniciar = widget.NewButton("Iniciar Simulación", func() {
        capacidad := 20
        fmt.Sscanf(capacidadEntry.Text, "%d", &capacidad)
        eventLogBinding.Set([]string{})
        vehiculosEstacionadosBinding.Set([]string{})
        vehiculosEnEsperaBinding.Set([]string{})
        ctx, cancel := context.WithCancel(context.Background())
        cancelFunc = cancel
        botonIniciar.Disable()

        go aplicacion.IniciarSimulacion(ctx, capacidad, func(vehiculosDentro int, evento string, estacionados []string, enEspera []string) {
            fyne.CurrentApp().SendNotification(&fyne.Notification{
                Title:   "Simulación Actualizada",
                Content: evento,
            })

            estadoBinding.Set(fmt.Sprintf("Vehículos en el estacionamiento: %d", vehiculosDentro))
            contadorBinding.Set(fmt.Sprintf("Contador de Vehículos: %d", vehiculosDentro))
            eventLogBinding.Append(evento)
            vehiculosEstacionadosBinding.Set(estacionados)
            vehiculosEnEsperaBinding.Set(enEspera)
        })
    })

    botonDetener := widget.NewButton("Detener Simulación", func() {
        if cancelFunc != nil {
            cancelFunc()
            cancelFunc = nil
            botonIniciar.Enable()
        }
    })

    // Configurar el contenido principal de la ventana con más espaciado y organización
    ventana.SetContent(container.NewVBox(
        widget.NewLabelWithStyle("Simulador de Estacionamiento", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true}),
        estadoLabel,
        contador,
        widget.NewLabel("Capacidad:"),
        capacidadEntry,
        container.NewHBox(botonIniciar, botonDetener),
        widget.NewLabel("Vehículos Estacionados:"),
        container.NewBorder(nil, nil, nil, nil, listaEstacionados), // Lista de vehículos estacionados expandida
        widget.NewLabel("Vehículos en Espera:"),
        container.NewBorder(nil, nil, nil, nil, listaEnEspera), // Lista de vehículos en espera expandida
        widget.NewLabel("Registro de Eventos:"),
        container.NewBorder(nil, nil, nil, nil, logList), // Registro de eventos expandido
    ))

    ventana.Resize(fyne.NewSize(900, 1200)) // Aumentar el tamaño vertical para más espacio
    ventana.ShowAndRun()
}
