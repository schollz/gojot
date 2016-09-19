
# Notes

# master branch

The `master` branch only contains README.md about what the repo is?

## Reading entries

Here, we'd like to pull particular entry without checking out (to avoid lots of file writing).

`git show branch_name:file` gets the current text of the `file` in `branch_name`
`git ls-tree --name-only branch_name` returns the names document in that branch
`git log --pretty=format:"%H -=- %ad -=- %s" branch_name` returns hash, date, and the subject for the latest in that branch


## Creating new entry

After document is written:

```
git checkout --orphan 'randombranchid'
git add document_a
git commit -am "GPG-encoded text of first line in document_a"
```

And then the commit can be pushed.

## Pulling latest

Things can be pulled with
```
# get list of branches on remote
git branch -r  
# track each branch
git branch --track 'branch_name1' 'origin/branch_name1'
git branch --track 'branch_name2' 'origin/branch_name2'
# fetch/pull all of them
git fetch --all
```

## Cache

Cache should look like:

```json
{
  "document_a.txt": {
    "branch_1": {
      "fulltext":"some text [encrytped]",
      "latest_commit": "asbaklsdjclasdifasdf"
    },
    "branch_2": {
      "fulltext":"some different text [encrytped]",
      "latest_commit": "asdfaskcoeckoekasec"
    },
  },
  "document_b.txt": {
    "branch_2": {
      "fulltext":"some text [encrytped]",
      "latest_commit": "asbaklsdjclasdifasdf"
    },
    "branch_3": {
      "fulltext":"some text [encrytped]",
      "latest_commit": "asbaklsdjclasdifasdf"
    },
  },
}
```

This cache can then keep track of all branches (all unique keys of the set of all keys in all documents).

Cache is updated on two conditions:

1. The branches listed by `git branch --list` contain branches that do not exist in cache. In this case, that branch is investigated to see which documents there are, and then adds them to the cache respectively.
2. 


## Deletion

Things can be deleted using `git branch -D <branch_name> && git push origin --delete <branch_name>`

## Stitching document

The whole document is recapitulated by looping over every branch and finding any file that matches and concatenating them in a file using `git show branch:document`:
```
>>> 45b5d58f89d415589a08c4c7f545f2804ef19aee Mon Sep 19 07:18:49 2016 -0400

Some text about the entry

>>> ced7eb3425b9bf409d47b4394eb75ba4605d6e75 Tue Sep 19 07:18:49 2016 -0400

Some other entry
```
Like original SDEES, each entry will be checked for changes, and if there is a change, then it will checkout that branch and commit a new encrypted entry.

## Config file

```json
{
  "editor":"vim",
  "remote":"github.com/somename/somerepo"
}
```

# Quetions
- Can git.exe be bundled for windows?
- libgit vs. command line. How fast is it to do command line stuff?
