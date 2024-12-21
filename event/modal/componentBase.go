package modal

import (
	"cake4everybot/modules/secretsanta"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Modal is an interface wrapper for all message components.
type Modal interface {
	// Function of a component.
	// All things that should happen after submitting a modal.
	HandleModal(*discordgo.Session, *discordgo.InteractionCreate)

	// Custom ID of the modal to identify the module
	ID() string
}

// ModalMap holds all active modals. It maps them from a unique string identifier to the
// corresponding [Modal].
var ModalMap = make(map[string]Modal)

// Register registers modals
func Register() {
	// This is the list of modals to use. Add a modal via
	// simply appending the struct (which must implement the
	// [Modal] interface) to the list, e.g.:
	//
	//  ModalList = append(ModalList, mymodule.MyComponent{})
	var ModalList []Modal

	ModalList = append(ModalList, secretsanta.Component{})

	if len(ModalList) == 0 {
		return
	}
	for _, c := range ModalList {
		ModalMap[c.ID()] = c
	}
	log.Printf("Added %d modal handler(s)!", len(ModalMap))
}
