About
=========

Why *sdees*
-----------

**Problem:** I would like to write into a document (e.g., a journal with
time stamped entries) and also have it be available on all computers I
use, and I want it to be locally encrypted. Unfortunately, typical cloud
storage synchronization utilities (e.g. Dropbox) do not play well with
syncing a single document that is encrypted locally.

**Solution:** The basic solution is to create a lot of files - each
encrypted - where each file contains the text for one entry in the
document. You don’t need any special tool for this - you can just use a
text editor, ``gpg``, and synchronization software like ``git`` or
Dropbox.

My program, *sdees*, just makes this solution easier to attain.
*sdees* comes as a single executable file with the text-editor and
``gpg`` bulit-in - the only system requirement is the installation of
``git`` (which is pretty easy to get on any system). The other benefit
of *sdees* is that it will automatically combine all the time-stamped
entries so it appears that you are editing a single document, and it
will also resolve merges that can occur if you edit the same entry
offline on two computers.


How *sdees* works
--------------------

When *sdees* starts for the first time it will request a ``git``
repository which can be `local or remote`_. Then *sdees* provides an
option to write text using the editor of choice into a new *entry*. Each
new entry is inserted into a new orphan branch in the supplied ``git``
repository. The benefit of placing each entry into its own orphan branch
is that merge collisions are avoided after creating new entries on
different machines. Thus, *sdees* makes it perfectly safe to make
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
is encrypted using OpenPGP with a symmetric cipher with a user-provided passphrase.
The passphrase is not stored anywhere on the machine or repo.
All the filenames and branch names are encrypted using ChaCha20 with a key generated
upon initialization. When editing an encrypted document, a decrypted temp file is
stored and then shredded (random bytes written and then deleted) after
use.

.. _local or remote: https://github.com/schollz/sdees/blob/master/INFO.md#setting-up-git-server

Alternatives to *sdees*
------------------------

Here are some software which operate very similarly to *sdees*, but often require other software or system-specific utilities.
I enjoy these software, and the use of some of them was the inpsiration for *sdees* (i.e. ``jrnl.sh``):

*  **Cryptpad** `[site] <https://beta.cryptpad.fr/pad/>`_ - zero knowledge realtime encrypted editor in your web browser
    Requires internet access, and a browser, difficult to reconstruct many documents.
*  **jrnl.sh** `[site] <http://jrnl.sh/>`_ - command line journal application with encryption
    Requires Python. Syncing only available via Dropbox which won't support merging encrypted files if editing offline.
*  **vimwiki** `[site] <http://vimwiki.github.io/>`_ - command line editor with `capability of distributed encryption <http://www.stochasticgeometry.ie/2012/11/23/vimwiki/>`_
    Requires system-specific filesystem encryption (e.g. `eCryptFS <http://ecryptfs.org/>`_). Works with any DVCS, but merges are not handled.
*  **Org mode** `[site] <http://orgmode.org/>`_ - fast and effective plain-text system for authoring documents
    Requires ``emacs``, requires adding DVCS later

Limitations
------------

*sdees* is not meant as an encrypted file system, as it has limits to
the number of entries that can be stored.

Entry names need to be unique
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Entry names allow you quick access to a specific entry without having to recall
a d ate. Entry names should be unique, as currently **sdees** will load a random
entry if their are two entries with the same name.

Collision of entry names
~~~~~~~~~~~~~~~~~~~~~~~~
Currently there are only 14,260,682,650 adjective+verb combinations
available for random entry names. Thus, a collision probability of 50%
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

.. _is slated to be resolved: https://github.com/schollz/sdees/issues/73
