package persistence

type LocalStorage struct {
	todoListAutoincrement uint32
	todoAutoincrement     uint32
	TodoListTable         map[uint32]TodoList
	TodoTable             map[uint32]Todo
	// This relationship maps each TodoList ID to the respective list of Todo ID's
	TodoListRelationship map[uint32][]uint32
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		todoListAutoincrement: 0,
		todoAutoincrement:     0,
		TodoListTable:         map[uint32]TodoList{},
		TodoTable:             map[uint32]Todo{},
		TodoListRelationship:  map[uint32][]uint32{},
	}
}

func (ls *LocalStorage) InsertTodoList(todoList TodoList) error {
	ls.TodoListTable[ls.todoListAutoincrement] = todoList
	ls.todoListAutoincrement++
	return nil
}

func (ls *LocalStorage) InsertTodo(todo Todo) error {
	if _, ok := ls.TodoListTable[todo.ListID]; !ok {
		return ErrTodoListNotFound
	}
	ls.TodoListRelationship[todo.ListID] = append(ls.TodoListRelationship[todo.ListID], todo.ID)
	return nil
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
	return nil, nil
}

func (ls *LocalStorage) UpdateTodoList(todoList *TodoList) (*TodoList, error) {
	return nil, nil
}

func (ls *LocalStorage) UpdateTodo(todo *Todo) (*Todo, error) {
	return nil, nil
}

func (ls *LocalStorage) DeleteTodoList(todoList *TodoList) error {
	return nil
}

func (ls *LocalStorage) DeleteTodo(todo *Todo) error {
	return nil
}
