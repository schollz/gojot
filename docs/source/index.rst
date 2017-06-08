Overview
=========

*gojot* is a modern command-line journal that is distributed and encrypted by default.


Ok. But, really, *gojot* is just a fancy wrapper for ``git`` that allows
you to make time-stamped entries to encrypted documents while keeping
the entire document synchronized. *gojot* is great for a *distributed
and encrypted* journal and is compatible with all major operating
systems. *gojot* is a single binary (with a text-editor built-in!), so
you only need ``git`` to get started and is compatible any local repo or
hosted service (`Gitlab`_/`Bitbucket`_/`Github`_).

Here's what it looks like in action (`check if its encrypted`_):

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-overview.json" async autoplay="true" size="small" speed="2.0" theme="asciinema"></asciinema-player>
   </center>



Features
------------

-  *Single* dependency
    Requires `git`_ (version 2.5+).
-  *Single* binary
    Available on `all major systems and architectures`_.
-  Built-in text editor
    `micro`_ is built-in by default, but can also uses `vim`_/`emacs`_/`nano`_.
-  Fulltext encryption
    Uses `OpenPGP`_, compatible with ``gpg``.
-  Filename encryption
    Uses `ChaCha20`_.
-  Built-in version control
    All versions are saved, currently only newest is shown.
-  Lots of other neat features.
    Searching, summarizing, synchronized deletion, self-updating, collision management, and more.
-  Open-source
     `Source on Github`_, written in `Go`_ and licensed under the MIT, with exception of a few open source components from third parties.

.. _check if its encrypted: https://github.com/schollz/demo
.. _see mine: https://github.com/schollz/demo
.. _OpenPGP: https://en.wikipedia.org/wiki/Pretty_Good_Privacy#OpenPGP
.. _ChaCha20: https://en.wikipedia.org/wiki/Salsa20#ChaCha_variant
.. _all major systems and architectures: /install.html
.. _Source on Github: https://github.com/schollz/gojot
.. _Gitlab: https://gitlab.com/users/sign_in
.. _Bitbucket: https://bitbucket.org/account/signin/
.. _Github: https://github.com/
.. _micro: https://github.com/zyedidia/micro
.. _vim: http://www.vim.org/download.php
.. _nano: https://www.nano-editor.org/
.. _emacs: https://www.gnu.org/software/emacs/
.. _Go: https://golang.org/
.. _git: https://git-scm.com/



.. toctree::
    :hidden:
    :maxdepth: 2

    Overview <index>
    About <about>
    Install <install>
    Help <examples>
    Advanced <advanced>
