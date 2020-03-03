package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ebml-go/webm"
	"github.com/rylio/ytdl"
)

type MusicPlayer struct {
	IsPlaying       bool
	ActiveSong      *PlaylistItem
	SongQueue       []*PlaylistItem
	loopQueue       bool
	loopSong        bool
	skip            chan bool
	replay          chan bool
	voiceConnection *discordgo.VoiceConnection
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		IsPlaying:       false,
		SongQueue:       make([]*PlaylistItem, 0),
		loopQueue:       false,
		loopSong:        false,
		skip:            make(chan bool),
		replay:          make(chan bool),
		voiceConnection: nil,
	}
}

func (p *MusicPlayer) Join(vc *discordgo.VoiceConnection) {
	p.voiceConnection = vc
}

func (p *MusicPlayer) AddPlaylistToQueue(url string) (*PlaylistInfo, error) {
	playlist, err := getPlaylistInfoFromURL(url)
	if err != nil {
		log.Printf("Failed to get playlist info: %v", err)
		return nil, err
	}

	if playlist == nil {
		return nil, nil
	}

	p.SongQueue = append(p.SongQueue, playlist.Items...)

	return playlist, nil
}

func (p *MusicPlayer) AddSongToQueue(url string) (*PlaylistItem, error) {
	video, err := getVideoFromURL(url)
	if err != nil {
		log.Printf("Failed to get video info: %v", err)
		return nil, err
	}

	if video == nil {
		return nil, nil
	}

	item := &PlaylistItem{
		VideoID:      video.ID,
		Title:        video.Title,
		Duration:     video.Duration,
		IsPlayable:   true,
		ThumbnailURL: video.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String(),
		VideoInfo:    video,
	}

	p.SongQueue = append(p.SongQueue, item)

	return item, nil
}

func (p *MusicPlayer) Play() {
	if p.IsPlaying {
		return
	}

	p.IsPlaying = true

	for len(p.SongQueue) > 0 && p.IsPlaying {
		p.playCurrentSong()
	}

	p.IsPlaying = false
}

func (p *MusicPlayer) Shutdown() {
	p.ClearQueue()
	p.Skip()
	p.voiceConnection.Disconnect()
}

func (p *MusicPlayer) Skip() {
	p.skip <- true
}

func (p *MusicPlayer) Replay() {
	p.replay <- true
}

func (p *MusicPlayer) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.SongQueue), func(i, j int) { p.SongQueue[i], p.SongQueue[j] = p.SongQueue[j], p.SongQueue[i] })
	if p.ActiveSong != p.SongQueue[0] {
		go PrepareSong(p.SongQueue[0])
	}
}

func (p *MusicPlayer) ClearQueue() {
	p.SongQueue = p.SongQueue[:1]
}

func (p *MusicPlayer) RemoveDuplicates() {
	keys := make(map[string]bool)
	list := make([]*PlaylistItem, 0)
	for _, entry := range p.SongQueue {
		if _, value := keys[entry.VideoID]; !value {
			keys[entry.VideoID] = true
			list = append(list, entry)
		}
	}
	p.SongQueue = list
}

func (p *MusicPlayer) RemoveSongFromQueue(videoID string) {
	if len(p.SongQueue) == 0 {
		return
	}

	sIdx := p.findSongIndex(videoID)

	var songExtension string
	if p.SongQueue[sIdx] != nil && p.SongQueue[sIdx].GetSongFormat() != nil {
		songExtension = p.SongQueue[sIdx].GetSongFormat().Extension
	}

	p.SongQueue = append(p.SongQueue[:sIdx], p.SongQueue[sIdx+1:]...)

	fileName := fmt.Sprintf("%s.%s", videoID, songExtension)
	os.Remove(fileName)
}

func (p *MusicPlayer) playCurrentSong() {
	p.ActiveSong = p.SongQueue[0]
	defer p.postSongHandling(p.ActiveSong)

	err := PrepareSong(p.ActiveSong)
	if err != nil {
		log.Printf("Failed to prepare song: %v", err)
		return
	}

	if len(p.SongQueue) > 1 {
		go PrepareSong(p.SongQueue[1])
	}

	file, err := GetSongFile(p.ActiveSong)
	if err != nil {
		log.Printf("Failed to get song file: %v", err)
		return
	}
	defer file.Close()

	reader, err := LoadSong(file)
	if err != nil {
		log.Printf("Failed to load song: %v", err)
		return
	}

	p.voiceConnection.Speaking(true)
	defer p.voiceConnection.Speaking(false)

	p.sendSongData(reader)
}

func (p *MusicPlayer) postSongHandling(item *PlaylistItem) {
	p.ActiveSong = nil

	if p.loopQueue {
		sIdx := p.findSongIndex(item.VideoID)
		p.SongQueue = append(p.SongQueue[:sIdx], p.SongQueue[sIdx+1:]...)
		p.SongQueue = append(p.SongQueue, item)
		return
	}

	if !p.loopSong {
		p.RemoveSongFromQueue(item.VideoID)
	}
}

func (p *MusicPlayer) sendSongData(reader *webm.Reader) {
	reader.Seek(0)

	for {
		select {
		case packet := <-reader.Chan:
			if packet.Timecode == webm.BadTC {
				return
			}

			p.voiceConnection.OpusSend <- packet.Data
		case <-p.skip:
			p.loopSong = false
			return
		case <-p.replay:
			reader.Seek(0)
		}
	}
}

func (p *MusicPlayer) findSongIndex(videoID string) int {
	for i := range p.SongQueue {
		if p.SongQueue[i].VideoID == videoID {
			return i
		}
	}
	return -1
}
