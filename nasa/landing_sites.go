package nasa

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

type World struct {
	Radius float64
}

type LandingSite struct {
	Rover    string
	Location Location
}

/*
Reads csv file line by line and tries to convert each record into
a LandingSite object
*/
func ReadLandingSitesFromFile(filename string) ([]LandingSite, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var landingSites []LandingSite
	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true

	i := 0
	for {
		// Do not read the first line of the csv file
		if i == 0 {
			i++
			continue
		}

		record, err := csvReader.Read()
		if err != nil {
			// Stop reading the file if End of file reached
			if err == io.EOF {
				break
			}

			log.Default().Printf("Failed to read line on csv file")
			continue
		}

		landingSite, err := createLandingSite(record)
		if err != nil {
			// If there is an error converting a record in a line, we log it
			// and continue to the next line
			log.Default().Println(err)
			continue
		}

		landingSites = append(landingSites, landingSite)
	}

	return landingSites, nil
}

/*
Creates a LandingSite object given a landingSites record
*/
func createLandingSite(record []string) (LandingSite, error) {
	var (
		landingSite          LandingSite
		expectedRecordLength int = 4
	)

	if len(record) != expectedRecordLength {
		return landingSite, fmt.Errorf("Invalid record format. Expected record length of %v, but got %v\n", expectedRecordLength, len(record))
	}

	latitude, err := coordToDecimal(record[1])
	if err != nil {
		return landingSite, err
	}

	longitude, err := coordToDecimal(record[2])
	if err != nil {
		return landingSite, err
	}

	landingSite.Rover = record[0]
	landingSite.Location.Latitude = latitude
	landingSite.Location.Longitude = longitude

	return landingSite, nil
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

/*
Calculates the distance between two locations using the Spherical Law of Cosines.
*/
func DistanceBetweenTwoLocations(world World, point1, point2 Location) float64 {
	s1, c1 := math.Sincos(degreesToRadians(point1.Latitude))
	s2, c2 := math.Sincos(degreesToRadians(point2.Latitude))

	clong := math.Cos(degreesToRadians(point1.Longitude - point2.Longitude))
	return world.Radius * math.Acos(s1*s2+c1*c2*clong)
}

func CalcLandingSiteDistances(world World, landingSites []LandingSite) map[string]float64 {
	distances := make(map[string]float64)
	var computed []string

	for _, siteOne := range landingSites {
		for _, siteTwo := range landingSites {
			if siteOne.Rover == siteTwo.Rover {
				// Refrain from calculating distances between a Location and itself
				continue
			}

			toFrom := fmt.Sprintf("%v-%v", siteOne.Rover, siteTwo.Rover)
			fromTo := fmt.Sprintf("%v-%v", siteTwo.Rover, siteOne.Rover)

			// Computing the from-to and to-from hashes and saving them as
			// computed values so we can know which distances are already calculated
			hash := md5.New()
			value1 := hash.Sum([]byte(toFrom))
			value2 := hash.Sum([]byte(fromTo))

			isComputed := slices.Contains(computed, string(value1)) || slices.Contains(computed, string(value2))
			if isComputed {
				continue
			}

			distance := DistanceBetweenTwoLocations(world, siteOne.Location, siteTwo.Location)
			distances[toFrom] = distance

			computed = append(computed, string(value1), string(value2))
		}
	}

	return distances
}

func PrintLandingSites(landingSites []LandingSite) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Rover\tLatitude\tLongitude")

	for _, landingSite := range landingSites {
		fmt.Fprintf(w, "%v\t%v\t%v\n", landingSite.Rover, landingSite.Location.Latitude, landingSite.Location.Longitude)
	}
	w.Flush()
}

func PrintLandingSiteDistances(landingSiteDistances map[string]float64) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "From\tTo\tDistance")

	for rovers, distance := range landingSiteDistances {
		parts := strings.Split(rovers, "-")
		fmt.Fprintf(w, "%v\t%v\t%v\n", parts[0], parts[1], distance)
	}
	w.Flush()
}
