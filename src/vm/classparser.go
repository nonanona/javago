package vm

import (
  "io/ioutil"
  "encoding/binary"
  "bytes"
)

type Constant struct {
  parent* ClassFile
  tag uint8
  data []byte
}

func (c Constant) getClassName() string {
  if c.tag != 7 {
    panic("This constant is not a Class info")
  }
  return c.parent.getConstant(binary.BigEndian.Uint16(c.data)).getString()
}

func (c Constant) getString() string {
  if c.tag != 1 {
    panic("This constant is not a UTF-8 info")
  }
  return string(c.data)
}

func (c Constant) getMethodRef() (string, string, string) {  // class, name, type
  if c.tag != 10 {
    panic("This constant is not a Method ref")
  }
  cn := c.parent.getConstant(binary.BigEndian.Uint16(c.data)).getClassName()
  n, t := c.parent.getConstant(binary.BigEndian.Uint16(c.data[2:])).getNameAndType()
  return cn, n, t
}

func (c Constant) getNameAndType() (string, string) {
  if c.tag != 12 {
    panic("This constant is not a NameAndType")
  }
  n := c.parent.getConstant(binary.BigEndian.Uint16(c.data)).getString()
  t := c.parent.getConstant(binary.BigEndian.Uint16(c.data[2:])).getString()
  return n, t
}

type Field struct {
  parent* ClassFile
  accessFlag uint16
  nameIndex uint16
  descriptorIndex uint16
  attributes []Attribute
}

type Method struct {
  parent* ClassFile
  accessFlag uint16
  nameIndex uint16
  descriptorIndex uint16
  attributes []Attribute
}

func (m Method) getMethodName() string {
  return m.parent.constantPool[m.nameIndex - 1].getString()
}

func (m Method) getDescriptor() string {
  return m.parent.constantPool[m.descriptorIndex - 1].getString()
}

func (m Method) getAttribute(tag string) Attribute {
  for _, attrib := range m.attributes {
    if m.parent.getString(attrib.attributeNameIndex) == tag {
      return attrib
    }
  }
  panic("Attribute for " + tag + " not found")
}

func (m Method) GetCode() CodeAttribute {
  return m.getAttribute("Code").toCodeAttribute()
}

type Attribute struct {
  parent* ClassFile
  attributeNameIndex uint16
  data []byte
}

type ExceptionTable struct {
  startPc uint16
  endPc uint16
  handlerPc uint16
  catchType uint16
}

type CodeAttribute struct {
  maxStack uint16
  maxLocals uint16
  code []byte
  exceptionTables []ExceptionTable
  attributes []Attribute
}

func (a Attribute) toCodeAttribute() CodeAttribute {
  if a.parent.getString(a.attributeNameIndex) != "Code" {
    panic("This attribute is not a Code attribute")
  }
  ca := CodeAttribute{}
  buf := bytes.NewBuffer(a.data)
  binary.Read(buf, binary.BigEndian, &ca.maxStack)
  binary.Read(buf, binary.BigEndian, &ca.maxLocals)

  var codeLength uint32
  binary.Read(buf, binary.BigEndian, &codeLength)
  ca.code = make([]byte, codeLength)
  binary.Read(buf, binary.BigEndian, &ca.code)

  var exLength uint16
  binary.Read(buf, binary.BigEndian, &exLength)
  ca.exceptionTables = make([]ExceptionTable, exLength)
  binary.Read(buf, binary.BigEndian, &ca.exceptionTables)

  var attribLength uint16
  binary.Read(buf, binary.BigEndian, &attribLength)
  ca.attributes = make([]Attribute, attribLength)
  for i := uint16(0); i < attribLength; i++ {
    ca.attributes[i] = readAttribute(buf, a.parent)
  }

  return ca
}

type ClassFile struct {
  classLoader *ClassLoader
  magic uint32
  minorVersion uint16
  majorVersion uint16
  constantPool []Constant
  accessFlag uint16
  thisClass uint16
  superClass uint16
  interfaces []uint16
  fields []Field
  methods []Method
  attributes []Attribute
}

func (cf ClassFile) getConstant(idx uint16) Constant {
  return cf.constantPool[idx - 1]
}

func (cf ClassFile) getThisClassName() string {
  return cf.getConstant(cf.thisClass).getClassName();
}

func (cf ClassFile) getString(idx uint16) string {
  return cf.getConstant(idx).getString()
}

func (cf ClassFile) FindMethod(name string, descriptor string) *Method {
  for _, method := range cf.methods {
    if method.getMethodName() == name && method.getDescriptor() == descriptor {
      return &method
    }
  }
  return nil
}

func readAttribute(buffer* bytes.Buffer, cf *ClassFile) Attribute {
  a := Attribute{parent: cf}
  binary.Read(buffer, binary.BigEndian, &a.attributeNameIndex)
  var length uint32
  binary.Read(buffer, binary.BigEndian, &length)
  a.data = make([]byte, length)
  binary.Read(buffer, binary.BigEndian, &a.data)
  return a
}

func readField(buffer* bytes.Buffer, cf *ClassFile) Field {
  f := Field{parent: cf}
  binary.Read(buffer, binary.BigEndian, &f.accessFlag)
  binary.Read(buffer, binary.BigEndian, &f.nameIndex)
  binary.Read(buffer, binary.BigEndian, &f.descriptorIndex)
  var attribCount uint16
  binary.Read(buffer, binary.BigEndian, &attribCount)
  f.attributes = make([]Attribute, attribCount)
  for i := uint16(0); i < attribCount; i++ {
    f.attributes[i] = readAttribute(buffer, cf)
  }
  return f
}

func readMethod(buffer* bytes.Buffer, cf *ClassFile) Method {
  m := Method{parent: cf}
  binary.Read(buffer, binary.BigEndian, &m.accessFlag)
  binary.Read(buffer, binary.BigEndian, &m.nameIndex)
  binary.Read(buffer, binary.BigEndian, &m.descriptorIndex)
  var attribCount uint16
  binary.Read(buffer, binary.BigEndian, &attribCount)
  m.attributes = make([]Attribute, attribCount)
  for i := uint16(0); i < attribCount; i++ {
    m.attributes[i] = readAttribute(buffer, cf)
  }
  return m
}

func readConstant(buffer* bytes.Buffer, cf *ClassFile) Constant {
  var tag uint8
  var length int
  binary.Read(buffer, binary.BigEndian, &tag)
  switch tag {
    case 1:  // Utf8
      var strLength uint16
      binary.Read(buffer, binary.BigEndian, &strLength)
      length = int(strLength)
    case 3: length = 4 // Integer
    case 4: length = 4 // Float
    case 5: length = 8 // Long
    case 6: length = 8 // Double
    case 7: length = 2 // Class
    case 8: length = 2 // String
    case 9: length = 4 // Filedref
    case 10: length = 4 // Methodref
    case 11: length = 4 // InterfaceMethodref
    case 12: length = 4 // NameAndType
    case 15: length = 3 // MethodHandle
    case 16: length = 2 // MethodType
    case 18: length = 4 // InvokeDynamic
    default:
      panic("Unknown constant pool type:" + string(tag))
  }
  b := make([]byte, length)
  binary.Read(buffer, binary.BigEndian, &b)
  return Constant{cf, tag, b}
}

// Parse the class file
func ParseClassFile(path string, cl *ClassLoader) ClassFile {
  dat, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err)
  }
  classFile := ClassFile{classLoader: cl}

  buf := bytes.NewBuffer(dat)
  binary.Read(buf, binary.BigEndian, &classFile.magic)
  if classFile.magic != 0xCAFEBABE {
    panic("Unknown magic number")
  }
  binary.Read(buf, binary.BigEndian, &classFile.minorVersion)
  binary.Read(buf, binary.BigEndian, &classFile.majorVersion)

  // Read constant pools
  var constantPoolCount uint16
  binary.Read(buf, binary.BigEndian, &constantPoolCount)
  constantPoolCount--;
  classFile.constantPool = make([]Constant, constantPoolCount)
  for i := uint16(0); i < constantPoolCount; i++ {
    classFile.constantPool[i] = readConstant(buf, &classFile)
  }

  binary.Read(buf, binary.BigEndian, &classFile.accessFlag);
  binary.Read(buf, binary.BigEndian, &classFile.thisClass);
  binary.Read(buf, binary.BigEndian, &classFile.superClass);

  // Read interfaces
  var interfaceCount uint16
  binary.Read(buf, binary.BigEndian, &interfaceCount)
  classFile.interfaces = make([]uint16, interfaceCount)
  binary.Read(buf, binary.BigEndian, &classFile.interfaces)

  // Read fields
  var fieldCount uint16
  binary.Read(buf, binary.BigEndian, &fieldCount)
  classFile.fields = make([]Field, fieldCount)
  for i := uint16(0); i < fieldCount; i++ {
    classFile.fields[i] = readField(buf, &classFile)
  }

  // Read Methods
  var methodCount uint16
  binary.Read(buf, binary.BigEndian, &methodCount)
  classFile.methods = make([]Method, methodCount)
  for i := uint16(0); i < methodCount; i++ {
    classFile.methods[i] = readMethod(buf, &classFile)
  }

  // Read Attributes
  var attribCount uint16
  binary.Read(buf, binary.BigEndian, &attribCount)
  classFile.attributes = make([]Attribute, attribCount)
  for i := uint16(0); i < attribCount; i++ {
    classFile.attributes[i] = readAttribute(buf, &classFile)
  }

  return classFile
}
