# -*- coding: utf-8 -*-

import click
from gojot import gojot


@click.command()
@click.argument('repo', nargs=1)
@click.option('--doc', default=None)
@click.option('--imp', default=None)
def main(repo, doc, imp, args=None):
    """Console script for gojot"""
    click.echo("Replace this message by putting your code into "
               "gojot.cli.main")
    click.echo("See click documentation at http://click.pocoo.org/")
    if imp != None:
        gojot.run_import(repo, imp)
    else:
        gojot.run(repo, doc)
    # TODO
    # Add flag for editing a single entry (boolean)


if __name__ == "__main__":
    main()
