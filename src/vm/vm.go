package vm

import (
)

type VirtualMachine struct {
  classLoaders []ClassLoader
}

func NewVM() VirtualMachine {
  vm := VirtualMachine{}
  vm.classLoaders = make([]ClassLoader, 1)
  // Bootstrap class loader
  vm.classLoaders[0] = NewClassLoader()
  return vm
}

func (vm VirtualMachine) CallMethod(m *Method, args []uint32) uint32 {
  return ProcessInstruction(m, args, &vm)
}

func (vm VirtualMachine) CallStaticInitializer(cf ClassFile) {
  method := cf.FindMethod("<clinit>", "()V")
  if method != nil {
    vm.CallMethod(method, make([]uint32, 0))
  }
}

func (vm VirtualMachine) Execute(classFilePath string) {
  cf := vm.classLoaders[0].LoadClassByFile(classFilePath, &vm)
  method := cf.FindMethod("main", "([Ljava/lang/String;)V")
  vm.CallMethod(method, make([]uint32, 0))
}
