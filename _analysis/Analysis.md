# Analysis Report


### Word Statistics

* Lexicon - based on default
* Islands - are words that changing any letter will not form another word
* Doublets - are words that changing any letter will only form one other word
* LDS% (local decay smoothness) – the longest consecutive run of ladder lengths for which each word count is at least 95% of the previous ladder length’s count, expressed as a percentage of all ladder lengths in the dictionary
* Var – variance of consecutive count ratios, measuring how smoothly (low) or unevenly (high) word counts decay across ladder lengths
* Drop% – percentage decrease from ladder length 2 to 3, representing initial connectivity decay
* Longest - is the longest possible ladder length for the dictionary
* Numeric columns - are the number of words that can form that ladder length

| Letters | Words | Islands |   %  | Doublets |   %  | LDS% |  Var | Drop% | Longest |     2 |     3 |     4 |     5 |     6 |     7 |     8 |     9 |    10 |    11 |    12 |    13 |    14 |    15 |    16 |    17 |    18 |    19 |    20 |    21 |    22 |    23 |    24 |    25 |    26 |    27 |    28 |    29 |    30 |    31 |    32 |    33 |    34 |    35 |    36 |    37 |    38 |    39 |    40 |    41 |    42 |    43 |    44 |    45 |    46 |    47 |    48 |    49 |    50 |    51 |    52 |    53 |    54 |    55 |    56 |    57 |    58 |    59 |    60 |    61 |
|--------:|------:|--------:|-----:|---------:|-----:|-----:|-----:|------:|--------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|
|       2 |   127 |       0 |  0.0 |        0 |  0.0 | 75.0 | 0.12 |   0.0 |       5 |   127 |   127 |   127 |    33 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       3 |  1285 |       1 |  0.1 |        0 |  0.0 | 75.0 | 0.12 |   0.0 |       9 |  1284 |  1284 |  1284 |  1284 |  1284 |  1245 |   723 |    18 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       4 |  5191 |      65 |  1.3 |       16 |  0.3 | 62.5 | 0.12 |   0.3 |      17 |  5126 |  5110 |  5103 |  5101 |  5101 |  5101 |  5101 |  5101 |  5101 |  5052 |  4365 |  2572 |   653 |    95 |    25 |     5 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       5 | 11798 |     861 |  7.3 |      310 |  2.6 | 58.6 | 0.07 |   2.8 |      30 | 10937 | 10627 | 10536 | 10486 | 10449 | 10432 | 10430 | 10427 | 10427 | 10427 | 10427 | 10427 | 10427 | 10427 | 10427 | 10427 | 10393 |  9839 |  7967 |  5583 |  3230 |  1488 |   547 |   215 |    87 |    46 |    21 |     9 |     2 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       6 | 20952 |    3925 | 18.7 |     1898 |  9.1 | 55.1 | 0.03 |  11.1 |      50 | 17027 | 15129 | 14312 | 13865 | 13617 | 13476 | 13379 | 13320 | 13274 | 13242 | 13229 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13222 | 13217 | 13177 | 13001 | 12518 | 11267 |  8929 |  6160 |  4150 |  2853 |  1942 |  1317 |   907 |   633 |   434 |   269 |   162 |    98 |    62 |    42 |    28 |    18 |    13 |     7 |     3 |       |       |       |       |       |       |       |       |       |       |       |
|       7 | 31498 |    9605 | 30.5 |     4874 | 15.5 | 57.7 | 0.03 |  22.3 |      53 | 21893 | 17019 | 14758 | 13501 | 12849 | 12437 | 12163 | 11999 | 11885 | 11783 | 11690 | 11614 | 11552 | 11493 | 11432 | 11382 | 11351 | 11324 | 11312 | 11298 | 11285 | 11275 | 11264 | 11256 | 11251 | 11251 | 11251 | 11251 | 11250 | 11234 | 11176 | 10987 | 10529 |  9791 |  8868 |  7895 |  7074 |  6269 |  5242 |  3815 |  2441 |  1516 |   941 |   588 |   369 |   221 |   141 |    85 |    50 |    26 |    13 |     4 |       |       |       |       |       |       |       |       |
|       8 | 39095 |   16959 | 43.4 |     8331 | 21.3 | 25.0 | 0.02 |  37.6 |      61 | 22136 | 13805 | 10237 |  8523 |  7594 |  7014 |  6646 |  6387 |  6169 |  6006 |  5863 |  5717 |  5570 |  5394 |  5187 |  4988 |  4797 |  4600 |  4348 |  4006 |  3639 |  3327 |  3031 |  2777 |  2544 |  2371 |  2247 |  2160 |  2094 |  2042 |  1992 |  1923 |  1859 |  1806 |  1753 |  1698 |  1653 |  1607 |  1560 |  1509 |  1437 |  1327 |  1149 |   865 |   579 |   418 |   340 |   290 |   256 |   223 |   179 |   137 |   102 |    78 |    69 |    58 |    46 |    31 |    16 |     4 |
|       9 | 39967 |   22021 | 55.1 |     9654 | 24.2 |  0.0 | 0.03 |  53.8 |      32 | 17946 |  8292 |  4382 |  2845 |  2128 |  1674 |  1352 |  1138 |   971 |   876 |   809 |   731 |   662 |   598 |   542 |   488 |   423 |   379 |   353 |   329 |   299 |   264 |   229 |   188 |   143 |   103 |    67 |    34 |    19 |    10 |     3 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      10 | 34327 |   21823 | 63.6 |     7943 | 23.1 |  0.0 | 0.01 |  63.5 |      11 | 12504 |  4561 |  1568 |   535 |   227 |   113 |    68 |    36 |    14 |     3 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      11 | 26061 |   17772 | 68.2 |     5818 | 22.3 |  0.0 | 0.00 |  70.2 |       7 |  8289 |  2471 |   671 |   152 |    36 |     7 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      12 | 18679 |   13275 | 71.1 |     3885 | 20.8 |  0.0 | 0.00 |  71.9 |       7 |  5404 |  1519 |   374 |    83 |    32 |     7 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      13 | 12541 |    9135 | 72.8 |     2655 | 21.2 |  0.0 | 0.01 |  78.0 |       5 |  3406 |   751 |   152 |     7 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      14 |  7908 |    5821 | 73.6 |     1577 | 19.9 |  0.0 | 0.01 |  75.6 |       6 |  2087 |   510 |   125 |     9 |     2 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      15 |  4773 |    3526 | 73.9 |     1002 | 21.0 |  0.0 | 0.00 |  80.4 |       4 |  1247 |   245 |    31 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |

Observation notes:
1. word - islands = ladder length 2 words
2. word - islands - doublets = ladder length 3 words


![Chart](analysis.png)

### Adjacent Words Counts

This table shows the spread of adjacent word counts for each word in the dictionary.
* Words are considered adjacent if changing just one letter in one word forms the other word.

| Letters |     0 |     1 |     2 |     3 |     4 |     5 |     6 |     7 |     8 |     9 |    10 |    11 |    12 |    13 |    14 |    15 |    16 |    17 |    18 |    19 |    20 |    21 |    22 |    23 |    24 |    25 |    26 |    27 |    28 |    29 |    30 |    31 |    32 |    33 |    34 |    35 |    36 |    37 |    38 |    39 |
|--------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|------:|
|       2 |       |       |       |       |       |     1 |     1 |     3 |     6 |     6 |     3 |     5 |     3 |     2 |     5 |     6 |    17 |    19 |    14 |    12 |     9 |     5 |     2 |       |       |       |     2 |       |     1 |     3 |     1 |       |     1 |       |       |       |       |       |       |       |
|       3 |     1 |     5 |    11 |    15 |    24 |    22 |    39 |    25 |    33 |    33 |    32 |    28 |    37 |    29 |    44 |    52 |    37 |    42 |    57 |    52 |    50 |    60 |    59 |    55 |    48 |    56 |    55 |    49 |    36 |    40 |    47 |    26 |    35 |    24 |    14 |     6 |     2 |     3 |     2 |       |
|       4 |    65 |   137 |   167 |   205 |   215 |   256 |   280 |   258 |   236 |   260 |   227 |   246 |   263 |   222 |   202 |   184 |   205 |   181 |   162 |   183 |   131 |   133 |   124 |    98 |   108 |    75 |    53 |    66 |    60 |    37 |    51 |    23 |    25 |    17 |    15 |     7 |     9 |     3 |     1 |     1 |
|       5 |   861 |  1167 |  1223 |  1119 |   990 |   897 |   765 |   634 |   613 |   483 |   473 |   401 |   323 |   294 |   274 |   237 |   176 |   174 |   140 |   108 |    90 |    76 |    70 |    58 |    40 |    33 |    23 |    17 |     7 |     9 |     9 |     1 |     6 |     3 |     2 |       |       |       |     1 |     1 |
|       6 |  3925 |  4051 |  3143 |  2219 |  1570 |  1173 |   893 |   679 |   565 |   478 |   415 |   373 |   320 |   279 |   225 |   184 |   142 |   109 |    84 |    52 |    28 |    24 |    15 |     5 |     1 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       7 |  9605 |  7522 |  4554 |  2697 |  1689 |  1242 |   929 |   785 |   583 |   446 |   357 |   269 |   232 |   161 |   124 |    97 |    74 |    59 |    33 |    23 |    11 |     3 |     1 |     2 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       8 | 16959 | 10529 |  5311 |  2778 |  1468 |   860 |   481 |   309 |   180 |   118 |    60 |    27 |    11 |     2 |     1 |     1 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|       9 | 22021 | 10984 |  4403 |  1697 |   547 |   185 |    68 |    29 |    21 |     6 |     1 |     4 |       |     1 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      10 | 21823 |  8503 |  2793 |   946 |   193 |    49 |    15 |     3 |     1 |     1 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      11 | 17772 |  6051 |  1663 |   468 |    76 |    25 |     4 |     2 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      12 | 13275 |  3981 |  1047 |   313 |    55 |     6 |     2 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      13 |  9135 |  2707 |   548 |   126 |    22 |     3 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      14 |  5821 |  1600 |   392 |    73 |    19 |     3 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |
|      15 |  3526 |  1026 |   194 |    22 |     5 |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |       |

![Chart](adjacents.png)

### Longest Ladders


#### 2-letter words

* _**33**_ words at maximum ladder length _**5**_
* `BA`, `BO`, `DA`, `DO`, `EL`, `ET`, `FA`, `GO`, `GU`, `IT`, `JA`, `JO`, `KA`, `KO`, `LA`, `LO`, `NA`, `NO`, `NU`, `PA`, `PO`, `TA`, `TO`, `UG`, `UP`, `UT`, `WO`, `XU`, `YA`, `YO`, `YU`, `ZA`, `ZO`


#### 3-letter words

* _**18**_ words at maximum ladder length _**9**_
* `ECO`, `GNU`, `ISM`, `IVY`, `JAB`, `JAW`, `JIB`, `KEB`, `KEF`, `KEX`, `KEY`, `MIB`, `MIZ`, `TWO`, `UEY`, `UFO`, `VAU`, `ZZZ`


#### 4-letter words

* _**5**_ words at maximum ladder length _**17**_
* `ODIC`, `OTIC`, `QUEY`, `UNAU`, `URIC`


#### 5-letter words

* _**2**_ words at maximum ladder length _**30**_
* `APPEL`, `INPUT`


#### 6-letter words

* _**3**_ words at maximum ladder length _**50**_
* `FOREST`, `MEREST`, `UNFEED`


#### 7-letter words

* _**4**_ words at maximum ladder length _**53**_
* `BARMILY`, `DESCEND`, `INJECTS`, `SHIVVED`


#### 8-letter words

* _**4**_ words at maximum ladder length _**61**_
* `ENDITING`, `ENROBING`, `ROWELING`, `TOWNLING`


#### 9-letter words

* _**3**_ words at maximum ladder length _**32**_
* `FLOORINGS`, `SLEEKINGS`, `SLEEVINGS`


#### 10-letter words

* _**3**_ words at maximum ladder length _**11**_
* `BLISTERING`, `GLISTENING`, `SNOTTERING`


#### 11-letter words

* _**7**_ words at maximum ladder length _**7**_
* `NOSOLOGISTS`, `OENOLOGISTS`, `OVERGANGING`, `OVERLOADING`, `RECOMMENDER`, `RECOMPENSED`, `RECOMPENSER`


#### 12-letter words

* _**7**_ words at maximum ladder length _**7**_
* `MATERIALIZED`, `MATERIALIZER`, `METROLOGISTS`, `PATERNALISTS`, `PATHOLOGIZED`, `ROOFLESSNESS`, `WORKLESSNESS`


#### 13-letter words

* _**7**_ words at maximum ladder length _**5**_
* `BASIFICATIONS`, `EXPANDABILITY`, `EXTENSIBILITY`, `FACTIONALISMS`, `FACTIONALISTS`, `FICTIONALIZED`, `RATIFICATIONS`


#### 14-letter words

* _**2**_ words at maximum ladder length _**6**_
* `DEHYDROGENATED`, `DEHYDROGENIZED`


#### 15-letter words

* _**31**_ words at maximum ladder length _**4**_
* `DECOLORISATIONS`, `DECOLORIZATIONS`, `HYPERCRITICISMS`, `HYPERCRITICIZED`, `LATERALISATIONS`, `LATERALIZATIONS`, `LIBERALISATIONS`, `LIBERALIZATIONS`, `PROCESSIONALIST`, `PROFESSIONALIZE`, `PROLETARIANISMS`, `PROLETARIANIZED`, `RECOLONISATIONS`, `RECOLONIZATIONS`, `SENSATIONALISMS`, `SENSATIONALISTS`, `SENSATIONALIZED`, `SENTIMENTALISMS`, `SENTIMENTALISTS`, `SENTIMENTALIZED`, `SUBSTANTIALISMS`, `SUBSTANTIALISTS`, `SUBSTANTIALIZED`, `TERRITORIALISMS`, `TERRITORIALISTS`, `TERRITORIALIZED`, `TRADITIONALISMS`, `TRADITIONALISTS`, `TRADITIONALIZED`, `UTILITARIANISMS`, `UTILITARIANIZED`

