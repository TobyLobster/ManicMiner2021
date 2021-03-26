# Manic Miner (BBC)

This is a disassembly/reassembly of Manic Miner for the BBC Micro.
I thought I would write up some thoughts on it while I'm here.

    constants.a
    MINER1.a
    MINER2.a
    MINER3.a
    MINER4.a

This code assembles (using the ACME assembler) to create a byte for byte identical copy of the original game.

I should mention that this is the version from http://bbcmicro.co.uk/game.php?id=188 so it differs slightly from the original in that it has instructions before the game loads, and the copy protection has been circumvented.

# Copy Protection

The original asks for a four digit code from a random grid location found on a physical sheet of paper that was supplied with the game. The paper has groups of four colours laid out in a grid, each colour corresponding to a number 1-4 that must be entered correctly to continue. This system has previously been circumvented, but the original code to handle this is still present.

# Loading
After the instructions, there are three more title screens!

The first shows the large dancing letters of MANIC MINER and a scrolling message. The cassette version of the game animates this screen while loading, which I think is fairly unusual.

The second was used for the entry code, but here is displayed for a brief time before the third title screen with a 'Penrose triangle' (the triangular 3d optical illusion) drawn via OSWRCH PLOT commands. At this point the game is fully loaded, and RETURN starts a new game.

# Obfuscation
*FX 200,3 is sprinkled a few times while loading / initialising to ensure memory is cleared on BREAK.

Code and data are split between files MINER1 through MINER4, and moved around in memory, sometimes EOR'd with $55 (code) or $AA (data).

A couple of short routines are stashed in zero page, perhaps to make them harder to find, to hinder cheating. The first is the code that triggers a switch, and the second is the code that resets your score at the end of the game.

A single byte of the level data is changed (the position of a key on level 20) during initialisation. Perhaps to make it harder to make a completely perfect copy of the game, or perhaps it's just a last minute 'fix'. It doesn't make a significant difference to the difficulty of the last level.

# Implementation issues

## Slow
The games runs slowly compared to the Spectrum original. This is partly down to the method of plotting. Horizontal guardians are drawn using OSWRCH to write user defined characters to the screen. The conveyors and keys are animated using the same method.

The second major reason for being slow is that there is no collision map. When the game needs to know what is on the screen (e.g. the squares surrounding the player) it reads the character from the screen using OSBYTE 135. The OS has to read the 8x8 pixels from the screen and then compare them against each character 32-255 in the current character set in turn until if finds a match. None of this is efficient!

A further example: each key is animated (one at a time round robin style) by moving the text cursor to the X,Y position of a key (3 calls to OSWRCH), reading the character from the screen (OSBYTE 135) and if it's still a key (not taken) then plot the key in a new colour (three more calls to OSWRCH).

The player and vertical guardians are drawing using a custom plot routine (which is also not greatly efficient).

In several places in the code (including the plot routine), multiplication is required but this is implemented with a loop of repeated addition. This is sub-optimal.

There is also an explicit short delay loop in the main loop, but this is a negligible factor in the speed. This was probably used more in early development when (with fewer features implemented) the speed would otherwise be too fast.

## Flicker
The player in particular is flickery when moving. While the player is not shown, the game copies bytes at the player position on screen to a cache (in a non-visible area of the screen). This cache just has the background graphics around the player and is used for collision detection. It is later copied back to the visible screen to erase the player.

## Collision detection
Sadly there is only box collision here, not pixel perfect collision. The boxes around the guardians are adjusted to give a little leeway.

## Music
Music is basic (similar to the Spectrum) and a note can be missed when overridden by playing a sound effect.

## Compatibility
Doesn't work on the either the Master or Electron.

## Other variations from the Spectrum version
* The BBC version does not occupy the full screen.

* Only four colours are available in this MODE 1 screen, instead of the Spectrum's 15.

* The BBC cannot control the 'border colour' as the Spectrum does.

* The title screen is rudimentary on the BBC. The Spectrum version has a visual scene, a piano keyboard and The Blue Danube playing.

* The layout of the main game on screen is different, with the Air bar at the side rather than below the play area, and the Spectrum version shows the lives visually as player sprites walking, but it's just a number on the BBC.

* The individual graphics are close but sometimes not perfectly identical to the Spectrum.

* The final two levels are different. The Spectrum's "Solar Power Generator" has been replaced with a new design "The Meteor Storm", and the final level is a different design too (but still called "The Final Barrier"). They both include a new feature that I am calling 'energy fields', barriers that turn on / off at regular intervals and are deadly when on.

* A few vertical guardians move at four pixels per frame on the Spectrum, but not on the BBC version: One of the Skylabs, a guardian in The Warehouse, and a guardian in Amoebatrons' Revenge.

# The level teleport cheat
Type 'A SECRET' on the pause screen. After resuming play, the fn keys teleport you to the different levels, and SHIFT+fn keys go to later levels. The code is a little hidden since it doesn't use one of the regular methods of reading keys. It doesn't call the OS, nor does it read the keyboard directly. Instead it reads memory location $ec, which the OS uses to hold the code of the key currently pressed. I wonder if this was deliberate obfuscation.

# Level format

All bytes relating to the level data are EOR'd with $55 in the binary files on disk.

## Strips
Levels are constructed from horizontal strips and rectangles, stored at 'levelDefinitions' (memory address $6c00, or offset $1680 into MINER4 binary file).

The separator between each level is $ff, $ff which occurs before and after every level.
This is followed by header information:

    byte offset     meaning
         0          two nybbles, high is palette colour 0, low is palette colour 1
         1          low nybble is palette colour 2
         2          high nybble is the regular floor sprite, low nybble is the crumble floor
         3          lower 3 bits are the conveyor sprite and the side wall sprite
         4          lower nybble is the key sprite
         5          the exit sprite - top two bits are the type:

           00 = 8x8 exit sprite repeated four times (lower 3 bits are sprite number)
           10 = 8x16 exit sprite mirrored about the Y axis (lower 3 bits are sprite number)
           01 = }
           11 = } 16x16 exit sprite (lower 6 bits are the sprite number)

    Then the level data itself starts, consisting of a sequence of commands:

    command         description
        $ff         Increments the levelFeatureIndex. Once it reaches 5, the level is done.
        $fe         The next four bytes determine a rectangle to draw:
                        <x_min> <y_min> <width> <y_max>
        $fd         Sets the levelFeatureIndex to the next byte value.
        else        Next three bytes determine a horizontal strip to draw:
                        <x_min> <y_min> <width>

## Single Items
Items that typically occupy a single cell are stored at 'levelSingleItemDefinitions' (memory address $7170, or offset $1bf0 into MINER4 binary file)

The separator between each level is $ff which occurs before and after every level.

A single byte header holds the 'number of keys - 2' in the top two bits, and the sprite increment amount in the remaining lower bits.

Then there is a sequence of commands:

    command         description
        $ff         end of level
        $fe         all X coordinates have 15 added from now on until the next change of type.
        $fd         increment current type by the sprite increment amount. All X coordinates have 5 added.
        else        top nybble is X coordinate, bottom nybble is Y coordinate. Plot the item.

At plot time, the current type determines the outcome:

        $eb         switch to colour 1 and draw <something>
        $ee         spider and thread (next byte is length of thread)
        $f0         key (colour 3)
        else        sprite number (see below)

## Standard Sprites

    Sprite numbers start at $80 and repeat after 32, so $80 is the same as $a0 $c0 and $e0.
    Alternative sprite pages are swapped in and out as needed e.g. to draw horizontal guardians.

      sprite        description
        ($5f         underscore, which is identical to the most crumbled floor sprite)
        $80         platform
        $81         alternative platform
        $82         crumble floor
        $83         crumble floor
        $84         crumble floor
        $85         crumble floor
        $86         crumble floor
        $87         crumble floor
        $88         crumble floor
        $89         crumble floor
        $8a         (empty)
        $8b         conveyor
        $8c         conveyor
        $8d         conveyor
        $8e         conveyor
        $8f         wall
        $90         key
        $91         exit top left
        $92         exit top right
        $93         exit bottom left
        $94         exit bottom right
        $95         (empty)
        $96         (empty)
        $97         (empty)
        $98         unswitched switch
        $99         switched switch
        $9a         (empty)
        $9b         spike 1
        $9c         spike 2
        $9e         thread
        $9f         (empty)

## Vertical Guardians
Level 11 upwards can change the sprites used for vertical guardians.
This is stored in the 'verticalGuardiansSpritesArray' array (memory address $1da2, or offset $14a2 into MINER2 binary file).
One byte each from level 11 onwards.

Levels 9 and 11 upwards can have up to four vertical guardians.
This is stored at 'verticalGuardians' (memory address $2478, or offset $1B78 into MINER2 binary file) as eight bytes per level (two bytes per guardian. Use $ff, $ff for unused guardian slots).
The array holds data from level 8 upwards, but levels 8 and 10 are not used.

Each guardian is stored in two bytes:
    <x coordinate>    x position in cells, top bit specifies initial direction, (clear = up, set = down)
    <y coordinates>   top nybble is the initial Y and also the second Y extent, bottom nybble is the first Y extent

## Horizontal Guardians
Stored at 'guardianLevelData' (memory address $2300, or offset $1A00 in MINER2 binary file). Level separator before and after each level is '$ff', followed by three bytes per guardian specifying initial position and direction.
       X1, Y + top bit, X2
       'top bit' indicates initial direction (set = moving left), then bouncing between X1 and X2

The array 'guardianSetForEachLevel' (memory address $20e6, or offset $17e6 in MINER2 binary file) holds the index into the horizontal guardians sprites (00-0f) for the level.

## Conveyor directions
The direction of the conveyor on each level is stored in the top bit of the 20 bytes at 'fixedText' (memory address $25a8, or offset $1ca8 into MINER2 binary file). The bottom 7 bits must remain unchanged.

## Player start positions
The X pixel coordinate of the start position is stored in the array 'playerStartPositions' (memory address $6b88, or offset $1608 into MINER4 binary file), one byte per level.

## Limitations (Hardcoded features)
Many of the level specific features are not data driven, and are hard-coded. Eugene's and Kong's plummets into the exit for example. Level 14's Skylab plummeting to earth. Level 19's Meteors. The energy fields on levels 19 and 20. Level 10 having no conveyor. Level 16 not having any vertical guardians. Switches are at fixed positions on level 8 and 12 only (also the two nearby spikes are hardcoded). A byte in the level 20 data is poked at initialisation time 'pokeSingleItem'. Level 6's exit must remain in the same position, since the code to exit the level is hardcoded to this position. Palette colours are not flexible.
