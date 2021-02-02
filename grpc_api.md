# Todoer gRPC API

- [Core Concepts](#core-concepts)
- [Error Handling](#error-handling)
- [Todo List](#todo-list)
    - [Creating a todo list](#creating-a-todo-list)
    - [Retrieving a todo list](#retrieving-a-todo-list)
    - [Retrieving all todo lists](#retrieving-all-todo-lists)
    - [Updating a todo list](#updating-a-todo-list)
    - [Deleting a todo list](#deleting-a-todo-list)
- [Todo](#todo)
    - [Creating a todo](#creating-a-todo)
    - [Retrieving a todo](#retrieving-a-todo)
    - [Retrieving all todo's from a todo list](#retrieving-all-todos-from-a-todo-list)
    - [Updating a todo](#updating-a-todo)
    - [Deleting a todo](#deleting-a-todo)

The todoer API provides services related to todos, like
creating todo lists.

## Core Concepts

This API uses [protocol buffers](https://developers.google.com/protocol-buffers)
as the representation for resources. All the API functions and message types should
be specified on [todoer.proto](pb/todoer.proto).

## Error Handling

When an error occurs you can always expect an error message from the function
giving some more information on what went wrong (when appropriate).

## Todo List

A `todolist` object is the list containing `todo`s.

```protobuf
message TodoList {
  uint32 id = 1;
  string title = 2;
}
```

Fields:
- `id`: The identifier for this todo list;
- `title`: Title of the todo list;
  - **Can not be empty**;

### Creating a todo list

To create a todo list, use the following function:

```
  rpc CreateTodoList (CreateTodoListRequest) returns (CreateTodoListReply) {}
```

With the following request object:

```protobuf
message CreateTodoListRequest {
  string title = 1;
}
```

Example of Go request object:

```go
CreateTodoListRequest{
    Title: "Routine",
}
```

In case of success you can expect the following response:

```protobuf
message CreateTodoListReply {
  TodoList todo_list = 1;
}
```

Example of Go response object:

```go
CreateTodoListReply{
    TodoList: TodoList{
        Id:    0,
        Title: "Routine",
    },
}
```

### Retrieving a todo list

To retrieve a todo list, use the following function:

```
  rpc GetTodoList (GetTodoListRequest) returns (GetTodoListReply) {}
```

With the following request object:

```protobuf
message GetTodoListRequest {
  uint32 id = 1;
}
```

Example of Go request object:

```go
GetTodoListRequest{
    Id: 0,
}
```

In case of success you can expect the following response:

```protobuf
message GetTodoListReply {
  TodoList todo_list = 1;
}
```

Example of Go response object:

```go
GetTodoListReply{
    TodoList: TodoList{
        Id:    0,
        Title: "Routine",
    },
}
```

### Retrieving all todo lists

To get a todo list, use the following function:

```
  rpc GetAllTodoLists (Empty) returns (GetAllTodoListsReply) {}
```

In case of success you can expect the following response:

```protobuf
message GetAllTodoListsReply {
  repeated TodoList todo_lists = 1;
}
```

Example of Go response object:

```go
GetAllTodoListsReply{
    TodoLists: []TodoList{
        TodoList{
            Id:    0,
            Title: "Routine",
        },
        TodoList{
            Id:    1,
            Title: "Work",
        },
    },
}
```

### Updating a todo list

To update a todo list, use the following function:

```
  rpc UpdateTodoList (UpdateTodoListRequest) returns (Empty) {}
```

With the following request object:

```protobuf
message UpdateTodoListRequest {
  TodoList todo_list = 1;
}
```

Example of Go request object:

```go
UpdateTodoListRequest{
    TodoList: TodoList{
        Id:    0,
        Title: "Rot",
    },
}
```

In case of success you can expect no error to be returned.

### Deleting a todo list

To delete a todo list, use the following function:

```
  rpc DeleteTodoList (DeleteTodoListRequest) returns (Empty) {}
```

With the following request object:

```protobuf
message DeleteTodoListRequest {
  uint32 id = 1;
}
```

Example of Go request object:

```go
DeleteTodoListRequest{
    Id: 0,
}
```

In case of success you can expect no error to be returned.

## Todo

A `todo` object is what contains details about a task that you need **to do**.

```protobuf
message Todo {
  uint32 id = 1;
  uint32 list_id = 2;
  string description = 3;
  string comments = 4;
  string due_date = 5;
  repeated string labels = 6;
  bool done = 7;
}
```

Fields:
- `id`: The identifier for this todo item;
- `list_id`: The `id` from a [TodoList](#todo-list) which this todo is associated;
- `description`: Description of task you need to do;
  - **Can not be empty**;
- `comments`: General comments about the task;
- `due_date`: The final date when this task needs to be done;
  - You can expect an string representation
    of date following the [RFC 3339](https://tools.ietf.org/html/rfc3339),
    for example: "2021-01-01T00:00:01Z".
- `labels`: A list of labels to mark your task with;
- `done`: If the task is done or not.

### Creating a todo

To create a todo, use the following function:

```
  rpc CreateTodo (CreateTodoRequest) returns (CreateTodoReply) {}
```

With the following request object:

```protobuf
message CreateTodoRequest {
  uint32 list_id = 1;
  string description = 2;
  string comments = 4;
  string due_date = 5;
  repeated string labels = 6;
  bool done = 7;
}
```

Example of Go request object:

```go
CreateTodoRequest{
    ListId:      0,
    Description: "Make the bed."
    Comments:    "Will be easy",
    DueDate:     "2021-01-01T00:00:01Z",
    Labels:      []string{"bed", "bedroom"},
    Done:        false,
}
```

In case of success you can expect an the following reply:

```protobuf
message CreateTodoReply {
  Todo todo = 1;
}
```

Example of Go reply object:

```go
CreateTodoReply{
    Todo: Todo{
        Id:          0,
        ListId:      0,
        Description: "Make the bed."
        Comments:    "Will be easy",
        DueDate:     "2021-01-01T00:00:01Z",
        Labels:      []string{"bed", "bedroom"},
        Done:        false,
    },
}
```

### Retrieving a todo

To retrieve a todo, use the following function:

```
  rpc GetTodo (GetTodoRequest) returns (GetTodoReply) {}
```

With the following request object:

```protobuf
message GetTodoRequest {
  uint32 id = 1;
}
```

Example of Go request object:

```go
GetTodoRequest{
    Id:      0,
}
```

In case of success you can expect an the following reply:

```protobuf
message GetTodoReply {
  Todo todo = 1;
}
```

Example of Go reply object:

```go
GetTodoReply{
    Todo: Todo{
        Id:          0,
        ListId:      0,
        Description: "Make the bed."
        Comments:    "Will be easy",
        DueDate:     "2021-01-01T00:00:01Z",
        Labels:      []string{"bed", "bedroom"},
        Done:        false,
    },
}
```

### Retrieving all todo's from a todo list

To retrieve all todo's from a list, use the following function:

```
  rpc GetTodosByList (GetTodosByListRequest) returns (GetTodosByListReply) {}
```

With the following request object:

```protobuf
message GetTodosByListRequest {
  uint32 list_id = 1;
}
```

Example of Go request object:

```go
GetTodosByListRequest{
    ListId:      0,
}
```

In case of success you can expect an the following reply:

```protobuf
message GetTodosByListReply {
  repeated Todo todos = 1;
}
```

Example of Go reply object:

```go
GetTodosByListReply{
    Todos: []Todo{
        Todo{
            Id:          0,
            ListId:      0,
            Description: "Make the bed."
            Comments:    "Will be easy",
            DueDate:     "2021-01-01T00:00:01Z",
            Labels:      []string{"bed", "bedroom"},
            Done:        true,
        },
        Todo{
            Id:          1,
            ListId:      0,
            Description: "Wash the dishes."
            Comments:    "Will be hard",
            DueDate:     "2021-01-03T00:00:01Z",
            Labels:      []string{"sink", "kitchen"},
            Done:        false,
        },
    },
}
```

### Updating a todo

To update a todo, use the following function:

```
  rpc UpdateTodo (UpdateTodoRequest) returns (Empty) {}
```

With the following request object:

```protobuf
message UpdateTodoRequest {
  Todo todo = 1;
}
```

Example of Go request object:

```go
UpdateTodoRequest{
    Todo: Todo{
        Id:          0,
        ListId:      0,
        Description: "Make the bed."
        Comments:    "Was easy",
        DueDate:     "2021-01-01T00:00:01Z",
        Labels:      []string{"bed", "bedroom"},
        Done:        done,
    }
}
```

In case of success you can expect no error to be returned.

### Deleting a todo

To delete a todo, use the following function:

```
  rpc DeleteTodo (DeleteTodoRequest) returns (Empty) {}
```

With the following request object:

```protobuf
message DeleteTodoRequest {
  Todo todo = 1;
}
```

Example of Go request object:

```go
DeleteTodoRequest{
    Id:          0,
    ListId:      0,
}
```

In case of success you can expect no error to be returned.
