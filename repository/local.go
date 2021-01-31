package repository

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

// TodoList

func (ls *LocalStorage) InsertTodoList(todoList TodoList) (*TodoList, error) {
	if todoList.Title == "" {
		return nil, ErrEmptyTitle
	}
	todoList.ID = ls.TodoListAutoincrement
	ls.TodoListTable[ls.TodoListAutoincrement] = todoList
	ls.TodoListAutoincrement++
	return &todoList, nil
}

func (ls *LocalStorage) GetAllTodoLists() ([]TodoList, error) {
	result := []TodoList{}
	for _, todoList := range ls.TodoListTable {
		result = append(result, todoList)
	}

	return result, nil
}

func (ls *LocalStorage) GetTodoListByID(id uint32) (*TodoList, error) {
	todoList, ok := ls.TodoListTable[id]
	if !ok {
		return nil, ErrTodoListNotFound
	}

	return &todoList, nil
}

func (ls *LocalStorage) UpdateTodoList(todoList TodoList) error {
	if _, ok := ls.TodoListTable[todoList.ID]; !ok {
		return ErrTodoListNotFound
	}

	if todoList.Title == "" {
		return ErrEmptyTitle
	}

	ls.TodoListTable[todoList.ID] = todoList
	return nil
}

func (ls *LocalStorage) DeleteTodoListByID(id uint32) error {
	if _, ok := ls.TodoListTable[id]; !ok {
		return ErrTodoListNotFound
	}

	delete(ls.TodoListTable, id)
	return nil
}

// Todo

func (ls *LocalStorage) InsertTodo(todo Todo) (*Todo, error) {
	if _, ok := ls.TodoListTable[todo.ListID]; !ok {
		return nil, ErrTodoListNotFound
	}

	if todo.Description == "" {
		return nil, ErrEmptyDescription
	}

	todo.ID = ls.TodoAutoincrement
	ls.TodoListRelationship[todo.ListID] = append(ls.TodoListRelationship[todo.ListID], todo.ID)
	ls.TodoTable[ls.TodoAutoincrement] = todo
	ls.TodoAutoincrement++
	return &todo, nil
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

	todos := []Todo{}
	todoIDs, ok := ls.TodoListRelationship[listID]
	if !ok {
		return todos, nil
	}

	for _, id := range todoIDs {
		todo := ls.TodoTable[id]
		todos = append(todos, todo)
	}

	return todos, nil
}

func (ls *LocalStorage) UpdateTodo(todo Todo) error {
	if _, ok := ls.TodoTable[todo.ID]; !ok {
		return ErrTodoNotFound
	}

	if todo.Description == "" {
		return ErrEmptyDescription
	}

	ls.TodoTable[todo.ID] = todo
	return nil
}

func (ls *LocalStorage) DeleteTodo(todo Todo) error {
	if _, ok := ls.TodoTable[todo.ID]; !ok {
		return ErrTodoNotFound
	}
	delete(ls.TodoTable, todo.ID)
	// If there is a relationship between the list, delete from there
	todoListRelationship, ok := ls.TodoListRelationship[todo.ListID]
	if ok {
		newRelationship := removeID(todoListRelationship, todo.ID)
		if len(newRelationship) > 0 {
			ls.TodoListRelationship[todo.ListID] = newRelationship
		} else {
			delete(ls.TodoListRelationship, todo.ListID)
		}
	}

	return nil
}

func removeID(oldTodoIDs []uint32, todoID uint32) []uint32 {
	newTodoIDs := []uint32{}
	for _, id := range oldTodoIDs {
		if id != todoID {
			newTodoIDs = append(newTodoIDs, id)
		}
	}
	return newTodoIDs
}
