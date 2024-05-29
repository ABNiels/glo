/*
Glo Rating System

Details about the algorithm can be found here: 
https://github.com/ABNiels/DGRDatabase/blob/main/RatingEquation_Notes.md
*/

package glo

import (
	"math"
)

const (
	K_HOLE = 35
	K_PLAYER_DEFAULT = 12
	R_WEIGHT = 0.2
	RD = 360
)

/* Convert strokes (-inf, inf) to score (0, 1) */
func toScore(strokes float64) float64 {
	// Could just use a lookup table if no use case for converting float strokes
	return 1 / (1 + math.Pow(10, strokes/2))
}

/* Convert score to strokes */
func toStrokes(score float64) float64 {
	return 2 * math.Log10((1-score)/score)
}

/* Calculate expected score */
func calcExpectedScore(holeRating float64, playerRating float64) float64 {
	return 1 / (1 + math.Pow(10, (holeRating-playerRating)/RD))
}

type performanceRatingData struct {
	holeRatings []float64
	totalScore float64
	min_return float64
	max_return float64
	iterations int
}
func calcPerformanceRating(data performanceRatingData) float64 {

	if data.iterations == 0 {
		data.iterations = 8
	}
	if data.max_return == 0 {
		data.max_return = 3000
	}
	// TODO: Add input validation

	sum := 0.0
	offset := (data.max_return - data.min_return) / 2
	performanceRating := data.min_return + offset

	for i := 0; i < data.iterations; i++ {
		offset /= 2
		sum = 0
		for _, holeRating := range data.holeRatings {
			sum += calcExpectedScore(holeRating, performanceRating)
		}
		if sum < data.totalScore {
			performanceRating += offset
		} else if sum > data.totalScore {
			performanceRating -= offset
		} else { // Unlikely
			return performanceRating
		}
		// TODO: Add tolerance for early return near min/max
	}
	return performanceRating
}

func modifyPlayerRating(playerRating float64, performanceRating float64) float64 {
	return playerRating + R_WEIGHT * (performanceRating - playerRating)
}

func modifyHoleRating(holeRating float64, details ...float64) float64 {
	// TODO: Decide what conditions/values to apply
	return holeRating
}

func calcPlayerKFactor(playerRating float64) float64 {
	// TODO: Optimize/rework this equation
	if playerRating < 1900 {
		return 16 * math.Sqrt(0.5625 + math.Pow(1900 - playerRating, 2)/250000)
	} 
	return K_PLAYER_DEFAULT
}

type RatingResult struct {
	playerRating float64
	holeRating float64
}
func calcRatingUpdates(playerRating float64, holeRating float64,
                       strokes float64, performanceRating float64) RatingResult {

	modifiedHoleRating := modifyHoleRating(holeRating)
	modifiedPlayerRating := modifyPlayerRating(playerRating, performanceRating)

	expectedScore := calcExpectedScore(modifiedHoleRating, modifiedPlayerRating)
	actualScore := toScore(strokes)

	playerKFactor := calcPlayerKFactor(playerRating)

	newPlayerRating := playerRating + playerKFactor * (actualScore - expectedScore)
	newHoleRating := holeRating + K_HOLE * (expectedScore - actualScore)

	return RatingResult{newPlayerRating, newHoleRating}
}