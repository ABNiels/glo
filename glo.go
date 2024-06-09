/*
Glo Rating System

Details about the algorithm can be found here:
https://github.com/ABNiels/glo/blob/main/README.md
*/

package glo

import (
	"math"
)

type GloRating = float64
type GloScore = float64

const (
	K_HOLE           = 35
	K_PLAYER_DEFAULT = 12
	R_WEIGHT         = 0.2
	RD               = 360
)

/* Convert strokes (-inf, inf) to score (0, 1) */
func ToScore(strokes float64) GloScore {
	// Could just use a lookup table if no use case for converting float strokes
	return 1 / (1 + math.Pow(10, strokes/2))
}

/* Convert score (0, 1)to strokes (-inf, inf) */
func ToStrokes(score GloScore) float64 {
	return 2 * math.Log10((1-score)/score)
}

/* Calculate expected score */
func CalcExpectedScore(holeRating GloRating, playerRating GloRating) GloScore {
	return 1 / (1 + math.Pow(10, (holeRating-playerRating)/RD))
}

type PerformanceRatingData struct {
	HoleRatings []GloRating
	TotalScore  GloScore
	MinReturn  GloRating
	MaxReturn  GloRating
	Tolerance   float64
}

func CalcPerformanceRating(data PerformanceRatingData) GloScore {

	if data.Tolerance == 0 {
		data.Tolerance = 0.25
	}
	if data.MaxReturn == 0 {
		data.MaxReturn = 3000
	}
	// TODO: Add input validation

	sum := 0.0
	offset := (data.MaxReturn - data.MinReturn) / 2
	performanceRating := data.MinReturn + offset

	for offset > data.Tolerance {
		offset /= 2
		sum = 0
		for _, holeRating := range data.HoleRatings {
			sum += CalcExpectedScore(holeRating, performanceRating)
		}
		if sum < data.TotalScore {
			performanceRating += offset
		} else if sum > data.TotalScore {
			performanceRating -= offset
		} else { // Unlikely
			return performanceRating
		}
	}
	return performanceRating
}

func ModifyPlayerRating(playerRating GloRating, performanceRating GloRating) float64 {
	return playerRating + R_WEIGHT*(performanceRating-playerRating)
}

func ModifyHoleRating(holeRating GloRating, details ...float64) float64 {
	// TODO: Decide what conditions/values to apply
	return holeRating
}

func CalcPlayerKFactor(playerRating GloRating) float64 {
	// TODO: Optimize/rework this equation
	if playerRating < 1900 {
		return 16 * math.Sqrt(0.5625+math.Pow(1900-playerRating, 2)/250000)
	}
	return K_PLAYER_DEFAULT
}

type StreamRatingData struct {
	PlayerRating      GloRating
	HoleRating        GloRating
	PerformanceRating GloRating
	Strokes           float64
}
type RatingResult struct {
	PlayerRating GloRating
	HoleRating   GloRating
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
	PlayerRating       GloRating
	HoleRatings        []GloRating
	PerformanceRatings []GloRating
	Strokes            []float64
}

func BatchRatingUpdate(data BatchRatingData) GloRating {
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
