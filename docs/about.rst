About
======================================

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
Dropbox. *gojot* just makes this solution easier to attain.

.. _which is easy to get on any system: https://git-scm.com/downloads

How *gojot* works
--------------------

When *gojot* starts for the first time it will request a ``git``
repository which can be local or remote. Then *gojot* provides an
option to write text using the editor of choice into a new *entry*. Each
new entry is committed to the supplied ``git``
repository as new file. Thus, *gojot* makes it perfectly safe to make
*new* entries without internet access without worrying of overwriting a document.



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