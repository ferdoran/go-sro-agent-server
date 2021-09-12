package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Simple Gtk3 Application written in go.
// This application creates a window on the application callback activate.
// More GtkApplication info can be found here -> https://wiki.gnome.org/HowDoI/GtkApplication

func main() {
	// Create Gtk Application, change appID to your application domain name reversed.
	const appID = "org.gtk.example"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	// Check to make sure no errors when creating Gtk Application
	if err != nil {
		log.Fatal("Could not create application.", err)
	}
	// Application signals available
	// startup -> sets up the application when it first starts
	// activate -> shows the default first window of the application (like a new document). This corresponds to the application being launched by the desktop environment.
	// open -> opens files and shows them in a new window. This corresponds to someone trying to open a document (or documents) using the application from the file browser, or similar.
	// shutdown ->  performs shutdown tasks
	// Setup Gtk Application callback signals
	application.Connect("activate", func() { onActivate(application) })
	// Run Gtk application
	os.Exit(application.Run(os.Args))
}

// Callback signal from Gtk Application
func onActivate(application *gtk.Application) {
	// Create ApplicationWindow
	appWindow, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}
	// Set ApplicationWindow Properties
	appWindow.SetTitle("Basic Application.")
	appWindow.SetDefaultSize(400, 400)
	appWindow.Show()
}
