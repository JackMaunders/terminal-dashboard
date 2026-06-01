package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/qeesung/image2ascii/convert"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PokemonTickMsg time.Time

type PokemonWidget struct {
	name     string
	asciiArt string
	hasError bool
}

type pokemonResponse struct {
	Name    string `json:"name"`
	Sprites struct {
		Other struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`
}

func NewPokemonWidget() PokemonWidget {
	// Get random number (original Pokemon)
	id := rand.Intn(151) + 1

	// Get pokemon name and sprite official artwork from API
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching Pokemon data: %v", err)
		return PokemonWidget{name: "Unknown"}
	}
	defer resp.Body.Close()
	var pokemonData pokemonResponse
	err = json.NewDecoder(resp.Body).Decode(&pokemonData)
	if err != nil {
		log.Printf("Error decoding Pokemon data: %v", err)
		return PokemonWidget{name: "Unknown"}
	}

	capitalisedName := cases.Title(language.English).String(pokemonData.Name)

	// Temp write the image in the URL to file so that image2ascii can read (it doesn't support reading from URLs directly)
	imageResp, err := http.Get(pokemonData.Sprites.Other.OfficialArtwork.FrontDefault)
	if err != nil {
		log.Printf("Error fetching Pokemon image: %v", err)
		return PokemonWidget{name: capitalisedName}
	}
	defer imageResp.Body.Close()

	tempFile := fmt.Sprintf("temp_pokemon_%d.png", id)
	out, err := os.Create(tempFile)
	if err != nil {
		log.Printf("Error creating temp file for Pokemon image: %v", err)
		return PokemonWidget{name: capitalisedName}
	}
	defer out.Close()
	_, err = io.Copy(out, imageResp.Body)
	if err != nil {
		log.Printf("Error saving Pokemon image to temp file: %v", err)
		return PokemonWidget{name: capitalisedName}
	}

	// Clean up the temp file once we have the ASCII string
	defer os.Remove(tempFile)

	// Use image2ascii to convert the sprite to ASCII art
	converter := convert.NewImageConverter()
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 76
	convertOptions.FixedHeight = 38
	asciiArt := converter.ImageFile2ASCIIString(tempFile, &convertOptions)

	return PokemonWidget{
		name:     capitalisedName,
		asciiArt: asciiArt,
		hasError: false,
	}
}

func (p PokemonWidget) View() string {
	pokemonBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.BrightCyan)

	asciiArtStyle := lipgloss.NewStyle()

	return pokemonBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, nameStyle.Render(p.name), asciiArtStyle.Render(p.asciiArt)))
}
