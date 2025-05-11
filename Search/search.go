package Search

import (
	"SEProject/Restaurant/types"
	"strings"
)

func SearchRestaurants(restaurants []types.Restaurant, keyword string) []types.Restaurant {
	var result []types.Restaurant
	for _, r := range restaurants {
		if strings.Contains(strings.ToLower(r.Name), strings.ToLower(keyword)) {
			result = append(result, r)
		}
	}
	return result
}

func FilterByCuisine(restaurants []types.Restaurant, cuisine string) []types.Restaurant {
	var result []types.Restaurant
	for _, r := range restaurants {
		if strings.EqualFold(r.Cuisine, cuisine) {
			result = append(result, r)
		}
	}
	return result
}

func FilterByLocation(restaurants []types.Restaurant, location string) []types.Restaurant {
	var result []types.Restaurant
	for _, r := range restaurants {
		if strings.EqualFold(r.Location, location) {
			result = append(result, r)
		}
	}
	return result
}
func FilterByPrice(restaurants []types.Restaurant, maxPrice int) []types.Restaurant {
	var result []types.Restaurant
	for _, r := range restaurants {
		if r.AvgPrice <= maxPrice {
			result = append(result, r)
		}
	}
	return result
}

func FilterByRating(restaurants []types.Restaurant, minRating float64) []types.Restaurant {
	var result []types.Restaurant
	for _, r := range restaurants {
		if r.Rating >= minRating {
			result = append(result, r)
		}
	}
	return result
}
