Installation Guide
==================

Requirements
----------------

Before you install *sdees* make sure you have the following requirement:

* **git** (version 2.5+)
    *git* can be downloaded for free from `git-scm.com`_.

If you are on Windows/OSX just `download the latest version of git`_.

If you are on Ubuntu, make sure to add the latest repo:

::

    add-apt-repository ppa:git-core/ppa -y
    apt-get update
    apt-get install git -y

.. _download the latest version of git: https://git-scm.com/downloads
.. _git-scm.com: https://git-scm.com/downloads


.. note::

    To determine branches that are ahead/behind, this program uses ``git for-each-ref``
    with the ``push:track`` option, which is not introduced until
    `version 2.5.0, released March, 2016`_.
    The alternative to this is ``git branch -vv`` but that is not considered stable.

Download
---------

Then, to install *sdees*, simply download `the latest release binary`_

*OR*

use ``go get`` if you have `installed Go`_:

::

    go get -u github.com/schollz/sdees

.. _the latest release binary: https://github.com/schollz/sdees/releases/latest
.. _installed Go: https://golang.org/dl/
.. _version 2.5.0, released March, 2016: https://git-scm.com/docs/git-for-each-ref/2.5.0
