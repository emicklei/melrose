/* Melrōse web frontend: Monaco editor + the Melrōse HTTP API. */

const FUNCTIONS = [
    "at", "bare", "bars", "beats", "biab", "bpm", "channel", "chord", "chordsequence",
    "device", "duration", "dynamic", "dynamicmap", "euclidean", "export", "fraction",
    "fractionmap", "group", "if", "import", "index", "interval", "iterator", "join",
    "joinmap", "key", "knob", "listen", "loop", "map", "merge", "midi", "midi_send",
    "multitrack", "next", "notemap", "note", "octave", "octavemap", "onbar", "onkey",
    "onoff", "play", "print", "prob", "progression", "random", "record", "repeat",
    "replace", "resequence", "reverse", "rotate", "scale", "sequence", "set", "stop",
    "stretch", "sync", "tabs", "track", "transpose", "transposemap", "trim", "undynamic",
    "ungroup", "value", "velocitymap"
];

const STORAGE_KEY = "melrose.web.source";
const FILENAME = "melrose-web.mel";

const DEFAULT_SOURCE = [
    "// Melrōse",
    "//",
    "// cmd/ctrl + e : evaluate",
    "// cmd/ctrl + 3 : play",
    "// cmd/ctrl + 5 : stop",
    "",
    'bpm(120)',
    "",
    's1 = sequence("C E G")',
    "",
    "play(s1)",
    ""
].join("\n");

const outputEl = document.getElementById("output");
let editor = null;
const stoppableMarkers = new Map();

function log(action, message, isError) {
    const entry = document.createElement("div");
    entry.className = "entry" + (isError ? " error" : "");
    const tag = document.createElement("span");
    tag.className = "action";
    tag.textContent = "[" + action + "]";
    entry.appendChild(tag);
    entry.appendChild(document.createTextNode(message));
    outputEl.appendChild(entry);
    outputEl.scrollTop = outputEl.scrollHeight;
}

/* Returns the selected text, or the contiguous block of non-empty lines
    around the cursor, together with its one-based end line. */
function currentStatement() {
    const model = editor.getModel();
    const selection = editor.getSelection();
    if (selection && !selection.isEmpty()) {
          return { source: model.getValueInRange(selection), line: selection.endLineNumber };
    }
    const cursor = editor.getPosition().lineNumber;
    const isBlank = (n) => model.getLineContent(n).trim().length === 0;
    if (isBlank(cursor)) {
        return { source: "", line: cursor };
    }
    let start = cursor;
    while (start > 1 && !isBlank(start - 1)) start--;
    let end = cursor;
    while (end < model.getLineCount() && !isBlank(end + 1)) end++;
    const lines = [];
    for (let n = start; n <= end; n++) lines.push(model.getLineContent(n));
    return { source: lines.join("\n"), line: end };
}

async function perform(action, targetLine) {
    const breakpointLine = action === "stop" ? (targetLine ?? Array.from(stoppableMarkers.keys()).pop()) : undefined;
    const { source, line } = breakpointLine === undefined
        ? currentStatement()
        : { source: editor.getModel().getLineContent(breakpointLine), line: breakpointLine };
    if (source.trim().length === 0) {
        log(action, "nothing to send", true);
        return;
    }
    const url = location.origin + "/v1/statements?action=" + encodeURIComponent(action) +
        "&file=" + encodeURIComponent(FILENAME) + "&line=" + line;
    let response;
    try {
        response = await fetch(url, {
            method: "POST",
            // text/plain keeps this a CORS simple request (no preflight).
            headers: { "Content-Type": "text/plain" },
            body: source
        });
    } catch (err) {
        log(action, "cannot reach " + location.origin + " : " + err, true);
        return;
    }
    let result;
    const raw = await response.text();
    try {
        result = JSON.parse(raw);
    } catch (e) {
        log(action, raw || response.statusText, !response.ok);
        return;
    }
    const failed = result["is-error"] === true || !response.ok;
    if (!failed) {
        const previous = stoppableMarkers.get(result.line) || [];
        const markers = action !== "stop" && result.stoppable === true ? [{
            range: new monaco.Range(result.line, 1, result.line, 1),
            options: { glyphMarginClassName: "stoppable-glyph", glyphMarginHoverMessage: { value: "Stop" } }
        }] : [];
        const updated = editor.deltaDecorations(previous, markers);
        if (markers.length) stoppableMarkers.set(result.line, updated);
        else stoppableMarkers.delete(result.line);
    }
    const parts = [];
    //if (result.type) parts.push(result.type);
    if (result.message) parts.push(result.message);
    if (!parts.length && result.object != null) parts.push(JSON.stringify(result.object));
    log(action, parts.join(" : ") || "ok", failed);
}

async function fetchVersion() {
    const el = document.getElementById("version");
    try {
        const response = await fetch(location.origin + "/version");
        const info = await response.json();
        el.textContent = "version " + info.BuildTag;
    } catch (err) {
        el.textContent = "offline";
    }
}

function registerMelroseLanguage(monaco) {
    monaco.languages.register({ id: "melrose" });
    monaco.languages.setMonarchTokensProvider("melrose", {
        keywords: FUNCTIONS,
        tokenizer: {
            root: [
                [/\/\/.*$/, "comment"],
                [/"([^"\\]|\\.)*"/, "string"],
                [/[a-z_][\w]*(?=\s*\()/, { cases: { "@keywords": "keyword", "@default": "identifier" } }],
                [/[a-zA-Z_]\w*/, "identifier"],
                [/\d+(\.\d+)?/, "number"],
                [/[=(),.]/, "delimiter"]
            ]
        }
    });
    monaco.languages.setLanguageConfiguration("melrose", {
        comments: { lineComment: "//" },
        brackets: [["(", ")"], ["[", "]"]],
        autoClosingPairs: [
            { open: "(", close: ")" },
            { open: "[", close: "]" },
            { open: '"', close: '"' }
        ]
    });
    monaco.languages.registerCompletionItemProvider("melrose", {
        provideCompletionItems: (model, position) => {
            const word = model.getWordUntilPosition(position);
            const range = {
                startLineNumber: position.lineNumber,
                endLineNumber: position.lineNumber,
                startColumn: word.startColumn,
                endColumn: word.endColumn
            };
            return {
                suggestions: FUNCTIONS.map((name) => ({
                    label: name,
                    kind: monaco.languages.CompletionItemKind.Function,
                    insertText: name + "($0)",
                    insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                    range: range
                }))
            };
        }
    });
}

require.config({ paths: { vs: "https://cdn.jsdelivr.net/npm/monaco-editor@0.52.2/min/vs" } });
require(["vs/editor/editor.main"], function () {
    registerMelroseLanguage(monaco);

    editor = monaco.editor.create(document.getElementById("editor"), {
        value: localStorage.getItem(STORAGE_KEY) || DEFAULT_SOURCE,
        language: "melrose",
        theme: "vs-dark",
        automaticLayout: true,
        glyphMargin: true,
        fontSize: 14,
        minimap: { enabled: false },
        scrollBeyondLastLine: false
    });

    editor.onDidChangeModelContent(() => {
        localStorage.setItem(STORAGE_KEY, editor.getValue());
    });

    editor.onMouseDown((event) => {
        if (event.target.type !== monaco.editor.MouseTargetType.GUTTER_GLYPH_MARGIN) return;
        const line = event.target.position?.lineNumber;
        if (stoppableMarkers.has(line)) perform("stop", line);
    });

    const bind = (keyCode, action) => {
        const keybindings = [monaco.KeyMod.CtrlCmd | keyCode];
        if (navigator.platform.startsWith("Mac")) keybindings.push(monaco.KeyMod.WinCtrl | keyCode);
        editor.addAction({
            id: "melrose." + action,
            label: "Melrōse: " + action,
            keybindings: keybindings,
            run: () => perform(action)
        });
    };
    bind(monaco.KeyCode.KeyE, "eval");
    bind(monaco.KeyCode.Digit3, "play");
    bind(monaco.KeyCode.Digit5, "stop");

    editor.focus();
});

document.getElementById("btn-eval").onclick = () => perform("eval");
document.getElementById("btn-play").onclick = () => perform("play");
document.getElementById("btn-stop").onclick = () => perform("stop");
document.getElementById("btn-clear").onclick = () => { outputEl.textContent = ""; };

// Same shortcuts while focus is outside the editor; Monaco handles them itself.
window.addEventListener("keydown", (e) => {
    if (!(e.metaKey || e.ctrlKey)) return;
    if (editor && editor.hasTextFocus()) return;
    const action = { "3": "play", "5": "stop", "e": "eval" }[e.key.toLowerCase()];
    if (!action) return;
    e.preventDefault();
    perform(action);
}, true);

fetchVersion();
