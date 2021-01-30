package persistence

type LocalStorage struct {
	TodoListAutoincrement uint32
	TodoAutoincrement     uint32
	TodoListTable         map[uint32]TodoList
	TodoTable             map[uint32]Todo
	// This relationship maps each TodoList ID to the respective list of Todo ID's
	TodoListRelationship map[uint32][]uint32
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		TodoListAutoincrement: 0,
		TodoAutoincrement:     0,
		TodoListTable:         map[uint32]TodoList{},
		TodoTable:             map[uint32]Todo{},
		TodoListRelationship:  map[uint32][]uint32{},
	}
}

func (ls *LocalStorage) InsertTodoList(todoList TodoList) (*TodoList, error) {
	todoList.ID = ls.TodoListAutoincrement
	ls.TodoListTable[ls.TodoListAutoincrement] = todoList
	ls.TodoListAutoincrement++
	return &todoList, nil
}

func (ls *LocalStorage) InsertTodo(todo Todo) (*Todo, error) {
	if _, ok := ls.TodoListTable[todo.ListID]; !ok {
		return nil, ErrTodoListNotFound
	}
	ls.TodoListRelationship[todo.ListID] = append(ls.TodoListRelationship[todo.ListID], todo.ID)
	todo.ID = ls.TodoAutoincrement
	ls.TodoTable[ls.TodoAutoincrement] = todo
	ls.TodoAutoincrement++
	return &todo, nil
}

func (ls *LocalStorage) GetTodoListByID(id uint32) (*TodoList, error) {
	todoList, ok := ls.TodoListTable[id]
	if !ok {
		return nil, ErrTodoListNotFound
	}

	return &todoList, nil
}

func (ls *LocalStorage) GetTodoByID(id uint32) (*Todo, error) {
	todo, ok := ls.TodoTable[id]
	if !ok {
		return nil, ErrTodoNotFound
	}

	return &todo, nil
}

func (ls *LocalStorage) GetTodosByListID(listID uint32) ([]Todo, error) {
	if _, ok := ls.TodoListTable[listID]; !ok {
		return nil, ErrTodoListNotFound
	}
	todoIDs, ok := ls.TodoListRelationship[listID]
	if !ok {
		return nil, ErrEmptyTodoList
	}

	todos := []Todo{}
	for _, id := range todoIDs {
		todo := ls.TodoTable[id]
		todos = append(todos, todo)
	}

	return todos, nil
}

func (ls *LocalStorage) UpdateTodoList(todoList TodoList) error {
	if _, ok := ls.TodoListTable[todoList.ID]; !ok {
		return ErrTodoListNotFound
	}

	ls.TodoListTable[todoList.ID] = todoList
	return nil
}

func (ls *LocalStorage) UpdateTodo(todo Todo) error {
	if _, ok := ls.TodoTable[todo.ID]; !ok {
		return ErrTodoNotFound
	}

	ls.TodoTable[todo.ID] = todo
	return nil
}

func (ls *LocalStorage) DeleteTodoList(todoList TodoList) error {
	if _, ok := ls.TodoListTable[todoList.ID]; !ok {
		return ErrTodoListNotFound
	}

	delete(ls.TodoListTable, todoList.ID)
	return nil
}

func (ls *LocalStorage) DeleteTodo(todo Todo) error {
	if _, ok := ls.TodoTable[todo.ID]; !ok {
		return ErrTodoNotFound
	}
	if _, ok := ls.TodoListTable[todo.ListID]; !ok {
		return ErrTodoListNotFound
	}

	delete(ls.TodoTable, todo.ID)
	return nil
}
