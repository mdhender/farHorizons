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

// createGalaxyCmd implements the create galaxy command
var createGalaxyCmd = &cobra.Command{
	Use:   "galaxy",
	Short: "A galaxy manager?",
	Long: `The command line interface to manage a galaxy.

The number of species must be between ` + fmt.Sprintf("%d and %d", fh.MIN_SPECIES, fh.MAX_SPECIES) + `, inclusive.

The number of stars is based on the number of species, and will be somewhere between ` + fmt.Sprintf("%d and %d", fh.MIN_STARS, fh.MAX_STARS) + `.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		if len(args) != 0 {
			return fmt.Errorf("create galaxy: unknown arguments")
		}
		numberOfSpecies, err := cmd.Flags().GetInt("number-of-species")
		if err != nil {
			return err
		}
		lessCrowded, err := cmd.Flags().GetBool("less-crowded")
		if err != nil {
			return err
		}
		g, err := fh.GenerateGalaxy(numberOfSpecies, lessCrowded)
		if err != nil {
			return err
		}
		return g.Write(name)
	},
}

func init() {
	createCmd.AddCommand(createGalaxyCmd)
	createGalaxyCmd.Flags().IntP("number-of-species", "n", 9, "number of species to create in the galaxy")
	_ = createGalaxyCmd.MarkFlagRequired("number-of-species")
	createGalaxyCmd.Flags().BoolP("less-crowded", "l", false, "create a less crowded galaxy")
}
