package glo

import (
	"math"
	"testing"
)

func Test_ToStrokes(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		score   float64
		strokes float64
	}
	tests := []testArgs{
		testArgs{0.5, 0},
		testArgs{0.759746927, -1},
		testArgs{0.909090909, -2},
		testArgs{0.05, 2.557507202},
		testArgs{0.1, 1.908485019},
		testArgs{0.2, 1.204119983},
		testArgs{0.3, 0.735953571},
		testArgs{0.4, 0.352182518},
		testArgs{0.95, -2.557507202},
		testArgs{0.9, -1.908485019},
		testArgs{0.8, -1.204119983},
		testArgs{0.7, -0.735953571},
		testArgs{0.6, -0.352182518},
	}

	for _, test := range tests {
		result := ToStrokes(test.score)
		expected := test.strokes
		if math.Abs(result-expected) > tolerance {
			t.Errorf("ToStrokes(%f) = %f, want %f +-%f", test.score, result, expected, tolerance)
		}
	}
}

func Test_ToScore(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		strokes float64
		score   float64
	}
	tests := []testArgs{
		testArgs{-2, 0.909090909},
		testArgs{-1, 0.759746927},
		testArgs{0, 0.5},
		testArgs{1, 0.240253073},
		testArgs{2, 0.090909091},
		testArgs{3, 0.03065343},
		testArgs{4, 0.00990099},
		testArgs{5, 0.003152309},
	}

	for _, test := range tests {
		result := ToScore(test.strokes)
		expected := test.score
		if math.Abs(result-expected) > tolerance {
			t.Errorf("ToScore(%f) = %f, want %f +-%f", test.strokes, result, expected, tolerance)
		}
	}
}

func Test_CalcExpectedScore(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		holeRating   float64
		playerRating float64
		score        float64
	}
	tests := []testArgs{
		testArgs{1500, 1500, 0.5},
		testArgs{2000, 2000, 0.5},
		testArgs{1500, 1680, 0.759746927},
		testArgs{1500, 1860, 0.909090909},
		testArgs{1680, 1500, 0.240253073},
		testArgs{1860, 1500, 0.090909091},
	}

	for _, test := range tests {
		result := CalcExpectedScore(test.holeRating, test.playerRating)
		expected := test.score
		if math.Abs(result-expected) > tolerance {
			t.Errorf("CalcExpectedScore(%f, %f) = %f, want %f +-%f", test.holeRating, test.playerRating, result, expected, tolerance)
		}
	}
}

func Test_CalcPerformanceRating(t *testing.T) {
	tolerance := 0.005
	type testArgs struct {
		data              performanceRatingData
		performanceRating float64
	}
	tests := []testArgs{
		testArgs{
			performanceRatingData{
				[]float64{1500},
				0.5, 0, 0, 0,
			},
			1500,
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				1, 0, 0, 0,
			},
			1500,
		},
		testArgs{
			performanceRatingData{
				[]float64{1300, 1400},
				1, 0, 0, 0,
			},
			1350,
		},
		testArgs{
			performanceRatingData{
				[]float64{1300, 1400, 1800},
				1.5, 0, 0, 0,
			},
			1480,
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				1.519493854, 0, 0, 0,
			},
			1680,
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				2, 0, 0, 0,
			},
			3000,
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				2, 0, 2000, 0,
			},
			2000,
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				0, 0, 0, 0,
			},
			5.859, // Lower with more iterations
		},
		testArgs{
			performanceRatingData{
				[]float64{1500, 1500},
				0, 1000, 0, 0,
			},
			1000,
		},
	}

	for index, test := range tests {
		result := CalcPerformanceRating(test.data)
		expected := test.performanceRating
		if math.Abs(result-expected) > tolerance*expected {
			t.Errorf("Test %d: CalcPerformanceRating() = %f, want %f +-%f", index, result, expected, expected*tolerance)
		}
	}
}

func Test_ModifyPlayerRating(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		playerRating         float64
		performanceRating    float64
		modifiedPlayerRating float64
	}
	tests := []testArgs{
		testArgs{1500, 1500, 1500},
		testArgs{2000, 2000, 2000},
		testArgs{1500, 1680, 1536},
		testArgs{1500, 1320, 1464},
	}

	for _, test := range tests {
		result := ModifyPlayerRating(test.playerRating, test.performanceRating)
		expected := test.modifiedPlayerRating
		if math.Abs(result-expected) > tolerance*expected {
			t.Errorf("ModifyPlayerRating(%f, %f) = %f, want %f +-%f",
				test.playerRating, test.performanceRating, result, expected, expected*tolerance)
		}
	}
}

func Test_ModifyHoleRating(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.001
	type testArgs struct {
		holeRating           float64
		modifiedPlayerRating float64
	}
	tests := []testArgs{
		testArgs{1500, 1500},
	}

	// ------------------------ Tests ------------------------
	for _, test := range tests {
		result := ModifyHoleRating(test.holeRating)
		expected := test.modifiedPlayerRating
		if math.Abs(result-expected) > tolerance*expected {
			t.Errorf("ModifyHoleRating(%f) = %f, want %f +-%f",
				test.holeRating, result, expected, expected*tolerance)
		}
	}
}

func Test_CalcPlayerKFactor(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.001
	type testArgs struct {
		playerRating  float64
		playerKFactor float64
	}
	tests := []testArgs{
		testArgs{1100, 28.27},
		testArgs{1200, 25.41},
		testArgs{1300, 22.64},
		testArgs{1500, 17.54},
	}

	// ------------------------ Tests ------------------------
	for _, test := range tests {
		result := CalcPlayerKFactor(test.playerRating)
		expected := test.playerKFactor
		if math.Abs(result-expected) > tolerance*expected {
			t.Errorf("CalcPlayerKFactor(%f) = %f, want %f +-%f",
				test.playerRating, result, expected, expected*tolerance)
		}
	}
}

func Test_CalcRatingUpdates(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.0005
	type testArgs struct {
		playerRating      float64
		holeRating        float64
		strokes           float64
		performanceRating float64
		newRatings        RatingResult
	}
	tests := []testArgs{
		testArgs{
			playerRating:      1500,
			holeRating:        1500,
			strokes:           0,
			performanceRating: 1500,
			newRatings: RatingResult{
				1500, 1500,
			},
		},
		testArgs{
			playerRating:      1680,
			holeRating:        1500,
			strokes:           -1,
			performanceRating: 1680,
			newRatings: RatingResult{
				1680, 1500,
			},
		},
		testArgs{
			playerRating:      1500,
			holeRating:        1500,
			strokes:           -1,
			performanceRating: 1700,
			newRatings: RatingResult{
				1503.44, 1493.13,
			},
		},
		testArgs{
			playerRating:      1480,
			holeRating:        1300,
			strokes:           -1,
			performanceRating: 1480,
			newRatings: RatingResult{
				1480, 1300,
			},
		},
		testArgs{
			playerRating:      1480,
			holeRating:        1300,
			strokes:           1,
			performanceRating: 1300,
			newRatings: RatingResult{
				1471.44, 1316.62,
			},
		},
	}

	// ------------------------ Tests ------------------------
	for index, test := range tests {
		result := CalcRatingUpdates(test.playerRating, test.holeRating, test.strokes, test.performanceRating)
		expected := test.newRatings
		if math.Abs(result.PlayerRating-expected.PlayerRating) > tolerance*expected.PlayerRating {
			t.Errorf("Test %d Player: CalcRatingUpdates() = %f, want %f +-%f",
				index, result, expected, expected.PlayerRating*tolerance)
		}
		if math.Abs(result.HoleRating-expected.HoleRating) > tolerance*expected.HoleRating {
			t.Errorf("Test %d Hole: CalcRatingUpdates() = %f, want %f +-%f",
				index, result, expected, expected.HoleRating*tolerance)
		}
	}
}
