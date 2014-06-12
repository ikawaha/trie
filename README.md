Trie
============

Golang implementation of the Double-Array trie.

Usage
-----

###Sample: Build a trie form the keyword list.

####Code
```
package main

import (
        "github.com/ikawaha/trie"

        "fmt"
 )

 func main() {
      keywords :=[]string{
                 "hello",
                 "world",
                 "関西",
                 "国際",
                 "国際空港",
                 "関西国際空港",
      }
      t, err := trie.NewDoubleArrayTrie(keywords)
      if err != nil {
         panic(err)
      }
      fmt.Println(t.Search("hello"))
      fmt.Println(t.Search("world"))
      fmt.Println(t.Search("goodby"))
      fmt.Println(t.CommonPrefixSearch("関西国際空港"))
}
```

####Result
```
true
true
false
[関西 関西国際空港]
```

###Sample: Build a trie from the file.

#### Input file (keyword_list.txt)
```
逓信大
電気通信大学
東京電気大学
電通大
電気通信大学大学院
電気通信大学大学院大学
情報工学科
```
#### Code
```
package main

import (
        "github.com/ikawaha/trie"

        "fmt"
        "os"
 )

 func main() {
      file, err := os.Open("keyword_list.txt")
      if err != nil {
         panic(err)
      }
      defer file.Close()
      t, err := trie.NewDoubleArrayTrie(file)
      if err != nil {
          panic(err)
      }
      fmt.Println(t.CommonPrefixSearch("電気通信大学大学院大学"))
}
```

####Result
```
[電気通信大学 電気通信大学大学院] [3 4]
```

Copyright and license
---------------------

Copyright (c) 2014 ikawaha Rights Reserved.

This software is released under the MIT License.
See LICENSE.txt
