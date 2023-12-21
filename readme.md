# wen

一款基于递归下降解析器的编程语言，语法类似于 BrainFuck，命名为 `wen-ilovebuaa`。

## tokens

### 1. 空白字符

会忽略换行符（ \n ） 、回车符（ \r ） 、 横向制表符（ \ t ）、纵向制表符（ \v ）、换页符（ \f ）和空格符这 6 个空白字符 。

### 2

## Lexer

## Parser

## AST

## Backend

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
