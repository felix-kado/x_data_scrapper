package service

import (
    "context"
    "time"

    "github.com/felix-kado/x_data_scrapper/internal/model"
    "github.com/felix-kado/x_data_scrapper/internal/twitter"
)

// Сервис для работы с пользователями Twitter
type UserService struct {
    tw *twitter.Client
}

func NewUserService(tw *twitter.Client) *UserService {
    return &UserService{tw: tw}
}

// Получить профиль пользователя
func (s *UserService) GetProfile(ctx context.Context, username string) (*model.User, error) {
    resp, err := s.tw.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }
    d := resp.Data
    u := &model.User{
        ID:           d.ID,
        Name:         d.Name,
        Username:     d.Username,
        Description:  d.Description,
        Followers:    d.PublicMetrics.Followers,
        Following:    d.PublicMetrics.Following,
        Listed:       d.PublicMetrics.Listed,
        TweetCount:   d.PublicMetrics.TweetCount,
        ProfileImage: d.ProfileImageURL,
        Verified:     d.Verified,
    }
    return u, nil
}

// Получить твиты пользователя (до лимита)
func (s *UserService) GetTweets(ctx context.Context, username string, limit int) ([]model.Tweet, error) {
    // first resolve id
prof, err := s.tw.GetUserByUsername(ctx, username)
if err != nil {
    return nil, err
}

pages, err := s.tw.GetTweets(ctx, prof.Data.ID, limit)
    if err != nil {
        return nil, err
    }
    tweets := make([]model.Tweet, 0)
    for _, p := range pages {
        for _, t := range p.Data {
            tweets = append(tweets, model.Tweet{
                ID:           t.ID,
                Text:         t.Text,
                CreatedAt:    t.CreatedAt,
                LikeCount:    t.PublicMetric.LikeCount,
                RetweetCount: t.PublicMetric.RetweetCount,
                ReplyCount:   t.PublicMetric.ReplyCount,
                QuoteCount:   t.PublicMetric.QuoteCount,
            })
        }
    }
    return tweets, nil
}

// Получить engagement-метрики пользователя
func (s *UserService) ComputeMetrics(ctx context.Context, username string, sample int) (*model.Metrics, error) {
    prof, err := s.GetProfile(ctx, username)
    if err != nil {
        return nil, err
    }

    tweets, err := s.GetTweets(ctx, username, sample)
    if err != nil {
        return nil, err
    }

    var likes, rts, replies int
    for _, t := range tweets {
        likes += t.LikeCount
        rts += t.RetweetCount
        replies += t.ReplyCount
    }
    n := len(tweets)
    if n == 0 {
        n = 1 // avoid div by zero
    }

    // compute posts per day based on created_at of earliest tweet
    var postsPerDay float64
    if len(tweets) > 1 {
        oldest, _ := time.Parse(time.RFC3339, tweets[len(tweets)-1].CreatedAt)
        newest, _ := time.Parse(time.RFC3339, tweets[0].CreatedAt)
        dur := newest.Sub(oldest).Hours() / 24
        if dur > 0 {
            postsPerDay = float64(len(tweets)) / dur
        }
    }

    m := &model.Metrics{
        UserID:      username,
        FFRatio:     float64(prof.Followers) / float64(max(1, prof.Following)),
        AvgLikes:    float64(likes) / float64(n),
        AvgRetweets: float64(rts) / float64(n),
        AvgReplies:  float64(replies) / float64(n),
        PostsPerDay: postsPerDay,
    }
    return m, nil
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
