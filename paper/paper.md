---
title: 'bíogo/ncbi: interfaces to NCBI services for the Go language'
tags:
  - bioinformatics
  - toolkit
  - golang
authors:
 - name: R Daniel Kortschak
   orcid: 0000-0001-8295-2301
   affiliation: 1
 - name: David L Adelson
   orcid: 0000-0003-2404-5636
   affiliation: 1
affiliations:
 - name: School of Biological Sciences, The University of Adelaide
   index: 1
date: 3 March 2017
bibliography: paper.bib
---

# Summary

The National Center for Biotechnology Information makes available BLAST sequence similarity search [@BLAST] and health science database search through the Entrez service [@Entrez].
In addition to an interactive web interface, BLAST and Entrez provide an application programmer interface to allow programmatic use of these services via the BLAST URL API [@BLASTURLAPI] and the Entrez EUtilities [@EUtilities].
The BioPerl suite [@BioPerl] provides access to the BLAST API via Bio::Tools::Run::StandAloneBlastPlus [@BioPerlRUN] and to Entrez via Bio::Tools::EUtilities [@BioPerlEUtilities].
Similarly, Biopython [@Biopython] provides access via the NCBIWWW function in the Bio.Blast module and functions in Bio.Entrez for EUtilities.
Packages within bíogo/ncbi provide Go application programmer interfaces to the NCBI BLAST and EUtilities services.
The design of bíogo/ncbi is light weight, allowing the user to make use of the Go language's control structures and data types, rather than imposing a library-specific access approach.
In addition to allowing remote BLAST searches, BioPerl and Biopython provide mechanisms to parse XML output from local BLAST search via BioPerl's Bio::SearchIO and Biopythons Bio.Blast NCBIXML.
Because of the modular design of bíogo/ncbi, XML encoded output from local BLAST searches or remote downloads can be parsed using the Go standard library's XML package.

# References
