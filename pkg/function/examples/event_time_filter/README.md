# Event Time Filter

This is a simple User Defined Function example which receives a message, applies the following data transformation, and returns the message.

## Data Transformation
* If the message event time is **before** year 2022, drop the message.
* If it's **within** year 2022, update the key to "within_year_2022" and update the message event time to Jan 1st 2022.
* Otherwise(exclusively after year 2022), update the key to "after_year_2022" and update the message event time to Jan 1st 2023.
