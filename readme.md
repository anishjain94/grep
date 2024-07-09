## Grep like CLI in go

#### Problem statement

Write a command line program that implements Unix command `grep` like functionality.

#### Features required and status

- [x] Ability to search for a string in a file

```
$ ./grep "search_string" filename.txt
I found the search_string in the file.
```

- [x] Ability to search for a string from standard input

```
$ ./grep foo
bar
barbazfoo
Foobar
food
^D
```

output -

```
barbazfoo
food
```

- [x] Ability to write output to a file instead of a standard out.

```
$ ./grep lorem loreipsum.txt -o out.txt
```

should create an out.txt file with the output from `grep`. for example,

```
$ cat out.txt
lorem ipsum
a dummy text usually contains lorem ipsum
```

- [-] Ability to search for a string recursively in any of the files in a given directory. When searching in multiple files, the output should indicate the file name and all the output from one file should be grouped together in the final output. (in other words, output from two files shouldn't be interleaved in the final output being printed)

```
$ ./grep "test" tests
tests/test1.txt:this is a test file
tests/test1.txt:one can test a program by running test cases
tests/inner/test2.txt:this file contains a test line
```

- [-] Package test
