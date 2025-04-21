from pyspark.sql import SparkSession

if __name__ == "__main__":  
  spark = SparkSession.builder.config(
      "spark.jars.packages",
      "com.google.cloud.spark:spark-bigquery-with-dependencies_2.12:0.42.1",
  ).getOrCreate()

  spark.conf.set("credentialsFile", "./credentials.json")

  # Considering we already have testschema dataset and sparttesttable table!!!
  df = spark.read.format("bigquery").load("testschema.sparktesttable")
  df.show()