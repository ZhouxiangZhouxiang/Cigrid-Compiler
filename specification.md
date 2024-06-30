# 一个类C编译器的说明

## 1 语法

### 1.1 基本语法

```
 unop->"!"|"-"|"*"|"&"
binop->"+"|"-"|"*"|"/"|"<"|">"|"<="|">="|"=="|"!="|"&&"|"||"
   ty->"void"|"int"|"string"|ty "*"
expr->Ident|UInt|String
     |expr binop expr
     |unop expr
     |Ident"("[expr {","expr}]")"
     |Ident"["expr"]"{"["expr"]"}
     |"("expr")"
stmt->ty Ident ";" // int i;
     |ty Ident "=" expr ";" // int i = 1;
     |Ident "=" expr";" // i = 1;
     |ty Ident"["UInt"]" { "["UInt"]" }";" // int array[3][3];
     |ty Ident"["UInt"]" { "["UInt"]" } "=" "{"..."}"" ";" 
      // int array[3][3] = {{1,2,3},{4,5,6},{7,8,9}};
     |Ident"["expr"]" { "["expr"]" } "=" expr ";" // array[0][0] = 2;
     |Ident"("[expr{","expr}]")"";" // functioncall(a, b);
     |"if""("expr")""{"stmt"}"["else""{"stmt"}"]
     |"while""("expr")""{"stmt"}"
     |"return"[expr]";"
global->ty Ident "=" expr ";"
       |ty Ident"["UInt"]" { "["UInt"]" } "=" "{"..."}"" ";"
       |ty Ident"(" [{ty Ident {"," ty Ident}] ")" "{" stmt "}"
program->{global}
```

### 1.2 关键词

```c++
void string int if else while return
```

$$
a_{i} = \alpha^{ab} \times v
$$

| s<sub>c</sub> | v<sup>d</sup> | $aa_b$ |
| ------------- | ------------- | ------ |
|               |               |        |
|               |               |        |
|               |               |        |
