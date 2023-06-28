package responses

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stablecog/sc-go/log"
)

// Send a message that only the user can see as a response to an interaction
func PrivateInteractionResponseWithComponents(s *discordgo.Session, i *discordgo.InteractionCreate, title, content, footer string, components []discordgo.MessageComponent) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				NewEmbed(title, content, footer),
			},
			Components: components,
		},
	})
	if err != nil {
		log.Errorf("Failed to respond to interaction: %v", err)
	}
	return err
}

func PrivateInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, title, content, footer string) error {
	return PrivateInteractionResponseWithComponents(s, i, title, content, footer, nil)
}

func UnknownErrorPrivateInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return PrivateInteractionResponseWithComponents(s, i, "😔", "An unknown error occurred. Please try again later.", "", nil)
}

// Public messages
func PublicImageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, url string, components []discordgo.MessageComponent) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewImageEmbed(url),
			},
			Components: components,
		},
	})
	if err != nil {
		log.Errorf("Failed to respond to interaction: %v", err)
	}
	return err
}
