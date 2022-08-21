# HTTP API server

Melrōse starts a HTTP server on port 8118 and evaluates expressions or programs providing the source as the payload (HTTP Body).
This server is used by the [Melrōse Plugin for Visual Studio Code](https://github.com/emicklei/melrose-for-vscode).


The API is documented in `openapi.yaml` which can be [viewed](https://editor-next.swagger.io/?url=https://raw.githubusercontent.com/emicklei/melrose/master/docs/openapi.yaml).


### example

    curl -d "note('c')" http://localhost:8118/v1/statements?action=play

### 200 OK

If the request was successful then the response looks like:

    {
        "type": "core.Note",
        "is-error": false,
        "message": "note('C')",
        "file": "",
        "line": 0,
        "column": 0
    }