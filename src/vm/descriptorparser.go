package vm

func ParseDescriptor(desc string) (string, []string) {
  i := 0
  if desc[i] != '(' {
    panic("Invalid descriptor: " + desc)
  }
  i++;

  args := make([]string, 0)

  for desc[i] != ')' {
    switch (desc[i]) {
    case 'I':
      args = append(args, "I")
    default:
      panic("Unimplemented arg descriptor symbol: " + desc)
    }
    i++;
  }

  if desc[i] != ')' {
    panic("Invlaid descriptor " + desc)
  }
  i++;

  switch (desc[i]) {
  case 'I':
    return "I", args
  default:
    panic("Unimpelmented return descriptor symbol: " + desc)
  }
}
