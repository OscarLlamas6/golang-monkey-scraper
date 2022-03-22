package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Lectura               string
	EntryPoint            string
	FileName              string
	TimeOutStatus         bool = false
	TimeOutValue          int64
	MonkeysAmmount        int64
	QueueSize             int64
	NrValue               int64
	Success, Failed, Send int64
)

//LimpiarPantalla fuction
func LimpiarPantalla() {

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	fmt.Println(string(getColor("white")), "--------------------- PRACTICA 2 - G12 SOPES2 ----------------------")
	fmt.Println(string(getColor("purple")), "------------------------ MONKEY WRAPPER CLI ------------------------")
	fmt.Println()
}

// Get Color Name
func getColor(colorName string) string {

	colors := map[string]string{
		"reset":  "\033[0m",
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"blue":   "\033[34m",
		"purple": "\033[35m",
		"cyan":   "\033[36m",
		"white":  "\033[37m",
	}
	return colors[colorName]
}

func main() {

	continuar := true
	LimpiarPantalla()

	for continuar {

		fmt.Println(string(getColor("cyan")), "- Ingrese cantidad de monos buscadores")
		fmt.Print(string(getColor("green")), "SOPES2 ")
		fmt.Print(string(getColor("yellow")), ">> ")
		fmt.Scanln(&Lectura)
		Lectura = strings.TrimSuffix(Lectura, "/")
		Lectura = strings.TrimSpace(Lectura)
		MonkeysAmmount, _ := strconv.ParseInt(Lectura, 10, 64)

		fmt.Println("")
		fmt.Println(string(getColor("cyan")), "- Ingrese tamaño de la cola")
		fmt.Print(string(getColor("green")), "SOPES2 ")
		fmt.Print(string(getColor("yellow")), ">> ")
		fmt.Scanln(&Lectura)
		Lectura = strings.TrimSuffix(Lectura, "/")
		Lectura = strings.TrimSpace(Lectura)
		QueueSize, _ := strconv.ParseInt(Lectura, 10, 64)

		fmt.Println("")
		fmt.Println(string(getColor("cyan")), "- Ingrese valor de Nr")
		fmt.Print(string(getColor("green")), "SOPES2 ")
		fmt.Print(string(getColor("yellow")), ">> ")
		fmt.Scanln(&Lectura)
		Lectura = strings.TrimSuffix(Lectura, "/")
		Lectura = strings.TrimSpace(Lectura)
		NrValue, _ := strconv.ParseInt(Lectura, 10, 64)

		fmt.Println("")
		fmt.Println(string(getColor("cyan")), "- Ingrese búsqueda inicial")
		fmt.Print(string(getColor("green")), "SOPES2 ")
		fmt.Print(string(getColor("yellow")), ">> ")
		fmt.Scanln(&EntryPoint)
		EntryPoint = strings.TrimSuffix(EntryPoint, "/")
		EntryPoint = strings.TrimSpace(EntryPoint)

		fmt.Println("")
		fmt.Println(string(getColor("cyan")), "- Ingrese nombre del archivo donde se almacenará el resultado.")
		fmt.Print(string(getColor("green")), "SOPES2 ")
		fmt.Print(string(getColor("yellow")), ">> ")
		fmt.Scanln(&FileName)
		FileName = strings.TrimSuffix(FileName, "/")
		FileName = strings.TrimSpace(FileName)

		fmt.Println("")
		fmt.Println("")
		fmt.Print(string(getColor("yellow")), "Resumen de configuración:")
		fmt.Print(string(getColor("cyan")), "Cantidad de monos -> ")
		fmt.Println(string(getColor("red")), MonkeysAmmount)
		fmt.Print(string(getColor("cyan")), "Tamaño de la cola -> ")
		fmt.Println(string(getColor("red")), QueueSize)
		fmt.Print(string(getColor("cyan")), "Valor Nr -> ")
		fmt.Println(string(getColor("red")), NrValue)
		fmt.Print(string(getColor("cyan")), "URL inicial -> ")
		fmt.Println(string(getColor("red")), EntryPoint)
		fmt.Print(string(getColor("cyan")), "Nombre del archivo -> ")
		fmt.Println(string(getColor("red")), FileName)

		fmt.Print(string(getColor("cyan")), "Cargando configuración:")

		for i := 0; i < 53; i++ {
			fmt.Print(string(getColor("yellow")), "#")
			time.Sleep(25 * time.Millisecond)
		}

		fmt.Println("")
		fmt.Print(string(getColor("cyan")), "Configuración cargada correctamente.")
		fmt.Println(string(getColor("yellow")), "Presione ENTER para iniciar el scrapping. :D")
		var wait string
		fmt.Scanln(&wait)
		fmt.Println(string(getColor("red")), "Ejecutando Monkey Wrapper... ")
		/*
			AQUI TOCA HACER EL PROCESO, EN ESTE PUNTO
			LAS VARIABLES GLOBALES YA ESTAN SETEADAS PARA PODER USARLAS
		*/
		continuar = false

	}

	fmt.Print(string(getColor("yellow")), "Hasta la próxima! :)")
}
