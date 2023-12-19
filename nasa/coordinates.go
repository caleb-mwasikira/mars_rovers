package nasa

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type DMSCoordinate struct {
	degrees   float64
	minutes   float64
	seconds   float64
	direction string
}

/*
Converts a coordinate string; in the format degrees, minutes, seconds, direction
to its decimal equivalent.
For example:

	coordToDecimal("135°54'0\" E")
	Returns: 135.90000, nil

And:

	coordToDecimal("1234")
	Returns: 0, fmt.Error("Invalid coordinate format")
*/
func coordToDecimal(coord string) (float64, error) {
	dmsCoord, err := parseCoordinate(coord)
	if err != nil {
		return 0, err
	}

	decimalDegrees := float64(dmsCoord.degrees) + (float64(dmsCoord.minutes) / 60) + (float64(dmsCoord.seconds) / 3600)

	// Reverse the sign for S and W directions
	if dmsCoord.direction == "S" || dmsCoord.direction == "W" {
		decimalDegrees = -decimalDegrees
	}

	return decimalDegrees, nil
}

/*
Converts coord string into its equivalent DMSCoordinate
Returns an error if the coord string is not a valid coordinate
*/
func parseCoordinate(coord string) (DMSCoordinate, error) {
	var dmsCoord DMSCoordinate

	/*
		Check if coord has a valid coordinate format
		The regex used:

			^: 						Start of the string
		    [+-]?:					An optional plus or minus sign for positive or negative sign.
		    (\d{1,3}): 				Capturing group for 1 to 3 digits representing degrees.
		    °?: 					An optional degree symbol.
		    \s?: 					An optional space character to separate degrees from minutes.
			[+-]?:					An optional plus or minus sign for positive or negative sign.
		    (\d{1,2}): 				Capturing group for 1 to 2 digits representing minutes.
		    '?: 					An optional single quote character (apostrophe) to separate minutes from seconds.
		    (?: ... )?: 			A non-capturing group to make the seconds and the second's symbol optional.
		    \s?:					An optional space character to separate minutes from seconds (if seconds are present).
			[+-]?:					An optional plus or minus sign for positive or negative sign.
		    (\d{1,2}(?:\.\d+)?)?: 	Capturing group for 1 to 2 digits with an optional decimal point and additional digits, representing seconds.
			The (?: ... ) is a non-capturing group to include the decimal point and additional digits, if present. The ? at the end makes the seconds and its decimal part optional.
		    "?: 					An optional double quote character to represent seconds (if seconds are present).
		    \s?: 					An optional space character before the direction (if seconds are present).
		    [NSWE]?: 				An optional character representing the direction (N, S, E, or W).
		    $: 						End of the string.
			*/
			pattern := `^[+-]?(\d{1,3})°?\s?[+-]?(\d{1,2})'?(?:\s?[+-]?(\d{1,2}(?:\.\d+)?)"?)?\s?[NSWE]?$`
			regex := regexp.MustCompile(pattern)
			if !regex.Match([]byte(coord)) {
				return dmsCoord, fmt.Errorf("Invalid coordinate pattern for coord %v\n", coord)
	}
	
	expectedCoordLen := 4
	replacer := strings.NewReplacer("°", " ", "'", " ", "\"", " ")
	coord = replacer.Replace(strings.Trim(coord, " "))
	parts := strings.Fields(coord)
	if len(parts) != expectedCoordLen {
		return dmsCoord, fmt.Errorf("Invalid coordinate length")
	}


	// Packing degrees, minutes, seconds, direction values into a
	// DMSCoordinate struct
	for i, item := range parts {
		isLastIndex := i == len(parts)-1

		if isLastIndex {
			// Check if string is a valid compass direction
			compassDirections := []string{"N", "S", "E", "W"}
			isCompassDirection := slices.Contains[[]string](compassDirections, item)

			if isCompassDirection {
				dmsCoord.direction = item
			}
		} else {
			val, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return dmsCoord, err
			}

			switch i {
			case 0:
				dmsCoord.degrees = val
			case 1:
				dmsCoord.minutes = val
			default:
				dmsCoord.seconds = val
			}
		}
	}

	return dmsCoord, nil
}
