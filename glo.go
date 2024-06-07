/*
Glo Rating System

Details about the algorithm can be found here:
https://github.com/ABNiels/glo/blob/main/README.md
*/

package glo

import (
	"math"
)

const (
	K_HOLE           = 35
	K_PLAYER_DEFAULT = 12
	R_WEIGHT         = 0.2
	RD               = 360
)

/* Convert strokes (-inf, inf) to score (0, 1) */
func ToScore(strokes float64) float64 {
	// Could just use a lookup table if no use case for converting float strokes
	return 1 / (1 + math.Pow(10, strokes/2))
}

/* Convert score to strokes */
func ToStrokes(score float64) float64 {
	return 2 * math.Log10((1-score)/score)
}

/* Calculate expected score */
func CalcExpectedScore(holeRating float64, playerRating float64) float64 {
	return 1 / (1 + math.Pow(10, (holeRating-playerRating)/RD))
}

type performanceRatingData struct {
	holeRatings []float64
	totalScore  float64
	min_return  float64
	max_return  float64
	iterations  int
}

func CalcPerformanceRating(data performanceRatingData) float64 {

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
			sum += CalcExpectedScore(holeRating, performanceRating)
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

func ModifyPlayerRating(playerRating float64, performanceRating float64) float64 {
	return playerRating + R_WEIGHT*(performanceRating-playerRating)
}

func ModifyHoleRating(holeRating float64, details ...float64) float64 {
	// TODO: Decide what conditions/values to apply
	return holeRating
}

func CalcPlayerKFactor(playerRating float64) float64 {
	// TODO: Optimize/rework this equation
	if playerRating < 1900 {
		return 16 * math.Sqrt(0.5625+math.Pow(1900-playerRating, 2)/250000)
	}
	return K_PLAYER_DEFAULT
}

type StreamRatingData struct {
	PlayerRating      float64
	HoleRating        float64
	PerformanceRating float64
	Strokes           float64
}
type RatingResult struct {
	PlayerRating float64
	HoleRating   float64
}

func StreamRatingUpdate(data StreamRatingData) RatingResult {

	modifiedHoleRating := ModifyHoleRating(data.HoleRating)
	modifiedPlayerRating := ModifyPlayerRating(data.PlayerRating, data.PerformanceRating)

	expectedScore := CalcExpectedScore(modifiedHoleRating, modifiedPlayerRating)
	actualScore := ToScore(data.Strokes)

	playerKFactor := CalcPlayerKFactor(data.PlayerRating)

	newPlayerRating := data.PlayerRating + playerKFactor*(actualScore-expectedScore)
	newHoleRating := data.HoleRating + K_HOLE*(expectedScore-actualScore)

	return RatingResult{newPlayerRating, newHoleRating}
}

type BatchRatingData struct {
	PlayerRating       float64
	HoleRatings        []float64
	PerformanceRatings []float64
	Strokes            []float64
}

func BatchRatingUpdate(data BatchRatingData) float64 {
	totalExpectedScore := 0.0
	totalActualScore := 0.0
	modifiedPlayerRating := 0.0
	modifiedHoleRating := 0.0

	for i := 0; i < len(data.HoleRatings); i++ {
		modifiedPlayerRating = ModifyPlayerRating(data.PlayerRating, data.PerformanceRatings[i])
		modifiedHoleRating = ModifyHoleRating(data.HoleRatings[i])

		expectedScore := CalcExpectedScore(modifiedHoleRating, modifiedPlayerRating)
		actualScore := ToScore(data.Strokes[i])

		totalExpectedScore += expectedScore
		totalActualScore += actualScore
	}

	playerKFactor := CalcPlayerKFactor(data.PlayerRating)

	newPlayerRating := data.PlayerRating + playerKFactor*(totalActualScore-totalExpectedScore)
	return newPlayerRating
}
