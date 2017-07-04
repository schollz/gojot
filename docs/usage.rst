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




