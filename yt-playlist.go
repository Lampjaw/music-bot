package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var (
	youtubeBaseURL            = "https://www.youtube.com/"
	regexpInitialPlaylistData = regexp.MustCompile(`\["ytInitialData"\] = (.+);`)
)

func (c *Client) GetPlaylistInfoFromID(id string) (*PlaylistInfo, error) {
	body, err := c.httpGetAndCheckResponseReadBody(youtubeBaseURL + "playlist?list=" + id)

	if err != nil {
		return nil, err
	}
	return getPlaylistInfoFromHTML(body)
}

func getPlaylistInfoFromHTML(html []byte) (*PlaylistInfo, error) {
	if matches := regexpInitialPlaylistData.FindSubmatch(html); len(matches) > 0 {
		data := initialPlaylistData{}

		if err := json.Unmarshal(matches[1], &data); err != nil {
			return nil, err
		}

		info := data.Microformat.MicroformatDataRenderer

		playlistInfo := &PlaylistInfo{
			Title:       info.Title,
			Description: info.Description,
		}

		if len(info.Thumbnail.Thumbnails) > 0 {
			playlistInfo.ThumbnailURL = info.Thumbnail.Thumbnails[0].URL
		}

		tabSection := data.Contents.TwoColumnBrowseResultsRenderer.Tabs[0]
		sectionList := tabSection.TabRenderer.Content.SectionListRenderer.Contents[0]
		itemSection := sectionList.ItemSectionRenderer.Contents[0]
		playlistItems := itemSection.PlaylistVideoListRenderer.Contents

		playlist := make([]*PlaylistItem, 0)

		for _, item := range playlistItems {
			vid := item.PlaylistVideoRenderer
			p := &PlaylistItem{
				VideoID:    vid.VideoID,
				Title:      vid.Title.SimpleText,
				Duration:   time.Duration(int64(vid.LengthSeconds) * int64(time.Second)),
				IsPlayable: vid.IsPlayable,
			}
			if len(vid.Thumbnail.Thumbnails) > 0 {
				p.ThumbnailURL = vid.Thumbnail.Thumbnails[0].URL
			}
			playlist = append(playlist, p)
		}

		playlistInfo.Items = playlist

		return playlistInfo, nil
	}

	return nil, nil
}

func (c *Client) httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// Youtube responses depend on language and user agent
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:70.0) Gecko/20100101 Firefox/70.0")

	return c.HTTPClient.Do(req)
}

func (c *Client) httpGetAndCheckResponse(url string) (*http.Response, error) {
	resp, err := c.httpGet(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return resp, nil
}

func (c *Client) httpGetAndCheckResponseReadBody(url string) ([]byte, error) {
	resp, err := c.httpGetAndCheckResponse(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
