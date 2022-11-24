package vcnewsdaily

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/allokate-ai/environment"
	"github.com/allokate-ai/feeds/app/internal/event"
	"github.com/g8rswimmer/go-twitter/v2"

	"github.com/spf13/cobra"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func ExtractHashtags(text string) []string {
	r := regexp.MustCompile(`(#|＃)([a-zA-Z0-9_\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}-\x{0254}\x{0256}-\x{0257}\x{0300}-\x{036f}\x{1e00}-\x{1eff}\x{0400}-\x{04ff}\x{0500}-\x{0527}\x{2de0}-\x{2dff}\x{a640}-\x{a69f}\x{0591}-\x{05bf}\x{05c1}-\x{05c2}\x{05c4}-\x{05c5}\x{05d0}-\x{05ea}\x{05f0}-\x{05f4}\x{fb12}-\x{fb28}\x{fb2a}-\x{fb36}\x{fb38}-\x{fb3c}\x{fb40}-\x{fb41}\x{fb43}-\x{fb44}\x{fb46}-\x{fb4f}\x{0610}-\x{061a}\x{0620}-\x{065f}\x{066e}-\x{06d3}\x{06d5}-\x{06dc}\x{06de}-\x{06e8}\x{06ea}-\x{06ef}\x{06fa}-\x{06fc}\x{0750}-\x{077f}\x{08a2}-\x{08ac}\x{08e4}-\x{08fe}\x{fb50}-\x{fbb1}\x{fbd3}-\x{fd3d}\x{fd50}-\x{fd8f}\x{fd92}-\x{fdc7}\x{fdf0}-\x{fdfb}\x{fe70}-\x{fe74}\x{fe76}-\x{fefc}\x{200c}-\x{200c}\x{0e01}-\x{0e3a}\x{0e40}-\x{0e4e}\x{1100}-\x{11ff}\x{3130}-\x{3185}\x{a960}-\x{a97f}\x{ac00}-\x{d7af}\x{d7b0}-\x{d7ff}\x{ffa1}-\x{ffdc}\x{30a1}-\x{30fa}\x{30fc}-\x{30fe}\x{ff66}-\x{ff9f}\x{ff10}-\x{ff19}\x{ff21}-\x{ff3a}\x{ff41}-\x{ff5a}\x{3041}-\x{3096}\x{3099}-\x{309e}\x{3400}-\x{4dbf}\x{4e00}-\x{9fff}\x{20000}-\x{2a6df}\x{2a700}-\x{2b73f}\x{2b740}-\x{2b81f}\x{2f800}-\x{2fa1f}]*[a-z_\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}-\x{0254}\x{0256}-\x{0257}\x{0300}-\x{036f}\x{1e00}-\x{1eff}\x{0400}-\x{04ff}\x{0500}-\x{0527}\x{2de0}-\x{2dff}\x{a640}-\x{a69f}\x{0591}-\x{05bf}\x{05c1}-\x{05c2}\x{05c4}-\x{05c5}\x{05d0}-\x{05ea}\x{05f0}-\x{05f4}\x{fb12}-\x{fb28}\x{fb2a}-\x{fb36}\x{fb38}-\x{fb3c}\x{fb40}-\x{fb41}\x{fb43}-\x{fb44}\x{fb46}-\x{fb4f}\x{0610}-\x{061a}\x{0620}-\x{065f}\x{066e}-\x{06d3}\x{06d5}-\x{06dc}\x{06de}-\x{06e8}\x{06ea}-\x{06ef}\x{06fa}-\x{06fc}\x{0750}-\x{077f}\x{08a2}-\x{08ac}\x{08e4}-\x{08fe}\x{fb50}-\x{fbb1}\x{fbd3}-\x{fd3d}\x{fd50}-\x{fd8f}\x{fd92}-\x{fdc7}\x{fdf0}-\x{fdfb}\x{fe70}-\x{fe74}\x{fe76}-\x{fefc}\x{200c}-\x{200c}\x{0e01}-\x{0e3a}\x{0e40}-\x{0e4e}\x{1100}-\x{11ff}\x{3130}-\x{3185}\x{a960}-\x{a97f}\x{ac00}-\x{d7af}\x{d7b0}-\x{d7ff}\x{ffa1}-\x{ffdc}\x{30a1}-\x{30fa}\x{30fc}-\x{30fe}\x{ff66}-\x{ff9f}\x{ff10}-\x{ff19}\x{ff21}-\x{ff3a}\x{ff41}-\x{ff5a}\x{3041}-\x{3096}\x{3099}-\x{309e}\x{3400}-\x{4dbf}\x{4e00}-\x{9fff}\x{20000}-\x{2a6df}\x{2a700}-\x{2b73f}\x{2b740}-\x{2b81f}\x{2f800}-\x{2fa1f}][a-z0-9_\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}-\x{0254}\x{0256}-\x{0257}\x{0300}-\x{036f}\x{1e00}-\x{1eff}\x{0400}-\x{04ff}\x{0500}-\x{0527}\x{2de0}-\x{2dff}\x{a640}-\x{a69f}\x{0591}-\x{05bf}\x{05c1}-\x{05c2}\x{05c4}-\x{05c5}\x{05d0}-\x{05ea}\x{05f0}-\x{05f4}\x{fb12}-\x{fb28}\x{fb2a}-\x{fb36}\x{fb38}-\x{fb3c}\x{fb40}-\x{fb41}\x{fb43}-\x{fb44}\x{fb46}-\x{fb4f}\x{0610}-\x{061a}\x{0620}-\x{065f}\x{066e}-\x{06d3}\x{06d5}-\x{06dc}\x{06de}-\x{06e8}\x{06ea}-\x{06ef}\x{06fa}-\x{06fc}\x{0750}-\x{077f}\x{08a2}-\x{08ac}\x{08e4}-\x{08fe}\x{fb50}-\x{fbb1}\x{fbd3}-\x{fd3d}\x{fd50}-\x{fd8f}\x{fd92}-\x{fdc7}\x{fdf0}-\x{fdfb}\x{fe70}-\x{fe74}\x{fe76}-\x{fefc}\x{200c}-\x{200c}\x{0e01}-\x{0e3a}\x{0e40}-\x{0e4e}\x{1100}-\x{11ff}\x{3130}-\x{3185}\x{a960}-\x{a97f}\x{ac00}-\x{d7af}\x{d7b0}-\x{d7ff}\x{ffa1}-\x{ffdc}\x{30a1}-\x{30fa}\x{30fc}-\x{30fe}\x{ff66}-\x{ff9f}\x{ff10}-\x{ff19}\x{ff21}-\x{ff3a}\x{ff41}-\x{ff5a}\x{3041}-\x{3096}\x{3099}-\x{309e}\x{3400}-\x{4dbf}\x{4e00}-\x{9fff}\x{20000}-\x{2a6df}\x{2a700}-\x{2b73f}\x{2b740}-\x{2b81f}\x{2f800}-\x{2fa1f}]*)`)
	hashtags := r.FindAllString(text, -1)
	if hashtags == nil {
		return []string{}
	}
	return hashtags
}

func ExtractUsernames(text string) []string {
	r := regexp.MustCompile(`^@(\w){1,15}$`)
	usernames := r.FindAllString(text, -1)
	if usernames == nil {
		return []string{}
	}
	return usernames
}

func BuildTweetEvent(item *twitter.TweetDictionary) event.Tweet {
	date, err := time.Parse(time.RFC3339, item.Tweet.CreatedAt)
	if err != nil {
		log.Panic("Failed to parse date")
	}

	return event.Tweet{
		Name:     item.Author.Name,
		UserName: item.Author.UserName,
		Date:     date,
		Content:  strings.TrimLeft(item.Tweet.Text, "."),
		Mentions: ExtractUsernames(item.Tweet.Text),
		Hashtags: ExtractHashtags(item.Tweet.Text),
	}
}

// rootCmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use: "vcnewsdaily",
	// Short: "A tool various data feeds",
	// Long:  `This is a tool used to scrape various RSS news feeds, API, and other sources for data used by Allokate.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		userId := 323853589

		client := &twitter.Client{
			Authorizer: authorize{
				Token: environment.MustGet("TWITTER_BEARER_TOKEN"),
			},
			Client: http.DefaultClient,
			Host:   "https://api.twitter.com",
		}

		startOfDay := time.Now().UTC().Round(24 * time.Hour)

		opts := twitter.UserTweetTimelineOpts{
			StartTime:   startOfDay,
			TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldAuthorID, twitter.TweetFieldConversationID, twitter.TweetFieldPublicMetrics, twitter.TweetFieldContextAnnotations},
			UserFields:  []twitter.UserField{twitter.UserFieldUserName},
			Expansions:  []twitter.Expansion{twitter.ExpansionAuthorID},
			MaxResults:  100,
		}

		timeline, err := client.UserTweetTimeline(context.Background(), strconv.Itoa(userId), opts)
		if err != nil {
			log.Panicf("user tweet timeline error: %v", err)
		}

		dictionaries := timeline.Raw.TweetDictionaries()

		for _, item := range dictionaries {
			tweet := BuildTweetEvent(item)

			// Send it!!
			if _, err := event.EmitTweetEvent("feeds.twitter.vcnewsdaily", tweet); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("%s: %s", tweet.Date.Format(time.RFC3339), tweet.Content)
			}
		}

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feeds.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
