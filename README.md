Trie
============

Golang implementation of the Double-Array trie.

Usage
-----

###Sample: Build a trie form the keyword list.

####Code
```
 import (
        "github.com/ikawaha/trie"
 
        "fmt"
        "os"
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
      t := trie.NewDoubleArrayTrieFromKeywords(keywords)
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
 奈良
 奈良先端
 奈良先端科学技術大学
 奈良奈良先端科学技術大学院大学
```
#### Code
```
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
      t, err = trie.NewDoubleArrayTrieFromFile(file)
      if err != nil {
          panic(err)
      }
      fmt.Println(t.CommonPrefixSearch("奈良先端科学技術大学院大学"))
 }
```

####Result
```
 [奈良 奈良先端 奈良先端科学技術大学]
```

Copyright and license
---------------------

Copyright (c) 2014 Ikuo Kawaharada All Rights Reserved.

This software is released under the MIT License.
See LICENSE.txt
