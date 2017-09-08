gmend
=====

gmend is cli tool helps to amend history commit.

### How to get it

```console
$ go get github.com/zhouhaibing089/gmend
```

### How to use it

```console
$ cd path/to/git/repository
$ gmend <commit>
```

### How it works

1.  create a backup branch in case of any failure.
1.  save all the commits `<commit>..HEAD`.
1.  reset `HEAD` to `<commit>`.
1.  make your changes and press Enter.
1.  save you commit via `git commit -a --amend --no-edit`.
1.  apply all the saved commits via `git cherry-pick`.