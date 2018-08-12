package vm

import (
  "encoding/binary"
  "fmt"
)

type Frame struct {
  pc int
  stack []uint32
  locals []uint32
  returnValue uint32
}

func debugPrint(str string) {
  fmt.Println(str)
}

func ProcessInstruction(method *Method, args []uint32, vm *VirtualMachine) uint32 {
  code := method.GetCode()
  frame := Frame {
    pc : 0,
    stack : make([]uint32, 0, code.maxStack),
    locals : args,
    returnValue : 0,
  }
  for frame.pc < len(code.code) {
    op := code.code[frame.pc];
    switch op {
    case 0x04:
      iconst_1_Op(code.code, method, vm, &frame)
    case 0x10:
      bipush_Op(code.code, method, vm, &frame)
    case 0x1a:
      iload_0_Op(code.code, method, vm, &frame)
    case 0x57:
      pop_Op(code.code, method, vm, &frame)
    case 0x12:
      ldc_Op(code.code, method, vm, &frame)
    case 0x64:
      isub_Op(code.code, method, vm, &frame)
    case 0x68:
      imul_Op(code.code, method, vm, &frame)
    case 0xa0:
      if_icmpne_Op(code.code, method, vm, &frame)
    case 0xa7:
      goto_Op(code.code, method, vm, &frame)
    case 0xac:
      ireturn_Op(code.code, method, vm, &frame)
    case 0xb2:
      getstatic_Op(code.code, method, vm, &frame)
    case 0xb8:
      invokestatic_Op(code.code, method, vm, &frame)
    case 0xb1:
      return_Op(code.code, method, vm, &frame)
    default:
      panic("Unimplemented instruction: " + fmt.Sprintf("0x%02x", op))
    }
  }
  return frame.returnValue
}

func bipush_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  b := code[frame.pc + 1]
  frame.stack = append(frame.stack, uint32(b))
  debugPrint(fmt.Sprintf("bipush(%d) : stack %+v", b, frame.stack))
  frame.pc += 2
}

func iconst_1_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  frame.stack = append(frame.stack, uint32(1))
  debugPrint(fmt.Sprintf("iconst_1 : stack %+v", frame.stack))
  frame.pc += 1
}

func iload_0_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  frame.stack = append(frame.stack, frame.locals[0])
  debugPrint(fmt.Sprintf("iload_0 : stack %+v", frame.stack))
  frame.pc += 1
}

func pop_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  debugPrint(fmt.Sprintf("pop()"))
  frame.stack = frame.stack[:len(frame.stack) - 1]
  debugPrint(fmt.Sprintf("pop : stack %+v", frame.stack))
  frame.pc += 1
}

func isub_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  value2 := frame.stack[len(frame.stack) - 1]
  value1 := frame.stack[len(frame.stack) - 2]
  frame.stack[len(frame.stack) - 2] = value1 - value2
  frame.stack = frame.stack[:len(frame.stack) - 1]
  debugPrint(fmt.Sprintf("isub : %d - %d, stack : %+v", value1, value2, frame.stack))
  frame.pc += 1
}

func imul_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  value2 := frame.stack[len(frame.stack) - 1]
  value1 := frame.stack[len(frame.stack) - 2]
  frame.stack[len(frame.stack) - 2] = value1 * value2
  frame.stack = frame.stack[:len(frame.stack) - 1]
  debugPrint(fmt.Sprintf("imulb : %d * %d, stack : %+v", value1, value2, frame.stack))
  frame.pc += 1
}

func goto_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  offset := binary.BigEndian.Uint16(code[frame.pc + 1:])
  debugPrint(fmt.Sprintf("goto(%d)", offset))
  frame.pc += int(offset)
}

func if_icmpne_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  offset := binary.BigEndian.Uint16(code[frame.pc + 1:])
  value2 := frame.stack[len(frame.stack) - 1]
  value1 := frame.stack[len(frame.stack) - 2]
  frame.stack = frame.stack[:len(frame.stack) - 2]
  debugPrint(fmt.Sprintf("if_icmpne(0x%02x) : %d != %d, stack %+v", offset, value1, value2, frame.stack))
  if (value1 != value2) {
      frame.pc += int(offset)
  } else {
      frame.pc += 3
  }
}

func getstatic_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  index := binary.BigEndian.Uint16(code[frame.pc + 1:])
  frame.stack = append(frame.stack, uint32(index))
  debugPrint(fmt.Sprintf("getstatic(%d) : stack %+v", index, frame.stack))
  frame.pc += 3
}

func ldc_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  index := code[frame.pc + 1]
  frame.stack = append(frame.stack, uint32(index))
  debugPrint(fmt.Sprintf("ldc(%d) : stack %+v", index, frame.stack))
  frame.pc += 2
}

func invokestatic_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  index := binary.BigEndian.Uint16(code[frame.pc + 1:])
  self := m.parent
  className, methodName, descriptor := self.getConstant(index).getMethodRef()
  target := m.parent.classLoader.FindClass(className, vm)
  method := target.FindMethod(methodName, descriptor)
  retDesc, argsDesc := ParseDescriptor(descriptor)
  args := make([]uint32, len(argsDesc))
  for i := 0; i < len(argsDesc); i++ {
    args[i] = frame.stack[len(frame.stack) - i - 1]
  }
  frame.stack = frame.stack[:len(frame.stack) - len(args)]
  debugPrint(fmt.Sprintf("invokevirtual(%d): %s.%s:%s", index, className, methodName, descriptor))
  returnValue := vm.CallMethod(method, args)
  if retDesc != "V" {
    frame.stack = append(frame.stack, returnValue)
  }
  frame.pc += 3
}

func return_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  debugPrint(fmt.Sprintf("return"))
  frame.returnValue = 0
  frame.pc = len(code) // move to end of the code
}

func ireturn_Op(code []byte, m *Method, vm *VirtualMachine, frame *Frame) {
  frame.returnValue = frame.stack[len(frame.stack) - 1]
  frame.stack = frame.stack[:len(frame.stack) - 1]
  frame.pc = len(code) // move to end of the code
  debugPrint(fmt.Sprintf("ireturn : %d", frame.returnValue))
}
