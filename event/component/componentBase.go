package component

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/modules/adventcalendar"
	"github.com/cake4everyone/cake4everybot/modules/faq"
	"github.com/cake4everyone/cake4everybot/modules/random"
	"github.com/cake4everyone/cake4everybot/modules/secretsanta"
)

var log = logger.New("Event/Component")

// Component is an interface wrapper for all message components.
type Component interface {
	// Function of a component.
	// All things that should happen after submitting or pressing a button.
	Handle(*discordgo.Session, *discordgo.InteractionCreate)

	// Custom ID of the modal to identify the module
	ID() string
}

// ComponentMap holds all active components. It maps them from a unique string identifier to the
// corresponding Component.
var ComponentMap = make(map[string]Component)

// Register registers add message components
func Register() {
	// This is the list of components to use. Add a component via
	// simply appending the struct (which must implement the
	// interface command.Component) to the list, e.g.:
	//
	//  componentList = append(componentList, mymodule.MyComponent{})
	var componentList []Component

	componentList = append(componentList, adventcalendar.Component{})
	componentList = append(componentList, faq.Component{})
	componentList = append(componentList, GenericComponents{})
	componentList = append(componentList, random.Component{})
	componentList = append(componentList, secretsanta.Component{})

	if len(componentList) == 0 {
		return
	}
	for _, c := range componentList {
		ComponentMap[c.ID()] = c
	}
	log.Printf("Added %d component handler(s)!", len(ComponentMap))
}
