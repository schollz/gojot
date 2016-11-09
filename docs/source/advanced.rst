Advanced
=========

Contributing
-----------------

Contributions are welcome!

Please feel free to `add an Issue`_ if you find a bug or have a feature
request.

Pull requests are welcome. If you’d like to test your changes, make sure
that you change ``GITHUB_TEST_REPO`` in ``src/git_test.go`` to a Git
repo that you have that can be used for testing. Then simply use

::

    make update
    make test

to perform tests.

.. _add an Issue: https://github.com/schollz/sdees/issues/new

How to setup ``git`` server
----------------------------



The usage of *sdees* requires having a ``git`` repository hosted somewhere. You can have these repositories hosted locally or remotely, through a 3rd party service, or on your own server. Here's some examples on how to do either.

*Easiest way:* Just dont. You can just use a `Github`_ or `Bitbucket`_
repository and skip this.

*Alternatively:* You can also make your own personal ``git`` server by
following these steps.

Remote server
~~~~~~~~~~~~~~~~~~

On the remote server ``remote.com``, install the ``git`` server
dependencies and make a new user:

::

    sudo apt-get install git-core && \
      useradd git && \
      passwd git && \
      mkdir /home/git && \
      chown git:git /home/git

The log into the client and add your key to the server:

::

    cat ~/.ssh/id_rsa.pub | ssh git@remote.com \
      "mkdir -p ~/.ssh && cat >>  ~/.ssh/authorized_keys"

Then, still on the client, create a new repo, ``newrepo.git``, on the
server:

::

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

which you can add to ``sdees`` as ``git@remote:newrepo.git``.

.. _Github: https://github.com/
.. _Bitbucket: https://bitbucket.org/


Local server
~~~~~~~~~~~~~~~~~~

If you want to just run locally, simply run:

::

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

to make a new repo ``newrepo.git``, and then add it to ``sdees`` as
``/folder/to/put/repo/newrepo.git``.
