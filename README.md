# Glo (G-lo) Rating System
This is an implementation of the Glo rating system for stroke play sports based on the Elo rating system

# Description of Rating System
This rating system is a modified Elo system in which an Expected Score is calculated and used to update the ratings of the hole and players.
The key difference is in how this method converts a number of strokes relative to par to a score in the range (0, 1). It also uses offsets to the players'
and holes' ratings based on external factors that provide advantages or disadvantages. These might be weather effects or hole suitability and familiarity as examples.

The system is designed to make it easy to mentally calculate 18 hole handicaps, although this is arbitrary, and the constants can be adjusted as the use case demands.
For the constants below, the a player can find their 18 hole handicap relative to another rated player by taking the difference in ratings divided by 10.
For example, a 1700 rated player would be expected to end 10 strokes better than a 1600 rated player over 18 holes.
This relationship can be easily altered via the Rating Differential Factor. Rating Difference multiplied by 2 divided by the Ratied Differential Factor gives the expected performance.
For example, with a factor of 360 and pitting a player rated 1700 against a player rated 1600, the higher rated player would expect to perform 200/360 strokes better on any given hole.
Across 18 holes, this turns into the easy mental math from before of 1/10 of the Rating Difference, but it is fairly difficult to do mentally for one single hole.
A factor of 400 would mean players with those ratings would be separated by 200/400, or half a stroke per hole. This would make per hole handicaps easier (1/20 Rating Difference per hole), but per course handicaps slightly
more difficult to mentally calculate (18/20 Rating Difference per course).

# Calculating A Rating Update
## Calculating the Expected Score
The Expected Score ($E$) is the expected value for the Score of a player at a hole. It is calculated from a Modified Rating Difference ($R$) using the Rating Differential Factor ($d$) as follows

$$ d = 360 \quad E = \frac{1}{1+10^{\frac{R}{d}}} $$

### Calculating the Modified Rating Difference
The Modified Rating Difference ($R$) is the difference between the Modified Hole Rating ($R_{mh}$) and the Modified Player Rating ($R_{mp}$). 

$$ R = R_{mh} - R_{mp} $$

The Modified Hole Rating is calculated using a static offset for each and any weather conditions present during the round.
This is hueristically assigned until data becomes available to make a statistically relevant assigment.

A Modified Player's Rating ($R_{mp}$) is a weighted average ($a$) between the Player Rating ($R_p$) and the player's Hole Performance Rating ($R_{hp}$).

$$ a = 0.2 \quad R_{mp} =  R_p + a(R_{hp} - R_p) $$

The Hole Performance Rating is the player's past Performance Rating at that specific hole. If the player has no history at the hole, it can be a static offset, the Player Rating, or the player's Performance Rating at a 
similar hole. The Performance Rating is the rating that would yield the player's Total Score ($S_t$) as the sum of the Expected Scores for that rating at that hole.
It can be limited to a certain number of historical games ($i$) to make it more recently biased and to within a certain range around the Player Rating to make the extremes less impactful.
This is the equation that must be satisfied (using the Expected Score equation with the historical Modified Hole Ratings $R_{mhi}$), which has no analytical solution but can be solved numerically for $R_{hp}$:

$$ S_t = \sum_i \frac{1}{1+10^{\frac{R_{mhi}-R_{hp}}{d}}} $$

## Calculating the New Rating

The new Rating ($R_n$) for a player is calculated by scaling ($k_p$) the difference between the Score ($S$) and Expected Score ($E$) and adding it to the old Rating ($R_o$).

$$ R_n = R_o + k_p\left(S-E\right) $$

The calculation is similar for a hole, but it uses a different scaling factor (k_h) and reverses the Score and Expected Score

$$ R_n = R_o + k_h\left(E-S\right) $$

### The Scaling Factors

Determining a Player Scaling Factor ($k_p$) is still in progress. This value determines how much a rating changes on each update with a larger K Factor leading to a larger change in rating. Ideally, a new player has a larger factor as they improve quicker and have a larger variance in scores.
A rating deviation similar to the glicko method could be created or based on the number of holes played similar to the USCF K factor, but it currently is based off of the current rating of the player ($R_p$)
using the following equation 


$$
k_p = \begin{cases} 
      16 \sqrt{0.5625 + \frac{1}{250,000} \left(1900-R_p\right)^2} & ; R_p < 1900 \\
      12 & ; R_p >= 1900 
   \end{cases}
$$

For holes, the scaling factor can be set to a constant value. This is because holes as a rule do not change in this system. If the hole changes in real life, it is considered a new hole.

### Calculating the Score

The equation for conversions between Strokes Relative to Par (P) and Scores (S) is as follows

$$ S = \frac{1}{1+10^{0.5P}} $$

$$ P = 2\log_{10}\left(\frac{1-S}{S}\right) $$

# Sample Calculation

## Expected Score
Alice has a rating of 1550 and is playing a par 3 hole rated 1600 for the first time in clear weather (negating any specific modification to the player / hole rating you want to do).
First, we calculate her Expected Score as

$$ E = \frac{1}{1+10^{\frac{1550 - 1600}{360}}} = 0.42 $$

Now we can see how many strokes above/below par that is

$$ P = 2\log_{10}\left(\frac{1-0.42}{0.42}\right) = 0.278 $$

So that means we expected Alice to average 3.278 strokes playing this hole.

## Updating the ratings

Let's say Alice overperforms and gets a birdie on her first try. The first step is convert her strokes (-1) to a score like this.

$$ S = \frac{1}{1+10^{0.5\left(-1\right)}} =  0.76 $$

Now we can find out what her new rating should be. Using the K Factor equation above, her K Factor would be 16.41, but you can use any consistent K Factor for your own purposes.

$$ R_n = 1550 + 16.41\left(0.76-0.42\right) = 1555.56 $$

And for the hole (again using a sample K Factor)

$$ R_n = 1600 + 35\left(0.42-0.76\right) = 1588.13 $$

## Mental Math
Alice and Bob want to play against each other at the local 18 hole course. Alice just got her new rating of 1555.56 and Bob is rated a whopping 1643. Because he is rated ~90 points higher, he should take a handicap of +9 to even things out for the round, or an extra stroke every other hole.

# Example Library Usage
```
/* Example usage of glo.

Calculates the performance rating and new player rating after a single hole
*/

package main
import (
	"fmt"
	"github.com/ABNiels/glo"
)

func main() {
	holeHistory := glo.PerformanceRatingData{
		HoleRatings: []glo.GloRating{1600, 1625, 1650, 1600, 1625},
		TotalScore: glo.ToScore(1) + glo.ToScore(1) + glo.ToScore(2) + glo.ToScore(1) + glo.ToScore(3),
	}

	oldPlayerRating := 1550.0
	holeRating := 1600.0
	performanceRating := glo.CalcPerformanceRating(holeHistory)

	ratingData := glo.StreamRatingData{
		PlayerRating: oldPlayerRating,
		HoleRating: holeRating,
		PerformanceRating: performanceRating,
		Strokes: 0,
	}

	ratingUpdate := glo.StreamRatingUpdate(ratingData)
	fmt.Println(fmt.Sprintf("%0.2f (%0.2f) x %0.2f -> %0.2f", oldPlayerRating, performanceRating, holeRating, ratingUpdate.PlayerRating))
}
```

# Uses
This can be used in a scorekeeping app or spreadsheet to even out players with varying ratings. Each player can instead play against their own Expected Score to keep everyone's numbers close together from the start to the end of the round.
