package funcs

import "github.com/buro9/microcosm/models"

// Stat will return the value of a statistic given it's name
func Stat(stats []models.Stat, name string) int64 {
	for _, stat := range stats {
		if stat.Metric == name {
			return stat.Value
		}
	}

	return 0
}
