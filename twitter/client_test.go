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
	os.Setenv("DB_HOST", "127.0.0.1")
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

func TestCountLang(t *testing.T) {
	os.Setenv("CONSUMER_API_KEY", "U3tuMr9txi4jtqtaZFHaKC2RO")
	os.Setenv("CONSUMER_API_SECRET", "sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc")
	os.Setenv("ACCESS_TOKEN_KEY", "89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy")
	os.Setenv("ACCESS_TOKEN_SECRET", "htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j")
	os.Setenv("DB_HOST", "127.0.0.1")
	c := config.NewConfig()
	tweets, err := LookupHashtags(c)
	if err != nil {
		t.Errorf("Error to Lookup on twitter: %v", err)
	}

	for k, i := range tweets {
		if len(tweets[k].Tweets.Tweet) != 0 {
			t.Logf("hashtag: %s tweets: %d", i.Hashtag, len(tweets[k].Tweets.Tweet))
		} else {
			t.Errorf("No tweets with this hashtag: %s", i.Hashtag)
		}
	}

	hashtags := CountLang(tweets)
	for _, i := range hashtags {
		for _, v := range i.LangCount {
			t.Logf("Hashtag: %s - Lang: %s - Total: %d", i.Hashtag, v.Lang, v.Total)
		}
	}

}

func TestTopFollowers(t *testing.T) {
	os.Setenv("CONSUMER_API_KEY", "U3tuMr9txi4jtqtaZFHaKC2RO")
	os.Setenv("CONSUMER_API_SECRET", "sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc")
	os.Setenv("ACCESS_TOKEN_KEY", "89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy")
	os.Setenv("ACCESS_TOKEN_SECRET", "htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j")
	os.Setenv("DB_HOST", "127.0.0.1")
	c := config.NewConfig()
	tweets, err := LookupHashtags(c)
	if err != nil {
		t.Errorf("Error to Lookup on twitter: %v", err)
	}

	for k, i := range tweets {
		if len(tweets[k].Tweets.Tweet) != 0 {
			t.Logf("hashtag: %s tweets: %d", i.Hashtag, len(tweets[k].Tweets.Tweet))
		} else {
			t.Errorf("No tweets with this hashtag: %s", i.Hashtag)
		}
	}

	sorted := TopFollowers(tweets)
	for _, i := range sorted {
		t.Logf("user: %v followers: %d", i.User, i.Followers)
	}
}

func TestCountByHour(t *testing.T) {
	os.Setenv("CONSUMER_API_KEY", "U3tuMr9txi4jtqtaZFHaKC2RO")
	os.Setenv("CONSUMER_API_SECRET", "sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc")
	os.Setenv("ACCESS_TOKEN_KEY", "89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy")
	os.Setenv("ACCESS_TOKEN_SECRET", "htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j")
	os.Setenv("DB_HOST", "127.0.0.1")
	c := config.NewConfig()
	tweets, err := LookupHashtags(c)
	if err != nil {
		t.Errorf("Error to Lookup on twitter: %v", err)
	}

	for k, i := range tweets {
		if len(tweets[k].Tweets.Tweet) != 0 {
			t.Logf("hashtag: %s tweets: %d", i.Hashtag, len(tweets[k].Tweets.Tweet))
		} else {
			t.Errorf("No tweets with this hashtag: %s", i.Hashtag)
		}
	}

	hourly := CountByHour(tweets)
	for _, i := range hourly {
		t.Logf("Hour: %d - Total: %d", i.Hour, i.Total)
	}
}
