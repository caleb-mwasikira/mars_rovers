package main

import (
	"fmt"
	"log"

	nasa "bellweathertech.com/mars_rover/nasa"
)

func main() {
	var (
		err error
	)

	filename := "data/rover_data_10.csv"
	landingSites, err := nasa.ReadLandingSitesFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	nasa.PrintLandingSites(landingSites)
	fmt.Println()

	mars := nasa.World{
		Radius: 3389.5,
	}
	landingSiteDistances := nasa.CalcLandingSiteDistances(mars, landingSites)
	nasa.PrintLandingSiteDistances(landingSiteDistances)

	return
}
