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

func (s *StarData) At(x, y, z int) bool {
	return s != nil && s.X == x && s.Y == y && s.Z == z
}

func (s *StarData) ConvertToHomeSystem(src []*PlanetData) {
	// convert the system at the given coordinates
	fmt.Printf("Converting system %d, %d, %d\n", s.X, s.Y, s.Z)

	// update the star with values from the source template
	for i, planet := range src {
		s.Planets[i] = planet.Clone()
	}

	// make minor random changes to the planets
	for _, planet := range s.Planets {
		if planet.TemperatureClass == 0 {
			// no changes
		} else if planet.TemperatureClass > 12 {
			planet.TemperatureClass -= rnd(3) - 1
		} else {
			planet.TemperatureClass += rnd(3) - 1
		}
		if planet.PressureClass == 0 {
			// no changes
		} else if planet.PressureClass > 12 {
			planet.PressureClass -= rnd(3) - 1
		} else {
			planet.PressureClass += rnd(3) - 1
		}
		if len(planet.Gases) > 2 {
			j := rnd(25) + 10
			a, b := 1, 2
			if planet.Gases[b].Percentage > 50 {
				planet.Gases[a].Percentage += j
				planet.Gases[b].Percentage -= j
			} else if planet.Gases[a].Percentage > 50 {
				planet.Gases[a].Percentage -= j
				planet.Gases[b].Percentage += j
			}
		}
		if planet.Diameter > 12 {
			planet.Diameter -= rnd(3) - 1
		} else {
			planet.Diameter += rnd(3) - 1
		}
		if planet.Gravity > 100 {
			planet.Gravity -= rnd(10)
		} else {
			planet.Gravity += rnd(10)
		}
		if planet.MiningDifficulty > 100 {
			planet.MiningDifficulty -= rnd(10)
		} else {
			planet.MiningDifficulty += rnd(10)
		}
	}
}

func (s *StarData) DistanceSquaredTo(to *StarData) int {
	deltaX, deltaY, deltaZ := s.X-to.X, s.Y-to.Y, s.Z-to.Z
	return (deltaX)*(deltaX) + (deltaY)*(deltaY) + (deltaZ)*(deltaZ)
}

// returns number, not index
func (s *StarData) HomePlanetNumber() int {
	for i, planet := range s.Planets {
		if planet.Special == IDEAL_HOME_PLANET {
			return i + 1
		}
	}
	return 0
}
