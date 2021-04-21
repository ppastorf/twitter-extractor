package main

import (
	"log"
	"strconv"
	"worker-extractor/extractor"
	batch "worker-extractor/extractor/tweet_batch"
	"worker-extractor/producer"
)

type Credentials struct {
	Twitter *extractor.TwitterCredentials `yaml:"twitter"`
}

type Configs struct {
	Extractor *extractor.ExtractorConfig `yaml:"extractor"`
	Kafka     *producer.KafkaConfig      `yaml:"kafka"`
}

func main() {
	// Configs
	creds := &Credentials{}
	confs := &Configs{}

	twitterSecretFile := readEnvVariable("TWITTER_CREDENTIALS", "secret/twitter.yaml")
	log.Println(twitterSecretFile)
	err := readYamlFile(twitterSecretFile, creds)
	if err != nil {
		log.Fatalf("failed to open Twitter credentials file: %v\n", err)
		panic(err)
	}
	log.Printf("%+v", creds)

	extractorConfigFile := readEnvVariable("EXTRACTOR_CONFIG", "config/extractor.yaml")
	// log.Println(twitterSecretFile)
	err = readYamlFile(extractorConfigFile, confs)
	if err != nil {
		log.Fatalf("failed to open Extractor config file: %v\n", err)
		panic(err)
	}

	kafkaConfigFile := readEnvVariable("KAFKA_CONFIG", "config/kafka.yaml")
	// log.Println(twitterSecretFile)
	err = readYamlFile(kafkaConfigFile, confs)
	if err != nil {
		log.Fatalf("failed to open Kafka config file: %v\n", err)
		panic(err)
	}

	confs.Kafka.Topics.Workload.Name = readEnvVariable("WORKLOAD_TOPIC", confs.Kafka.Topics.Workload.Name)
	confs.Kafka.Topics.Workload.SchemaFile = readEnvVariable("WORKLOAD_SCHEMA", confs.Kafka.Topics.Workload.SchemaFile)
	confs.Kafka.EnableTls, _ = strconv.ParseBool(readEnvVariable("KAFKA_ENABLE_TLS", strconv.FormatBool(confs.Kafka.EnableTls)))
	confs.Kafka.EnableLogs, _ = strconv.ParseBool(readEnvVariable("SARAMA_LOGS", strconv.FormatBool(confs.Kafka.EnableLogs)))
	confs.Kafka.Address = readEnvVariable("KAFKA_ADDRESS", confs.Kafka.Address)
	confs.Extractor.Interval = readEnvVariable("EXTRACTOR_INTERVAL", confs.Extractor.Interval)

	// Init
	twitterExtractor, err := extractor.NewExtractor(creds.Twitter, confs.Extractor)
	if err != nil {
		log.Fatalf("failed to create new Twitter api extractor: %v\n", err)
		panic(err)
	}

	batchCodec, err := newCodecFromSchema(confs.Kafka.Topics.Workload.SchemaFile)
	if err != nil {
		log.Fatalf("failed to create new Avro Codec: %v\n", err)
		panic(err)
	}

	kafkaProducer, err := producer.NewProducer(confs.Kafka)
	if err != nil {
		log.Fatalf("failed to create new Kafka producer: %v\n", err)
	}
	go kafkaProducer.LogEvents()

	// Inicia extracao
	tweets := make(chan batch.TweetBatch)
	go twitterExtractor.Start(tweets)

	// Consome resultado das extracoes, encoda e produz para o Kafka
	for batch := range tweets {
		bytes, err := avroEncode(batch, batchCodec)
		if err != nil {
			log.Printf("failed to encode a tweet batch: %v\n", err)
		} else {
			go kafkaProducer.PublishAsync(confs.Kafka.Topics.Workload.Name, bytes)
		}
	}
}
