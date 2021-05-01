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

type PlanetData struct {
	TemperatureClass int        /* Temperature class, 1-30. */
	PressureClass    int        /* Pressure class, 0-29. */
	Special          int        /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Gases            []*GasData /* Gas in atmosphere. Nil if none. */
	Diameter         int        /* Diameter in thousands of kilometers. */
	Density          int
	Gravity          int /* Surface gravity. Multiple of Earth gravity times 100. */
	MiningDifficulty int /* Mining difficulty times 100. */
	EconEfficiency   int /* Economic efficiency. Always 100 for a home planet. */
	MDIncrease       int /* Increase in mining difficulty. */
	Message          int /* Message associated with this planet, if any. */
}

type GasData struct {
	Type       GasType
	Percentage int
}

func GeneratePlanet(num_planets int) ([]*PlanetData, error) {
	var planets []*PlanetData

	// Values for the planets of Earth's solar system will be used as starting values.
	// Diameters are in thousands of kilometers.
	// The zeroth element of each array is a placeholder and is not used.
	// The fifth element corresponds to the asteroid belt, and is pure fantasy on my part.
	// I omitted Pluto because it is probably a captured planet, rather than an original member of our solar system.
	earth := []struct{ diameter, temperatureClass int }{
		{0, 0},   // unused
		{5, 29},  // Mercury
		{12, 27}, // Venus
		{13, 11}, // Earth
		{7, 9},   // Mars
		{20, 8},  // Asteroid Belt
		{143, 6}, // Jupiter
		{121, 5}, // Saturn
		{51, 5},  // Uranus
		{49, 3},  // Neptune
	}

	/* Main loop. Generate one planet at a time. */
	for planet_number := 1; planet_number <= num_planets; planet_number++ {
		planet := &PlanetData{}
		planets = append(planets, planet)

		/* Start with diameters, temperature classes and pressure classes based on the planets in Earth's solar system. */
		var startOffset int
		if num_planets <= 3 {
			startOffset = 2*planet_number + 1
		} else {
			startOffset = (9 * planet_number) / num_planets
		}
		planet.Diameter = planet.GenerateDiameter(earth[startOffset].diameter)
		planet.TemperatureClass = earth[startOffset].temperatureClass

		/* If diameter is greater than 40,000 km, assume the planet is a gas giant. */
		gas_giant := (planet.Diameter > 40)

		planet.Density = planet.GenerateDensity(gas_giant)

		// Gravitational acceleration is proportional to the mass divided by the radius-squared.
		// The radius is proportional to the diameter, and the mass is proportional to the density times the radius-cubed.
		// The net result is that "g" is proportional to the density times the diameter.
		// Our value for "g" will be a multiple of Earth gravity, and will be further multiplied by 100 to allow us to use integer arithmetic.
		// The factor 72 ensures that "g" will be 100 for Earth (density=550 and diameter=13).
		planet.Gravity = (planet.Density * planet.Diameter) / 72

		planet.TemperatureClass = planet.GenerateTemperatureClass(num_planets, planet_number, gas_giant, earth[startOffset].temperatureClass)
		/* Make sure that planets farther from the sun are not warmer than planets closer to the sun. */
		if planet_number > 1 && planets[planet_number-1].TemperatureClass < planet.TemperatureClass {
			planet.TemperatureClass = planets[planet_number-1].TemperatureClass - (rnd(3) - 1)
			if planet.TemperatureClass < 1 {
				planet.TemperatureClass = 1
			}
		}

		/* Pressure class depends primarily on gravity. Calculate an approximate value and randomize it. */
		planet.PressureClass = planet.GeneratePressureClass(planet.Gravity, planet.TemperatureClass, gas_giant)

		/* Generate gases, if any, in the atmosphere. */
		for _, gas := range planet.GenerateGases(planet.PressureClass, planet.TemperatureClass) {
			planet.Gases = append(planet.Gases, gas)
		}

		// Get mining difficulty.
		planet.MiningDifficulty = planet.GenerateMiningDifficulty(planet.Diameter)
	}

	return planets, nil
}

// GenerateDensity
// Density will depend on whether or not the planet is a gas giant.
// Again ignoring Pluto, densities range from 0.7 to 1.6 times the density of water for the gas giants, and from 3.9 to 5.5 for the others.
// We will expand this range slightly and use 100 times the actual density so that we can use integer arithmetic.
func (p *PlanetData) GenerateDensity(gasGiant bool) int {
	var base, sigma int
	if gasGiant {
		/* Final values from 60 thru 170. */
		base, sigma = 58, 56
	} else {
		/* Final values from 370 thru 570. */
		base, sigma = 368, 101
	}
	return base + rnd(sigma) + rnd(sigma)
}

// GenerateDiameter
func (p *PlanetData) GenerateDiameter(baseDiameter int) int {
	diameter, die_size := baseDiameter, baseDiameter/4
	if die_size < 2 {
		die_size = 2
	}
	for i := 1; i <= 4; i++ {
		if rnd(100) > 50 {
			diameter += rnd(die_size)
		} else {
			diameter -= rnd(die_size)
		}
	}

	// Minimum allowable diameter is 3,000 km.
	// Note that the maximum diameter we can generate is 283,000 km.
	for diameter < 3 {
		diameter += rnd(4)
	}

	return diameter
}

// GenerateGases
/* Generate gases, if any, in the atmosphere. */
func (p *PlanetData) GenerateGases(pressureClass, temperatureClass int) []*GasData {
	if pressureClass == 0 {
		// no atmosphere, no gases
		return nil
	}
	var gases []*GasData

	// Convert planet's temperature class to a value between 1 and 9.
	// We will use it as the start index into the list of 13 potential gases.
	var firstGas GasType
	switch 100 * temperatureClass / 225 {
	case 0:
		firstGas = H2 /* Hydrogen */
	case 1:
		firstGas = H2 /* Hydrogen */
	case 2:
		firstGas = CH4 /* Methane */
	case 3:
		firstGas = HE /* Helium */
	case 4:
		firstGas = NH3 /* Ammonia */
	case 5:
		firstGas = N2 /* Nitrogen */
	case 6:
		firstGas = CO2 /* Carbon Dioxide */
	case 7:
		firstGas = O2 /* Oxygen */
	case 8:
		firstGas = HCL /* Hydrogen Chloride */
	default:
		firstGas = CL2 /* Chlorine */
	}

	/* The following algorithm is something I tweaked until it worked well. */
	num_gases_wanted := (rnd(4) + rnd(4)) / 2
	for len(gases) == 0 {
		for i := firstGas; i <= firstGas+4 && len(gases) < num_gases_wanted; i++ {
			if i == HE && temperatureClass > 5 {
				// too hot for Helium
				continue
			}
			// skip to the next gas about one-third of the time
			// (unless we're on Helium, then it's two-thirds of the time)
			switch rnd(3) {
			case 2:
				if i == HE {
					/* Don't want too many Helium planets. */
					continue
				}
			case 3:
				continue
			}

			gas := &GasData{Type: i}
			switch i {
			case HE:
				// Helium is self-limiting
				gas.Percentage = rnd(20)
			case O2:
				// Oxygen is self-limiting
				gas.Percentage = rnd(50)
			default:
				gas.Percentage = rnd(100)
			}
			gases = append(gases, gas)
		}
	}

	// determine total quantity of gases in the atmosphere
	var total_quantity int
	for _, gas := range gases {
		total_quantity += gas.Percentage
	}
	// convert gas quantities to percentages
	var total_percent int
	for _, gas := range gases {
		gas.Percentage = 100 * gas.Percentage / total_quantity
		total_percent += gas.Percentage
	}

	// give leftover to first gas
	gases[0].Percentage += 100 - total_percent

	return gases
}

// GenerateMiningDifficulty
// Basically, mining difficulty is proportional to planetary diameter with randomization and an occasional big surprise.
// Actual values will range between 0.80 and 10.00.
// Again, the actual value will be multiplied by 100 to allow use of integer arithmetic.
func (p *PlanetData) GenerateMiningDifficulty(diameter int) int {
	mining_dif := 0
	for mining_dif < 40 || mining_dif > 500 {
		mining_dif = (rnd(3)+rnd(3)+rnd(3)-rnd(4))*rnd(diameter) + rnd(30) + rnd(30)
	}

	mining_dif = (mining_dif * 11) / 5 /* Fudge factor. */

	return mining_dif
}

// GeneratePressureClass
func (p *PlanetData) GeneratePressureClass(gravity, temperatureClass int, gasGiant bool) int {
	if gravity < 10 {
		// gravity is too low to retain an atmosphere
		return 0
	} else if temperatureClass < 2 || temperatureClass > 27 {
		// Planets outside this temperature range have no atmosphere
		return 0
	}

	pressureClass := gravity / 10
	die_size := pressureClass / 4
	if die_size < 2 {
		die_size = 2
	}
	for i, nRolls := 1, rnd(3)+rnd(3)+rnd(3); i <= nRolls; i++ {
		if rnd(100) > 50 {
			pressureClass += rnd(die_size)
		} else {
			pressureClass -= rnd(die_size)
		}
	}

	if gasGiant {
		for pressureClass < 11 {
			pressureClass += rnd(3)
		}
		for pressureClass > 29 {
			pressureClass -= rnd(3)
		}
	} else {
		for pressureClass < 0 {
			pressureClass += rnd(3)
		}
		for pressureClass > 12 {
			pressureClass -= rnd(3)
		}
	}

	return pressureClass
}

// GenerateTemperatureClass
func (p *PlanetData) GenerateTemperatureClass(numPlanets, orbit int, gasGiant bool, baseTemperatureClass int) int {
	/* Randomize the temperature class obtained earlier. */
	temperatureClass, die_size := baseTemperatureClass, baseTemperatureClass/4
	if die_size < 2 {
		die_size = 2
	}
	n_rolls := rnd(3) + rnd(3) + rnd(3)
	for i := 1; i <= n_rolls; i++ {
		if rnd(100) > 50 {
			temperatureClass += rnd(die_size)
		} else {
			temperatureClass -= rnd(die_size)
		}
	}

	if gasGiant {
		for temperatureClass < 3 {
			temperatureClass += rnd(2)
		}
		for temperatureClass > 7 {
			temperatureClass -= rnd(2)
		}
	} else {
		for temperatureClass < 1 {
			temperatureClass += rnd(3)
		}
		for temperatureClass > 30 {
			temperatureClass -= rnd(3)
		}
	}

	// Sometimes, an inner planet in star systems with less than four planets are too cold.
	// Warm them up a little.
	if numPlanets < 4 && orbit < 3 {
		for temperatureClass < 12 {
			temperatureClass += rnd(4)
		}
	}

	return temperatureClass
}
