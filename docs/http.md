---
title: Melrōse HTTP server
---

[Home](https://emicklei.github.io/melrose)


# HTTP API server

Melrōse starts a HTTP server on port 8118 and evaluates expressions or programs providing the source as the payload (HTTP Body).
This server is used by the [Melrōse Plugin for Visual Studio Code](https://github.com/emicklei/melrose-for-vscode).


## HTTP Request

    POST http://localhost:8118/v1/statements?action={action}

## HTTP response

### 200 OK

If the request was successful processed then the response looks like:

    {
        "type": "core.Note",
        "is-error": false,
        "message": "note('C')",
        "file": "",
        "line": 0,
        "column": 0
    }

### 500 Internal Server Error

If the request could not be processed then the response looks like:

    {
        "type": "*file.Error",
        "is-error": true,
        "message": "literal not terminated (1:9)\n | note(C\")\n | ........^",
        "file": "yours.mel",
        "line": 1,
        "column": 8
    }

### 400 Bad Request

If the request is malformed then the response will have the error message.

## HTTP Request parameters

### tracing

If the HTTP URL has the query parameter `debug=true` then `melrōse` will produce extra logging.

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

### version

    GET http://localhost:8118/version

return version information such as

    {"APIVersion":"v1","SyntaxVersion":"0.30","BuildTag":"v1.0.1"}