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

# connect knob/slider to variale value

k = knob('c3')
l = loop(pitch(k,sequence('c e g')))

or editor text syncs with knob. A "4" turns into "0" if the knob is turned down.
somehow need to "save" and "load" knob settings.
other names:
- range('c3')
- slider('c2')
- or allow aliases
- input('b2')
- something that looks like "onkey" ?
- control('f1')
- how to specify device?  
    - control(device(1), 'f1')
    - control(device(1,'e3'))
    - control(1,'e3')
- how to specify channel?
    - k = control(device(1,channel(2, 'c2')))
- currently listen cannot specify channel
- k = valueof(1,'c3')
- k = control(1,'c3',0)  // device 1, cc name c3, start with 0
- k = valueof(1,2,'c3')
- k = valueof(device(1,channel(2,range('c3'))))
- k = device(1).channel(2).range('c3')

In the listener, keep a map of knobListeners with the known value.