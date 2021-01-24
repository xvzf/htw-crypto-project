package crypt

import (
	"fmt"
	"math"
	"sort"

	"github.com/xvzf/htw-crypto-project/pkg/crypt"
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

type Analyse struct {
	ExpectedDimension image.Dimension
	Total             int
	Frequency         map[crypt.PixelPosition]int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func expand(a, b crypt.PixelPosition) crypt.PixelPosition {
	return crypt.PixelPosition{
		Width:  max(a.Width, b.Width),
		Height: max(a.Height, b.Height),
	}
}

// Load ciphertext
func Load(in []crypt.PixelPosition) *Analyse {

	a := &Analyse{
		Total:     len(in),
		Frequency: make(map[crypt.PixelPosition]int),
	}

	tmpDim := crypt.PixelPosition{Width: 0, Height: 0}

	// Retrieve pixelposition frequency & potential image dimensions
	for _, ppos := range in {

		// Update image dimensions
		tmpDim = expand(tmpDim, ppos)

		// Frequency
		if v, ok := a.Frequency[ppos]; ok {
			a.Frequency[ppos] = v + 1
		} else {
			a.Frequency[ppos] = 1
		}
	}

	// Update expected image dimension
	a.ExpectedDimension = image.Dimension{Width: tmpDim.Width + 1, Height: tmpDim.Height + 1}

	return a
}

// stat returns mean & standard deviation
func stat(in []PixelPositionFrequency) (float64, float64) {
	var sum int
	var sumSquared int
	var stdev float64
	var mean float64

	sum = 0
	sumSquared = 0
	for _, v := range in {
		sum += v.Count
		sumSquared += v.Count * v.Count
	}

	mean = float64(sum) / float64(len(in))
	stdev = math.Sqrt((float64(sumSquared) / float64(len(in))) - mean*mean)

	return mean, stdev
}

func sum(in []PixelPositionFrequency) int {
	sum := 0
	for _, v := range in {
		sum += v.Count
	}
	return sum
}

func extractDatapoints(in []PixelPositionFrequency) []int {
	var out []int
	for _, v := range in {
		out = append(out, v.Count)
	}
	return out
}

func cluster(in []PixelPositionFrequency, minStdev float64) [][]PixelPositionFrequency {
	var clusters [][]PixelPositionFrequency

	var cluster []PixelPositionFrequency

	// Sort input data for clustering based on standard deviation
	sort.SliceStable(in, func(i, j int) bool {
		return in[i].Count < in[j].Count
	})

	// Iterate sorted PixelPositions based on their count
	for _, v := range in {

		// First two values go into a cluster
		if len(cluster) < 2 {
			cluster = append(cluster, v)
			continue
		}

		// Calculate mean&stdev of cluster
		mean, stdev := stat(cluster)
		if stdev < minStdev {
			stdev = minStdev
		}

		if math.Abs(mean-float64(v.Count)) > stdev {
			clusters = append(clusters, cluster)
			fmt.Println(extractDatapoints(cluster))
			cluster = []PixelPositionFrequency{}
		} else {
			cluster = append(cluster, v)
		}
	}

	return clusters
}

// Extract Groups by clustering the pixel frequency into clusters
// It exploits the uniform distribution of a pseudo-random generator.
func (a *Analyse) ExtractGroups() [][]PixelGroupFrequency {
	var sorted []PixelPositionFrequency
	for k, v := range a.Frequency {
		sorted = append(sorted, PixelPositionFrequency{
			PixelPosition: k,
			Count:         v,
		})
	}

	scalingFactor := float64(len(sorted)) * 0.0025

	fmt.Println(len(cluster(sorted, scalingFactor)))

	return nil
}
