package leaderboard

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/trueskill"
	"github.com/danmane/abalone/go/game"
	"math"
	"sort"
)

type Rating struct {
	Mean   float64
	Stddev float64
}

var gameInfo = skills.DefaultGameInfo

type Ratings map[int]Rating

type Ranking struct {
	PlayerID int
	Rating   Rating
	Rank     int
}

type Rankings []Ranking

func (r Rankings) Len() int {
	return len(r)
}

func (r Rankings) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rankings) Less(i, j int) bool {
	if r[i].Rating.Mean == r[j].Rating.Mean {
		// if mean is the same, the better player (i.e. one with the lower rank) will be the one we are more confident about
		return r[i].Rating.Stddev < r[j].Rating.Stddev
	}
	// otherwise, the better player (i.e. one with lower rank) is the one with the higher mean
	return r[i].Rating.Mean > r[j].Rating.Mean
}

func DefaultRatings(numPlayers int) Ratings {
	out := make(Ratings)
	for i := 0; i < numPlayers; i++ {
		out[i] = Rating{Mean: gameInfo.InitialMean, Stddev: gameInfo.InitialStddev}
	}
	return out
}

type Result struct {
	whiteID int
	blackID int
	outcome game.Outcome
}

func outcomeToRanks(o game.Outcome) []int {
	var result []int
	switch o {
	case game.WhiteWins:
		result = []int{1, 2}
	case game.BlackWins:
		result = []int{2, 1}
	case game.Tie:
		result = []int{1, 1}
	default:
		panic("got null outcome or other invalid outcome")
	}
	return result
}

func rating2srating(r Rating) skills.Rating {
	return skills.NewRating(r.Mean, r.Stddev)
}

func srating2rating(s skills.Rating) Rating {
	return Rating{Mean: s.Mean(), Stddev: s.Stddev()}
}
func (r Rating) String() string {
	return fmt.Sprintf("{μ:%.6g σ:%.6g}", r.Mean, r.Stddev)
}

func RateGames(numPlayers int, games []Result) Rankings {
	// This function is messy because the API we are depending on (GoSkills) is pretty weird. Could be cleaned up by just moving the calculation logic into our own impl
	ratings := DefaultRatings(numPlayers)
	for _, r := range games {
		whiteTeam := skills.NewTeam()
		whiteTeam.AddPlayer(r.whiteID, rating2srating(ratings[r.whiteID]))
		blackTeam := skills.NewTeam()
		blackTeam.AddPlayer(r.blackID, rating2srating(ratings[r.blackID]))

		var twoPlayerCalc trueskill.TwoPlayerCalc
		newSkills := twoPlayerCalc.CalcNewRatings(gameInfo, []skills.Team{whiteTeam, blackTeam}, outcomeToRanks(r.outcome)...)
		ratings[r.whiteID] = srating2rating(newSkills[r.whiteID])
		ratings[r.blackID] = srating2rating(newSkills[r.blackID])
	}
	ranks := make(Rankings, numPlayers)
	for id, rating := range ratings {
		ranks[id] = Ranking{PlayerID: id, Rating: rating, Rank: -1}
	}
	sort.Sort(ranks)
	var prevRating Rating
	for i, _ := range ranks {
		if prevRating == ranks[i].Rating {
			ranks[i].Rank = ranks[i-1].Rank
		} else {
			ranks[i].Rank = i + 1
		}
		prevRating = ranks[i].Rating
	}
	return ranks
}

func (rankings Rankings) ProposeGame() []int {
	if len(rankings) == 0 {
		return nil
	}
	maxUncertainty := math.Inf(-1)
	id1 := -1
	idx1 := -1
	for i, rank := range rankings {
		if rank.Rating.Stddev > maxUncertainty {
			maxUncertainty = rank.Rating.Stddev
			id1 = rank.PlayerID
			idx1 = i
		}
	}
	var id2 int
	if idx1 == 0 {
		id2 = rankings[1].PlayerID
	} else {
		id2 = rankings[idx1-1].PlayerID
	}
	return []int{id1, id2}
}
