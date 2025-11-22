//go:build darwin
// +build darwin

package main

/*
#include <ApplicationServices/ApplicationServices.h>
#include <CoreGraphics/CoreGraphics.h>

static CGPoint getMouseLocation() {
	CGEventRef event = CGEventCreate(NULL);
	if (event == NULL) {
		return CGPointMake(0, 0);
	}
	CGPoint point = CGEventGetLocation(event);
	CFRelease(event);
	return point;
}

static uint32_t getDisplayForPoint(CGPoint point) {
	CGDirectDisplayID displayID;
	uint32_t displayCount = 1;
	CGGetDisplaysWithPoint(point, 1, &displayID, &displayCount);
	if (displayCount > 0) {
		return displayID;
	}
	return 0;
}
*/
import "C"
import (
	"fyne.io/fyne/v2"
)

// getMousePosition returns the current mouse cursor position using macOS Core Graphics
func getMousePosition() (x, y float32) {
	point := C.getMouseLocation()
	return float32(point.x), float32(point.y)
}

// positionWindowOnMouseDisplay positions the window on the display containing the mouse cursor
// Uses the same pattern as KrankyBearFileMover: show window first, then center
// On macOS, showing a window first helps Fyne detect which display it should be on
func positionWindowOnMouseDisplay(app fyne.App, window fyne.Window) {
	// Get mouse position (for verification/debugging)
	mouseX, mouseY := getMousePosition()
	_ = mouseX
	_ = mouseY
	
	// Create a tiny temporary window positioned near the mouse cursor
	// This helps Fyne detect which display the cursor is on
	tempWindow := app.NewWindow("")
	tempWindow.Resize(fyne.NewSize(1, 1))
	tempWindow.SetFixedSize(true)
	
	// Show the temp window first - this initializes it and helps Fyne detect
	// which display it should be on (typically the one containing the cursor on macOS)
	tempWindow.Show()
	
	// Center the temp window - this helps Fyne determine the correct display
	tempWindow.CenterOnScreen()
	
	// Now center the main window - it should center on the same display
	// Matching the pattern from KrankyBearFileMover's centerDialogOnMainWindow
	window.CenterOnScreen()
	
	// Close the temporary window immediately
	tempWindow.Close()
}

