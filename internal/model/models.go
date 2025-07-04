package model

// Публичный профиль пользователя Twitter
type User struct {
    ID            string `json:"id"`
    Name          string `json:"name"`
    Username      string `json:"username"`
    Description   string `json:"description,omitempty"`
    Followers     int    `json:"followers_count"`
    Following     int    `json:"following_count"`
    Listed        int    `json:"listed_count"`
    TweetCount    int    `json:"tweet_count"`
    ProfileImage  string `json:"profile_image_url,omitempty"`
    Verified      bool   `json:"verified"`
}

// Твит с основными метриками
type Tweet struct {
    ID            string `json:"id"`
    Text          string `json:"text"`
    LikeCount     int    `json:"like_count"`
    ReplyCount    int    `json:"reply_count"`
    RetweetCount  int    `json:"retweet_count"`
    QuoteCount    int    `json:"quote_count"`
    CreatedAt     string `json:"created_at"`
}

// Связь между двумя пользователями (граф)
type Edge struct {
    From string `json:"from"`
    To   string `json:"to"`
}

// Graph is the response model for /expand consisting of all discovered nodes
// and edges.
type Graph struct {
    Nodes []User `json:"nodes"`
    Edges []Edge `json:"edges"`
}

// Metrics aggregates engagement stats for a single user.
type Metrics struct {
    UserID         string  `json:"user_id"`
    FFRatio        float64 `json:"ff_ratio"` // followers / following
    AvgLikes       float64 `json:"avg_likes"`
    AvgRetweets    float64 `json:"avg_retweets"`
    AvgReplies     float64 `json:"avg_replies"`
    PostsPerDay    float64 `json:"posts_per_day"`
}
