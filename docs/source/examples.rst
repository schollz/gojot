Help
==========

First time startup
--------------------

On the first use you will be asked to point *gojot* to a repository.
You will also be able to select which editor you'd like to use.
Currently only four editors are supported: `micro`_, `vim`_, `emacs`_, and `nano`_.
All released versions of *gojot* (2.0.0+) include `micro`_ built-in.

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-firsttime.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>

Adding an entry
--------------------

Adding an entry is as simple as starting the program and selecting the document to put the entry in.

Here is a simple example of adding an entry:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-adding.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>

Each entry is assigned a random name.
The random name is only used to keep track of things internally, but you can also
use the name to quickly edit a specific entry. To do this, simply delete the randomized entry name and type in your own.

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-changingentry.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>

Then, to edit a specific entry, simply type in the entry name after running the program:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-editingentry.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>


``--delete``
--------------

You can delete whole documents, or single entries. In both cases you just name the item you want to delete.

.. warning::

    Deleting is permanent! All files locally and remotely will be deleted.

Here is an example for deleting a document:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-deletedoc.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>


And here is an example of deleting an entry:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-deleteentry.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>



``--all``
-------------------------------------

Its often useful to see old entries, which you can edit as well. To load the whole document, simply add the flag `--all`. Or you can simply select `y` when prompted about loading the whole document.

Here is an example of loading a whole document:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-all.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>



``--import``
------------------------

You can import prevoius *gojot* journals using `import`.

.. raw:: html

    <center>
    <asciinema-player src="/_static/asciicast/asciinema-import.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
    </center>


``--importold`` (DayOne importing)
-----------------------------------

You can import a DayOne type journal using `--importold`.

.. raw:: html

    <center>
    <asciinema-player src="/_static/asciicast/asciinema-importold.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
    </center>



``--export``
------------------------

You can export your whole document as a text-file using the `--export` flag. Here is an example:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-export.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>



``--stats``
----------------------------

If you'd like to get information about wordcounts and entry counts in all your documents, simply use the `--stats` flag:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-stats.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>



``--config``
----------------------------

With the `--config` flag you can change the repository that is being used and the editor that is being used.

.. raw:: html

    <center>
    <asciinema-player src="/_static/asciicast/asciinema-config.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
    </center>




``--clean``
----------------------------

With the `--clean` flag you can erase all the *gojot* folders. This includes the cache of currently known ``git`` repositories, ``$HOME/.cache/gojot``,
as well as any configuration files in ``$HOME/.config/gojot``. You will be prompted to verify that this is what you want.

.. raw:: html

    <center>
    <asciinema-player src="/_static/asciicast/asciinema-clean.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
    </center>


``--summary``
----------------------------

All entries can be succintly summarized using the `--summary` flag. This will show the date, the entry name, number of words, and then the first few words in the entry.

Here is an example:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-summary.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>

Problems?
----------

If you have any problems at all, please `submit an Issue`_ and someone can help you sort it out.

.. _submit an Issue: https://github.com/schollz/gojot/issues/new
.. _micro: https://github.com/zyedidia/micro
.. _vim: http://www.vim.org/download.php
.. _nano: https://www.nano-editor.org/
.. _emacs: https://www.gnu.org/software/emacs/
