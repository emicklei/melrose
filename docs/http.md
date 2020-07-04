---
title: Melrōse HTTP server
---

# HTTP API server

Melrōse starts a HTTP server on port 8118 and evaluates programs on `POST /v1/statements` providing the source as the payload (HTTP Body).
This server is used by the [Melrōse Plugin for Visual Studio Code](https://github.com/emicklei/melrose-for-vscode).

### HTTP response

#### 200 OK

If the request was successful processed then the response looks like:

    {
        "type": "core.Note",
        "is-error": false,
        "message": "note('C')",
        "file": "",
        "line": 0,
        "column": 0,
        "object": null
    }

#### 500 Internal Server Error

If the request could not be processed then the response looks like:

    {
        "type": "*file.Error",
        "is-error": true,
        "message": "literal not terminated (1:9)\n | note(C\")\n | ........^",
        "file": "yours.mel",
        "line": 1,
        "column": 0,
        "object": {
            "Line": 1,
            "Column": 8,
            "Message": "literal not terminated",
            "Snippet": "\n | note(C\")\n | ........^"
        }
    }

#### 400 Bad Request

If the request is malformed then the response will have the error message.

### HTTP port

The port can be changed to e.g. 8000 with the program option `-http :8000`.

### tracing

If the HTTP URL has the query parameter `trace=true` then `melrōse` will produce extra logging.

### play

If the HTTP URL has the query parameter `action=play` then `melrōse` will try to play the result of the selected expression(s).

### begin

If the HTTP URL has the query parameter `action=begin` then `melrōse` will try to `begin` the loop of the selected expression.

### end

If the HTTP URL has the query parameter `action=end` then `melrōse` will try to `end` the loop of the selected expression.

### kill

If the HTTP URL has the query parameter `action=kill` then `melrōse` will stop playing any sound.

### inspecting

If the HTTP URL has the query parameter `action=inspect` then `melrōse` will print inspection details of the selected expression.