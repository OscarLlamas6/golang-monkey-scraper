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
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Work struct {
	SHAPadre string
	URL      string
	NR       int64
}

type Monkey struct {
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

type Cache struct {
	Mu    sync.Mutex
	Slots []Work
}

var (
	SeguirScrapper     bool = true
	Lectura            string
	EntryPoint         string
	FileName           string
	MonkeysAmmount     int64
	QueueSize          int64
	NrValue            int64
	contadorEscrituras int64 = 0
	continuar          bool  = true
	myQueue            *Cache
	myResults          []Resultado
	myMonkeys          []Monkey
)

// checkear constantemente si hay trabajos nuevos para agregar a la cola
func CheckJobs(canal *chan Work, cache *Cache) {
	for SeguirScrapper {

		myJob := <-*canal
		// fmt.Println("Se recibio un nuevo work en el canal")
		// fmt.Println(myJob)
		Queue(&myJob)
	}
}

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
	fmt.Println(string(getColor("purple")), "------------------------ MONKEY SCRAPER CLI ------------------------")
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

// Escribir resultados en el json
func WriteResult() {

	file, err := os.OpenFile(FileName, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("could not create the file, err :%q", err)
		panic(err)
	}

	file.Truncate(0)
	file.Seek(0, 0)

	if _, err = file.WriteString("["); err != nil {
		panic(err)
	}

	for i, result := range myResults {

		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetIndent("", "\t")

		err = encoder.Encode(result)
		if err != nil {
			panic(err)
		}

		if _, err = file.Write(buffer.Bytes()); err != nil {
			panic(err)
		}

		if i < (len(myResults) - 1) {
			if _, err = file.WriteString(",\n"); err != nil {
				panic(err)
			}
		}

	}

	if _, err = file.WriteString("]"); err != nil {
		panic(err)
	}

	file.Close()
}

// Metodo para realizar el scrapping
/*
nr = numero de enlaces que debe buscar en la pagina
url = direccion de la pagina
mono = id del mono que esta haciendo el scrapping
origen = es el SHA del contenido de la pagina donde se encontro ese enlace, si es el primer link origen='0'
*/
func RunScraper(nr int64, url string, mono string, origen string, works chan Work, monkeyIndex int) {

	contenido := ""
	var cantidadEnlaces int64 = 0
	var contadorNR int64 = 0
	var contentSHA string = ""
	c := colly.NewCollector()

	c.OnHTML("p", func(e *colly.HTMLElement) {
		contenido += e.Text + "\n"
	})

	// Buscando todas los enlaces dentro de las etiquetas p
	c.OnHTML("p", func(e *colly.HTMLElement) {

		e.ForEach("p", func(_ int, text *colly.HTMLElement) {
			contenido += text.Text + "\n"
		})

		h := sha1.New()
		h.Write([]byte(contenido))
		// variable con el SHA para el atributo SHAPadre de los nuevos works
		contentSHA = hex.EncodeToString(h.Sum(nil))

		e.ForEach("a[href]", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			if elem.Request.AbsoluteURL(link) != "" {

				cantidadEnlaces++

				if nr > 0 && contadorNR < nr {

					// variable con el enlace para el nuevo work que se ingresara en la cola
					linkNuevoWork := elem.Request.AbsoluteURL(link)

					// variable con el valor de nr - 1
					newNR := nr - 1

					contadorNR++

					/* AQUI SE MANDARIAN LOS NUEVOS WORKS A LA COLA*/
					// fmt.Printf("Se encontro el enlace #%v: %s - nuevo Nr = %v \n ", contadorNR, linkNuevoWork, newNR)

					foundJob := Work{
						SHAPadre: contentSHA,
						URL:      linkNuevoWork,
						NR:       newNR,
					}

					works <- foundJob

				}
			}
		})
	})

	c.Visit(url)

	JsonResult := Resultado{
		Origen:   origen,
		Palabras: int64(WordCount(contenido)),
		Enlaces:  cantidadEnlaces,
		SHA:      contentSHA,
		URL:      url,
		Mono:     mono,
	}

	myResults = append(myResults, JsonResult)
	contadorEscrituras++
	numberOfMiliseconds := WordCount(contenido)

	mensaje := fmt.Sprintf("El mono %s esta descanzando", myMonkeys[monkeyIndex].ID)
	fmt.Println(string(getColor("yellow")), mensaje)
	time.Sleep(time.Millisecond * time.Duration(numberOfMiliseconds) * 2)
	myMonkeys[monkeyIndex].Disponible = true
	mensaje = fmt.Sprintf("El mono %s ya esta disponible", myMonkeys[monkeyIndex].ID)
	fmt.Println(string(getColor("cyan")), mensaje)
}

// Metodo para contar palabras
func WordCount(value string) int {
	re := regexp.MustCompile(`[\S]+`)
	results := re.FindAllString(value, -1)
	return len(results) + 1
}

// Agregar un trabajo a la cola
func Queue(work *Work) {
	if len(myQueue.Slots) < cap(myQueue.Slots) {
		myQueue.Mu.Lock()
		myQueue.Slots = append(myQueue.Slots, *work)
		myQueue.Mu.Unlock()
	}
}

// Quitar un trabajo de la cola
func DeQueue() *Work {
	myQueue.Mu.Lock()
	auxWork := myQueue.Slots[0] // x|2|3|4|5|...
	myQueue.Slots = myQueue.Slots[1:]
	myQueue.Mu.Unlock()

	return &auxWork
}

// Buscar un mono disponible
func FindAvailableMoney(monos *[]Monkey) int {

	for i, mono := range *monos {

		if mono.Disponible {
			return i
		}
	}

	return -1
}

// Setear parametros de scrapper
func SetScraper(firstJob *Work, monosValue int64, queueSize int64) {

	// haciendo un canal donde pasaran Works
	jobs := make(chan Work, 1000)

	// Haciendo un arreglo donde estan mis N monos
	myMonkeys = make([]Monkey, monosValue)

	// Seteando ID para cada mono
	for i := 0; i < int(monosValue); i++ {
		myMonkeys[i].ID = "monkey_0" + strconv.Itoa(i+1) //monkey_01, monkey_02
		myMonkeys[i].Disponible = true
	}

	// Creando cola de trabajos
	myQueue = &Cache{Slots: make([]Work, 0, queueSize)}
	myQueue.Mu.Lock()
	// agregando el trabajo inicial a la cola
	myQueue.Slots = append(myQueue.Slots, *firstJob)
	myQueue.Mu.Unlock()

	// Go routine para leer canal y agregar trabajos a la cola
	go CheckJobs(&jobs, myQueue)

	for len(myQueue.Slots) > 0 {

		if len(myQueue.Slots) > 0 {

			newJob := DeQueue()
			searchMonkey := true
			monkeyIndex := -1

			for searchMonkey {

				monkeyIndex = FindAvailableMoney(&myMonkeys)
				if monkeyIndex != -1 {
					searchMonkey = false
				}

			}

			// si llego aqui, es porque hay un mono disponible
			myMonkeys[monkeyIndex].Disponible = false
			monkeyID := myMonkeys[monkeyIndex].ID
			go RunScraper(newJob.NR, newJob.URL, monkeyID, newJob.SHAPadre, jobs, monkeyIndex)
			time.Sleep(time.Millisecond * 1000)
		}

	}

	WriteResult()

	keepScraping := true

	for keepScraping {

		contadorMonosDisponibles := 0

		for _, mono := range myMonkeys {

			if mono.Disponible {
				contadorMonosDisponibles++
			}

		}

		if contadorMonosDisponibles == len(myMonkeys) {
			keepScraping = false
		}

	}

	time.Sleep(time.Second * 1)
}

func main() {

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
		fmt.Println(string(getColor("yellow")), "Presione ENTER para iniciar el scraping. :D")
		var wait string
		fmt.Scanln(&wait)
		fmt.Println(string(getColor("red")), "Ejecutando Monkey Scraper... ")

		myFirstJob := Work{
			SHAPadre: "0",
			URL:      EntryPoint,
			NR:       NrValue,
		}

		SetScraper(&myFirstJob, MonkeysAmmount, QueueSize)

		// Realizando un scrapper de prueba, esto se haria iterando en la cola
		//RunScrapper(NrValue, EntryPoint, "mono_01", "0")
		fmt.Println(string(getColor("green")), "Scraping terminado. Presione ENTER para salir. :D")
		fmt.Scanln(&Lectura)
		continuar = false

	}

	fmt.Print(string(getColor("yellow")), "Hasta la próxima! :)")
}
