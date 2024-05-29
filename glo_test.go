package glo

import (
	"testing"
	"math"
)


func Test_toStrokes(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		score float64
		strokes float64
	}
	tests := []testArgs {
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
		result := toStrokes(test.score)
		expected := test.strokes
		if math.Abs(result - expected) > tolerance {
			t.Errorf("toStrokes(%f) = %f, want %f +-%f", test.score, result, expected, tolerance)
		}
	}
}

func Test_toScore(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		strokes float64
		score float64
	}
	tests := []testArgs {
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
		result := toScore(test.strokes)
		expected := test.score
		if math.Abs(result - expected) > tolerance {
			t.Errorf("toScore(%f) = %f, want %f +-%f", test.strokes, result, expected, tolerance)
		}
	}
}

func Test_calcExpectedScore(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		holeRating float64
		playerRating float64
		score float64
	}
	tests := []testArgs {
		testArgs{1500, 1500, 0.5},
		testArgs{2000, 2000, 0.5},
		testArgs{1500, 1680, 0.759746927},
		testArgs{1500, 1860, 0.909090909},
		testArgs{1680, 1500, 0.240253073},
		testArgs{1860, 1500, 0.090909091},
	}


	for _, test := range tests {
		result := calcExpectedScore(test.holeRating, test.playerRating)
		expected := test.score
		if math.Abs(result - expected) > tolerance {
			t.Errorf("calcExpectedScore(%f, %f) = %f, want %f +-%f", test.holeRating, test.playerRating, result, expected, tolerance)
		}
	}
}

func Test_calcPerformanceRating(t *testing.T) {
	tolerance := 0.005
	type testArgs struct {
		data performanceRatingData
		performanceRating float64
	}
	tests := []testArgs {
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
		result := calcPerformanceRating(test.data)
		expected := test.performanceRating
		if math.Abs(result - expected) > tolerance * expected {
			t.Errorf("Test %d: calcPerformanceRating() = %f, want %f +-%f", index, result, expected, expected*tolerance)
		}
	}
}

func Test_modifyPlayerRating(t *testing.T) {
	tolerance := 0.001
	type testArgs struct {
		playerRating float64
		performanceRating float64
		modifiedPlayerRating float64
	}
	tests := []testArgs {
		testArgs{1500, 1500, 1500},
		testArgs{2000, 2000, 2000},
		testArgs{1500, 1680, 1536},
		testArgs{1500, 1320, 1464},
	}


	for _, test := range tests {
		result := modifyPlayerRating(test.playerRating, test.performanceRating)
		expected := test.modifiedPlayerRating
		if math.Abs(result - expected) > tolerance * expected {
			t.Errorf("modifyPlayerRating(%f, %f) = %f, want %f +-%f",
			test.playerRating, test.performanceRating, result, expected, expected*tolerance)
		}
	}
}

func Test_modifyHoleRating(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.001
	type testArgs struct {
		holeRating float64
		modifiedPlayerRating float64
	}
	tests := []testArgs {
		testArgs{1500, 1500},
	}

	// ------------------------ Tests ------------------------
	for _, test := range tests {
		result := modifyHoleRating(test.holeRating)
		expected := test.modifiedPlayerRating
		if math.Abs(result - expected) > tolerance * expected {
			t.Errorf("modifyHoleRating(%f) = %f, want %f +-%f", 
			test.holeRating, result, expected, expected*tolerance)
		}
	}
}

func Test_calcPlayerKFactor(t *testing.T) {
	// ------------------------ Setup ------------------------
	tolerance := 0.001
	type testArgs struct {
		playerRating float64
		playerKFactor float64
	}
	tests := []testArgs {
		testArgs{1100, 28.27},
		testArgs{1200, 25.41},
		testArgs{1300, 22.64},
		testArgs{1500, 17.54},
	}

	// ------------------------ Tests ------------------------
	for _, test := range tests {
		result := calcPlayerKFactor(test.playerRating)
		expected := test.playerKFactor
		if math.Abs(result - expected) > tolerance * expected {
			t.Errorf("calcPlayerKFactor(%f) = %f, want %f +-%f", 
			test.playerRating, result, expected, expected*tolerance)
		}
	}
}