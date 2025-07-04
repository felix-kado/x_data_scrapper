package twitter

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "time"

    "github.com/google/go-querystring/query"
)

// Клиент для Twitter API v2 (только нужные методы)
type Client struct {
    http *http.Client
    host string
    key  string // bearer token
}

func NewClient(httpClient *http.Client, bearerToken string) *Client {
    if httpClient == nil {
        httpClient = &http.Client{Timeout: 15 * time.Second}
    }

    return &Client{
        http: httpClient,
        host: "https://api.twitter.com/2",
        key:  bearerToken,
    }
}

// Выполнить HTTP-запрос к Twitter API и распарсить ответ
func (c *Client) do(ctx context.Context, method, path string, params url.Values, dest interface{}) error {
    u := fmt.Sprintf("%s%s", c.host, path)

    if len(params) > 0 {
        u += "?" + params.Encode()
    }

    req, err := http.NewRequestWithContext(ctx, method, u, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+c.key)
    req.Header.Set("User-Agent", "x_data_scrapper/1.0")

    res, err := c.http.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode >= 400 {
        return fmt.Errorf("twitter api error: %s", res.Status)
    }

    return json.NewDecoder(res.Body).Decode(dest)
}

//--------------------------------------------------------------------
// Domain-specific helpers
//--------------------------------------------------------------------

type userResponse struct {
    Data struct {
        ID              string `json:"id"`
        Name            string `json:"name"`
        Username        string `json:"username"`
        Description     string `json:"description"`
        ProfileImageURL string `json:"profile_image_url"`
        Verified        bool   `json:"verified"`
        PublicMetrics   struct {
            Followers  int `json:"followers_count"`
            Following  int `json:"following_count"`
            Listed     int `json:"listed_count"`
            TweetCount int `json:"tweet_count"`
        } `json:"public_metrics"`
    } `json:"data"`
}

type tweetsResponse struct {
    Data []struct {
        ID           string `json:"id"`
        Text         string `json:"text"`
        CreatedAt    string `json:"created_at"`
        PublicMetric struct {
            LikeCount    int `json:"like_count"`
            RetweetCount int `json:"retweet_count"`
            ReplyCount   int `json:"reply_count"`
            QuoteCount   int `json:"quote_count"`
        } `json:"public_metrics"`
    } `json:"data"`
    Meta struct {
        NextToken string `json:"next_token"`
        ResultCnt int    `json:"result_count"`
    } `json:"meta"`
}

// Получить пользователя по id (внутренний метод)
// func (c *Client) getUserByID(ctx context.Context, id string) (*userResponse, error) {
//     params := url.Values{}
//     params.Set("user.fields", "public_metrics,description,profile_image_url,verified")
//     var res userResponse
//     if err := c.do(ctx, http.MethodGet, "/users/"+id, params, &res); err != nil {
//         return nil, err
//     }
//     return &res, nil
// }

// Получить пользователя по username
func (c *Client) GetUserByUsername(ctx context.Context, username string) (*userResponse, error) {
    params := url.Values{}
    params.Set("user.fields", "public_metrics,description,profile_image_url,verified")
    var res userResponse
    if err := c.do(ctx, http.MethodGet, "/users/by/username/"+username, params, &res); err != nil {
        return nil, err
    }
    return &res, nil
}

// Получить твиты пользователя по id (до лимита)
func (c *Client) GetTweets(ctx context.Context, id string, limit int) ([]tweetsResponse, error) {
    // The official API returns max 100 per page. We'll page until limit or no next_token
    const pageSize = 100

    params := struct {
        MaxResults   int    `url:"max_results,omitempty"`
        Pagination   string `url:"pagination_token,omitempty"`
        TweetFields  string `url:"tweet.fields"`
        Expansions   string `url:"expansions,omitempty"`
        MediaFields  string `url:"media.fields,omitempty"`
        StartTime    string `url:"start_time,omitempty"`
        EndTime      string `url:"end_time,omitempty"`
        SinceID      string `url:"since_id,omitempty"`
        UntilID      string `url:"until_id,omitempty"`
        Exclude      string `url:"exclude,omitempty"`
        PlaceFields  string `url:"place.fields,omitempty"`
        PollFields   string `url:"poll.fields,omitempty"`
        UserFields   string `url:"user.fields,omitempty"`
    }{
        MaxResults:  pageSize,
        TweetFields: "created_at,public_metrics",
    }

    collected := 0
    var pages []tweetsResponse
    next := ""

    for collected < limit || limit == 0 {
        params.Pagination = next
        v, _ := query.Values(params)
        var resp tweetsResponse
        if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/users/%s/tweets", id), v, &resp); err != nil {
            return nil, err
        }
        pages = append(pages, resp)
        collected += len(resp.Data)
        if resp.Meta.NextToken == "" || (limit != 0 && collected >= limit) {
            break
        }
        next = resp.Meta.NextToken
    }

    return pages, nil
}
