package weatherdata

import (
	"fmt"
	"server/pkg/models"
)

// AverageTemperature calculates the average temperature given a CityClimateData struct.
func AverageTemperature(data models.CityClimateData) float64 {
    var sum float64
    for _, feature := range data.Features {
        sum += feature.Properties.Values
    }
    return sum / float64(len(data.Features))
}

func TemperatureNext3H(data []float64) (float64, error) {
    if len(data) < 3 {
        return 0, fmt.Errorf("not enough data points to calculate the next 3 hours")
    }

    sum := 0.0
    for i := 0; i < 3; i++ {
        sum += data[i]
    }

    return sum / 3, nil
}


// temperatureNext1H calculates the average temperature for the next 3 hours given an array of float64 temperatures.
func TemperatureNext6H(data []float64) (float64, error) {
    if len(data) < 6 {
        return 0, fmt.Errorf("not enough data points to calculate the next 6 hours")
    }

    sum := 0.0
    for i := 0; i < 6; i++ {
        sum += data[i]
    }

    return sum / 6, nil
}

// PeakMeteoWindspeed returns the maximum windspeed in the data.
func PeakMeteoWindspeed(data models.MeteoBlueData) float64 {
    var max float64
    for i, windspeed := range data.Data1H.Windspeed {
        if windspeed > max {
            max = windspeed

            // stop after 12 hours
            if i > 12 {
                break
            }
        }
    }
    return max
}

// PeakMeteoTemperature returns the maximum temperature in the data and the corresponding times.
func PeakMeteoTemperature(data models.MeteoBlueData) (float64, string) {
    var max float64
    var timeOfMax string

    for i, temp := range data.Data1H.Temperature {
        if i == 0 || temp > max {  // Initialize max with the first element or update it
            max = temp
            timeOfMax = data.Data1H.Time[i]  // Assuming a corresponding Time slice

            // stop after 12 hours
            if i > 12 {
                break
            }
        }
    }

    return max, timeOfMax
}


// willItRain returns a slice of timestamps when the rain probability exceeds 50%.
func WillItRain(data models.MeteoBlueData) ([]string) {
    var times []string
    for i, probability := range data.Data1H.PrecipitationProbability {
        if probability > 50 {

            times = append(times, data.Data1H.Time[i])
            if len(times) == 8 {
                break
            }
        }
    }
    return times
}

// willItSnow returns a slice of timestamps when the snow fraction is more than 0.5.
func WillItSnow(data models.MeteoBlueData) ([]string) {
    var times []string
    for i, snowFraction := range data.Data1H.SnowFraction {
        if snowFraction > 0.5 {
            times = append(times, data.Data1H.Time[i])
            if len(times) == 8 {
                break
            }
        }
    }
    return times
}

// willItBeFoggy returns a slice of timestamps when foggy conditions are detected (pictocode == 3).
func WillItBeFoggy(data models.MeteoBlueData) ([]string) {
    var times []string
    for i, pictocode := range data.Data1H.Pictocode {
        if pictocode == 3 {
            times = append(times, data.Data1H.Time[i])
            if len(times) == 8 {
                break
            }
        }
    }
    return times
}

// willItBeWindy returns a slice of timestamps when the windspeed exceeds 8.
func WillItBeWindy(data models.MeteoBlueData) ([]string) {
    var times []string
    for i, windspeed := range data.Data1H.Windspeed {
        if windspeed > 6 {
            times = append(times, data.Data1H.Time[i])
            if len(times) == 8 {
                break
            }
        }
    }
    return times
}


// WillHaveHighUVIndex returns a slice of timestamps when the UV index exceeds 4.
func WillHaveHighUVIndex(data models.MeteoBlueData) ([]string) {
    var times []string
    for i, uvIndex := range data.Data1H.UVIndex {
        if uvIndex > 5 {
            times = append(times, data.Data1H.Time[i])
            if len(times) == 8 {
                break
            }
        }
    }
    return times
}