# TODOs

Found 58 TODO mentions across 44 files.

## api / control

- [api/impl.go](../api/impl.go#L30): handle channel
- [control/knob.go](../control/knob.go#L51): check channel
- [control/recording.go](../control/recording.go#L55): temporary test recording storage

## core

- [core/beatmaster.go](../core/beatmaster.go#L74): testing behavior; NEEDSFIX
- [core/beatmaster.go](../core/beatmaster.go#L105): move checks to `SetBIAB`
- [core/chord.go](../core/chord.go#L78): `aug` handling
- [core/chord.go](../core/chord.go#L183): handle inversion 3
- [core/chord_sequence.go](../core/chord_sequence.go#L23): scanning workaround
- [core/chord_test.go](../core/chord_test.go#L78): commented test value
- [core/chord_test.go](../core/chord_test.go#L188): unfinished test
- [core/inspect.go](../core/inspect.go#L79): get port from flag
- [core/inspect_test.go](../core/inspect_test.go#L6): replace test context initialization
- [core/interfaces.go](../core/interfaces.go#L128): require `Play` context condition
- [core/iterator.go](../core/iterator.go#L50): return value needed?
- [core/loop.go](../core/loop.go#L12): protect with mutex
- [core/note.go](../core/note.go#L193): unfinished note logic
- [core/print.go](../core/print.go#L21): check condition
- [core/scale.go](../core/scale.go#L78): unfinished scale logic
- [core/sequence_builder.go](../core/sequence_builder.go#L57): duration assumption
- [core/sequence_builder.go](../core/sequence_builder.go#L81): move inside `sequenceBuilder`
- [core/sequence_ops.go](../core/sequence_ops.go#L62): create operation
- [core/timeline.go](../core/timeline.go#L226): warn
- [core/timeline.go](../core/timeline.go#L236): warn
- [core/track.go](../core/track.go#L61): unfinished track logic
- [core/valueholder.go](../core/valueholder.go#L162): used?

## docs / dsl

- [docs/Notes.md](Notes.md#L16): TODO list heading
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L26): allow fractions
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L184): handle loop
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L860): channel
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L874): channel
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L928): default input-device channel
- [dsl/eval_funcs.go](../dsl/eval_funcs.go#L1314): check type
- [dsl/eval_funcs_utils.go](../dsl/eval_funcs_utils.go#L73): use `core.ValueOf`
- [dsl/eval_funcs_utils.go](../dsl/eval_funcs_utils.go#L83): use `core.ValueOf`
- [dsl/evaluator.go](../dsl/evaluator.go#L48): tab handling
- [dsl/language_test.go](../dsl/language_test.go#L255): skipped test

## midi

- [midi/file/write.go](../midi/file/write.go#L42): make `4` configurable
- [midi/file/write.go](../midi/file/write.go#L74): create signature function
- [midi/midi_event.go](../midi/midi_event.go#L30): verify note-on check
- [midi/output_device.go](../midi/output_device.go#L174): longest TODO in core?
- [midi/play.go](../midi/play.go#L12): check `DeviceSelector`
- [midi/registry_cli.go](../midi/registry_cli.go#L111): unfinished logic
- [midi/registry_device.go](../midi/registry_device.go#L188): used?
- [midi/transport/m_listener.go](../midi/transport/m_listener.go#L121): BPM handling
- [midi/transport/rt_test.go](../midi/transport/rt_test.go#L6): create mock MIDI input

## op / ui

- [op/bare.go](../op/bare.go#L14): remove pedals and bindings
- [op/dynamic_map.go](../op/dynamic_map.go#L73): silently ignored error
- [op/fraction_map.go](../op/fraction_map.go#L112): move validation
- [op/joinmap.go](../op/joinmap.go#L37): determine rest duration
- [op/joinmap.go](../op/joinmap.go#L46): determine rest duration
- [op/notemap_test.go](../op/notemap_test.go#L82): unfinished test
- [op/random.go](../op/random.go#L50): `Replaceable`
- [op/undynamic_test.go](../op/undynamic_test.go#L23): unfinished test
- [ui/cli/app.go](../ui/cli/app.go#L76): liner handling of Ctrl+C
- [ui/cli/complete.go](../ui/cli/complete.go#L39): closest completion
- [ui/img/draw_test.go](../ui/img/draw_test.go#L86): recorded test fixture
- [ui/img/draw_test.go](../ui/img/draw_test.go#L88): stored recording
- [ui/img/notes_view.go](../ui/img/notes_view.go#L19): BIAB support
