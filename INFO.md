# How it works

When `sdees` starts for the first time it will request a `git` repository which can be [local or remote](https://github.com/schollz/sdees/blob/master/INFO.md#setting-up-git-server). Then `sdees` provides an option to write text using the editor of choice into a new *entry*. Each new entry is inserted into a new orphan branch in the supplied `git` repository. The benefit of placing each entry into its own orphan branch is that merge collisions are avoided after creating new entries on different machines. Thus, `sdees` makes it perfectly safe to make *new* entries without internet access.

The combination of entries be displayed as a *document*. The document is reconstructed by 1) fetching all remote branches, 2) filtering out which ones contain entries for the document of interest, 3) decrypting each entry (if needed), and then 4) sorting the entries by the commit date. Merge conflicts should only occur when simultaneous edits are made to the same entry (which is not the use case here, in general). Multiple documents can be stored in a single `git` repository.

Optionally, all information saved in the `git` repo can be encrypted using a symmetric cipher with a user-provided passphrase. The passphrase is not stored anywhere on the machine or repo. When enabled, each entry in the `git` repo is encrypted. When editing an encrypted document, a decrypted temp file is stored and then shredded (random bytes written and then deleted) after use.


# Setting up `git` server

You can just use a [Github](https://github.com/) or [Bitbucket](https://bitbucket.org/) repository and skip this.

You can also make your own personal `git` server by following these steps.

## Remote server

On the remote server `remote.com`, install the `git` server dependencies and make a new user:

```
sudo apt-get install git-core && \
  useradd git && \
  passwd git && \
  mkdir /home/git && \
  chown git:git /home/git
```

Then, from the client, add your key:

```
cat ~/.ssh/id_rsa.pub | ssh git@remote.com \
  "mkdir -p ~/.ssh && cat >>  ~/.ssh/authorized_keys"
```

and then to create a new repo, `newrepo.git`:

```
ssh git@remote.com "\
  mkdir -p newrepo.git && \
  git init --bare newrepo.git/ && \
  rm -rf clonetest && \
  git clone newrepo.git clonetest && \
  cd clonetest && \
  touch .new && \
  git add . && \
  git commit -m 'added master' && \
  git push origin master && cd .. && \
  rm -rf clonetest"
```
which you can add to `sdees` as `git@remote:newrepo.git`.

## Local server

If you want to just run locally, simply run:
```
cd /folder/to/put/repo
mkdir -p newrepo.git && \
  git init --bare newrepo.git/ && \
  rm -rf clonetest && \
  git clone newrepo.git clonetest && \
  cd clonetest && \
  touch .new && \
  git add . && \
  git commit -m 'added master' && \
  git push origin master && cd .. && \
  rm -rf clonetest
```
to make a new repo `newrepo.git`, and then add it to `sdees` as `/folder/to/put/repo/newrepo.git`.
