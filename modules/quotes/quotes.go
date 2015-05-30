package quotes

import (
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/jamesclonk-io/stdlib/web"
)

var quotes = []string{
	`valar morghulis`,
	`valar dohaeris`,
	`winter is coming`,
	`the laughing man`,
	`I thought what I'd do was<br/>I'd pretend I was one of those deaf-mutes`,
	`the cosmos is all that is or ever was or ever will be`,
	`imagination will often carry us to worlds that never were,<br/>but without it we go nowhere`,
	`if you wish to make an apple pie from scratch,<br/>you must first invent the universe`,
	`matter is composed mainly of nothing`,
	`all mass is interaction`,
	`..and you will find someday that, after all,<br/>it isn't as horrible as it looks`,
	`ask me when it's all over..`,
	`the universe seems neither benign nor hostile,<br/>merely indifferent`,
	`haste makes waste`,
	`take it with a grain of salt`,
	`sell a man a fish, he eats for a day,<br/>teach a man to fish, he eats for his lifetime`,
	`a journey of a thousand miles<br/>starts with a single step`,
	`holy cow!`,
	`cogito ergo sum`,
	`great spirits have always encountered<br/>violent opposition from mediocre minds`,
	`damned if you do,<br/>damned if you don't`,
	`ask me no questions,<br/>I'll tell you no lies`,
	`we are star stuff<br/>which has taken its destiny into its own hands`,
	`holy heart failure, batman!`,
	`there are only two hard things in computer science:<br/>cache invalidation and naming things`} // don't make quotes longer than this one

type Quotes struct {
	frontend *web.Frontend
}

func init() {
	rand.Seed(time.Now().Unix())
}

func getRandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}

func NewQuoteMiddleware(frontend *web.Frontend) *Quotes {
	return &Quotes{frontend}
}

func (c *Quotes) ServeHTTP(http.ResponseWriter, *http.Request) {
	// quote middleware uses pagemaster data field to store quotes in, so that templates can read it from page.Data
	c.frontend.PageMaster.Data = template.HTML(getRandomQuote())
}
