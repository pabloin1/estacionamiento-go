package aplicacion

import (
    "context"
    "main/models"
    "math/rand"
    "sync"
    "time"
)

var (
    estacionamiento   models.Estacionamiento
    mutexID           sync.Mutex
    vehiculoID        int
    vehiculosEnEspera []string
    mutexEspera       sync.Mutex
)

func IniciarSimulacion(ctx context.Context, capacidad int, actualizarEstado func(int, string, []string, []string)) {
    estacionamiento = models.NuevoEstacionamiento(capacidad)
    vehiculoID = 0
    vehiculosEnEspera = []string{}

    var wg sync.WaitGroup

    for i := 0; i < 100; i++ { // Simular 100 vehículos
        select {
        case <-ctx.Done():
            return
        default:
            tiempoLlegada := rand.ExpFloat64() / 2.0
            time.Sleep(time.Duration(tiempoLlegada) * time.Second)

            mutexID.Lock()
            vehiculoID++
            id := vehiculoID
            mutexID.Unlock()

            wg.Add(1)
            v := models.NuevoVehiculo(id, estacionamiento, actualizarEstado, ctx, &vehiculosEnEspera, &mutexEspera)
            go func() {
                defer wg.Done()
                v.Estacionar()
            }()
        }
    }

    // Esperar a que todos los vehículos terminen
    wg.Wait()
}
