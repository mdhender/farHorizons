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

import "fmt"

type StarData struct {
	X, Y, Z             int                    /* Coordinates. */
	Type                StarType               /* Dwarf, degenerate, main sequence or giant. */ // was `type`
	Color               StarColor              /* Star color. Blue, blue-white, etc. */
	Size                int                    /* Star size, from 0 thru 9 inclusive. */
	NumPlanets          int                    /* Number of usable planets in star system. */
	HomeSystem          bool                   /* TRUE if this is a good potential home system. */
	WormHere            bool                   /* TRUE if wormhole entry/exit. */
	WormX, WormY, WormZ int                    /* Coordinates. */
	Message             int                    /* Message associated with this star system, if any. */
	VisitedBy           [NUM_CONTACT_WORDS]int /* A bit is set if corresponding species has been here. */
	PlanetIndex         int                    /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	Planets             []*PlanetData
}

func GenerateStar(x, y, z int) (*StarData, error) {
	fmt.Printf("Generating star (%3d, %3d, %3d)\n", x, y, z)

	/* Set coordinates. */
	star := &StarData{
		X:           x,
		Y:           y,
		Z:           z,
		NumPlanets:  -2, // default value to initialize the planet generator
		PlanetIndex: -1,
	}

	/* Determine type of star. Make MAIN_SEQUENCE the most common star type. */
	switch rnd(GIANT + 6) {
	case 1:
		star.Type = DWARF
	case 2:
		star.Type = DEGENERATE
	case 3:
		star.Type = GIANT
	default:
		star.Type = MAIN_SEQUENCE
	}

	/* Determine the number of planets in orbit around the star. The algorithm is something I tweaked until I liked it. It's weird, but it works. */
	/* Color and size of star are totally random. */
	star.Size = rnd(10) - 1
	switch c := rnd(RED); c {
	case BLUE:
		star.Color = BLUE
	case BLUE_WHITE:
		star.Color = BLUE_WHITE
	case WHITE:
		star.Color = WHITE
	case YELLOW_WHITE:
		star.Color = YELLOW_WHITE
	case YELLOW:
		star.Color = YELLOW
	case ORANGE:
		star.Color = ORANGE
	case RED:
		star.Color = RED
	default:
		return nil, fmt.Errorf("assert(StarColor != %d)", c)
	}

	/* Size of die. Big stars (blue, blue-white) roll bigger dice. Smaller stars (orange, red) roll smaller dice. */
	var sizeOfDie int
	switch star.Color {
	case BLUE:
		sizeOfDie = 8
	case BLUE_WHITE:
		sizeOfDie = 7
	case WHITE:
		sizeOfDie = 6
	case YELLOW_WHITE:
		sizeOfDie = 5
	case YELLOW:
		sizeOfDie = 4
	case ORANGE:
		sizeOfDie = 3
	case RED:
		sizeOfDie = 2
	}

	/* Number of rolls: dwarves have 1 roll, degenerates and main sequence stars have 2 rolls, and giants have 3 rolls. */
	var numberOfDice int
	switch star.Type {
	case DWARF:
		numberOfDice = 1
	case DEGENERATE:
		numberOfDice = 2
	case MAIN_SEQUENCE:
		numberOfDice = 2
	case GIANT:
		numberOfDice = 3
	default:
		panic(fmt.Sprintf("assert(star.Type != %d)", star.Type))
	}

	for i := 1; i <= numberOfDice; i++ {
		star.NumPlanets += rnd(sizeOfDie)
	}
	// adjust if too few or too many planets
	for star.NumPlanets < 1 {
		star.NumPlanets += rnd(2)
	}
	for star.NumPlanets > 9 {
		star.NumPlanets -= rnd(3)
	}

	fmt.Printf("Generating star (%3d, %3d, %3d) (type %-13s) (planets %d)\n", x, y, z, star.Type, star.NumPlanets)

	// generate planets
	var err error
	star.Planets, err = GeneratePlanet(star.NumPlanets)
	if err != nil {
		return nil, err
	}

	return star, nil
}
