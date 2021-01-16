## notes about the design of the melrose tool and its language

### error handling
the explicit creation of a musical object, e.g note,sequence,joinmap,etc will panic if the arguments are invalid

if the implicit creation of sequences fails while playing, any error will result in a warning, no panic.

### map

use key:value pairs, separated by comma

# trigger start and stop of a playable by pressing a note
trigger('c1',play(sequence('c e f')))
fun = play(sequence('c e f'))
trigger(device(2,'c1'),fun)

fun = loop(scale(2,'c4'))
trigger('d1',fun)
