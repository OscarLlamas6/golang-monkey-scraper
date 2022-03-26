package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Work struct {
	SHAPadre string
	URL      string
	NR       int64
}

type Mono struct {
	ID         string
	Disponible bool
}

type Resultado struct {
	Origen   string
	Palabras int64
	Enlaces  int64
	SHA      string
	URL      string
	Mono     string
}

var (
	Lectura            string
	EntryPoint         string
	FileName           string
	TimeOutStatus      bool = false
	TimeOutValue       int64
	MonkeysAmmount     int64
	QueueSize          int64
	NrValue            int64
	contadorEscrituras int64 = 0
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

// Metodo para realizar el scrapping
func RunScrapper(nr int64, url string, mono string, origen string) {

	contenido := ""
	var cantidadEnlaces int64 = 0
	var contadorNR int64 = 0
	file, err := os.OpenFile(FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("could not create the file, err :%q", err)
		panic(err)
	}

	defer file.Close()

	c := colly.NewCollector()

	// Buscando todas los enlaces dentro de las etiquetas p
	c.OnHTML("p", func(e *colly.HTMLElement) {
		contenido += e.Text + "\n"
		e.ForEach("a[href]", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			if elem.Request.AbsoluteURL(link) != "" {

				cantidadEnlaces++

				if nr > 0 && contadorNR < nr {

					/* AQUI SE MANDARIAN LOS NUEVOS WORKS A LA COLA*/

					// variable con el enlace para el nuevo work que se ingresara en la cola
					linkNuevoWork := elem.Request.AbsoluteURL(link)

					// variable con el valor de nr - 1
					newNR := nr - 1

					contadorNR++

					fmt.Printf("Se encontro el enlace #%v: %s - nuevo Nr = %v \n ", contadorNR, linkNuevoWork, newNR)

				}
			}
		})
	})

	c.Visit(url)

	h := sha1.New()
	h.Write([]byte(contenido))
	contentSHA := hex.EncodeToString(h.Sum(nil))

	JsonResult := Resultado{
		Origen:   origen,
		Palabras: int64(WordCount(contenido)),
		Enlaces:  cantidadEnlaces,
		SHA:      contentSHA,
		URL:      url,
		Mono:     mono,
	}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(JsonResult)
	if err != nil {
		panic(err)
	}

	if _, err = file.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

	contadorEscrituras++

	if _, err = file.WriteString(",\n"); err != nil {
		panic(err)
	}
}

func WordCount(value string) int {
	re := regexp.MustCompile(`[\S]+`)
	results := re.FindAllString(value, -1)
	return len(results) + 1
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
		fmt.Println(string(getColor("yellow")), "Resumen de configuración:")
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

		// Realizando un scrapper de prueba, esto se haria iterando en la cola
		RunScrapper(NrValue, EntryPoint, "mono_01", "0")

		continuar = false

	}

	fmt.Print(string(getColor("yellow")), "Hasta la próxima! :)")
}
