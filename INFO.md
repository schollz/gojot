# How-to setup git server

You can definetly use github.com or bitbucket.org or similar. If you want to make a personal one you can also do it locally or remotely.

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
