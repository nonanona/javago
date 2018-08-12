# javago
java with golang

$ GOPATH=$PWD go run src/main.go src/Factorial.class 

bipush(10) : stack [10]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [10]  
iconst_1 : stack [10 1]  
if_icmpne(0x07) : 10 != 1, stack []  
iload_0 : stack [10]  
iload_0 : stack [10 10]  
iconst_1 : stack [10 10 1]  
isub : 10 - 1, stack : [10 9]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [9]  
iconst_1 : stack [9 1]  
if_icmpne(0x07) : 9 != 1, stack []  
iload_0 : stack [9]  
iload_0 : stack [9 9]  
iconst_1 : stack [9 9 1]  
isub : 9 - 1, stack : [9 8]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [8]  
iconst_1 : stack [8 1]  
if_icmpne(0x07) : 8 != 1, stack []  
iload_0 : stack [8]  
iload_0 : stack [8 8]  
iconst_1 : stack [8 8 1]  
isub : 8 - 1, stack : [8 7]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [7]  
iconst_1 : stack [7 1]  
if_icmpne(0x07) : 7 != 1, stack []  
iload_0 : stack [7]  
iload_0 : stack [7 7]  
iconst_1 : stack [7 7 1]  
isub : 7 - 1, stack : [7 6]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [6]  
iconst_1 : stack [6 1]  
if_icmpne(0x07) : 6 != 1, stack []  
iload_0 : stack [6]  
iload_0 : stack [6 6]  
iconst_1 : stack [6 6 1]  
isub : 6 - 1, stack : [6 5]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [5]  
iconst_1 : stack [5 1]  
if_icmpne(0x07) : 5 != 1, stack []  
iload_0 : stack [5]  
iload_0 : stack [5 5]  
iconst_1 : stack [5 5 1]  
isub : 5 - 1, stack : [5 4]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [4]  
iconst_1 : stack [4 1]  
if_icmpne(0x07) : 4 != 1, stack []  
iload_0 : stack [4]  
iload_0 : stack [4 4]  
iconst_1 : stack [4 4 1]  
isub : 4 - 1, stack : [4 3]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [3]  
iconst_1 : stack [3 1]  
if_icmpne(0x07) : 3 != 1, stack []  
iload_0 : stack [3]  
iload_0 : stack [3 3]  
iconst_1 : stack [3 3 1]  
isub : 3 - 1, stack : [3 2]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [2]  
iconst_1 : stack [2 1]  
if_icmpne(0x07) : 2 != 1, stack []  
iload_0 : stack [2]  
iload_0 : stack [2 2]  
iconst_1 : stack [2 2 1]  
isub : 2 - 1, stack : [2 1]  
invokevirtual(2): Factorial.factorial:(I)I  
iload_0 : stack [1]  
iconst_1 : stack [1 1]  
if_icmpne(0x07) : 1 != 1, stack []  
iconst_1 : stack [1]  
goto(11)  
ireturn : 1  
imulb : 2 * 1, stack : [2]  
ireturn : 2  
imulb : 3 * 2, stack : [6]  
ireturn : 6  
imulb : 4 * 6, stack : [24]  
ireturn : 24  
imulb : 5 * 24, stack : [120]  
ireturn : 120  
imulb : 6 * 120, stack : [720]  
ireturn : 720  
imulb : 7 * 720, stack : [5040]  
ireturn : 5040  
imulb : 8 * 5040, stack : [40320]  
ireturn : 40320  
imulb : 9 * 40320, stack : [362880]  
ireturn : 362880  
imulb : 10 * 362880, stack : [3628800]  
ireturn : 3628800  
pop()  
pop : stack []  
return  
