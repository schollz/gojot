Overview
=========

*sdees* is for *distributed* *editing* of *encrypted* *stuff*.

*sdees does editing, encryption* and *synchronization*.


But, really, *sdees* is just a fancy wrapper for ``git`` that
allows you to make time-stamped entries to encrypted documents while
keeping the entire document synchronized. *sdees* is great for a
*distributed and encrypted* journal and is compatible any local repo or
hosted service (Gitlab/Bitbucket/Github).

Here's what it looks like in action:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/cwc20isy2dz7vepn05sghsdgn.js" id="asciicast-cwc20isy2dz7vepn05sghsdgn" async data-autoplay="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

All local files are encrypted, and all temporary files are shredded. You can browse the local files, but this is what you will see:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/7ph7viak75hgbeqnfvayeafu3.js" id="asciicast-7ph7viak75hgbeqnfvayeafu3" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

Features
--------

-  *Only one* dependency: ``git`` (version 2.5+)
-  Cross-compatibility (Windows/Linux/OS X).
-  Built-in encryption, compatible with ``gpg``.
-  Built-in text editor, ``micro`` (but options for
   ``emacs``/``vim``/``nano``).
-  Built-in version control (all versions are saved, currently only
   newest is shown).
-  Encrypted filenamesmostly.
-  Searching, summarizing, synchronized deletion, self-updating,
   collision management, and more.


.. toctree::
    :hidden:
    :maxdepth: 2

    Overview <index>
    About <about>
    Install <install>
    Usage <examples>
    Advanced <advanced>
