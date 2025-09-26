<p align="center">
  <img src="assets/spark_logo_dark_background.png" alt="Company Logo" width="300" height=100%/>
</p>

# Non Function Testing Document

**Version:** 1.0.0.0  
**Prepared By:** Spark Industries  
**Prepared For:** Demo4 Purposes 

## Content
* [Introduction](#introduction)
* [Performance Test](#performance-testing)

## Introduction 
The following document outlines some of the non functional testing we performed to ensure the quality requirements as outlined in our documentation.
It is broken up into specific non functional tests that we performed. It details why these tests were needed, what quality requirement it aims to test, how the tests were performed and the test results.

## Performance Testing

To ensure the adequate performance of our application we test 2 features specifically namely our clustering (semantic sorting) and keyword (used for smart search) endpoints. Both of the mentioned endpoints were tested due to them being the most performance intensive.

### Expected Results
We aim to ensure that both of these endpoints scales at least linearly in execution time for an increasing number of files. We cannot place a numerical value on exactly how long these tests should take, as it is heavily dependant on the hardware of the user using the service. We note the hardware of the machine used to run these tests in our experimental procedure.

### Experimental Procedure

To properly test these endpoints we conduct an experiment running both the clustering and keyword service for an increasing number of files. We make use of the same file types containing the same content (in this test a textfile). Note that different file types could slightly influence the results. To ensure accurate results we also run each test 3 times and take the average of the results to report on. 