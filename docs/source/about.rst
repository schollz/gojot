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

All information saved in the ``git`` repo is encrypted using a symmetric
cipher with a user-provided passphrase. The passphrase is not stored
anywhere on the machine or repo. Each entry in the ``git`` repo is
encrypted. When editing an encrypted document, a decrypted temp file is
stored and then shredded (random bytes written and then deleted) after
use. All filenames are encrypted using a one-time pad whose pages are
generated at the initialization.

.. _local or remote: https://github.com/schollz/sdees/blob/master/INFO.md#setting-up-git-server



Limitations
------------

*sdees* is not meant as an encrypted file system, as it has limits to
the number of entries that can be stored.

Possible collisions in entry names
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Currently there are only 14,260,682,650 adjective+verb combinations
available for random entry names. Thus, a collision probability of 50%
will occur after ~120,000 entries. Collisions are not detrimental, but
it will only allow one document to be loaded with the same entry name.
The reason that this happens is technical, and `is slated to be
resolved`_.

.. warning::

    Do not store sensitive information in the entry or document names! They are encrypted,
    but they are not encrypted as securely as the text of the entry.
    
Weak encryption of filenames
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

While the text of each entry is securely encrypted using a
GPG-compatible symmetric cipher, the filenames are encrypted using a OTP
from a static source of random bytes. The OTP is useful for filenames
since encrypted strings are short enough to be used for branch names and
filenames (which are limited to 255 characters). Pads are used randomly
(since usage cannot be synced), and generally only 20-30 bytes will be
used at a time. Still, this means a probability of 50% to overlap could
start to occur after ~600 entries and a complete collision could start
occurring with 50% probability after ~300 documents. This would only
allow an attacker to reveal the names of two documents, though, and not
any of the information inside the documents (as that is stored under
GPG). So do not store credit-card information in the names of your
files!

.. _is slated to be resolved: https://github.com/schollz/sdees/issues/73
