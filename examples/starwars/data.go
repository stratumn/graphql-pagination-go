package starwars

import "fmt"

/*
This defines a basic set of data for our Star Wars Schema.

This data is hard coded for the sake of the demo, but you could imagine
fetching this data from a backend service rather than from hardcoded
JSON objects in a more complex demo.
*/

// Ship describes a ship
type Ship struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Faction describes a faction
type Faction struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Ships []string `json:"ships"`
}

// XWing is a X-Wing
var XWing = &Ship{"1", "X-Wing"}

// YWing is a Y-Wing
var YWing = &Ship{"2", "Y-Wing"}

// AWing is a A-Wing
var AWing = &Ship{"3", "A-Wing"}

// Yeah, technically it's Corellian. But it flew in the service of the rebels,
// so for the purposes of this demo it's a rebel ship.

// Falcon is the Millenium Falcon
var Falcon = &Ship{"4", "Millenium Falcon"}

// HomeOne is the Home One
var HomeOne = &Ship{"5", "Home One"}

// TIEFighter is a TIE Fighter
var TIEFighter = &Ship{"6", "TIE Fighter"}

// TIEInterceptor is a TIE Interceptor
var TIEInterceptor = &Ship{"7", "TIE Interceptor"}

// Executor is a Executor
var Executor = &Ship{"8", "Executor"}

// Rebels is the rebels faction
var Rebels = &Faction{
	"1",
	"Alliance to Restore the Republic",
	[]string{"1", "2", "3", "4", "5"},
}

// Empire is the empire faction
var Empire = &Faction{
	"2",
	"Galactic Empire",
	[]string{"6", "7", "8"},
}

var factions = map[string]*Faction{
	"1": Rebels,
	"2": Empire,
}
var ships = map[string]*Ship{
	"1": XWing,
	"2": YWing,
	"3": AWing,
	"4": Falcon,
	"5": HomeOne,
	"6": TIEFighter,
	"7": TIEInterceptor,
	"8": Executor,
}
var nextShip = 9

// CreateShip creates a ship
func CreateShip(shipName string, factionID string) *Ship {
	nextShip = nextShip + 1
	newShip := &Ship{
		fmt.Sprintf("%v", nextShip),
		shipName,
	}
	ships[newShip.ID] = newShip

	faction := GetFaction(factionID)
	if faction != nil {
		faction.Ships = append(faction.Ships, newShip.ID)
	}
	return newShip
}

// GetShip gets a ship
func GetShip(id string) *Ship {
	if ship, ok := ships[id]; ok {
		return ship
	}
	return nil
}

// GetFaction gets a faction
func GetFaction(id string) *Faction {
	if faction, ok := factions[id]; ok {
		return faction
	}
	return nil
}

// GetRebels gets the rebels faction
func GetRebels() *Faction {
	return Rebels
}

// GetEmpire gets the empire faction
func GetEmpire() *Faction {
	return Empire
}
