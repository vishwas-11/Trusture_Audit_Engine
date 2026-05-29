package service

import "math"

type GeoPoint struct {
	Lat float64
	Lng float64
}

// Haversine formula
func DistanceKM(a, b GeoPoint) float64 {
	const R = 6371 // Earth radius in KM

	dLat := (b.Lat - a.Lat) * math.Pi / 180
	dLng := (b.Lng - a.Lng) * math.Pi / 180

	lat1 := a.Lat * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180

	x := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLng/2)*math.Sin(dLng/2)*math.Cos(lat1)*math.Cos(lat2)

	return R * 2 * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}
