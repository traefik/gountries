package gountries

import "math"

// earthRadius Earth Radius in Kilometers.
const earthRadius = 6372.8

// Deg2Rad Degree to radian.
func Deg2Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// CalculatePythagorasEquirectangular returns equirectangular projection.
func CalculatePythagorasEquirectangular(lat1, lon1, lat2, lon2 float64) float64 {
	lat1 = Deg2Rad(lat1)
	lon1 = Deg2Rad(lon1)

	lat2 = Deg2Rad(lat2)
	lon2 = Deg2Rad(lon2)

	r := 6371.0 // km
	x := (lon2 - lon1) * math.Cos((lat1+lat2)/2)
	y := lat2 - lat1

	// Return Distance in Kilometers
	return math.Sqrt(x*x+y*y) * r
}

// CalculateHaversine returns distance in Kilometers.
func CalculateHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := Deg2Rad(lat2 - lat1)
	dLon := Deg2Rad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(Deg2Rad(lat1))*math.Cos(Deg2Rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
