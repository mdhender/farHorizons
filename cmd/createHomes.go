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

// createHomesCmd implements the create homes command
var createHomesCmd = &cobra.Command{
	Use:   "homes",
	Short: "Create home systems",
	Long: `This command creates the set of templates used to populate systems
that have a home planet. It randomly populates a template for systems
containing from 3 to 9 planets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		g, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}

		// seed random number generator
		fh.Seed(0xC0FFEE)

		for num_planets := 3; num_planets < 10; num_planets++ {
			fmt.Printf("Creating home system with %d planets...\n", num_planets)
			var planets []*fh.PlanetData
			for planets == nil {
				planets = fh.GenerateEarthLikePlanet(num_planets)
			}
			g.Templates.Homes[num_planets] = planets
		}

		return g.Write(name)
	},
}

func init() {
	createCmd.AddCommand(createHomesCmd)
}
