package stack

type Stack[T any] struct {
	items []T
}

func (stack *Stack[T]) Size() int {
	return len(stack.items)
}

func (stack *Stack[T]) Push(item T) {
	stack.items = append(stack.items, item)
}

func (stack *Stack[T]) Pop() T {
	i := len(stack.items) - 1
	val := stack.items[i]
	stack.items = stack.items[:i]
	return val
}

func (stack *Stack[T]) Peek() T {
	i := len(stack.items) - 1
	val := stack.items[i]
	return val
}
