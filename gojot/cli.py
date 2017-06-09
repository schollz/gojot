# -*- coding: utf-8 -*-

import click
from gojot import gojot


@click.command()
@click.option('--repo', '-r', default=None, help='Repo to open.')
@click.option('--doc', '-d', default=None, help='Document to open.')
@click.option('--imp', '-i', default=None, help='File to import.')
@click.option('--all', '-a', 'load_all', is_flag=True, help='Open all entries.')
@click.option('--edit', '-e', 'edit_one', is_flag=True, help='Edit an entry.')
@click.option('--export', '-x', 'export', is_flag=True, help='Export a document.')
def main(repo, doc, imp, load_all, edit_one, export, args=None):
    """gojot is a modern command-line journal that is distributed and encrypted by default."""
    click.echo("")
    click.echo("See gojot documentation at https://gojot.schollz.com/")
    if imp != None:
        gojot.run_import( imp)
    else:
        gojot.run(repo, doc, load_all=load_all, edit_one=edit_one,export=export)


if __name__ == "__main__":
    main()
