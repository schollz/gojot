Installation Guide
==================

Requirements
----------------

Before you install *sdees* make sure you have the following requirement:

* **git** (version 2.5+)


If you are on Windows/OS X just `download the latest version of git`_.

If you are on Ubuntu, you can install the latest from apt:

::

    add-apt-repository ppa:git-core/ppa -y
    apt-get update
    apt-get install git -y

.. _download the latest version of git: https://git-scm.com/downloads
.. _git-scm.com: https://git-scm.com/downloads


.. include:: ./downloads.rst



Install from source
----------------------

First `install Go`_ on your system. Then just use

::

    go get -u github.com/schollz/sdees

.. _the latest release binary: https://github.com/schollz/sdees/releases/latest
.. _install Go: https://golang.org/dl/
.. _version 2.5.0, released March, 2016: https://git-scm.com/docs/git-for-each-ref/2.5.0


