//go:generate fyne bundle -o bundled.go -a Resources/Images/KrankyBearVikingHelmet.png Resources/Images/KrankyBearHardHat.png Resources/Sounds/boing.mp3 Resources/Sounds/KrankyBearGrowl.mp3 Resources/Sounds/uhOh.mp3 Resources/Sounds/wheeHoo.mp3 Resources/Sounds/zip.mp3

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type appTheme struct {
	fyne.Theme
}

func (a *appTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameHeadingText {
		return a.Theme.Size(n) * 1.5
	}

	return a.Theme.Size(n)
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
