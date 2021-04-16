# Manic Miner 2021 (BBC Micro)

Having recently completely disassembled Manic Miner [https://github.com/TobyLobster/ManicMiner] it seemed irresistible to improve it. I've written up some notes below.

    NEWMINER.ssd

# Improvements

- Fast
- Flicker free player movement
- Fixes to cavern layouts, graphics, and colours
- Better air bar and colours
- Fixed the shape of the jump to match the Spectrum
- Better collision detection
- Better music
- 'GAME    OVER' added
- Master compatible

# What I did

The first job was to gather the code into one source file as the original was split over multiple files that called each other back and forth. I also added a short loader program to load and copy memory around before execution. My mantra throughout was to keep the game running at all times because it's way too easy to introduce bugs in this process, so having a running game means I find and fix them faster.

I could see there wasn't much memory available so I started making various optimisations for space. I removed the latent copy protection code, removed the redundant delay loop, shortened the code that draws the title screen (adding a better Penrose triangle while I was there).

Now I had a bit of memory to play with. I started to optimise for speed, starting with the plot routine. There's one plot routine that draws the player and the guardians that move vertically, so I started by using a table lookup for calculating the screen address for each character row. The plot routine also has a 'mode' indicating how the pixels should combine with the screen, and I improved the performance of handling this too.

In a few places (including plot) there were calculations where a multiply by a constant factor was implemented in a loop of repeated addition. I fixed these to use more efficient code. In several places the use of `cmp #0` could be removed:

    lda someValue
    cmp #0              ; <-- redundant. The lda instruction sets the Z flag accordingly already
    beq somePlace

Next up, I optimised the routine that copies from the screen memory around the player to a cache and back again.
Another small example of optimisation, here's a routine that reverses the order of the bits in a byte (from the original code):

    reverseBits
        sta byteToReverse

        ldy #0
        ldx #0
    reverseBitsLoop
        lda byteToReverse
        clc
        rol
        sta byteToReverse
        bcs bitSet
        txa
        clc

    continueReverseBits
        ror
        tax
        iny
        cpy #8
        bne reverseBitsLoop
        rts

    bitSet
        txa
        sec
        jmp continueReverseBits

and here's the optimised version [from https://sites.google.com/site/h2obsession/programming/6502]:

    reverseBits
        sta byteToReverse
        lda #1
    -
        lsr byteToReverse
        rol
        bcc -
        rts

Which is clearly much smaller and faster.

By the way, nothing written here is meant to denigrate the original work from 1984, which remains a remarkable achievement in the world before such things as the internet, source control, fast reliable storage media, extensive information and a community about the inner workings of the BBC Micro, etc. It all just highlights the benefits we have today.

Now it was time to replace the routine that had the single largest impact on speed: OSBYTE 135. This reads a character at the current text cursor position and decodes the pixels to it's ASCII value. I replaced the calls to the OS routine with my own copy, but this new version was only interested in checking character 32 (space) and the user defined characters (128-160). By restricting the range of characters to test and also hardcoding to reading pixels in MODE 1, a significant speed up was made.

This added more code to the program, so I found space by moving common code into subroutines, e.g. moving the text cursor, setting the text colour, checking if a key is pressed etc.

Then back to optimisation for speed, this time animating the conveyors. The same sprite is copied multiple times along it's length, so this plotting is now basically a simple memory copy routine. Previously this was implemented with repeated OSWRCH calls.

With these speed improvements, I wanted to make the player character flicker free. When updating the player on screen, the code first reinstates the screen background (6 character cells around the player) from cache, animates the conveyor by one step, records (copies) the new screen background around the new player position into cache, then draws the player at the new position. All this needs to happen in the time between the electron gun passing the current player position on screen and before it returns on the next frame to the new player position.

I set up a timer that interrupts as the electron beam reaches every third character row on the screen, incrementing a counter as it goes. In the main program when about to update the player as described above I look at this counter. If I estimate the update code will not complete in time (before the electron gun reaches the position of the player on screen) then I wait until the counter (and therefore the electron gun) has passed the player and it is safe to continue to update the player. This is fiddly to get right, but totally worth it.

The original code used an Event (timed with the vertical sync) to update the music and some sound effects. I now use my new timer based interrupt to do these same jobs instead.

I make more space by putting the 'foot' sprite into the (formerly unused) printer buffer. I also shorten code that copies sprites for the current level into the user defined characters at $0c00-$0cff.

Now some gameplay: Miner Willy falls much slower on the BBC version (2 pixels per frame) compared to the Spectrum original (4 pixels per frame), so I reinstated the correct falling speed. The air runs out much faster on the original BBC version, so I made the air depletion rate much closer to the Spectrum version.

Eugene on Eugene's Lair falls toward the exit when the last key is taken. This was too fast on the original BBC version so I've made him fall slower to the exit now to match the Spectrum version.

I draw the air bar at the bottom of the screen rather than down the side, which makes the display look more balanced and matches the Spectrum version. This does use more screen memory since one more character row need to be visible, but I've regained enough memory to do that. I also use my interrupt routine to change the palette for the bottom three character rows of the screen allowing different colours to be used in the air bar and scores text compared to the main play area. Draining the air at the end of a level more closely matches the look, timing, and score of the Spectrum version. The sound and visual effect has also been tweaked.

Time to improve another feature: Music. I've changed the short staccato beeps reminiscent of the Spectrum original to be less harsh on the ears and expanded it to play a much more complete tune. Some notes are accented (played slightly louder). Although this is quite subtle I think it helps the overall experience. I repeat the main section of tune at a faster pace, then play the coda. I'm still using one channel only for the music. This means the 'jump' and 'key taken' sound effects can live on two separate channels and not cut each other short.

The code to animate the (up to) five keys used to use OSBYTE 135 to see if they were still present, but I now store this information in a separate array to avoid the need. I only animate one key at a time to further improve performance.

I have written a second plot routine (cellPlot) to draw at character positions on screen and this is used for the horizontal guardians, flashing exits, colour cycling the keys etc. This replaces all of the OSWRCH calls during regular play for a further performance boost.

There are several differences between the graphics definitions for the Spectrum and BBC, and I wanted to restore them to the proper Spectrum versions. I fixed the exits (e.g. Skull and crossbones exit in The Ore Refinery; Endorian Forest exit; and Kong levels exit among others). I added the telephone box exit from the Spectrum version that was missing from the BBC.

I added more 'spikes' to correspond with the Spectrum version (e.g. level 2's spike, a second type of spider sprite etc). I fixed the banana, apple, money collectibles, platform definitions, walls (including adding a missing graphic), switches, the pedestal etc. Some of these are subtle, but I thought I'd get it right. Thankfully the player and guardian sprites were already pixel perfect.

I fixed the room title so that it doesn't overlap with the bottom cell of the gameplay area.

By this time The Central Cavern (being a relatively simple level with only one guardian) was running too fast(!) so I added some code to regulate time. I've set the game speed to be a bit faster than the Spectrum version as that always felt slightly slow to me. By the end of my optimisations all levels run at pretty much regulation speed.

The walls are 'filled in' as per the Spectrum now (e.g. in Central Cavern they are yellow and red rather then yellow and black)

More memory was recovered by moving lots of variables zero page. Each time any variable is accessed they use up one less byte now. Further variables moved to page one. I saved more memory by having a 'printFollowingMessage':

    ... some code ...
    jsr printFollowingMessage
    !byte 5          ; message length
    !text "Hello"
    ... more code ...

Inlining the message with the code means the address of the text doesn't need to be stored anywhere.

I've moved the initialisation code to $0400. This code only ever executes once on startup, so once initialisation is complete that memory can be reused (in our case to cache the vertical sprite definitions for the current level).

The look of the title screen has been tidied up, and now also displays the high score.

Somewhere around this point I did a pass looking through the levels comparing them to the original Spectrum version to make sure they had the right colours and graphics. This also involved adding the ability to set individual colours for strips and single items in a level. We are still constrained by the limit of only four colours in the play area, but we do pretty nicely within that. This would be the first of several passes where I would find more tweaks, such as The Sixteenth Carven's layout was wrong with the right hand side of the screen moved one cell too far right. This is now fixed. The player start positions are corrected (e.g. level 16 again), and the horizontal guardians now have the correct start positions and extents to match the Spectrum.

I had a go at adding true pixel based collision code, but it was too complex, large and slow. So I compromised. I now do collision with octagonal boxes - boxes with the corners chopped off to a greater or smaller extent for each guardian. This proved to be pretty good in practice, and certainly an improvement over the original BBC version (box collision).

For Master compatibility the issue was that the BBC Micro stores the user defined characters at $0c00 to $0cff, but the Master doesn't. It squirrels them away in another bank of memory. So when the $0c00 to $0cff memory has changed (at the start of a level) I call the OS to define each character 'officially'. This would be too slow to do during gameplay but fine at level start.

I reimplemented the Meteors and Energy Fields of levels 19 and 20 to use cellPlot not OSWRCH and made the code shorter too.

The Game Over screen now says 'GAME OVER' and colour cycles it similarly to the Spectrum.

At the last minute, I realised the shape of the jump was subtly different from the Spectrum, so I fixed this.

Finally I tested each level in turn making sure each was completable.

TobyLobster
