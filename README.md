# fintools-public ![release](https://github.com/filmil/fintools-public/actions/workflows/release.yml/badge.svg) ![test](https://github.com/filmil/fintools-public/actions/workflows/test.yml/badge.svg)

A collection of tools for financial calculations that can be publicized.

## Google paystub converter for beancount

This directory contains the source code to a small program I wrote to convert
Google paystubs to beancount transaction records.

[bc]: https://beancount.furius.ca

As you may know Google's paystubs are only available as PDF files to employees.
For folks who are like me and like doing their own text-based accounting, this
means having to retype manually the PDF reports into beancount.

I thought there may be a way to automate this a bit, and wrote this program.

The program is far from finished at the moment but the bits that work seem to
work robustly.

## Prerequisites

### go development environment

### pdf2txt from `python-pdfminer` package in debian.

Install with:

```
sudo apt-get install python-pdfminer
```

## Installation

```
git clone git@github.com:filmil/fintools.git
cd fintools
go install ./cmd/...
```

The above commands will install the program `paystub`. This is a very generic
name.  While better naming scheme is welcome, the reason for this is that my
beancount setup (unpublished) is ran in a docker container, with only the 
needed utilities.

### Testing

```
cd fintools
go test ./...
```

At the moment expect above steps to fail because tests were removed on purpose.
:(

## Using `paystub`

### Download the paystub file as a PDF

Download from Google's paystub front-end the PDF file of your pay stub.  Suppose
you download it as to `paystub.pdf`.  The actual name will be different and it
is *always* different for some reason but what to do.

### Convert the paystub file from PDF to XML

This step requires `pdf2txt`:

```
pdf2txt -t xml -o paystub.xml paystub.pdf
```

### Convert the XML file into a beancount transaction

```
paystub -input=paystub..xml
```

This will print a transaction in beancount format.  The transaction by default
uses account names that I use in my own beancount files.  You may want to tweak
the account names.  There is currently no documentation for those flags, please
see [paystub/main.go][pmg] for all the flags that are defined, or use:

```
paystub --help
```

[pmg]: https://github.com/filmil/fintools/tools/cnmd/paystub/main.go

## Using `payxml`

The program `payxml` produces a bounding box drawing of the paystub. I wrote
it to visualize what is being analyzed by `paystub`. I will assume that you
have the `paystub.xml` from when you tested `paystub`


## How the paystub is parsed

The idea is very simple: `pdf2txt` produces an XML file in which the paystub
text is laid out together with bounding boxes that describe where on the paystub
page the text goes.  The coordinate system begins in the lower-left corner of
the page and X axis grows to the right and Y axis grows upwards. The unit is
1/72th of an inch. 

Parsing the bounding box layout is more robust than parsing the text output
produced by`pdf2txt -t txt)`, since the paystub generator uses unstable
iteration to typeset the paystub so text output could be arbitrarily
rearranged.

Luckily with XML and bounding boxes we don't need to worry about that. Instead
there is a small query language that allows you to look up a label or a
bounding box, and then filter the message based on that.  See the file `xml.go`
and `query.go` for details.

This query language allows you to find a number of "anchors" , like "Pay type"
or "deductions".  You can also make bounding boxes to say: extend this bounding
box from the place where "Pay type" text label appears to the right until you
hit the label "Earnings" and extend it down to the level of "Taxes".

Once you do that, you can do things like: find column "Pay", and get all the
numbers below it, find all labels below "Type" and get all text below it and
match the two up.  Once you have a matching, add to the `Transaction`.

Once this structure is built out, it is written using go text templates.

## Bugs and Limitations

* I had to remove tests in order to publish this program. :( Ideas welcome on
how to make the program self-contained with anonymous tests.

* Currently there is no automation for the conversion process.  I built that
part for me personally, but it's currently difficult to extract and publicize.

* Not all parts of the paystub are currently supported. I haven't finished that
bit, as I had a limited number of paystubs to import and only wrote code that
was enough to import those paystubs.

* This is based on my own paystubs which are for a particular Google US office.
I have no idea what paystubs look like for other offices or countries.

* Multiple pay-to accounts are not supported simply because I don't do that. It
would be easy to add if you are so inclined.

* Needs more documentation and testing, for sure.

# Contributions

Pull requests are welcome.
