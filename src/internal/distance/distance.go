package distance

import (
	"ex00/internal/structures"
	"github.com/umahmood/haversine"
	"sort"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mgwyness"
	password = "etototsamiymysh"
	dbname   = "postgres"
	limit    = 20
)

func FindDistance(rests structures.Restaurant, myCoord haversine.Coord) float64 {
	restCoord := haversine.Coord{Lat: rests.Location.Latitude, Lon: rests.Location.Longitude}
	_, km := haversine.Distance(myCoord, restCoord)
	return km
}

func FindThreeRests(rests []structures.Restaurant, lat float64, lon float64) []structures.Restaurant {

	myCoord := haversine.Coord{Lat: lat, Lon: lon}

	restsMap := make(map[float64]structures.Restaurant)
	for _, rest := range rests {
		dist := FindDistance(rest, myCoord)
		restsMap[dist] = rest
	}
	restsArray := make([]float64, 0, len(restsMap))
	for key := range restsMap {
		restsArray = append(restsArray, key)
	}
	sort.Float64s(restsArray)

	var closest []structures.Restaurant

	for _, key := range restsArray {
		closest = append(closest, restsMap[key])
	}
	return closest[:3]
}
