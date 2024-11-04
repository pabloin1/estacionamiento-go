package models

import (
    "context"
    "sync"
)

type Estacionamiento interface {
    IntentarEntrada()
    SalirEstacionamiento()
    LiberarEntrada()
    LiberarSalida()
    OcuparCajon(ctx context.Context) (int, error)
    LiberarCajon(int)
    GetVehiculosDentro() int
    AgregarVehiculoEstacionado(string)
    RemoverVehiculoEstacionado(string)
    GetVehiculosEstacionados() []string
}

type estacionamiento struct {
    Capacidad             int
    Cajones               []bool
    VehiculosDentro       int
    mutexCajones          sync.Mutex
    mutexEntradaSalida    sync.Mutex
    condEntrada           *sync.Cond
    condSalida            *sync.Cond
    condEspacio           *sync.Cond
    entrando              int
    saliendo              int
    vehiculosEstacionados []string
    mutexVehiculos        sync.Mutex
}

func NuevoEstacionamiento(capacidad int) Estacionamiento {
    e := &estacionamiento{
        Capacidad: capacidad,
        Cajones:   make([]bool, capacidad),
    }
    e.condEntrada = sync.NewCond(&e.mutexEntradaSalida)
    e.condSalida = sync.NewCond(&e.mutexEntradaSalida)
    e.condEspacio = sync.NewCond(&e.mutexCajones)
    return e
}

func (e *estacionamiento) IntentarEntrada() {
    e.mutexEntradaSalida.Lock()
    for e.saliendo > 0 {
        e.condEntrada.Wait()
    }
    e.entrando++
    e.mutexEntradaSalida.Unlock()
}

func (e *estacionamiento) SalirEstacionamiento() {
    e.mutexEntradaSalida.Lock()
    for e.entrando > 0 {
        e.condSalida.Wait()
    }
    e.saliendo++
    e.mutexEntradaSalida.Unlock()
}

func (e *estacionamiento) LiberarEntrada() {
    e.mutexEntradaSalida.Lock()
    e.entrando--
    if e.entrando == 0 {
        e.condSalida.Broadcast()
    }
    e.mutexEntradaSalida.Unlock()
}

func (e *estacionamiento) LiberarSalida() {
    e.mutexEntradaSalida.Lock()
    e.saliendo--
    if e.saliendo == 0 {
        e.condEntrada.Broadcast()
    }
    e.mutexEntradaSalida.Unlock()
}

func (e *estacionamiento) OcuparCajon(ctx context.Context) (int, error) {
    e.mutexCajones.Lock()
    defer e.mutexCajones.Unlock()
    for {
        // Buscar un cajón disponible
        for i, ocupado := range e.Cajones {
            if !ocupado {
                e.Cajones[i] = true
                e.VehiculosDentro++
                return i, nil
            }
        }
        // No hay cajones disponibles, esperar
        select {
        case <-ctx.Done():
            return -1, ctx.Err()
        default:
            e.condEspacio.Wait()
        }
    }
}

func (e *estacionamiento) LiberarCajon(cajon int) {
    e.mutexCajones.Lock()
    e.Cajones[cajon] = false
    e.VehiculosDentro--
    e.mutexCajones.Unlock()
    // Notificar a los vehículos que esperan espacio
    e.condEspacio.Broadcast()
}

func (e *estacionamiento) GetVehiculosDentro() int {
    e.mutexCajones.Lock()
    defer e.mutexCajones.Unlock()
    return e.VehiculosDentro
}

func (e *estacionamiento) AgregarVehiculoEstacionado(id string) {
    e.mutexVehiculos.Lock()
    e.vehiculosEstacionados = append(e.vehiculosEstacionados, id)
    e.mutexVehiculos.Unlock()
}

func (e *estacionamiento) RemoverVehiculoEstacionado(id string) {
    e.mutexVehiculos.Lock()
    index := -1
    for i, vid := range e.vehiculosEstacionados {
        if vid == id {
            index = i
            break
        }
    }
    if index != -1 {
        e.vehiculosEstacionados = append(e.vehiculosEstacionados[:index], e.vehiculosEstacionados[index+1:]...)
    }
    e.mutexVehiculos.Unlock()
}

func (e *estacionamiento) GetVehiculosEstacionados() []string {
    e.mutexVehiculos.Lock()
    defer e.mutexVehiculos.Unlock()
    copia := make([]string, len(e.vehiculosEstacionados))
    copy(copia, e.vehiculosEstacionados)
    return copia
}
