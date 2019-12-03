package astitwitter

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type UserTimelineItem struct {
	CreatedAt UserTimelineTime     `json:"created_at"`
	Entities  UserTimelineEntities `json:"entities"`
	ID        int                  `json:"id"`
	Text      string               `json:"text"`
}

type UserTimelineTime time.Time

func (t *UserTimelineTime) UnmarshalText(b []byte) error {
	tt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", string(b))
	if err != nil {
		return err
	}
	*t = UserTimelineTime(tt)
	return nil
}

type UserTimelineEntities struct {
	URLs []UserTimelineURL `json:"urls"`
}

type UserTimelineURL struct {
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
}

type UserTimelineOptions struct {
	Count          *int
	ExcludeReplies *bool
	IncludeRts     *bool
	MaxID          *int
	ScreenName     string
	SinceID        *int
	TrimUser       *bool
	UserID         *int
}

func (c *Client) UserTimeline(o UserTimelineOptions) (is []UserTimelineItem, err error) {
	// Create query parameters
	qs := url.Values{}
	if o.Count != nil {
		qs.Set("count", strconv.Itoa(*o.Count))
	}
	if o.ExcludeReplies != nil {
		qs.Set("exclude_replies", fmt.Sprintf("%v", *o.ExcludeReplies))
	}
	if o.IncludeRts != nil {
		qs.Set("include_rts", fmt.Sprintf("%v", *o.IncludeRts))
	}
	if o.MaxID != nil {
		qs.Set("max_id", strconv.Itoa(*o.MaxID))
	}
	if o.ScreenName != "" {
		qs.Set("screen_name", o.ScreenName)
	}
	if o.SinceID != nil {
		qs.Set("since_id", strconv.Itoa(*o.SinceID))
	}
	if o.TrimUser != nil {
		qs.Set("trim_user", fmt.Sprintf("%v", *o.TrimUser))
	}
	if o.UserID != nil {
		qs.Set("user_id", strconv.Itoa(*o.UserID))
	}

	// Send
	if err = c.sendAuthenticated(http.MethodGet, "/1.1/statuses/user_timeline.json?"+qs.Encode(), nil, nil, &is); err != nil {
		err = errors.Wrap(err, "astitwitter: sending authenticated failed")
		return
	}
	return
}
