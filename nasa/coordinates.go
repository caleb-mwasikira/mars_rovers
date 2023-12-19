package nasa

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type DecimalLocation float64

type Location struct {
	Latitude  float64
	Longitude float64
}

/*
Locations can be represented as degrees° minutes' seconds"(DMS) format
*/

type DMSLocation struct {
	Degrees   float64
	Minutes   float64
	Seconds   float64
	Direction string
}

/* Converts a DMS Location to its decimal equivalent */
func (coord *DMSLocation) toDecimal() float64 {
	// Calculate the decimal degrees based on the DMS values
	decimalDegrees := float64(coord.Degrees) + (float64(coord.Minutes) / 60) + (float64(coord.Seconds) / 3600)

	// Reverse the sign for S and W directions
	if coord.Direction == "S" || coord.Direction == "W" {
		decimalDegrees = -decimalDegrees
	}

	return decimalDegrees
}

/*
Converts decimal Locations to degrees, minutes, seconds
*/
func decimaltoDMS(decimalCoord float64) *DMSLocation {
	// Extract degrees
	degrees := math.Floor(decimalCoord)

	// Calculate minutes
	minutesDecimal := (decimalCoord - degrees) * 60
	minutes := math.Floor(minutesDecimal)

	// Calculate seconds
	seconds := (minutesDecimal - minutes) * 60
	return &DMSLocation{Degrees: degrees, Minutes: minutes, Seconds: seconds}
}

/*
Takes in a Location string and checks if it is in the
degrees,minutes,seconds(DMS) format

The regex used:

		^: 						Start of the string
	    [+-]?:					An optional plus or minus sign for positive or negative Locations.
	    (\d{1,3}): 				Capturing group for 1 to 3 digits representing degrees.
	    °?: 					An optional degree symbol.
	    \s?: 					An optional space character to separate degrees from minutes.
	    (\d{1,2}): 				Capturing group for 1 to 2 digits representing minutes.
	    '?: 					An optional single quote character (apostrophe) to separate minutes from seconds.
	    (?: ... )?: 			A non-capturing group to make the seconds and the second's symbol optional.
	    \s?:					An optional space character to separate minutes from seconds (if seconds are present).
	    (\d{1,2}(?:\.\d+)?)?: 	Capturing group for 1 to 2 digits with an optional decimal point and additional digits, representing seconds.
								The (?: ... ) is a non-capturing group to include the decimal point and additional digits, if present. The ? at the end makes the seconds and its decimal part optional.
	    "?: 					An optional double quote character to represent seconds (if seconds are present).
	    \s?: 					An optional space character before the direction (if seconds are present).
	    [NSWE]?: 				An optional character representing the direction (N, S, E, or W).
	    $: 						End of the string.
*/
func isDMSCoord(coord string) bool {
	pattern := `^[+-]?(\d{1,3})°?\s?(\d{1,2})'?(?:\s?(\d{1,2}(?:\.\d+)?)"?)?\s?[NSWE]?$`
	regex := regexp.MustCompile(pattern)
	return regex.Match([]byte(coord))
}

/*
Checks if a direction string is a valid compass direction (N, S, E, W)
*/
func isCompassDirection(direction string) bool {
	compassDirections := []string{"N", "S", "E", "W"}
	isCompass := false

	for _, _direction := range compassDirections {
		if _direction == strings.ToUpper(direction) {
			isCompass = true
			break
		}
	}

	return isCompass
}

/*
Converts a coord string into the degrees, minutes, seconds format
*/
func coordToDMS(coord string) (*DMSLocation, error) {
	replacer := strings.NewReplacer("°", " ", "'", " ", "\"", " ")
	coord = replacer.Replace(strings.Trim(coord, " "))

	parts := strings.Fields(coord)
	DMSCoord := &DMSLocation{}

	for i, item := range parts {
		isLastIndex := i == len(parts)-1

		if isLastIndex {
			if isCompassDirection(item) {
				DMSCoord.Direction = item
			}
		} else {
			switch i {
			case 0: // Degrees
				val, err := strconv.ParseFloat(item, 64)
				if err != nil {
					return nil, err
				}
				DMSCoord.Degrees = val
			case 1: // Minutes
				val, err := strconv.ParseFloat(item, 64)
				if err != nil {
					return nil, err
				}
				DMSCoord.Minutes = val
			default: // Seconds
				val, err := strconv.ParseFloat(item, 64)
				if err != nil {
					return nil, err
				}
				DMSCoord.Seconds = val
			}
		}
	}

	return DMSCoord, nil
}

/*
Converts a coord string it to its decimal form

For example:

	coordToDecimal("135°54'0\" E")
	Returns: 135.90000

Or:

	coordToDecimal("135.90000")
	Returns: 135.90000
*/
func coordToDecimal(coord string) (float64, error) {
	if isDMSCoord(coord) {
		DMSCoord, err := coordToDMS(coord)
		if err != nil {
			return 0, err
		}

		return DMSCoord.toDecimal(), nil
	}

	decimalCoord, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid Location format: %v", err)
	}

	return decimalCoord, nil
}

/*
Takes in a location string and converts it into its
corresponding Location{} struct with latitude and longitude.

For example:

	london has the location `(51°30'N, 0°08'W)`
		51°30'N -> latitude
		0°08'W  -> longitude

	This will be converted into a
		Location{
			Latitude: <>
			Longitude: <>
		}
	struct
*/
func parseLocation(location string) (*Location, error) {
	location = strings.Trim(location, "()")
	latAndLong := strings.Split(location, ",")

	_latitude := strings.Trim(latAndLong[0], " ")
	_longitude := strings.Trim(latAndLong[1], " ")

	latitude, err := coordToDecimal(_latitude)
	if err != nil {
		return nil, err
	}

	longitude, err := coordToDecimal(_longitude)
	if err != nil {
		return nil, err
	}

	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
