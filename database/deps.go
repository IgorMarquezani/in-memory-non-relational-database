package database

import (
	"app/instruction"
)

var InstructionQueue = make(chan instruction.Instruction, 1000)
