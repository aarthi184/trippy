package atkins

import (
	"trippy/slotmachine"
)

const (
	_EMPTY slotmachine.Symbol = iota
	_ATKINS
	_STEAK
	_HAM
	_BUFFALO_WINGS
	_SAUSAGE
	_EGGS
	_BUTTER
	_CHEESE
	_BACON
	_MAYONNAISE
	_SCALE

	_SCATTER_COUNT_FOR_FREE_SPIN = 3
	_FREE_SPINS                  = 10
)

var (
	PayTable = slotmachine.PayTable{
		_ATKINS: slotmachine.Pays{
			5: 5000,
			4: 500,
			3: 50,
			2: 5,
		},
		_STEAK: slotmachine.Pays{
			5: 1000,
			4: 200,
			3: 40,
			2: 3,
		},
		_HAM: slotmachine.Pays{
			5: 500,
			4: 150,
			3: 30,
			2: 2,
		},
		_BUFFALO_WINGS: slotmachine.Pays{
			5: 300,
			4: 100,
			3: 25,
			2: 2,
		},
		_SAUSAGE: slotmachine.Pays{
			5: 200,
			4: 75,
			3: 20,
		},
		_EGGS: slotmachine.Pays{
			5: 200,
			4: 75,
			3: 20,
		},
		_BUTTER: slotmachine.Pays{
			5: 100,
			4: 50,
			3: 15,
		},
		_CHEESE: slotmachine.Pays{
			5: 100,
			4: 50,
			3: 15,
		},
		_BACON: slotmachine.Pays{
			5: 50,
			4: 25,
			3: 10,
		},
		_MAYONNAISE: slotmachine.Pays{
			5: 50,
			4: 25,
			3: 10,
		},
	}

	Reels = slotmachine.Reels{
		{_SCALE, _MAYONNAISE, _HAM, _HAM, _BACON},
		{_MAYONNAISE, _BUFFALO_WINGS, _BUTTER, _CHEESE, _SCALE},
		{_HAM, _STEAK, _EGGS, _ATKINS, _STEAK},
		{_SAUSAGE, _SAUSAGE, _SCALE, _SCALE, _HAM},
		{_BACON, _CHEESE, _CHEESE, _BUTTER, _CHEESE},
		{_EGGS, _MAYONNAISE, _MAYONNAISE, _BACON, _SAUSAGE},
		{_CHEESE, _HAM, _BUTTER, _CHEESE, _BUTTER},
		{_MAYONNAISE, _BUTTER, _HAM, _SAUSAGE, _BACON},
		{_SAUSAGE, _BACON, _SAUSAGE, _STEAK, _BUFFALO_WINGS},
		{_BUTTER, _STEAK, _BACON, _EGGS, _CHEESE},
		{_BUFFALO_WINGS, _SAUSAGE, _STEAK, _BACON, _SAUSAGE},
		{_BACON, _MAYONNAISE, _BUFFALO_WINGS, _MAYONNAISE, _HAM},
		{_EGGS, _HAM, _BUTTER, _SAUSAGE, _BUTTER},
		{_MAYONNAISE, _ATKINS, _MAYONNAISE, _CHEESE, _STEAK},
		{_STEAK, _BUTTER, _CHEESE, _BUTTER, _MAYONNAISE},
		{_BUFFALO_WINGS, _EGGS, _SAUSAGE, _HAM, _EGGS},
		{_BUTTER, _CHEESE, _EGGS, _MAYONNAISE, _SAUSAGE},
		{_CHEESE, _BACON, _BACON, _BACON, _HAM},
		{_EGGS, _SAUSAGE, _MAYONNAISE, _BUFFALO_WINGS, _ATKINS},
		{_ATKINS, _BUFFALO_WINGS, _BUFFALO_WINGS, _SAUSAGE, _BUTTER},
		{_BACON, _SCALE, _HAM, _CHEESE, _BUFFALO_WINGS},
		{_MAYONNAISE, _MAYONNAISE, _SAUSAGE, _EGGS, _MAYONNAISE},
		{_HAM, _BUTTER, _BACON, _BUTTER, _EGGS},
		{_CHEESE, _CHEESE, _CHEESE, _BUFFALO_WINGS, _HAM},
		{_EGGS, _BACON, _EGGS, _BACON, _BACON},
		{_SCALE, _EGGS, _ATKINS, _MAYONNAISE, _BUTTER},
		{_BUTTER, _BUFFALO_WINGS, _BUFFALO_WINGS, _EGGS, _STEAK},
		{_BACON, _MAYONNAISE, _BACON, _HAM, _MAYONNAISE},
		{_SAUSAGE, _STEAK, _BUTTER, _SAUSAGE, _SAUSAGE},
		{_BUFFALO_WINGS, _HAM, _CHEESE, _STEAK, _EGGS},
		{_STEAK, _CHEESE, _MAYONNAISE, _MAYONNAISE, _CHEESE},
		{_BUTTER, _BACON, _STEAK, _BACON, _BUFFALO_WINGS},
	}

	PayLines = slotmachine.PayLines{
		{2, 2, 2, 2, 2},
		{1, 1, 1, 1, 1},
		{3, 3, 3, 3, 3},
		{1, 2, 3, 2, 1},
		{3, 2, 1, 2, 3},
		{2, 1, 1, 1, 2},
		{2, 3, 3, 3, 2},
		{1, 1, 2, 3, 3},
		{3, 3, 2, 1, 1},
		{2, 1, 2, 3, 2},
		{2, 3, 2, 1, 2},
		{1, 2, 2, 2, 1},
		{3, 2, 2, 2, 3},
		{1, 2, 1, 2, 1},
		{3, 2, 3, 2, 3},
		{2, 2, 1, 2, 2},
		{2, 2, 3, 2, 2},
		{1, 1, 3, 1, 1},
		{3, 3, 1, 3, 3},
		{1, 3, 3, 3, 1},
	}
)
