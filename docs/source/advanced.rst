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

.. _add an Issue: https://github.com/schollz/gojot/issues/new

How to setup ``git`` server
----------------------------



The usage of *gojot* requires having a ``git`` repository hosted somewhere. You can have these repositories hosted locally or remotely, through a 3rd party service, or on your own server. Here's some examples on how to do either.

*Easiest way:* Just dont. You can just use a `Github`_ / `Bitbucket`_ / `Gitlab`_
repository and skip this.

*Alternatively:* You can also make your own personal ``git`` server by
following these steps.

Local server
~~~~~~~~~~~~~~~~~~

If you want to just run locally, simply run:

::

    git init --bare /folder/to/newrepo.git

to make a new repo ``newrepo.git``, and then add it to ``gojot`` as
``/folder/to/newrepo.git``.


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
      git init --bare newrepo.git

which you can add to ``gojot`` as ``git@remote.com:newrepo.git``.

.. _Github: https://github.com/
.. _Bitbucket: https://bitbucket.org/
.. _Gitlab: https://gitlab.com/users/sign_in
.. _micro: https://github.com/zyedidia/micro
.. _vim: http://www.vim.org/download.php
.. _nano: https://www.nano-editor.org/
.. _emacs: https://www.gnu.org/software/emacs/
