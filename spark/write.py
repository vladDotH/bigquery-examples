from pyspark.sql import SparkSession

if __name__ == "__main__":
  spark = SparkSession.builder.config(
      "spark.jars.packages",
      "com.google.cloud.spark:spark-bigquery-with-dependencies_2.12:0.42.1",
  ).getOrCreate()

  spark.conf.set("credentialsFile", "./credentials.json")

  df = spark.createDataFrame(
      [
        (1234, "some-str"), 
        (0, "null"), 
        (10000, "1000-string")
        ], 
      ["num", "str"]
  )

  # Creates new table (for appending use `.mode('append')`)
  df.write.format("bigquery").option("writeMethod", "direct").save(
      "testschema.sparktesttable"
  )
