# Notations

This document describes the various notations used in Melrose.

## Note

A note is represented by a string with the following components:

`[fraction][.]name[accidental][octave][velocity][~]`

- **fraction**: The duration of the note. Can be `1`, `2`, `4`, `8`, `16`, `32`. Defaults to `4` (quarter note).
- **.**: A dot after the fraction indicates a dotted note.
- **name**: The name of the note, from `A` to `G` (case-insensitive).
  - `=`: Represents a rest.
  - `^`: Pedal up and down.
  - `<`: Pedal up.
  - `>`: Pedal down.
- **accidental**:
  - `#` or `♯`: Sharp.
  - `_` or `♭` or `b`: Flat.
- **octave**: The octave of the note, e.g., `C4`. Defaults to `4`.
- **velocity**: The loudness of the note, represented by a sequence of `+` (louder) and `-` (softer) characters. `o` is normal.
- **~**: A tie to the next note.

### Examples

- `C`: A C4 quarter note.
- `8.F#++`: A dotted F-sharp eight note with increased velocity.
- `16G_2`: A G-flat sixteenth note in the second octave.

## Sequence

A sequence is a series of notes separated by spaces. Notes can be grouped using parentheses `()`.

`note note (note note) note`

A fraction before a group applies to all notes in the group.

### Examples

- `C D E F`: A sequence of four quarter notes.
- `8(C D E F)`: A sequence of four eight notes.
- `C (D E) F`: A sequence with a group.

## Chord

A chord is represented by a string with the following components:

`note[/quality][/interval][/inversion]`

- **note**: The root note of the chord.
- **quality**:
  - `m`: Minor.
  - `M`: Major.
  - `dim`: Diminished.
  - `aug`: Augmented.
  - `sus2`: Suspended 2nd.
  - `sus4`: Suspended 4th.
  - `7`: Dominant 7th.
- **interval**:
  - `6`: Sixth.
  - `7`: Seventh.
- **inversion**:
  - `1`: First inversion.
  - `2`: Second inversion.
  - `3`: Third inversion.

### Examples

- `C/m`: A C-minor chord.
- `F#/M7/1`: An F-sharp major seventh chord in the first inversion.

## Chord Progression

A chord progression is a sequence of chords in a scale, represented by Roman numerals.

`[fraction][.]roman[modifier] [velocity]`

- **fraction**: The duration of the chord.
- **.**: Dotted chord.
- **roman**: A Roman numeral from `I` to `VII` (or `i` to `vii`).
- **modifier**:
  - `maj`: Major.
  - `m`: Minor.
  - `dim`: Diminished.
  - `7`: Seventh.
- **velocity**: The velocity of the chord.

### Examples

- `I IV V`: A simple progression.
- `ii m7 V7 I maj7`: A more complex progression.

## Scale

A scale is defined by a starting note and a type.

`[style] note[/m]`

- **style**: The type of scale, e.g., "major". Defaults to major.
- **note**: The starting note of the scale.
- **/m**: Indicates a minor scale.

### Examples

- `C`: C major scale.
- `E/m`: E minor scale.
- `dorian F#`: F-sharp dorian scale.

## Tabs

A bass tablature is a sequence of notes, each with the following format:

`[fraction][.]name[fret][velocity]`

- **fraction**: The duration of the note.
- **.**: A dot after the fraction indicates a dotted note.
- **name**: The string name: `E`, `A`, `D`, `G`. `=` can be used for a rest.
- **fret**: The fret number, from 0 to 24.
- **velocity**: The loudness of the note.

### Example

`E3 A2 A5 D5 A5 A2 E3`
