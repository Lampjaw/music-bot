package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ebml-go/webm"
	"github.com/rylio/ytdl"
)

func PrepareSong(item *PlaylistItem) error {
	if item.VideoInfo == nil {
		vid, err := getVideoFromID(item.VideoID)
		if err != nil {
			return err
		}

		item.VideoInfo = vid
	}

	dlFormat := item.GetSongFormat()
	if dlFormat == nil {
		return fmt.Errorf("No suitable audio formats found")
	}

	fileName := getFileName(item)
	if fileExists(fileName) {
		return nil
	}

	os.Mkdir("tmp", os.ModeTemporary)
	file, _ := os.Create(fileName)
	err := item.VideoInfo.Download(dlFormat, file)
	if err != nil {
		return err
	}

	return nil
}

func GetSongFile(item *PlaylistItem) (*os.File, error) {
	fileName := getFileName(item)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func LoadSong(file *os.File) (*webm.Reader, error) {
	var w webm.WebM
	reader, err := webm.Parse(file, &w)
	return reader, err
}

func RemoveSong(item *PlaylistItem) error {
	fileName := getFileName(item)

	if fileName != "" && fileExists(fileName) {
		return os.Remove(fileName)
	}
}

func getPlaylistInfoFromURL(surl string) (*PlaylistInfo, error) {
	playlistURL, err := url.ParseRequestURI(surl)
	if err != nil {
		return nil, err
	}

	plID := extractPlaylistID(playlistURL)

	if plID == "" {
		return nil, nil
	}

	playlistInfo, err := DefaultClient.GetPlaylistInfoFromID(plID)
	if err != nil {
		return nil, err
	}

	return playlistInfo, nil
}

func getVideoFromURL(surl string) (*ytdl.VideoInfo, error) {
	videoURL, err := url.ParseRequestURI(surl)
	if err != nil {
		return nil, err
	}

	vID := extractVideoID(videoURL)

	if vID == "" {
		return nil, nil
	}

	return getVideoFromID(vID)
}

func getVideoFromID(id string) (*ytdl.VideoInfo, error) {
	videoInfo, err := ytdl.GetVideoInfoFromID(id)
	if err != nil {
		return nil, err
	}

	return videoInfo, nil
}

func extractVideoID(u *url.URL) string {
	switch u.Host {
	case "www.youtube.com", "youtube.com", "m.youtube.com":
		if u.Path == "/watch" {
			return u.Query().Get("v")
		}
		if strings.HasPrefix(u.Path, "/embed/") {
			return u.Path[7:]
		}
	case "youtu.be":
		if len(u.Path) > 1 {
			return u.Path[1:]
		}
	}
	return ""
}

func extractPlaylistID(u *url.URL) string {
	switch u.Host {
	case "www.youtube.com", "youtube.com", "m.youtube.com":
		if u.Path == "/playlist" {
			return u.Query().Get("list")
		}
	}
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() && info.Size() > 0
}

func getFileName(item *PlaylistItem) string {
	return fmt.Sprintf("tmp/%s.%s", item.VideoID, item.GetSongFormat().Extension)
}

func (vi *PlaylistItem) GetSongFormat() *ytdl.Format {
	var dlFormat *ytdl.Format
	if vi.VideoInfo != nil {
		for _, f := range vi.VideoInfo.Formats {
			if f.AudioEncoding == "opus" && (dlFormat == nil || f.AudioBitrate > dlFormat.AudioBitrate) {
				dlFormat = f
			}
		}
	}
	return dlFormat
}
