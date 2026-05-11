# 可用的 Step

## Test

输出等于输入，可以用来测试动态表达式

```` yaml
- step: Test
  in:
  - (( Split(input["configIds"]) ))
  - (( ParseInt(input["minLevel"])*2 ))
````

#### MatrixNew

从输入创建矩阵

```` yaml
- step: MatrixNew
  in:
	- ";" # 分隔符
	- a; b; c; d # 第一行数据
	- 1; 2; 3; 4 # 第二行数据
````

#### MatrixTranspose

矩阵转置，行变成列，列变成行

```` yaml
- step: MatrixTranspose
  in:
	- prev.out # 矩阵名，内置 prev 指向同 workflow 中的上一个 step
````

#### MatrixInsert

插入一行或一列到矩阵

```` yaml
- step: MatrixInsert
  in:
	- youStep.in # 矩阵名，定义了 name 的 step 的输入矩阵，可以是任意 workflow 中的 step
	- row # 行（row）还是列（col）
	- 0 # 第一行
	- a # 第一行第一列的数据
	- c # 第一行第二列的数据
	- 1 # 第一行第三列的数据
````

#### MatrixSelect

选择矩阵里的部分行或列

```` yaml
- step: MatrixSelect
  in:
	- prev.out # 矩阵名
	- col # 行（row）还是列（col）
	- 0 # 第一列
	- 2 # 第三列
	- 5 # 第六列
````

#### MatrixVLookUp

垂直查找并替换矩阵的一列

```` yaml
- step: MatrixVLookUp
  in:
	- prev.out # 源矩阵名
	- 1 # 源矩阵的第二列
	- youStep.out # 查找替换的目标矩阵名
	- 0 # 匹配目标矩阵的第一列
	- 1 # 替换为目标矩阵的第二列
````

#### MatrixApply

对矩阵的行或列应用动态表达式

``value`` 代表矩阵元素的值

**注意**：这里的动态表达式不要使用 ``(( ))`` 包裹

```` yaml
- step: MatrixApply
  in:
	- prev.out # 矩阵名
	- col # 行（row）还是列（col）
	- 1 # 矩阵的第二列
	- TimeConvert(value, "" , "2006-01-02 15:04") # 动态表达式
````

#### MysqlExecute

执行 MySQL 查询

```` yaml
- step: MysqlExecute
  in:
	- env["mysqlUri"] # 连接字符串
	- select count(id) from hero # SQL
````