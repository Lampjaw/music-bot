package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lampjaw/discordgobot"
)

type MusicPlugin struct {
	discordgobot.Plugin
	player *MusicPlayer
}

func NewMusicPlugin() discordgobot.IPlugin {
	return &MusicPlugin{}
}

func (p *MusicPlugin) Name() string {
	return "Music"
}

func (p *MusicPlugin) Commands() []*discordgobot.CommandDefinition {
	return []*discordgobot.CommandDefinition{
		&discordgobot.CommandDefinition{
			CommandID: "music-play",
			Triggers: []string{
				"play",
			},
			Arguments: []discordgobot.CommandDefinitionArgument{
				discordgobot.CommandDefinitionArgument{
					Optional: false,
					Pattern:  ".+",
					Alias:    "url",
				},
			},
			Description: "Plays a song or playlist with the given url",
			Callback:    p.runPlayMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-disconnect",
			Triggers: []string{
				"disconnect",
				"dc",
			},
			Description: "Disconnect the bot from the voice channel it is in",
			Callback:    p.runDisconnectMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-nowplaying",
			Triggers: []string{
				"nowplaying",
				"np",
			},
			Description: "Shows what song the bot is currently playing",
			Callback:    p.runNowPlayingMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-skip",
			Triggers: []string{
				"skip",
			},
			Description: "Skips the currently playing song",
			Callback:    p.runSkipMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-remove",
			Triggers: []string{
				"remove",
			},
			Arguments: []discordgobot.CommandDefinitionArgument{
				discordgobot.CommandDefinitionArgument{
					Optional: false,
					Pattern:  ".+",
					Alias:    "videoID",
				},
			},
			Description: "Removes a certain entry from the queue",
			Callback:    p.runRemoveMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-loopqueue",
			Triggers: []string{
				"loopqueue",
				"lq",
			},
			Description: "Loops the whole queue",
			Callback:    p.runLoopQueueMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-loop",
			Triggers: []string{
				"loop",
			},
			Description: "Loop the currently playing song",
			Callback:    p.runLoopMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-resume",
			Triggers: []string{
				"resume",
			},
			Description: "Resume paused music",
			Callback:    p.runResumeMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-skipto",
			Triggers: []string{
				"skipto",
			},
			Arguments: []discordgobot.CommandDefinitionArgument{
				discordgobot.CommandDefinitionArgument{
					Optional: false,
					Pattern:  ".+",
					Alias:    "queuePosition",
				},
			},
			Description: "Skips to a certain position in the queue",
			Callback:    p.runSkipToMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-clear",
			Triggers: []string{
				"clear",
			},
			Description: "Clears the queue",
			Callback:    p.runClearMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-replay",
			Triggers: []string{
				"replay",
			},
			Description: "Reset the progress of the current song",
			Callback:    p.runReplayMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-pause",
			Triggers: []string{
				"pause",
			},
			Description: "Pauses the currently playing track",
			Callback:    p.runPauseMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-removedupes",
			Triggers: []string{
				"removedupes",
			},
			Description: "Removes duplicate songs from the queue",
			Callback:    p.runRemoveDupesMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-shuffle",
			Triggers: []string{
				"shuffle",
			},
			Description: "Shuffles the queue",
			Callback:    p.runShuffleMusicCommand,
		},
		&discordgobot.CommandDefinition{
			CommandID: "music-queue",
			Triggers: []string{
				"queue",
			},
			Description: "View the queue",
			Callback:    p.runQueueMusicCommand,
		},
	}
}

func (p *MusicPlugin) runPlayMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	userID := payload.Message.UserID()
	guildID, _ := payload.Message.ResolveGuildID()

	voiceState := findVoiceChannel(client.Session, guildID, userID)
	if voiceState == nil {
		client.SendMessage(payload.Message.Channel(), "You must be in a voice channel to use this command.")
		return
	}

	if p.player == nil {
		p.player = NewMusicPlayer()
	}

	if payload.Arguments["url"] != "" {
		ytURL := payload.Arguments["url"]
		if strings.Contains(ytURL, "playlist") {
			playlist, _ := p.player.AddPlaylistToQueue(ytURL)

			client.SendMessage(payload.Message.Channel(), fmt.Sprintf("Adding %v songs to the queue from `%s`", len(playlist.Items), playlist.Title))
		} else {
			vid, _ := p.player.AddSongToQueue(ytURL)

			client.SendMessage(payload.Message.Channel(), fmt.Sprintf("Adding `%s` to the queue", vid.Title))
		}
	}

	go p.playMusicInChannel(client.Session, voiceState.GuildID, voiceState.ChannelID)
}

func (p *MusicPlugin) runDisconnectMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.Shutdown()
		p.player = nil
	}
}

func (p *MusicPlugin) runNowPlayingMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {

}

func (p *MusicPlugin) runSkipMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.Skip()
	}
}

func (p *MusicPlugin) runRemoveMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {

}

func (p *MusicPlugin) runLoopQueueMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.loopQueue = !p.player.loopQueue
		if p.player.loopQueue {
			client.SendMessage(payload.Message.Channel(), "Queue looping enabled!")
		} else {
			client.SendMessage(payload.Message.Channel(), "Queue looping disabled")
		}
	}
}

func (p *MusicPlugin) runLoopMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.loopSong = !p.player.loopSong
		if p.player.loopSong {
			client.SendMessage(payload.Message.Channel(), "Song looping enabled!")
		} else {
			client.SendMessage(payload.Message.Channel(), "Song looping enabled!")
		}
	}
}

func (p *MusicPlugin) runResumeMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		guildID, _ := payload.Message.ResolveGuildID()
		userID := payload.Message.UserID()

		voiceState := findVoiceChannel(client.Session, guildID, userID)
		if voiceState == nil {
			client.SendMessage(payload.Message.Channel(), "You must be in a voice channel to use this command.")
			return
		}

		go p.playMusicInChannel(client.Session, voiceState.GuildID, voiceState.ChannelID)
	}
}

func (p *MusicPlugin) runSkipToMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {

}

func (p *MusicPlugin) runClearMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.ClearQueue()
		client.SendMessage(payload.Message.Channel(), "Queue cleared!")
	}
}

func (p *MusicPlugin) runReplayMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.Replay()
	}
}

func (p *MusicPlugin) runPauseMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {

}

func (p *MusicPlugin) runRemoveDupesMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.RemoveDuplicates()
		client.SendMessage(payload.Message.Channel(), "Duplicates removed!")
	}
}

func (p *MusicPlugin) runShuffleMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	if p.player != nil {
		p.player.Shuffle()
		client.SendMessage(payload.Message.Channel(), "Songs shuffled!")
	}
}

func (p *MusicPlugin) runQueueMusicCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	var sb strings.Builder

	np := p.player.ActiveSong

	sb.WriteString("Now Playing:\n")
	sb.WriteString(fmt.Sprintf("`%s | %v`", np.Title, np.Duration))
	sb.WriteString("\n\nUpNext\n")

	var totalDuration time.Duration
	rowCount := 1
	for i, s := range p.player.SongQueue {
		if rowCount < 11 && !(i == 0 && s.VideoID == np.VideoID) {
			sb.WriteString(fmt.Sprintf("`%v. %s | %v`\n", rowCount, s.Title, s.Duration))
			rowCount++
		}

		totalDuration += s.Duration
	}

	sb.WriteString(fmt.Sprintf("\n\n**%v songs in queue | %v total length**", len(p.player.SongQueue), totalDuration))

	embed := &discordgo.MessageEmbed{
		Title:       "Queue",
		Color:       0x070707,
		Description: sb.String(),
	}

	client.SendEmbedMessage(payload.Message.Channel(), embed)
}

func findVoiceChannel(s *discordgo.Session, guildID string, userID string) *discordgo.VoiceState {
	guild, _ := s.Guild(guildID)

	for _, s := range guild.VoiceStates {
		if s.UserID == userID {
			return s
		}
	}

	return nil
}

func (p *MusicPlugin) playMusicInChannel(s *discordgo.Session, guildID string, channelID string) {
	if p.player.voiceConnection == nil {
		vc, _ := s.ChannelVoiceJoin(guildID, channelID, false, true)
		p.player.Join(vc)
	}

	p.player.Play()
}
