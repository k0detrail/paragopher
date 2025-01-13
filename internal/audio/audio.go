package audio

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const sampleRate = 196000

var audioContext *audio.Context = audio.NewContext(sampleRate)

func PlaySound(soundBytes []byte) {
	decoded, err := vorbis.DecodeWithSampleRate(
		sampleRate,
		bytes.NewReader(soundBytes),
	)
	if err != nil {
		log.Fatalf("Faied to decode OGG: %v", err)
	}
	player, err := audioContext.NewPlayer(decoded)
	if err != nil {
		log.Fatalf("Failed to create player: %v", err)
	}
	player.Play()
}
