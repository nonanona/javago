package vm

import (
)

type ClassLoader struct {
  loadedClasses map[string]ClassFile
}

func NewClassLoader() ClassLoader {
  cl := ClassLoader{}
  cl.loadedClasses = make(map[string]ClassFile)
  return cl
}

func (cl ClassLoader) FindClass(className string, vm *VirtualMachine) ClassFile {
  loaded, ok := cl.loadedClasses[className]
  if ok {
    return loaded
  }
  panic("Loading class is not yet implemented")
}

func (cl ClassLoader) LoadClassByFile(path string, vm *VirtualMachine) ClassFile {
  cf := ParseClassFile(path, &cl)
  cl.loadedClasses[cf.getThisClassName()] = cf
  vm.CallStaticInitializer(cf)
  return cf
}
