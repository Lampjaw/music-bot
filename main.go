package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/lampjaw/discordgobot"
)

const token = ""

func main() {
	log.Println("Running...")

	config := &discordgobot.GobotConf{
		CommandPrefix: "?",
	}

	b, err := discordgobot.NewBot(token, config, nil)
	if err != nil {
		log.Println(err)
	}

	b.RegisterPlugin(NewMusicPlugin())

	b.Open()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

out:
	for {
		select {
		case <-c:
			break out
		}
	}
}
