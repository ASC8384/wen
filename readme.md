# wen

一款基于递归下降解析器的编程语言，语法类似于 BrainFuck，命名为 `wen-ilovebuaa`。

## tokens

### 1. 空白字符

会忽略换行符（ \n ） 、回车符（ \r ） 、 横向制表符（ \ t ）、纵向制表符（ \v ）、换页符（ \f ）和空格符这 6 个空白字符 。

### 2. 操作符

|运算符| 作用|
|-|-|
|a| 将指针左移一个位置|
|i| 将指针右移一个位置|
|l| 将指针指向的单元的数值增加1|
|e| 将指针指向的单元的数值减少1|
|o| 输出指针指向的单元的数值（一般为ASCII字符）|
|v| 输入一个数值到指针指向的单元（一般为ASCII字符）|
|b| 如果指针指向的单元的数值为0，跳过匹配的 ]|
|u| 如果指针指向的单元的数值不为0，跳回到匹配的 [|
|+| 将指针处的字节与存储中的字节相加,结果存储在寄存器的字节中。|
|-| 将指针处的字节与存储中的字节相加,结果存储在寄存器的字节中。|
|*| 将指针处的字节与存储中的字节相乘,结果存储在寄存器的字节中。|
|/| 将指针处的字节与存储中的字节相除,结果存储在寄存器的字节中。|
|@| 修饰符存储处的数值不用转换ASCII码，以数字类型输入以及输出|
|A| 表示寄存器|
|?| 读取寄存器的值到指针处的字节 |
|!|存储寄存器的值到指针处的字节|
|$| 变量前缀|
|\[\]\(\)| 限制变量操作范围|
|&| 初始化数据|
|^| 指针指向变量的起始位置|
|B| 根据当前指针指向的值进行不同跳转|
|#| 条件语句结束|
|V| 跳转到条件语句结束|
|E| 程序立刻终止|

## Parser

使用递归下降解析器。

## Backend

使用 Go 语言解释执行。同时也提供了 C 语言后端，作为编译执行的可选项。

## Run

```bash
cd wen/
go build
wen sourcename.wen
```

或者

```bash
go run "wen\main.go" -c "wen\1.wen"
```

或者

```bash
go run main.go -d "llllllllbillllbillilllilllilaaaaeuililieiilbauaeuiioieeeollllllloollloiioaeoaollloeeeeeeoeeeeeeeeoiiloillo"
```

其中，`-c` 代表使用C语言后端，`-d` 代表输出抽象语法树。
