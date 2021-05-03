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

type GalaxyData struct {
	Secret            string
	DNumSpecies       int
	LessCrowded       bool
	NumSpecies        int
	Radius            int
	NumberOfStars     int
	NumberOfWormHoles int
	NumberOfPlanets   int
	TurnNumber        int
	Stars             []*StarData
}

func GenerateGalaxy(ns int, lessCrowded bool) (*GalaxyData, error) {
	// seed random number generator
	Seed(0xC0FFEE)

	if ns < MIN_SPECIES || MAX_SPECIES < ns {
		return nil, fmt.Errorf("number of species must be between %d and %d, inclusive", MIN_SPECIES, MAX_SPECIES)
	}
	d_num_species := ns
	if lessCrowded {
		// add 50% more species to the mix as a way to trick the program into adding more stars
		d_num_species = (ns * 3) / 2
		if MAX_SPECIES < d_num_species {
			d_num_species = MAX_SPECIES
		}
	}

	/* Get approximate number of star systems to generate. */
	desired_num_stars := (d_num_species * STANDARD_NUMBER_OF_STAR_SYSTEMS) / STANDARD_NUMBER_OF_SPECIES
	if lessCrowded {
		desired_num_stars = (desired_num_stars * 3) / 2
	}
	if MAX_STARS < desired_num_stars {
		return nil, fmt.Errorf("number of stars must be between 1 and %d, inclusive", MAX_STARS)
	}
	if lessCrowded {
		fmt.Printf("For %d species, a less crowded game needs about %d stars.\n", d_num_species, desired_num_stars)
	} else {
		fmt.Printf("For %d species, a game needs about %d stars.\n", d_num_species, desired_num_stars)
	}

	/* Get size of galaxy to generate. */
	volume := desired_num_stars * STANDARD_GALACTIC_RADIUS * STANDARD_GALACTIC_RADIUS * STANDARD_GALACTIC_RADIUS / STANDARD_NUMBER_OF_STAR_SYSTEMS
	galactic_radius := 2
	for galactic_radius*galactic_radius*galactic_radius < volume {
		galactic_radius++
	}
	fmt.Printf("For %d stars, the galaxy should have a radius of about %d parsecs.\n", desired_num_stars, galactic_radius)
	galactic_diameter := 2 * galactic_radius

	/* Get the number of cubic parsecs within a sphere with a radius of galactic_radius parsecs. */
	volume = (4 * 314 * galactic_radius * galactic_radius * galactic_radius) / 300

	/* The probability of a star system existing at any particular set of x,y,z coordinates is one in chance_of_star. */
	chance_of_star := volume / desired_num_stars
	if chance_of_star < 50 {
		return nil, fmt.Errorf("galactic radius is too small for %d stars", desired_num_stars)
	} else if chance_of_star > 3200 {
		return nil, fmt.Errorf("galactic radius is too large for %d stars", desired_num_stars)
	}

	/* Initialize star location data. */
	var star_here [MAX_DIAMETER][MAX_DIAMETER]int
	for x := 0; x < galactic_diameter; x++ {
		for y := 0; y < galactic_diameter; y++ {
			star_here[x][y] = -1
		}
	}

	// randomly place stars
	for num_stars := 0; num_stars < desired_num_stars; {
		// generate coordinates randomly
		x, y, z := rnd(galactic_diameter)-1, rnd(galactic_diameter)-1, rnd(galactic_diameter)-1
		// verify the coordinates are within the galactic boundary
		real_x, real_y, real_z := x-galactic_radius, y-galactic_radius, z-galactic_radius
		sq_distance_from_center := (real_x * real_x) + (real_y * real_y) + (real_z * real_z)
		if sq_distance_from_center >= galactic_radius*galactic_radius {
			continue
		}
		// verify that we don't already have a star here
		if star_here[x][y] != -1 {
			continue
		}
		// add the star at these coordinates
		star_here[x][y] = z /* z-coordinate. */
		num_stars++
	}

	galaxy := &GalaxyData{
		Secret:      "your-private-key-belongs-here",
		DNumSpecies: ns,
		LessCrowded: lessCrowded,
		Radius:      galactic_radius,
		TurnNumber:  0,
	}

	for x := 0; x < galactic_diameter; x++ {
		for y := 0; y < galactic_diameter; y++ {
			// verify that we have a star at these coordinates
			z := star_here[x][y]
			if z == -1 {
				continue
			}

			star, err := GenerateStar(x, y, z)
			if err != nil {
				return nil, err
			}
			galaxy.Stars = append(galaxy.Stars, star)
		}
	}
	galaxy.NumberOfStars = len(galaxy.Stars)

	// generate natural wormholes
	minWormholeLength := 20 // galactic_radius + 3 // in parsecs
	//if minWormholeLength > 20 {
	//	minWormholeLength = 20
	//}
	for _, star := range galaxy.Stars {
		if star.HomeSystem || star.WormHere || rnd(100) < 92 {
			continue
		}

		// we want to put a wormhole here if we can find a star at least that minimum distance away that doesn't already have a worm hole
		var worm_star *StarData
		for k, f := 0, rnd(desired_num_stars); k < desired_num_stars && worm_star == nil; k++ {
			ps := galaxy.Stars[(k+f)%desired_num_stars]
			if ps == star || ps.HomeSystem || ps.WormHere {
				continue
			}
			// eliminate wormholes less than the minimum
			dx, dy, dz := star.X-ps.X, star.Y-ps.Y, star.Z-ps.Z
			if distance_squared := (dx * dx) + (dy * dy) + (dz * dz); distance_squared < minWormholeLength*minWormholeLength {
				continue
			}
			worm_star = ps
		}
		if worm_star == nil {
			// wow. none of the existing stars met the criteria
			continue
		}

		star.WormHere = true
		star.WormX, star.WormY, star.WormZ = worm_star.X, worm_star.Y, worm_star.Z

		worm_star.WormHere = true
		worm_star.WormX, worm_star.WormY, worm_star.WormZ = star.X, star.Y, star.Z

		// todo: consider making a number of the wormholes one-way
		galaxy.NumberOfWormHoles++
	}

	for _, star := range galaxy.Stars {
		galaxy.NumberOfPlanets += len(star.Planets)
	}

	fmt.Printf("This galaxy contains a total of %d stars and %d planets.\n", len(galaxy.Stars), galaxy.NumberOfPlanets)
	if galaxy.NumberOfWormHoles == 1 {
		fmt.Printf("The galaxy contains %d natural wormhole.\n\n", galaxy.NumberOfWormHoles)
	} else {
		fmt.Printf("The galaxy contains %d natural wormholes.\n\n", galaxy.NumberOfWormHoles)
	}

	return galaxy, nil
}

// GetGalaxy loads data from a JSON file.
func GetGalaxy(name string) (*GalaxyData, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var galaxy GalaxyData
	if err := json.Unmarshal(data, &galaxy); err != nil {
		return nil, err
	}
	return &galaxy, nil
}

func (g *GalaxyData) List(listPlanets, listWormholes bool) error {
	// initialize counts
	total_planets := 0
	total_wormstars := 0
	var type_count [10]int
	for i := DWARF; i <= GIANT; i++ {
		type_count[i] = 0
	}

	// for each star, list info
	for star_index, star := range g.Stars {
		if !listWormholes {
			if listPlanets {
				fmt.Printf("System #%d:\t", star_index+1)
			}
			fmt.Printf("x = %d\ty = %d\tz = %d", star.X, star.Y, star.Z)
			fmt.Printf("\tstellar type = %s%s%s", star.Type.Char(), star.Color.Char(), StarSizeChar[star.Size])
			if listPlanets {
				fmt.Printf("\t%d planets.", star.NumPlanets)
			}
			fmt.Printf("\n")

			if star.NumPlanets == 0 {
				fmt.Printf("\tStar #%d went nova!", star_index+1)
				fmt.Printf(" All planets were blown away!\n")
			} else if star.NumPlanets != len(star.Planets) {
				return fmt.Errorf("assert(numPlanets == lenPlanets)")
			}
		}

		total_planets += star.NumPlanets
		type_count[star.Type] += 1

		if star.WormHere {
			total_wormstars++
			if listPlanets {
				fmt.Printf("!!! Natural wormhole from here to %d %d %d\n", star.WormX, star.WormY, star.WormZ)
			} else if listWormholes {
				fmt.Printf("Wormhole #%d: from %d %d %d to %d %d %d\n", total_wormstars, star.X, star.Y, star.Z, star.WormX, star.WormY, star.WormZ)
				// turn off the target's worm flag to avoid double-reporting
				for _, worm_star := range g.Stars {
					if star.WormX == worm_star.X && star.WormY == worm_star.Y && star.WormZ == worm_star.Z {
						worm_star.WormHere = false
						break
					}
				}
			}
		}

		var home_planet *PlanetData
		if listPlanets {
			/* Check if system has a home planet. */
			for _, planet := range star.Planets {
				if planet.Special == IDEAL_HOME_PLANET || planet.Special == IDEAL_COLONY_PLANET {
					home_planet = planet
					break
				}
			}
		}

		if listPlanets {
			for i, planet := range star.Planets {
				switch planet.Special {
				case NOT_SPECIAL:
					fmt.Printf("     ")
				case IDEAL_HOME_PLANET:
					fmt.Printf(" HOM ")
				case IDEAL_COLONY_PLANET:
					fmt.Printf(" COL ")
				case RADIOACTIVE_HELLHOLE:
					fmt.Printf("     ")
				}
				fmt.Printf("#%d dia=%3d g=%d.%02d tc=%2d pc=%2d md=%d.%02d", i,
					planet.Diameter,
					planet.Gravity/100,
					planet.Gravity%100,
					planet.TemperatureClass,
					planet.PressureClass,
					planet.MiningDifficulty/100,
					planet.MiningDifficulty%100)

				if home_planet != nil {
					fmt.Printf("%4d ", LSN(planet, home_planet))
				} else {
					fmt.Printf("  ")
				}

				num_gases := len(planet.Gases)
				for i, gas := range planet.Gases {
					if gas.Percentage > 0 {
						if i > 0 {
							fmt.Printf(",")
						}
						fmt.Printf("%s(%d%%)", gas.Type.Char(), gas.Percentage)
					}
				}
				if num_gases == 0 {
					fmt.Printf("No atmosphere")
				}

				fmt.Printf("\n")
			}
		}

		if listPlanets {
			fmt.Printf("\n")
		}
	}

	if !listWormholes {
		fmt.Printf("The galaxy has a radius of %d parsecs.\n", g.Radius)
		fmt.Printf("It contains %d dwarf stars, %d degenerate stars, %d main sequence stars,\n", type_count[DWARF], type_count[DEGENERATE], type_count[MAIN_SEQUENCE])
		fmt.Printf("    and %d giant stars, for a total of %d stars.\n", type_count[GIANT], g.NumberOfStars)
		if listPlanets {
			fmt.Printf("The total number of planets in the galaxy is %d.\n", total_planets)
			fmt.Printf("The total number of natural wormholes in the galaxy is %d.\n", total_wormstars/2)
			fmt.Printf("The galaxy was designed for %d species.\n", g.DNumSpecies)
			fmt.Printf("A total of %d species have been designated so far.\n\n", g.NumSpecies)
		}
	}

	/* Internal test. */
	if g.NumberOfPlanets != total_planets {
		return fmt.Errorf("WARNING!  Program error!  Internal inconsistency!")
	}

	return nil
}
