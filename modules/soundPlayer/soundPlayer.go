package soundPlayer

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

const (
	moduleName = "SoundPlayer"
)

var (
	// sound encoding settings
	bitRate      = 128
	maxQueueSize = 6
)

//SoundPlayer is a struct of the SoundPlayer module
type SoundPlayer struct {
	session *discordgo.Session
	owner   string

	// Map of Guild id's to *play channels, used for queuing and rate-limiting guilds
	queues map[string]chan *play
}

// play represents an individual use of the !airhorn command
type play struct {
	guildID   string
	channelID string
	userID    string
	sound     *sound

	// The next play to occur after this, only used forshould have chaining sounds like anotha
	next *play

	// If true, this was a forced play using a specific airhorn sound name
	forced bool
}

type soundCollection struct {
	prefix    string
	commands  []string
	sounds    []*sound
	chainWith *soundCollection

	soundRange int
}

// sound represents a sound clip
type sound struct {
	name string

	// Weight adjust how likely it is this song will play, higher = more likely
	weight int

	// Delay (in milliseconds) for the bot to wait before sending the disconnect request
	partDelay int

	// Buffer to store encoded PCM packets
	buffer [][]byte
}

// Array of all the sounds we have
var AIRHORN *soundCollection = &soundCollection{
	prefix: "airhorn",
	commands: []string{
		"!airhorn",
	},
	sounds: []*sound{
		createSound("default", 1000, 250),
		createSound("reverb", 800, 250),
		createSound("spam", 800, 0),
		createSound("tripletap", 800, 250),
		createSound("fourtap", 800, 250),
		createSound("distant", 500, 250),
		createSound("echo", 500, 250),
		createSound("clownfull", 250, 250),
		createSound("clownshort", 250, 250),
		createSound("clownspam", 250, 0),
		createSound("highfartlong", 200, 250),
		createSound("highfartshort", 200, 250),
		createSound("midshort", 100, 250),
		createSound("truck", 10, 250),
	},
}

var KHALED *soundCollection = &soundCollection{
	prefix:    "another",
	chainWith: AIRHORN,
	commands: []string{
		"!anotha",
		"!anothaone",
	},
	sounds: []*sound{
		createSound("one", 1, 250),
		createSound("one_classic", 1, 250),
		createSound("one_echo", 1, 250),
	},
}

var CENA *soundCollection = &soundCollection{
	prefix: "jc",
	commands: []string{
		"!johncena",
		"!cena",
	},
	sounds: []*sound{
		createSound("airhorn", 1, 250),
		createSound("echo", 1, 250),
		createSound("full", 1, 250),
		createSound("jc", 1, 250),
		createSound("nameis", 1, 250),
		createSound("spam", 1, 250),
	},
}

var ETHAN *soundCollection = &soundCollection{
	prefix: "ethan",
	commands: []string{
		"!ethan",
		"!eb",
		"!ethanbradberry",
		"!h3h3",
	},
	sounds: []*sound{
		createSound("areyou_classic", 100, 250),
		createSound("areyou_condensed", 100, 250),
		createSound("areyou_crazy", 100, 250),
		createSound("areyou_ethan", 100, 250),
		createSound("classic", 100, 250),
		createSound("echo", 100, 250),
		createSound("high", 100, 250),
		createSound("slowandlow", 100, 250),
		createSound("cuts", 30, 250),
		createSound("beat", 30, 250),
		createSound("sodiepop", 1, 250),
	},
}

var COW *soundCollection = &soundCollection{
	prefix: "cow",
	commands: []string{
		"!stan",
		"!stanislav",
	},
	sounds: []*sound{
		createSound("herd", 10, 250),
		createSound("moo", 10, 250),
		createSound("x3", 1, 250),
	},
}

var BIRTHDAY *soundCollection = &soundCollection{
	prefix: "birthday",
	commands: []string{
		"!birthday",
		"!bday",
	},
	sounds: []*sound{
		createSound("horn", 50, 250),
		createSound("horn3", 30, 250),
		createSound("sadhorn", 25, 250),
		createSound("weakhorn", 25, 250),
	},
}

var WOW *soundCollection = &soundCollection{
	prefix: "wow",
	commands: []string{
		"!wowthatscool",
		"!wtc",
	},
	sounds: []*sound{
		createSound("thatscool", 50, 250),
	},
}

var HYPE *soundCollection = &soundCollection{
	prefix: "hype",
	commands: []string{
		"!hype",
		"!na",
		"!cs",
	},
	sounds: []*sound{
		createSound("bestteam", 100, 250),
		createSound("brabrabravo", 100, 250),
		createSound("brabrabravo2", 100, 250),
		createSound("cans", 100, 250),
		createSound("cans2", 100, 250),
		createSound("givenoise", 100, 250),
		createSound("givenoise2", 100, 250),
		createSound("givenoise3", 100, 250),
		createSound("goingtowar", 100, 250),
		createSound("herewegoagain", 100, 250),
		createSound("millions", 100, 250),
		createSound("poland", 100, 250),
		createSound("ready", 100, 250),
		createSound("show", 100, 250),
		createSound("veryexciting", 100, 250),
		createSound("warisover", 100, 250),
	},
}

var COLLECTIONS []*soundCollection = []*soundCollection{
	AIRHORN,
	KHALED,
	CENA,
	ETHAN,
	COW,
	BIRTHDAY,
	WOW,
	HYPE,
}

// GetModuleName returns the name of the module
func GetModuleName() string {
	return moduleName
}

// Create a sound struct
func createSound(name string, weight int, partDelay int) *sound {
	return &sound{
		name:      name,
		weight:    weight,
		partDelay: partDelay,
		buffer:    make([][]byte, 0),
	}
}

func (sc *soundCollection) load() {
	for _, sound := range sc.sounds {
		sc.soundRange += sound.weight
		sound.load(sc)
	}
}

func (sc *soundCollection) random() *sound {
	var i int
	number := randomRange(0, sc.soundRange)

	for _, sound := range sc.sounds {
		i += sound.weight

		if number < i {
			return sound
		}
	}
	return nil
}

// load attempts to load an encoded sound file from disk
// DCA files are pre-computed sound files that are easy to send to Discord.
// If you would like to create your own DCA files, please use:
// https://github.com/nstafie/dca-rs
// eg: dca-rs --raw -i <input wav file> > <output file>
func (s *sound) load(c *soundCollection) error {
	path := fmt.Sprintf("audio/%v_%v.dca", c.prefix, s.name)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// read opus frame length from dca file
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			fmt.Println("error reading from dca file :", err)
			return err
		}

		// read encoded pcm from dca file
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("error reading from dca file :", err)
			return err
		}

		// append encoded pcm data to the buffer
		s.buffer = append(s.buffer, InBuf)
	}
}

// Plays this sound over the specified VoiceConnection
func (s *sound) play(vc *discordgo.VoiceConnection) {
	vc.Speaking(true)
	defer vc.Speaking(false)

	for _, buff := range s.buffer {
		vc.OpusSend <- buff
	}
}

// Attempts to find the current users voice channel inside a given guild
func (p *SoundPlayer) getCurrentVoiceChannel(user *discordgo.User, guild *discordgo.Guild) *discordgo.Channel {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			channel, _ := p.session.State.Channel(vs.ChannelID)
			return channel
		}
	}
	return nil
}

// Returns a random integer between min and max
func randomRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// Prepares a play
func (p *SoundPlayer) createPlay(user *discordgo.User, guild *discordgo.Guild, coll *soundCollection, sound *sound) *play {
	// Grab the users voice channel
	channel := p.getCurrentVoiceChannel(user, guild)
	if channel == nil {
		log.WithFields(log.Fields{
			"user":  user.ID,
			"guild": guild.ID,
		}).Warning("Failed to find channel to play sound in")
		return nil
	}

	// Create the play
	newPlay := &play{
		guildID:   guild.ID,
		channelID: channel.ID,
		userID:    user.ID,
		sound:     sound,
		forced:    true,
	}

	// If we didn't get passed a manual sound, generate a random one
	if newPlay.sound == nil {
		newPlay.sound = coll.random()
		newPlay.forced = false
	}

	// If the collection is a chained one, set the next sound
	if coll.chainWith != nil {
		newPlay.next = &play{
			guildID:   newPlay.guildID,
			channelID: newPlay.channelID,
			userID:    newPlay.userID,
			sound:     coll.chainWith.random(),
			forced:    newPlay.forced,
		}
	}

	return newPlay
}

// Prepares and enqueues a play into the ratelimit/buffer guild queue
func (p *SoundPlayer) enqueuePlay(user *discordgo.User, guild *discordgo.Guild, coll *soundCollection, sound *sound) {
	newPlay := p.createPlay(user, guild, coll, sound)
	if newPlay == nil {
		return
	}

	// Check if we already have a connection to this guild
	//   yes, this isn't threadsafe, but its "OK" 99% of the time
	_, exists := p.queues[guild.ID]

	if exists {
		if len(p.queues[guild.ID]) < maxQueueSize {
			p.queues[guild.ID] <- newPlay
		}
	} else {
		p.queues[guild.ID] = make(chan *play, maxQueueSize)
		p.playSound(newPlay, nil)
	}
}

// play a sound
func (p *SoundPlayer) playSound(play *play, vc *discordgo.VoiceConnection) (err error) {
	log.WithFields(log.Fields{
		"play": play,
	}).Info("Playing sound")

	if vc == nil {
		vc, err = p.session.ChannelVoiceJoin(play.guildID, play.channelID, false, false)
		// vc.Receive = false
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to play sound")
			delete(p.queues, play.guildID)
			return err
		}
	}

	// If we need to change channels, do that now
	if vc.ChannelID != play.channelID {
		vc.ChangeChannel(play.channelID, false, false)
		time.Sleep(time.Millisecond * 125)
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(time.Millisecond * 32)

	// play the sound
	play.sound.play(vc)

	// If this is chained, play the chained sound
	if play.next != nil {
		p.playSound(play.next, vc)
	}

	// If there is another song in the queue, recurse and play that
	if len(p.queues[play.guildID]) > 0 {
		play = <-p.queues[play.guildID]
		p.playSound(play, vc)
		return nil
	}

	// If the queue is empty, delete it
	time.Sleep(time.Millisecond * time.Duration(play.sound.partDelay))
	delete(p.queues, play.guildID)
	vc.Disconnect()
	return nil
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable != nil {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			s.ChannelMessageSend(channel.ID, "**AIRHORN BOT READY FOR HORNING. TYPE `!AIRHORN` WHILE IN A VOICE CHANNEL TO ACTIVATE**")
			return
		}
	}
}

func scontains(key string, options ...string) bool {
	for _, item := range options {
		if item == key {
			return true
		}
	}
	return false
}

func (p *SoundPlayer) displayBotStats(cid string) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	users := 0
	for _, guild := range p.session.State.Ready.Guilds {
		users += len(guild.Members)
	}

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "Discordgo: \t%s\n", discordgo.VERSION)
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	fmt.Fprintf(w, "Memory: \t%s / %s (%s total allocated)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Tasks: \t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "Servers: \t%d\n", len(p.session.State.Ready.Guilds))
	fmt.Fprintf(w, "Users: \t%d\n", users)
	fmt.Fprintf(w, "```\n")
	w.Flush()
	p.session.ChannelMessageSend(cid, buf.String())
}

func utilGetMentioned(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.User {
	for _, mention := range m.Mentions {
		if mention.ID != s.State.Ready.User.ID {
			return mention
		}
	}
	return nil
}

func (p *SoundPlayer) airhornBomb(cid string, guild *discordgo.Guild, user *discordgo.User, cs string) {
	count, _ := strconv.Atoi(cs)
	p.session.ChannelMessageSend(cid, ":ok_hand:"+strings.Repeat(":trumpet:", count))

	// Cap it at something
	if count > 100 {
		return
	}

	play := p.createPlay(user, guild, AIRHORN, nil)
	vc, err := p.session.ChannelVoiceJoin(play.guildID, play.channelID, true, true)
	if err != nil {
		return
	}

	for i := 0; i < count; i++ {
		AIRHORN.random().play(vc)
	}

	vc.Disconnect()
}

// Handles bot operator messages, should be refactored (lmao)
func (p *SoundPlayer) handleBotControlMessages(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) {
	if scontains(parts[1], "status") {
		p.displayBotStats(m.ChannelID)
	} else if scontains(parts[1], "bomb") && len(parts) >= 4 {
		p.airhornBomb(m.ChannelID, g, utilGetMentioned(s, m), parts[3])
	}
}

func (p *SoundPlayer) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= 0 || (m.Content[0] != '!' && len(m.Mentions) < 1) {
		return
	}

	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(strings.ToLower(msg), " ")

	channel, _ := s.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := s.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}

	// If this is a mention, it should come from the owner (otherwise we don't care)
	if len(m.Mentions) > 0 && m.Author.ID == p.owner && len(parts) > 0 {
		mentioned := false
		for _, mention := range m.Mentions {
			mentioned = (mention.ID == s.State.Ready.User.ID)
			if mentioned {
				break
			}
		}

		if mentioned {
			p.handleBotControlMessages(s, m, parts, guild)
		}
		return
	}

	// Find the collection for the command we got
	for _, coll := range COLLECTIONS {
		if scontains(parts[0], coll.commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var sound *sound
			if len(parts) > 1 {
				for _, s := range coll.sounds {
					if parts[1] == s.name {
						sound = s
					}
				}

				if sound == nil {
					return
				}
			}

			go p.enqueuePlay(m.Author, guild, coll, sound)
			return
		}
	}
}

// Initialize starts the module
func Initialize(s *discordgo.Session) {
	var (
		//	Token = flag.String("t", "", "Discord Authentication Token")
		owner = flag.String("o", "", "Owner ID")
	)

	//Create new SoundPlayer instance
	player := new(SoundPlayer)
	player.session = s
	if *owner != "" {
		player.owner = *owner
	}

	// Preload all the sounds
	for _, coll := range COLLECTIONS {
		coll.load()
	}

	player.session.AddHandler(onGuildCreate)
	player.session.AddHandler(player.onMessageCreate)

	// Fully initialzed
	log.Info("Initialized SoundPlayer module")
}
