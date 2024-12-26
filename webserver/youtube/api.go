// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package youtube

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// feed is the object that holds a incomming notification feed from
// youtube. This could be a new video (upload/publish) or an update
// of an existing one.
type feed struct {
	Video     feedVideo   `xml:"entry>videoId"`
	Channel   feedChannel `xml:"entry>channelId"`
	Published time.Time   `xml:"entry>published"`
	Updated   time.Time   `xml:"entry>updated"`
}

// feedVideo is part of the feed xml struct and contains the videoId
// field.
type feedVideo struct {
	XMLName xml.Name `xml:"videoId"`
	ID      string   `xml:",chardata"`
}

// feedChannel is part of the xml feed struct and contains the
// channelId field.
type feedChannel struct {
	XMLName xml.Name `xml:"channelId"`
	ID      string   `xml:",chardata"`
}

// HandleGet is the HTTP/GET handler for the YouTube PubSubHubBub
// endpoint.
//
// It is used to accept new webhook subscriptions for YouTube video
// news feed, like publish a new video or editing an existing one.
func HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	topic := r.FormValue("hub.topic")
	challenge := r.FormValue("hub.challenge")
	mode := r.FormValue("hub.mode")

	if topic == "" || challenge == "" || mode == "" {
		log.Println("Missing at least one of topic, challenge, mode in query parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok, _ := regexp.MatchString("(?:un)?subscribe", mode); !ok {
		log.Printf("Unsupported mode '%s'", mode)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	topicURL, err := url.Parse(topic)
	if err != nil {
		log.Printf("Error on parse topic url '%s': %v\n", topic, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// only accept youtube video feed
	if topicURL.Host != "www.youtube.com" {
		log.Printf("Topic host is not youtube: %s\n", topicURL.Host)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if topicURL.Path != "/xml/feeds/videos.xml" {
		log.Printf("Topic path is not for videos: %s\n", topicURL.Path)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	channelID := topicURL.Query().Get("channel_id")

	// check for valid actions
	if !subscribtions[channelID] && mode == "subscribe" {
		log.Printf("Requested subscription for unknown channel: %s\n", channelID)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if subscribtions[channelID] && mode == "unsubscribe" {
		log.Printf("Requested unsubscribe for used channel: %s\n", channelID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(challenge))
	log.Printf("Accepted '%s' from %s for channel %s\n", mode, topicURL.Host, channelID)
}

// HandlePost is the HTTP/POST handler for the YouTube PubSubHubBub
// endpoint.
//
// It is used to handle a notification feed comming from a newly
// published video of a subscribed channel.
func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// only accept atom feed
	content := r.Header.Get("Content-Type")
	if content != "application/atom+xml" {
		log.Printf("Content-Type '%s' not supported\n", content)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	feed := feed{}
	err = xml.Unmarshal(buf, &feed)
	if err != nil {
		log.Printf("Error on parse XML body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// need yt namespace i.e. <yt:videoId>1a2b3c4d</yt:videoId> as
	// well as the xmlns:yt attribute to be set in a parent xmlns tag.
	// The namespace needs to be a valid url with "www.youtube.com"
	// as host.
	xmlnsVideo, err := url.Parse(feed.Video.XMLName.Space)
	if err != nil {
		log.Printf("Error on video namespace url '%s': %v\n", feed.Video.XMLName.Space, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if xmlnsVideo.Host != "www.youtube.com" {
		log.Printf("xml video namespace is not from youtube '%s': %v\n", xmlnsVideo, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if feed.Video.XMLName.Space != feed.Channel.XMLName.Space {
		log.Printf("xml channel namespace ('%s') is not the same as video namespace ('%s')\n", feed.Channel.XMLName.Space, feed.Video.XMLName.Space)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Content will be checked and could be ignored on mismatch"))

	go func() {
		if age := feed.Updated.Sub(feed.Published); age > time.Minute*15 {
			return
		}
		video, ok := checkVideo(feed.Video.ID, feed.Channel.ID)
		if !ok {
			return
		}
		dcHandler(dcSession, video)
	}()
}
