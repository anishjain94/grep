## Grep like CLI in go

#### Problem statement

Write a command line program that implements Unix command `grep` like functionality.

#### Video solution and approach

https://github.com/user-attachments/assets/61de835e-db48-450f-9026-4e1f9872c863

https://github.com/user-attachments/assets/37226aa3-dedb-46de-a29b-e73ba7076f5d


#### Features required and status

- [x] Ability to search for a string in a file

```
$  ./grep anish test_files/testfile.txt
```

```
this is anish.
is this anish.
this is anish?
anish
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

```
barbazfoo
food
```

- [x] Ability to perform case-insensitive grep using -i flag

```
$ ./grep -i Test test_files
```

```
test_files/testfile2.txt: this is test 8
test_files/testfile2.txt: this is test 14
```

- [x] Ability to write output to a file instead of a standard out.

```
$ ./grep -o temp.txt anish test_files
```

should create an temp.txt file with the output from `grep`. for example,

```
$ cat out.txt
test_files/parent_dir1/child_dir1/child_dir1_file.txt: this is anish parent_dir1/child_dir1/child_dir1_file.txt
test_files/parent_dir1/child_dir1/child_dir1_file.txt: is this anish parent_dir1/child_dir1/child_dir1_file.txt
test_files/parent_dir1/child_dir1/child_dir1_file.txt: this is anish? parent_dir1/child_dir1/child_dir1_file.txt
test_files/parent_dir1/child_dir2/child_dir2_file.txt: this is anish parent_dir1/child_dir2/child_dir2_file.txt
test_files/parent_dir1/child_dir2/child_dir2_file.txt: is this anish parent_dir1/child_dir2/child_dir2_file.txt
test_files/parent_dir1/child_dir2/child_dir2_file.txt: this is anish? parent_dir1/child_dir2/child_dir2_file.txt
test_files/parent_dir2/parent_dir2_file1.txt: this is anish parent_dir2/parent_dir2_file1.txt
test_files/parent_dir2/parent_dir2_file1.txt: is this anish parent_dir2/parent_dir2_file1.txt
test_files/parent_dir2/parent_dir2_file1.txt: this is anish? parent_dir2/parent_dir2_file1.txt
test_files/testfile.txt: this is anish.
test_files/testfile.txt: is this anish.
test_files/testfile.txt: this is anish?
test_files/testfile.txt: anish
```

- [x] Ability to search for a string recursively in any of the files in a given directory. When searching in multiple files, the output should indicate the file name and all the output from one file should be grouped together in the final output. (in other words, output from two files shouldn't be interleaved in the final output being printed)

```
$ ./grep test test_files
test_files/testfile2.txt: this is test 8
test_files/testfile2.txt: this is test 14
```

- [x] Ability to print n lines before using `-b` flag

```
$ ./grep -b 2 test test_files
```

```
test_files/testfile2.txt: this is line 6
test_files/testfile2.txt: this is line 7
test_files/testfile2.txt: this is test 8
test_files/testfile2.txt: this is line 12
test_files/testfile2.txt: this is line 13
test_files/testfile2.txt: this is test 14
```

- [x] Ability to print n lines after using `-a` flag

```
$ ./grep -a 2 test test_files
```

```
test_files/testfile2.txt: this is test 8
test_files/testfile2.txt: this is line 9
test_files/testfile2.txt: this is line 10
test_files/testfile2.txt: this is test 14
test_files/testfile2.txt: this is line 15
test_files/testfile2.txt: this is line 16
```

- [x] Ability to print count of matches using `-c` flag

```
$ ./grep -c anish test_files/parent_dir1/child_dir1/child_dir1_file.txt
```

```
3
```

### Future Todo's

- Handle for condition when file limit opening is restricted by os. make is os independent.
- Load test and benchmarking grep
- Add other options from `GREP`
