package twitter

import (
	"os"
	"testing"

	"github.com/caiotedim/crawler-tags/config"
)

func TestSearchTweetByHashTag(t *testing.T) {
	os.Setenv("CONSUMER_API_KEY", "U3tuMr9txi4jtqtaZFHaKC2RO")
	os.Setenv("CONSUMER_API_SECRET", "sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc")
	os.Setenv("ACCESS_TOKEN_KEY", "89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy")
	os.Setenv("ACCESS_TOKEN_SECRET", "htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j")
	c := config.NewConfig()
	client := NewTwitterClient(c)
	hashtag := "devops"
	tweets, err := client.SearchTweetByHashTag(hashtag)
	if err != nil {
		t.Errorf("Error on search tweet using hashtag: %v", err)
		return
	}

	if tweets.Tweet[0].Msg == "" {
		t.Errorf("NOK")
	}

	if len(tweets.Tweet) != 0 {
		t.Logf("Tweets: %d", len(tweets.Tweet))
	} else {
		t.Errorf("No tweets with this hashtag: %s", hashtag)
	}

}

func TestLookupHashtags(t *testing.T) {
	os.Setenv("CONSUMER_API_KEY", "U3tuMr9txi4jtqtaZFHaKC2RO")
	os.Setenv("CONSUMER_API_SECRET", "sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc")
	os.Setenv("ACCESS_TOKEN_KEY", "89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy")
	os.Setenv("ACCESS_TOKEN_SECRET", "htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j")
	hashtag := ""
	tweets := LookupHashtags(hashtag)

	for k, i := range tweets {
		if len(tweets[k].Tweets.Tweet) != 0 {
			t.Logf("hashtag: %s tweets: %d", i.Hashtag, len(tweets[k].Tweets.Tweet))
		} else {
			t.Errorf("No tweets with this hashtag: %s", i.Hashtag)
		}
	}

	/*sorted := topFollowers(tweets)
	for _, i := range sorted {
		t.Errorf("user: %v followers: %d", i.User, i.Followers)
	}

	hourly := countByHour(tweets)
	for _, i := range hourly {
		t.Errorf("Hour: %d - Total: %d", i.Hour, i.Total)
	}*/

	hashtags := CountLang(tweets)
	for _, i := range hashtags {
		for _, v := range i.LangCount {
			t.Errorf("Hashtag: %s - Lang: %s - Total: %d", i.Hashtag, v.Lang, v.Total)
		}
	}

}
