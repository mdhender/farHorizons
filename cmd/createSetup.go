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

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/farHorizons/internal/fh"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type Setup struct {
	Galaxy struct {
		Name  string `json:"name"`
		Setup struct {
			LowDensity            bool `json:"low_density"`
			ForbidNearbyWormholes bool `json:"forbid_nearby_wormholes"`
			MinimumDistance       int  `json:"minimum_distance"`
		} `json:"setup"`
	} `json:"galaxy"`
	Players []Player `json:"players"`
}
type Player struct {
	Email          string `json:"email"`
	SpName         string `json:"species_name"`
	HomePlanetName string `json:"home_planet_name"`
	GovName        string `json:"government_name"`
	GovType        string `json:"government_type"`
	ML             int    `json:"military_level"`
	GV             int    `json:"gravitics_level"`
	LS             int    `json:"life_support_level"`
	BI             int    `json:"biology_level"`
	Load           bool   `json:"load"`
}

// createSetupCmd implements the create setup command
var createSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create a new setup file",
	Long: `This command creates a new setup file ready to be
filled out with player information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("file-name")
		if err != nil {
			return err
		}
		if filename == "" {
			return fmt.Errorf("you must specify a valid file name to create")
		}
		galaxyName, err := cmd.Flags().GetString("galaxy-name")
		if err != nil {
			return err
		}
		if galaxyName == "" {
			return fmt.Errorf("you must specify a valid galaxy name")
		}
		numberOfPlayers, err := cmd.Flags().GetInt("number-of-players")
		if err != nil {
			return err
		}
		if numberOfPlayers < fh.MIN_SPECIES || numberOfPlayers > fh.MAX_SPECIES {
			return fmt.Errorf("number of players must be between %d and %d", fh.MIN_SPECIES, fh.MAX_SPECIES)
		}

		var s Setup
		s.Galaxy.Name = galaxyName
		for i := 0; i < numberOfPlayers; i++ {
			ml, gv, ls, bi := 1, 1, 1, 1
			for k := 5; k <= 15; k++ {
				switch fh.Roll(4) {
				case 1:
					ml++
				case 2:
					gv++
				case 3:
					ls++
				case 4:
					bi++
				}
			}
			s.Players = append(s.Players, Player{
				Email:          fmt.Sprintf("email%02d.example.com", i),
				SpName:         fmt.Sprintf("spName%02d", i),
				HomePlanetName: fmt.Sprintf("hpName%02d", i),
				GovName:        fmt.Sprintf("gName%02d", i),
				GovType:        fmt.Sprintf("gType%02d", i),
				ML:             ml,
				GV:             gv,
				LS:             ls,
				BI:             bi,
				Load:           false,
			})
		}
		if b, err := json.MarshalIndent(s, "  ", "  "); err != nil {
			return err
		} else if err := ioutil.WriteFile(filename, b, 0644); err != nil {
			return err
		}
		fmt.Printf("Created %q.\n", filename)

		return nil
	},
}

func init() {
	createCmd.AddCommand(createSetupCmd)
	createSetupCmd.Flags().StringP("galaxy-name", "g", "", "name of galaxy to be setup")
	_ = createSetupCmd.MarkFlagRequired("galaxy-name")
	createSetupCmd.Flags().StringP("file-name", "f", "", "name of file to create")
	_ = createSetupCmd.MarkFlagRequired("file-name")
	createSetupCmd.Flags().IntP("number-of-players", "n", 0, "number of player entries to create")
	_ = createSetupCmd.MarkFlagRequired("number-of-players")
}
