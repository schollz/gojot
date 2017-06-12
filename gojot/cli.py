# -*- coding: utf-8 -*-

import click
from gojot import gojot


@click.command()
@click.option('--repo', '-r', default=None, help='Repo to open.')
@click.option('--document', '-d', default=None, help='Document to open.')
@click.option('--import', '-i', 'imp', default=None, help='File to import.')
@click.option('--all', '-a', 'load_all', is_flag=True, help='Open all entries.')
@click.option('--edit', '-e', 'edit_one', is_flag=True, help='Edit an entry.')
@click.option('--export', '-x', 'export', is_flag=True, help='Export a document.')
@click.option('--stats', '-s', 'show_stats', is_flag=True, help='Show document stats.')
@click.option('--version', '-v', 'version', is_flag=True, help='Print version.')
def main(repo, document, imp, load_all, edit_one, export, show_stats, version, args=None):
    """gojot is a modern command-line journal that is distributed and encrypted by default."""
    click.echo("")
    click.echo("See gojot documentation at https://gojot.schollz.com/")
    if imp != None:
        gojot.run_import(repo, imp)
    elif version != False:
        print("gojot version 3.0.2")
    else:
        gojot.run(repo, document, load_all=load_all,
                  edit_one=edit_one, export=export, show_stats=show_stats)


if __name__ == "__main__":
    main()
