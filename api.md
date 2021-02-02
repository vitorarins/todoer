# Todoer API

The todoer API provides services related to todos, like
creating todo lists.

## Core Concepts

For example, with this specification:

```
{
    "field" : <boolean>,
    "optionalField" : <string>(optional)
}
```

You can expect a JSON like this:

```json
{
    "field": true,
    "optionalField" : "data"
}
```

All JSON fields documented as part of request/response bodies are
to be considered obligatory, unless they are explicitly
documented as optional.

When a field has type "<date>" you can expect an string representation
of date following the [RFC 3339](https://tools.ietf.org/html/rfc3339),
for example: "2018-01-01T00:00:01Z".


## Error Handling

When an error occurs you can always expect an HTTP status code indicating the
nature of the failure and also a response body with an error message
giving some more information on what went wrong (when appropriate).

It follows this schema:

```
{
    "error": {
        "message" : <string>
    }
}
```

The **message** is intended for human inspection, no programmatic decision
should be made using their contents. Services integrating with this API
can depend on the error response schema, but the contents of the
message itself should be handled as opaque strings.

## Todo List

A `todolist` object is the list containing `todo`s.

```json
{
    "id":    <int>,
    "title": <string>
}
```

### Creating a todo list

To create a todo list, send the following request:

```
POST /todolist
```

With the following request body:

```json
{
    "title": <string>
}
```

Example of request body:

```json
{
    "title": "Routine"
}
```

In case of success you can expect an status code 200/OK and the following response:

```json
{
    "id":    <int>,
    "title": <string>
}
```

Example of response body:

```json
{
    "id":    0,
    "title": "Routine"
}
```

### Retrieving a todo list

To retrieve a todo list, send the following request:

```
GET /todolist/{id}
```

In case of success you can expect an status code 200/OK and the following response:

```json
{
    "id":    <int>,
    "title": <string>
}
```
Example of response body:

```json
{
    "id":    0,
    "title": "Routine"
}
```


### Retrieving all todo lists

To retrieve all todo lists, send the following request:

```
GET /todolist
```

In case of success you can expect an status code 200/OK and the following response:

```json
[
    <todolist>,
    ...
]
```

Example of response body:

```json
[
    {
        "id":    0,
        "title": "Routine"
    },
    {
        "id":    1,
        "title": "Work"
    }
]
```

### Updating a todo list

To update a todo list, send the following request:

```
PUT /todolist/{id}
```

With the following request body:

```json
{
    "title": <string>
}
```

Example of request body:

```json
{
    "title": "Routine"
}
```

In case of success you can expect an status code 200/OK.

### Deleting a todo list

To delete a todo list, send the following request:

```
DELETE /todolist/{id}
```

In case of success you can expect an status code 200/OK.

## Todo

A `todo` object is what contains details about a task that you need **to do**.

```json
{
    "id":          <int>,
    "list_id":     <int>,
    "description": <string>,
    "done":        <boolean>,
    "comments":    <string>,
    "due_date":    <date>,
    "labels":      [<string>,...]
}
```

### Creating a todo

To create a todo, send the following request:

```
POST /todolist/{list_id}/todo
```

With the following request body:

```json
{
    "description": <string>,
    "done":        <boolean>,
    "comments":    <string>(optional),
    "due_date":    <date>(optional),
    "labels":      [<string>,...](optional)
}
```

Example of request body:

```json
{
    "description": "Make the bed",
    "done":        false
    "comments":    "Will be easy",
    "due_date":    "2021-02-01T00:00:01Z",
    "labels":      ["bed", "bedroom"],
}
```

In case of success you can expect an status code 200/OK and the following response:

```json
{
    "id":          <int>,
    "list_id":     <int>,
    "description": <string>,
    "done":        <boolean>,
    "comments":    <string>,
    "due_date":    <date>,
    "labels":      [<string>,...]
}
```

Example of response body:

```json
{
    "id":          0,
    "list_id":     0,
    "description": "Make the bed",
    "done":        false
    "comments":    "Will be easy",
    "due_date":    "2021-02-01T00:00:01Z",
    "labels":      ["bed", "bedroom"]
}
```

### Retrieving a todo

To retrieve a todo, send the following request:

```
GET /todolist/{list_id}/todo/{id}
```

In case of success you can expect an status code 200/OK and the following response:

```json
{
    "id":          <int>,
    "list_id":     <int>,
    "description": <string>,
    "done":        <boolean>,
    "comments":    <string>,
    "due_date":    <date>,
    "labels":      [<string>,...]
}
```
Example of response body:

```json
{
    "id":          0,
    "list_id":     0,
    "description": "Make the bed",
    "done":        false
    "comments":    "Will be easy",
    "due_date":    "2021-02-01T00:00:01Z",
    "labels":      ["bed", "bedroom"]
}
```


### Retrieving all todo's from a todo list

To retrieve all todo's from a todo list, send the following request:

```
GET /todolist/{list_id}/todo
```

In case of success you can expect an status code 200/OK and the following response:

```json
[
    <todo>,
    ...
]
```
Example of response body:

```json
[
    {
        "id":          0,
        "list_id":     0,
        "description": "Make the bed",
        "comments":    "Will be easy",
        "due_date":    "2021-02-01T00:00:01Z",
        "labels":      ["bed", "bedroom"],
        "done":        false
    },
    {
        "id":          1,
        "list_id":     0,
        "description": "Type stuff",
        "comments":    "Will be hard",
        "due_date":    "2021-02-01T00:00:01Z",
        "labels":      ["computer", "office"],
        "done":        false
    }
]
```

### Updating a todo

To update a todo, send the following request:

```
PUT /todolist/{list_id}/todo/{id}
```

With the following request body:

```json
{
    "description": <string>,
    "done":        <boolean>,
    "comments":    <string>(optional),
    "due_date":    <date>(optional),
    "labels":      [<string>,...](optional)
}
```

Example of request body:

```json
{
    "description": "Make the bed",
    "done":        true
    "comments":    "Was easy",
    "due_date":    "2021-02-01T00:00:01Z",
    "labels":      ["bed", "bedroom"],
}
```

In case of success you can expect an status code 200/OK.

### Deleting a todo

To delete a todo, send the following request:

```
DELETE /todolist/{list_id}/todo/{id}
```

In case of success you can expect an status code 200/OK.
