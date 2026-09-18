package domain

import "math"

const earthRadiusM = 6371000.0

// HaversineMeters devuelve la distancia en metros entre dos coordenadas.
func HaversineMeters(a, b Coordinates) float64 {
	dLat := (b.Lat - a.Lat) * math.Pi / 180
	dLng := (b.Long - a.Long) * math.Pi / 180
	lat1 := a.Lat * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180

	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/2)*math.Sin(dLng/2)

	return 2 * earthRadiusM * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
}
