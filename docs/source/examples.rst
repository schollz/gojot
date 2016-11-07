Help
==========

First time startup
--------------------

On the first use you will be asked to point *sdees* to a repository. You will also be able to select which editor you'd like to use. Currently only four editors are supported: micro, vim, emacs, and nano.

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/4afttloazk5indyvszueds7ot.js" id="asciicast-4afttloazk5indyvszueds7ot" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

Adding an entry
--------------------

Adding an entry is as simple as starting the program and selecting the document to put the entry in.


Here is a simple example of adding an entry:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/bwgsyl52a7vmtsbwwvbjxm4ii.js" id="asciicast-bwgsyl52a7vmtsbwwvbjxm4ii" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

Each entry is assigned a random name, which you can also change to whatever you want. To do this, simply delete the randomized entry name and type in your own.

.. warning::

    Do not store sensitive information in the entry names, as the OTP encryption could be broken if you store thousands of entries. Also make sure the entry names are unique!

Here is an example:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/0jky6qw2cccm54zcfs99awpwq.js" id="asciicast-0jky6qw2cccm54zcfs99awpwq" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>



``--delete``
--------------

You can delete whole documents, or single entries. In both cases you just name the item you want to delete. 

.. warning::

    Deleting is permanent! All files locally and remotely will be deleted.

Here is an example for deleting a document:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/8y4l3uj4ygodujf5kzpdaj3in.js" id="asciicast-8y4l3uj4ygodujf5kzpdaj3in" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

And here is an example of deleting an entry:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/7gegbsn3eiw81t3qwf6urzgxh.js" id="asciicast-7gegbsn3eiw81t3qwf6urzgxh" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>

``--all``
-------------------------------------

Its often useful to see old entries, which you can edit as well. To load the whole document, simply add the flag `--all`. Or you can simply select `y` when prompted about loading the whole document.

Here is an example of loading a whole document:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/drgeyew22fagwmezn7ff9qdxb.js" id="asciicast-drgeyew22fagwmezn7ff9qdxb" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>



``--export``
------------------------

You can export your whole document as a text-file using the `--export` flag. Here is an example:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/boocc1jgv4quydvhyom35pr4p.js" id="asciicast-boocc1jgv4quydvhyom35pr4p" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>


``--stats``
----------------------------

If you'd like to get information about wordcounts and entry counts in all your documents, simply use the `--stats` flag:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/8ifdjbwe2ujgcqlzvjw7uzkl5.js" id="asciicast-8ifdjbwe2ujgcqlzvjw7uzkl5" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>


``--summary``
----------------------------

All entries can be succintly summarized using the `--summary` flag. This will show the date, the entry name, number of words, and then the first few words in the entry.

Here is an example:

.. raw:: html

   <center>
   <script type="text/javascript" src="https://asciinema.org/a/478ig89xc0dsr0bmsk8nco6mz.js" id="asciicast-478ig89xc0dsr0bmsk8nco6mz" async data-autoplay="false" data-preload="true" data-size="small" data-speed="0.9" data-theme="asciinema"></script>
   </center>
