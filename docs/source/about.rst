About
=========

Why *gojot*
-----------

**Problem:** I would like to write into a document (e.g., a journal with
time stamped entries) and also have it to be compatible with all the computers
I use (Windows + Linux, work + home),
and I want it to be locally encrypted. Unfortunately, typical cloud
storage synchronization utilities (e.g. Dropbox) do not play well with
syncing a single document that is encrypted locally.

**Solution:** The basic solution is to create a lot of files - each
encrypted - where each file contains the text for one entry in the
document. You don’t need any special tool for this - you can just use a
text editor, ``gpg``, and synchronization software like ``git`` or
Dropbox.

*gojot* just makes this solution easier to attain.
*gojot* is a single executable file with the text-editor and
``gpg`` bulit-in, compatible with all major OS/architectures.
The only system requirement is the installation of
the distributed version control software ``git``,
`which is easy to get on any system`_. The other benefit
of *gojot* is that it will automatically combine all the time-stamped
entries so it appears that you are editing a single document, and it
will also resolve merges that can occur if you edit the same entry
offline on two computers.

.. _which is easy to get on any system: https://git-scm.com/downloads


How *gojot* works
--------------------

When *gojot* starts for the first time it will request a ``git``
repository which can be `local or remote`_. Then *gojot* provides an
option to write text using the editor of choice into a new *entry*. Each
new entry is inserted into a new orphan branch in the supplied ``git``
repository. The benefit of placing each entry into its own orphan branch
is that merge collisions are avoided after creating new entries on
different machines. Thus, *gojot* makes it perfectly safe to make
*new* entries without internet access.

The combination of entries be displayed as a *document*. The document is
reconstructed by 1) fetching all remote branches, 2) filtering out which
ones contain entries for the document of interest, 3) decrypting each
entry, and then 4) sorting the entries by the commit date. Merge
conflicts should only occur when simultaneous edits are made to the same
entry and they are resolved by combining the diffed versions of the
entries. Multiple documents can be stored in a single ``git``
repository.

All information saved in the ``git`` repo is encrypted. The text of each entry
is encrypted using `OpenPGP`_ with a symmetric cipher with a user-provided passphrase.
The passphrase is not stored anywhere on the machine or repo.
All the filenames and branch names are encrypted using `ChaCha20`_ with a key generated
upon initialization. When editing an encrypted document, a decrypted temp file is
stored and then shredded (random bytes written and then deleted) after
use. Thus, you can use gojot with a public git repository (for example, `see mine`_) without
revealing information since it will look like gibberish, for example:

.. raw:: html

   <center>
  <asciinema-player src="/_static/asciicast/asciinema-about.json" async preload="true" size="small" speed="0.9" theme="asciinema"></asciinema-player>
   </center>



Alternatives to *gojot*
------------------------

Here is some software which are similar to *gojot*, but often require other software
or system-specific utilities.
I engojoy these software, and used a lot of inpsiration from them, but ultimately I found
that *gojot* could provide some functionality or utility that was still absent.

*  **13 lines of shell** `[site] <https://gist.github.com/schollz/27b4ffe562b0b74bf8ee1e8055680d22>`_ - git-based journal
    No encryption, no editing past entries, no version control, no deletion - but only 13 lines!
*  **Cryptpad** `[site] <https://beta.cryptpad.fr/pad/>`_ - zero knowledge realtime encrypted editor in your web browser
    Requires internet access, and a browser, difficult to reconstruct many documents.
*  **jrnl.sh** `[site] <http://jrnl.sh/>`_ - command line journal application with encryption
    Requires Python. Syncing only available via Dropbox which won't support merging encrypted files if editing offline.
*  **ntbk** `[site] <hhttps://www.npmjs.com/package/ntbk>`_ - command line journal
    Requires Node.js. Syncing only available via Dropbox which won't support merging encrypted files if editing offline.
*  **vimwiki** `[site] <http://vimwiki.github.io/>`_ - command line editor with `capability of distributed encryption <http://www.stochasticgeometry.ie/2012/11/23/vimwiki/>`_
    Requires system-specific filesystem encryption (e.g. `eCryptFS <http://ecryptfs.org/>`_). Works with any DVCS, but merges are not handled.
*  **Org mode** `[site] <http://orgmode.org/>`_ - fast and effective plain-text system for authoring documents
    Requires ``emacs``, requires adding DVCS later

Limitations
------------

*gojot* is not meant as an encrypted file system, as it has limits to
the number of entries that can be stored.

Entry names need to be unique
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Entry names allow you quick access to a specific entry without having to recall
a date. Entry names should be unique, as currently **gojot** will load a random
entry if their are two entries with the same name.

Collision of entry names
~~~~~~~~~~~~~~~~~~~~~~~~
Currently there are only 14,260,682,650 adjective+verb combinations
available for random entry names. Thus, a `collision probability`_ of 50%
will occur after ~120,000 entries. Collisions are not detrimental, but
it will only allow one document to be loaded with the same entry name.
The reason that this happens is technical, and `is slated to be
resolved`_.


``git`` version 2.5+
~~~~~~~~~~~~~~~~~~~~~~

To determine branches that are ahead/behind, this program uses ``git for-each-ref``
with the ``push:track`` option, which is not introduced until
`version 2.5.0, released March, 2016`_.
The alternative to this is ``git branch -vv`` but that is not considered stable.

.. _local or remote: https://github.com/schollz/gojot/blob/master/INFO.md#setting-up-git-server
.. _see mine: https://github.com/schollz/demo
.. _all major systems and architectures: /install.html
.. _is slated to be resolved: https://github.com/schollz/gojot/issues/73
.. _version 2.5.0, released March, 2016: https://git-scm.com/docs/git-for-each-ref/2.5.0
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
.. _collision probability: https://en.wikipedia.org/wiki/Birthday_problem#Approximation_of_number_of_people

Acknowledgements
-----------------

*gojot* was written by `schollz`_. There are number of third-party code snippets and imports
(see source for attribution and License information for each),
and I am very grateful to these authors for their code:
`mholt`_, `jbenet`_, and `aryann`_.

.. _schollz: https://schollz.com
.. _mholt: httsp://github.com/mholt
.. _jbenet: https://github.com/jbenet
.. _aryann: https://github.com/aryann
