package glo

import (
	"math"
	"testing"
)

func Test_ToStrokes(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		score   GloScore
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
		score   GloScore
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
		holeRating   GloRating
		playerRating GloRating
		score        GloScore
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
	type testArgs struct {
		data              performanceRatingData
		performanceRating GloRating
	}
	tests := []testArgs{
		testArgs{
			performanceRatingData{
				[]GloRating{1500},
				0.5, 0, 0, 0,
			},
			1500,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				1, 0, 0, 10,
			},
			1500,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1300, 1400},
				1, 0, 0, 0,
			},
			1350,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1300, 1400, 1800},
				1.5, 0, 0, 0,
			},
			1480,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				1.519493854, 0, 0, 0,
			},
			1680,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				2, 0, 0, 0,
			},
			3000,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				2, 0, 2000, 0,
			},
			2000,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				0, 0, 0, 0,
			},
			0,
		},
		testArgs{
			performanceRatingData{
				[]GloRating{1500, 1500},
				0, 1000, 0, 0,
			},
			1000,
		},
	}

	for index, test := range tests {
		result := CalcPerformanceRating(test.data)
		expected := test.performanceRating
		tolerance := test.data.tolerance
		if tolerance == 0 {
			tolerance = 0.25
		}

		if math.Abs(result-expected) > tolerance {
			t.Errorf("Test %d: CalcPerformanceRating() = %f, want %f +-%f", index, result, expected, tolerance)
		}
	}
}

func Test_ModifyPlayerRating(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		playerRating         GloRating
		performanceRating    GloRating
		modifiedPlayerRating GloRating
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
		holeRating           GloRating
		modifiedPlayerRating GloRating
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
		playerRating  GloRating
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

func Test_StreamRatingUpdate(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.0005
	type testArgs struct {
		data       StreamRatingData
		newRatings RatingResult
	}
	tests := []testArgs{
		testArgs{
			data: StreamRatingData{
				PlayerRating:      1500,
				HoleRating:        1500,
				PerformanceRating: 1500,
				Strokes:           0,
			},
			newRatings: RatingResult{
				1500, 1500,
			},
		},
		testArgs{
			data: StreamRatingData{
				PlayerRating:      1680,
				HoleRating:        1500,
				PerformanceRating: 1680,
				Strokes:           -1,
			},
			newRatings: RatingResult{
				1680, 1500,
			},
		},
		testArgs{
			data: StreamRatingData{
				PlayerRating:      1480,
				HoleRating:        1300,
				PerformanceRating: 1480,
				Strokes:           -1,
			},
			newRatings: RatingResult{
				1480, 1300,
			},
		},
		testArgs{
			data: StreamRatingData{
				PlayerRating:      1480,
				HoleRating:        1300,
				PerformanceRating: 1300,
				Strokes:           1,
			},
			newRatings: RatingResult{
				1471.44, 1316.62,
			},
		},
	}

	// ------------------------ Tests ------------------------
	for index, test := range tests {
		result := StreamRatingUpdate(test.data)
		expected := test.newRatings
		if math.Abs(result.PlayerRating-expected.PlayerRating) > tolerance*expected.PlayerRating {
			t.Errorf("Test %d Player: StreamRatingUpdate() = %f, want %f +-%f",
				index, result, expected, expected.PlayerRating*tolerance)
		}
		if math.Abs(result.HoleRating-expected.HoleRating) > tolerance*expected.HoleRating {
			t.Errorf("Test %d Hole: StreamRatingUpdate() = %f, want %f +-%f",
				index, result, expected, expected.HoleRating*tolerance)
		}
	}
}

func Test_BatchRatingUpdate(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.0005
	type testArgs struct {
		data      BatchRatingData
		newRating GloRating
	}
	tests := []testArgs{
		testArgs{
			data: BatchRatingData{
				PlayerRating:       1500,
				HoleRatings:        []GloRating{1500, 1500, 1500},
				PerformanceRatings: []GloRating{1500, 1500, 1500},
				Strokes:            []GloScore{0, 0, 0},
			},
			newRating: 1500,
		},
		testArgs{
			data: BatchRatingData{
				PlayerRating:       1500,
				HoleRatings:        []GloRating{1680, 1500},
				PerformanceRatings: []GloRating{1500, 1500},
				Strokes:            []GloScore{0, 0},
			},
			newRating: 1504.55,
		},
		testArgs{
			data: BatchRatingData{
				PlayerRating:       1500,
				HoleRatings:        []GloRating{1680, 1500, 1500, 1500, 1500, 1500},
				PerformanceRatings: []GloRating{1500, 1500, 1500, 1500, 1500, 1500},
				Strokes:            []GloScore{0, 0, 0, 0, 0, 0},
			},
			newRating: 1504.55,
		},
		testArgs{
			data: BatchRatingData{
				PlayerRating:       1500,
				HoleRatings:        []GloRating{1500, 1500, 1500, 1500, 1500, 1500},
				PerformanceRatings: []GloRating{1500, 1500, 1500, 1500, 1500, 1500},
				Strokes:            []GloScore{-1, 1, -1, 1, -1, 1},
			},
			newRating: 1500,
		},
	}

	// ------------------------ Tests ------------------------
	for index, test := range tests {
		result := BatchRatingUpdate(test.data)
		expected := test.newRating
		if math.Abs(result-expected) > tolerance*expected {
			t.Errorf("Test %d Player: BatchRatingUpdate() = %f, want %f +-%f",
				index, result, expected, expected*tolerance)
		}
	}
}
