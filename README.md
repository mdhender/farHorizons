# fh
This repository implements a new server for the 7th edition of Rick Morneau's *FAR HORIZONS*.
The code is a port of [Ramblurr's repository](https://github.com/Ramblurr/Far-Horizons).
The server is under active development and should be treated cautiously.

# FAR HORIZONS
FAR HORIZONS is a strategic role-playing game of galactic exploration, trade, diplomacy, and conquest.
This is the seventh edition of a game made by Rick Morneau.
It is NOT a beta game and has been thoroughly play tested for years.
It is closed-ended, but we accept replacement players to play species that drop out.

A new game is beginning as soon as we have enough players.
If you are interested, just head to the [website](https://example.com/index.html) and fill out the online entry form.

Rules and turn schedule are on the site.
Turns are weekly and computer interpreted.
I monitor everything closely, and will not be playing, of course.

DESCRIPTION:
At the start of a game, each player controls an intelligent species and the home planet on which it lives.
As the game progresses, you can explore nearby regions of the galaxy and establish colonies.
As you range farther and farther from home, you will encounter other intelligent species.
These encounters can be hostile, neutral, or friendly, depending on the participants.
Interstellar war is a distinct possibility.

FAR HORIZONS, unlike some similar games, has been designed to make role-playing as easy and practical as possible.
In addition to being a rich and realistic simulation, there are no true victory conditions - the game is played solely for enjoyment.
However, at the end of the last turn, final statistics for all species will be sent to all the players so that they can compare their relative strengths and weaknesses.
Thus, rather than requiring a massive bloodletting as in some other similar games, it's possible for a peace-loving species to effectively "win".

Still, those who enjoy a more aggressive game, or those who wish to
role-play an "evil" or warlike species will not be disappointed.
FAR HORIZONS does not discriminate against anyone - it simply tries to be as realistic as possible.

# GameMaster
## Sequence
The following steps apply only to normal turns.

For the setup turn, first create the galaxy using the NewGalaxy and MakeHomes programs.
Then generate the star list (ListGalaxy -p) and map files (MapGalaxy).
Next, for each player run the HomeSystem and AddSpecies programs.
After this has been done for all species, run Finish and Report.
Finally, continue with step 3 below.
Before you do, verify the galaxy does not contain duplicate star
systems or planets:
```shell
ListGalaxy |sort|uniq -cd|less
```
If it does - regenerate.

1. As orders come in, copy them to the appropriate `spNN.ord` files in the `fh/game` directory.
   (You may want to use script `fhorders` to do this automatically.)

2. After all orders have been received, run NoOrders, Combat, PreDeparture, Jump, Production, PostArrival, Locations, Strike, Finish, and Report (in that order).
   The scripts `fhtest` and `fhsave` were written to automate this process.

3. Run the `fhreports` script.
   It will mail the reports to the players.

4. Run the `fhclean` script.
   It will delete all temporary files that were used during the turn, and copy all data and report files to backup directories.

## Commands

```sh
# run NewGalaxy
# Galaxy size is calculated based on default star density and
# fixed number of stars per race. Add the --less-crowded flag
# to create a galaxy with more stars.
$ fh create galaxy --number-of-species 8
# run MakeHomes
$ fh create homes
# run ListGalaxy -p
# check for duplicate systems
$ fh list galaxy -p
# run HomeSystemAuto
$ fh set home-system --auto --species 'Borgia' --system 'Foo'
# run AddSpecies
$ fh create species --name 'Borgia'
# run Finish
$ fh run finish
# run Report
$ fh run report

$ fh run no-orders                       ## NoOrders
$ fh run combat                          ## Combat
$ fh run pre-departure                   ## PreDeparture
$ fh run jump                            ## Jump
$ fh run production                      ## Production
$ fh run post-arrival                    ## PostArrival
$ fh run locations                       ## Locations
$ fh run strike                          ## Strike
$ fh run finish                          ## Finish
$ fh run report                          ## Report
$ fh map galaxy                          ## MapGalaxy
$ fh show turn                           ## TurnNumber
```

# Acknowledgments

The original source to [Far Horizons](https://github.com/Ramblurr/Far-Horizons) is licensed under the GPL v2.
Per comments at that repository, the original ANSI C source code and rules are Copyright (c) 1999 by Richard A. Morneau.
