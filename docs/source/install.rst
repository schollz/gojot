Getting started
==================

Before you begin
----------------

Before you install *jot* make sure you have the following requirement:

* **git** (version 2.5+)


If you are on Windows/OS X just `download the latest version of git`_.

If you are on Ubuntu, you can install the latest from apt:

::

    add-apt-repository ppa:git-core/ppa -y
    apt-get update
    apt-get install git -y

.. _download the latest version of git: https://git-scm.com/downloads
.. _git-scm.com: https://git-scm.com/downloads

Install
-----------------

.. include:: ./downloads.rst

Alternative downloads
~~~~~~~~~~~~~~~~~~~~~~~

Previous versions of *jot* and also downloads for FreeBSD, OpenBSD, and 32-bit versions can be found on the `Github releases <https://github.com/schollz/jot/releases/latest>`_.

Install from source
~~~~~~~~~~~~~~~~~~~~~

First `install Go`_ on your system. Then just use

::

    go get -u github.com/schollz/jot

.. _the latest release binary: https://github.com/schollz/jot/releases/latest
.. _install Go: https://golang.org/dl/
.. _version 2.5.0, released March, 2016: https://git-scm.com/docs/git-for-each-ref/2.5.0

Run *jot*
-------------

When you first run *jot* you will be prompted for a ``git`` repository. If you have an account
on `Github`_, `Bitbucket`_, `Gitlab`_, you can make a new repo and use that (make sure you have exchanged SSH keys). Alternatively you
can make a local repository using

::

    git init --bare /folder/to/newrepo.git

which you can add to ``jot`` as
``/folder/to/newrepo.git``. You can also `host your own repo quite easily`_ if you'd like (totally optional).


.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciicast-92591.json" async autoplay="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>

That's it! If you'd like to see help, simply use ``--help`` or check out `usage guide`_ for more information about what *jot* can do.

Problems?
----------

If you have any problems at all, please `submit an Issue`_ and someone can help you sort it out.

.. _submit an Issue: https://github.com/schollz/jot/issues/new
.. _usage guide: /examples.html
.. _host your own repo quite easily: /advanced.html#how-to-setup-git-server
.. _Gitlab: https://gitlab.com/users/sign_in
.. _Bitbucket: https://bitbucket.org/account/signin/
.. _Github: https://github.com/
