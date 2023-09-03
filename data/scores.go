package data

import (
	"errors"
	"strconv"
)

func Log(id int64, status string) {
	redisDB.LPush(ctx, "user:"+strconv.FormatInt(id, 10)+":history", status)
}

func AddScore(id int64, score float64) {
	redisDB.ZIncrBy(ctx, "users:leaderboard", score, strconv.FormatInt(int64(id), 10))
}

func FetchScore(id int64) float64 {
	out, _ := redisDB.ZScore(ctx, "users:leaderboard", strconv.FormatInt(int64(id), 10)).Result()
	return out
}

type LeaderboardPlacement struct {
	ID    int
	Score float64
	Rank  int
}

func LeaderboardPage(page int, count int) ([]LeaderboardPlacement, error) {
	if page < 0 || count <= 0 {
		return nil, errors.New("negative number")
	}

	pageindex := page * count
	result, err := redisDB.ZRevRangeWithScores(ctx, "users:leaderboard", int64(pageindex), int64(pageindex+count-1)).Result()
	if err != nil {
		return nil, err
	}
	output := make([]LeaderboardPlacement, 0)
	for rank, z := range result {
		newLeaderboardPlacement := LeaderboardPlacement{}
		newLeaderboardPlacement.ID, err = strconv.Atoi(z.Member.(string))
		if err != nil {
			return nil, err
		}
		newLeaderboardPlacement.Score = z.Score
		newLeaderboardPlacement.Rank = rank + 1
		output = append(output, newLeaderboardPlacement)
	}

	return output, nil
}
