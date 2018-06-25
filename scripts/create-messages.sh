#!/bin/sh

# This scripts creates new PO files and fill them with existing translations
#
# Must be run from the misc/translator directory

tempdir=$(mktemp -d)
echo $tempdir

./translate create

languages=$(find po -name '*-pref.po' -printf "%f " | sed 's/-pref.po//g')
for lang in $languages; do
    echo $lang
    for i in marc-MARC21 marc-NORMARC marc-UNIMARC opac-bootstrap pref staff-help staff-prog; do
        msgattrib --force-po -o $tempdir/$lang-$i-no-fuzzy.po --translated --no-fuzzy po/$lang-$i.po
        msgattrib --force-po -o $tempdir/$lang-$i-only-fuzzy.po --translated --only-fuzzy po/$lang-$i.po
    done

    msgcat --force-po -o $tempdir/$lang.po --use-first $tempdir/$lang-*-no-fuzzy.po $tempdir/$lang-*-only-fuzzy.po
    msgmerge --force-po -o $tempdir/$lang-messages.po $tempdir/$lang.po po/$lang-messages.po
    msgattrib --force-po -o po/$lang-messages.po --no-obsolete $tempdir/$lang-messages.po
done

