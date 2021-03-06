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
		forbidNearbyWormholes, err := cmd.Flags().GetBool("forbid-nearby-wormholes")
		if err != nil {
			return err
		}
		galaxyName, err := cmd.Flags().GetString("galaxy-name")
		if err != nil {
			return err
		}
		if galaxyName == "" {
			return fmt.Errorf("you must specify a valid galaxy name")
		}
		lowDensity, err := cmd.Flags().GetBool("low-density")
		if err != nil {
			return err
		}
		minDistance, err := cmd.Flags().GetInt("minimum-distance")
		if err != nil {
			return err
		}
		numberOfPlayers, err := cmd.Flags().GetInt("number-of-players")
		if err != nil {
			return err
		}
		if numberOfPlayers < fh.MIN_SPECIES || numberOfPlayers > fh.MAX_SPECIES {
			return fmt.Errorf("number of players must be between %d and %d", fh.MIN_SPECIES, fh.MAX_SPECIES)
		}

		var s fh.SetupData
		s.Galaxy.Name = galaxyName
		s.Galaxy.ForbidNearbyWormholes = forbidNearbyWormholes
		s.Galaxy.LowDensity = lowDensity
		s.Galaxy.MinimumDistance = minDistance
		for i := 1; i <= numberOfPlayers; i++ {
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
			s.Players = append(s.Players, fh.PlayerData{
				Email:          fmt.Sprintf("email%02d.example.com", i),
				SpeciesName:    fmt.Sprintf("spName%02d", i),
				HomePlanetName: fmt.Sprintf("hpName%02d", i),
				GovName:        fmt.Sprintf("gName%02d", i),
				GovType:        fmt.Sprintf("gType%02d", i),
				ML:             ml,
				GV:             gv,
				LS:             ls,
				BI:             bi,
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
	createSetupCmd.Flags().Bool("forbid-nearby-wormholes", false, "forbid wormholes to be neighbors")
	createSetupCmd.Flags().Bool("low-density", false, "increase the radius by 50%")
	createSetupCmd.Flags().IntP("minimum-distance", "d", 10, "minimum distance between home systems")
}
