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
	"encoding/json"
	"fmt"
	"io"
	logger "log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var log logger.Logger = *logger.New(logger.Writer(), "[WebYouTube] ", logger.LstdFlags|logger.Lmsgprefix)

type listResponse struct {
	Item []item `json:"items,omitempty"`
}

type item struct {
	ID    string `json:"id,omitempty"`
	Video Video  `json:"snippet,omitempty"`
}

// Video represents a YouTube video
type Video struct {
	// The video ID
	ID string `json:"id,omitempty"`
	// Time publication
	Time time.Time `json:"publishedAt,omitempty"`
	// The ID of the corresponding YouTube channel
	ChannelID string `json:"channelId,omitempty"`
	// The video title
	Title string `json:"title,omitempty"`
	// The video description
	Description string `json:"description,omitempty"`
	// Available thumbnails
	Thumbnails map[string]Thumbnail `json:"thumbnails,omitempty"`
	// The name of the corresponding YouTube channel
	Channel string `json:"channelTitle,omitempty"`
	// Tags of the video
	Tags []string `json:"tags,omitempty"`
	// The video category
	Category Category `json:"categoryId,omitempty"`
}

// Thumbnail is the presentation image of a YouTube video. This
// struct is part of the Video struct.
type Thumbnail struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Category represents a YouTube video category. This is part of the
// Video struct.
type Category string

// Category types
const (
	FilmAnimation       Category = "1"
	AutosVehicles       Category = "2"
	Music               Category = "10"
	PetsAnimals         Category = "15"
	Sports              Category = "17"
	ShortMovies         Category = "18"
	TravelEvents        Category = "19"
	Gaming              Category = "20"
	Videoblogging       Category = "21"
	PeopleBlogs         Category = "22"
	Comedy              Category = "23"
	Entertainments      Category = "24"
	NewsPolitics        Category = "25"
	HowtoStyles         Category = "26"
	Education           Category = "27"
	ScienceTechnology   Category = "28"
	NonprofitsActivisim Category = "29"
	Movies              Category = "30"
	AnimeAnimation      Category = "31"
	ActionAdventure     Category = "32"
	Classics            Category = "33"
	Comedy2             Category = "34"
	Documentary         Category = "35"
	Drama               Category = "36"
	Family              Category = "37"
	Foreign             Category = "38"
	Horro               Category = "39"
	SciFiFantasy        Category = "40"
	Thriller            Category = "41"
	Shorts              Category = "42"
	Shows               Category = "43"
	Trailers            Category = "44"
)

func (c Category) String() string {
	switch c {
	case FilmAnimation:
		return "Film & Animation"
	case AutosVehicles:
		return "Autos & Vehicles"
	case Music:
		return "Music"
	case PetsAnimals:
		return "Pets & Animals"
	case Sports:
		return "Sports"
	case ShortMovies:
		return "Short Movies"
	case TravelEvents:
		return "Travel & Events"
	case Gaming:
		return "Gaming"
	case Videoblogging:
		return "Videoblogging"
	case PeopleBlogs:
		return "People & Blogs"
	case Comedy:
		return "Comedy"
	case Entertainments:
		return "Entertainments"
	case NewsPolitics:
		return "News & Politics"
	case HowtoStyles:
		return "HowtoStyles"
	case Education:
		return "Education"
	case ScienceTechnology:
		return "Science & Technology"
	case NonprofitsActivisim:
		return "Nonprofits &Activisim"
	case Movies:
		return "Movies"
	case AnimeAnimation:
		return "Anime/Animation"
	case ActionAdventure:
		return "Action/Adventure"
	case Classics:
		return "Classics"
	case Comedy2:
		return "Comedy"
	case Documentary:
		return "Documentary"
	case Drama:
		return "Drama"
	case Family:
		return "Family"
	case Foreign:
		return "Foreign"
	case Horro:
		return "Horro"
	case SciFiFantasy:
		return "Sci-Fi/Fantasy"
	case Thriller:
		return "Thriller"
	case Shorts:
		return "Shorts"
	case Shows:
		return "Shows"
	case Trailers:
		return "Trailers"
	default:
		return "<Unknown Category>"
	}
}

const (
	youtubeAPIBaseURL string = "https://youtube.googleapis.com/youtube/v3"
	// CALLBACKPATH is the endpoint that youtube will call to send events.
	CALLBACKPATH string = "/api/yt_pubsubhubbub/"
	// CALLBACKURL is the full URL to the CALLBACKPATH.
	CALLBACKURL string = "https://webhook.cake4everyone.de" + CALLBACKPATH
)

// checkVideo checks if a video really is from the provided channel
// by making an API call back to youtube an trying to get the video
// by id. It also returns some other video details.
func checkVideo(id string, channelID string) (v *Video, ok bool) {
	if dcSession == nil || dcHandler == nil {
		log.Printf("Error: got video event '%s', but discord is not set up!", id)
		return nil, false
	}

	// check if the channel is actually in the subscription list
	if !subscribtions[channelID] {
		log.Println("Channel not subscribed to")
		return nil, false
	}

	// get video by by from youtube
	url := fmt.Sprintf("%s?part=%s&id=%s&key=%s",
		youtubeAPIBaseURL+"/videos",
		"snippet",
		id,
		viper.GetString("google.apiKey"),
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error calling YouTube API for video '%s': %v\n", id, err)
		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error calling YouTube API for video '%s': YouTube responded with %d but expected 200", id, resp.StatusCode)
		return nil, false
	}

	// parse youtube video
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing YouTube API response for video '%s': %v", id, err)
		return nil, false
	}

	listResponse := &listResponse{}
	err = json.Unmarshal(data, listResponse)
	if err != nil {
		log.Printf("Error parsing YouTube API response for video '%s': %v", id, err)
		return nil, false
	}

	if len(listResponse.Item) == 0 {
		log.Printf("Got no video for id %s", id)
		return nil, false
	}

	v = &listResponse.Item[0].Video
	v.ID = listResponse.Item[0].ID

	// check against expected IDs
	if v.ID != id {
		log.Printf("Requested video does not match with online video: requested id '%s', got '%s'", id, v.ID)
		return nil, false
	}
	if v.ChannelID != channelID {
		log.Printf("Requested video does not match with online video: requested channel '%s', video is from '%s'", channelID, v.ChannelID)
		return nil, false
	}

	return v, true
}
