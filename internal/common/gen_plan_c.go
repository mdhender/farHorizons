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

package common

/* Generate planets. */

type PlanetData struct {
	TemperatureClass int    /* Temperature class, 1-30. */
	PressureClass    int    /* Pressure class, 0-29. */
	Special          int    /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Gas              [4]int /* Gas in atmosphere. Zero if none. */
	GasPercent       [4]int /* Percentage of gas in atmosphere. */
	Diameter         int    /* Diameter in thousands of kilometers. */
	Gravity          int    /* Surface gravity. Multiple of Earth gravity times 100. */
	MiningDifficulty int    /* Mining difficulty times 100. */
	EconEfficiency   int    /* Economic efficiency. Always 100 for a home planet. */
	MDIncrease       int    /* Increase in mining difficulty. */
	Message          int    /* Message associated with this planet, if any. */
}

var start_diameter = [10]int{0, 5, 12, 13, 7, 20, 143, 121, 51, 49}
var start_temp_class = [10]int{0, 29, 27, 11, 9, 8, 6, 5, 5, 3}

/* Values for the planets of Earth's solar system will be used
   as starting values. Diameters are in thousands of kilometers.
   The zeroth element of each array is a placeholder and is not
   used. The fifth element corresponds to the asteroid belt, and
   is pure fantasy on my part. I omitted Pluto because it is probably
   a captured planet, rather than an original member of our solar
   system. */

func GeneratePlanets(first_planet *planet_data, num_planets int) []*PlanetData {
	var density, temp, total_percent int
	var first_gas, gas_quantity, num_gases_found, num_gases_wanted int
	var diameter, g, mining_difficulty, pressure_class, temperature_class [10]int
	var gas, gas_percent [10][5]int

	planets := make([]*PlanetData, num_planets+1, num_planets+1)

	/* Main loop. Generate one planet at a time. */
	for planet_number := 1; planet_number <= num_planets; planet_number++ {
		/* Start with diameters, temperature classes and pressure classes based on the planets in Earth's solar system. */
		startOffset := 2*planet_number + 1
		if num_planets > 3 {
			startOffset = (9 * planet_number) / num_planets
		}
		dia := start_diameter[startOffset]
		tc := start_temp_class[startOffset]

		/* Randomize the diameter. */
		die_size := dia / 4
		if die_size < 2 {
			die_size = 2
		}
		for i := 1; i <= 4; i++ {
			if rnd(100) > 50 {
				dia = dia + rnd(die_size)
			} else {
				dia = dia - rnd(die_size)
			}
		}

		/* Minimum allowable diameter is 3,000 km. Note that the maximum diameter we can generate is 283,000 km. */
		for dia < 3 {
			dia += rnd(4)
		}

		diameter[planet_number] = dia

		/* If diameter is greater than 40,000 km, assume the planet is a gas giant. */
		gas_giant := (dia > 40)

		/* Density will depend on whether or not the planet is a gas giant.
		   Again ignoring Pluto, densities range from 0.7 to 1.6 times the
		   density of water for the gas giants, and from 3.9 to 5.5 for the
		   others. We will expand this range slightly and use 100 times the
		   actual density so that we can use integer arithmetic. */
		if gas_giant {
			/* Final values from 60 thru 170. */
			density = 58 + rnd(56) + rnd(56)
		} else {
			/* Final values from 370 thru 570. */
			density = 368 + rnd(101) + rnd(101)
		}

		/* Gravitational acceleration is proportional to the mass divided
		   by the radius-squared. The radius is proportional to the
		   diameter, and the mass is proportional to the density times the
		   radius-cubed. The net result is that "g" is proportional to
		   the density times the diameter. Our value for "g" will be
		   a multiple of Earth gravity, and will be further multiplied
		   by 100 to allow us to use integer arithmetic. */
		grav := (density * diameter[planet_number]) / 72
		/* The factor 72 ensures that "g" will be 100 for Earth (density=550 and diameter=13). */
		g[planet_number] = grav

		/* Randomize the temperature class obtained earlier. */
		die_size = tc / 4
		if die_size < 2 {
			die_size = 2
		}
		n_rolls := rnd(3) + rnd(3) + rnd(3)
		for i := 1; i <= n_rolls; i++ {
			if rnd(100) > 50 {
				tc = tc + rnd(die_size)
			} else {
				tc = tc - rnd(die_size)
			}
		}

		if gas_giant {
			for tc < 3 {
				tc += rnd(2)
			}
			for tc > 7 {
				tc -= rnd(2)
			}
		} else {
			for tc < 1 {
				tc += rnd(3)
			}
			for tc > 30 {
				tc -= rnd(3)
			}
		}

		/* Sometimes, planets close to the sun in star systems with less than four planets are too cold. Warm them up a little. */
		if num_planets < 4 && planet_number < 3 {
			for tc < 12 {
				tc += rnd(4)
			}
		}

		/* Make sure that planets farther from the sun are not warmer than planets closer to the sun. */
		if planet_number > 1 {
			if temperature_class[planet_number-1] < tc {
				tc = temperature_class[planet_number-1]
			}
		}

		temperature_class[planet_number] = tc

		/* Pressure class depends primarily on gravity. Calculate an approximate value and randomize it. */
		pc := g[planet_number] / 10
		die_size = pc / 4
		if die_size < 2 {
			die_size = 2
		}
		n_rolls = rnd(3) + rnd(3) + rnd(3)
		for i := 1; i <= n_rolls; i++ {
			if rnd(100) > 50 {
				pc = pc + rnd(die_size)
			} else {
				pc = pc - rnd(die_size)
			}
		}

		if gas_giant {
			for pc < 11 {
				pc += rnd(3)
			}
			for pc > 29 {
				pc -= rnd(3)
			}
		} else {
			for pc < 0 {
				pc += rnd(3)
			}
			for pc > 12 {
				pc -= rnd(3)
			}
		}

		if grav < 10 {
			/* Planet's gravity is too low to retain an atmosphere. */
			pc = 0
		}
		if tc < 2 || tc > 27 {
			/* Planets outside this temperature range have no atmosphere. */
			pc = 0
		}

		pressure_class[planet_number] = pc

		/* Generate gases, if any, in the atmosphere. */
		for i := 1; i <= 4; i++ { /* Initialize. */
			gas[planet_number][i] = 0
			gas_percent[planet_number][i] = 0
		}
		if pc == 0 {
			/* No atmosphere. */
			goto done_gases
		}

		/* Convert planet's temperature class to a value between 1 and 9.
		   We will use it as the start index into the list of 13 potential
		   gases. */
		first_gas = 100 * tc / 225
		if first_gas < 1 {
			first_gas = 1
		}
		if first_gas > 9 {
			first_gas = 9
		}

		/* The following algorithm is something I tweaked until it worked well. */
		num_gases_wanted = (rnd(4) + rnd(4)) / 2
		num_gases_found = 0
		gas_quantity = 0

	get_gases:
		for i := first_gas; i <= first_gas+4; i++ {
			if num_gases_wanted == num_gases_found {
				break
			}

			if i == HE { /* Treat Helium specially. */
				if rnd(3) > 1 {
					/* Don't want too many He planets. */
					continue
				}
				if tc > 5 {
					/* Too hot for helium. */
					continue
				}
				num_gases_found++
				gas[planet_number][num_gases_found] = HE
				temp = rnd(20)
				gas_percent[planet_number][num_gases_found] = temp
				gas_quantity += temp
			} else { /* Not Helium. */
				if rnd(3) == 3 {
					continue
				}
				num_gases_found++
				gas[planet_number][num_gases_found] = i
				if i == O2 {
					temp = rnd(50) /* Oxygen is self-limiting. */
				} else {
					temp = rnd(100)
				}
				gas_percent[planet_number][num_gases_found] = temp
				gas_quantity += temp
			}
		}

		if num_gases_found == 0 {
			/* Try again. */
			goto get_gases
		}

		/* Now convert gas quantities to percentages. */
		total_percent = 0
		for i := 1; i <= num_gases_found; i++ {
			gas_percent[planet_number][i] =
				100 * gas_percent[planet_number][i] / gas_quantity
			total_percent += gas_percent[planet_number][i]
		}

		/* Give leftover to first gas. */
		gas_percent[planet_number][1] += 100 - total_percent

	done_gases:

		/* Get mining difficulty. Basically, mining difficulty is
		   proportional to planetary diameter with randomization and an
		   occasional big surprise. Actual values will range between 0.80
		   and 10.00. Again, the actual value will be multiplied by 100
		   to allow use of integer arithmetic. */
		mining_dif := 0
		for mining_dif < 40 || mining_dif > 500 {
			mining_dif = (rnd(3)+rnd(3)+rnd(3)-rnd(4))*rnd(dia) + rnd(30) + rnd(30)
		}

		mining_dif *= 11 /* Fudge factor. */
		mining_dif /= 5

		mining_difficulty[planet_number] = mining_dif
	}

	/* Copy planet data to structure. */
	for i := 1; i <= num_planets; i++ {
		/* Initialize all bytes of record to zero. */
		current_planet := &PlanetData{}

		current_planet.Diameter = diameter[i]
		current_planet.Gravity = g[i]
		current_planet.MiningDifficulty = mining_difficulty[i]
		current_planet.TemperatureClass = temperature_class[i]
		current_planet.PressureClass = pressure_class[i]
		current_planet.Special = 0

		for n := 0; n < 4; n++ {
			current_planet.Gas[n] = gas[i][n+1]
			current_planet.GasPercent[n] = gas_percent[i][n+1]
		}

		planets[i] = current_planet
	}

	return planets
}
