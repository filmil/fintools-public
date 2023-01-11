# Testing

Sadly, compensation data is private. This means, I can't really check my paystub
data here, and I can not reliably get a test-only paystub.

Instead, I decided to allow each contributor to initialize their own test data
repository.  The directory `testdata_private` is gitignored here, so you can
even add your own private repository there.

Instead, you can do this:

1. Create a subdirectory called `testdata_private` in this directory. 
2. Create a file `test_spec.json` in `testdata_private`. Sample contents are below.
3. Add the test data file `testdata_private/out.txt`.

Sample `test_spec.json` is below. All data is faked.

```
[
  {
    "Input": "testdata_private/out.txt",
    "Expected": {
      "Date": "2019-01-18",
      "DocNum": "13541270",
      "NetPay": 1000
      "AnnualBonus": 300,
      "Bonus401kPre": 400,
      "FederalIncomeTax": 400,
      "EmployeeMedicare": 500
      "SocialSecurityEmployeeTax": 600,
      "CAStateIncomeTax": 700.25,
      "CAPrivateDisabilityEmployee": 800,
      "Employer": {
        "Bonus401kPre": 900
      }
    }
  }
]
```
