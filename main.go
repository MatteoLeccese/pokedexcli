package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"errors"
	"net/http"
	"time"
	"io"
	"encoding/json"
	"github.com/matteoleccese/pokedexcli/internal/pokecache"
)

// Function config struct to hold the next and previous map locations
type FunctionConfig struct {
	next			string
	previous	string
}

// Command struct to hold the command name and its function
type CliCommand struct {
	name        string
	description string
	callback    func(*FunctionConfig, string) error
}

// Location result slice struct
type LocationResult struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
}

// Location struct from the PokeAPI
type Location struct {
	Count 		int								`json:"count"`
	Next		  string						`json:"next"`
	Previous	string						`json:"previous"`
	Results		[]LocationResult	`json:"results"`
}

// EncounterLocation struct
type EncounterLocation struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
}

// Location struct
type EncounterVersionDetails struct {
	Rate			int								`json:"rate"`
	Version   EncounterLocation	`json:"version"`
}

// EncounterMethodRates struct
type EncounterMethodRates struct {
	Encounter_method	EncounterLocation					`json:"encounter_method"`
	Version_details		[]EncounterVersionDetails	`json:"version_details"`
}

// EncounterNames struct
type EncounterNames struct {
	Language	[]EncounterLocation `json:"language"`
	Name			string							`json:"name"`
}

// EncounterDetails struct
type EncounterDetails struct {
	Chance						int									`json:"chance"`
	Condition_values	[]interface{}				`json:"condition_values"`
	Max_level		      int									`json:"max_level"`
	Method						EncounterLocation		`json:"method"`
	Min_level					int									`json:"min_level"`
}

// VersionDetails struct
type VersionDetails struct {
	Encounter_details	[]EncounterDetails	`json:"encounter_details"`
	Max_chance				int									`json:"max_chance"`
	Version						EncounterLocation		`json:"version"`
}

// PokemonEncounters struct
type PokemonEncounters struct {
	Pokemon					EncounterLocation	`json:"pokemon"`
	Version_details []VersionDetails	`json:"version_details"`
}

// Pokemon Location Area struct
type LocationArea struct {
	Encounter_method_rates	[]EncounterMethodRates	`json:"encounter_method_rates"`
	Game_index							int											`json:"game_index"`
	Id											int 										`json:"id"`
	Location								EncounterLocation				`json:"location"`
	Name										string									`json:"name"`
	Names										[]EncounterNames				`json:"names"`
	Pokemon_encounters			[]PokemonEncounters			`json:"pokemon_encounters"`
}

// Constants for the PokeAPI URL
const POKE_API_URL = "https://pokeapi.co/api/v2/location-area";

// Function to clean the input string
func cleanInput(test string) []string {
	// Parsed string
	test = strings.ToLower(test);

	// Slice to hold the cleaned strings
	slice := strings.Fields(test);

	// slice := make([]string, 0);

	return slice;
}

func commandExit(fc *FunctionConfig, argument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!");

	os.Exit(0);
	return nil;
}

func commandHelp(fc *FunctionConfig, argument string) error {
	multiString := `
	Welcome to the Pokedex!
	Usage:

	help: Displays a help message
	exit: Exit the Pokedex
	`;

	// Print the help message
	fmt.Println(multiString);

	return nil;
}

func commandMap(fc *FunctionConfig, argument string) error {
	// Getting the next map location
	url := fc.next;

	// Create a new cache with a 5 second interval
	const cacheInterval = time.Second * 5;

	// Create a new cache
	cache := pokecache.NewCache(cacheInterval);

	// Checking if the URL is in the cache
	cacheData, ok := cache.Get(url);

	if ok {
		// Initialize a new Location struct
		location := Location{}

		// Unmarshal the JSON response into the Location struct
		err := json.Unmarshal(cacheData, &location)

		// Check if there was an error unmarshalling the JSON
		if err != nil {
			errors.New("Error unmarshalling JSON: " + err.Error())
		}

		// Update the FunctionConfig struct with the next and previous map locations
		fc.next = location.Next;
		fc.previous = location.Previous;

		// Now we access to the results
		for _, result := range location.Results {
			// Print the name of the location
			fmt.Println(result.Name);
		}

		return nil;
	}

	// If the next map location is empty, set it to the PokeAPI URL
	if url == "" {
		url = POKE_API_URL;
	}

	// Make a GET request to the PokeAPI URL
	res, err := http.Get(url)

	// Check if the request was successful
	if err != nil {
		errors.New("Error making GET request: " + err.Error())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	// Close the response body
	defer res.Body.Close()

	// Check if there was an error reading the response body
	if err != nil {
		errors.New("Error reading response body: " + err.Error())
	}

	// Initialize a new Location struct
	location := Location{}

	// Add the data to the cache
	cache.Add(url, body);

	// Unmarshal the JSON response into the Location struct
	err = json.Unmarshal(body, &location)

	// Check if there was an error unmarshalling the JSON
	if err != nil {
		errors.New("Error unmarshalling JSON: " + err.Error())
	}

	// Update the FunctionConfig struct with the next and previous map locations
	fc.next = location.Next;
	fc.previous = location.Previous;

	// Now we access to the results
	for _, result := range location.Results {
		// Print the name of the location
		fmt.Println(result.Name);
	}

	return nil;
}

func commandMapB(fc *FunctionConfig, argument string) error {
	// Getting the next map location
	url := fc.previous;

	// Create a new cache with a 5 second interval
	const cacheInterval = time.Second * 5;

	// Create a new cache
	cache := pokecache.NewCache(cacheInterval);

	// Checking if the URL is in the cache
	cacheData, ok := cache.Get(url);

	if ok {
		// Initialize a new Location struct
		location := Location{}

		// Unmarshal the JSON response into the Location struct
		err := json.Unmarshal(cacheData, &location)

		// Check if there was an error unmarshalling the JSON
		if err != nil {
			errors.New("Error unmarshalling JSON: " + err.Error())
		}

		// Update the FunctionConfig struct with the next and previous map locations
		fc.next = location.Next;
		fc.previous = location.Previous;

		// Now we access to the results
		for _, result := range location.Results {
			// Print the name of the location
			fmt.Println(result.Name);
		}

		return nil;
	}

	// If the next map location is empty, set it to the PokeAPI URL
	if url == "" {
		url = POKE_API_URL;
	}

	// Make a GET request to the PokeAPI URL
	res, err := http.Get(url)

	// Check if the request was successful
	if err != nil {
		errors.New("Error making GET request: " + err.Error())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	// Close the response body
	defer res.Body.Close()

	// Check if there was an error reading the response body
	if err != nil {
		errors.New("Error reading response body: " + err.Error())
	}

	// Initialize a new Location struct
	location := Location{}

	// Add the data to the cache
	cache.Add(url, body);

	// Unmarshal the JSON response into the Location struct
	err = json.Unmarshal(body, &location)

	// Check if there was an error unmarshalling the JSON
	if err != nil {
		errors.New("Error unmarshalling JSON: " + err.Error())
	}

	// Update the FunctionConfig struct with the next and previous map locations
	fc.next = location.Next;
	fc.previous = location.Previous;

	// Now we access to the results
	for _, result := range location.Results {
		// Print the name of the location
		fmt.Println(result.Name);
	}

	return nil;
}

// Function to get the list of all the Pokémon located on a given area
func commandExplore(fc *FunctionConfig, argument string) error {
	// Making the url to the PokeAPI
	url := fmt.Sprintf("%s/%s", POKE_API_URL, argument);

	// Create a new cache with a 5 second interval
	// const cacheInterval = time.Second * 5;

	// // Create a new cache
	// cache := pokecache.NewCache(cacheInterval);

	// // Checking if the URL is in the cache
	// cacheData, ok := cache.Get(url);

	// if ok {
	// 	// CODE HERE
	// }

	// Make a GET request to the PokeAPI URL
	res, err := http.Get(url)

	// Check if the request was successful
	if err != nil {
		errors.New("Error making GET request: " + err.Error())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	// Close the response body
	defer res.Body.Close()

	// Check if there was an error reading the response body
	if err != nil {
		errors.New("Error reading response body: " + err.Error())
	}

	// Initialize a new Location struct
	locationArea := LocationArea{}

	// Add the data to the cache
	// cache.Add(url, body);

	// Unmarshal the JSON response into the Location struct
	err = json.Unmarshal(body, &locationArea)

	// Check if there was an error unmarshalling the JSON
	if err != nil {
		errors.New("Error unmarshalling JSON: " + err.Error())
	}

	// Telling the user what we are exploring
	fmt.Println("Exploring %s...", argument);

	pokemonSlice := locationArea.Pokemon_encounters;

	if len(pokemonSlice) == 0 {
		fmt.Println("No Pokémon found in this area.");
		return nil;
	}

	fmt.Println("Found Pokemon:");

	for _, pokemonEncounter := range pokemonSlice {
		// Check if Pokemon is nil
		if pokemonEncounter.Pokemon.Name == "" {
			continue;
			// fmt.Println("No Pokemon found in this area.");
			// return nil;
		}

		// Print the name of the locationArea
		fmt.Println("- %s", pokemonEncounter.Pokemon.Name);
	}
	// // Now we access to the results
	// for _, pokemonEncounter := range pokemonSlice {
	// 	// Print the name of the locationArea
	// 	fmt.Println("- %s", pokemonEncounter.Pokemon.Name);
	// }

	return nil;
}

// Main function of the program
func main() {

	// Create a new FunctionConfig struct
	fc := &FunctionConfig{
		next: "",
		previous: "",
	};

	// List of supported commands
	supportedCommands := map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the name of the next 20 map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the name of the previous 20 map locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays the list of all the Pokémon located on a given area.",
			callback:    commandExplore,
		},
	}

	// Create a new scanner to read from the input
	scanner := bufio.NewScanner(os.Stdin);

	// Infinite loop to read input until "exit" is entered
	for {
		// Prompt the user for input
		fmt.Print("Pokedex > ");

		// Read the input
		scanner.Scan();
		// Get the text from the scanner
		test := scanner.Text();

		// Clean the user input
		clean_input := cleanInput(test);

		// Check if the input is "exit"
		user_command := clean_input[0];

		// Variable to hold the argument
		argument := "";

		// Check if there is a second argument
		if len(clean_input) > 1 {
			argument = clean_input[1];
		}

		// Check if the command is supported
		command, supported_command := supportedCommands[user_command];

		// If the command is supported, call the callback function
		if !supported_command {
			fmt.Println("Unknown command.");
			return; // Return early if the command is not supported
		}

		// Call the command's callback function
		command.callback(fc, argument);
	}
}
