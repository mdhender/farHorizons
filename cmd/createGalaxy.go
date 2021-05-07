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
	"os"
	"time"
)

// createGalaxyCmd implements the create galaxy command
var createGalaxyCmd = &cobra.Command{
	Use:   "galaxy",
	Short: "Create a new galaxy",
	Long: `This commands loads setup data from a
configuration file, then creates a new galaxy file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		fh.Seed(0xC0FFEE) // seed random number generator

		galaxyFileName, err := cmd.Flags().GetString("galaxy-file")
		if err != nil {
			return err
		} else if galaxyFileName == "" {
			return fmt.Errorf("you must specify a valid file name to create")
		}
		setupFileName, err := cmd.Flags().GetString("setup-file")
		if err != nil {
			return err
		} else if setupFileName == "" {
			return fmt.Errorf("you must specify a valid setup file name")
		}

		setupData, err := fh.GetSetup(setupFileName)
		if err != nil {
			return err
		}

		// NewGalaxy step in setup_game.py
		g, err := fh.GenerateGalaxy(setupData)
		if err != nil {
			return err
		}

		if setupData.Galaxy.MinimumDistance < 1 || setupData.Galaxy.MinimumDistance > g.Radius*2 {
			return fmt.Errorf("minimum-distance must be between 1 and %d", g.Radius*2)
		}

		// MakeHomes step in setup_game.py
		for num_planets := 3; num_planets < 10; num_planets++ {
			fmt.Printf("Creating home system with %d planets...\n", num_planets)
			var planets []*fh.PlanetData
			for planets == nil {
				planets = fh.GenerateEarthLikePlanet(num_planets)
			}
			g.Templates.Homes[num_planets] = planets
		}

		// skip ListGalaxy step in setup_game.py

		for i, player := range setupData.Players {
			var spec fh.SpeciesData
			spec.Number = i + 1
			spec.Name = player.SpeciesName
			spec.HomeNampla = &fh.NamedPlanetData{}

			home_nampla := spec.HomeNampla
			home_nampla.Name = player.HomePlanetName
			spec.GovtName = player.GovName
			spec.GovtType = player.GovType

			// HomeSystemAuto step in setup_game.py
			forbidNearbyWormholes := setupData.Galaxy.ForbidNearbyWormholes
			minDistance := setupData.Galaxy.MinimumDistance
			x, y, z, err := g.GetFirstXYZ(minDistance, forbidNearbyWormholes)
			if err != nil {
				return err
			}
			// convert the system at those coordinates to a home system
			star, err := g.GetStarAt(x, y, z)
			if err != nil {
				return err
			}
			// fetch the home system template and update the star with values from the template
			star.ConvertToHomeSystem(g.Templates.Homes[star.NumPlanets])
			pn := star.HomePlanetNumber()
			fmt.Printf("Converted system %d %d %d, home planet %d\n", x, y, z, pn)

			// get pointer to home planet
			spec.HomePlanet = star.Planets[star.HomePlanetIndex()]

			// AddSpecies step in setup_game.py
			spec.X, spec.Y, spec.Z = x, y, z
			spec.PN = pn
			home_nampla.X, home_nampla.Y, home_nampla.Z = x, y, z
			home_nampla.PN = pn

			fmt.Printf("Scan of star system:\n\n")
			star.Scan(os.Stdout, nil)
			fmt.Printf("\n")

			/* Check tech levels. */
			totalTechLevels := 0
			totalTechLevels += player.BI
			totalTechLevels += player.GV
			totalTechLevels += player.LS
			totalTechLevels += player.ML
			if totalTechLevels != 15 {
				fmt.Printf("\n\tERROR! ML + GV + LS + BI is not equal to 15!\n\n")
				return fmt.Errorf("total tech levels must sum up to 15")
			}
			// set player-specified tech levels (mining and manufacturing are each 10)
			spec.TechLevel[fh.BI] = player.BI
			spec.TechLevel[fh.GV] = player.GV
			spec.TechLevel[fh.LS] = player.LS
			spec.TechLevel[fh.MA] = 10
			spec.TechLevel[fh.MI] = 10
			spec.TechLevel[fh.ML] = player.ML

			// initialize other tech stuff
			for i := fh.MI; i <= fh.BI; i++ {
				j := spec.TechLevel[i]
				spec.TechKnowledge[i] = j
				spec.InitTechLevel[i] = j
				spec.TechEps[i] = 0
			}

			// confirm that required gas is present
			spec.RequiredGas = fh.O2 // (we're biased towards oxygen breathers?)
			for _, gas := range spec.HomePlanet.Gases {
				if gas.Type == spec.RequiredGas {
					spec.RequiredGasMin = gas.Percentage / 2
					if spec.RequiredGasMin < 1 {
						spec.RequiredGasMin = 1
					}
					spec.RequiredGasMax = 2 * gas.Percentage
					if spec.RequiredGasMax < 20 {
						spec.RequiredGasMax += 20
					} else if spec.RequiredGasMax > 100 {
						// TODO: i prefer 99% for the max
						spec.RequiredGasMax = 100
					}
				}
			}
			if spec.RequiredGasMax == 0 {
				fmt.Printf("\n\tERROR! Planet does not have %s(%s)!\n", spec.RequiredGas.String(), spec.RequiredGas.Char())
				return fmt.Errorf("planet does not have required gas %s", spec.RequiredGas.Char())
			}

			// all home planet gases are either required or neutral
			num_neutral := len(spec.HomePlanet.Gases)
			var goodGas [14]bool
			for _, gas := range spec.HomePlanet.Gases {
				goodGas[gas.Type] = true
			}
			if !goodGas[fh.HE] {
				// Helium must always be neutral since it is a noble gas.
				goodGas[fh.HE] = true
				num_neutral++
			}
			if !goodGas[fh.H2O] {
				// This game is biased towards oxygen breathers, so make H2O neutral also.
				goodGas[fh.H2O] = true
				num_neutral++
			}
			// Start with the good_gas array and add neutral gases until there are exactly seven of them.
			// One of the seven gases will be the required gas.
			for num_neutral < 7 {
				if n := fh.Roll(13); !goodGas[n] {
					goodGas[n] = true
					num_neutral++
				}
			}

			// add the neutral and poison gases
			for n := 1; n <= 13; n++ {
				t := fh.GasType(n)
				if !goodGas[n] {
					spec.PoisonGas = append(spec.PoisonGas, t)
				} else if t != spec.RequiredGas { // required gas isn't neutral!
					spec.NeutralGas = append(spec.NeutralGas, t)
				}
			}

			// Do mining and manufacturing bases of home planet.
			// Initial mining and production capacity will be 25 times sum of MI and MA plus a small random amount.
			// Mining and manufacturing base will be reverse-calculated from the capacity.
			levels := spec.TechLevel[fh.MI] + spec.TechLevel[fh.MA]
			n := (25 * levels) + fh.Roll(levels) + fh.Roll(levels) + fh.Roll(levels)
			home_nampla.MIBase = (n * spec.HomePlanet.MiningDifficulty) / (10 * spec.TechLevel[fh.MI])
			home_nampla.MABase = (10 * n) / spec.TechLevel[fh.MA]

			// initialize contact/ally/enemy masks
			spec.Contact = make([]bool, g.DNumSpecies+1, g.DNumSpecies+1)
			spec.Ally = make([]bool, g.DNumSpecies+1, g.DNumSpecies+1)
			spec.Enemy = make([]bool, g.DNumSpecies+1, g.DNumSpecies+1)

			spec.NumNamplas = 1 // just the home planet for now ("nampla" means "named planet")
			home_nampla.Status = fh.HOME_PLANET | fh.POPULATED
			home_nampla.PopUnits = fh.HP_AVAILABLE_POP
			home_nampla.Shipyards = 1

			/* Print summary. */
			fmt.Printf("\n  Summary for species #%d:\n", spec.Number)
			fmt.Printf("\tName of species: %s\n", spec.Name)
			fmt.Printf("\tName of home planet: %s\n", home_nampla.Name)
			fmt.Printf("\t\tCoordinates: %d %d %d #%d\n", spec.X, spec.Y, spec.Z, spec.PN)
			fmt.Printf("\tName of government: %s\n", spec.GovtName)
			fmt.Printf("\tType of government: %s\n\n", spec.GovtType)

			fmt.Printf("\tTech levels: %s = %d,  %s = %d,  %s = %d\n",
				fh.TechName[fh.MI], spec.TechLevel[fh.MI],
				fh.TechName[fh.MA], spec.TechLevel[fh.MA],
				fh.TechName[fh.ML], spec.TechLevel[fh.ML])
			fmt.Printf("\t             %s = %d,  %s = %d,  %s = %d\n",
				fh.TechName[fh.MI], spec.TechLevel[fh.GV],
				fh.TechName[fh.MA], spec.TechLevel[fh.LS],
				fh.TechName[fh.ML], spec.TechLevel[fh.BI])

			fmt.Printf("\n\n\tFor this species, the required gas is %s (%d%%-%d%%).\n",
				spec.RequiredGas.String(),
				spec.RequiredGasMin, spec.RequiredGasMax)

			fmt.Printf("\tGases neutral to species:")
			for _, gasType := range spec.NeutralGas {
				fmt.Printf(" %s ", gasType.String())
			}

			fmt.Printf("\n\tGases poisonous to species:")
			for _, gasType := range spec.PoisonGas {
				fmt.Printf(" %s ", gasType.String())
			}

			fmt.Printf("\n\n\tInitial mining base = %d.%d. Initial manufacturing base = %d.%d.\n",
				home_nampla.MIBase/10, home_nampla.MIBase%10,
				home_nampla.MABase/10, home_nampla.MABase%10)
			fmt.Printf("\tIn the first turn, %d raw material units will be produced,\n",
				(10*spec.TechLevel[fh.MI]*home_nampla.MIBase)/spec.HomePlanet.MiningDifficulty)
			fmt.Printf("\tand the total production capacity will be %d.\n\n",
				(spec.TechLevel[fh.MA]*home_nampla.MABase)/10)

			// set visited_by bit in star data
			star.VisitedBy[spec.Number] = true

			g.Species[spec.Number] = &spec

			/* Create log file for first turn. Write home star system data to it. */
			logFile := fmt.Sprintf("D:/GoLand/farHorizons/testdata/sp%02d.log", spec.Number)
			w, err := os.Create(logFile)
			if err != nil {
				return err
			}

			fmt.Fprintf(w, "\nScan of home star system for SP %s:\n\n", spec.Name)
			star.Scan(w, &spec)
			fmt.Fprintf(w, "\n")

			fmt.Printf("Created file %q\n", logFile)
		}

		err = g.Write(galaxyFileName)
		if err != nil {
			return err
		}

		fmt.Printf("Created file %q in %v\n", galaxyFileName, time.Now().Sub(started))
		return nil
	},
}

func init() {
	createCmd.AddCommand(createGalaxyCmd)
	createGalaxyCmd.Flags().StringP("galaxy-file", "g", "", "name of galaxy file to create")
	_ = createGalaxyCmd.MarkFlagRequired("galaxy-file")
	createGalaxyCmd.Flags().StringP("setup-file", "i", "", "name of configuration file to load")
	_ = createGalaxyCmd.MarkFlagRequired("setup-file")
}
