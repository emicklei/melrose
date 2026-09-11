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
const ENDPOINT_KEY = "melrose.web.endpoint";
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

const endpointInput = document.getElementById("endpoint");
const outputEl = document.getElementById("output");
let editor = null;

// When served by melrose itself the API lives on the same origin.
if (location.protocol.startsWith("http")) {
    endpointInput.value = location.origin;
}

function endpoint() {
    return endpointInput.value.replace(/\/+$/, "");
}

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
   around the cursor, together with its one-based start line. */
function currentStatement() {
    const model = editor.getModel();
    const selection = editor.getSelection();
    if (selection && !selection.isEmpty()) {
        return { source: model.getValueInRange(selection), line: selection.startLineNumber };
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
    return { source: lines.join("\n"), line: start };
}

async function perform(action) {
    const { source, line } = currentStatement();
    if (source.trim().length === 0) {
        log(action, "nothing to send", true);
        return;
    }
    const url = endpoint() + "/v1/statements?action=" + encodeURIComponent(action) +
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
        log(action, "cannot reach " + endpoint() + " : " + err, true);
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
    const parts = [];
    if (result.type) parts.push(result.type);
    if (result.message) parts.push(result.message);
    if (!parts.length && result.object != null) parts.push(JSON.stringify(result.object));
    log(action, parts.join(" : ") || "ok", failed);
}

async function fetchVersion() {
    const el = document.getElementById("version");
    try {
        const response = await fetch(endpoint() + "/version");
        const info = await response.json();
        el.textContent = "api " + info.APIVersion + " · syntax " + info.SyntaxVersion + " · " + info.BuildTag;
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
        fontSize: 14,
        minimap: { enabled: false },
        scrollBeyondLastLine: false
    });

    editor.onDidChangeModelContent(() => {
        localStorage.setItem(STORAGE_KEY, editor.getValue());
    });

    const bind = (keybinding, action) => {
        editor.addAction({
            id: "melrose." + action,
            label: "Melrōse: " + action,
            keybindings: [keybinding],
            run: () => perform(action)
        });
    };
    bind(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyE, "eval");
    bind(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Digit3, "play");
    bind(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Digit5, "stop");

    editor.focus();
});

document.getElementById("btn-eval").onclick = () => perform("eval");
document.getElementById("btn-play").onclick = () => perform("play");
document.getElementById("btn-stop").onclick = () => perform("stop");
document.getElementById("btn-clear").onclick = () => { outputEl.textContent = ""; };

endpointInput.value = localStorage.getItem(ENDPOINT_KEY) || endpointInput.value;
endpointInput.onchange = () => {
    localStorage.setItem(ENDPOINT_KEY, endpoint());
    fetchVersion();
};

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
