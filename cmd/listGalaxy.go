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
	"fmt"
	"github.com/mdhender/farHorizons/internal/fh"
	"github.com/spf13/cobra"
)

// listGalaxyCmd implements the list galaxy command
var listGalaxyCmd = &cobra.Command{
	Use:   "galaxy",
	Short: "list galaxy properties",
	Long:  `List details about the galaxy, writing the results to stdout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		if len(args) != 0 {
			return fmt.Errorf("list galaxy: unknown arguments")
		}

		listPlanets, listWormholes := true, false
		noListPlanets, err := cmd.Flags().GetBool("no-list-planets")
		if err != nil {
			return err
		}
		if noListPlanets {
			listPlanets = false
		}
		onlyWormholes, err := cmd.Flags().GetBool("only-wormholes")
		if err != nil {
			return err
		}
		if onlyWormholes {
			listPlanets = false
			listWormholes = true
		}

		// load all the data
		galaxy, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}
		err = galaxy.List(listPlanets, listWormholes)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listGalaxyCmd)
	listGalaxyCmd.Flags().BoolP("no-list-planets", "p", false, "do not list planets")
	listGalaxyCmd.Flags().BoolP("only-wormholes", "w", false, "list only wormholes")
}
