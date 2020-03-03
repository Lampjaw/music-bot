package main

import (
	"time"

	"github.com/rylio/ytdl"
)

type PlaylistInfo struct {
	Title        string
	Description  string
	ThumbnailURL string
	Items        []*PlaylistItem
}

type PlaylistItem struct {
	VideoID      string
	Title        string
	Duration     time.Duration
	IsPlayable   bool
	ThumbnailURL string
	VideoInfo    *ytdl.VideoInfo
}

type initialPlaylistData struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								ItemSectionRenderer struct {
									Contents []struct {
										PlaylistVideoListRenderer struct {
											Contents []struct {
												PlaylistVideoRenderer struct {
													VideoID   string `json:"videoId"`
													Thumbnail struct {
														Thumbnails []struct {
															URL    string `json:"url"`
															Width  int    `json:"width"`
															Height int    `json:"height"`
														} `json:"thumbnails"`
													} `json:"thumbnail"`
													Title struct {
														SimpleText string `json:"simpleText"`
													} `json:"title"`
													LengthSeconds int  `json:"lengthSeconds,string"`
													IsPlayable    bool `json:"isPlayable"`
												} `json:"playlistVideoRenderer"`
											} `json:"contents"`
										} `json:"playlistVideoListRenderer"`
									} `json:"contents"`
								} `json:"itemSectionRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Microformat struct {
		MicroformatDataRenderer struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnail   struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
		} `json:"microformatDataRenderer"`
	} `json:"microformat"`
}
