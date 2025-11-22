package main

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"fyne.io/fyne/v2"
)

// SoundManager manages sound effects
type SoundManager struct {
	enabled      bool
	sampleRate   beep.SampleRate
	placeSound   *beep.Buffer
	lineSound    *beep.Buffer
	multiLineSound *beep.Buffer
	gameOverSound *beep.Buffer
}

// NewSoundManager creates a new sound manager
func NewSoundManager() *SoundManager {
	sm := &SoundManager{
		enabled:    true, // Default to enabled
		sampleRate: beep.SampleRate(44100),
	}
	
	// Initialize speaker
	speaker.Init(sm.sampleRate, sm.sampleRate.N(time.Second/10))
	
	// Load sounds from bundled resources
	sm.loadSounds()
	
	return sm
}

// loadSounds loads all sound effects from bundled resources
func (sm *SoundManager) loadSounds() {
	// Load zip.mp3 for piece placement
	sm.placeSound = sm.loadSound(resourceZipMp3)
	
	// Load wheeHoo.mp3 for single line clear
	sm.lineSound = sm.loadSound(resourceWheeHooMp3)
	
	// Load KrankyBearGrowl.mp3 for multiple lines cleared
	sm.multiLineSound = sm.loadSound(resourceKrankyBearGrowlMp3)
	
	// Load uhOh.mp3 for game over
	sm.gameOverSound = sm.loadSound(resourceUhOhMp3)
}

// loadSound loads a sound from a Fyne resource
func (sm *SoundManager) loadSound(resource fyne.Resource) *beep.Buffer {
	if resource == nil {
		return nil
	}
	
	// Create a reader from the resource
	reader := bytes.NewReader(resource.Content())
	
	// Decode MP3
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		// Try to load from file system as fallback
		return sm.loadSoundFromFile(resource.Name())
	}
	defer streamer.Close()
	
	// Create buffer with target sample rate
	bufferFormat := beep.Format{
		SampleRate:  sm.sampleRate,
		NumChannels: format.NumChannels,
		Precision:   format.Precision,
	}
	buffer := beep.NewBuffer(bufferFormat)
	
	// Resample if needed and append to buffer
	if format.SampleRate != sm.sampleRate {
		resampled := beep.Resample(4, format.SampleRate, sm.sampleRate, streamer)
		buffer.Append(resampled)
	} else {
		buffer.Append(streamer)
	}
	
	return buffer
}

// loadSoundFromFile loads a sound from file system (fallback)
func (sm *SoundManager) loadSoundFromFile(filename string) *beep.Buffer {
	file, err := os.Open("Resources/Sounds/" + filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return nil
	}
	defer streamer.Close()
	
	// Create buffer with target sample rate
	bufferFormat := beep.Format{
		SampleRate:  sm.sampleRate,
		NumChannels: format.NumChannels,
		Precision:   format.Precision,
	}
	buffer := beep.NewBuffer(bufferFormat)
	
	// Resample if needed and append to buffer
	if format.SampleRate != sm.sampleRate {
		resampled := beep.Resample(4, format.SampleRate, sm.sampleRate, streamer)
		buffer.Append(resampled)
	} else {
		buffer.Append(streamer)
	}
	
	return buffer
}

// PlayPlace plays the sound for piece placement
func (sm *SoundManager) PlayPlace() {
	if !sm.enabled || sm.placeSound == nil {
		return
	}
	speaker.Play(sm.placeSound.Streamer(0, sm.placeSound.Len()))
}

// PlayLineClear plays the sound for line clear
func (sm *SoundManager) PlayLineClear() {
	if !sm.enabled || sm.lineSound == nil {
		return
	}
	speaker.Play(sm.lineSound.Streamer(0, sm.lineSound.Len()))
}

// PlayMultiLineClear plays the sound for multiple lines cleared
func (sm *SoundManager) PlayMultiLineClear() {
	if !sm.enabled || sm.multiLineSound == nil {
		return
	}
	speaker.Play(sm.multiLineSound.Streamer(0, sm.multiLineSound.Len()))
}

// PlayGameOver plays the game over sound
func (sm *SoundManager) PlayGameOver() {
	if !sm.enabled || sm.gameOverSound == nil {
		return
	}
	speaker.Play(sm.gameOverSound.Streamer(0, sm.gameOverSound.Len()))
}

// SetEnabled enables or disables sound
func (sm *SoundManager) SetEnabled(enabled bool) {
	sm.enabled = enabled
}

// IsEnabled returns whether sound is enabled
func (sm *SoundManager) IsEnabled() bool {
	return sm.enabled
}

