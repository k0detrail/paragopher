package game

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/audio"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/utils"

	"bytes"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var textFaceSource *text.GoTextFaceSource

type Game struct {
	Score        int
	gameData     *utils.GameData
	soundProfile *audio.SoundProfile

	selectingUser   bool
	userMenuIndex   int
	userInputName   string
	creatingUser    bool
	deletingUser    bool
	lastKeyNavTime  time.Time
	showLeaderboard bool

	showIntro     bool
	introStep     int
	lastIntroStep time.Time

	showExitDialog     bool
	showGameOverDialog bool

	barrelAngle            float64
	barrelImage            *ebiten.Image
	turretBaseImage        *ebiten.Image
	helicopterImage        *ebiten.Image
	paratrooperImage       *ebiten.Image
	paratrooperLandedImage *ebiten.Image
	paratrooperFellImage   *ebiten.Image
	bulletImage            *ebiten.Image

	bullets      []*Bullet
	lastShot     time.Time
	helicopters  []*Helicopter
	paratroopers []*Paratrooper
}

func NewGame() *Game {

	if textFaceSource == nil {
		var err error
		textFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
		if err != nil {
			log.Fatal(err)
		}
	}

	gameData, err := utils.LoadData()
	if err != nil {
		log.Println("Error loading game data!")
		gameData = &utils.GameData{}
	}
	game := &Game{
		bullets:        make([]*Bullet, 0),
		lastShot:       time.Now(),
		gameData:       gameData,
		soundProfile:   audio.NewSoundProfile(),
		showIntro:      true,
		selectingUser:  true,
		lastKeyNavTime: time.Now(),
	}
	game.initTurretImage()
	game.initBarrelImage()
	game.initBulletImage()
	game.initHelicopterImage()
	game.initParatrooperImage()
	game.initParatrooperLandedImage()
	game.initParatrooperFellImage()
	game.initIntro()

	return game
}

// Ebiten Game Interface
func (g *Game) Draw(screen *ebiten.Image) {
	if g.showLeaderboard {
		g.drawLeaderboard(screen)
		return
	}
	if g.selectingUser {
		g.drawUserMenu(screen)
		return
	}
	if g.showIntro {
		g.drawIntro(screen)
		return
	}
	g.drawTurret(screen)
	g.drawBullets(screen)
	g.drawHelicopters(screen)
	g.drawParatroopers(screen)

	// display score for current user
	currentHiScore := 0
	for _, u := range g.gameData.Users {
		if u.Name == g.gameData.CurrentUser {
			currentHiScore = u.HiScore
			break
		}
	}

	// Display Score
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("USER: %s   SCORE: %d    HI-SCORE: %d",
			g.gameData.CurrentUser, g.Score, currentHiScore),
	)

	if g.showExitDialog {
		showYesNoDialog(screen, "Do you want to exit the game?")
	}

	if g.showGameOverDialog {
		showYesNoDialog(screen, "GAME OVER!\nWould you like to start again?")
	}
}

func (g *Game) Update() error {
	if g.selectingUser {
		return g.updateUserMenu()
	}
	if g.showIntro {
		return nil
	}
	// update HiScore for current user
	for i := range g.gameData.Users {
		if g.gameData.Users[i].Name == g.gameData.CurrentUser {
			if g.Score > g.gameData.Users[i].HiScore {
				g.gameData.Users[i].HiScore = g.Score
			}
			break
		}
	}
	if g.showLeaderboard && ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.showLeaderboard = false
	}
	if g.showExitDialog {
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			if err := utils.SaveData(g.gameData); err != nil {
				log.Fatalf("Failed to save game dada: %v", err)
			}
			return config.ErrQuit
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			g.showExitDialog = false
		}
		return nil
	}
	if g.showGameOverDialog {
		if err := utils.SaveData(g.gameData); err != nil {
			log.Fatalf("Failed to save game dada: %v", err)
		}
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			g.Reset()
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			return config.ErrQuit
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.showExitDialog = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.barrelAngle > config.BarrelAngleMin {
			g.barrelAngle -= config.BarrelAngleStep
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.barrelAngle < config.BarrelAngleMax {
			g.barrelAngle += config.BarrelAngleStep
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if time.Since(g.lastShot).Milliseconds() > config.ShotCooldown {
			g.shoot()
		}
	}

	g.updateBullets()
	g.spawnHelicopter()
	g.updateHelicopters()
	g.updateParatroopers()
	g.checkHits()

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) Reset() {
	g.soundProfile.GameOverPlayer.Pause()
	g.Score = 0
	g.showExitDialog = false
	g.showGameOverDialog = false
	g.barrelAngle = 0.0
	g.bullets = nil
	g.helicopters = nil
	g.paratroopers = nil
}

func showYesNoDialog(screen *ebiten.Image, message string) {
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(config.SemiTransparentBlack)
	screen.DrawImage(overlay, nil)

	dialogWidth, dialogHeight := 300, 150
	dialogX := (screen.Bounds().Dx() - dialogWidth) / 2
	dialogY := (screen.Bounds().Dy() - dialogHeight) / 2
	dialog := ebiten.NewImage(dialogWidth, dialogHeight)
	dialog.Fill(config.ColourDarkGrey)

	vector.DrawFilledRect(
		dialog,
		0,
		0,
		float32(dialogWidth),
		5,
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		0,
		float32(dialogHeight-5),
		float32(dialogWidth),
		5,
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		0,
		0,
		5,
		float32(dialogHeight),
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		float32(dialogWidth-5),
		0,
		5,
		float32(dialogHeight),
		config.ColourBlack,
		false,
	)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dialogX), float64(dialogY))
	screen.DrawImage(dialog, op)

	textX := dialogX + 50
	textY := dialogY + 40
	ebitenutil.DebugPrintAt(screen, message, textX, textY)

	yesText := "Y: Yes"
	noText := "N: No"
	ebitenutil.DebugPrintAt(screen, yesText, dialogX+50, dialogY+90)
	ebitenutil.DebugPrintAt(screen, noText, dialogX+200, dialogY+90)
}

func (g *Game) drawUserMenu(screen *ebiten.Image) {
	screen.Fill(config.ColourBlack)

	fontFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   16,
	}

	title := "SELECT USER"
	titleW, _ := text.Measure(title, fontFace, 1.0)

	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate((config.ScreenWidth-titleW)/2.0, 50)
	titleOp.ColorScale.ScaleWithColor(config.ColourTeal)
	text.Draw(screen, title, fontFace, titleOp)

	// draw user list
	baseY := 100.0
	for i, user := range g.gameData.Users {
		name := fmt.Sprintf("%s (HiScore: %d)", user.Name, user.HiScore)
		if i == g.userMenuIndex {
			name = "→ " + name
		} else {
			name = "  " + name
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(100, baseY+float64(i*24))
		if i == g.userMenuIndex {
			op.ColorScale.ScaleWithColor(config.ColourPink)
		} else {
			op.ColorScale.ScaleWithColor(config.ColourWhite)
		}
		text.Draw(screen, name, fontFace, op)
	}

	controlStartY := float64(config.ScreenHeight - 160)
	// show controls
	controls := []string{
		"↑ / ↓ : Navigate",
		"ENTER : Select",
		"N     : New User",
		"D     : Delete User",
		"ESC   : Cancel",
		"L     : View Leaderboard",
	}

	for i, line := range controls {
		op := &text.DrawOptions{}
		op.GeoM.Translate(100, controlStartY+float64(i*18))
		op.ColorScale.ScaleWithColor(config.ColourWhite)
		text.Draw(screen, line, fontFace, op)
	}

	// show name input field
	if g.creatingUser {
		cursor := "_"
		if time.Now().UnixNano()/1e8%2 == 0 {
			cursor = " "
		}
		inputLabel := "ENTER USERNAME: " + g.userInputName + cursor
		inputW, _ := text.Measure(inputLabel, fontFace, 1.0)

		inputOp := &text.DrawOptions{}
		inputOp.GeoM.Translate((config.ScreenWidth-inputW)/2.0, config.ScreenHeight-30)
		inputOp.ColorScale.ScaleWithColor(config.ColourTeal)
		text.Draw(screen, inputLabel, fontFace, inputOp)
	}
}

func (g *Game) updateUserMenu() error {
	if g.creatingUser {
		// type username
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			trimmed := strings.TrimSpace(g.userInputName)

			if trimmed == "" {
				return nil
			}

			for _, u := range g.gameData.Users {
				if u.Name == trimmed {
					return nil
				}
			}

			newUser := utils.User{Name: g.userInputName}
			g.gameData.Users = append(g.gameData.Users, newUser)
			g.gameData.CurrentUser = newUser.Name
			g.selectingUser = false
			utils.SaveData(g.gameData)
		} else {
			for _, char := range ebiten.InputChars() {
				if len(g.userInputName) < 12 && (char >= 32 && char <= 126) {
					g.userInputName += string(char)
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.userInputName) > 0 {
				g.userInputName = g.userInputName[:len(g.userInputName)-1]
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.creatingUser = false
			g.userInputName = ""
			return nil
		}
		// leaderboards
		if ebiten.IsKeyPressed(ebiten.KeyL) {
			g.showLeaderboard = true
		}
		return nil
	}

	// navigate menu
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		time.Since(g.lastKeyNavTime) > 200*time.Millisecond {
		g.lastKeyNavTime = time.Now()
		if g.userMenuIndex < len(g.gameData.Users)-1 {
			g.userMenuIndex++
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
		time.Since(g.lastKeyNavTime) > 200*time.Millisecond {
		g.lastKeyNavTime = time.Now()
		if g.userMenuIndex > 0 {
			g.userMenuIndex--
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.showLeaderboard = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyN) {
		g.creatingUser = true
		g.userInputName = ""
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && len(g.gameData.Users) > 0 {
		// delete user
		idx := g.userMenuIndex
		g.gameData.Users = append(g.gameData.Users[:idx], g.gameData.Users[idx+1:]...)
		if g.userMenuIndex >= len(g.gameData.Users) && g.userMenuIndex > 0 {
			g.userMenuIndex--
		}
		utils.SaveData(g.gameData)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.gameData.Users) > 0 {
		// prevent transition if in the middle of creating a user
		if g.creatingUser {
			return nil
		}
		g.selectingUser = false
		g.showIntro = true
		g.introStep = 0
		g.lastIntroStep = time.Now()
		g.soundProfile.IntroPlayer.Rewind() // replay intro audio
		g.gameData.CurrentUser = g.gameData.Users[g.userMenuIndex].Name
		utils.SaveData(g.gameData)
		return nil
	}
	return nil
}

func (g *Game) drawLeaderboard(screen *ebiten.Image) {
	screen.Fill(config.ColourBlack)

	fontFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   16,
	}

	// sort users by HiScore, this is descending
	users := make([]utils.User, len(g.gameData.Users))
	copy(users, g.gameData.Users)
	sort.Slice(users, func(i, j int) bool {
		return users[i].HiScore > users[j].HiScore
	})

	// title
	title := "L E A D E R B O A R D"
	titleW, _ := text.Measure(title, fontFace, 1.0)
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate((config.ScreenWidth-titleW)/2.0, 50)
	titleOp.ColorScale.ScaleWithColor(config.ColourTeal)
	text.Draw(screen, title, fontFace, titleOp)

	// entries
	startY := 100.0
	lineHeight := 22.0

	for i, user := range users {
		entry := fmt.Sprintf("%2d. %-10s  -  %d", i+1, user.Name, user.HiScore)
		entryW, _ := text.Measure(entry, fontFace, 1.0)
		op := &text.DrawOptions{}
		op.GeoM.Translate((config.ScreenWidth-entryW)/2.0, startY+float64(i)*lineHeight)

		// color top 3 entries
		switch i {
		case 0:
			op.ColorScale.ScaleWithColor(config.ColourPink)
		case 1:
			op.ColorScale.ScaleWithColor(config.ColourTeal)
		case 2:
			op.ColorScale.ScaleWithColor(config.ColourWhite)
		default:
			op.ColorScale.ScaleWithColor(config.ColourDarkGrey)
		}
		text.Draw(screen, entry, fontFace, op)
	}

	// instructions
	escText := "ESC: Back"
	escW, _ := text.Measure(escText, fontFace, 1.0)
	escOp := &text.DrawOptions{}
	escOp.GeoM.Translate((config.ScreenWidth-escW)/2.0, config.ScreenHeight-40)
	escOp.ColorScale.ScaleWithColor(config.ColourWhite)
	text.Draw(screen, escText, fontFace, escOp)
}
