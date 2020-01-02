package twitter

import (
	"flag"
	"os"
	"sort"
	"time"

	"github.com/caiotedim/crawler-tags/config"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/golang/glog"
)

type (
	// Client return Twitter Client
	Client struct {
		twitterClient *twitter.Client
	}

	// Hashtags struct
	Hashtags struct {
		Hashtag string `json:"hashtag"`
		Tweets  Tweets
	}

	// Tweets return a tweet
	Tweets struct {
		Tweet []Tweet `json:"tweets"`
	}

	// Tweet return message, user, created_by, followers, lang, country
	Tweet struct {
		Msg       string    `json:"message"`
		User      string    `json:"user"`
		Created   time.Time `json:"created"`
		Followers int       `json:"followers_count"`
		Lang      string    `json:"language"`
	}

	// Followers struct to summarize ByFollowers
	Followers struct {
		Followers int    `json:"followers_count"`
		User      string `json:"user"`
	}

	// ByFollowers implements sort.Interface based on the Followers field
	ByFollowers []Followers

	// Hourly struct to summarize
	Hourly struct {
		Total int
		Hour  int
	}

	//HashtagCount struct to summarize
	HashtagCount struct {
		Hashtag   string
		LangCount []LangCount
	}

	//LangCount struct to summarize
	LangCount struct {
		Lang  string
		Total int
	}
)

var hashtags = []string{"openbanking", "apifirst", "devops", "cloudfirst", "microservices", "apigateway", "oauth", "swagger", "raml", "openapis"}

// NewTwitterClient return a configured client to access Twitter API
func NewTwitterClient(c *config.Config) Client {
	client := new(Client)
	config := oauth1.NewConfig(c.Twitter.ConsumerAPIKey, c.Twitter.ConsumerAPISecret)
	token := oauth1.NewToken(c.Twitter.AccessTokenKey, c.Twitter.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client.twitterClient = twitter.NewClient(httpClient)
	return *client
}

// SearchTweetByHashTag Search tweets using some hashtag
func (client Client) SearchTweetByHashTag(hashtag string) (Tweets, error) {

	var tweet Tweet
	var tweets Tweets
	searches, _, err := client.twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query: hashtag,
	})
	if err != nil {
		glog.Errorf("Error to search on Twitter API: %v", err)
		return tweets, err
	}

	for _, v := range searches.Statuses {
		tweet.Msg = v.Text
		tweet.User = v.User.Name
		tweet.Created, err = v.CreatedAtTime()
		if err != nil {
			glog.Errorf("Error to get CreatedTime on a tweet: %v", err)
		}
		tweet.Followers = v.User.FollowersCount
		tweet.Lang = v.Lang
		tweets.Tweet = append(tweets.Tweet, tweet)
	}

	return tweets, nil

}

// LookupHashtags ...
func LookupHashtags(c *config.Config) ([]Hashtags, error) {
	//c := config.NewConfig()
	client := NewTwitterClient(c)

	var tweetsHashtags []Hashtags
	var tweetHashtag Hashtags
	var err error
	for _, value := range hashtags {
		tweetHashtag.Hashtag = value
		tweetHashtag.Tweets, err = client.SearchTweetByHashTag(value)
		if err != nil {
			glog.Errorf("Error on search tweet using hashtag: %s -> %v", value, err)
			return tweetsHashtags, err
		}
		tweetsHashtags = append(tweetsHashtags, tweetHashtag)
	}
	//var summarized []Hashtags
	//summarizeData(tweetsHashtags, summarized)
	return tweetsHashtags, nil
}

/*func summarizeData(tweetsHashtags, summarized []Hashtags) []Hashtags {
	topFollowers(tweetsHashtags)
	return tweetsHashtags
}*/

func (f ByFollowers) Len() int           { return len(f) }
func (f ByFollowers) Less(i, j int) bool { return f[i].Followers < f[j].Followers }
func (f ByFollowers) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// TopFollowers return 5 Top followers
func TopFollowers(tweetsHashtags []Hashtags) []Followers {
	var followers []Followers
	var follow Followers
	for _, v := range tweetsHashtags {
		for _, i := range v.Tweets.Tweet {
			follow.Followers = i.Followers
			follow.User = i.User
			followers = append(followers, follow)
		}
	}
	sort.Sort(sort.Reverse(ByFollowers(followers)))
	return followers[:5]
}

/*func (h ByHourly) Len() int           { return len(h) }
func (h ByHourly) Less(i, j int) bool { return h[i].Hour < h[j].Hour }
func (h ByHourly) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func topHourlyPost(tweetsHashtags []Hashtags) []Hourly {
	var hourly []Hourly
	var hour Hourly
	for _, v := range tweetsHashtags {
		for _, i := range v.Tweets.Tweet {
			hour.Hour = i.Created.Hour()
			hourly = append(hourly, hour)
		}
	}
	sort.Sort(ByHourly(hourly))
	return hourly
}*/

// CountByHour data summarized by hour
func CountByHour(tweetsHashtags []Hashtags) []Hourly {
	var hourly []Hourly
	var hour Hourly

	hourCount := make(map[int]int)
	for _, i := range tweetsHashtags {
		for _, v := range i.Tweets.Tweet {
			hourCount[v.Created.Hour()]++
		}

	}

	for k, v := range hourCount {
		hour.Hour = k
		hour.Total = v
		hourly = append(hourly, hour)
	}
	return hourly
}

// CountLang total of posts by lang per hashtag
func CountLang(tweetsHashtags []Hashtags) []HashtagCount {
	var lang []HashtagCount

	for _, i := range tweetsHashtags {
		var langAux []LangCount
		var hashtagAux HashtagCount
		langCount := make(map[string]int)
		for _, v := range i.Tweets.Tweet {
			langCount[v.Lang]++
		}
		hashtagAux.Hashtag = i.Hashtag
		for k, v := range langCount {
			var aux LangCount
			aux.Lang = k
			aux.Total = v
			langAux = append(langAux, aux)
		}
		hashtagAux.LangCount = langAux
		lang = append(lang, hashtagAux)
	}

	return lang
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "false")
	flag.Set("log_dir", "/tmp/crawler-tags")
	flag.Parse()
}
