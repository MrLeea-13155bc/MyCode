package main

func init() {
	addresses := []string{"http://127.0.0.1:9200", "http://127.0.0.1:9201"}
	config := elasticsearch.Config{
		Addresses: addresses,
		Username:  "",
		Password:  "",
		CloudID:   "",
		APIKey:    "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	failOnError(err, "Error creating the client")

}

func main() {

}
