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
	"fmt"
)

type GalaxyData struct {
	Secret      string
	DNumSpecies int
	NumSpecies  int
	Radius      int
	TurnNumber  int
	Stars       []*StarData
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
	//last_random = time(NULL);
	n := rnd(100) + rnd(200) + rnd(300)
	for i := 0; i < n; i++ {
		rnd(10)
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

	return galaxy, nil
}
