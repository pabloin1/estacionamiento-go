package models

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Vehiculo interface {
	Estacionar()
}

type vehiculo struct {
	ID                int
	Estacionamiento   Estacionamiento
	ActualizarEstado  func(int, string, []string, []string)
	Ctx               context.Context
	VehiculosEnEspera *[]string
	MutexEspera       *sync.Mutex
}

func NuevoVehiculo(id int, e Estacionamiento, actualizarEstado func(int, string, []string, []string), ctx context.Context, vehiculosEnEspera *[]string, mutexEspera *sync.Mutex) Vehiculo {
	return &vehiculo{
		ID:                id,
		Estacionamiento:   e,
		ActualizarEstado:  actualizarEstado,
		Ctx:               ctx,
		VehiculosEnEspera: vehiculosEnEspera,
		MutexEspera:       mutexEspera,
	}
}

func (v *vehiculo) Estacionar() {
	idStr := fmt.Sprintf("Vehículo %d", v.ID)

	// Añadir a la lista de espera
	v.MutexEspera.Lock()
	*v.VehiculosEnEspera = append(*v.VehiculosEnEspera, idStr)
	v.ActualizarEstado(v.Estacionamiento.GetVehiculosDentro(), fmt.Sprintf("%s llegó y está en espera", idStr), v.Estacionamiento.GetVehiculosEstacionados(), *v.VehiculosEnEspera)
	v.MutexEspera.Unlock()

	// Intentar ocupar un cajón
	cajon, err := v.Estacionamiento.OcuparCajon(v.Ctx)
	if err != nil {
		return
	}

	// Remover de la lista de espera
	v.MutexEspera.Lock()
	index := -1
	for i, vid := range *v.VehiculosEnEspera {
		if vid == idStr {
			index = i
			break
		}
	}
	if index != -1 {
		*v.VehiculosEnEspera = append((*v.VehiculosEnEspera)[:index], (*v.VehiculosEnEspera)[index+1:]...)
	}
	v.MutexEspera.Unlock()

	// Actualizar estado
	v.Estacionamiento.AgregarVehiculoEstacionado(idStr)
	v.ActualizarEstado(v.Estacionamiento.GetVehiculosDentro(), fmt.Sprintf("%s está entrando", idStr), v.Estacionamiento.GetVehiculosEstacionados(), *v.VehiculosEnEspera)

	// Simulación de entrada
	v.Estacionamiento.IntentarEntrada()
	fmt.Printf("%s entrando por la puerta...\n", idStr)
	time.Sleep(1 * time.Second)
	v.Estacionamiento.LiberarEntrada()

	fmt.Printf("%s estacionado en el cajón %d\n", idStr, cajon)
	v.ActualizarEstado(v.Estacionamiento.GetVehiculosDentro(), fmt.Sprintf("%s estacionado en el cajón %d", idStr, cajon), v.Estacionamiento.GetVehiculosEstacionados(), *v.VehiculosEnEspera)

	// Estacionar por un tiempo aleatorio
	tiempoEstacionado := rand.Intn(3) + 3 // Entre 3 y 5 segundos
	select {
	case <-time.After(time.Duration(tiempoEstacionado) * time.Second):
		// Continúa
	case <-v.Ctx.Done():
		return
	}

	// Intentar salir
	v.Estacionamiento.SalirEstacionamiento()
	fmt.Printf("%s saliendo...\n", idStr)
	v.ActualizarEstado(v.Estacionamiento.GetVehiculosDentro(), fmt.Sprintf("%s saliendo", idStr), v.Estacionamiento.GetVehiculosEstacionados(), *v.VehiculosEnEspera)
	time.Sleep(1 * time.Second)
	v.Estacionamiento.LiberarSalida()

	// Liberar el cajón
	v.Estacionamiento.LiberarCajon(cajon)
	v.Estacionamiento.RemoverVehiculoEstacionado(idStr)
	fmt.Printf("%s ha salido del estacionamiento\n", idStr)
	v.ActualizarEstado(v.Estacionamiento.GetVehiculosDentro(), fmt.Sprintf("%s salió del estacionamiento", idStr), v.Estacionamiento.GetVehiculosEstacionados(), *v.VehiculosEnEspera)
}
