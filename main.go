package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	definitions := parseDefinitions()
	boundaries := parseBoundaries(definitions)

	for _, boundary := range boundaries {

		reversedBoundary := rewindBoundary(&boundary)

		boundaryJson, err := json.MarshalIndent(&reversedBoundary, "", "	")

		if err != nil {
			fmt.Println("Error marshalling json", err)
			continue
		}

		os.WriteFile("./output/"+boundary.Properties.Prefix+".json", []byte(boundaryJson), 0644)
	}
}

type Definition struct {
	Icao   string `json:"icao"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func parseDefinitions() []Definition {
	file, err := os.Open("./data/definitions.txt")

	if err != nil {
		log.Fatal("Error opening definitions", err)
	}

	scanner := bufio.NewScanner(file)

	var output []Definition

	for scanner.Scan() {
		// ICAO|NAME|PREFIX POSITION|
		elements := strings.Split(scanner.Text(), "|")

		definition := Definition{Icao: elements[0], Name: elements[1], Prefix: elements[2]}
		output = append(output, definition)
	}

	return output
}

type Geometry struct {
	Kind        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type Boundary struct {
	Kind       string     `json:"type"`
	Properties Definition `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

func parseBoundaries(definitions []Definition) []Boundary {
	file, err := os.Open("./data/coordinates.txt")

	if err != nil {
		log.Fatal("Error opening boundaries", err)
	}

	scanner := bufio.NewScanner(file)

	var output []Boundary

	for scanner.Scan() {
		elements := strings.Split(scanner.Text(), "|")

		isCoordinate := len(elements) == 2

		if isCoordinate {
			lat, err := strconv.ParseFloat(elements[0], 64)

			if err != nil {
				fmt.Println("Error parsing coords", err)
				continue
			}

			lon, err := strconv.ParseFloat(elements[1], 64)

			if err != nil {
				fmt.Println("Error parsing coords", err)
				continue
			}

			output[len(output)-1].Geometry.Coordinates[0] = append(output[len(output)-1].Geometry.Coordinates[0], []float64{lon, lat})
		} else {
			// add first coordinate to end of previous coordinate array (circular geojson)
			if len(output) != 0 {
				previous := output[len(output)-1].Geometry.Coordinates[0]
				output[len(output)-1].Geometry.Coordinates[0] = append(previous, previous[0])
			}

			var boundaryDefinition Definition

			for _, definition := range definitions {
				if definition.Icao == elements[0] {
					boundaryDefinition = definition
				}
			}

			if len(boundaryDefinition.Icao) == 0 {
				fmt.Println("Missing", elements[0])
				continue
			}

			newBoundary := Boundary{
				Kind:       "Feature",
				Properties: boundaryDefinition,
				Geometry: Geometry{
					Kind:        "Polygon",
					Coordinates: [][][]float64{{}},
				},
			}

			output = append(output, newBoundary)
		}

	}
	return output
}
