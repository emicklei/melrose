var midi = (function() {
    "use strict";

    var module = {};

    var midiAccess, port = null,
        channel = -1;


    module.init = function() {
        if (navigator.requestMIDIAccess) {
            navigator.requestMIDIAccess().then(onMIDIInit, onMIDIReject);
        } else {
            showError('MIDI support is not present in your browser. You can still use ' +
                'your computer\'s keyboard.');
        }
    };

    var onMIDIInit = function(mAccess) {
        midiAccess = mAccess;

        // Get last used port ID from local storage.
        var preferredPort = storage.get('midiPort');

        // Add all MIDI ports to the dropdown box.
        midiAccess.inputs.forEach(function(port) {
            addMIDIPort(port);
            if (port.id == preferredPort) {
                $('#midi-port').val(port.id);
            }
        });

        midiAccess.addEventListener('statechange', MIDIConnectionEventListener);
        $('#midi-port').change(onMIDIPortChange);
        $('#midi-channel').change(onMIDIChannelChange);
        onMIDIPortChange();
    };

    var onMIDIReject = function(err) {
        console.log('Failed to obtain access to MIDI.');
    };

    var onMIDIPortChange = function() {
        var id = $('#midi-port').val();
        var currentId = (port != null ? port.id : null);

        if (id != currentId) {
            if (port != null) {
                port.removeEventListener('midimessage', MIDIMessageEventListener);
                console.log('panic');
            }

            port = midiAccess.inputs.get(id);

            if (port != null) {
                port.addEventListener('midimessage', MIDIMessageEventListener);
                console.log('Listening on port ' + port.name);

                storage.set('midiPort', port.id);
            }
        }
    };

    var onMIDIChannelChange = function() {
        var currentChannel = channel;
        channel = Number($('#midi-channel').val());

        if (channel != currentChannel)
            console.log('panic');
    };

    var MIDIConnectionEventListener = function(event) {
        var port = event.port;
        if (port.type != 'input') return;

        var portOption = $('#midi-port option').filter(function() {
            return $(this).attr('value') == port.id;
        });

        if (portOption.length > 0 && port.state == 'disconnected') {
            showWarning(port.name + ' was disconnected.');

            portOption.remove();
            onMIDIPortChange();
        } else if (portOption.length == 0 && port.state == 'connected') {
            showSuccess(port.name + ' is connected.');

            addMIDIPort(port);
            onMIDIPortChange();
        }
    };

    var addMIDIPort = function(port) {
        $('#midi-port')
            .append($("<option></option>")
                .attr("value", port.id)
                .text(port.name));
    };

    var MIDI_NOTE_ON = 0x90,
        MIDI_NOTE_OFF = 0x80,
        MIDI_CONTROL_CHANGE = 0xB0,

        MIDI_CC_SUSTAIN = 64,
        MIDI_CC_ALL_CONTROLLERS_OFF = 121,
        MIDI_CC_ALL_NOTES_OFF = 123;

    var ALL_CHANNELS = -1,
        ALL_EXCEPT_DRUMS = -10;

    var MIDIMessageEventListener = function(event) {
        var msg = event.data;
        var msgType = msg[0] & 0xF0;
        var msgChannel = msg[0] & 0x0F;

        if ((channel >= 0 && msgChannel != channel) ||
            (channel == ALL_EXCEPT_DRUMS && msgChannel == 9))
            return;

        switch (msgType) {
            case MIDI_NOTE_ON:
                if (msg[2] != 0) {
                    tonnetz.noteOn(msgChannel, msg[1]);
                    break;
                }
                // velocity == 0:  note off
            case MIDI_NOTE_OFF:
                //tonnetz.noteOff(msgChannel, msg[1]);
                break;
            case MIDI_CONTROL_CHANGE:
                switch (msg[1]) {
                    case MIDI_CC_SUSTAIN:
                        if (msg[2] >= 64) {
                            tonnetz.sustainOn(msgChannel);
                        } else {
                            tonnetz.sustainOff(msgChannel);
                        }
                        break;
                    case MIDI_CC_ALL_CONTROLLERS_OFF:
                        //tonnetz.sustainOff(msgChannel);
                        break;
                    case MIDI_CC_ALL_NOTES_OFF:
                        //tonnetz.allNotesOff(msgChannel);
                        break;
                }
                break;
        }
    };

    return module;
})();