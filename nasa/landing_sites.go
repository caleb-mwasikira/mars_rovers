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

type LandingSite struct {
	Rover    string
	Site     string
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

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			// If there is an error reading a single line, we skip it
			// and continue to the next line
			log.Default().Printf("Failed to read line on csv file: %v\n", err)
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

	latitude, err := coordToDecimal(record[2])
	if err != nil {
		return landingSite, err
	}

	longitude, err := coordToDecimal(record[3])
	if err != nil {
		return landingSite, err
	}

	landingSite.Rover = record[0]
	landingSite.Site = record[1]
	landingSite.Location.Latitude = latitude
	landingSite.Location.Longitude = longitude

	return landingSite, nil
}

func CalcLandingSiteDistances(world World, landingSites []LandingSite) map[string]float64 {
	distances := make(map[string]float64)
	var computed []string

	for _, siteOne := range landingSites {
		for _, siteTwo := range landingSites {
			if siteOne.Site == siteTwo.Site {
				// Refrain from calculating distances between a Location and itself
				continue
			}

			site1 := fmt.Sprintf("%v-%v", siteOne.Site, siteTwo.Site)
			site2 := fmt.Sprintf("%v-%v", siteTwo.Site, siteOne.Site)

			// Computing the from-to and to-from hashes and saving them as
			// computed values so we can know which distances are already calculated
			hash := md5.New()
			value1 := hash.Sum([]byte(site1))
			value2 := hash.Sum([]byte(site2))

			isComputed := slices.Contains(computed, string(value1)) || slices.Contains(computed, string(value2))
			if isComputed {
				continue
			}

			distance := world.distanceBetweenTwoLocations(siteOne.Location, siteTwo.Location)
			distances[site1] = distance

			computed = append(computed, string(value1), string(value2))
		}
	}
	return distances
}

func CalcClosestLandingSites(distances map[string]float64) {
	closestSites := ""
	smallestDistance := math.MaxFloat64

	for Location, distance := range distances {
		if distance < smallestDistance {
			smallestDistance = distance
			closestSites = Location
		}
	}

	sites := strings.Split(closestSites, "-")
	fmt.Printf("Closest landing sites are '%v' and '%v' with a distance of %v\n", sites[0], sites[1], smallestDistance)
}

func CalcFarthestLandingSites(distances map[string]float64) {
	farthestSites := ""
	largestDistance := -math.MaxFloat64

	for Location, distance := range distances {
		if distance > largestDistance {
			largestDistance = distance
			farthestSites = Location
		}
	}

	sites := strings.Split(farthestSites, "-")
	fmt.Printf("Farthest landing sites are '%v' and '%v' with a distance of %v\n", sites[0], sites[1], largestDistance)
}

func PrintLandingSites(landingSites []LandingSite) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Rover\tSite\tLatitude\tLongitude")

	for _, landingSite := range landingSites {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", landingSite.Rover, landingSite.Site, landingSite.Location.Latitude, landingSite.Location.Longitude)
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
