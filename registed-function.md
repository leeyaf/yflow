# 可用的 Function

## Cell

从矩阵中读值

```` yaml
# 矩阵名, 行下标, 列下标
(( Cell("prev.out", 0, 1) ))
````

## Row

从矩阵中读取一行

```` yaml
# 矩阵名, 行下标
# return string
(( Row("prev.out", 0) ))
````

## Col

从矩阵中读取一列

```` yaml
# 矩阵名, 列下标
# return []string
(( Col("prev.out", 0) ))
````

## ParseInt

字符串转 int64

```` yaml
# 字符串
# return int64
(( ParseInt("1234") ))
````

## ParseFloat

字符串转 float64

```` yaml
# 字符串
# return float64
(( ParseFloat("1234") ))
````

## Join

字符串数组转字符串

元素的分隔符与包裹符的类型是 ``rune``，因为符号转换问题，这里使用十进制数字进行填充，虽然降低了可读性，但保护了灵活性

|常用符号|ASCII十进制|
|---|----------|
|"|34|
|'|39|
|,|44|
|-|45|
|.|46|
|:|58|
|;|59|
|`|96|

包裹符的默认值为 0，视为无包裹符

```` yaml
# 字符串数组，分隔符，元素包裹符
# return string
(( Join([]string{"a","b","c"}, 44, 0) ))
````

## Split

字符串分隔为字符串数组

```` yaml
# 字符串，分隔符
# return []string
(( Split("a, b, c", 44) ))
````

## TimeConvert

变化时间字符串的显示

```` yaml
# 时间字符串, 原格式化模版, 新格式化模版
# return string
(( TimeConvert("2025-06-04 10:46:47 +0800 CST", "", "2006-01-02 15:04") ))
````