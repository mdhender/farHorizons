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
	NumSpecies        int
	Radius            int
	TurnNumber        int
	NumberOfWormHoles int
	Stars             []*StarData
}

func GenerateGalaxy(ns int) (*GalaxyData, error) {
	if ns < MIN_SPECIES || MAX_SPECIES < ns {
		return nil, fmt.Errorf("number of species must be between %d and %d, inclusive", MIN_SPECIES, MAX_SPECIES)
	}

	d_num_species := ns

	/* Get approximate number of star systems to generate. */
	desired_num_stars := (d_num_species * STANDARD_NUMBER_OF_STAR_SYSTEMS) / STANDARD_NUMBER_OF_SPECIES
	if MAX_STARS < desired_num_stars {
		return nil, fmt.Errorf("number of stars must be between 1 and %d, inclusive", MAX_STARS)
	}
	fmt.Printf("For %d species, a game needs about %d stars.\n", d_num_species, desired_num_stars)

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

	/* Seed random number generator. */
	Seed(0xC0FFEE)

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
		DNumSpecies: d_num_species,
		NumSpecies:  0,
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

	numPlanets := 0
	for _, star := range galaxy.Stars {
		numPlanets += len(star.Planets)
	}

	fmt.Printf("This galaxy contains a total of %d stars and %d planets.\n", len(galaxy.Stars), numPlanets)
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
