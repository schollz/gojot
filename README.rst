=====
gojot
=====


.. image:: https://img.shields.io/pypi/v/gojot.svg
        :target: https://pypi.python.org/pypi/gojot

.. image:: https://img.shields.io/travis/schollz/gojot.svg
        :target: https://travis-ci.org/schollz/gojot

.. image:: https://readthedocs.org/projects/gojot/badge/?version=latest
        :target: https://gojot.readthedocs.io/en/latest/?badge=latest
        :alt: Documentation Status


.. image:: https://pyup.io/repos/github/schollz/gojot/shield.svg
     :target: https://pyup.io/repos/github/schollz/gojot/
     :alt: Updates

*gojot* is a modern command-line journal that is distributed and encrypted by default.


OK. But, really, *gojot* is just a fancy wrapper for ``git``, ``gpg`` and ``vim`` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized in it a ``git`` repository. 


Here's what it looks like in action (`check if its encrypted`_):

.. image:: /docs/_static/demo2.gif
     :alt: Updates

=======
Install
=======

First make sure you have ``gpg``, ``git``, and ``vim`` installed:

.. code-block:: console

    $ sudo apt-get install gpg git vim


Then you can install gojot using ``pip``:

.. code-block:: console

	$ pip install gojot



.. _check if its encrypted: https://github.com/schollz/demo


=====
Usage
=====

First-time use
---------------

For the first time setup, just use::

    gojot

If you do not have any GPG keys you should first generate one with::

	gpg --gen-key


Examples
---------

To edit a specific entry in a specific document you can do (you'll be prompted with a list of entries to choose from)::

	gojot --document notes --edit

To see stats about your writing::

	gojot --stats

To change to a different repo::

	gojot --repo git@github.com/NAME/REPO.git

To export your file::

	gojot --export







