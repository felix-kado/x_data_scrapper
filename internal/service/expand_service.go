package service

import (
    "context"

    "github.com/felix-kado/x_data_scrapper/internal/model"
    "github.com/felix-kado/x_data_scrapper/internal/twitter"
)

// Направление обхода графа
type ExpandDirection int

const (
    FollowersDir ExpandDirection = iota // Направление обхода: подписчики
    FollowingDir // Направление обхода: подписки
)

// Параметры для построения графа
type ExpandParams struct {
    SeedIDs   []string        `json:"seed_ids"` // Идентификаторы стартовых пользователей
    Depth     int             `json:"depth"` // Глубина построения графа
    Direction ExpandDirection `json:"direction"` // Направление обхода графа (0 = followers, 1 = following)
    CollectTweets  bool `json:"collect_tweets"` // Собирать ли твиты
    CollectMetrics bool `json:"collect_metrics"` // Собирать ли метрики
}

// Сервис построения подграфа Twitter
type ExpandService struct {
    tw *twitter.Client
    userSvc *UserService
}

func NewExpandService(tw *twitter.Client, user *UserService) *ExpandService {
    return &ExpandService{tw: tw, userSvc: user}
}

func (s *ExpandService) Expand(ctx context.Context, p ExpandParams) (*model.Graph, error) {
    // For brevity we only collect nodes (profiles) without edges due to limited API access (followers/ids require elevated).
    graph := &model.Graph{}

    visited := map[string]bool{}
    queue := []string{}
    queue = append(queue, p.SeedIDs...)

    depth := 0
    for len(queue) > 0 && depth <= p.Depth {
        nextQ := []string{}
        for _, id := range queue {
            if visited[id] {
                continue
            }
            visited[id] = true

            prof, err := s.userSvc.GetProfile(ctx, id)
            if err != nil {
                return nil, err
            }
            graph.Nodes = append(graph.Nodes, *prof)

            // NOTE: follower/following expansions require Twitter API endpoints that may be unavailable for basic access.
            // This is left as TODO – here we just stop after seed level.
        }
        queue = nextQ
        depth++
    }

    return graph, nil
}
