/*
 * farHorizons - a clone of Far Horizons
 * Copyright (C) 2021  Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package fh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SetupData struct {
	Galaxy struct {
		Name      string `json:"name"`
		Overrides struct {
			UseOverrides  bool `json:"use_overrides"`
			Radius        int  `json:"radius"`
			NumberOfStars int  `json:"number_of_stars"`
		}
		LowDensity            bool `json:"low_density"`
		ForbidNearbyWormholes bool `json:"forbid_nearby_wormholes"`
		MinimumDistance       int  `json:"minimum_distance"`
	} `json:"galaxy"`
	Players []PlayerData `json:"players"`
}

type PlayerData struct {
	Email          string `json:"email"`
	SpeciesName    string `json:"species_name"`
	HomePlanetName string `json:"home_planet_name"`
	GovName        string `json:"government_name"`
	GovType        string `json:"government_type"`
	ML             int    `json:"military_level"`
	GV             int    `json:"gravitics_level"`
	LS             int    `json:"life_support_level"`
	BI             int    `json:"biology_level"`
}

func GetSetup(name string) (*SetupData, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var setup SetupData
	if err := json.Unmarshal(data, &setup); err != nil {
		return nil, err
	}
	emails := make(map[string]bool)
	homePlanetName := make(map[string]bool)
	speciesNames := make(map[string]bool)
	for i, player := range setup.Players {
		if exists := emails[player.Email]; exists {
			return nil, fmt.Errorf("player %d: duplicate email address %q", i+1, player.Email)
		}
		emails[player.Email] = true
		if exists := homePlanetName[player.HomePlanetName]; exists {
			return nil, fmt.Errorf("player %d: duplicate home planet name %q", i+1, player.HomePlanetName)
		}
		homePlanetName[player.HomePlanetName] = true
		if exists := speciesNames[player.SpeciesName]; exists {
			return nil, fmt.Errorf("player %d: duplicate species name %q", i+1, player.SpeciesName)
		}
		speciesNames[player.Email] = true
		if player.BI+player.GV+player.LS+player.ML != 15 {
			return nil, fmt.Errorf("player %d: the tech levels must sum to 15", i+1)
		}
	}
	return &setup, nil
}
