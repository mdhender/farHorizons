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

// convertCmd implements the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a system to a Home System",
	Long: `Converts star system to one suitable for a home system.
The star system may be specified by giving the X, Y, Z co-ordinates,
or one may be picked at random. If picking a random system, the program
that ensure that is at least a certain distance from all other home
systems.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allSystems, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		forbidNearbyWormholes, err := cmd.Flags().GetBool("forbid-nearby-wormholes")
		if err != nil {
			return err
		}
		oneSystem, err := cmd.Flags().GetBool("one")
		if err != nil {
			return err
		}
		if allSystems && oneSystem {
			return fmt.Errorf("specify either all or one, not both")
		}
		reset, err := cmd.Flags().GetBool("reset")
		if err != nil {
			return err
		}
		addUpTo, err := cmd.Flags().GetInt("add-up-to")
		if err != nil {
			return err
		} else if addUpTo < 0 {
			addUpTo = 0
		}
		minDistance, err := cmd.Flags().GetInt("minimum-distance")
		if err != nil {
			return err
		}
		x, err := cmd.Flags().GetInt("x-origin")
		if err != nil {
			return err
		}
		y, err := cmd.Flags().GetInt("y-origin")
		if err != nil {
			return err
		}
		z, err := cmd.Flags().GetInt("z-origin")
		if err != nil {
			return err
		}
		if allSystems && (x != -1 || y != -1 || z != -1) {
			return fmt.Errorf("specify either all or co-ordinates, not both")
		}
		if oneSystem && (x != -1 || y != -1 || z != -1) {
			return fmt.Errorf("specify either one or co-ordinates, not both")
		}

		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		g, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}

		if minDistance < 1 || minDistance > g.Radius*2 {
			return fmt.Errorf("minimum-distance must be between 1 and %d", g.Radius*2)
		}

		// seed random number generator
		fh.Seed(0xC0FFEE)

		if reset {
			for _, star := range g.Stars {
				if star.HomeSystem {
					star.HomeSystem = false
				}
			}
		}

		systemsToConvert := 1
		if addUpTo != 0 {
			systemsToConvert = addUpTo
			for _, star := range g.Stars {
				if star.HomeSystem {
					systemsToConvert--
				}
			}
		} else if allSystems {
			systemsToConvert = g.DNumSpecies
			for _, star := range g.Stars {
				if star.HomeSystem {
					systemsToConvert--
				}
			}
		}

		systemsConverted := 0
		for ; systemsToConvert > 0; systemsToConvert-- {
			if oneSystem || allSystems || addUpTo != 0 {
				x, y, z, err = g.GetFirstXYZ(minDistance, forbidNearbyWormholes)
				if err != nil {
					return err
				}
			}

			// convert the system at the given coordinates
			fmt.Printf("Converting system %d %d %d\n", x, y, z)
			star, err := g.GetStarAt(x, y, z)
			if err != nil {
				return err
			}

			// fetch the home system template and update the star with values from the template
			star.ConvertToHomeSystem(g.Templates.Homes[star.NumPlanets])
			star.HomeSystem = true
			fmt.Printf("Converted system %d %d %d, home planet %d\n", x, y, z, star.HomePlanetNumber())
			systemsConverted++
		}

		fmt.Printf("Converted %d systems.\n", systemsConverted)
		return g.Write(name)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().Bool("all", false, "convert randomly picked systems, up to the species limit")
	convertCmd.Flags().Bool("forbid-nearby-wormholes", false, "forbid wormholes to be neighbors")
	convertCmd.Flags().Bool("one", false, "convert one randomly picked system")
	convertCmd.Flags().Bool("reset", false, "reset existing home systems first")
	convertCmd.Flags().IntP("add-up-to", "n", 0, "add up to a maximum number of systems")
	convertCmd.Flags().IntP("minimum-distance", "d", 10, "minimum distance between home systems")
	convertCmd.Flags().IntP("x-origin", "x", -1, "x coordinate of system to convert")
	convertCmd.Flags().IntP("y-origin", "y", -1, "y coordinate of system to convert")
	convertCmd.Flags().IntP("z-origin", "z", -1, "z coordinate of system to convert")
}
